#!/usr/bin/env base

publish_directory="/tmp/publish"
while [ ! -e "$publish_directory" ]; do
    sleep 1
done
