
import gevent
from gevent.pywsgi import WSGIServer
from gevent.queue import Queue
from gevent.monkey import patch_all
from gevent.subprocess import Popen, PIPE, STDOUT
patch_all()

import sys
import os
import json

from flask import Flask, request, render_template, send_from_directory, send_file, make_response, jsonify, Response
app = Flask(__name__)

MYDIR = os.path.abspath(os.path.dirname('__file__'))
sys.path.append(os.path.join(MYDIR, 'lib'))

CACHEDIR = os.path.join(MYDIR, 'cache')

from geopy.geocoders import Nominatim #, Mapzen
#geomapzen = Mapzen("mapzen-RBNbmcZ") # Nominatim()
geoosm = Nominatim(timeout=7, user_agent="wttrin-geo/0.0.2")

import airports

# from tzwhere import tzwhere
import timezonefinder
tf = timezonefinder.TimezoneFinder()

def load_cache(location_string):
    try:
        location_string = location_string.replace('/', '_')
        cachefile = os.path.join(CACHEDIR, location_string)

        return json.loads(open(cachefile, 'r').read())
    except:
        return None

def shorten_full_address(address):
    parts = address.split(',')
    if len(parts) > 6:
        parts = parts[:2] + [x for x in parts[-4:] if len(x) < 20]
        return ','.join(parts)
    return address
    

def save_cache(location_string, answer):
    location_string = location_string.replace('/', '_')
    cachefile = os.path.join(CACHEDIR, location_string)
    open(cachefile, 'w').write(json.dumps(answer))

def query_osm(location_string):
    try:
        location = geoosm.geocode(location_string)
        return {
            'address':  location.address,
            'latitude': location.latitude,
            'longitude':location.longitude,
        }

    except Exception as e:
        print(e)
        return None

def add_timezone_information(geo_data):
    # tzwhere_ = tzwhere.tzwhere()
    # timezone_str = tzwhere_.tzNameAt(geo_data["latitude"], geo_data["longitude"])
    timezone_str = tf.certain_timezone_at(lat=geo_data["latitude"], lng=geo_data["longitude"])

    answer = geo_data.copy()
    answer["timezone"] = timezone_str

    return answer

@app.route("/<string:location>")
def find_location(location):

    airport_gps_location = airports.get_airport_gps_location(location.upper().lstrip('~'))
    is_airport = False
    if airport_gps_location is not None:
        location = airport_gps_location
        is_airport = True

    location = location.replace('+', ' ')
    answer = load_cache(location)
    loaded_answer = None

    if answer is not None:
        loaded_answer = answer.copy()
        print("cache found: %s" % location)

    if answer is None:
        answer = query_osm(location)
    
    if is_airport:
        answer['address'] = shorten_full_address(answer['address'])

    if "timezone" not in answer:
       answer = add_timezone_information(answer)

    if answer is not None and loaded_answer != answer:
        save_cache(location, answer)

    if answer is None:
        return ""
    else:
        r = Response(json.dumps(answer)) 
        r.headers["Content-Type"] = "text/json; charset=utf-8"
        return r

app.config['JSON_AS_ASCII'] = False
server = WSGIServer(("127.0.0.1", 8004), app)
server.serve_forever()
