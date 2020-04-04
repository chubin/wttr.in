import sys

import os
import dateutil.parser

from gevent.subprocess import Popen, PIPE

sys.path.insert(0, "..")
import constants
import parse_query
import globals

def get_moon(location, html=False, lang=None, query=None):
    if query is None:
        query = {}

    date = None
    if '@' in location:
        date = location[location.index('@')+1:]
        location = location[:location.index('@')]

    cmd = [globals.PYPHOON]
    if date:
        try:
            dateutil.parser.parse(date)
        except Exception as e:
            print("ERROR: %s" % e)
        else:
            cmd += [date]

    env = os.environ.copy()
    if lang:
        env['LANG'] = lang
    p = Popen(cmd, stdout=PIPE, stderr=PIPE, env=env)
    stdout = p.communicate()[0]
    stdout = stdout.decode("utf-8")

    if query.get('no-terminal', False):
        stdout = globals.remove_ansi(stdout)

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
