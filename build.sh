#!/bin/bash

# KeepBlog Docker 构建脚本
# 自动获取版本信息

set -e

# 获取版本信息
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u '+%Y-%m-%d')

# 镜像名称
IMAGE_NAME="jieepre/keepblog"

echo "================================================"
echo "构建 KeepBlog Docker 镜像"
echo "================================================"
echo "版本: $VERSION"
echo "提交: $GIT_COMMIT"
echo "日期: $BUILD_DATE"
echo "================================================"

# 构建镜像
docker build \
    --build-arg VERSION="$VERSION" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_DATE="$BUILD_DATE" \
    -t "$IMAGE_NAME:$VERSION" \
    -t "$IMAGE_NAME:latest" \
    .

echo ""
echo "✅ 构建完成!"
echo "   镜像: $IMAGE_NAME:$VERSION"
echo "   镜像: $IMAGE_NAME:latest"
