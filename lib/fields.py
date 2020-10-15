"""
Human readable description of the available data fields
describing current weather, weather forecast, and astronomical data
"""

DESCRIPTION = {
    # current condition fields
    "FeelsLikeC": (
        "Feels Like Temperature in Celsius",
        "temperature_feels_like_celsius"),
    "FeelsLikeF": (
        "Feels Like Temperature in Fahrenheit",
        "temperature_feels_like_fahrenheit"),
    "cloudcover": (
        "Cloud Coverage in Percent",
        "cloudcover_percentage"),
    "humidity": (
        "Humidity in Percent",
        "humidity_percentage"),
    "precipMM": (
        "Precipitation (Rainfall) in mm",
        "precipitation_mm"),
    "pressure": (
        "Air pressure in hPa",
        "pressure_hpa"),
    "temp_C": (
        "Temperature in Celsius",
        "temperature_celsius"),
    "temp_F": (
        "Temperature in Fahrenheit",
        "temperature_fahrenheit"),
    "uvIndex": (
        "Ultaviolet Radiation Index",
        "uv_index"),
    "visibility": (
        "Visible Distance in Kilometres",
        "visibility"),
    "weatherCode": (
        "Code to describe Weather Condition",
        "weather_code"),
    "winddirDegree": (
        "Wind Direction in Degree",
        "winddir_degree"),
    "windspeedKmph": (
        "Wind Speed in Kilometres per Hour",
        "windspeed_kmph"),
    "windspeedMiles": (
        "Wind Speed in Miles per Hour",
        "windspeed_mph"),
    "observation_time": (
        "Minutes since start of the day the observation happened",
        "observation_time"),

    # fields with `description`
    "weatherDesc": (
        "Weather Description",
        "weather_desc"),
    "winddir16Point": (
        "Wind Direction on a 16-wind compass rose",
        "winddir_16_point"),

    # forecast fields
    "maxtempC": (
        "Maximum Temperature in Celsius",
        "temperature_celsius_maximum"),
    "maxtempF": (
        "Maximum Temperature in Fahrenheit",
        "temperature_fahrenheit_maximum"),
    "mintempC": (
        "Minimum Temperature in Celsius",
        "temperature_celsius_minimum"),
    "mintempF": (
        "Minimum Temperature in Fahrenheit",
        "temperature_fahrenheit_minimum"),
    "sunHour":(
        "Hours of sunlight",
        "sun_hour"),
    "totalSnow_cm":(
        "Total snowfall in cm",
        "snowfall_cm"),

    # astronomy fields
    "moon_illumination": (
        "Percentage of the moon illuminated",
        "astronomy_moon_illumination"),

    # astronomy fields with description
    "moon_phase": (
        "Phase of the moon",
        "astronomy_moon_phase"),

    # astronomy fields with time
    "moonrise": (
        "Minutes since start of the day untill the moon appears above the horizon",
        "astronomy_moonrise_min"),
    "moonset": (
        "Minutes since start of the day untill the moon disappears below the horizon",
        "astronomy_moonset_min"),
    "sunrise": (
        "Minutes since start of the day untill the sun appears above the horizon",
        "astronomy_sunrise_min"),
    "sunset": (
        "Minutes since start of the day untill the moon disappears below the horizon",
        "astronomy_sunset_min"),
    }
