"""
All location related functions and converters.

The main entry point is `location_processing`
which gets `location` and `source_ip_address`
and basing on this information generates
precise location description.

"""
from __future__ import print_function

import sys
import os
import json
import socket
import requests
import geoip2.database
import pycountry

from globals import GEOLITE, GEOLOCATOR_SERVICE, IP2LCACHE, IP2LOCATION_KEY, NOT_FOUND_LOCATION, \
                    ALIASES, BLACKLIST, IATA_CODES_FILE, IPLOCATION_ORDER, IPINFO_TOKEN

GEOIP_READER = geoip2.database.Reader(GEOLITE)

COUNTRY_MAP = {"Russian Federation": "Russia"}

def ascii_only(string):
    "Check if `string` contains only ASCII symbols"

    try:
        for _ in range(5):
            string = string.encode('utf-8')
        return True
    except UnicodeDecodeError:
        return False

def is_ip(ip_addr):
    """
    Check if `ip_addr` looks like an IP Address
    """

    if sys.version_info[0] < 3:
        ip_addr = ip_addr.encode("utf-8")

    try:
        socket.inet_pton(socket.AF_INET, ip_addr)
        return True
    except socket.error:
        try:
            socket.inet_pton(socket.AF_INET6, ip_addr)
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



def geolocator(location):
    """
    Return a GPS pair for specified `location` or None
    if nothing can be found
    """

    try:
        geo = requests.get('%s/%s' % (GEOLOCATOR_SERVICE, location)).text
    except requests.exceptions.ConnectionError as exception:
        print("ERROR: %s" % exception)
        return None

    if geo == "":
        return None

    try:
        answer = json.loads(geo.encode('utf-8'))
        return answer
    except ValueError as exception:
        print("ERROR: %s" % exception)
        return None

    return None


def ipcachewrite(ip_addr, location):
    """ Write a retrieved ip+location into cache
        Can stress some filesystems after long term use, see
        https://stackoverflow.com/questions/466521/how-many-files-can-i-put-in-a-directory
    """
    cachefile = os.path.join(IP2LCACHE, ip_addr)
    if not os.path.exists(IP2LCACHE):
        os.makedirs(IP2LCACHE)

    with open(cachefile, 'w') as file:
        file.write(location[3] + ';' + location[2] + ';' + location[1] + ';' + location[0] + ';' + location[4] + ';' + location[5])
        # like ip2location format

def ipcache(ip_addr):
    """ Retrieve a location from cache by ip addr
        Returns a triple of (CITY, REGION, COUNTRY) or None
        TODO: When cache becomes more robust, transition to using latlong
    """
    cachefile = os.path.join(IP2LCACHE, ip_addr)
    if not os.path.exists(IP2LCACHE):
        os.makedirs(IP2LCACHE)

    if os.path.exists(cachefile):
        try:
            _, country, region, city, *_ = open(cachefile, 'r').read().split(';')
            return city, region, country
        except ValueError:
            # cache entry is malformed: should be
            # [ccode];country;region;city;[lat];[long];...
            return None
    return None


def ip2location(ip_addr):
    """Convert IP address `ip_addr` to a location name"""
    # if IP2LOCATION_KEY is not set, do not query,
    # because the query wont be processed anyway
    if not IP2LOCATION_KEY:
        return None
    try:
        r = requests.get(
            'http://api.ip2location.com/?ip=%s&key=%s&package=WS5'
            % (ip_addr, IP2LOCATION_KEY))
        r.raise_for_status()
        location = r.text
        if location and ';' in location:
            ccode, country, region, city, lat, long, *_ = location.split(';')
    except requests.exceptions.RequestException:
        return None
    return city, region, country, ccode, lat, long


def ipinfo(ip_addr):
    if not IPINFO_TOKEN:
        return None
    try:
        r = requests.get(
            'https://ipinfo.io/%s/json?token=%s'
            % (ip_addr, IPINFO_TOKEN))
        r.raise_for_status()
        r_json = r.json()
        # can't do two unpackings on one line
        city, region, country, ccode  = r_json["city"], r_json["region"], '', r_json["country"],
        lat, long = r_json["loc"].split(',')
        # NOTE: ipinfo only provides ISO codes for countries
        country = pycountry.countries.get(alpha_2=ccode).name
    except (requests.exceptions.RequestException, ValueError):
        # latter is thrown by failure to parse json in reponse
        return None
    return city, region, country, ccode, lat, long


def geoip(ip_addr):
    try:
        response = GEOIP_READER.city(ip_addr)
        city, region, country, ccode, lat, long = response.city.name, response.subdivisions.name, response.country.name, response.country.iso_code, response.location.latitude, response.location.longitude
    except geoip2.errors.AddressNotFoundError:
        return None
    return city, region, country, ccode, lat, long


def workaround(country):
    # workaround for strange bug with the country name
    # maybe some other countries has this problem too
    country = COUNTRY_MAP.get(country) or country
    return country

def get_location(ip_addr):
    """
    Return location triple (CITY, REGION, COUNTRY) for `ip_addr`
    """
    location = ipcache(ip_addr)
    if location:
        return location

    # location from iplocators have the following order:
    # (CITY, REGION, COUNTRY, CCODE, LAT, LONG)
    for method in IPLOCATION_ORDER:
        if method == 'geoip':
            location = geoip(ip_addr)
        elif method == 'ip2location':
            location = ip2location(ip_addr)
        elif method == 'ipinfo':
            location = ipinfo(ip_addr)
        else:
            print("ERROR: invalid iplocation method specified: %s" % method)
        if location is not None:
            break

    if location is not None and all(location):
        ipcachewrite(ip_addr, location)
        # cache write used to happen before workaround, preserve that
        location[2] = workaround(location[2])
        return location[:3]  # city, region, country
        # ccode is cached but not needed for location

    # temporary disabled it because of geoip services capcacity
    #
    #if city is None and response.location:
    #    coord = "%s, %s" % (response.location.latitude, response.location.longitude)
    #    try:
    #        location = geolocator.reverse(coord, language='en')
    #        city = location.raw.get('address', {}).get('city')
    #    except:
    #        pass

    # No methods resulted in a location - return default
    return NOT_FOUND_LOCATION, None, None


def location_canonical_name(location):
    "Find canonical name for `location`"

    location = location_normalize(location)
    if location.lower() in LOCATION_ALIAS:
        return LOCATION_ALIAS[location.lower()]
    return location

def load_aliases(aliases_filename):
    """
    Load aliases from the aliases file
    """
    aliases_db = {}
    with open(aliases_filename, 'r') as f_aliases:
        for line in f_aliases.readlines():
            try:
                from_, to_ = line.decode('utf-8').split(':', 1)
            except AttributeError:
                from_, to_ = line.split(':', 1)

            aliases_db[location_normalize(from_)] = location_normalize(to_)
    return aliases_db

def load_iata_codes(iata_codes_filename):
    """
    Load IATA codes from the IATA codes file
    """
    with open(iata_codes_filename, 'r') as f_iata_codes:
        result = []
        for line in f_iata_codes.readlines():
            result.append(line.strip())
    return set(result)

LOCATION_ALIAS = load_aliases(ALIASES)
LOCATION_BLACK_LIST = [x.strip() for x in open(BLACKLIST, 'r').readlines()]
IATA_CODES = load_iata_codes(IATA_CODES_FILE)

def is_location_blocked(location):
    """
    Return True if this location is blocked
    or False if it is allowed
    """
    return location is not None and location.lower() in LOCATION_BLACK_LIST


def get_hemisphere(location):
    """
    Return hemisphere of the location (True = North, False = South).
    Assume North and return True if location can't be found.
    """
    if all(location):
        location_string = ", ".join(location)

    geolocation = geolocator(location_string)
    if geolocation is None:
        return True
    return geolocation["latitude"] > 0

def location_processing(location, ip_addr):
    """
    """

    # if location is starting with ~
    # or has non ascii symbols
    # it should be handled like a search term (for geolocator)
    override_location_name = None
    full_address = None
    hide_full_address = False
    force_show_full_address = location is not None and location.startswith('~')

    # location ~ means that it should be detected automatically,
    # and shown in the location line below the report
    if location == '~':
        location = None

    if location and location.lstrip('~ ').startswith('@'):
        try:
            location, region, country = get_location(
                socket.gethostbyname(
                    location.lstrip('~ ')[1:]))
            location = '~' + location
            if region and country:
                location += ", %s, %s" % (region, country)
            hide_full_address = not force_show_full_address
        except:
            location, region, country = NOT_FOUND_LOCATION, None, None

    query_source_location = get_location(ip_addr)

    # For moon queries, hemisphere must be found
    # True for North, False for South
    hemisphere = False
    if location is not None and (location.lower()+"@").startswith("moon@"):
        hemisphere = get_hemisphere(query_source_location)

    country = None
    if not location or location == 'MyLocation':
        location = ip_addr

    if is_ip(location):
        location, region, country = get_location(location)
        # location is just city here

        # here too
        if location:
            location = '~' + location
            if region and country:
                location += ", %s, %s" % (region, country)
            hide_full_address = not force_show_full_address

    if location and not location.startswith('~'):
        tmp_location = location_canonical_name(location)
        if tmp_location != location:
            override_location_name = location
            location = tmp_location

    # up to this point it is possible that the name
    # contains some unicode symbols
    # here we resolve them
    if location is not None: # and not ascii_only(location):
        location = "~" + location.lstrip('~ ')
        if not override_location_name:
            override_location_name = location.lstrip('~')

    # if location is not None and location.upper() in IATA_CODES:
    #     location = '~%s' % location

    if location is not None and not location.startswith("~-,") and location.startswith('~'):
        geolocation = geolocator(location_canonical_name(location[1:]))
        if geolocation is not None:
            if not override_location_name:
                override_location_name = location[1:].replace('+', ' ')
            location = "%s,%s" % (geolocation['latitude'], geolocation['longitude'])
            country = None
            if not hide_full_address:
                full_address = geolocation['address']
            else:
                full_address = None
        else:
            location = NOT_FOUND_LOCATION #location[1:]


    return location, \
            override_location_name, \
            full_address, \
            country, \
            query_source_location, \
            hemisphere
