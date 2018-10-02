#vim: encoding=utf-8

import gevent
from gevent.wsgi import WSGIServer
from gevent.queue import Queue
from gevent.monkey import patch_all
from gevent.subprocess import Popen, PIPE, STDOUT
patch_all()

import os
import time
import json
import glob

from flask import Flask, request, render_template, send_from_directory, Response
app = Flask(__name__)

import requests
# Disable InsecureRequestWarning "Unverified HTTPS request is being..."
from requests.packages.urllib3.exceptions import InsecureRequestWarning, InsecurePlatformWarning, SNIMissingWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)
requests.packages.urllib3.disable_warnings(InsecurePlatformWarning)
requests.packages.urllib3.disable_warnings(SNIMissingWarning)

import cyrtranslit

CACHEDIR = "api-cache"


def load_translations():
    """
    load all translations
    """
    translations = {}

    langs = ['az', 'bs', 'ca', 'cy', 'eo', 'he', 'hr', 'hy', 'id', 'is', 'it', 'ja', 'kk', 'lv', 'mk', 'nb', 'nn', 'sl', 'uz']
    for f_name in langs:
        f_name = 'share/translations/%s.txt' % f_name
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


def find_srv_for_query( path, query ):
    return 'http://api.worldweatheronline.com/'

def load_content_and_headers( path, query ):
    timestamp = time.strftime( "%Y%m%d%H", time.localtime() )
    p = os.path.join( CACHEDIR, timestamp, path, query )
    try:
        return open( p, 'r' ).read(), json.loads( open( p+".headers", 'r' ).read() )
    except Exception, e:
        return None, None

def save_content_and_headers( path, query, content, headers ):
    timestamp = time.strftime( "%Y%m%d%H", time.localtime() )
    p = os.path.join( CACHEDIR, timestamp, path, query )
    d = os.path.dirname( p )
    if not os.path.exists( d ): os.makedirs( d )
    open( p+".headers", 'w' ).write( json.dumps( headers ) )
    open( p, 'w' ).write( content )

def translate(text, lang):
    translated = TRANSLATIONS.get(lang, {}).get(text, text)
    if text.encode('utf-8') == translated:
        print "%s: %s" % (lang, text)
    return translated

def cyr(to_translate):
    return cyrtranslit.to_cyrillic(to_translate)

def patch_greek(original):
    return original.decode('utf-8').replace(u"Ηλιόλουστη/ο", u"Ηλιόλουστη").encode('utf-8')

def add_translations(content, lang):
    languages_to_translate = TRANSLATIONS.keys()
    #print type(content['data']) #['current_condition'][0]['lang_xx'] = [{'value':'XXX'}]
    try:
        d = json.loads(content)
    except Exception as e:
        print "---"
        print e
        print "---"

    try:
        weather_condition = d['data']['current_condition'][0]['weatherDesc'][0]['value']
        if lang in languages_to_translate:
            d['data']['current_condition'][0]['lang_%s' % lang] = [{'value':translate(weather_condition, lang)}]
            print translate(weather_condition, lang)
        elif lang == 'sr':
            d['data']['current_condition'][0]['lang_%s' % lang] = [{'value':cyr(d['data']['current_condition'][0]['lang_%s' % lang][0]['value'].encode('utf-8'))}]
        elif lang == 'el':
            d['data']['current_condition'][0]['lang_%s' % lang] = [{'value':patch_greek(d['data']['current_condition'][0]['lang_%s' % lang][0]['value'].encode('utf-8'))}]
        elif lang == 'sr-lat':
            d['data']['current_condition'][0]['lang_%s' % lang] = [{'value':d['data']['current_condition'][0]['lang_sr'][0]['value'].encode('utf-8')}]

        fixed_weather = []
        for w in d['data']['weather']:
            fixed_hourly = []
            for h in w['hourly']:
                weather_condition = h['weatherDesc'][0]['value']
                if lang in languages_to_translate:
                    h['lang_%s' % lang] = [{'value':translate(weather_condition, lang)}]
                elif lang == 'sr':
                    h['lang_%s' % lang] = [{'value':cyr(h['lang_%s' % lang][0]['value'].encode('utf-8'))}]
                elif lang == 'el':
                    h['lang_%s' % lang] = [{'value':patch_greek(h['lang_%s' % lang][0]['value'].encode('utf-8'))}]
                elif lang == 'sr-lat':
                    h['lang_%s' % lang] = [{'value':h['lang_sr'][0]['value'].encode('utf-8')}]
                fixed_hourly.append(h)
            w['hourly'] = fixed_hourly
            fixed_weather.append(w)
        d['data']['weather'] = fixed_weather 

        content = json.dumps(d)
    except Exception as e:
        print e
    return content

@app.route("/<path:path>")
def proxy(path):
    lang = request.args.get('lang', 'en')
    query_string = request.query_string
    query_string = query_string.replace('sr-lat', 'sr')
    content, headers = load_content_and_headers(path,query_string)

    if content is None:
        srv = find_srv_for_query(path, query_string)
        url = '%s/%s?%s' % (srv, path, query_string)
        print url

        attempts = 5
        while attempts:
            r = requests.get(url, timeout=10)
            try:
                json.loads(r.content)
                break
            except:
                attempts -= 1

        headers = {}
        headers['Content-Type'] = r.headers['content-type']
        content = add_translations(r.content, lang)
        try:
          save_content_and_headers( path, query_string, content, headers )
        except Exception, e:
          print e

    return content, 200, headers

if __name__ == "__main__":
    #app.run(host='0.0.0.0', port=5001, debug=False)
    #app.debug = True
    server = WSGIServer(("", 5001), app)
    server.serve_forever()

