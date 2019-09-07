"""
LRU-Cache implementation for formatted (`format=`) answers
"""

import datetime
import re
import time
import pylru
import pytz

CACHE_SIZE = 10000
CACHE = pylru.lrucache(CACHE_SIZE)

def _update_answer(answer):
    def _now_in_tz(timezone):
        return datetime.datetime.now(pytz.timezone(timezone)).strftime("%H:%M:%S%z")

    if "%{{NOW(" in answer:
        answer = re.sub(r"%{{NOW\(([^}]*)\)}}", lambda x: _now_in_tz(x.group(1)), answer)

    return answer

def get_signature(user_agent, query_string, client_ip_address, lang):
    """
    Get cache signature based on `user_agent`, `url_string`,
    `lang`, and `client_ip_address`
    """

    timestamp = int(time.time()) / 1000
    signature = "%s:%s:%s:%s:%s" % \
        (user_agent, query_string, client_ip_address, lang, timestamp)
    return signature

def get(signature):
    """
    If `update_answer` is not True, return answer as it is
    stored in the cache. Otherwise update it, using
    the `_update_answer` function.
    """

    if signature in CACHE:
        return _update_answer(CACHE[signature])
    return None

def store(signature, value):
    """
    Store in cache `value` for `signature`
    """
    CACHE[signature] = value
    return _update_answer(value)
