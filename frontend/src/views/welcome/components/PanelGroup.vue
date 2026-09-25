<template>
  <div class="panel-group">
    <!-- 第一梯队：核心流量与内容重点 -->
    <div class="stat-grid primary-grid">
      <!-- 今日访问 -->
      <div
        class="stat-card group"
        @click="emit('filterChart', 'today')"
      >
        <div class="stat-card-header">
          <span class="stat-title">今日访问</span>
          <div class="icon-box icon-cyan">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </div>
        </div>
        <div class="stat-card-body">
          <div class="stat-value">
            <count-to
              :start-val="0"
              :end-val="panelData.todayVisit"
              :duration="2200"
              separator=","
            />
            <span class="stat-unit">人次</span>
          </div>
          <div
            v-if="panelData.yesterdayVisit !== undefined"
            class="stat-badge"
            :class="growthClass"
          >
            <span class="badge-arrow">{{ (panelData.visitGrowth ?? 0) >= 0 ? '↑' : '↓' }}</span>
            <span>{{ Math.abs(panelData.visitGrowth ?? 0) }}%</span>
          </div>
        </div>
        <div class="stat-card-footer">
          <span class="footer-desc">
            昨日访问: {{ panelData.yesterdayVisit ?? '-' }} 人次
          </span>
        </div>
      </div>

      <!-- 总访问量 -->
      <div
        class="stat-card group"
        @click="emit('filterChart', 'visits')"
      >
        <div class="stat-card-header">
          <span class="stat-title">总访问量</span>
          <div class="icon-box icon-blue">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="2" y1="12" x2="22" y2="12" />
              <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
            </svg>
          </div>
        </div>
        <div class="stat-card-body">
          <div class="stat-value">
            <count-to
              :start-val="0"
              :end-val="panelData.visit"
              :duration="2600"
              separator=","
            />
            <span class="stat-unit">IP</span>
          </div>
        </div>
        <div class="stat-card-footer">
          <span class="footer-desc">累计独立访客沉淀</span>
        </div>
      </div>

      <!-- 文章总数 -->
      <div
        class="stat-card group is-clickable"
        @click="router.push('/manage/post')"
      >
        <div class="stat-card-header">
          <span class="stat-title">文章总数</span>
          <div class="icon-box icon-emerald">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14 2 14 8 20 8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
              <polyline points="10 9 9 9 8 9" />
            </svg>
          </div>
        </div>
        <div class="stat-card-body">
          <div class="stat-value">
            <count-to
              :start-val="0"
              :end-val="panelData.postTotal"
              :duration="2000"
              separator=","
            />
            <span class="stat-unit">篇</span>
          </div>
          <div v-if="panelData.weekPostTotal" class="stat-badge badge-green">
            本周 +{{ panelData.weekPostTotal }}
          </div>
        </div>
        <div class="stat-card-footer">
          <span class="footer-desc flex-between">
            <span>已发布技术博文</span>
            <span class="link-hint">管理 &rarr;</span>
          </span>
        </div>
      </div>

      <!-- 总阅读量 -->
      <div class="stat-card group">
        <div class="stat-card-header">
          <span class="stat-title">总阅读量</span>
          <div class="icon-box icon-amber">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
              <circle cx="12" cy="12" r="3" />
            </svg>
          </div>
        </div>
        <div class="stat-card-body">
          <div class="stat-value">
            <count-to
              :start-val="0"
              :end-val="panelData.totalReadCount"
              :duration="2400"
              separator=","
            />
            <span class="stat-unit">次</span>
          </div>
        </div>
        <div class="stat-card-footer">
          <span class="footer-desc">文章累计受关注度</span>
        </div>
      </div>
    </div>

    <!-- 第二梯队：资产与聚合指标 -->
    <div class="stat-grid secondary-grid">
      <!-- 总字数 -->
      <div class="stat-card compact-card">
        <div class="compact-left">
          <div class="icon-box icon-purple">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="4 7 4 4 20 4 20 7" />
              <line x1="9" y1="20" x2="15" y2="20" />
              <line x1="12" y1="4" x2="12" y2="20" />
            </svg>
          </div>
          <div class="compact-info">
            <div class="compact-title">总字数</div>
            <div class="compact-value">
              <count-to
                :start-val="0"
                :end-val="panelData.totalWords"
                :duration="2400"
                separator=","
              />
            </div>
          </div>
        </div>
        <div class="compact-desc">持续思考的文字沉淀</div>
      </div>

      <!-- 分类数 -->
      <div
        class="stat-card compact-card is-clickable"
        @click="router.push('/manage/category')"
      >
        <div class="compact-left">
          <div class="icon-box icon-indigo">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
            </svg>
          </div>
          <div class="compact-info">
            <div class="compact-title">全部分类</div>
            <div class="compact-value">
              <count-to
                :start-val="0"
                :end-val="panelData.categoryTotal"
                :duration="1800"
              />
              <span class="compact-unit">个</span>
            </div>
          </div>
        </div>
        <div class="compact-desc">架构知识体系</div>
      </div>

      <!-- 标签数 -->
      <div
        class="stat-card compact-card is-clickable"
        @click="router.push('/manage/tag')"
      >
        <div class="compact-left">
          <div class="icon-box icon-rose">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
              <line x1="7" y1="7" x2="7.01" y2="7" />
            </svg>
          </div>
          <div class="compact-info">
            <div class="compact-title">标签聚合</div>
            <div class="compact-value">
              <count-to
                :start-val="0"
                :end-val="panelData.tagTotal"
                :duration="1800"
              />
              <span class="compact-unit">个</span>
            </div>
          </div>
        </div>
        <div class="compact-desc">技术脉络速查</div>
      </div>

      <!-- 音乐收藏 -->
      <div
        class="stat-card compact-card is-clickable"
        @click="router.push('/manage/music')"
      >
        <div class="compact-left">
          <div class="icon-box icon-teal">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 18V5l12-2v13" />
              <circle cx="6" cy="18" r="3" />
              <circle cx="18" cy="16" r="3" />
            </svg>
          </div>
          <div class="compact-info">
            <div class="compact-title">背景音乐</div>
            <div class="compact-value">
              <count-to
                :start-val="0"
                :end-val="panelData.totalMusic"
                :duration="1600"
              />
              <span class="compact-unit">首</span>
            </div>
          </div>
        </div>
        <div class="compact-desc">博文律动旋律</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { CountTo } from "vue3-count-to";
import type { DashboardPanel } from "@/api/dashboard";

defineOptions({
  name: "PanelGroup"
});

const props = defineProps<{
  panelData: DashboardPanel;
}>();

const emit = defineEmits<{
  (e: "filterChart", type: string): void;
}>();

const router = useRouter();

const growthClass = computed(() => {
  const g = props.panelData?.visitGrowth ?? 0;
  if (g > 0) return "badge-up";
  if (g < 0) return "badge-down";
  return "badge-flat";
});
</script>

<style lang="scss" scoped>
.panel-group {
  margin-bottom: 20px;

  .stat-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;

    &.secondary-grid {
      margin-top: 16px;
    }
  }

  .stat-card {
    position: relative;
    padding: 20px;
    border-radius: 12px;
    background: var(--el-bg-color-overlay);
    border: 1px solid var(--el-border-color-lighter);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
    transition: all 0.28s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;

    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 10px 24px -4px rgba(0, 0, 0, 0.08);
      border-color: var(--el-color-primary-light-5);
    }

    &.is-clickable {
      cursor: pointer;
    }

    .stat-card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .stat-title {
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-secondary);
      }
    }

    .icon-box {
      width: 42px;
      height: 42px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: transform 0.25s ease;

      svg {
        width: 22px;
        height: 22px;
      }

      &.icon-cyan {
        background: rgba(64, 201, 198, 0.12);
        color: #17b3b0;
      }

      &.icon-blue {
        background: rgba(54, 163, 247, 0.12);
        color: #268fe6;
      }

      &.icon-emerald {
        background: rgba(103, 194, 58, 0.12);
        color: #52a829;
      }

      &.icon-amber {
        background: rgba(230, 162, 60, 0.12);
        color: #d68714;
      }

      &.icon-purple {
        background: rgba(155, 89, 182, 0.12);
        color: #8e44ad;
      }

      &.icon-indigo {
        background: rgba(84, 110, 122, 0.12);
        color: #455a64;
      }

      &.icon-rose {
        background: rgba(244, 81, 108, 0.12);
        color: #e03855;
      }

      &.icon-teal {
        background: rgba(52, 191, 163, 0.12);
        color: #2ba28a;
      }
    }

    &:hover .icon-box {
      transform: scale(1.08);
    }

    .stat-card-body {
      display: flex;
      align-items: baseline;
      justify-content: space-between;
      margin-top: 14px;

      .stat-value {
        font-size: 26px;
        font-weight: 700;
        color: var(--el-text-color-primary);
        letter-spacing: -0.5px;
        line-height: 1.1;

        .stat-unit {
          font-size: 13px;
          font-weight: 400;
          color: var(--el-text-color-placeholder);
          margin-left: 4px;
        }
      }

      .stat-badge {
        font-size: 12px;
        font-weight: 600;
        padding: 2px 8px;
        border-radius: 12px;
        display: inline-flex;
        align-items: center;
        gap: 2px;

        &.badge-up {
          background: rgba(103, 194, 58, 0.12);
          color: #67c23a;
        }

        &.badge-down {
          background: rgba(245, 108, 108, 0.12);
          color: #f56c6c;
        }

        &.badge-flat {
          background: var(--el-fill-color);
          color: var(--el-text-color-placeholder);
        }

        &.badge-green {
          background: rgba(103, 194, 58, 0.12);
          color: #67c23a;
        }
      }
    }

    .stat-card-footer {
      margin-top: 10px;
      padding-top: 10px;
      border-top: 1px dashed var(--el-border-color-extra-light);

      .footer-desc {
        font-size: 12px;
        color: var(--el-text-color-placeholder);

        &.flex-between {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .link-hint {
          color: var(--el-color-primary);
          font-weight: 500;
        }
      }
    }
  }

  // 紧凑卡片样式
  .compact-card {
    padding: 16px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    .compact-left {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .compact-info {
      .compact-title {
        font-size: 13px;
        color: var(--el-text-color-secondary);
      }

      .compact-value {
        font-size: 20px;
        font-weight: 700;
        color: var(--el-text-color-primary);
        line-height: 1.2;
        margin-top: 2px;

        .compact-unit {
          font-size: 12px;
          font-weight: 400;
          color: var(--el-text-color-placeholder);
          margin-left: 2px;
        }
      }
    }

    .compact-desc {
      font-size: 12px;
      color: var(--el-text-color-placeholder);
      margin-top: 12px;
      padding-top: 8px;
      border-top: 1px dashed var(--el-border-color-extra-light);
    }
  }
}

// 响应式布局
@media (max-width: 1200px) {
  .panel-group {
    .stat-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
}

@media (max-width: 640px) {
  .panel-group {
    .stat-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 10px;
    }

    .stat-card {
      padding: 14px;

      .stat-card-body {
        margin-top: 8px;

        .stat-value {
          font-size: 20px;
        }
      }

      .icon-box {
        width: 34px;
        height: 34px;

        svg {
          width: 18px;
          height: 18px;
        }
      }
    }

    .compact-card {
      padding: 12px;
    }
  }
}
</style>
