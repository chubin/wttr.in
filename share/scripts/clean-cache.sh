#!/bin/bash

LOGFILE=/tmp/clean-cache.log

_log() {
  echo "$(date +"[%Y-%m-%d %H:%M:%S]") $*" >> "$LOGFILE"
}

_log_pipe() {
  while read -r line; do
    _log "$line"
  done
}

CACHEDIR="/wttr.in/cache"

_log Cleaning up the cache.
_log Before:
df -h "$CACHEDIR" | _log_pipe

for dir in wego proxy-wwo png lru
do
  mv "${CACHEDIR}/${dir}" "${CACHEDIR}/${dir}.old"
  mkdir "${CACHEDIR}/${dir}"
  rm -rf "${CACHEDIR}/${dir}.old"
done

_log After:
df -h "$CACHEDIR" | _log_pipe
_log "============================="

cd /wttr.in/log
mv main.log main.log.1
touch main.log

