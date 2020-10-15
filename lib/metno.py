#!/bin/env python
# vim: fileencoding=utf-8
from datetime import datetime, timedelta
import json
import logging
import os
import re
import sys

import timezonefinder
from pytz import timezone

from constants import WWO_CODE

logging.basicConfig(level=os.environ.get("LOGLEVEL", "INFO"))
logger = logging.getLogger(__name__)


def metno_request(path, query_string):
    # We'll need to sanitize the inbound request - ideally the
    # premium/v1/weather.ashx portion would have always been here, though
    # it seems as though the proxy was built after the majority of the app
    # and not refactored. For WAPI we'll strip this and the API key out,
    # then manage it on our own.
    logger.debug('Original path: ' + path)
    logger.debug('Original query: ' + query_string)

    path = path.replace('premium/v1/weather.ashx',
                        'weatherapi/locationforecast/2.0/complete')
    query_string = re.sub(r'key=[^&]*&', '', query_string)
    query_string = re.sub(r'format=[^&]*&', '', query_string)
    days = int(re.search(r'num_of_days=([0-9]+)&', query_string).group(1))
    query_string = re.sub(r'num_of_days=[0-9]+&', '', query_string)
    # query_string = query_string.replace('key=', '?key=' + WAPI_KEY)
    # TP is for hourly forecasting, which isn't available in the free api.
    query_string = re.sub(r'tp=[0-9]*&', '', query_string)
    # This assumes lang=... is at the end. Also note that the API doesn't
    # localize, and we're not either. TODO: add language support
    query_string = re.sub(r'lang=[^&]*$', '', query_string)
    query_string = re.sub(r'&$', '', query_string)

    logger.debug('qs: ' + query_string)
    # Deal with coordinates. Need to be rounded to 4 decimals for metno ToC
    # and in a different query string format
    coord_match = re.search(r'q=[^&]*', query_string)
    coords_str = coord_match.group(0)
    coords = re.findall(r'[-0-9.]+', coords_str)
    lat = str(round(float(coords[0]), 4))
    lng = str(round(float(coords[1]), 4))
    logger.debug('lat: ' + lat)
    logger.debug('lng: ' + lng)
    query_string = re.sub(r'q=[^&]*', 'lat=' + lat + '&lon=' + lng + '&',
                          query_string)
    logger.debug('Return path: ' + path)
    logger.debug('Return query: ' + query_string)

    return path, query_string, days


def celsius_to_f(celsius):
    return round((1.8 * celsius) + 32, 1)


def to_weather_code(symbol_code):
    logger.debug(symbol_code)
    code = re.sub(r'_.*', '', symbol_code)
    logger.debug(code)
    # symbol codes: https://api.met.no/weatherapi/weathericon/2.0/documentation
    # they also have _day, _night and _polartwilight variants
    # See json from https://api.met.no/weatherapi/weathericon/2.0/legends
    # WWO codes: https://github.com/chubin/wttr.in/blob/master/lib/constants.py
    #            http://www.worldweatheronline.com/feed/wwoConditionCodes.txt
    weather_code_map = {
            "clearsky": 113,
            "cloudy": 119,
            "fair": 116,
            "fog": 143,
            "heavyrain": 302,
            "heavyrainandthunder": 389,
            "heavyrainshowers": 305,
            "heavyrainshowersandthunder": 386,
            "heavysleet": 314,  # There's a ton of 'LightSleet' in WWO_CODE...
            "heavysleetandthunder": 377,
            "heavysleetshowers": 362,
            "heavysleetshowersandthunder": 374,
            "heavysnow": 230,
            "heavysnowandthunder": 392,
            "heavysnowshowers": 371,
            "heavysnowshowersandthunder": 392,
            "lightrain": 266,
            "lightrainandthunder": 200,
            "lightrainshowers": 176,
            "lightrainshowersandthunder": 386,
            "lightsleet": 281,
            "lightsleetandthunder": 377,
            "lightsleetshowers": 284,
            "lightsnow": 320,
            "lightsnowandthunder": 392,
            "lightsnowshowers": 368,
            "lightssleetshowersandthunder": 365,
            "lightssnowshowersandthunder": 392,
            "partlycloudy": 116,
            "rain": 293,
            "rainandthunder": 389,
            "rainshowers": 299,
            "rainshowersandthunder": 386,
            "sleet": 185,
            "sleetandthunder": 392,
            "sleetshowers": 263,
            "sleetshowersandthunder": 392,
            "snow": 329,
            "snowandthunder": 392,
            "snowshowers": 230,
            "snowshowersandthunder": 392,
    }
    if code not in weather_code_map:
        logger.debug('not found')
        return -1  # not found
    logger.debug(weather_code_map[code])
    return weather_code_map[code]


def to_description(symbol_code):
    desc = WWO_CODE[str(to_weather_code(symbol_code))]
    logger.debug(desc)
    return desc


def to_16_point(degrees):
    # 360 degrees / 16 = 22.5 degrees of arc or 11.25 degrees around the point
    if degrees > (360 - 11.25) or degrees <= 11.25:
        return 'N'
    if degrees > 11.25 and degrees <= (11.25 + 22.5):
        return 'NNE'
    if degrees > (11.25 + (22.5 * 1)) and degrees <= (11.25 + (22.5 * 2)):
        return 'NE'
    if degrees > (11.25 + (22.5 * 2)) and degrees <= (11.25 + (22.5 * 3)):
        return 'ENE'
    if degrees > (11.25 + (22.5 * 3)) and degrees <= (11.25 + (22.5 * 4)):
        return 'E'
    if degrees > (11.25 + (22.5 * 4)) and degrees <= (11.25 + (22.5 * 5)):
        return 'ESE'
    if degrees > (11.25 + (22.5 * 5)) and degrees <= (11.25 + (22.5 * 6)):
        return 'SE'
    if degrees > (11.25 + (22.5 * 6)) and degrees <= (11.25 + (22.5 * 7)):
        return 'SSE'
    if degrees > (11.25 + (22.5 * 7)) and degrees <= (11.25 + (22.5 * 8)):
        return 'S'
    if degrees > (11.25 + (22.5 * 8)) and degrees <= (11.25 + (22.5 * 9)):
        return 'SSW'
    if degrees > (11.25 + (22.5 * 9)) and degrees <= (11.25 + (22.5 * 10)):
        return 'SW'
    if degrees > (11.25 + (22.5 * 10)) and degrees <= (11.25 + (22.5 * 11)):
        return 'WSW'
    if degrees > (11.25 + (22.5 * 11)) and degrees <= (11.25 + (22.5 * 12)):
        return 'W'
    if degrees > (11.25 + (22.5 * 12)) and degrees <= (11.25 + (22.5 * 13)):
        return 'WNW'
    if degrees > (11.25 + (22.5 * 13)) and degrees <= (11.25 + (22.5 * 14)):
        return 'NW'
    if degrees > (11.25 + (22.5 * 14)) and degrees <= (11.25 + (22.5 * 15)):
        return 'NNW'


def meters_to_miles(meters):
    return round(meters * 0.00062137, 2)


def mm_to_inches(mm):
    return round(mm / 25.4, 2)


def hpa_to_mb(hpa):
    return hpa


def hpa_to_in(hpa):
    return round(hpa * 0.02953, 2)


def group_hours_to_days(lat, lng, hourlies, days_to_return):
    tf = timezonefinder.TimezoneFinder()
    timezone_str = tf.certain_timezone_at(lat=lat, lng=lng)
    logger.debug('got TZ: ' + timezone_str)
    tz = timezone(timezone_str)
    start_day_gmt = datetime.fromisoformat(hourlies[0]['time']
                                           .replace('Z', '+00:00'))
    start_day_local = start_day_gmt.astimezone(tz)
    end_day_local = (start_day_local + timedelta(days=days_to_return - 1)).date()
    logger.debug('series starts at gmt time: ' + str(start_day_gmt))
    logger.debug('series starts at local time: ' + str(start_day_local))
    logger.debug('series ends on day: ' + str(end_day_local))
    days = {}

    for hour in hourlies:
        current_day_gmt = datetime.fromisoformat(hour['time']
                                                 .replace('Z', '+00:00'))
        current_local = current_day_gmt.astimezone(tz)
        current_day_local = current_local.date()
        if current_day_local > end_day_local:
            continue
        if current_day_local not in days:
            days[current_day_local] = {'hourly': []}
        hour['localtime'] = current_local.time()
        days[current_day_local]['hourly'].append(hour)

    # Need a second pass to build the min/max/avg data
    for date, day in days.items():
        minTempC = -999
        maxTempC = 1000
        avgTempC = None
        n = 0
        maxUvIndex = 0
        for hour in day['hourly']:
            temp = hour['data']['instant']['details']['air_temperature']
            if temp > minTempC:
                minTempC = temp
            if temp < maxTempC:
                maxTempC = temp
            if avgTempC is None:
                avgTempC = temp
                n = 1
            else:
                avgTempC = ((avgTempC * n) + temp) / (n + 1)
                n = n + 1

            uv = hour['data']['instant']['details']
            if 'ultraviolet_index_clear_sky' in uv:
                if uv['ultraviolet_index_clear_sky'] > maxUvIndex:
                    maxUvIndex = uv['ultraviolet_index_clear_sky']
        day["maxtempC"] = str(maxTempC)
        day["maxtempF"] = str(celsius_to_f(maxTempC))
        day["mintempC"] = str(minTempC)
        day["mintempF"] = str(celsius_to_f(minTempC))
        day["avgtempC"] = str(round(avgTempC, 1))
        day["avgtempF"] = str(celsius_to_f(avgTempC))
        # day["totalSnow_cm": "not implemented",
        # day["sunHour": "12",  # This would come from astonomy data
        day["uvIndex"] = str(maxUvIndex)

    return days

def _convert_hour(hour):
    # Whatever is upstream is expecting data in the shape of WWO. This method will
    # morph from metno to hourly WWO response format.
    # Note that WWO is providing data every 3 hours. Metno provides every hour
    #   {
    # "time": "0",
    # "tempC": "19",
    # "tempF": "66",
    # "windspeedMiles": "6",
    # "windspeedKmph": "9",
    # "winddirDegree": "276",
    # "winddir16Point": "W",
    # "weatherCode": "119",
    # "weatherIconUrl": [
    #   {
    #     "value": "http://cdn.worldweatheronline.com/images/wsymbols01_png_64/wsymbol_0003_white_cloud.png"
    #   }
    # ],
    # "weatherDesc": [
    #   {
    #     "value": "Cloudy"
    #   }
    # ],
    # "precipMM": "0.0",
    # "precipInches": "0.0",
    # "humidity": "62",
    # "visibility": "10",
    # "visibilityMiles": "6",
    # "pressure": "1017",
    # "pressureInches": "31",
    # "cloudcover": "66",
    # "HeatIndexC": "19",
    # "HeatIndexF": "66",
    # "DewPointC": "12",
    # "DewPointF": "53",
    # "WindChillC": "19",
    # "WindChillF": "66",
    # "WindGustMiles": "8",
    # "WindGustKmph": "13",
    # "FeelsLikeC": "19",
    # "FeelsLikeF": "66",
    # "chanceofrain": "0",
    # "chanceofremdry": "93",
    # "chanceofwindy": "0",
    # "chanceofovercast": "89",
    # "chanceofsunshine": "18",
    # "chanceoffrost": "0",
    # "chanceofhightemp": "0",
    # "chanceoffog": "0",
    # "chanceofsnow": "0",
    # "chanceofthunder": "0",
    # "uvIndex": "1"
    details = hour['data']['instant']['details']
    if 'next_1_hours' in hour['data']:
        next_hour = hour['data']['next_1_hours']
    elif 'next_6_hours' in hour['data']:
        next_hour = hour['data']['next_6_hours']
    elif 'next_12_hours' in hour['data']:
        next_hour = hour['data']['next_12_hours']
    else:
        next_hour = {}

    # Need to dig out symbol_code and precipitation_amount
    symbol_code = 'clearsky_day'  # Default to sunny
    if 'summary' in next_hour and 'symbol_code' in next_hour['summary']:
        symbol_code = next_hour['summary']['symbol_code']
    precipitation_amount = 0  # Default to no rain
    if 'details' in next_hour and 'precipitation_amount' in next_hour['details']:
        precipitation_amount = next_hour['details']['precipitation_amount']

    uvIndex = 0  # default to 0 index
    if 'ultraviolet_index_clear_sky' in details:
        uvIndex = details['ultraviolet_index_clear_sky']
    localtime = ''
    if 'localtime' in hour:
        localtime = "{h:02.0f}".format(h=hour['localtime'].hour) + \
                    "{m:02.0f}".format(m=hour['localtime'].minute)
        logger.debug(str(hour['localtime']))
    # time property is local time, 4 digit 24 hour, with no :, e.g. 2100
    return {
        'time': localtime,
        'observation_time': hour['time'],  # Need to figure out WWO TZ
        # temp_C is used in we-lang.go calcs in such a way
        # as to expect a whole number
        'temp_C': str(int(round(details['air_temperature'], 0))),
        # temp_F can be more precise - not used in we-lang.go calcs
        'temp_F': str(celsius_to_f(details['air_temperature'])),
        'weatherCode': str(to_weather_code(symbol_code)),
        'weatherIconUrl': [{
            'value': 'not yet implemented',
        }],
        'weatherDesc': [{
            'value': to_description(symbol_code),
        }],
        # similiarly, windspeedMiles is not used by we-lang.go, but kmph is
        "windspeedMiles": str(meters_to_miles(details['wind_speed'])),
        "windspeedKmph": str(int(round(details['wind_speed'], 0))),
        "winddirDegree": str(details['wind_from_direction']),
        "winddir16Point": to_16_point(details['wind_from_direction']),
        "precipMM": str(precipitation_amount),
        "precipInches": str(mm_to_inches(precipitation_amount)),
        "humidity": str(details['relative_humidity']),
        "visibility": 'not yet implemented',  # str(details['vis_km']),
        "visibilityMiles": 'not yet implemented',  # str(details['vis_miles']),
        "pressure": str(hpa_to_mb(details['air_pressure_at_sea_level'])),
        "pressureInches": str(hpa_to_in(details['air_pressure_at_sea_level'])),
        "cloudcover": 'not yet implemented',  # Convert from cloud_area_fraction?? str(details['cloud']),
        # metno doesn't have FeelsLikeC, but we-lang.go is using it in calcs,
        # so we shall set it to temp_C
        "FeelsLikeC": str(int(round(details['air_temperature'], 0))),
        "FeelsLikeF": 'not yet implemented',  # str(details['feelslike_f']),
        "uvIndex": str(uvIndex),
    }


def _convert_hourly(hours):
    converted_hours = []
    for hour in hours:
        converted_hours.append(_convert_hour(hour))
    return converted_hours


# Whatever is upstream is expecting data in the shape of WWO. This method will
# morph from metno to WWO response format.
def create_standard_json_from_metno(content, days_to_return):
    try:
        forecast = json.loads(content)         # pylint: disable=invalid-name
    except (ValueError, TypeError) as exception:
        logger.error("---")
        logger.error(exception)
        logger.error("---")
        return {}, ''
    hourlies = forecast['properties']['timeseries']
    current = hourlies[0]
    # We are assuming these units:
    # "units": {
    #   "air_pressure_at_sea_level": "hPa",
    #   "air_temperature": "celsius",
    #   "air_temperature_max": "celsius",
    #   "air_temperature_min": "celsius",
    #   "cloud_area_fraction": "%",
    #   "cloud_area_fraction_high": "%",
    #   "cloud_area_fraction_low": "%",
    #   "cloud_area_fraction_medium": "%",
    #   "dew_point_temperature": "celsius",
    #   "fog_area_fraction": "%",
    #   "precipitation_amount": "mm",
    #   "relative_humidity": "%",
    #   "ultraviolet_index_clear_sky": "1",
    #   "wind_from_direction": "degrees",
    #   "wind_speed": "m/s"
    # }
    content = {
        'data': {
            'request': [{
                'type': 'feature',
                'query': str(forecast['geometry']['coordinates'][1]) + ',' +
                        str(forecast['geometry']['coordinates'][0])
            }],
            'current_condition': [
                _convert_hour(current)
            ],
            'weather': []
        }
    }

    days = group_hours_to_days(forecast['geometry']['coordinates'][1],
                               forecast['geometry']['coordinates'][0],
                               hourlies, days_to_return)

    # TODO: Astronomy needs to come from this:
    # https://api.met.no/weatherapi/sunrise/2.0/.json?lat=40.7127&lon=-74.0059&date=2020-10-07&offset=-05:00
    # and obviously can be cached for a while
    # https://api.met.no/weatherapi/sunrise/2.0/documentation
    # Note that full moon/new moon/first quarter/last quarter aren't returned
    # and the moonphase value should match these from WWO:
    # New Moon
    # Waxing Crescent
    # First Quarter
    # Waxing Gibbous
    # Full Moon
    # Waning Gibbous
    # Last Quarter
    # Waning Crescent

    for date, day in days.items():
        content['data']['weather'].append({
                "date": str(date),
                "astronomy": [],
                "maxtempC": day['maxtempC'],
                "maxtempF": day['maxtempF'],
                "mintempC": day['mintempC'],
                "mintempF": day['mintempF'],
                "avgtempC": day['avgtempC'],
                "avgtempF": day['avgtempF'],
                "totalSnow_cm": "not implemented",
                "sunHour": "12",  # This would come from astonomy data
                "uvIndex": day['uvIndex'],
                'hourly': _convert_hourly(day['hourly']),
        })

    # for day in forecast.
    return json.dumps(content)


if __name__ == "__main__":
    # if len(sys.argv) == 1:
    #     for deg in range(0, 360):
    #         print('deg: ' + str(deg) + '; 16point: ' + to_16_point(deg))
    if len(sys.argv) == 2:
        req = sys.argv[1].split('?')
        # to_description(sys.argv[1])
        metno_request(req[0], req[1])
    elif len(sys.argv) == 3:
        with open(sys.argv[1], 'r') as contentf:
            content = create_standard_json_from_metno(contentf.read(),
                                                      int(sys.argv[2]))
        print(content)
    else:
        print('usage: metno <content file> <days>')
