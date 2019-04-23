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
from astral import Astral, Location
from constants import WWO_CODE, WEATHER_SYMBOL, WIND_DIRECTION
from weather_data import get_weather_data

PRECONFIGURED_FORMAT = {
    '1':    u'%c %t',
    '2':    u'%c 🌡️%t 🌬️%w',
    '3':    u'%l: %c %t',
    '4':    u'%l: %c 🌡️%t 🌬️%w',
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

def render_condition(data, query):
    """
    condition (c)
    """

    weather_condition = WEATHER_SYMBOL[WWO_CODE[data['weatherCode']]]
    return weather_condition

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

def render_humidity(data, query):
    """
    humidity (h)
    """

    humidity = data.get('humidity', '')
    if humidity:
        humidity += '%'
    return humidity

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
        wind_direction = WIND_DIRECTION[((degree+22)%360)/45]
    else:
        wind_direction = ""

    if query.get('use_ms_for_wind', False):
        unit = ' m/s'
        wind = u'%s%.1f%s' % (wind_direction, float(data['windspeedKmph'])/36.0*10.0, unit)
    elif query.get('use_imperial', False):
        unit = ' mph'
        wind = u'%s%s%s' % (wind_direction, data['windspeedMiles'], unit)
    else:
        unit = ' km/h'
        wind = u'%s%s%s' % (wind_direction, data['windspeedKmph'], unit)

    return wind

def render_location(data, query):
    """
    location (l)
    """

    return (data['override_location'] or data['location']) # .title()

def render_moonphase(_, query):
    """
    A symbol describing the phase of the moon
    """
    astral = Astral()
    moon_index = int(
        int(32.0*astral.moon_phase(date=datetime.datetime.today())/28+2)%32/4
    )
    return MOON_PHASES[moon_index]

def render_moonday(_, query):
    """
    An number describing the phase of the moon (days after the New Moon)
    """
    astral = Astral()
    return str(int(astral.moon_phase(date=datetime.datetime.today())))

def render_sunset(data, query):
    location = data['location']
    city_name = location
    astral = Astral()
    location = Location(('Nuremberg', 'Germany',
              49.453872, 11.077298, 'Europe/Berlin', 0))
    sun = location.sun(date=datetime.datetime.today(), local=True)


    return str(sun['sunset'])

FORMAT_SYMBOL = {
    'c':    render_condition,
    'C':    render_condition_fullname,
    'h':    render_humidity,
    't':    render_temperature,
    'w':    render_wind,
    'l':    render_location,
    'm':    render_moonphase,
    'M':    render_moonday,
    's':    render_sunset,
    }

def render_line(line, data, query):
    """
    Render format `line` using `data`
    """

    def render_symbol(match):
        """
        Render one format symbol from re `match`
        using `data` from external scope.
        """

        symbol_string = match.group(0)
        symbol = symbol_string[-1]

        if symbol not in FORMAT_SYMBOL:
            return ''

        render_function = FORMAT_SYMBOL[symbol]
        return render_function(data, query)

    return re.sub(r'%[^%]*[a-zA-Z]', render_symbol, line)

def format_weather_data(format_line, location, override_location, data, query):
    """
    Format information about current weather `data` for `location`
    with specified in `format_line` format
    """

    if 'data' not in data:
        return 'Unknow location; please try ~%s' % location
    current_condition = data['data']['current_condition'][0]
    current_condition['location'] = location
    current_condition['override_location'] = override_location
    output = render_line(format_line, current_condition, query)
    return output

def wttr_line(location, override_location_name, query, lang):
    """
    Return 1line weather information for `location`
    in format `line_format`
    """

    format_line = query.get('format', '')

    if format_line in PRECONFIGURED_FORMAT:
        format_line = PRECONFIGURED_FORMAT[format_line]

    weather_data = get_weather_data(location, lang)

    output = format_weather_data(format_line, location, override_location_name, weather_data, query)
    output = output.rstrip("\n")+"\n"
    return output

def main():
    """
    Function for standalone module usage
    """

    location = sys.argv[1]
    query = {
        'line': sys.argv[2],
        }

    sys.stdout.write(wttr_line(location, location, query, 'en'))

if __name__ == '__main__':
    main()
