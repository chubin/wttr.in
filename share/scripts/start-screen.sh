#!/bin/bash

SESSION_NAME=wttr.in
SCREENRC_PATH=$(dirname $(dirname "$0"))/screenrc

screen -dmS "$SESSION_NAME" -c "$SCREENRC_PATH"

