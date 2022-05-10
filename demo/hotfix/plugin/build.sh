#!/bin/bash

NAME="plugin_`date +%s`.so"

find . -name "*.so" | xargs rm -f

echo $NAME

go build --buildmode=plugin -ldflags "-X 'main.Version=${NAME}'" -o $NAME plugin.go

echo "${NAME}" > ./version