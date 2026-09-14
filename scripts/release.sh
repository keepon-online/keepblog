#!/bin/bash
#
# KeepBlog 发版脚本（git tag 驱动）
#
# 用法：
#   sh scripts/release.sh v5           # 发版 v5：校验 → 打 annotated tag → push → 构建
#   sh scripts/release.sh v5 --no-build # 只打 tag 和 push，不构建镜像
#   sh scripts/release.sh              # 不带参数：提示下一个版本号
#
# 版本号约定：v + 数字（v1, v2, ... v5）
# 打 tag 后 build.sh 输出干净版本号（如 v5），tag 之后的开发提交会带 -N-g<hash> 后缀

set -e

cd "$(git rev-parse --show-toplevel)"

# ---------- 参数解析 ----------
VERSION=""
NO_BUILD=0
if [ $# -eq 0 ]; then
    # 无参数：提示下一个建议版本号
    LATEST=$(git tag -l 'v*' | tr -d 'v' | sort -n | tail -1)
    if [ -z "$LATEST" ]; then
        echo "当前无版本 tag，建议从 v1 开始"
    else
        echo "当前最新 tag: v$LATEST，建议下一个版本: v$((LATEST + 1))"
    fi
    echo ""
    echo "用法: sh scripts/release.sh vN [--no-build]"
    exit 0
fi

for arg in "$@"; do
    case "$arg" in
        --no-build) NO_BUILD=1 ;;
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

# ---------- 构建镜像 ----------
IMAGE_NAME="jieepre/keepblog"
if [ "$NO_BUILD" -eq 1 ]; then
    echo "⏭️  --no-build，跳过镜像构建与推送"
else
    echo ""
    echo "🔨 构建镜像..."
    sh build.sh

    echo ""
    echo "⬆️  推送镜像到 Docker Hub..."
    docker push "$IMAGE_NAME:$VERSION"
    docker push "$IMAGE_NAME:latest"
    echo "✅ 已推送 $IMAGE_NAME:$VERSION 和 $IMAGE_NAME:latest"
fi

echo ""
echo "🎉 发版完成: $VERSION"
