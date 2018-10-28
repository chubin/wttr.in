#!/usr/bin/env python
# vim: set encoding=utf-8

"""
Main wttr.in rendering function implementation
"""

import logging
import os
from flask import render_template, send_file, make_response

import wttrin_png
import parse_query
from translations import get_message, FULL_TRANSLATION, PARTIAL_TRANSLATION, SUPPORTED_LANGS
from buttons import add_buttons
from globals import get_help_file, log, \
                    BASH_FUNCTION_FILE, TRANSLATION_FILE, LOG_FILE, \
                    NOT_FOUND_LOCATION, \
                    MALFORMED_RESPONSE_HTML_PAGE, \
                    PLAIN_TEXT_AGENTS, PLAIN_TEXT_PAGES, \
                    MY_EXTERNAL_IP
from location import is_location_blocked, location_processing
from limits import Limits
from wttr import get_wetter, get_moon
from wttr_line import wttr_line

if not os.path.exists(os.path.dirname(LOG_FILE)):
    os.makedirs(os.path.dirname(LOG_FILE))
logging.basicConfig(filename=LOG_FILE, level=logging.DEBUG, format='%(asctime)s %(message)s')

LIMITS = Limits(whitelist=[MY_EXTERNAL_IP], limits=(30, 60, 100))

def show_text_file(name, lang):
    """
    show static file `name` for `lang`
    """
    text = ""
    if name == ":help":
        text = open(get_help_file(lang), 'r').read()
        text = text.replace('FULL_TRANSLATION', ' '.join(FULL_TRANSLATION))
        text = text.replace('PARTIAL_TRANSLATION', ' '.join(PARTIAL_TRANSLATION))
    elif name == ":bash.function":
        text = open(BASH_FUNCTION_FILE, 'r').read()
    elif name == ":translation":
        text = open(TRANSLATION_FILE, 'r').read()
        text = text\
                .replace('NUMBER_OF_LANGUAGES', str(len(SUPPORTED_LANGS)))\
                .replace('SUPPORTED_LANGUAGES', ' '.join(SUPPORTED_LANGS))
    return text.decode('utf-8')

def client_ip_address(request):
    """
    Return client ip address for `request`.
    Flask related
    """

    if request.headers.getlist("X-Forwarded-For"):
        ip_addr = request.headers.getlist("X-Forwarded-For")[0]
        if ip_addr.startswith('::ffff:'):
            ip_addr = ip_addr[7:]
    else:
        ip_addr = request.remote_addr

    return ip_addr

def get_answer_language(request):
    """
    Return preferred answer language based on
    domain name, query arguments and headers
    """

    def _parse_accept_language(accept_language):
        languages = accept_language.split(",")
        locale_q_pairs = []

        for language in languages:
            try:
                if language.split(";")[0] == language:
                    # no q => q = 1
                    locale_q_pairs.append((language.strip(), "1"))
                else:
                    locale = language.split(";")[0].strip()
                    weight = language.split(";")[1].split("=")[1]
                    locale_q_pairs.append((locale, weight))
            except IndexError:
                pass

        return locale_q_pairs

    def _find_supported_language(accepted_languages):
        for lang_tuple in accepted_languages:
            lang = lang_tuple[0]
            if '-' in lang:
                lang = lang.split('-', 1)[0]
            if lang in SUPPORTED_LANGS:
                return lang
        return None

    lang = None
    hostname = request.headers['Host']
    if hostname != 'wttr.in' and hostname.endswith('.wttr.in'):
        lang = hostname[:-8]

    if 'lang' in request.args:
        lang = request.args.get('lang')

    header_accept_language = request.headers.get('Accept-Language', '')
    if lang is None and header_accept_language:
        lang = _find_supported_language(
            _parse_accept_language(header_accept_language))

    return lang

def get_output_format(request):
    """
    Return preferred output format: ansi, text, html or png
    based on arguments and headers in `request`.
    Return new location (can be rewritten)
    """

    # FIXME
    user_agent = request.headers.get('User-Agent', '').lower()
    html_output = not any(agent in user_agent for agent in PLAIN_TEXT_AGENTS)
    return html_output


def wttr(location, request):
    """
    Main rendering function, it processes incoming weather queries.
    Depending on user agent it returns output in HTML or ANSI format.

    Incoming data:
        request.args
        request.headers
        request.remote_addr
        request.referrer
        request.query_string
    """

    if is_location_blocked(location):
        return ""

    ip_addr = client_ip_address(request)

    try:
        LIMITS.check_ip(ip_addr)
    except RuntimeError, exception:
        return str(exception)

    png_filename = None
    if location is not None and location.lower().endswith(".png"):
        png_filename = location
        location = location[:-4]

    lang = get_answer_language(request)
    query = parse_query.parse_query(request.args)
    html_output = get_output_format(request)
    user_agent = request.headers.get('User-Agent', '').lower()

    if location in PLAIN_TEXT_PAGES:
        help_ = show_text_file(location, lang)
        if html_output:
            return render_template('index.html', body=help_)
        return help_

    orig_location = location

    location, override_location_name, full_address, country, query_source_location = \
            location_processing(location, ip_addr)

    us_ip = query_source_location[1] == 'United States' and 'slack' not in user_agent
    query = parse_query.metric_or_imperial(query, lang, us_ip=us_ip)

    # logging query
    orig_location_utf8 = (orig_location or "").encode('utf-8')
    location_utf8 = location.encode('utf-8')
    use_imperial = query.get('use_imperial', False)
    log(" ".join(map(str,
                     [ip_addr, user_agent, orig_location_utf8, location_utf8, use_imperial, lang])))

    if country and location != NOT_FOUND_LOCATION:
        location = "%s,%s" % (location, country)

    # We are ready to return the answer
    try:
        if 'format' in query:
            return wttr_line(location, query)

        if png_filename:
            options = {
                'lang': None,
                'location': location}
            options.update(query)

            cached_png_file = wttrin_png.make_wttr_in_png(png_filename, options=options)
            response = make_response(send_file(cached_png_file,
                                               attachment_filename=png_filename,
                                               mimetype='image/png'))
            for key, value in {
                    'Cache-Control': 'no-cache, no-store, must-revalidate',
                    'Pragma': 'no-cache',
                    'Expires': '0',
                }.items():
                response.headers[key] = value

            # Trying to disable github caching
            return response

        if location == 'moon' or location.startswith('moon@'):
            output = get_moon(location, html=html_output, lang=lang)
        else:
            output = get_wetter(location, ip_addr,
                                html=html_output,
                                lang=lang,
                                query=query,
                                location_name=override_location_name,
                                full_address=full_address,
                                url=request.url,
                               )

        if html_output:
            output = add_buttons(output)
        else:
            if query.get('days', '3') != '0':
                #output += '\n' + get_message('NEW_FEATURE', lang).encode('utf-8')
                output += '\n' + get_message('FOLLOW_ME', lang).encode('utf-8') + '\n'
        return output

    except RuntimeError, exception:
        if 'Malformed response' in str(exception) \
                or 'API key has reached calls per day allowed limit' in str(exception):
            if html_output:
                return MALFORMED_RESPONSE_HTML_PAGE
            return get_message('CAPACITY_LIMIT_REACHED', lang).encode('utf-8')
        logging.error("Exception has occured", exc_info=1)
        return "ERROR"
