#!/usr/bin/env python
# vim: set encoding=utf-8

"""
Main wttr.in rendering function implementation
"""

import logging
import os
import re
import socket
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

if not os.path.exists(os.path.dirname(LOG_FILE)):
    os.makedirs(os.path.dirname(LOG_FILE))
logging.basicConfig(filename=LOG_FILE, level=logging.DEBUG, format='%(asctime)s %(message)s')

LIMITS = Limits(whitelist=[MY_EXTERNAL_IP], limits=(30, 60, 100))

def is_ip(ip_addr):
    """
    Check if `ip_addr` looks like an IP Address
    """

    if re.match(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', ip_addr) is None:
        return False
    try:
        socket.inet_aton(ip_addr)
        return True
    except socket.error:
        return False

def location_normalize(location):
    """
    Normalize location name `location`
    """
    #translation_table = dict.fromkeys(map(ord, '!@#$*;'), None)
    def _remove_chars(chars, string):
        return ''.join(x for x in string if x not in chars)

    location = location.lower().replace('_', ' ').replace('+', ' ').strip()
    if not location.startswith('moon@'):
        location = _remove_chars(r'!@#$*;:\\', location)
    return location

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
    if request.referrer:
        print request.referrer

    hostname = request.headers['Host']
    lang = None
    if hostname != 'wttr.in' and hostname.endswith('.wttr.in'):
        lang = hostname[:-8]

    if request.headers.getlist("X-Forwarded-For"):
        ip_addr = request.headers.getlist("X-Forwarded-For")[0]
        if ip_addr.startswith('::ffff:'):
            ip_addr = ip_addr[7:]
    else:
        ip_addr = request.remote_addr

    try:
        LIMITS.check_ip(ip_addr)
    except RuntimeError, e:
        return str(e)
    except Exception, e:
        logging.error("Exception has occured", exc_info=1)
        return "ERROR"

    if location is not None and location.lower() in LOCATION_BLACK_LIST:
        return ""

    png_filename = None
    if location is not None and location.lower().endswith(".png"):
        png_filename = location
        location = location[:-4]

    query = parse_query.parse_query(request.args)

    if 'lang' in request.args:
        lang = request.args.get('lang')
    if lang is None and 'Accept-Language' in request.headers:
        lang = find_supported_language(
            parse_accept_language(
                request.headers.get('Accept-Language', '')))

    user_agent = request.headers.get('User-Agent', '').lower()
    html_output = not any(agent in user_agent for agent in PLAIN_TEXT_AGENTS)
    if location in PLAIN_TEXT_PAGES:
        help_ = show_help(location, lang)
        if html_output:
            return render_template('index.html', body=help_)
        return help_

    orig_location = location

    if request.headers.getlist("X-Forwarded-For"):
        ip_addr = request.headers.getlist("X-Forwarded-For")[0]
        if ip_addr.startswith('::ffff:'):
            ip_addr = ip_addr[7:]
    else:
        ip_addr = request.remote_addr

    try:
        # if location is starting with ~
        # or has non ascii symbols
        # it should be handled like a search term (for geolocator)
        override_location_name = None
        full_address = None

        if location is not None and not ascii_only(location):
            location = "~" + location

        if location is not None and location.upper() in IATA_CODES:
            location = '~%s' % location

        if location is not None and location.startswith('~'):
            geolocation = geolocator(location_canonical_name(location[1:]))
            if geolocation is not None:
                override_location_name = location[1:].replace('+', ' ')
                location = "%s,%s" % (geolocation['latitude'], geolocation['longitude'])
                full_address = geolocation['address']
                print full_address
            else:
                location = NOT_FOUND_LOCATION #location[1:]
        try:
            query_source_location = get_location(ip_addr)
        except:
            query_source_location = NOT_FOUND_LOCATION, None

        # what units should be used
        # metric or imperial
        # based on query and location source (imperial for US by default)
        print "lang = %s" % lang
        if query.get('use_metric', False) and not query.get('use_imperial', False):
            query['use_imperial'] = False
            query['use_metric'] = True
        elif query.get('use_imperial', False) and not query.get('use_metric', False):
            query['use_imperial'] = True
            query['use_metric'] = False
        elif lang == 'us':
            # slack uses m by default, to override it speciy us.wttr.in
            query['use_imperial'] = True
            query['use_metric'] = False
        else:
            if query_source_location[1] in ['US'] and 'slack' not in user_agent:
                query['use_imperial'] = True
                query['use_metric'] = False
            else:
                query['use_imperial'] = False
                query['use_metric'] = True

        country = None
        if location is None or location == 'MyLocation':
            location, country = query_source_location

        if is_ip(location):
            location, country = get_location(location)
        if location.startswith('@'):
            try:
                location, country = get_location(socket.gethostbyname(location[1:]))
            except:
                query_source_location = NOT_FOUND_LOCATION, None

        location = location_canonical_name(location)
        log("%s %s %s %s %s %s" \
            % (ip_addr, user_agent, orig_location, location,
               query.get('use_imperial', False), lang))

        # We are ready to return the answer
        if png_filename:
            options = {}
            if lang is not None:
                options['lang'] = lang

            options['location'] = "%s,%s" % (location, country)
            options.update(query)

            cached_png_file = wttrin_png.make_wttr_in_png(png_filename, options=options)
            response = make_response(send_file(cached_png_file,
                                               attachment_filename=png_filename,
                                               mimetype='image/png'))

            # Trying to disable github caching
            response.headers['Cache-Control'] = 'no-cache, no-store, must-revalidate'
            response.headers['Pragma'] = 'no-cache'
            response.headers['Expires'] = '0'
            return response

        if location == 'moon' or location.startswith('moon@'):
            output = get_moon(location, html=html_output, lang=lang)
        else:
            if country and location != NOT_FOUND_LOCATION:
                location = "%s, %s" % (location, country)
            output = get_wetter(location, ip_addr,
                                html=html_output,
                                lang=lang,
                                query=query,
                                location_name=override_location_name,
                                full_address=full_address,
                                url=request.url,
                               )

        if 'Malformed response' in str(output) \
                or 'API key has reached calls per day allowed limit' in str(output):
            if html_output:
                return MALFORMED_RESPONSE_HTML_PAGE
            return get_message('CAPACITY_LIMIT_REACHED', lang).encode('utf-8')

        if html_output:
            output = output.replace('</body>',
                                    (TWITTER_BUTTON
                                     + GITHUB_BUTTON
                                     + GITHUB_BUTTON_3
                                     + GITHUB_BUTTON_2
                                     + GITHUB_BUTTON_FOOTER) + '</body>')
        else:
            if query.get('days', '3') != '0':
                #output += '\n' + get_message('NEW_FEATURE', lang).encode('utf-8')
                output += '\n' + get_message('FOLLOW_ME', lang).encode('utf-8') + '\n'
        return output

    #except RuntimeError, e:
    #    return str(e)
    except Exception, e:
        if 'Malformed response' in str(e) \
                or 'API key has reached calls per day allowed limit' in str(e):
            if html_output:
                return MALFORMED_RESPONSE_HTML_PAGE
            return get_message('CAPACITY_LIMIT_REACHED', lang).encode('utf-8')
        logging.error("Exception has occured", exc_info=1)
        return "ERROR"

SERVER = WSGIServer((LISTEN_HOST, LISTEN_PORT), APP)
SERVER.serve_forever()
