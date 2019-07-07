"""
Weather data source
"""

import json
import requests
from globals import WWO_KEY

def get_weather_data(location, lang):
    """
    Get weather data for `location`
    """
    key = WWO_KEY
    url = ('/premium/v1/weather.ashx'
           '?key=%s&q=%s&format=json'
           '&num_of_days=3&tp=3&lang=%s') % (key, location, lang)
    url = 'http://127.0.0.1:5001' + url

    response = requests.get(url, timeout=1)
    try:
        data = json.loads(response.content)
    except ValueError:
        data = {}
    return data
