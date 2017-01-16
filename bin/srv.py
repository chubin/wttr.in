import sys
import logging
import os
import re
import requests
import socket
import subprocess
import time
import traceback

import geoip2.database
from geopy.geocoders import Nominatim
import jinja2

import gevent
from gevent.wsgi import WSGIServer
from gevent.queue import Queue
from gevent.monkey import patch_all
from gevent.subprocess import Popen, PIPE, STDOUT
patch_all()

from flask import Flask, request, render_template, send_from_directory
app = Flask(__name__)

MYDIR = os.environ.get('WTTR_MYDIR', os.path.abspath(os.path.dirname( os.path.dirname('__file__') )))
GEOLITE = os.environ.get('WTTR_GEOLITE', os.path.join( MYDIR, "GeoLite2-City.mmdb" ))
WEGO = os.environ.get('WTTR_WEGO', "/home/igor/go/bin/wego")

CACHEDIR  = os.path.join( MYDIR, "cache" )
IP2LCACHE = os.path.join( MYDIR, "cache/ip2l" )
ALIASES   = os.path.join( MYDIR, "share/aliases" )
ANSI2HTML = os.path.join( MYDIR, "share/ansi2html.sh" )
HELP_FILE = os.path.join( MYDIR, 'share/help.txt' )
LOG_FILE  = os.path.join( MYDIR, 'log/main.log' )
TEMPLATES = os.path.join( MYDIR, 'share/templates' )
STATIC    = os.path.join( MYDIR, 'share/static' )

NOT_FOUND_LOCATION = "NOT_FOUND"
DEFAULT_LOCATION = "Oymyakon"

NOT_FOUND_MESSAGE = """
We were unable to find your location,
so we have brought you to Oymyakon,
one of the coldest permanently inhabited locales on the planet.
"""

if not os.path.exists(os.path.dirname( LOG_FILE )):
    os.makedirs( os.path.dirname( LOG_FILE ) )
logging.basicConfig(filename=LOG_FILE, level=logging.DEBUG)

reader = geoip2.database.Reader(GEOLITE)
geolocator = Nominatim()

my_loader = jinja2.ChoiceLoader([
    app.jinja_loader,
    jinja2.FileSystemLoader(TEMPLATES),
])
app.jinja_loader = my_loader


class Limits:
    def __init__( self ):
        self.intervals = [ 'min', 'hour', 'day' ]
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
            'min':      10,
            'hour':     20,
            'day':      100,
            }
        self.last_update = {
            'min':      0,
            'hour':     0,
            'day':      0,
            }
        self.clear_counters()

    def check_ip( self, ip ):
        self.clear_counters()
        for interval in self.intervals:
            if ip not in self.counter[interval]:
                self.counter[interval][ip] = 0
            self.counter[interval][ip] += 1
            if self.limit[interval] <= self.counter[interval][ip]:
                log("Too many queries: %s in %s for %s" % (self.limit[interval], interval, ip) )
                raise RuntimeError("Not so fast! Number of queries per %s is limited to %s" % (interval, self.limit[interval]))
            print self.counter

    def clear_counters( self ):
        t = int( time.time() )
        for interval in self.intervals:
            if t / self.divisor[interval] != self.last_update[interval]:
                self.counter[interval] = {}
                self.last_update[interval] = t / self.divisor[interval]
        

limits = Limits()

def error( text ):
    print text
    raise RuntimeError(text)

def log( text ):
    print text
    logging.info( text )

def is_ip( ip ):
    if re.match('\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', ip) is None:
        return False
    try:
        socket.inet_aton(ip)
        return True
    except socket.error:
        return False

def save_weather_data( location, filename ):

    if location == NOT_FOUND_LOCATION:
        location_not_found = True
        location = DEFAULT_LOCATION
    else:
        location_not_found = False
    
    p = Popen( [ WEGO, '--city=%s' % location ], stdout=PIPE, stderr=PIPE )
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


def ip2location( ip ):
    cached = os.path.join( IP2LCACHE, ip )

    if os.path.exists( cached ):
        return open( cached, 'r' ).read()

    try:
        t = requests.get( 'http://api.ip2location.com/?ip=%s&key=demo&package=WS10' % ip ).text
        if ';' in t:
            location = t.split(';')[3]
            if not os.path.exists( IP2LCACHE ):
                os.makedirs( IP2LCACHE )
            open( cached, 'w' ).write( location )
            return location
    except:
            pass

def get_location( ip_addr ):
    response = reader.insights( ip_addr )
    city = response.city.name + ',' + response.subdivisions.most_specific.name + ',' + response.country.name

    if city is None and response.location:
        coord = "%s, %s" % (response.location.latitude, response.location.longitude)
        location = geolocator.reverse(coord, language='en')
        city = location.raw.get('address', {}).get('address')
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

def show_help():
    return open(HELP_FILE, 'r').read()

@app.route('/files/<path:path>')
def send_static(path):
    return send_from_directory(STATIC, path)

@app.route('/favicon.ico')
def send_favicon():
    return send_from_directory(STATIC, 'favicon.ico')

@app.route("/")
@app.route("/<string:location>")
def wttr(location = None):
    user_agent = request.headers.get('User-Agent').lower()

    html_output = True
    if 'curl' in user_agent or 'wget' in user_agent or 'httpie' in user_agent or 'lwp-request' in user_agent:
        html_output = False

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
            location = get_location( socket.gethostbyname( location[1:] ) )

        location = location_canonical_name( location )
        log("%s %s %s %s" % (ip, user_agent, orig_location, location))
        return get_wetter( location, ip, html=html_output )
    except Exception, e:
        logging.error("Exception has occured", exc_info=1)
        return str(e).rstrip()+"\n"

server = WSGIServer(("", 8002), app)
server.serve_forever()

