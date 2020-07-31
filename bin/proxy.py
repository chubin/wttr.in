#vim: fileencoding=utf-8

"""

The proxy server acts as a backend for the wttr.in service.

It caches the answers and handles various data sources transforming their
answers into format supported by the wttr.in service.

If WTTRIN_TEST is specified, it works in a special test mode:
it does not fetch and does not store the data in the cache,
but is using the fake data from "test/proxy-data".

"""
from __future__ import print_function

from gevent.pywsgi import WSGIServer
from gevent.monkey import patch_all
patch_all()

# pylint: disable=wrong-import-position,wrong-import-order
import sys
import os
import time
import json
import hashlib

import requests
import cyrtranslit

from flask import Flask, request
APP = Flask(__name__)


MYDIR = os.path.abspath(
    os.path.dirname(os.path.dirname(__file__)))
sys.path.append("%s/lib/" % MYDIR)

from globals import PROXY_CACHEDIR, PROXY_HOST, PROXY_PORT
from translations import PROXY_LANGS
# pylint: enable=wrong-import-position

def is_testmode():
    """Server is running in the wttr.in test mode"""

    return "WTTRIN_TEST" in os.environ

def load_translations():
    """
    load all translations
    """
    translations = {}

    for f_name in PROXY_LANGS:
        f_name = '/root/app/share/translations/%s.txt' % f_name
        translation = {}
        lang = f_name.split('/')[-1].split('.', 1)[0]
        with open(f_name, "r") as f_file:
            for line in f_file:
                if ':' not in line:
                    continue
                if line.count(':') == 3:
                    _, trans, orig, _ = line.strip().split(':', 4)
                else:
                    _, trans, orig = line.strip().split(':', 3)
                trans = trans.strip()
                orig = orig.strip()

                translation[orig] = trans
        translations[lang] = translation
    return translations
TRANSLATIONS = load_translations()

def _find_srv_for_query(path, query):   # pylint: disable=unused-argument
    return 'http://api.worldweatheronline.com'

def _cache_file(path, query):
    """Return cache file name for specified `path` and `query`
    and for the current time.

    Do smooth load on the server, expiration time
    is slightly varied basing on the path+query sha1 hash digest.
    """

    digest = hashlib.sha1(("%s %s" % (path, query)).encode("utf-8")).hexdigest()
    digest_number = ord(digest[0].upper())
    expiry_interval = 60*(digest_number+30)

    timestamp = "%010d" % (int(time.time())//expiry_interval*expiry_interval)
    filename = os.path.join(PROXY_CACHEDIR, timestamp, path, query)

    return filename


def _load_content_and_headers(path, query):
    if is_testmode():
        cache_file = "test/proxy-data/data1"
    else:
        cache_file = _cache_file(path, query)
    try:
        return (open(cache_file, 'r').read(),
                json.loads(open(cache_file+".headers", 'r').read()))
    except IOError:
        return None, None

def _touch_empty_file(path, query):
    cache_file = _cache_file(path, query)
    cache_dir = os.path.dirname(cache_file)
    if not os.path.exists(cache_dir):
        os.makedirs(cache_dir)
    open(cache_file, 'w').write("")

def _save_content_and_headers(path, query, content, headers):
    cache_file = _cache_file(path, query)
    cache_dir = os.path.dirname(cache_file)
    if not os.path.exists(cache_dir):
        os.makedirs(cache_dir)
    open(cache_file + ".headers", 'w').write(json.dumps(headers))
    open(cache_file, 'wb').write(content)

def translate(text, lang):
    """
    Translate `text` into `lang`
    """
    translated = TRANSLATIONS.get(lang, {}).get(text, text)
    if text == translated:
        print("%s: %s" % (lang, text))
    return translated

def cyr(to_translate):
    """
    Transliterate `to_translate` from latin into cyrillic
    """
    return cyrtranslit.to_cyrillic(to_translate)

def _patch_greek(original):
    return original.replace(u"Ηλιόλουστη/ο", u"Ηλιόλουστη")

def add_translations(content, lang):
    """
    Add `lang` translation to `content` (JSON)
    returned by the data source
    """

    if content == "{}":
        return {}

    languages_to_translate = TRANSLATIONS.keys()
    try:
        d = json.loads(content)         # pylint: disable=invalid-name
    except (ValueError, TypeError) as exception:
        print("---")
        print(exception)
        print("---")
        return {}

    try:
        weather_condition = d['data']['current_condition'
                ][0]['weatherDesc'][0]['value'].capitalize()
        d['data']['current_condition'][0]['weatherDesc'][0]['value'] = \
            weather_condition
        if lang in languages_to_translate:
            d['data']['current_condition'][0]['lang_%s' % lang] = \
                [{'value': translate(weather_condition, lang)}]
        elif lang == 'sr':
            d['data']['current_condition'][0]['lang_%s' % lang] = \
                [{'value': cyr(
                    d['data']['current_condition'][0]['lang_%s' % lang][0]['value']\
                    )}]
        elif lang == 'el':
            d['data']['current_condition'][0]['lang_%s' % lang] = \
                [{'value': _patch_greek(
                    d['data']['current_condition'][0]['lang_%s' % lang][0]['value']\
                    )}]
        elif lang == 'sr-lat':
            d['data']['current_condition'][0]['lang_%s' % lang] = \
                [{'value':d['data']['current_condition'][0]['lang_sr'][0]['value']\
                            }]

        fixed_weather = []
        for w in d['data']['weather']:  # pylint: disable=invalid-name
            fixed_hourly = []
            for h in w['hourly']:       # pylint: disable=invalid-name
                weather_condition = h['weatherDesc'][0]['value']
                if lang in languages_to_translate:
                    h['lang_%s' % lang] = \
                        [{'value': translate(weather_condition, lang)}]
                elif lang == 'sr':
                    h['lang_%s' % lang] = \
                        [{'value': cyr(h['lang_%s' % lang][0]['value'])}]
                elif lang == 'el':
                    h['lang_%s' % lang] = \
                        [{'value': _patch_greek(h['lang_%s' % lang][0]['value'])}]
                elif lang == 'sr-lat':
                    h['lang_%s' % lang] = \
                        [{'value': h['lang_sr'][0]['value']}]
                fixed_hourly.append(h)
            w['hourly'] = fixed_hourly
            fixed_weather.append(w)
        d['data']['weather'] = fixed_weather

        content = json.dumps(d)
    except (IndexError, ValueError) as exception:
        print(exception)
    return content

@APP.route("/<path:path>")
def proxy(path):
    """
    Main proxy function. Handles incoming HTTP queries.
    """

    lang = request.args.get('lang', 'en')
    query_string = request.query_string.decode("utf-8")
    query_string = query_string.replace('sr-lat', 'sr')
    query_string = query_string.replace('lang=None', 'lang=en')
    query_string += "&extra=localObsTime"
    query_string += "&includelocation=yes"
    content, headers = _load_content_and_headers(path, query_string)

    if content is None:
        srv = _find_srv_for_query(path, query_string)
        url = '%s/%s?%s' % (srv, path, query_string)

        attempts = 10
        response = None
        print(url)
        while attempts:
            try:
                response = requests.get(url, timeout=2)
            except requests.ReadTimeout:
                attempts -= 1
                continue
            try:
                json.loads(response.content)
                break
            except ValueError:
                attempts -= 1

        _touch_empty_file(path, query_string)
        if response:
            headers = {}
            headers['Content-Type'] = response.headers['content-type']
            _save_content_and_headers(path, query_string, response.content, headers)
            content = response.content
        else:
            content = "{}"
    else:
        print("cache found")

    content = add_translations(content, lang)

    return content, 200, headers

if __name__ == "__main__":
    #app.run(host='0.0.0.0', port=5001, debug=False)
    #app.debug = True
    bind_addr = "0.0.0.0"
    SERVER = WSGIServer((bind_addr, PROXY_PORT), APP)
    SERVER.serve_forever()
