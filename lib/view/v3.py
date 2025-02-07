import os
import time
import random
import shutil

import requests

from globals import remove_ansi, MYDIR

from . import v2

CACHE_DIR = "/wttr.in/cache/v3"

def _get_cache_filename(filetype, name, timestamp=True):
    if name.endswith(".%s" % filetype):
        name = name[:len(filetype)+1]

    if not timestamp:
        return os.path.join(CACHE_DIR, filetype, name+"."+filetype)
    timestamp = str(int(time.time()/3600))
    return os.path.join(CACHE_DIR, filetype, timestamp, name+"."+filetype)

def _process_request(filetype, location, filename):
    """Try to fetch response from one of the upstream workers.
    If response code is 200, save the fetched response to the cache file.
    If no reponse returned, don't save anything, so the file will not be created.
    """

    ports = [9500, 9501, 9502]
    random.shuffle(ports)
    for port in ports:
        url = "http://127.0.0.1:%s/%s/%s" % (port, filetype, location)
        try:
            response = requests.get(url, stream=True)
        except requests.exceptions.ConnectionError:
            continue

        if response.status_code == 200:
            with open(filename, "wb") as out_file:
                shutil.copyfileobj(response.raw, out_file)
                return

def _get_v3_file(filetype, location):
    cache_file = _get_cache_filename(filetype, location)
    if os.path.exists(cache_file):
        return cache_file

    dirname = os.path.dirname(cache_file)
    if not os.path.exists(dirname):
        os.makedirs(dirname)

    _process_request(filetype, location, cache_file)
    return cache_file

def v3_file(location):
    if location.endswith(".sxl"):
        filetype = "sxl"
        location = location[:-4]
    elif location.endswith(".png"):
        filetype = "png"
        location = location[:-4]
    elif location.endswith(".it2"):
        filetype = "it2"
        location = location[:-4]
    else:
        if "." in location:
            return "ERROR Unknown filetype: %s" % location
        else:
            filetype = "png"
    return _get_v3_file(filetype, location)

def _png_file(location):
    filename = _get_v3_file("png", location)
    with open(filename, "rb") as png_file:
        return png_file.read()

def _main_view(location):
    filetype = ""
    if location.endswith(".sxl"):
        filetype = "sxl"
        location = location[:-4]
    elif location.endswith(".it2"):
        filetype = "it2"
        location = location[:-4]

    if filetype != "":
        filename = _get_v3_file(filetype, location)
        with open(filename, "r") as f_data:
            return f_data.read()

    with open(os.path.join(MYDIR, "share", "v3.txt"), "r") as f_v3:
        v3_text = f_v3.read()
    return v3_text

def main(query, parsed_query, data):
    parsed_query["locale"] = "en_US"

    html_output = parsed_query["html_output"]

    location = parsed_query["orig_location"]
    if location is None:
        location = parsed_query["location"]

    if parsed_query.get("png_filename"):
        return _png_file(location)

    filename = location + ".png"


    if html_output:
        output = """
<html>
<head>
<title>Weather report for {location}</title>
<link rel="stylesheet" type="text/css" href="/files/style.css" />
<style>
figure {{
  position: relative;
}}

img {{
  max-width: 100%;
  position: absolute;
}}
</style>
</head>
<body>
<figure style="padding-bottom: calc((500/600)*100%)">
  <img src="/v3/{filename}" />
</figure>

<!--
  <div>
    <img src="/v3/{filename}" height="600"/>
  </div>
-->
</body>
</html>
""".format(location=location, filename=filename)
    else:
        output = _main_view(location)
        if parsed_query.get('no-terminal', False):
            output = remove_ansi(output)
    return output
