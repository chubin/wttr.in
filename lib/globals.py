import logging
import os

MYDIR = os.path.abspath(os.path.dirname( os.path.dirname('__file__') ))

GEOLITE = os.path.join(MYDIR, 'data', "GeoLite2-City.mmdb")
WEGO = "/home/igor/go/bin/we-lang"
PYPHOON = "/home/igor/wttr.in/pyphoon/bin/pyphoon-lolcat"

CACHEDIR  = os.path.join( MYDIR, "cache" )
IP2LCACHE = os.path.join( MYDIR, "cache/ip2l" )

ALIASES   = os.path.join( MYDIR, "share/aliases" )
ANSI2HTML = os.path.join( MYDIR, "share/ansi2html.sh" )
BLACKLIST = os.path.join( MYDIR, "share/blacklist" )

HELP_FILE           = os.path.join( MYDIR, 'share/help.txt' )
BASH_FUNCTION_FILE  = os.path.join( MYDIR, 'share/bash-function.txt' )
TRANSLATION_FILE    = os.path.join( MYDIR, 'share/translation.txt' )
TEST_FILE           = os.path.join( MYDIR, 'share/test-NAME.txt' )

IATA_CODES_FILE = os.path.join(MYDIR, 'share/list-of-iata-codes.txt')

LOG_FILE  = os.path.join( MYDIR, 'log/main.log' )
TEMPLATES = os.path.join( MYDIR, 'share/templates' )
STATIC    = os.path.join( MYDIR, 'share/static' )

NOT_FOUND_LOCATION = "not found"
DEFAULT_LOCATION = "oymyakon"

MALFORMED_RESPONSE_HTML_PAGE = open(os.path.join(STATIC, 'malformed-response.html')).read()

LISTEN_HOST = ""
LISTEN_PORT = 8002

PLAIN_TEXT_AGENTS = [
    "curl",
    "httpie",
    "lwp-request",
    "wget",
    "python-requests"
]

PLAIN_TEXT_PAGES = [':help', ':bash.function', ':translation']

IP2LOCATION_KEY = ''

def error(text):
    if not text.startswith('Too many queries'):
        print text
    logging.error("ERROR "+text)
    raise RuntimeError(text)

def log(text):
    if not text.startswith('Too many queries'):
        print text
        logging.info(text)

def get_help_file(lang):
    help_file = os.path.join(MYDIR, 'share/translations/%s-help.txt' % lang)
    print ">>>", help_file
    if os.path.exists(help_file):
        return help_file
    else:
        return HELP_FILE

