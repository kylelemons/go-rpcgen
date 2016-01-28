#!/bin/bash

ROOT="."
RPCGEN=$(dirname "$0")
PREFIX="github.com/bradhe/go-rpcgen"

set -e

for DIR in "webrpc"; do
  echo "# Copying ${RPCGEN}/${DIR} into ${ROOT}/${PREFIX}/${DIR}..."
  mkdir -p "${ROOT}/${PREFIX}/${DIR}"
  cp -R "${RPCGEN}/${DIR}/"* "${ROOT}/${PREFIX}/${DIR}/"
done
