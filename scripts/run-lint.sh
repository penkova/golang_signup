#!/usr/bin/env bash

go get -u github.com/alecthomas/gometalinter
"${GOPATH}/bin/gometalinter" --install > /dev/null
"${GOPATH}/bin/gometalinter" ./... && echo "✅ Your code is awesome!!!"