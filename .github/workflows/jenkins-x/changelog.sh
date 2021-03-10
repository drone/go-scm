#!/usr/bin/env sh

echo "REPO_NAME = $PULL_BASE_SHA"

export PULL_BASE_SHA=$(git rev-parse HEAD)

jx changelog create --verbose --version=$VERSION --rev=$PULL_BASE_SHA --output-markdown=changelog.md --update-release=false

