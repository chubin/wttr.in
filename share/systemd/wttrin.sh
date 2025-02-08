#!/usr/bin/env bash

SESSION_NAME=""
SRC_DIR=/home/igor/src/wttr.in
SERVICES_FILE=config/services/services.yaml

start_service() {
  local name="$1"
  local workdir="$2"
  local cmd="$3"

  if [[ -z "$SESSION_NAME" ]]; then
    echo Unknown SESSION_NAME. Exiting >&2
    exit 1
  fi

  local WINDOW_NAME="$name"
  local TEXT_TO_ENTER="cd $workdir && $cmd"

  if ! tmux has-session -t "$SESSION_NAME" >& /dev/null; then
    tmux new-session -d -s "$SESSION_NAME"
  fi

  # Create a new window if it doesn't exist
  tmux list-windows -t "$SESSION_NAME" | grep -q "^.\+ $WINDOW_NAME " || tmux new-window -t "$SESSION_NAME" -n "$WINDOW_NAME"

  sleep 0.05

  # Send text to the new window and press Enter
  tmux send-keys -t "$SESSION_NAME:$WINDOW_NAME" "$TEXT_TO_ENTER" C-m
}


main() {
  local name
  local cmd

  if [[ -n $1 ]]; then
    SESSION_NAME=$1
  else
    echo Usage: $0 SESSION_NAME >&2
    exit 1
  fi

  cd "$SRC_DIR" || exit 1

  set -x

  while read -r line; do
    name=$(jq -r .name <<< "$line")
    workdir=$(jq -r .workdir <<< "$line")
    cmd=$(jq -r .command <<< "$line")

    start_service "$name" "$workdir" "$cmd"
  done <<< "$(yq -c .services[] < "$SERVICES_FILE")"
}

main "$@"
