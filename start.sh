#!bash


docker run \
    -v $(pwd):/home/go-exp \
    --detach \
    --interactive \
    --name=go-exp my-go-exp:latest