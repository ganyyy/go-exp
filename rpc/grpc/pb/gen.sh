#!/usr/bin/env bash

protoc -I. --go_out=../ --go-grpc_out=../ --go_opt=paths=import --go-grpc_opt=paths=import *.proto