"""
Rendering weather data in the Prometheus format.

"""

EXPORTED_FIELDS = [
    "FeelsLikeC", "FeelsLikeF", "cloudcover", "humidity",
    "precipMM", "pressure", "temp_C", "temp_F", "uvIndex",
    "visibility", "winddirDegree", "windspeedKmph",
    "windspeedMiles",
    ]

def _render_current(data):

    output = []
    current_condition = data["current_condition"][0]
    for field in EXPORTED_FIELDS:
        try:
            output.append("%s{forecast=\"0h\"} %s" % (field, current_condition[field]))
        except IndexError:
            pass

    try:
        weather_desc = current_condition["weatherDesc"][0]["value"]
        output.append("weatherDesc{forecast=\"0h\", description=\"%s\"} 1" % weather_desc)
    except IndexError:
        pass

    try:
        winddir16point = current_condition["winddir16Point"]
        output.append("winddir16Point{forecast=\"0h\", description=\"%s\"} 1" % winddir16point)
    except IndexError:
        pass

    return "\n".join(output)+"\n"

def render_prometheus(data):
    """
    Convert `data` into Prometheus format
    and return it as string.
    """

    return _render_current(data)
