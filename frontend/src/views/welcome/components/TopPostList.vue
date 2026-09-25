<template>
  <div class="top-posts-container">
    <div v-if="posts && posts.length > 0" class="posts-list">
      <div
        v-for="(item, index) in posts"
        :key="item.id"
        class="post-item"
        @click="handleClick(item.id)"
      >
        <div class="rank-badge" :class="getRankClass(index)">
          {{ index + 1 }}
        </div>
        <div class="post-info">
          <div class="post-title" :title="item.title">
            {{ item.title }}
          </div>
          <div class="post-meta">
            <span v-if="item.categoryName" class="category-tag">
              {{ item.categoryName }}
            </span>
            <span class="meta-date">{{ formatDate(item.createTime) }}</span>
          </div>
        </div>
        <div class="read-count">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
            <circle cx="12" cy="12" r="3" />
          </svg>
          <span>{{ item.readCount.toLocaleString() }}</span>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <el-empty description="暂无热门文章数据" :image-size="60" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import type { TopPostItem } from "@/api/dashboard";

defineOptions({
  name: "TopPostList"
});

const props = defineProps<{
  posts?: TopPostItem[];
}>();

const router = useRouter();

const getRankClass = (index: number) => {
  if (index === 0) return "rank-gold";
  if (index === 1) return "rank-silver";
  if (index === 2) return "rank-bronze";
  return "rank-normal";
};

const formatDate = (timestamp: number) => {
  if (!timestamp) return "";
  const d = new Date(timestamp * 1000);
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${m}-${day}`;
};

const handleClick = (id: number) => {
  router.push(`/manage/preview/${id}`);
};
</script>

<style lang="scss" scoped>
.top-posts-container {
  height: 100%;
  display: flex;
  flex-direction: column;

  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .post-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 14px;
    border-radius: 8px;
    background: var(--el-fill-color-light);
    border: 1px solid var(--el-border-color-extra-light);
    cursor: pointer;
    transition: all 0.25s ease;

    &:hover {
      background: var(--el-fill-color);
      transform: translateX(4px);
      border-color: var(--el-color-primary-light-5);

      .post-title {
        color: var(--el-color-primary);
      }
    }

    .rank-badge {
      width: 24px;
      height: 24px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 700;
      flex-shrink: 0;

      &.rank-gold {
        background: linear-gradient(135deg, #ffd700, #ffaa00);
        color: #fff;
        box-shadow: 0 2px 8px rgba(255, 170, 0, 0.35);
      }

      &.rank-silver {
        background: linear-gradient(135deg, #b0bec5, #78909c);
        color: #fff;
        box-shadow: 0 2px 8px rgba(120, 144, 156, 0.3);
      }

      &.rank-bronze {
        background: linear-gradient(135deg, #d7ccc8, #a1887f);
        color: #fff;
        box-shadow: 0 2px 8px rgba(161, 136, 127, 0.3);
      }

      &.rank-normal {
        background: var(--el-fill-color-darker);
        color: var(--el-text-color-secondary);
      }
    }

    .post-info {
      flex: 1;
      min-width: 0;

      .post-title {
        font-size: 13px;
        font-weight: 500;
        color: var(--el-text-color-primary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        transition: color 0.2s ease;
      }

      .post-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-top: 3px;

        .category-tag {
          font-size: 11px;
          padding: 1px 6px;
          border-radius: 4px;
          background: rgba(64, 158, 255, 0.1);
          color: var(--el-color-primary);
        }

        .meta-date {
          font-size: 11px;
          color: var(--el-text-color-placeholder);
        }
      }
    }

    .read-count {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      font-weight: 600;
      color: var(--el-text-color-secondary);
      flex-shrink: 0;

      svg {
        width: 14px;
        height: 14px;
        color: var(--el-text-color-placeholder);
      }
    }
  }

  .empty-state {
    padding: 30px 0;
  }
}
</style>
