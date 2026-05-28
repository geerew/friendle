#!/bin/sh
VERSION=$(git describe --tags --exact-match 2>/dev/null || echo dev)
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)

go build \
  -tags dev \
  -ldflags "-X github.com/geerew/friendle/version.Version=$VERSION -X github.com/geerew/friendle/version.Commit=$COMMIT" \
  -o ./tmp/main \
  .