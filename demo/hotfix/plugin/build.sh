#!/bin/bash


find . -name "*.so" -print0 | xargs rm

NAME="plugin_$(date +%s).so"

echo "${NAME}"

go build --buildmode=plugin -ldflags "-X 'main.Version=${NAME}'" -o "$NAME" plugin.go

echo "${NAME}" > ./version