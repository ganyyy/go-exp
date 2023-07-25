#!/usr/bin/env bash

# 参数至少为一个, 可以是 update 或者 rollback

if [ $# -lt 1 ]; then
    echo "Usage: $0 <update|rollback>"
    exit 1
fi

# 获取进程id
PROCESS_ID=`pgrep hotfix`

# 如果进程id不存在, 则直接退出
if [ -z "$PROCESS_ID" ]; then
    echo "hotfix is not running"
    exit 1
fi

if [ "$1" == "update" ]; then
    # 发送 SIGUSR1 信号
    kill -SIGUSR1 $PROCESS_ID
elif [ "$1" == "rollback" ]; then
    # 发送 SIGUSR2 信号
    kill -SIGUSR2 $PROCESS_ID
else
    echo "Usage: $0 <update|rollback>"
    exit 1
fi