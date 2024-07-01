#!/usr/bin/env bash

LOG_DIR="/wttr.in/log"
LOG_FILE="$LOG_DIR/diskspace.log"

DISK=/wttr.in

log() {
  mkdir -p "$LOG_DIR"

  echo "$(date +"[%Y-%m-%d %H:%M:%S]") $*" | tee -a "$LOG_FILE"
}

log $(df -k "$DISK" | tail -1 | awk '{print $4}')
