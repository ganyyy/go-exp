#!/usr/bin/env bash


# docker compose -p nats2 down --volumes --remove-orphans
# docker compose -p nats2 up -d
# docker compose -p nats2 restart

# 必须输入一个参数, 可以是 start/stop/restart
if [ $# -ne 1 ]; then
    echo "Usage: $0 start|stop|restart"
    exit 1
fi

# 项目名称
PROJECT_NAME=nats2

# 判断输入的参数
case $1 in
    start)
        echo "Starting ${PROJECT_NAME}..."
        docker-compose -p ${PROJECT_NAME} up -d
        ;;
    stop)
        echo "Stopping ${PROJECT_NAME}..."
        docker-compose -p ${PROJECT_NAME} down --volumes --remove-orphans
        ;;
    restart)
        echo "Restarting ${PROJECT_NAME}..."
        docker-compose -p ${PROJECT_NAME} restart
        ;;
    *)
        echo "Usage: $0 start|stop|restart"
        exit 1
        ;;
esac