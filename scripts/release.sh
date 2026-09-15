#!/bin/bash
#
# KeepBlog 发版脚本（git tag 驱动）
#
# 用法：
#   sh scripts/release.sh v5           # 发版 v5：校验 → 打 annotated tag → push（Gitee + GitHub）
#   sh scripts/release.sh              # 不带参数：提示下一个版本号
#   镜像由 GitHub Actions 在 tag push 后自动构建并推送 Docker Hub
#
# 版本号约定：v + 数字（v1, v2, ... v5）
# 打 tag 后 Actions 构建注入干净版本号（如 v5），tag 之后的开发提交会带 -N-g<hash> 后缀

set -e

cd "$(git rev-parse --show-toplevel)"

# ---------- 参数解析 ----------
VERSION=""
if [ $# -eq 0 ]; then
    # 无参数：提示下一个建议版本号
    LATEST=$(git tag -l 'v*' | tr -d 'v' | sort -n | tail -1)
    if [ -z "$LATEST" ]; then
        echo "当前无版本 tag，建议从 v1 开始"
    else
        echo "当前最新 tag: v$LATEST，建议下一个版本: v$((LATEST + 1))"
    fi
    echo ""
    echo "用法: sh scripts/release.sh vN"
    exit 0
fi

for arg in "$@"; do
    case "$arg" in
        v*)         VERSION="$arg" ;;
        *)          echo "❌ 无效参数: $arg"; exit 1 ;;
    esac
done

if [ -z "$VERSION" ]; then
    echo "❌ 必须指定版本号，例如: sh scripts/release.sh v5"
    exit 1
fi

# ---------- 前置校验 ----------
# 1. 工作区必须干净（避免把未提交改动混入版本点）
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ 工作区不干净，请先 commit 或 stash："
    git status --short
    exit 1
fi

# 2. 必须在分支 HEAD 上打 tag（HEAD 是当前分支最新提交）
BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "📌 当前分支: $BRANCH"

# 3. 版本号格式校验：v + 数字
if ! echo "$VERSION" | grep -Eq '^v[0-9]+$'; then
    echo "❌ 版本号格式错误: $VERSION（应为 v + 数字，如 v5）"
    exit 1
fi

# 4. tag 不能已存在
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "❌ tag $VERSION 已存在，指向: $(git rev-list -n1 --oneline $VERSION)"
    echo "   如需重发版，先删除: git tag -d $VERSION && git push origin :refs/tags/$VERSION"
    exit 1
fi

# 5. 提示是否落后于远程（避免在旧 commit 上发版）
REMOTE_REF="origin/$BRANCH"
if git rev-parse --verify -q "$REMOTE_REF" >/dev/null; then
    BEHIND=$(git rev-list --count HEAD.."$REMOTE_REF" 2>/dev/null || echo 0)
    if [ "$BEHIND" -gt 0 ]; then
        echo "⚠️  本地落后远程 $BEHIND 个提交，建议先 git pull"
        read -p "继续发版? [y/N] " confirm
        [ "$confirm" != "y" ] && { echo "已取消"; exit 1; }
    fi
fi

COMMIT=$(git rev-parse --short HEAD)
echo "================================================"
echo "🚀 准备发版 $VERSION"
echo "   分支:   $BRANCH"
echo "   commit: $COMMIT"
echo "   主题:   $(git log -1 --format=%s)"
echo "================================================"
read -p "确认发版? [y/N] " confirm
[ "$confirm" != "y" ] && { echo "已取消"; exit 1; }

# ---------- 打 tag ----------
git tag -a "$VERSION" -m "Release $VERSION

commit: $COMMIT
branch: $BRANCH
date:   $(date -u '+%Y-%m-%d %H:%M:%S UTC')"

echo "✅ 已打 tag $VERSION → $COMMIT"

# ---------- push ----------
echo "⬆️  推送 tag 到 origin..."
git push origin "$VERSION"
echo "✅ 已推送 $VERSION"

# 同步主分支与 tag 到 GitHub 镜像仓库，由 Actions 构建并推送 Docker 镜像
if git remote | grep -q '^github$'; then
    echo "⬆️  同步到 GitHub 镜像仓库（Actions 将构建镜像并推送 Docker Hub）..."
    git push github "$BRANCH" "$VERSION"
    echo "✅ 构建进度: https://github.com/keepon-online/keepblog/actions"
else
    echo "⚠️  未配置 github 远端，跳过镜像构建"
    echo "   配置: git remote add github https://github.com/keepon-online/keepblog.git"
fi

echo ""
echo "🎉 发版完成: $VERSION（镜像由 GitHub Actions 构建推送，本地无需构建）"
