#/usr/bin/env bash

# 执行go-lint, 检查是否存在循环依赖
golangci-lint run --disable-all -E typecheck
