"""
global configuration of the project

External environment variables:

    WTTR_MYDIR
    WTTR_GEOLITE
    WTTR_WEGO
    WTTR_LISTEN_HOST
    WTTR_LISTEN_PORT
    WTTR_USER_AGENT

"""
from __future__ import print_function

import logging
import os
import re

MYDIR = os.path.abspath(os.path.dirname(os.path.dirname('__file__')))

if "WTTR_GEOLITE" in os.environ:
    GEOLITE = os.environ["WTTR_GEOLITE"]
else:
    GEOLITE = os.path.join(MYDIR, 'data', "GeoLite2-City.mmdb")

WEGO = os.environ.get("WTTR_WEGO", "/home/igor/go/bin/we-lang")
PYPHOON = "pyphoon-lolcat"

_DATADIR = "/wttr.in"
_LOGDIR = "/wttr.in/log"

IP2LCACHE = os.path.join(_DATADIR, "cache/ip2l/")
PNG_CACHE = os.path.join(_DATADIR, "cache/png")
LRU_CACHE = os.path.join(_DATADIR, "cache/lru")

LOG_FILE = os.path.join(_LOGDIR, 'main.log')
MISSING_TRANSLATION_LOG = os.path.join(_LOGDIR, 'missing-translation/%s.log')

ALIASES = os.path.join(MYDIR, "share/aliases")
ANSI2HTML = os.path.join(MYDIR, "share/ansi2html.sh")
BLACKLIST = os.path.join(MYDIR, "share/blacklist")

HELP_FILE = os.path.join(MYDIR, 'share/help.txt')
BASH_FUNCTION_FILE = os.path.join(MYDIR, 'share/bash-function.txt')
TRANSLATION_FILE = os.path.join(MYDIR, 'share/translation.txt')

IATA_CODES_FILE = os.path.join(MYDIR, 'share/list-of-iata-codes.txt')

TEMPLATES = os.path.join(MYDIR, 'share/templates')
STATIC = os.path.join(MYDIR, 'share/static')

NOT_FOUND_LOCATION = "not found"
DEFAULT_LOCATION = "oymyakon"

MALFORMED_RESPONSE_HTML_PAGE = open(os.path.join(STATIC, 'malformed-response.html')).read()

GEOLOCATOR_SERVICE = 'http://localhost:8004'

# number of queries from the same IP address is limited
# (minute, hour, day) limitations:
QUERY_LIMITS = (300, 3600, 24*3600)

LISTEN_HOST = os.environ.get("WTTR_LISTEN_HOST", "")
try:
    LISTEN_PORT = int(os.environ.get("WTTR_LISTEN_PORT"))
except (TypeError, ValueError):
    LISTEN_PORT = 8002

PROXY_HOST = "127.0.0.1"
PROXY_PORT = 5001
PROXY_CACHEDIR = os.path.join(_DATADIR, "cache/proxy-wwo/")

MY_EXTERNAL_IP = '5.9.243.187'

PLAIN_TEXT_AGENTS = [
    "curl",
    "httpie",
    "lwp-request",
    "wget",
    "python-requests",
    "openbsd ftp",
    "powershell",
    "fetch",
    "aiohttp",
]

PLAIN_TEXT_PAGES = [':help', ':bash.function', ':translation', ':iterm2']

_IPLOCATION_ORDER = os.environ.get(
    "WTTR_IPLOCATION_ORDER",
    'geoip,ip2location,ipinfo')
IPLOCATION_ORDER = _IPLOCATION_ORDER.split(',')

_IP2LOCATION_KEY_FILE = os.environ.get(
    "WTTR_IP2LOCATION_KEY_FILE",
    os.environ['HOME'] + '/.ip2location.key')
IP2LOCATION_KEY = None
if os.path.exists(_IP2LOCATION_KEY_FILE):
    IP2LOCATION_KEY = open(_IP2LOCATION_KEY_FILE, 'r').read().strip()

_IPINFO_KEY_FILE = os.environ.get(
    "WTTR_IPINFO_KEY_FILE",
    os.environ['HOME'] + '/.ipinfo.key')
IPINFO_TOKEN = None
if os.path.exists(_IPINFO_KEY_FILE):
    IPINFO_TOKEN = open(_IPINFO_KEY_FILE, 'r').read().strip()

_WWO_KEY_FILE = os.environ.get(
    "WTTR_WWO_KEY_FILE",
    os.environ['HOME'] + '/.wwo.key')
WWO_KEY = "key-is-not-specified"
USE_METNO = True
USER_AGENT = os.environ.get("WTTR_USER_AGENT", "")
if os.path.exists(_WWO_KEY_FILE):
    WWO_KEY = open(_WWO_KEY_FILE, 'r').read().strip()
    USE_METNO = False

def error(text):
    "log error `text` and raise a RuntimeError exception"

    if not text.startswith('Too many queries'):
        print(text)
    logging.error("ERROR %s", text)
    raise RuntimeError(text)

def log(text):
    "log error `text` and do not raise any exceptions"

    if not text.startswith('Too many queries'):
        print(text)
        logging.info(text)

def debug_log(text):
    """
    Write `text` to the debug log
    """

    with open('/tmp/wttr.in-debug.log', 'a') as f_debug:
        f_debug.write(text+'\n')

def get_help_file(lang):
    "Return help file for `lang`"

    help_file = os.path.join(MYDIR, 'share/translations/%s-help.txt' % lang)
    if os.path.exists(help_file):
        return help_file
    return HELP_FILE

def remove_ansi(sometext):
    ansi_escape = re.compile(r'(\x9B|\x1B\[)[0-?]*[ -\/]*[@-~]')
    return ansi_escape.sub('', sometext)
