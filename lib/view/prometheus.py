"""
Rendering weather data in the Prometheus format.

"""

from datetime import datetime

from fields import DESCRIPTION

def _render_current(data, for_day="current"):
    "Converts data into prometheus style format"

    output = []

    for field_name, val in DESCRIPTION.items():

        help, name = val

        try:
            value = data[field_name]
            if field_name == "weatherDesc":
                value = value[0]["value"]
        except (IndexError, KeyError):
            try:
                value = data["astronomy"][0][field_name]
                if value.endswith(" AM") or value.endswith(" PM"):
                    value = _convert_time_to_minutes(value)
            except (IndexError, KeyError, ValueError):
                continue

        try:
            if name == "observation_time":
                value = _convert_time_to_minutes(value)
        except ValueError:
            continue

        description = ""
        try:
            float(value)
        except ValueError:
            description = f", description=\"{value}\""
            value = "1"

        output.append(f"# HELP {name} {help}\n"
                      f"{name}{{forecast=\"{for_day}\"{description}}} {value}")

    return "\n".join(output)+"\n"

def _convert_time_to_minutes(time_str):
    "Convert time from midnight to minutes"
    return int((datetime.strptime(time_str, "%I:%M %p")
        - datetime.strptime("12:00 AM", "%I:%M %p")).total_seconds())//60

def render_prometheus(data):
    """
    Convert `data` into Prometheus format
    and return it as string.
    """

    answer = _render_current(data["current_condition"][0])
    for i in range(3):
        data2 = data["weather"][i]
        answer += _render_current(data2, for_day="%sd" % i)
    return answer
