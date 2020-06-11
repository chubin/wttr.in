"""
Rendering weather data in the Prometheus format.

"""


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
    "winddirDegree":("Wind Direction in Degree", "winddir_degree"),
    "windspeedKmph":("Wind Speed in Kilometres per Hour", "windspeed_kmph"),
    "windspeedMiles":("Wind Speed in Miles per Hour", "windspeed_mph"),
    }

EXPORTED_FIELDS_DESC = {
    "weatherDesc":("Weather Description", "weather_desc"),
    "winddir16Point":("Wind Direction on a 16-wind compass rose", "winddir_16_point"),
    }

def _render_current(data):

    output = []
    current_condition = data["current_condition"][0]
    for field in EXPORTED_FIELDS:
        try:
            output.append("# HELP %s %s\n%s{forecast=\"0h\"} %s" %
                          (EXPORTED_FIELDS[field][1],
                           EXPORTED_FIELDS[field][0],
                           EXPORTED_FIELDS[field][1],
                           current_condition[field]))
        except IndexError:
            pass
    for field in EXPORTED_FIELDS_DESC:
        try:
            output.append("# HELP %s %s\n%s{forecast=\"0h\", description=\"%s\"} 1" %
                          (EXPORTED_FIELDS_DESC[field][1],
                           EXPORTED_FIELDS_DESC[field][0],
                           EXPORTED_FIELDS_DESC[field][1],
                           current_condition[field]))
        except IndexError:
            pass

    return "\n".join(output)+"\n"

def render_prometheus(data):
    """
    Convert `data` into Prometheus format
    and return it as string.
    """

    return _render_current(data)
