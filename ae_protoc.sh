#!/bin/bash

set -e

if [[ "$#" -eq 0 ]]; then
  echo "Usage: ae_protoc.sh <file1.proto> <file2.proto>"
  echo "  This script will use the protoc in your path"
  echo "to compile each proto into go source and will"
  echo "edit each file to be usable on AppEngine by"
  echo "removing references to the goprotobuf library."
  exit 1
fi

if [[ -z "$GO_STUBS" ]]; then
  export GO_STUBS="web"
fi

for FILE in "$@"; do
  echo "Compiling $FILE..."
  protoc --go_out=. "$FILE"

  PB_FILE="${FILE%.proto}.pb.go"
  echo "Sanitizing $PB_FILE..."
  {
    echo "H"              # Display human-readable errors, if any
    echo "g/goprotobuf/d" # Delete lines containing goprotobuf
    echo "g/proto\./d"    # Delete lines calling into the library
    echo "w"              # Write
    echo "q"              # Quit
  } | ed -s "$PB_FILE"
done
