"""
LRU-Cache implementation for formatted (`format=`) answers
"""

import datetime
import re
import time
import os
import hashlib

import pytz
import pylru

from globals import LRU_CACHE

CACHE_SIZE = 10000
CACHE = pylru.lrucache(CACHE_SIZE)

# strings longer than this are stored not in ram
# but in the file cache
MIN_SIZE_FOR_FILECACHE = 80

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

    timestamp = int(time.time() / 1000)
    signature = "%s:%s:%s:%s:%s" % \
        (user_agent, query_string, client_ip_address, lang, timestamp)
    print(signature)
    return signature

def get(signature):
    """
    If `update_answer` is not True, return answer as it is
    stored in the cache. Otherwise update it, using
    the `_update_answer` function.
    """

    value = CACHE.get(signature)
    if value:
        if value.startswith("file:"):
            sighash = value[5:]
            value = _read_from_file(signature, sighash=sighash)
            if not value:
                return None
        return _update_answer(value)
    return None

def store(signature, value):
    """
    Store in cache `value` for `signature`
    """
    if len(value) < MIN_SIZE_FOR_FILECACHE:
        CACHE[signature] = value
    else:
        sighash = _store_in_file(signature, value)
        CACHE[signature] = "file:%s" % sighash
    return _update_answer(value)

def _hash(signature):
    return hashlib.md5(signature.encode("utf-8")).hexdigest()

def _store_in_file(signature, value):
    """Store `value` for `signature` in cache file.
    Return file name (signature_hash) as the result.
    """

    signature_hash = _hash(signature)
    filename = os.path.join(LRU_CACHE, signature_hash)
    if not os.path.exists(LRU_CACHE):
        os.makedirs(LRU_CACHE)
    with open(filename, "w") as f_cache:
        f_cache.write(value)
    return signature_hash

def _read_from_file(signature, sighash=None):
    """Read value for `signature` from cache file,
    or return None if file is not found.
    If `sighash` is specified, do not calculate file name
    from signature, but use `sighash` instead.
    """

    signature_hash = sighash or _hash(signature)
    filename = os.path.join(LRU_CACHE, signature_hash)
    if not os.path.exists(filename):
        return None

    with open(filename, "r") as f_cache:
        return f_cache.read()
