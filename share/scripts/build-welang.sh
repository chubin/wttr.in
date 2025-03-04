#!/usr/bin/env bash

# The scipr is used to build a standalone we-lang binary from
# the wttr.in Go source code. The script is needed as long as
# a standalone we-lang binary is used.
#
# The script requires a configured Go compiler in PATH.


show_usage() {
  cat <<EOF
Usage:

    $0 WELANG_PATH

Build we-lang from the source code and install it as a binary to WELANG_PATH.
EOF
}

fatal() {
  rm -rf "$DIR"
  echo "FATAL: $*"
  exit 1
}

add_main() {
  cat <<EOF > main.go
package main

func main() {
  Cmd()
}
EOF
}

if [[ -z $1 ]]; then
  show_usage
  exit 0
fi

DST_FILE="$1"
DIR=$(mktemp -d /tmp/build-welang-XXXXXXXXX)

cp -R "internal/view/v1" "$DIR/v1"

sed -i 's/^package .*/package main/' "$DIR"/v1/*.go

cd "$DIR"/v1 || fatal "can't change into the build directory"

go mod init github.com/chubin/we-lang || fatal "Can't do 'go mod init'"
go mod tidy || fatal "Can't do 'go mod tidy'"
add_main
go build -o "$DST_FILE" *.go || fatal "Building error"

rm -rf "$DIR"


