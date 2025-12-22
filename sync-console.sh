#!/bin/bash

# ================================================
# 前端构建并同步到后端 console 目录
# ================================================

set -e

# 路径配置
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONSOLE_DIR="../console"
BACKEND_CONSOLE_DIR="$SCRIPT_DIR/static/console"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  前端构建 & 同步脚本${NC}"
echo -e "${GREEN}================================================${NC}"

# 检查前端目录
if [ ! -d "$CONSOLE_DIR" ]; then
    echo -e "${RED}错误: 找不到前端目录 $CONSOLE_DIR${NC}"
    exit 1
fi

# 进入前端目录
cd "$CONSOLE_DIR"
echo -e "${YELLOW}➜ 进入前端目录: $(pwd)${NC}"

# 检查是否安装依赖
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}➜ 安装依赖...${NC}"
    pnpm install
fi

# 构建前端
echo -e "${YELLOW}➜ 构建前端项目...${NC}"
pnpm build

# 检查构建结果
if [ ! -d "dist" ]; then
    echo -e "${RED}错误: 构建失败，找不到 dist 目录${NC}"
    exit 1
fi

# 返回后端目录
cd "$SCRIPT_DIR"

# 备份旧文件（可选）
if [ -d "$BACKEND_CONSOLE_DIR" ]; then
    echo -e "${YELLOW}➜ 清理旧的 console 文件...${NC}"
    # 保留 static.go 文件
    find "$BACKEND_CONSOLE_DIR" -mindepth 1 -not -name "static.go" -delete 2>/dev/null || true
fi

# 同步新文件
echo -e "${YELLOW}➜ 同步构建产物到后端...${NC}"
cp -r "$CONSOLE_DIR/dist/"* "$BACKEND_CONSOLE_DIR/"

# 显示结果
echo ""
echo -e "${GREEN}✅ 同步完成!${NC}"
echo -e "   源目录: $CONSOLE_DIR/dist"
echo -e "   目标目录: $BACKEND_CONSOLE_DIR"
echo ""
echo -e "${YELLOW}文件列表:${NC}"
ls -la "$BACKEND_CONSOLE_DIR"
