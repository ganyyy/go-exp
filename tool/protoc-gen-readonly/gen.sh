#!/usr/bin/env bash
protoc --experimental_allow_proto3_optional \
    --go_out=.. \
    --readonly_out=.. \
    --readonly_opt=readonly_pkg="protoc-gen-readonly/readonly" \
    --readonly_opt=suffix="readonly" \
    -I. proto/**/*.proto #
