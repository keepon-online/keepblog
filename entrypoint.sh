#!/bin/bash
set -e

echo "检查配置文件..."
if [ ! -f /app/config.yaml ]; then
    echo "初始化配置文件..."
    # 确保源文件存在
    if [ -f /app/default/config.yaml ]; then
        cp /app/default/config.yaml /app/config.yaml
    else
        echo "错误：找不到默认配置文件 /app/default/config.yaml"
        ls -lR /app/default  # 调试：列出目录内容
        exit 1
    fi
fi

echo "检查IP数据库..."
if [ ! -f /app/data/ip2region.xdb ]; then
    echo "初始化IP数据库..."
    # 确保源文件存在
    if [ -f /app/default/ip2region.xdb ]; then
        cp /app/default/ip2region.xdb /app/data/
    else
        echo "错误：找不到默认IP数据库 /app/default/ip2region.xdb"
        ls -lR /app/default  # 调试：列出目录内容
        exit 1
    fi
fi

echo "确保日志目录权限..."
mkdir -p /app/logs
chmod -R 755 /app/logs

echo "启动应用程序: $@"
exec ./site "$@"
