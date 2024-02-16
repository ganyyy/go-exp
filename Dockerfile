FROM golang:1.21.6

# 这个需要在实际构建的时候指定
ARG REPO_DIR="/invalid/path"

ENV GOPROXY=https://goproxy.cn,direct

VOLUME  ${REPO_DIR}:/home/go-exp

WORKDIR /home/

CMD [ "bash" ]