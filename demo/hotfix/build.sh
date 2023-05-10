#!/usr/bin/env bash


# 在当前目录下执行go build并禁用函数内联
# 这里的影响的范围是当前目录下的所有go文件, 如果有其他目录下的go文件也需要禁用内联, 需要在对应的目录下执行go build -gcflags="all=-l"
go build -gcflags=-l  -trimpath -o hotfix

