#vim: fileencoding=utf-8

"""
One-line output mode.

Initial implementation of one-line output mode.

[ ] forecast
[ ] spark
[ ] several locations
[ ] location handling
[ ] more preconfigured format lines
[ ] add information about this mode to /:help
"""

import sys
import re
import datetime
import json
import requests

from astral import LocationInfo
from astral import moon
from astral.sun import sun

import pytz

from constants import WWO_CODE, WEATHER_SYMBOL, WEATHER_SYMBOL_WI_NIGHT, WEATHER_SYMBOL_WI_DAY, WIND_DIRECTION, WIND_DIRECTION_WI, WEATHER_SYMBOL_WIDTH_VTE, WEATHER_SYMBOL_PLAIN
from weather_data import get_weather_data
from . import v2
from . import v3
from . import prometheus

PRECONFIGURED_FORMAT = {
    '1':    r'%c %t\n',
    '2':    r'%c 🌡️%t 🌬️%w\n',
    '3':    r'%l: %c %t\n',
    '4':    r'%l: %c 🌡️%t 🌬️%w\n',
    '69':   r'nice',
}

MOON_PHASES = (
    u"🌑", u"🌒", u"🌓", u"🌔", u"🌕", u"🌖", u"🌗", u"🌘"
)

def convert_to_fahrenheit(temp):
    "Convert Celcius `temp` to Fahrenheit"

    return (temp*9.0/5)+32

def render_temperature(data, query):
    """
    temperature (t)
    """

    if query.get('use_imperial', False):
        temperature = u'%s°F' % data['temp_F']
    else:
        temperature = u'%s°C' % data['temp_C']

    if temperature[0] != '-':
        temperature = '+' + temperature

    return temperature

def render_feel_like_temperature(data, query):
    """
    feel like temperature (f)
    """

    if query.get('use_imperial', False):
        temperature = u'%s°F' % data['FeelsLikeF']
    else:
        temperature = u'%s°C' % data['FeelsLikeC']

    if temperature[0] != '-':
        temperature = '+' + temperature

    return temperature

def render_condition(data, query):
    """Emoji encoded weather condition (c)
    """

    if query.get("view") == "v2n":
        weather_condition = WEATHER_SYMBOL_WI_NIGHT.get(
                WWO_CODE.get(
                    data['weatherCode'], "Unknown"))
        spaces = "  "
    elif query.get("view") == "v2d":
        weather_condition = WEATHER_SYMBOL_WI_DAY.get(
                WWO_CODE.get(
                    data['weatherCode'], "Unknown"))
        spaces = "  "
    else:
        weather_condition = WEATHER_SYMBOL.get(
                WWO_CODE.get(
                    data['weatherCode'], "Unknown"))
        spaces = " "*(3 - WEATHER_SYMBOL_WIDTH_VTE.get(weather_condition, 1))

    return weather_condition + spaces

def render_condition_fullname(data, query):
    """
    condition_fullname (C)
    """

    found = None
    for key, val in data.items():
        if key.startswith('lang_'):
            found = val
            break
    if not found:
        found = data['weatherDesc']

    try:
        weather_condition = found[0]['value']
    except KeyError:
        weather_condition = ''

    return weather_condition

def render_condition_plain(data, query):
    """Plain text weather condition (x)
    """

    weather_condition = WEATHER_SYMBOL_PLAIN[WWO_CODE[data['weatherCode']]]

    return weather_condition

def render_humidity(data, query):
    """
    humidity (h)
    """

    humidity = data.get('humidity', '')
    if humidity:
        humidity += '%'
    return humidity

def render_precipitation(data, query):
    """
    precipitation (p)
    """

    answer = data.get('precipMM', '')
    if answer:
        answer += 'mm'
    return answer

def render_precipitation_chance(data, query):
    """
    precipitation chance (o)
    """

    answer = data.get('chanceofrain', '')
    if answer:
        answer += '%'
    return answer

def render_pressure(data, query):
    """
    pressure (P)
    """

    answer = data.get('pressure', '')
    if answer:
        answer += 'hPa'
    return answer

def render_uv_index(data, query):
    """
    UV Index (u)
    """

    answer = data.get('uvIndex', '')
    return answer

def render_wind(data, query):
    """
    wind (w)
    """

    try:
        degree = data["winddirDegree"]
    except KeyError:
        degree = ""

    try:
        degree = int(degree)
    except ValueError:
        degree = ""

    if degree:
        if query.get("view") in ["v2n", "v2d"]:
            wind_direction = WIND_DIRECTION_WI[int(((degree+22.5)%360)/45.0)]
        else:
            wind_direction = WIND_DIRECTION[int(((degree+22.5)%360)/45.0)]
    else:
        wind_direction = ""

    if query.get('use_ms_for_wind', False):
        unit = 'm/s'
        wind = u'%s%.1f%s' % (wind_direction, float(data['windspeedKmph'])/36.0*10.0, unit)
    elif query.get('use_imperial', False):
        unit = 'mph'
        wind = u'%s%s%s' % (wind_direction, data['windspeedMiles'], unit)
    else:
        unit = 'km/h'
        wind = u'%s%s%s' % (wind_direction, data['windspeedKmph'], unit)

    return wind

def render_location(data, query):
    """
    location (l)
    """

    return (data['override_location'] or data['location'])

def render_moonphase(_, query):
    """moonpahse(m)
    A symbol describing the phase of the moon
    """
    moon_phase = moon.phase(date=datetime.datetime.today())
    moon_index = int(int(32.0*moon_phase/28+2)%32/4)
    return MOON_PHASES[moon_index]

def render_moonday(_, query):
    """moonday(M)
    An number describing the phase of the moon (days after the New Moon)
    """
    moon_phase = moon.phase(date=datetime.datetime.today())
    return str(int(moon_phase))

##################################
# this part should be rewritten
# this is just a temporary solution

def get_geodata(location):
    # text = requests.get("http://localhost:8004/%s" % location).text
    text = requests.get("http://127.0.0.1:8083/:geo-location?location=%s" % location).text
    return json.loads(text)


def render_dawn(data, query, local_time_of):
    """dawn (D)
    Local time of dawn"""
    return local_time_of("dawn")

def render_dusk(data, query, local_time_of):
    """dusk (d)
    Local time of dusk"""
    return local_time_of("dusk")

def render_sunrise(data, query, local_time_of):
    """sunrise (S)
    Local time of sunrise"""
    return local_time_of("sunrise")

def render_sunset(data, query, local_time_of):
    """sunset (s)
    Local time of sunset"""
    return local_time_of("sunset")

def render_zenith(data, query, local_time_of):
    """zenith (z)
    Local time of zenith"""
    return local_time_of("noon")

def render_local_time(data, query, local_time_of):
    """local_time (T)
    Local time"""
    return "%{{NOW("+ local_time_of("TZ") +")}}"

def render_local_timezone(data, query, local_time_of):
    """local_time (Z)
    Local time"""
    return local_time_of("TZ")

##################################

FORMAT_SYMBOL = {
    'c':    render_condition,
    'C':    render_condition_fullname,
    'x':    render_condition_plain,
    'h':    render_humidity,
    't':    render_temperature,
    'f':    render_feel_like_temperature,
    'w':    render_wind,
    'l':    render_location,
    'm':    render_moonphase,
    'M':    render_moonday,
    'p':    render_precipitation,
    'o':    render_precipitation_chance,
    'P':    render_pressure,
    "u":    render_uv_index,
    }

FORMAT_SYMBOL_ASTRO = {
    'D':    render_dawn,
    'd':    render_dusk,
    'S':    render_sunrise,
    's':    render_sunset,
    'z':    render_zenith,

    'T':    render_local_time,
    'Z':    render_local_timezone,
}

def render_line(line, data, query):
    """
    Render format `line` using `data`
    """

    def get_local_time_of():

        location = data["location"]
        geo_data = get_geodata(location)

        city = LocationInfo()
        city.latitude = geo_data["latitude"]
        city.longitude = geo_data["longitude"]
        city.timezone = geo_data["timezone"]

        timezone = city.timezone

        local_tz = pytz.timezone(timezone)

        datetime_day_start = datetime.datetime.now()\
                .replace(hour=0, minute=0, second=0, microsecond=0)
        current_sun = sun(city.observer, date=datetime_day_start)

        local_time_of = lambda x: city.timezone if x == "TZ" else \
                                     current_sun[x]\
                                    .replace(tzinfo=pytz.utc)\
                                    .astimezone(local_tz)\
                                    .strftime("%H:%M:%S")
        return local_time_of

    def render_symbol(match):
        """
        Render one format symbol from re `match`
        using `data` from external scope.
        """

        symbol_string = match.group(0)
        symbol = symbol_string[-1]

        if symbol in FORMAT_SYMBOL:
            render_function = FORMAT_SYMBOL[symbol]
            return render_function(data, query)
        if symbol in FORMAT_SYMBOL_ASTRO and local_time_of is not None:
            render_function = FORMAT_SYMBOL_ASTRO[symbol]
            return render_function(data, query, local_time_of)

        return ''

    template_regexp = r'%[a-zA-Z]'
    for template_code in re.findall(template_regexp, line):
        if template_code.lstrip("%") in FORMAT_SYMBOL_ASTRO:
            local_time_of = get_local_time_of()
            break

    return re.sub(template_regexp, render_symbol, line)

def render_json(data):
    output = json.dumps(data, indent=4, sort_keys=True, ensure_ascii=False)

    output = "\n".join(
        re.sub('"[^"]*worldweatheronline[^"]*"', '""', line) if "worldweatheronline" in line else line
        for line in output.splitlines()) + "\n"

    return output

def format_weather_data(query, parsed_query, data):
    """
    Format information about current weather `data` for `location`
    with specified in `format_line` format
    """

    if 'data' not in data:
        return 'Unknown location; please try ~%s' % parsed_query["location"]

    format_line = parsed_query.get("view", "")
    if format_line in PRECONFIGURED_FORMAT:
        format_line = PRECONFIGURED_FORMAT[format_line]

    if format_line in ["j1", "j2"]:
        # j2 is a lightweight j1, without 'hourly' in 'weather' (weather forecast)
        if "weather" in data["data"] and format_line == "j2":
            for i in range(len(data["data"]["weather"])):
                del data["data"]["weather"][i]["hourly"]
        return render_json(data['data'])
    if format_line == "p1":
        return prometheus.render_prometheus(data['data'])
    if format_line[:2] == "v2":
        return v2.main(query, parsed_query, data)
    if format_line[:2] == "v3":
        return v3.main(query, parsed_query, data)

    current_condition = data['data']['current_condition'][0]
    current_condition['location'] = parsed_query["location"]
    current_condition['override_location'] = parsed_query["override_location_name"]
    output = render_line(format_line, current_condition, query)
    output = output.rstrip("\n").replace(r"\n", "\n")
    return output

def wttr_line(query, parsed_query):
    """
    Return 1line weather information for `location`
    in format `line_format`
    """
    location = parsed_query['location']
    lang = parsed_query['lang']

    data = get_weather_data(location, lang)
    output = format_weather_data(query, parsed_query, data)
    return output

def main():
    """
    Function for standalone module usage
    """

    location = sys.argv[1]
    query = {
        'line': sys.argv[2],
        }
    parsed_query = {
        "location": location,
        "orig_location": location,
        "language": "en",
        "format": "v2",
        }

    sys.stdout.write(wttr_line(query, parsed_query))

if __name__ == '__main__':
    main()
