# vim: set encoding=utf-8
# pylint: disable=wrong-import-position

"""
Main view (wttr.in) implementation.
The module is a wrapper for the modified Wego program.
"""

import sys
import re

from gevent.subprocess import Popen, PIPE

sys.path.insert(0, "..")
from translations import get_message, SUPPORTED_LANGS
from globals import WEGO, NOT_FOUND_LOCATION, DEFAULT_LOCATION, ANSI2HTML, \
                    error, remove_ansi


def get_wetter(parsed_query):

    location = parsed_query['location']
    html = parsed_query['html_output']
    lang = parsed_query['lang']

    location_not_found = False
    if location == NOT_FOUND_LOCATION:
        location_not_found = True

    stderr = ""
    returncode = 0
    if not location_not_found:
        stdout, stderr, returncode = _wego_wrapper(location, parsed_query)

    if location_not_found or \
        (returncode != 0 \
            and ('Unable to find any matching weather'
                 ' location to the parsed_query submitted') in stderr):
            stdout, stderr, returncode = _wego_wrapper(NOT_FOUND_LOCATION, parsed_query)
            location_not_found = True
            stdout += get_message('NOT_FOUND_MESSAGE', lang)

    if "\n" in stdout:
        first_line, stdout = _wego_postprocessing(location, parsed_query, stdout)
    else:
        first_line = ""

    # This is where you append date to stdout
    # stdout["timestamp"] = date().epochTime() - example

    if html:
        return _htmlize(stdout, first_line, parsed_query)
    return stdout

def _wego_wrapper(location, parsed_query):
    lang = parsed_query['lang']
    location_name = parsed_query['override_location_name']

    cmd = [WEGO, '-location=%s' % location]

    if lang and lang in SUPPORTED_LANGS:
        cmd += ['-owm-lang', lang, '-wwo-lang', lang, '-forecast-lang', lang]

    if parsed_query.get('use_imperial', False):
        cmd += ['-units', 'imperial']

    proc = Popen(cmd, stdout=PIPE, stderr=PIPE)
    stdout, stderr = proc.communicate()
    stdout = stdout.decode("utf-8")
    stderr = stderr.decode("utf-8")

    # Make another call to get maps: http://maps.openweathermap.org/maps/2.0/weather/{op}/{z}/{x}/{y}
    # Check for error on API call, if none, add weather maps converted to ascii, to stdout
    # https://github.com/SketchingDev/image-to-ascii-converter/blob/master/sketchingdev/image_to_ascii/converter.py#L50
    # May need to convert for web display

    return stdout, stderr, proc.returncode

def _wego_postprocessing(location, parsed_query, stdout):
    full_address = parsed_query['full_address']
    lang = parsed_query['lang']

    if 'days' in parsed_query:
        if parsed_query['days'] == '0':
            stdout = "\n".join(stdout.splitlines()[:7]) + "\n"
        if parsed_query['days'] == '1':
            stdout = "\n".join(stdout.splitlines()[:17]) + "\n"
        if parsed_query['days'] == '2':
            stdout = "\n".join(stdout.splitlines()[:27]) + "\n"


    first = stdout.splitlines()[0]
    rest = stdout.splitlines()[1:]
    if parsed_query.get('no-caption', False):
        if ':' in first:
            first = first.split(":", 1)[1]
            stdout = "\n".join([first.strip()] + rest) + "\n"

    if parsed_query.get('no-terminal', False):
        stdout = remove_ansi(stdout)

    if parsed_query.get('no-city', False):
        stdout = "\n".join(stdout.splitlines()[2:]) + "\n"

    if full_address \
        and parsed_query.get('format', 'txt') != 'png' \
        and (not parsed_query.get('no-city')
             and not parsed_query.get('no-caption')
             and not parsed_query.get('days') == '0'):
        line = "%s: %s [%s]\n" % (
            get_message('LOCATION', lang),
            full_address,
            location)
        stdout += line

    if parsed_query.get('padding', False):
        lines = [x.rstrip() for x in stdout.splitlines()]
        max_l = max(len(remove_ansi(x)) for x in lines)
        last_line = " "*max_l + "   .\n"
        stdout = " \n" + "\n".join("  %s  " %x for x in lines) + "\n" + last_line

    return first, stdout

def _htmlize(ansi_output, title, parsed_query):
    """Return HTML representation of `ansi_output`.
    Use `title` as the title of the page.
    Format page according to query parameters from `parsed_query`."""

    cmd = ["bash", ANSI2HTML, "--palette=solarized"]
    if not parsed_query.get('inverted_colors'):
        cmd += ["--bg=dark"]

    proc = Popen(cmd, stdin=PIPE, stdout=PIPE, stderr=PIPE)
    stdout, stderr = proc.communicate(ansi_output.encode("utf-8"))
    stdout = stdout.decode("utf-8")
    stderr = stderr.decode("utf-8")
    if proc.returncode != 0:
        error(stdout + stderr)

    if parsed_query.get('inverted_colors'):
        stdout = stdout.replace(
            '<body class="">', '<body class="" style="background:white;color:#777777">')

    title = "<title>%s</title>" % title
    opengraph = _get_opengraph(parsed_query)
    stdout = re.sub("<head>", "<head>" + title + opengraph, stdout)
    return stdout

def _get_opengraph(parsed_query):
    """Return OpenGraph data for `parsed_query`"""

    url = parsed_query['request_url'] or ""
    pic_url = url.replace('?', '_')

    return (
        '<meta property="og:image" content="%(pic_url)s_0pq.png" />'
        '<meta property="og:site_name" content="wttr.in" />'
        '<meta property="og:type" content="profile" />'
        '<meta property="og:url" content="%(url)s" />'
    ) % {
        'pic_url': pic_url,
        'url': url,
    }
