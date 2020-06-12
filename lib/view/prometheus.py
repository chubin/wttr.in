"""
Rendering weather data in the Prometheus format.

"""

from datetime import datetime


EXPORTED_FIELDS = {
    "FeelsLikeC":("Feels Like Temperature in Celsius", "temperature_feels_like_celsius"),
    "FeelsLikeF":("Feels Like Temperature in Fahrenheit", "temperature_feels_like_fahrenheit"),
    "cloudcover":("Cloud Coverage in Percent", "cloudcover_percentage"),
    "humidity":("Humidity in Percent", "humidity_percentage"),
    "precipMM":("Precipitation (Rainfall) in mm", "precipitation_mm"),
    "pressure":("Air pressure in hPa", "pressure_hpa"),
    "temp_C":("Temperature in Celsius", "temperature_celsius"),
    "temp_F":("Temperature in Fahrenheit", "temperature_fahrenheit"),
    "uvIndex":("Ultaviolet Radiation Index", "uv_index"),
    "visibility":("Visible Distance in Kilometres", "visibility"),
    "weatherCode":("Code to describe Weather Condition", "weather_code"),
    "winddirDegree":("Wind Direction in Degree", "winddir_degree"),
    "windspeedKmph":("Wind Speed in Kilometres per Hour", "windspeed_kmph"),
    "windspeedMiles":("Wind Speed in Miles per Hour", "windspeed_mph"),
    }

EXPORTED_FIELDS_DESC = {
    "weatherDesc":("Weather Description", "weather_desc"),
    "winddir16Point":("Wind Direction on a 16-wind compass rose", "winddir_16_point"),
    }

EXPORTED_FIELDS_CONV = {
    "observation_time":("Minutes since start of the day the observation happened", "obversation_time")
    }

EXPORTED_FIELDS_WEATHER = {
    "maxtempC":("Maximum Temperature in Celsius", "temperature_celsius_maximum"),
    "maxtempF":("Maximum Temperature in Fahrenheit", "temperature_fahrenheit_maximum"),
    "mintempC":("Minimum Temperature in Celsius", "temperature_celsius_minimum"),
    "mintempF":("Minimum Temperature in Fahrenheit", "temperature_fahrenheit_minimum"),
    "sunHour":("Hours of sunlight", "sun_hour"),
    "totalSnow_cm":("Total snowfall in cm", "snowfall_cm"),
    "uvIndex":("Ultaviolet Radiation Index", "uv_index"),
    }

EXPORTED_FIELDS_WEATHER_ASTRONOMY = {
    "moon_illumination":("Percentage of the moon illuminated", "astronomoy_moon_illumination"),
    }

EXPORTED_FIELDS_WEATHER_ASTRONOMY_DESC = {
    "moon_phase": ("Phase of the moon", "astronomoy_moon_phase"),
    }

EXPORTED_FIELDS_WEATHER_ASTRONOMY_CONV = {
    "moonrise":("Minutes since start of the day untill the moon appears above the horizon", "astronomoy_moonrise_min"),
    "moonset":("Minutes since start of the day untill the moon disappears below the horizon", "astronomoy_moonset_min"),
    "sunrise":("Minutes since start of the day untill the sun appears below the horizon", "astronomoy_sunrise_min"),
    "sunset":("Minutes since start of the day untill the moon disappears below the horizon", "astronomoy_sunset_min"),

    }

def _render_current(data):
    """
    Converts data into prometheus style format
    """

    output = []
    current_condition = data["current_condition"][0]
    for field in EXPORTED_FIELDS:
        try:
            output.append("# HELP %s %s\n%s{forecast=\"current\"} %s" %
                          (EXPORTED_FIELDS[field][1],
                           EXPORTED_FIELDS[field][0],
                           EXPORTED_FIELDS[field][1],
                           current_condition[field]))
        except IndexError:
            pass

    for field in EXPORTED_FIELDS_CONV:
        try:
            output.append("# HELP %s %s\n%s{forecast=\"current\"} %s" %
                          (EXPORTED_FIELDS[field][1],
                           EXPORTED_FIELDS[field][0],
                           EXPORTED_FIELDS[field][1],
                           current_condition[field]))
        except IndexError:
            pass

    for field in EXPORTED_FIELDS_DESC:
        try:
            output.append("# HELP %s %s\n%s{forecast=\"current\", description=\"%s\"} 1" %
                          (EXPORTED_FIELDS_DESC[field][1],
                           EXPORTED_FIELDS_DESC[field][0],
                           EXPORTED_FIELDS_DESC[field][1],
                           convert_time_to_minutes(current_condition[field])))
        except IndexError:
            pass

    weather = data["weather"]
    i = 0
    for day in weather:
        for field in EXPORTED_FIELDS_WEATHER:
            try:
                output.append("# HELP %s %s\n%s{forecast=\"%id\"} %s" %
                              (EXPORTED_FIELDS_WEATHER[field][1],
                               EXPORTED_FIELDS_WEATHER[field][0],
                               i,
                               EXPORTED_FIELDS_WEATHER[field][1],
                               day[field]))
            except IndexError:
                pass

        for field in EXPORTED_FIELDS_WEATHER_ASTRONOMY:
            try:
                output.append("# HELP %s %s\n%s{forecast=\"%dd\"} %s" %
                              (EXPORTED_FIELDS_WEATHER[field][1],
                               EXPORTED_FIELDS_WEATHER[field][0],
                               EXPORTED_FIELDS_WEATHER[field][1],
                               i,
                               day["astronomy"][field]))
            except IndexError:
                pass

        for field in EXPORTED_FIELDS_WEATHER_ASTRONOMY_CONV:
            try:
                output.append("# HELP %s %s\n%s{forecast=\"%dd\"} %s" %
                              (EXPORTED_FIELDS_WEATHER[field][1],
                               EXPORTED_FIELDS_WEATHER[field][0],
                               EXPORTED_FIELDS_WEATHER[field][1],
                               i,
                               convert_time_to_minutes(day["astronomy"][field])))
            except IndexError:
                pass

        for field in EXPORTED_FIELDS_WEATHER_ASTRONOMY_DESC:
            try:
                output.append("# HELP %s %s\n%s{forecast=\"%dd\", description=\"%s\"} 1" %
                              (EXPORTED_FIELDS_WEATHER[field][1],
                               EXPORTED_FIELDS_WEATHER[field][0],
                               EXPORTED_FIELDS_WEATHER[field][1],
                               i,
                               day["astronomy"][field]))
            except IndexError:
                pass
        i = i+1

    return "\n".join(output)+"\n"

def convert_time_to_minutes(time_str):
    """
    Convert time from midnight to minutes
    """
    return int((datetime.strptime(time_str, "%I:%M %p") - datetime.strptime("12:00 AM", "%I:%M %p")).total_seconds()/60)

def render_prometheus(data):
    """
    Convert `data` into Prometheus format
    and return it as string.
    """

    return _render_current(data)
