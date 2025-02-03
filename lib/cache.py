"""
LRU-Cache implementation for formatted (`format=`) answers
"""

import datetime
import re
import time
import os
import hashlib
import random

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

    if isinstance(answer, str) and "%{{NOW(" in answer:
        answer = re.sub(
            r"%{{NOW\(([^}]*)\)}}", lambda x: _now_in_tz(x.group(1)), answer
        )

    return answer


def get_signature(user_agent, query_string, client_ip_address, lang):
    """
    Get cache signature based on `user_agent`, `url_string`,
    `lang`, and `client_ip_address`
    Return `None` if query should not be cached.
    """

    if "?" in query_string:
        location = query_string.split("?", 1)[0]
    else:
        location = query_string
    if location.startswith("http://"):
        location = location[7:]
    elif location.startswith("https://"):
        location = location[8:]
    if ":" in location:
        return None

    signature = "%s:%s:%s:%s" % (user_agent, query_string, client_ip_address, lang)
    print(signature)
    return signature


def get(signature):
    """
    If `update_answer` is not True, return answer as it is
    stored in the cache. Otherwise update it, using
    the `_update_answer` function.
    """

    if not signature:
        return None

    value_record = CACHE.get(signature)
    if not value_record:
        return None

    value = value_record["val"]
    expiry = value_record["expiry"]
    if value and time.time() < expiry:
        if value.startswith("file:") or value.startswith("bfile:"):
            value = _read_from_file(signature, sighash=value)
            if not value:
                return None
        return _update_answer(value)
    return None


def _randint(minimum, maximum):
    return random.randrange(maximum - minimum)


def store(signature, value):
    """
    Store in cache `value` for `signature`
    """
    if not signature:
        return _update_answer(value)

    if len(value) >= MIN_SIZE_FOR_FILECACHE:
        value_to_store = _store_in_file(signature, value)
    else:
        value_to_store = value

    value_record = {
        "val": value_to_store,
        "expiry": time.time() + _randint(1000, 2000),
    }

    CACHE[signature] = value_record

    return _update_answer(value)


def _hash(signature):
    return hashlib.md5(signature.encode("utf-8")).hexdigest()


def _store_in_file(signature, value):
    """Store `value` for `signature` in cache file.
    Return file name (signature_hash) as the result.
    `value` can be string as well as bytes.
    Returned filename is prefixed with "file:" (for text files)
    or "bfile:" (for binary files).
    """

    signature_hash = _hash(signature)
    filename = os.path.join(LRU_CACHE, signature_hash)
    if not os.path.exists(LRU_CACHE):
        os.makedirs(LRU_CACHE)

    if isinstance(value, bytes):
        mode = "wb"
        signature_hash = "bfile:%s" % signature_hash
    else:
        mode = "w"
        signature_hash = "file:%s" % signature_hash

    with open(filename, mode) as f_cache:
        f_cache.write(value)
    return signature_hash


def _read_from_file(signature, sighash=None):
    """Read value for `signature` from cache file,
    or return None if file is not found.
    If `sighash` is specified, do not calculate file name
    from signature, but use `sighash` instead.

    `sigash` can be prefixed with "file:" (for text files)
    or "bfile:" (for binary files).
    """

    mode = "r"
    if sighash:
        if sighash.startswith("file:"):
            sighash = sighash[5:]
        elif sighash.startswith("bfile:"):
            sighash = sighash[6:]
            mode = "rb"
    else:
        sighash = _hash(signature)

    filename = os.path.join(LRU_CACHE, sighash)
    if not os.path.exists(filename):
        return None

    with open(filename, mode) as f_cache:
        return f_cache.read()
