#!/bin/bash

set -e -o pipefail

if [ "$DISABLE_LINTER" == "true" ]
then
  exit 0
fi

linterVersion="$(golangci-lint --version | awk '{print $4}')"

if [[ ! "${linterVersion}" =~ ^v?1\.6[0-9] ]]; then
	echo "Install GolangCI-Lint version 1.6x.x"
  exit 1
fi

export GO111MODULE=on
golangci-lint run \
  --verbose \
  --build-tags build
