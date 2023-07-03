#!/usr/bin/env bash

# 构建可执行文件

ELF_FILE_NAME="ebpf_trace"

go build -o ${ELF_FILE_NAME}

# 输出可执行文件的绝对路径

ELF_FILE_PATH=$(pwd)/${ELF_FILE_NAME}

echo "ELF_FILE_PATH: ${ELF_FILE_PATH}"

# 使用 sed -i 命令修改 ./trace/trace_fmt.py 文件中的 ELF_FILE_PATH 变量

sed -i "s#ELF_FILE_PATH = .*#ELF_FILE_PATH = b\"${ELF_FILE_PATH}\"#g" ./trace/trace_fmt.py
