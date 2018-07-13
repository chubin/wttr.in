from __future__ import print_function
import logging
import os

MYDIR = os.path.abspath(os.path.dirname( os.path.dirname('__file__') ))

GEOLITE = os.path.join( MYDIR, "GeoLite2-City.mmdb" )
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

LOG_FILE  = os.path.join( MYDIR, 'log/main.log' )
TEMPLATES = os.path.join( MYDIR, 'share/templates' )
STATIC    = os.path.join( MYDIR, 'share/static' )

NOT_FOUND_LOCATION = "not found"
DEFAULT_LOCATION = "oymyakon"

MALFORMED_RESPONSE_HTML_PAGE = open(os.path.join(STATIC, 'malformed-response.html')).read()

def error(text):
    if not text.startswith('Too many queries'):
        print(text)
    logging.error("ERROR "+text)
    raise RuntimeError(text)

def log(text):
    if not text.startswith('Too many queries'):
        print(text)
        logging.info(text)


