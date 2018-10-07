"""
All location related functions and converters.
"""

import json
import requests
import geoip2.database

GEOIP_READER = geoip2.database.Reader(GEOLITE)

def ascii_only(string):
    "Check if `string` contains only ASCII symbols"

    try:
        for _ in range(5):
            string = string.encode('utf-8')
        return True
    except UnicodeDecodeError:
        return False


def geolocator(location):
    """
    Return a GPS pair for specified `location` or None
    if nothing can't be found
    """

    try:
        geo = requests.get('%s/%s' % (GEOLOCATOR_SERVICE, location)).text
    except requests.exceptions.ConnectionError as exception:
        print "ERROR: %s" % exception
        return None

    if geo == "":
        return None

    try:
        answer = json.loads(geo.encode('utf-8'))
        return answer
    except ValueError as exception:
        print "ERROR: %s" % exception
        return None

    return None

def ip2location(ip_addr):
    "Convert IP address `ip_addr` to a location name"

    cached = os.path.join(IP2LCACHE, ip_addr)
    if not os.path.exists(IP2LCACHE):
        os.makedirs(IP2LCACHE)

    if os.path.exists(cached):
        location = open(cached, 'r').read()
        return location

    try:
        ip2location_response = requests\
                .get('http://api.ip2location.com/?ip=%s&key=%s&package=WS10' \
                        % (IP2LOCATION_KEY, ip_addr)).text
        if ';' in ip2location_response:
            location = ip2location_response.split(';')[3]
            open(cached, 'w').write(location)
            print "ip2location says: %s" % location
            return location
    except requests.exceptions.ConnectionError as exception:
        return None

def get_location(ip_addr):
    """
    Return location pair (CITY, COUNTRY) for `ip_addr`
    """

    response = GEOIP_READER.city(ip_addr)
    country = response.country.iso_code
    city = response.city.name

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
        city = ip2location(ip_addr)
    return (city or NOT_FOUND_LOCATION), country


def location_canonical_name(location):
    "Find canonical name for `location`"

    location = location_normalize(location)
    if location in LOCATION_ALIAS:
        return LOCATION_ALIAS[location.lower()]
    return location

def load_aliases(aliases_filename):
    """
    Load aliases from the aliases file
    """
    aliases_db = {}
    with open(aliases_filename, 'r') as f_aliases:
        for line in f_aliases.readlines():
            from_, to_ = line.decode('utf-8').split(':', 1)
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
    return location is not None and location.lower() in LOCATION_BLACK_LIST


def location_processing():

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

