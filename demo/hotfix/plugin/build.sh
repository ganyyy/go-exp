#!/bin/bash

NAME="plugin_`date +%s`.so"

find . -name "*.so" | xargs rm -f

echo $NAME

go build --buildmode=plugin  -o $NAME plugin.go

echo "./plugin/${NAME}" > ./version