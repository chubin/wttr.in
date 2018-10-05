#!/usr/bin/env python
# vim: set encoding=utf-8

import sys
import os
MYDIR = os.path.abspath(
    os.path.dirname(os.path.dirname('__file__')))
sys.path.append("%s/lib/" % MYDIR)
import srv
