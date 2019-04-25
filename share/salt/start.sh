#!/bin/sh
export WEGORC="/srv/ephemeral/.wegorc"
export GOPATH="/srv/ephemeral"

export WTTR_MYDIR="/srv/ephemeral/wttr.in"
export WTTR_GEOLITE="/srv/ephemeral/GeoLite2-City.mmdb"
export WTTR_WEGO="$GOPATH/bin/wego"

export WTTR_LISTEN_HOST="0.0.0.0"
export WTTR_LISTEN_PORT="80"

python $WTTR_MYDIR/bin/srv.py
