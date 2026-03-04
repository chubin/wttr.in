#!/usr/bin/env bash

build() {
  CGO_CFLAGS="-Wno-return-local-addr" go build -o main main.go
}

case "$1" in
  build)
    build "$@"
    ;;
  *)
    echo "Unknown command: $1" >&2
    exit 1
    ;;
esac
