#!/bin/bash

function err {
  echo "$@"
  exit 1
}

set -e

if ! which protoc >/dev/null; then
  err "Could not find 'protoc'"
fi

go fix ./...
go build ./...
go test ./...

( cd compiler && go build -o ./protoc-gen-go )

PATH="compiler/:$PATH"
which protoc-gen-go

for PROTO in $(find . -name "*.proto"); do
  echo "Compiling ${PROTO}..."
  protoc --go_out=. ${PROTO}
  go fix $(dirname "${PROTO}")
done
