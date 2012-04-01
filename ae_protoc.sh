#!/bin/bash

set -e

# This will make the pb only compile for appengine
AE_COPY=1
if [[ "$1" == "--ae-only" ]]; then
  AE_COPY=0
fi

if [[ "$#" -eq 0 ]]; then
  echo "Usage: ae_protoc.sh [--ae-only] <file1.proto> [<file2.proto> ...]"
  echo "  This script will use the protoc in your path"
  echo "to compile each proto into go source and will"
  echo "edit each file to be usable on AppEngine by"
  echo "removing references to the goprotobuf library."
  echo
  echo "Unless --ae-also is specified, the compiled protobuf"
  echo "output file will be duplicated to an appengine"
  echo "specific .ae.go and both will be guarded with the"
  echo "appropriate +build directive."
  exit 1
fi

if [[ -z "$GO_STUBS" ]]; then
  export GO_STUBS="web"
fi

for FILE in "$@"; do
  echo "Compiling $FILE..."
  protoc --go_out=. "$FILE"

  # Determine file names
  PB_FILE="${FILE%.proto}.pb.go"
  if [[ $AE_COPY -ne 0 ]]; then
    AE_FILE="${FILE%.proto}.ae.go"
    cp "$PB_FILE" "$AE_FILE"
  else
    AE_FILE="$PB_FILE"
  fi

  echo "Sanitizing $AE_FILE..."
  {
    echo "H"                    # Display human-readable errors, if any
    echo "g/goprotobuf/d"       # Delete lines containing goprotobuf
    echo "g/proto\./d"          # Delete lines calling into the library
    echo "w"                    # Write
    echo "q"                    # Quit
  } | ed -s "$AE_FILE"

  if [[ $AE_COPY -eq 0 ]]; then
    # If we are sharing the same proto, don't insert guards
    continue
  fi

  echo "Guarding $PB_FILE..."
  {
    echo "H"                    # Display human-readable errors, if any
    echo "1i"                   # Insert at the beginning of the file
    echo "// +build !appengine" # Don't compile this file under appengine
    echo                        # Blank line to not confuse anything
    echo "."                    # Exit insert mode
    echo "w"                    # Write
    echo "q"                    # Quit
  } | ed -s "$PB_FILE"

  echo "Guarding $AE_FILE..."
  {
    echo "H"                    # Display human-readable errors, if any
    echo "1i"                   # Insert at the beginning of the file
    echo "// +build appengine"  # Only compile this file under appengine
    echo                        # Blank line to not confuse anything
    echo "."                    # Exit insert mode
    echo "w"                    # Write
    echo "q"                    # Quit
  } | ed -s "$AE_FILE"
done
