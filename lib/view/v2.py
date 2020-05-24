# vim: fileencoding=utf-8
# vim: foldmethod=marker foldenable:

"""
[X] emoji
[ ] wego icon
[ ] v2.wttr.in
[X] astronomical (sunset)
[X] time
[X] frames
[X] colorize rain data
[ ] date + locales
[X] wind color
[ ] highlight current date
[ ] bind to real site
[ ] max values: temperature
[X] max value: rain
[ ] comment github
[ ] commit

"""

import sys

import re
import math
import json
import datetime
import io

import requests
import diagram
import pyjq
import pytz
import numpy as np
from astral import LocationInfo
from astral import moon, sun
from scipy.interpolate import interp1d
from babel.dates import format_datetime

from globals import WWO_KEY, remove_ansi
import constants
import translations
import parse_query
from . import line as wttr_line

if not sys.version_info >= (3, 0):
    reload(sys) # noqa: F821
    sys.setdefaultencoding("utf-8")

# data processing {{{

def get_data(config):
    """
    Fetch data for `query_string`
    """

    url = (
        'http://'
        'localhost:5001/premium/v1/weather.ashx'
        '?key=%s'
        '&q=%s&format=json&num_of_days=3&tp=3&lang=None'
    ) % (WWO_KEY, config["location"])
    text = requests.get(url).text
    parsed_data = json.loads(text)
    return parsed_data

def interpolate_data(input_data, max_width):
    """
    Resample `input_data` to number of `max_width` counts
    """

    input_data = list(input_data)
    input_data_len = len(input_data)
    x = list(range(input_data_len))
    y = input_data
    xvals = np.linspace(0, input_data_len-1, max_width)
    yinterp = interp1d(x, y, kind='cubic')
    return yinterp(xvals)

def jq_query(query, data_parsed):
    """
    Apply `query` to structued data `data_parsed`
    """

    pyjq_data = pyjq.all(query, data_parsed)
    data = list(map(float, pyjq_data))
    return data

# }}}
# utils {{{
def colorize(string, color_code, html_output=False):
    if html_output:
        return "<font color='#777777'>%s</font>" % (string)
    else:
        return "\033[%sm%s\033[0m" % (color_code, string)
# }}}
# draw_spark {{{


def draw_spark(data, height, width, color_data):
    """
    Spark-style visualize `data` in a region `height` x `width`
    """

    _BARS = u' _▁▂▃▄▅▇█'

    def _box(height, row, value, max_value):
        row_height = 1.0 * max_value / height
        if row_height * row >= value:
            return _BARS[0]
        if row_height * (row+1) <= value:
            return _BARS[-1]

        return _BARS[int(1.0*(value - row_height*row)/(row_height*1.0)*len(_BARS))]

    max_value = max(data)

    output = ""
    color_code = 20
    for i in range(height):
        for j in range(width):
            character = _box(height, height-i-1, data[j], max_value)
            if data[j] != 0:
                chance_of_rain = color_data[j]/100.0 * 2
                if chance_of_rain > 1:
                    chance_of_rain = 1
                color_index = int(5*chance_of_rain)
                color_code = 16 + color_index # int(math.floor((20-16) * 1.0 * (height-1-i)/height*(max_value/data[j])))
            output += "\033[38;5;%sm%s\033[0m" % (color_code, character)
        output += "\n"

    # labeling max value
    if max_value == 0:
        max_line = " "*width
    else:
        max_line = ""
        for j in range(width):
            if data[j] == max_value:
                max_line = "%3.2fmm|%s%%" % (max_value, int(color_data[j]))
                orig_max_line = max_line

                # aligning it
                if len(max_line)//2 < j and len(max_line)//2 + j < width:
                    spaces = " "*(j - len(max_line)//2)
                    max_line = spaces + max_line # + spaces
                    max_line = max_line + " "*(width - len(max_line))
                elif len(max_line)//2 + j >= width:
                    max_line = " "*(width - len(max_line)) + max_line

                max_line = max_line.replace(orig_max_line, colorize(orig_max_line, "38;5;33"))

                break

    if max_line:
        output = "\n" + max_line + "\n" + output + "\n"

    return output

# }}}
# draw_diagram {{{
def draw_diagram(data, height, width):

    option = diagram.DOption()
    option.size = diagram.Point([width, height])
    option.mode = 'g'

    stream = io.BytesIO()
    gram = diagram.DGWrapper(
        data=[list(data), range(len(data))],
        dg_option=option,
        ostream=stream)
    gram.show()
    return stream.getvalue().decode("utf-8")
# }}}
# draw_date {{{


def draw_date(config, geo_data):
    """
    """

    tzinfo = pytz.timezone(geo_data["timezone"])

    locale = config.get("locale", "en_US")
    datetime_day_start = datetime.datetime.utcnow()

    answer = ""
    for day in range(3):
        datetime_ = datetime_day_start + datetime.timedelta(hours=24*day)
        date = format_datetime(datetime_, "EEE dd MMM", locale=locale, tzinfo=tzinfo)

        spaces = ((24-len(date))//2)*" "
        date = spaces + date + spaces
        date = " "*(24-len(date)) + date
        answer += date
    answer += "\n"

    for _ in range(3):
        answer += " "*23 + u"╷"
    return answer[:-1] + " "


# }}}
# draw_time {{{


def draw_time(geo_data):
    """
    """

    tzinfo = pytz.timezone(geo_data["timezone"])

    line = ["", ""]

    for _ in range(3):
        part = u"─"*5 + u"┴" + u"─"*5
        line[0] += part + u"┼" + part + u"╂"
    line[0] += "\n"

    for _ in range(3):
        line[1] += "     6    12    18      "
    line[1] += "\n"

    # highlight current time
    hour_number = \
        (datetime.datetime.now(tzinfo)
         - datetime.datetime.now(tzinfo).replace(hour=0, minute=0, second=0, microsecond=0)
        ).seconds//3600

    for line_number, _ in enumerate(line):
        line[line_number] = \
                line[line_number][:hour_number] \
                + colorize(line[line_number][hour_number], "46") \
                + line[line_number][hour_number+1:]

    return "".join(line)


# }}}
# draw_astronomical {{{
def draw_astronomical(city_name, geo_data):
    datetime_day_start = datetime.datetime.now().replace(hour=0, minute=0, second=0, microsecond=0)

    city = LocationInfo()
    city.latitude = geo_data["latitude"]
    city.longitude = geo_data["longitude"]
    city.timezone = geo_data["timezone"]

    answer = ""
    moon_line = ""
    for time_interval in range(72):

        current_date = (
            datetime_day_start
            + datetime.timedelta(hours=1*time_interval)).replace(tzinfo=pytz.timezone(geo_data["timezone"]))

        try:
            dawn = sun.dawn(city.observer, date=current_date)
        except ValueError:
            dawn = current_date

        try:
            dusk = sun.dusk(city.observer, date=current_date)
        except ValueError:
            dusk = current_date + datetime.timedelta(hours=24)

        try:
            sunrise = sun.sunrise(city.observer, date=current_date)
        except ValueError:
            sunrise = current_date

        try:
            sunset = sun.sunset(city.observer, date=current_date)
        except ValueError:
            sunset = current_date + datetime.timedelta(hours=24)

        char = "."
        if current_date < dawn:
            char = " "
        elif current_date > dusk:
            char = " "
        elif dawn <= current_date and current_date <= sunrise:
            char = u"─"
        elif sunset <= current_date and current_date <= dusk:
            char = u"─"
        elif sunrise <= current_date and current_date <= sunset:
            char = u"━"

        answer += char

        # moon
        if time_interval in [0,23,47,69]: # time_interval % 3 == 0:
            moon_phase = moon.phase(
                date=datetime_day_start + datetime.timedelta(hours=time_interval))
            moon_phase_emoji = constants.MOON_PHASES[
                int(math.floor(moon_phase*1.0/28.0*8+0.5)) % len(constants.MOON_PHASES)]
        #    if time_interval in [0, 24, 48, 69]:
            moon_line += moon_phase_emoji # + " "
        elif time_interval % 3 == 0:
            if time_interval not in [24,28]: #se:
                moon_line += "   "
            else:
                moon_line += " "


    answer = moon_line + "\n" + answer + "\n"
    answer += "\n"
    return answer
# }}}
# draw_emoji {{{
def draw_emoji(data):
    answer = ""
    for i in data:
        emoji = constants.WEATHER_SYMBOL.get(
            constants.WWO_CODE.get(
                str(int(i)), "Unknown"))
        space = " "*(3-constants.WEATHER_SYMBOL_WIDTH_VTE.get(emoji))
        answer += emoji + space
    answer += "\n"
    return answer
# }}}
# draw_wind {{{
def draw_wind(data, color_data):

    def _color_code_for_wind_speed(wind_speed):

        color_codes = [
            (3,  241),  # 82
            (6,  242),  # 118
            (9,  243),  # 154
            (12, 246),  # 190
            (15, 250),  # 226
            (19, 253),  # 220
            (23, 214),
            (27, 208),
            (31, 202),
            (-1, 196)
        ]

        for this_wind_speed, this_color_code in color_codes:
            if wind_speed <= this_wind_speed:
                return this_color_code
        return color_codes[-1][1]

    answer = ""
    answer_line2 = ""

    for j, degree in enumerate(data):

        degree = int(degree)
        if degree:
            wind_direction = constants.WIND_DIRECTION[((degree+22)%360)//45]
        else:
            wind_direction = ""

        color_code = "38;5;%s" % _color_code_for_wind_speed(int(color_data[j]))
        answer += " %s " % colorize(wind_direction, color_code)

        # wind_speed
        wind_speed = int(color_data[j])
        wind_speed_str = colorize(str(wind_speed), color_code)
        if wind_speed < 10:
            wind_speed_str = " " + wind_speed_str + " "
        elif wind_speed < 100:
            wind_speed_str = " " + wind_speed_str
        answer_line2 += wind_speed_str

    answer += "\n"
    answer += answer_line2 + "\n"
    return answer
# }}}
# panel implementation {{{

def add_frame(output, width, config):
    """
    Add frame arond `output` that has width `width`
    """

    empty_line = " "*width
    output = "\n".join(u"│"+(x or empty_line)+u"│" for x in output.splitlines()) + "\n"

    weather_report = \
        translations.CAPTION[config.get("lang") or  "en"] \
        + " " \
        + (config["override_location_name"] or config["location"])

    caption = u"┤ " + " " + weather_report + " " + u" ├"
    output = u"┌" + caption + u"─"*(width-len(caption)) + u"┐\n" \
                + output + \
             u"└" + u"─"*width + u"┘\n"

    return output

def generate_panel(data_parsed, geo_data, config):
    """
    """

    max_width = 72

    precip_mm_query = "[.data.weather[] | .hourly[]] | .[].precipMM"
    precip_chance_query = "[.data.weather[] | .hourly[]] | .[].chanceofrain"
    feels_like_query = "[.data.weather[] | .hourly[]] | .[].FeelsLikeC"
    weather_code_query = "[.data.weather[] | .hourly[]] | .[].weatherCode"
    wind_direction_query = "[.data.weather[] | .hourly[]] | .[].winddirDegree"
    wind_speed_query = "[.data.weather[] | .hourly[]] | .[].windspeedKmph"

    output = ""

    output += "\n\n"

    output += draw_date(config, geo_data)
    output += "\n"
    output += "\n"
    output += "\n"

    data = jq_query(feels_like_query, data_parsed)
    data_interpolated = interpolate_data(data, max_width)
    output += draw_diagram(data_interpolated, 10, max_width)

    output += "\n"

    output += draw_time(geo_data)

    data = jq_query(precip_mm_query, data_parsed)
    color_data = jq_query(precip_chance_query, data_parsed)
    data_interpolated = interpolate_data(data, max_width)
    color_data_interpolated = interpolate_data(color_data, max_width)
    output += draw_spark(data_interpolated, 5, max_width, color_data_interpolated)
    output += "\n"

    data = jq_query(weather_code_query, data_parsed)
    output += draw_emoji(data)

    data = jq_query(wind_direction_query, data_parsed)
    color_data = jq_query(wind_speed_query, data_parsed)
    output += draw_wind(data, color_data)
    output += "\n"

    output += draw_astronomical(config["location"], geo_data)
    output += "\n"

    output = add_frame(output, max_width, config)
    return output


# }}}
# textual information {{{
def textual_information(data_parsed, geo_data, config, html_output=False):
    """
    Add textual information about current weather and
    astronomical conditions
    """

    def _shorten_full_location(full_location, city_only=False):

        def _count_runes(string):
            return len(string.encode('utf-16-le')) // 2

        words = full_location.split(",")

        output = words[0]
        if city_only:
            return output

        for word in words[1:]:
            if _count_runes(output + "," + word) > 50:
                return output
            output += "," + word

        return output

    def _colorize(text, color):
        return colorize(text, color, html_output=html_output)

    city = LocationInfo()
    city.latitude = geo_data["latitude"]
    city.longitude = geo_data["longitude"]
    city.timezone = geo_data["timezone"]

    output = []
    timezone = city.timezone

    datetime_day_start = datetime.datetime.now()\
            .replace(hour=0, minute=0, second=0, microsecond=0)

    format_line = "%c %C, %t, %h, %w, %P"
    current_condition = data_parsed['data']['current_condition'][0]
    query = {}
    weather_line = wttr_line.render_line(format_line, current_condition, query)
    output.append('Weather: %s' % weather_line)

    output.append('Timezone: %s' % timezone)

    local_tz = pytz.timezone(timezone)

    def _get_local_time_of(what):
        _sun = {
            "dawn": sun.dawn,
            "sunrise": sun.sunrise,
            "noon": sun.noon,
            "sunset": sun.sunset,
            "dusk": sun.dusk,
            }[what]

        current_time_of_what = _sun(city.observer, date=datetime_day_start)
        return current_time_of_what\
                .replace(tzinfo=pytz.utc)\
                .astimezone(local_tz)\
                .strftime("%H:%M:%S")

    local_time_of = {}
    for what in ["dawn", "sunrise", "noon", "sunset", "dusk"]:
        try:
            local_time_of[what] = _get_local_time_of(what)
        except ValueError:
            local_time_of[what] = "-"*8

    tmp_output = []

    tmp_output.append('  Now:    %%{{NOW(%s)}}' % timezone)
    tmp_output.append('Dawn:    %s' % local_time_of["dawn"])
    tmp_output.append('Sunrise: %s' % local_time_of["sunrise"])
    tmp_output.append('  Zenith: %s     ' % local_time_of["noon"])
    tmp_output.append('Sunset:  %s' % local_time_of["sunset"])
    tmp_output.append('Dusk:    %s' % local_time_of["dusk"])

    tmp_output = [
        re.sub("^([A-Za-z]*:)", lambda m: _colorize(m.group(1), "2"), x)
        for x in tmp_output]

    output.append(
        "%20s" % tmp_output[0] \
        + " | %20s " % tmp_output[1] \
        + " | %20s" % tmp_output[2])
    output.append(
        "%20s" % tmp_output[3] \
        + " | %20s " % tmp_output[4] \
        + " | %20s" % tmp_output[5])

    city_only = False
    suffix = ""
    if "Simferopol" in timezone:
        city_only = True
        suffix = ", Крым"

    if config["full_address"]:
        output.append('Location: %s%s [%5.4f,%5.4f]' \
                % (
                    _shorten_full_location(config["full_address"], city_only=city_only),
                    suffix,
                    geo_data["latitude"],
                    geo_data["longitude"],
                ))

    output = [
        re.sub("^( *[A-Za-z]*:)", lambda m: _colorize(m.group(1), "2"),
               re.sub("^( +[A-Za-z]*:)", lambda m: _colorize(m.group(1), "2"),
                      re.sub(r"(\|)", lambda m: _colorize(m.group(1), "2"), x)))
        for x in output]

    return "".join("%s\n" % x for x in output)

# }}}
# get_geodata {{{
def get_geodata(location):
    text = requests.get("http://localhost:8004/%s" % location).text
    return json.loads(text)
# }}}

def main(query, parsed_query, data):
    parsed_query["locale"] = "en_US"

    location = parsed_query["location"]
    html_output = parsed_query["html_output"]

    geo_data = get_geodata(location)
    if data is None:
        data_parsed = get_data(parsed_query)
    else:
        data_parsed = data

    if html_output:
        parsed_query["text"] = "no"
        filename = "b_" + parse_query.serialize(parsed_query) + ".png"
        output = """
<html>
<head>
<title>Weather report for {orig_location}</title>
<link rel="stylesheet" type="text/css" href="/files/style.css" />
</head>
<body>
  <img src="/{filename}" width="592" height="532"/>
<pre>
{textual_information}
</pre>
</body>
</html>
""".format(
        filename=filename, orig_location=parsed_query["orig_location"],
        textual_information=textual_information(
            data_parsed, geo_data, parsed_query, html_output=True))
    else:
        output = generate_panel(data_parsed, geo_data, parsed_query)
        if query.get("text") != "no" and parsed_query.get("text") != "no":
            output += textual_information(data_parsed, geo_data, parsed_query)
        if parsed_query.get('no-terminal', False):
            output = remove_ansi(output)
    return output

if __name__ == '__main__':
    sys.stdout.write(main(sys.argv[1]))
