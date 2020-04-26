#!/usr/bin/env python
# vim: set encoding=utf-8

"""
Main wttr.in rendering function implementation
"""

import logging
import io
import os
import time
from flask import render_template, send_file, make_response

import fmt.png

import parse_query
from translations import get_message, FULL_TRANSLATION, PARTIAL_TRANSLATION, SUPPORTED_LANGS
from buttons import add_buttons
from globals import get_help_file, \
                    BASH_FUNCTION_FILE, TRANSLATION_FILE, LOG_FILE, \
                    NOT_FOUND_LOCATION, \
                    MALFORMED_RESPONSE_HTML_PAGE, \
                    PLAIN_TEXT_AGENTS, PLAIN_TEXT_PAGES, \
                    MY_EXTERNAL_IP, QUERY_LIMITS
from location import is_location_blocked, location_processing
from limits import Limits
from view.wttr import get_wetter
from view.moon import get_moon
from view.line import wttr_line

import cache

if not os.path.exists(os.path.dirname(LOG_FILE)):
    os.makedirs(os.path.dirname(LOG_FILE))
logging.basicConfig(filename=LOG_FILE, level=logging.INFO, format='%(asctime)s %(message)s')

LIMITS = Limits(whitelist=[MY_EXTERNAL_IP], limits=QUERY_LIMITS)

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
    return text

def _client_ip_address(request):
    """Return client ip address for flask `request`.
    """

    if request.headers.getlist("X-PNG-Query-For"):
        ip_addr = request.headers.getlist("X-PNG-Query-For")[0]
        if ip_addr.startswith('::ffff:'):
            ip_addr = ip_addr[7:]
    elif request.headers.getlist("X-Forwarded-For"):
        ip_addr = request.headers.getlist("X-Forwarded-For")[0]
        if ip_addr.startswith('::ffff:'):
            ip_addr = ip_addr[7:]
    else:
        ip_addr = request.remote_addr

    return ip_addr

def _parse_language_header(header):
    """
    >>> _parse_language_header("en-US,en;q=0.9")
    >>> _parse_language_header("en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
    >>> _parse_language_header("xx, fr-CA;q=0.8, da-DK;q=0.9")
    'da'
    """

    def _parse_accept_language(accept_language):
        languages = accept_language.split(",")
        locale_q_pairs = []

        for language in languages:
            try:
                if language.split(";")[0] == language:
                    # no q => q = 1
                    locale_q_pairs.append((language.strip(), 1))
                else:
                    locale = language.split(";")[0].strip()
                    weight = float(language.split(";")[1].split("=")[1])
                    locale_q_pairs.append((locale, weight))
            except (IndexError, ValueError):
                pass

        return locale_q_pairs

    def _find_supported_language(accepted_languages):

        def supported_langs():
            """Yields all pairs in the Accept-Language header
            supported in SUPPORTED_LANGS or None if 'en' is the preferred"""
            for lang_tuple in accepted_languages:
                lang = lang_tuple[0]
                if '-' in lang:
                    lang = lang.split('-', 1)[0]
                if lang in SUPPORTED_LANGS:
                    yield lang, lang_tuple[1]
                elif lang == 'en':
                    yield None, lang_tuple[1]
        try:
            return max(supported_langs(), key=lambda lang_tuple: lang_tuple[1])[0]
        except ValueError:
            return None

    return _find_supported_language(_parse_accept_language(header))

def get_answer_language_and_view(request):
    """
    Return preferred answer language based on
    domain name, query arguments and headers
    """

    lang = None
    view_name = None
    hostname = request.headers['Host']
    if hostname != 'wttr.in' and hostname.endswith('.wttr.in'):
        lang = hostname[:-8]
        if lang.startswith("v2"):
            view_name = lang
            lang = None

    if 'lang' in request.args:
        lang = request.args.get('lang')
        if lang.lower() == 'none':
            lang = None

    header_accept_language = request.headers.get('Accept-Language', '')
    if lang is None and header_accept_language:
        lang = _parse_language_header(header_accept_language)

    return lang, view_name

def get_output_format(query, parsed_query):
    """
    Return preferred output format: ansi, text, html or png
    based on arguments and headers in `request`.
    Return new location (can be rewritten)
    """

    if ('view' in query and not query["view"].startswith("v2")) \
        or parsed_query.get("png_filename") \
        or query.get('force-ansi'):
        return False

    user_agent = parsed_query.get("user_agent", "").lower()
    html_output = not any(agent in user_agent for agent in PLAIN_TEXT_AGENTS)
    return html_output

def _cyclic_location_selection(locations, period):
    """Return one of `locations` (: separated list)
    basing on the current time and query interval `period`
    """

    locations = locations.split(':')
    max_len = max(len(x) for x in locations)
    locations = [x.rjust(max_len) for x in locations]

    try:
        period = int(period)
    except ValueError:
        period = 1

    index = int(time.time()/period) % len(locations)
    return locations[index]


def _response(parsed_query, query, fast_mode=False):
    """Create response text based on `parsed_query` and `query` data.
    If `fast_mode` is True, process only requests that can
    be handled very fast (cached and static files).
    """

    answer = None
    cache_signature = cache.get_signature(
        parsed_query["user_agent"],
        parsed_query["request_url"],
        parsed_query["ip_addr"],
        parsed_query["lang"])
    answer = cache.get(cache_signature)

    if parsed_query['orig_location'] in PLAIN_TEXT_PAGES:
        answer = show_text_file(parsed_query['orig_location'], parsed_query['lang'])
        if parsed_query['html_output']:
            answer = render_template('index.html', body=answer)

    if answer or fast_mode:
        return answer

    # at this point, we could not handle the query fast,
    # so we handle it with all available logic
    loc = (parsed_query['orig_location'] or "").lower()
    if parsed_query.get("view"):
        output = wttr_line(query, parsed_query)
    elif loc == 'moon' or loc.startswith('moon@'):
        output = get_moon(query, parsed_query)
    else:
        output = get_wetter(query, parsed_query)

    if parsed_query.get('png_filename'):
        output = fmt.png.render_ansi(
            output, options=parsed_query)
    else:
        if query.get('days', '3') != '0' and not query.get('no-follow-line'):
            if parsed_query['html_output']:
                output = add_buttons(output)
            else:
                output += '\n' + get_message('FOLLOW_ME', parsed_query['lang']) + '\n'

    return cache.store(cache_signature, output)

def parse_request(location, request, query, fast_mode=False):
    """Parse request and provided extended information for the query,
    including location data, language, output format, view, etc.

    Incoming data:

        `location`              location name extracted from the query url
        `request.args`
        `request.headers`
        `request.remote_addr`
        `request.referrer`
        `request.query_string`
        `query`                 parsed command line arguments

    Parameters priorities (from low to high):

        * HTTP-header
        * Domain name
        * URL
        * Filename

    Return: dictionary with parsed parameters
    """

    png_filename = None
    if location is not None and location.lower().endswith(".png"):
        png_filename = location
        location = location[:-4]
    if location and ':' in location and location[0] != ":":
        location = _cyclic_location_selection(location, query.get('period', 1))

    parsed_query = {
        'ip_addr': _client_ip_address(request),
        'user_agent': request.headers.get('User-Agent', '').lower(),
        'request_url': request.url,
        }

    if png_filename:
        parsed_query["png_filename"] = png_filename
        parsed_query.update(parse_query.parse_wttrin_png_name(png_filename))

    lang, _view = get_answer_language_and_view(request)

    parsed_query["view"] = parsed_query.get("view", query.get("view", _view))
    parsed_query["location"] = parsed_query.get("location", location)
    parsed_query["orig_location"] = parsed_query["location"]
    parsed_query["lang"] = parsed_query.get("lang", lang)

    parsed_query["html_output"] = get_output_format(query, parsed_query)

    if not fast_mode: # not png_filename and not fast_mode:
        location, override_location_name, full_address, country, query_source_location = \
                location_processing(parsed_query["location"], parsed_query["ip_addr"])

        us_ip = query_source_location[1] == 'United States' \
                and 'slack' not in parsed_query['user_agent']
        query = parse_query.metric_or_imperial(query, lang, us_ip=us_ip)

        if country and location != NOT_FOUND_LOCATION:
            location = "%s,%s" % (location, country)

        parsed_query.update({
            'location': location,
            'override_location_name': override_location_name,
            'full_address': full_address,
            'country': country,
            'query_source_location': query_source_location})

    return parsed_query


def wttr(location, request):
    """Main rendering function, it processes incoming weather queries,
    and depending on the User-Agent string and other paramters of the query
    it returns output in HTML, ANSI or other format.
    """

    def _wrap_response(response_text, html_output, png_filename=None):
        if not isinstance(response_text, str) and \
           not isinstance(response_text, bytes):
            return response_text

        if png_filename:
            response = make_response(send_file(
                io.BytesIO(response_text),
                attachment_filename=png_filename,
                mimetype='image/png'))

            for key, value in {
                    'Cache-Control': 'no-cache, no-store, must-revalidate',
                    'Pragma': 'no-cache',
                    'Expires': '0',
                }.items():
                response.headers[key] = value
        else:
            response = make_response(response_text)
            response.mimetype = 'text/html' if html_output else 'text/plain'
        return response

    if is_location_blocked(location):
        return ""

    try:
        LIMITS.check_ip(_client_ip_address(request))
    except RuntimeError as exception:
        return str(exception)

    query = parse_query.parse_query(request.args)

    # first, we try to process the query as fast as possible
    # (using the cache and static files),
    # and only if "fast_mode" was unsuccessful,
    # use the full track
    parsed_query = parse_request(location, request, query, fast_mode=True)
    response = _response(parsed_query, query, fast_mode=True)

    try:
        if not response:
            parsed_query = parse_request(location, request, query)
            response = _response(parsed_query, query)
    # pylint: disable=broad-except
    except Exception as exception:
        logging.error("Exception has occured", exc_info=1)
        if parsed_query['html_output']:
            response = MALFORMED_RESPONSE_HTML_PAGE
        else:
            response = get_message('CAPACITY_LIMIT_REACHED', parsed_query['lang'])

        # if exception is occured, we return not a png file but text
        if "png_filename" in parsed_query:
            del parsed_query["png_filename"]
    return _wrap_response(
        response, parsed_query['html_output'],
        png_filename=parsed_query.get('png_filename'))

if __name__ == "__main__":
    import doctest
    doctest.testmod()
