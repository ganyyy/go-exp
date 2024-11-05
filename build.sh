#!/usr/bin/env bash


docker buildx build --platform linux/x86_64 --load  --build-arg REPO_DIR="$(pwd)" -t ganyyy/my-go-exp:v2 . 
