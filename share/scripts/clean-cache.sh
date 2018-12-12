#!/bin/bash

CACHEDIR="/wttr.in/cache"

for dir in wego proxy-wwo png
do
  mv "${CACHEDIR}/${dir}" "${CACHEDIR}/${dir}.old"
  mkdir "${CACHEDIR}/${dir}"
  rm -rf "${CACHEDIR}/${dir}.old"
done
