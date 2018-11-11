#!/bin/bash

CACHEDIR="/wttr.in/cache"

mv "${CACHEDIR}/wego" "${CACHEDIR}/wego.old"
mkdir "${CACHEDIR}/wego"
rm -rf "${CACHEDIR}/wego.old"
