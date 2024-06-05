#!/bin/bash

set -e -o pipefail

if [ "$DISABLE_LINTER" == "true" ]
then
  exit 0
fi

linterVersion="$(golangci-lint --version | awk '{print $4}')"

if [[ ! "${linterVersion}" =~ ^1\.5[89] ]]; then
	echo "Install GolangCI-Lint version 1.58 or 1.59"
  exit 1
fi

export GO111MODULE=on
golangci-lint run \
  --verbose \
  --build-tags build
