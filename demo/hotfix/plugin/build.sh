#!/bin/bash

# 查找目录下所有的 .so 文件并删除
find . -name "plugin_*.so*" -exec rm -rf {} \;

NAME="plugin_$(date +%s).so"

echo "${NAME}"


mv locate/locate.go locate/locate.go.bak

go run locate/gen/gen_locate.go > locate/locate.go

go build --buildmode=plugin -gcflags="" -trimpath -ldflags "-X 'main.Version=${NAME}'" -o "$NAME" plugin.go

mv locate/locate.go.bak locate/locate.go

go run ../enctext/main.go -input "$NAME" -prefix "ganyyy.com/go-exp/demo/hotfix/plugin/fix"

echo "${NAME}" >./version


if [ "$1" = "test" ]; then
    echo -e "Testing the plugin...\n\n\n\n"
    go run -trimpath main/main.go -plugin "$NAME"
else
    echo "Building the plugin..."
fi