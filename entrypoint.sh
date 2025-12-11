#!/bin/bash
set -e

CONFIG_DIR="/app/config"
DEFAULT_CONFIG_DIR="/app/default"

echo "检查配置文件..."

# 确保配置目录存在
mkdir -p "${CONFIG_DIR}"

# 检查主要配置文件是否存在
if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
    echo "初始化配置文件..."
    # 如果目录中有示例配置文件，复制它
    if [ -f "${DEFAULT_CONFIG_DIR}/config.yaml" ]; then
        cp "${DEFAULT_CONFIG_DIR}/config.yaml" "${CONFIG_DIR}/config.yaml"
        echo "已从示例文件创建配置"
    else
        echo "错误：找不到默认配置文件"
        echo "默认配置目录内容："
        find "${DEFAULT_CONFIG_DIR}" -type f
        exit 1
    fi

    # 设置适当的权限
    chmod 644 "${CONFIG_DIR}/config.yaml"
fi

# 验证配置文件
if [ ! -s "${CONFIG_DIR}/config.yaml" ]; then
    echo "错误：配置文件为空"
    exit 1
fi

echo "配置文件检查完成"

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
