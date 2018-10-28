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
from constants import WWO_CODE, WEATHER_SYMBOL, WIND_DIRECTION
from weather_data import get_weather_data

PRECONFIGURED_FORMAT = {
    '1':    u'%c %t',
    '2':    u'%c 🌡️%t 🌬️%w',
    '3':    u'%l: %c %t',
    '4':    u'%l: %c 🌡️%t 🌬️%w',
}

def render_temperature(data):
    """
    temperature (t)
    """

    temperature = u'%s⁰C' % data['temp_C']
    if temperature[0] != '-':
        temperature = '+' + temperature
    return temperature

def render_condition(data):
    """
    condition (c)
    """

    weather_condition = WEATHER_SYMBOL[WWO_CODE[data['weatherCode']]]
    return weather_condition

def render_wind(data):
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

    unit = ' km/h'
    wind = u'%s%s%s' % (wind_direction, data['windspeedKmph'], unit)
    return wind

def render_location(data):
    """
    location (l)
    """

    return data['location']

FORMAT_SYMBOL = {
    'c':    render_condition,
    't':    render_temperature,
    'w':    render_wind,
    'l':    render_location,
    }

def render_line(line, data):
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
        return render_function(data)

    return re.sub(r'%[^%]*[a-z]', render_symbol, line)

def format_weather_data(format_line, location, data):
    """
    Format information about current weather `data` for `location`
    with specified in `format_line` format
    """

    current_condition = data['data']['current_condition'][0]
    current_condition['location'] = location
    output = render_line(format_line, current_condition)
    return output

def wttr_line(location, query):
    """
    Return 1line weather information for `location`
    in format `line_format`
    """

    format_line = query.get('line', '')

    if format_line in PRECONFIGURED_FORMAT:
        format_line = PRECONFIGURED_FORMAT[format_line]

    weather_data = get_weather_data(location)

    output = format_weather_data(format_line, location, weather_data)
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

    sys.stdout.write(wttr_line(location, query))

if __name__ == '__main__':
    main()
