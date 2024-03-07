#!bash


docker run \
    -v $(pwd):/home/go-exp \
    --detach \
    --interactive \
    --cap-add=NET_ADMIN \
    --name=go-exp my-go-exp:latest
