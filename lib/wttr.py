# vim: set encoding=utf-8

from __future__ import print_function
import gevent
from gevent.pywsgi import WSGIServer
from gevent.queue import Queue
from gevent.monkey import patch_all
from gevent.subprocess import Popen, PIPE, STDOUT
patch_all()

import os
import re
import time
import dateutil.parser

from translations import get_message, FULL_TRANSLATION, PARTIAL_TRANSLATION, SUPPORTED_LANGS
from globals import WEGO, PYPHOON, CACHEDIR, ANSI2HTML, \
                    NOT_FOUND_LOCATION, DEFAULT_LOCATION, TEST_FILE, \
                    log, error

def _is_invalid_location(location):
    if '.png' in location:
        return True

def remove_ansi(sometext):
    ansi_escape = re.compile(r'(\x9B|\x1B\[)[0-?]*[ -\/]*[@-~]')
    return ansi_escape.sub('', sometext)

def get_wetter(location, ip, html=False, lang=None, query=None, location_name=None, full_address=None, url=None):

    local_url = url
    local_location = location

    def get_opengraph():

        if local_url is None:
            url = ""
        else:
            url = local_url.encode('utf-8')

        if local_location is None:
            location = ""
        else:
            location = local_location.encode('utf-8')

        pic_url = url.replace('?', '_')

        return (
            '<meta property="og:image" content="%(pic_url)s_0pq.png" />'
            '<meta property="og:site_name" content="wttr.in" />'
            '<meta property="og:type" content="profile" />'
            '<meta property="og:url" content="%(url)s" />'
        ) % {
            'pic_url': pic_url,
            'url': url,
            'location': location,
        }

            # '<meta property="og:title" content="Weather report: %(location)s" />'
            # '<meta content="Partly cloudy // 6-8 °C // ↑ 9 km/h // 10 km // 0.4 mm" property="og:description" />'

    def get_filename(location, lang=None, query=None, location_name=None):
        location = location.replace('/', '_')
        timestamp = time.strftime( "%Y%m%d%H", time.localtime() )

        imperial_suffix = ''
        if query.get('use_imperial', False):
            imperial_suffix = '-imperial'

        lang_suffix = ''
        if lang is not None:
            lang_suffix = '-lang_%s' % lang

        if query != None:
            query_line = "_" + "_".join("%s=%s" % (key, value) for (key, value) in query.items())
        else:
            query_line = ""

        if location_name is None:
            location_name = ""
        return "%s/%s/%s%s%s%s%s" % (CACHEDIR, location, timestamp, imperial_suffix, lang_suffix, query_line, location_name)

    def save_weather_data(location, filename, lang=None, query=None, location_name=None, full_address=None):
     
        if _is_invalid_location( location ):
            error("Invalid location: %s" % location)
        
        NOT_FOUND_MESSAGE_HEADER = ""
        while True:
            location_not_found = False
            if location in [ "test-thunder" ]:
                test_name = location[5:]
                test_file = TEST_FILE.replace('NAME', test_name)
                stdout = open(test_file, 'r').read()
                stderr = ""
                break
            if location == NOT_FOUND_LOCATION:
                location_not_found = True
                location = DEFAULT_LOCATION
            
            cmd = [WEGO, '--city=%s' % location]

            if query.get('inverted_colors'):
                cmd += ['-inverse']

            if query.get('use_ms_for_wind'):
                cmd += ['-wind_in_ms']

            if query.get('narrow'):
                cmd += ['-narrow']

            if lang and lang in SUPPORTED_LANGS:
                cmd += ['-lang=%s'%lang]

            if query.get('use_imperial', False):
                cmd += ['-imperial']

            if location_name:
                cmd += ['-location_name', location_name]

            p = Popen(cmd, stdout=PIPE, stderr=PIPE)
            stdout, stderr = p.communicate()
            if p.returncode != 0:
                print("ERROR: location not found: %s" % location)
                if 'Unable to find any matching weather location to the query submitted' in stderr:
                    if location != NOT_FOUND_LOCATION:
                        NOT_FOUND_MESSAGE_HEADER = u"ERROR: %s: %s\n---\n\n" % (get_message('UNKNOWN_LOCATION', lang), location)
                        location = NOT_FOUND_LOCATION
                        continue
                error(stdout + stderr)
            break

        dirname = os.path.dirname(filename)
        if not os.path.exists(dirname):
            os.makedirs(dirname)
        
        if location_not_found:
            stdout += get_message('NOT_FOUND_MESSAGE', lang).encode('utf-8')
            stdout = NOT_FOUND_MESSAGE_HEADER.encode('utf-8') + stdout

        if 'days' in query:
            if query['days'] == '0':
                stdout = "\n".join(stdout.splitlines()[:7]) + "\n"
            if query['days'] == '1':
                stdout = "\n".join(stdout.splitlines()[:17]) + "\n"
            if query['days'] == '2':
                stdout = "\n".join(stdout.splitlines()[:27]) + "\n"

        first = stdout.splitlines()[0].decode('utf-8')
        rest = stdout.splitlines()[1:]
        if query.get('no-caption', False):

            separator = None
            if ':' in first:
                separator = ':'
            if u'：' in first:
                separator = u'：'

            if separator:
                first = first.split(separator,1)[1]
                stdout = "\n".join([first.strip().encode('utf-8')] + rest) + "\n"

        if query.get('no-terminal', False):
            stdout = remove_ansi(stdout)

        if query.get('no-city', False):
            stdout = "\n".join(stdout.splitlines()[2:]) + "\n"

        if full_address \
            and query.get('format', 'txt') != 'png' \
            and (not query.get('no-city') and not query.get('no-caption')):
            line = "%s: %s [%s]\n" % (
                get_message('LOCATION', lang).encode('utf-8'),
                full_address.encode('utf-8'),
                location.encode('utf-8'))
            stdout += line

        if query.get('padding', False):
            lines = [x.rstrip() for x in stdout.splitlines()]
            max_l = max(len(remove_ansi(x).decode('utf8')) for x in lines)
            last_line = " "*max_l + "   .\n"
            stdout = " \n" + "\n".join("  %s  " %x for x in lines) + "\n" + last_line 

        open(filename, 'w').write(stdout)

        cmd = ["bash", ANSI2HTML, "--palette=solarized"]
        if not query.get('inverted_colors'):
            cmd += ["--bg=dark"]

        p = Popen(cmd,  stdin=PIPE, stdout=PIPE, stderr=PIPE )
        stdout, stderr = p.communicate(stdout)
        if p.returncode != 0:
            error(stdout + stderr)

        if query.get('inverted_colors'):
            stdout = stdout.replace('<body class="">', '<body class="" style="background:white;color:#777777">')
        
        title = "<title>%s</title>" % first.encode('utf-8')
        opengraph = get_opengraph()
        stdout = re.sub("<head>", "<head>" + title + opengraph, stdout)
        open(filename+'.html', 'w').write(stdout)

    filename = get_filename(location, lang=lang, query=query, location_name=location_name)
    if not os.path.exists(filename):
        save_weather_data(location, filename, lang=lang, query=query, location_name=location_name, full_address=full_address)
    if html:
        filename += '.html'

    return open(filename).read()

def get_moon(location, html=False, lang=None, query=None):
    if query is None:
        query = {}

    date = None
    if '@' in location:
        date = location[location.index('@')+1:]
        location = location[:location.index('@')]

    cmd = [PYPHOON]
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

    if query.get('no-terminal', False):
        stdout = remove_ansi(stdout)

    if html:
        p = Popen(["bash", ANSI2HTML, "--palette=solarized", "--bg=dark"],  stdin=PIPE, stdout=PIPE, stderr=PIPE)
        stdout, stderr = p.communicate(stdout)
        if p.returncode != 0:
            error(stdout + stderr)

    return stdout

