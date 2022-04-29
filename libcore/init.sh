#!/bin/bash

source .github/env.sh

chmod -R 777 build 2>/dev/null
rm -rf build 2>/dev/null

go get -v -d

# Install gomobile
if [ ! -f "$GOPATH/bin/gomobile" ]; then
    go install -v golang.org/x/mobile/cmd/gomobile@latest
fi

gomobile init
