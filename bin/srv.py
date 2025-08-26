#!/usr/bin/env python
# vim: set encoding=utf-8

from gevent.pywsgi import WSGIServer
from gevent.monkey import patch_all

patch_all()

# pylint: disable=wrong-import-position,wrong-import-order
import sys
import os
import jinja2

from flask import Flask, request, send_from_directory, send_file

APP = Flask(__name__)

MYDIR = os.path.abspath(os.path.dirname(os.path.dirname("__file__")))
sys.path.append("%s/lib/" % MYDIR)

import wttr_srv
from globals import TEMPLATES, STATIC, LISTEN_HOST, LISTEN_PORT

# pylint: enable=wrong-import-position,wrong-import-order

# from view.v3 import v3_file

MY_LOADER = jinja2.ChoiceLoader(
    [
        APP.jinja_loader,
        jinja2.FileSystemLoader(TEMPLATES),
    ]
)

APP.jinja_loader = MY_LOADER


# @APP.route("/v3/<string:location>")
# def send_v3(location):
#     filepath = v3_file(location)
#     if filepath.startswith("ERROR"):
#         return filepath.rstrip("\n") + "\n"
#     return send_file(filepath)


@APP.route("/files/<path:path>")
def send_static(path):
    "Send any static file located in /files/"
    return send_from_directory(STATIC, path)


@APP.route("/favicon.ico")
def send_favicon():
    "Send static file favicon.ico"
    return send_from_directory(STATIC, "favicon.ico")


@APP.route("/malformed-response.html")
def send_malformed():
    "Send static file malformed-response.html"
    return send_from_directory(STATIC, "malformed-response.html")


@APP.route("/")
@APP.route("/<string:location>")
def wttr(location=None):
    "Main function wrapper"
    return wttr_srv.wttr(location, request)


SERVER = WSGIServer(
    (LISTEN_HOST, int(os.environ.get("WTTRIN_SRV_PORT", LISTEN_PORT))), APP
)
SERVER.serve_forever()
