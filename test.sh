#!/bin/bash

REPO="github.com/kylelemons/go-rpcgen"

function err {
  echo "$@"
  exit 1
}

set -e

if ! which protoc >/dev/null; then
  err "Could not find 'protoc'"
fi

echo "Building protoc-gen-go..."
go build -o protoc-gen-go/protoc-gen-go $REPO/protoc-gen-go
export PATH="protoc-gen-go/:$PATH"

echo "Building protobufs..."
for PROTO in $(find . -name "*.proto" | grep -v "option.proto"); do
  echo " - Compiling ${PROTO}..."
  GO_STUBS="rpc,web" protoc --go_out=. ${PROTO}
done

echo "Testing packages..."
PACKAGES=$(find ./* -type d | grep -v "ae" | grep -v "examples$")
go test -i ${PACKAGES}
go test ${PACKAGES}

go install ./protoc-gen-go
