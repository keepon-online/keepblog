<template>
  <div class="welcome-banner">
    <div class="banner-content">
      <div class="greeting-area">
        <el-avatar
          v-if="userAvatar"
          :size="64"
          :src="userAvatar"
          class="user-avatar"
        />
        <div v-else class="user-avatar-placeholder">
          <span>{{ (userName || "Admin").slice(0, 1).toUpperCase() }}</span>
        </div>

        <div class="greeting-text">
          <div class="greeting-title">
            <span class="greeting-time">{{ greetingTime }}</span>
            <span class="greeting-name">{{ userName }}</span>
            <span class="greeting-emoji">{{ greetingEmoji }}</span>
          </div>
          <div class="greeting-subtitle">
            {{ greetingQuote }}
          </div>
          <div class="badge-row">
            <div class="stat-badge">
              <span class="badge-dot" />
              <span>本周更新 {{ panelData?.weekPostTotal || 0 }} 篇</span>
            </div>
            <div class="stat-badge">
              <span class="badge-dot dot-cyan" />
              <span>今日访问 {{ panelData?.todayVisit || 0 }} 人</span>
            </div>
            <div class="stat-badge">
              <span class="badge-dot dot-purple" />
              <span>累计产出 {{ formatWords(panelData?.totalWords || 0) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="action-area">
        <el-button
          type="primary"
          class="action-btn"
          @click="router.push('/manage/editor')"
        >
          <svg
            class="btn-icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M12 20h9" />
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
          </svg>
          <span>写新文章</span>
        </el-button>

        <el-button
          class="action-btn-secondary"
          @click="router.push('/system/monitor')"
        >
          <svg
            class="btn-icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
            <line x1="8" y1="21" x2="16" y2="21" />
            <line x1="12" y1="17" x2="12" y2="21" />
          </svg>
          <span>系统监控</span>
        </el-button>

        <el-tooltip content="刷新仪表盘数据" placement="top">
          <el-button
            circle
            class="refresh-btn"
            :loading="loading"
            @click="emit('refresh')"
          >
            <svg
              class="refresh-icon"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <polyline points="23 4 23 10 17 10" />
              <polyline points="1 20 1 14 7 14" />
              <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" />
            </svg>
          </el-button>
        </el-tooltip>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/store/modules/user";
import type { DashboardPanel } from "@/api/dashboard";

defineOptions({
  name: "WelcomeBanner"
});

const props = defineProps<{
  panelData?: DashboardPanel;
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: "refresh"): void;
}>();

const router = useRouter();
const userStore = useUserStore();

const userName = computed(() => {
  return userStore.nickname || userStore.username || "管理员";
});

const userAvatar = computed(() => {
  return userStore.avatar;
});

// 根据当前时间计算时段问候
const timeGreeting = computed(() => {
  const hour = new Date().getHours();
  if (hour >= 5 && hour < 9) {
    return {
      greeting: "早上好",
      emoji: "🌅",
      quote: "新的一天开始了，愿你灵感涌动，笔耕不辍。"
    };
  } else if (hour >= 9 && hour < 12) {
    return {
      greeting: "上午好",
      emoji: "☀️",
      quote: "保持专注与热爱，每一行文字都在记录沉淀。"
    };
  } else if (hour >= 12 && hour < 14) {
    return {
      greeting: "中午好",
      emoji: "🍱",
      quote: "午餐时间到了，适当小憩让思维更清晰。"
    };
  } else if (hour >= 14 && hour < 19) {
    return {
      greeting: "下午好",
      emoji: "☕",
      quote: "泡一杯清茶或咖啡，享受充实的工作与创作时光。"
    };
  } else if (hour >= 19 && hour < 23) {
    return {
      greeting: "晚上好",
      emoji: "🌙",
      quote: "夜色宁静，是总结思考与整理心得的最佳时刻。"
    };
  } else {
    return {
      greeting: "夜深了",
      emoji: "🌌",
      quote: "夜深了，早点休息，明天又是精彩的一天。"
    };
  }
});

const greetingTime = computed(() => timeGreeting.value.greeting);
const greetingEmoji = computed(() => timeGreeting.value.emoji);
const greetingQuote = computed(() => timeGreeting.value.quote);

const formatWords = (words: number) => {
  if (words >= 10000) {
    return (words / 10000).toFixed(1) + " 万字";
  }
  return words + " 字";
};
</script>

<style lang="scss" scoped>
.welcome-banner {
  position: relative;
  border-radius: 14px;
  margin-bottom: 20px;
  padding: 24px 28px;
  background: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(
      90deg,
      var(--el-color-primary) 0%,
      var(--el-color-primary-light-3) 50%,
      #36a3f7 100%
    );
  }

  &::after {
    content: "";
    position: absolute;
    right: -60px;
    bottom: -60px;
    width: 220px;
    height: 220px;
    border-radius: 50%;
    background: radial-gradient(
      circle,
      var(--el-color-primary-light-8) 0%,
      transparent 70%
    );
    opacity: 0.6;
    pointer-events: none;
  }

  .banner-content {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 20px;
  }

  .greeting-area {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .user-avatar {
    border: 3px solid var(--el-color-primary-light-7);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
    flex-shrink: 0;
  }

  .user-avatar-placeholder {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(
      135deg,
      var(--el-color-primary) 0%,
      var(--el-color-primary-light-3) 100%
    );
    color: #fff;
    font-size: 26px;
    font-weight: 700;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
    flex-shrink: 0;
  }

  .greeting-text {
    .greeting-title {
      font-size: 22px;
      font-weight: 700;
      color: var(--el-text-color-primary);
      display: flex;
      align-items: center;
      gap: 8px;
      line-height: 1.3;

      .greeting-name {
        background: linear-gradient(
          120deg,
          var(--el-color-primary) 0%,
          var(--el-color-primary-light-2) 100%
        );
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
      }
    }

    .greeting-subtitle {
      margin-top: 6px;
      font-size: 13px;
      color: var(--el-text-color-secondary);
      line-height: 1.5;
    }

    .badge-row {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-top: 10px;
      flex-wrap: wrap;

      .stat-badge {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        padding: 3px 10px;
        border-radius: 20px;
        background: var(--el-fill-color-light);
        color: var(--el-text-color-regular);
        border: 1px solid var(--el-border-color-extra-light);

        .badge-dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #67c23a;

          &.dot-cyan {
            background: #409eff;
          }

          &.dot-purple {
            background: #9b59b6;
          }
        }
      }
    }
  }

  .action-area {
    display: flex;
    align-items: center;
    gap: 12px;

    .action-btn {
      height: 40px;
      padding: 0 18px;
      border-radius: 8px;
      font-weight: 600;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      box-shadow: 0 4px 12px -2px var(--el-color-primary-light-5);
      transition: all 0.25s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 16px -2px var(--el-color-primary-light-3);
      }
    }

    .action-btn-secondary {
      height: 40px;
      padding: 0 18px;
      border-radius: 8px;
      font-weight: 500;
      background: var(--el-fill-color-light);
      border-color: var(--el-border-color-light);
      color: var(--el-text-color-primary);
      display: inline-flex;
      align-items: center;
      gap: 6px;
      transition: all 0.25s ease;

      &:hover {
        background: var(--el-fill-color);
        border-color: var(--el-color-primary-light-5);
        color: var(--el-color-primary);
        transform: translateY(-2px);
      }
    }

    .refresh-btn {
      width: 40px;
      height: 40px;
      border-radius: 8px;
      background: var(--el-fill-color-light);
      border-color: var(--el-border-color-light);
      color: var(--el-text-color-secondary);
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.25s ease;

      &:hover {
        background: var(--el-fill-color);
        color: var(--el-color-primary);
        transform: rotate(45deg);
      }
    }

    .btn-icon,
    .refresh-icon {
      width: 16px;
      height: 16px;
    }
  }
}

@media (max-width: 768px) {
  .welcome-banner {
    padding: 18px;

    .banner-content {
      flex-direction: column;
      align-items: flex-start;
    }

    .action-area {
      width: 100%;
      justify-content: flex-start;
    }
  }
}
</style>
