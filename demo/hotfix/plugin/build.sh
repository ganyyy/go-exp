#!/bin/bash


# 查找目录下所有的 .so 文件并删除
find . -name "*.so" -exec rm -rf  {} \;

NAME="plugin_$(date +%s).so"

echo "${NAME}"

go build --buildmode=plugin -gcflags="ganyyy.com/go-exp/demo/hotfix/...=-l"  -trimpath -ldflags "-X 'main.Version=${NAME}'" -o "$NAME" plugin.go

echo "${NAME}" > ./version
