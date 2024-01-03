import sys

import os
import dateutil.parser

from gevent.subprocess import Popen, PIPE

sys.path.insert(0, "..")
import constants
import parse_query
import globals

def get_moon(parsed_query):

    location = parsed_query['orig_location']
    html = parsed_query['html_output']
    lang = parsed_query['lang']
    hemisphere = parsed_query['hemisphere']

    date = None
    if '@' in location:
        date = location[location.index('@')+1:]
        location = location[:location.index('@')]

    cmd = [globals.PYPHOON]
    if lang:
        cmd += ["-l", lang]

    if not hemisphere:
        cmd += ["-s", "south"]

    if date:
        try:
            dateutil.parser.parse(date)
        except Exception as e:
            print("ERROR: %s" % e)
        else:
            cmd += [date]

    p = Popen(cmd, stdout=PIPE, stderr=PIPE)
    stdout = p.communicate()[0]
    stdout = stdout.decode("utf-8")

    if parsed_query.get('no-terminal', False):
        stdout = globals.remove_ansi(stdout)

    if parsed_query.get('dumb', False):
        stdout = stdout.translate(globals.TRANSLATION_TABLE)

    if html:
        p = Popen(
            ["bash", globals.ANSI2HTML, "--palette=solarized", "--bg=dark"],
            stdin=PIPE, stdout=PIPE, stderr=PIPE)
        stdout, stderr = p.communicate(stdout.encode("utf-8"))
        stdout = stdout.decode("utf-8")
        stderr = stderr.decode("utf-8")
        if p.returncode != 0:
            globals.error(stdout + stderr)

    return stdout
