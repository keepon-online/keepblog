#!/bin/sh
set -e

CONFIG_DIR="/app/config"
DEFAULT_CONFIG_DIR="/app/default"

echo "检查配置文件..."

# 确保配置目录存在
mkdir -p "${CONFIG_DIR}"

# 检查主要配置文件是否存在
if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
    echo "初始化配置文件..."
    if [ -f "${DEFAULT_CONFIG_DIR}/config.yaml" ]; then
        cp "${DEFAULT_CONFIG_DIR}/config.yaml" "${CONFIG_DIR}/config.yaml"
        echo "已从示例文件创建配置"
    else
        echo "错误：找不到默认配置文件"
        exit 1
    fi
    chmod 644 "${CONFIG_DIR}/config.yaml"
fi

# 首次部署时自动生成随机 JWT 密钥，避免所有部署共享同一密钥。
# 已有配置（含挂载卷里的）不会被修改；也可用环境变量 JWT_SECRET 覆盖。
if ! grep -q '^jwt:' "${CONFIG_DIR}/config.yaml" && [ -z "${JWT_SECRET:-}" ]; then
    GENERATED_SECRET="$(head -c 48 /dev/urandom | base64 | tr -d '/+=\n' | head -c 43)"
    printf '\njwt:\n  secret: "%s"\n' "${GENERATED_SECRET}" >> "${CONFIG_DIR}/config.yaml"
    echo "已生成随机 JWT 密钥并写入配置文件"
fi

# 验证配置文件
if [ ! -s "${CONFIG_DIR}/config.yaml" ]; then
    echo "错误：配置文件为空"
    exit 1
fi

echo "配置文件检查完成"

echo "确保数据和日志目录..."
mkdir -p /app/data /app/logs
chmod -R 755 /app/logs 2>/dev/null || true

echo "启动应用程序..."
exec /app/keepblog "$@"
