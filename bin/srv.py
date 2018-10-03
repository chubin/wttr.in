#!/usr/bin/env python
# vim: set encoding=utf-8

import gevent
from gevent.wsgi import WSGIServer
from gevent.queue import Queue
from gevent.monkey import patch_all
from gevent.subprocess import Popen, PIPE, STDOUT
patch_all()

import sys
import logging
import os
import re
import requests
import socket
import time
import json

import geoip2.database
import jinja2

from flask import Flask, request, render_template, send_from_directory, send_file, make_response
app = Flask(__name__)

MYDIR = os.path.abspath(os.path.dirname( os.path.dirname('__file__') ))
sys.path.append("%s/lib/" % MYDIR)
import wttrin_png, parse_query
from translations import get_message, FULL_TRANSLATION, PARTIAL_TRANSLATION, SUPPORTED_LANGS
from buttons import TWITTER_BUTTON, GITHUB_BUTTON, GITHUB_BUTTON_2, GITHUB_BUTTON_3, GITHUB_BUTTON_FOOTER

from globals import GEOLITE, \
                    IP2LCACHE, ALIASES, BLACKLIST, \
                    get_help_file, BASH_FUNCTION_FILE, TRANSLATION_FILE, LOG_FILE, TEST_FILE, \
                    TEMPLATES, STATIC, \
                    NOT_FOUND_LOCATION, \
                    MALFORMED_RESPONSE_HTML_PAGE, \
                    IATA_CODES_FILE, \
                    log, error

from wttr import get_wetter, get_moon

if not os.path.exists(os.path.dirname(LOG_FILE)):
    os.makedirs(os.path.dirname(LOG_FILE))
logging.basicConfig(filename=LOG_FILE, level=logging.DEBUG, format='%(asctime)s %(message)s')

my_loader = jinja2.ChoiceLoader([
    app.jinja_loader,
    jinja2.FileSystemLoader(TEMPLATES),
])
app.jinja_loader = my_loader

class Limits:
    def __init__( self ):
        self.intervals = ['min', 'hour', 'day']
        self.divisor = {
            'min':      60,
            'hour':     3600,
            'day':      86400,
            }
        self.counter = {
            'min':      {},
            'hour':     {},
            'day':      {},  
            }
        self.limit = {
            'min':      30,
            'hour':     600,
            'day':      1000,
            }
        self.last_update = {
            'min':      0,
            'hour':     0,
            'day':      0,
            }
        self.clear_counters()

    def check_ip(self, ip):
        if ip == '5.9.243.177':
            return
        self.clear_counters()
        for interval in self.intervals:
            if ip not in self.counter[interval]:
                self.counter[interval][ip] = 0
            self.counter[interval][ip] += 1
            if self.limit[interval] <= self.counter[interval][ip]:
                log("Too many queries: %s in %s for %s" % (self.limit[interval], interval, ip) )
                raise RuntimeError("Not so fast! Number of queries per %s is limited to %s" % (interval, self.limit[interval]))

    def clear_counters( self ):
        t = int( time.time() )
        for interval in self.intervals:
            if t / self.divisor[interval] != self.last_update[interval]:
                self.counter[interval] = {}
                self.last_update[interval] = t / self.divisor[interval]

limits = Limits()

def error(text):
    print text
    raise RuntimeError(text)

def log(text):
    print text.encode('utf-8')
    logging.info( text.encode('utf-8') )

def is_ip(ip):
    if re.match('\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', ip) is None:
        return False
    try:
        socket.inet_aton(ip)
        return True
    except socket.error:
        return False

def location_normalize(location):
    #translation_table = dict.fromkeys(map(ord, '!@#$*;'), None)
    def remove_chars( c, s ):
        return ''.join(x for x in s if x not in c )

    location = location.lower().replace('_', ' ').replace('+', ' ').strip()
    if not location.startswith('moon@'):
        location = remove_chars(r'!@#$*;:\\', location)
    return location

def load_aliases(aliases_filename):
    aliases_db = {}
    with open(aliases_filename, 'r') as f:
        for line in f.readlines():
            from_, to_ = line.decode('utf-8').split(':', 1)
            aliases_db[location_normalize(from_)] = location_normalize(to_)
    return aliases_db

def load_iata_codes(iata_codes_filename):
    with open(iata_codes_filename, 'r') as f:
        result = []
        for line in f.readlines():
            result.append(line.strip())
    return set(result)

location_alias = load_aliases(ALIASES)
location_black_list = [x.strip() for x in open(BLACKLIST, 'r').readlines()]
iata_codes = load_iata_codes(IATA_CODES_FILE)
print "IATA CODES LOADED: %s" % len(iata_codes)

def location_canonical_name(location):
    location = location_normalize(location)
    if location in location_alias:
        return location_alias[location.lower()]
    return location

def ascii_only(s):
    try:
        for i in range(5):
            s = s.encode('utf-8')
        return True
    except:
        return False

def geolocator(location):
    try:
        geo = requests.get('http://localhost:8004/%s' % location).text
    except Exception as e:
        print "ERROR: %s" % e
        return

    if geo == "":
        return

    try:
        answer = json.loads(geo.encode('utf-8'))
        return answer
    except Exception as e:
        print "ERROR: %s" % e
        return None

def ip2location(ip):
    cached = os.path.join(IP2LCACHE, ip)
    if not os.path.exists(IP2LCACHE):
        os.makedirs(IP2LCACHE)

    if os.path.exists(cached):
        location = open(cached, 'r').read()
        return location

    try:
        t = requests.get( 'http://api.ip2location.com/?ip=%s&key=%s&package=WS10' % (IP2LOCATION_KEY, ip)).text
        if ';' in t:
            location = t.split(';')[3]
            open(cached, 'w').write(location)
            print "ip2location says: %s" % location
            return location
    except:
            pass

reader = geoip2.database.Reader(GEOLITE)
def get_location(ip_addr):
    response = reader.city(ip_addr)

    if location == NOT_FOUND_LOCATION:
        location_not_found = True
        location = DEFAULT_LOCATION
    else:
        location_not_found = False
    p = Popen( [ WEGO, '-location=%s' % location ], stdout=PIPE, stderr=PIPE )
    stdout, stderr = p.communicate()
    if p.returncode != 0:
        error( stdout + stderr )

    dirname = os.path.dirname( filename )
    if not os.path.exists( dirname ):
        os.makedirs( dirname )
    
    if location_not_found:
        stdout += NOT_FOUND_MESSAGE

    open( filename, 'w' ).write( stdout )

    p = Popen( [ "bash", ANSI2HTML, "--palette=solarized", "--bg=dark" ],  stdin=PIPE, stdout=PIPE, stderr=PIPE )
    stdout, stderr = p.communicate( stdout )
    if p.returncode != 0:
        error( stdout + stderr )

    open( filename+'.html', 'w' ).write( stdout )

def get_filename( location ):
    location = location.replace('/', '_')
    timestamp = time.strftime( "%Y%m%d%H", time.localtime() )
    return "%s/%s/%s" % ( CACHEDIR, location, timestamp )

def get_wetter(location, ip, html=False):
    filename = get_filename( location )
    if not os.path.exists( filename ):
        limits.check_ip( ip )
        save_weather_data( location, filename )
    if html:
        filename += '.html'
    return open(filename).read()



def get_location( ip_addr ):
    response = reader.city( ip_addr )
    city = response.city.name
    if city is None and response.location:
        coord = "%s, %s" % (response.location.latitude, response.location.longitude)
        location = geolocator.reverse(coord, language='en')
        city = location.raw.get('address', {}).get('city')
    if city is None:
        print ip_addr
        city = ip2location( ip_addr )
    return city or NOT_FOUND_LOCATION

def load_aliases( aliases_filename ):
    aliases_db = {}
    with open( aliases_filename, 'r' ) as f:
        for line in f.readlines():
            from_, to_ = line.split(':', 1)
            aliases_db[ from_.strip().lower() ] = to_.strip()
    return aliases_db

location_alias = load_aliases( ALIASES )
def location_canonical_name( location ):
    if location.lower() in location_alias:
        return location_alias[location.lower()]
    return location

def find_supported_language(accepted_languages):
    for p in accepted_languages:
        lang = p[0]
        if '-' in lang:
            lang = lang.split('-',1)[0]
        if lang in SUPPORTED_LANGS:
            return lang
    return None

def show_help(location, lang):
    text = ""
    if location == ":help":
        text = open(get_help_file(lang), 'r').read()
        text = text.replace('FULL_TRANSLATION', ' '.join(FULL_TRANSLATION))
        text = text.replace('PARTIAL_TRANSLATION', ' '.join(PARTIAL_TRANSLATION))
    elif location == ":bash.function":
        text = open(BASH_FUNCTION_FILE, 'r').read()
    elif location == ":translation":
        text = open(TRANSLATION_FILE, 'r').read()
        text = text.replace('NUMBER_OF_LANGUAGES', str(len(SUPPORTED_LANGS))).replace('SUPPORTED_LANGUAGES', ' '.join(SUPPORTED_LANGS))
    return text.decode('utf-8')
show_help.pages = [':help', ':bash.function', ':translation' ]

@app.route('/files/<path:path>')
def send_static(path):
    return send_from_directory(STATIC, path)

@app.route('/favicon.ico')
def send_favicon():
    return send_from_directory(STATIC, 'favicon.ico')

@app.route('/malformed-response.html')
def send_malformed():
    return send_from_directory(STATIC, 'malformed-response.html')

@app.route("/")
@app.route("/<string:location>")
def wttr(location = None):

    user_agent = request.headers.get('User-Agent').lower()

    if any(agent in user_agent for agent in PLAIN_TEXT_AGENTS):
        html_output = False
    else:
        html_output = True


    if location == ':help':
        help_ = show_help()
        if html_output:
            return render_template( 'index.html', body=help_ )
        else:
            return help_

    orig_location = location

    if request.headers.getlist("X-Forwarded-For"):
        ip = request.headers.getlist("X-Forwarded-For")[0]
        if ip.startswith('::ffff:'):
            ip = ip[7:]
    else:
        ip = request.remote_addr

    try:
        if location is None:
            location = get_location( ip )

        if is_ip( location ):
            location = get_location( location )
        if location.startswith('@'):
            try:
                loc = dns.resolver.query( location[1:], 'LOC' )
                location = str("%.7f,%.7f" % (loc[0].float_latitude, loc[0].float_longitude))
            except DNSException, e:
                location = get_location( socket.gethostbyname( location[1:] ) )

        location = location_canonical_name( location )
        log("%s %s %s %s" % (ip, user_agent, orig_location, location))
        return get_wetter( location, ip, html=html_output )
    except Exception, e:
        if 'Malformed response' in str(e) or 'API key has reached calls per day allowed limit' in str(e):
            if html_output:
                return MALFORMED_RESPONSE_HTML_PAGE
            else:
                return get_message('CAPACITY_LIMIT_REACHED', lang).encode('utf-8')
        logging.error("Exception has occured", exc_info=1)
        return "ERROR"

server = WSGIServer((LISTEN_HOST, LISTEN_PORT), app)
server.serve_forever()

