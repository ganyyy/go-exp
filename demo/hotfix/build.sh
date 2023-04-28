#!/usr/bin/env bash


# 在当前目录下执行go build并禁用函数内联
go build -gcflags=-l  -trimpath -o hotfix

