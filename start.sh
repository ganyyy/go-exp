#!/usr/bin/env bash


docker run \
    --platform linux/x86_64 \
    -v "$(pwd)":/home/go-exp \
    --detach \
    --interactive \
    --cap-add=NET_ADMIN \
    --name=go-exp ganyyy/my-go-exp:v2
