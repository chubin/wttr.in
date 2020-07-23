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

from globals import GEOLITE, GEOLOCATOR_SERVICE, IP2LCACHE, IP2LOCATION_KEY, NOT_FOUND_LOCATION, \
                    ALIASES, BLACKLIST, IATA_CODES_FILE

GEOIP_READER = geoip2.database.Reader(GEOLITE)

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
    if nothing can't be found
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

def ip2location(ip_addr):
    "Convert IP address `ip_addr` to a location name"

    cached = os.path.join(IP2LCACHE, ip_addr)
    if not os.path.exists(IP2LCACHE):
        os.makedirs(IP2LCACHE)

    location = None

    if os.path.exists(cached):
        location = open(cached, 'r').read()
    else:
        # if IP2LOCATION_KEY is not set, do not the query,
        # because the query wont be processed anyway
        if IP2LOCATION_KEY:
            try:
                ip2location_response = requests\
                        .get('http://api.ip2location.com/?ip=%s&key=%s&package=WS3' \
                                % (ip_addr, IP2LOCATION_KEY)).text
                if ';' in ip2location_response:
                    open(cached, 'w').write(ip2location_response)
                location = ip2location_response
            except requests.exceptions.ConnectionError:
                pass

    if location and ';' in location:
        location = location.split(';')[3], location.split(';')[1]
    else:
        location = location, None

    return location

def get_location(ip_addr):
    """
    Return location pair (CITY, COUNTRY) for `ip_addr`
    """

    try:
        response = GEOIP_READER.city(ip_addr)
        country = response.country.name
        city = response.city.name
    except geoip2.errors.AddressNotFoundError:
        country = None
        city = None

    #
    # temporary disabled it because of geoip services capcacity
    #
    #if city is None and response.location:
    #    coord = "%s, %s" % (response.location.latitude, response.location.longitude)
    #    try:
    #        location = geolocator.reverse(coord, language='en')
    #        city = location.raw.get('address', {}).get('city')
    #    except:
    #        pass
    if city is None:
        city, country = ip2location(ip_addr)

    # workaround for the strange bug with the country name
    # maybe some other countries has this problem too
    if country == 'Russian Federation':
        country = 'Russia'

    if city:
        return city, country
    else:
        return NOT_FOUND_LOCATION, None


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
            location, country = get_location(
                socket.gethostbyname(
                    location.lstrip('~ ')[1:]))
            location = '~' + location
            if country:
                location += ", %s" % country
            hide_full_address = not force_show_full_address
        except:
            location, country = NOT_FOUND_LOCATION, None

    query_source_location = get_location(ip_addr)

    country = None
    if not location or location == 'MyLocation':
        location = ip_addr

    if is_ip(location):
        location, country = get_location(location)

        # here too
        if location:
            location = '~' + location
            if country:
                location += ", %s" % country
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
            query_source_location
