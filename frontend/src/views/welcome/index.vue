<template>
  <div class="dashboard-container">
    <!-- 顶部欢迎与快速工作台 -->
    <WelcomeBanner
      :panel-data="dashboardData.panelGroup"
      :loading="loading"
      @refresh="handleRefresh"
    />

    <!-- 核心统计指标组 -->
    <PanelGroup
      :panel-data="dashboardData.panelGroup"
      @filter-chart="handleFilterChart"
    />

    <!-- 图表分析与内容生态区 -->
    <el-row :gutter="16" class="charts-container">
      <!-- 流量走势分析 -->
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="16"
        :xl="16"
        class="chart-col"
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 100 } }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <div class="header-left">
                <span class="header-icon">📈</span>
                <span class="chart-title">流量访问趋势</span>
              </div>
              <div class="header-right">
                <el-radio-group
                  v-model="selectedDays"
                  size="small"
                  class="days-radio-group"
                  @change="handleDaysChange"
                >
                  <el-radio-button :value="7">近7天</el-radio-button>
                  <el-radio-button :value="14">近14天</el-radio-button>
                  <el-radio-button :value="30">近30天</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>
          <el-skeleton animated :rows="7" :loading="loading">
            <template #default>
              <Line :line-data="dashboardData.lineChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 文章分类占比 -->
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="8"
        :xl="8"
        class="chart-col"
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 200 } }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <div class="header-left">
                <span class="header-icon">🎯</span>
                <span class="chart-title">文章分类生态</span>
              </div>
            </div>
          </template>
          <el-skeleton animated :rows="7" :loading="loading">
            <template #default>
              <Pie :pie-data="dashboardData.pieChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 近一周访问频次柱状图 -->
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="14"
        :xl="14"
        class="chart-col"
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 300 } }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <div class="header-left">
                <span class="header-icon">📊</span>
                <span class="chart-title">近7天访问频次</span>
              </div>
              <el-tag size="small" type="info" effect="plain" class="header-tag">
                每日 PV
              </el-tag>
            </div>
          </template>
          <el-skeleton animated :rows="6" :loading="loading">
            <template #default>
              <Bar :bar-data="dashboardData.barChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 热门博文 Top 5 阅读排行 -->
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="10"
        :xl="10"
        class="chart-col"
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 400 } }"
      >
        <el-card shadow="hover" class="chart-card top-posts-card">
          <template #header>
            <div class="chart-header">
              <div class="header-left">
                <span class="header-icon">🔥</span>
                <span class="chart-title">热门文章阅读排行</span>
              </div>
              <router-link to="/manage/post" class="more-link">
                更多 &rarr;
              </router-link>
            </div>
          </template>
          <el-skeleton animated :rows="6" :loading="loading">
            <template #default>
              <TopPostList :posts="dashboardData.topPosts" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 访客地理分布地图 -->
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="16"
        :xl="16"
        class="chart-col"
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 500 } }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <div class="header-left">
                <span class="header-icon">🗺️</span>
                <span class="chart-title">访客地域热力分布</span>
              </div>
              <el-tag size="small" type="primary" effect="plain" class="header-tag">
                国内地域聚合
              </el-tag>
            </div>
          </template>
          <el-skeleton animated :rows="9" :loading="loading">
            <template #default>
              <ChinaMap :map-data="dashboardData.mapChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 访客省份 Top 10 -->
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="8"
        :xl="8"
        class="chart-col"
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 600 } }"
      >
        <el-card shadow="hover" class="chart-card top10-card">
          <template #header>
            <div class="chart-header">
              <div class="header-left">
                <span class="header-icon">🏆</span>
                <span class="chart-title">访客 Top 10 省份</span>
              </div>
            </div>
          </template>
          <el-skeleton animated :rows="9" :loading="loading">
            <template #default>
              <div class="top10-list">
                <div
                  v-for="(item, index) in top10Provinces"
                  :key="index"
                  class="top10-item"
                  :style="{ animationDelay: `${index * 0.08}s` }"
                >
                  <div class="rank" :class="getRankClass(index)">
                    {{ index + 1 }}
                  </div>
                  <div class="province-name">{{ item.name }}</div>
                  <div class="visitor-count">
                    <span class="count">{{ item.value.toLocaleString() }}</span>
                    <span class="unit">人</span>
                  </div>
                  <div class="progress-bar">
                    <div
                      class="progress-fill"
                      :style="{ width: getProgressWidth(item.value) }"
                    />
                  </div>
                </div>
                <div v-if="top10Provinces.length === 0" class="no-data">
                  <el-empty description="暂无地域访问数据" :image-size="60" />
                </div>
              </div>
            </template>
          </el-skeleton>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import WelcomeBanner from "./components/WelcomeBanner.vue";
import PanelGroup from "./components/PanelGroup.vue";
import Line from "./components/Line.vue";
import Pie from "./components/Pie.vue";
import Bar from "./components/Bar.vue";
import TopPostList from "./components/TopPostList.vue";
import ChinaMap from "./components/ChinaMap.vue";
import { getDashboardData, type DashboardPanel, type TopPostItem } from "@/api/dashboard";
import { message } from "@/utils/message";

defineOptions({
  name: "Dashboard"
});

const loading = ref<boolean>(true);
const selectedDays = ref<number>(30);

// 仪表板数据结构
const dashboardData = ref<{
  panelGroup: DashboardPanel;
  lineChart: Array<{ name: string; pv: number; uv: number }>;
  pieChart: Array<{ name: string; value: number }>;
  barChart: Array<{ name: string; value: number }>;
  mapChart: Array<{ name: string; value: number }>;
  topPosts: Array<TopPostItem>;
}>({
  panelGroup: {
    postTotal: 0,
    categoryTotal: 0,
    tagTotal: 0,
    visit: 0,
    totalWords: 0,
    totalReadCount: 0,
    todayVisit: 0,
    yesterdayVisit: 0,
    visitGrowth: 0,
    weekPostTotal: 0,
    totalMusic: 0
  },
  lineChart: [],
  pieChart: [],
  barChart: [],
  mapChart: [],
  topPosts: []
});

// 获取仪表板数据
const fetchDashboardData = async (days: number = selectedDays.value) => {
  loading.value = true;
  try {
    const res = await getDashboardData({ days });
    if (res.code === 200 && res.payload) {
      dashboardData.value = {
        panelGroup: {
          postTotal: res.payload.panel?.postTotal || 0,
          categoryTotal: res.payload.panel?.categoryTotal || 0,
          tagTotal: res.payload.panel?.tagTotal || 0,
          visit: res.payload.panel?.visit || 0,
          totalWords: res.payload.panel?.totalWords || 0,
          totalReadCount: res.payload.panel?.totalReadCount || 0,
          todayVisit: res.payload.panel?.todayVisit || 0,
          yesterdayVisit: res.payload.panel?.yesterdayVisit || 0,
          visitGrowth: res.payload.panel?.visitGrowth || 0,
          weekPostTotal: res.payload.panel?.weekPostTotal || 0,
          totalMusic: res.payload.panel?.totalMusic || 0
        },
        lineChart: res.payload.line || [],
        pieChart: res.payload.pie || [],
        barChart: res.payload.bar || [],
        mapChart: res.payload.map || [],
        topPosts: res.payload.topPosts || []
      };
    }
  } catch (error) {
    console.error("获取仪表板数据异常:", error);
    message("获取仪表盘数据失败，请重试", { type: "error" });
  } finally {
    loading.value = false;
  }
};

const handleRefresh = () => {
  fetchDashboardData(selectedDays.value);
};

const handleDaysChange = (val: string | number | boolean | undefined) => {
  if (typeof val === "number") {
    fetchDashboardData(val);
  }
};

const handleFilterChart = (_type: string) => {
  // 卡片联动事件
};

onMounted(() => {
  fetchDashboardData();
});

// 中国省份列表
const chineseProvinces = [
  "北京",
  "天津",
  "上海",
  "重庆",
  "河北",
  "山西",
  "辽宁",
  "吉林",
  "黑龙江",
  "江苏",
  "浙江",
  "安徽",
  "福建",
  "江西",
  "山东",
  "河南",
  "湖北",
  "湖南",
  "广东",
  "海南",
  "四川",
  "贵州",
  "云南",
  "陕西",
  "甘肃",
  "青海",
  "台湾",
  "内蒙古",
  "广西",
  "西藏",
  "宁夏",
  "新疆",
  "香港",
  "澳门"
];

// Top 10 省份计算属性
const top10Provinces = computed(() => {
  if (!dashboardData.value.mapChart) return [];
  return [...dashboardData.value.mapChart]
    .filter(item => chineseProvinces.includes(item.name))
    .sort((a, b) => b.value - a.value)
    .slice(0, 10);
});

// 获取进度条宽度
const getProgressWidth = (value: number) => {
  const maxValue = top10Provinces.value[0]?.value || 1;
  return `${Math.min(100, Math.round((value / maxValue) * 100))}%`;
};

// 获取排名样式类
const getRankClass = (index: number) => {
  if (index === 0) return "rank-gold";
  if (index === 1) return "rank-silver";
  if (index === 2) return "rank-bronze";
  return "rank-normal";
};
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 120px);

  .charts-container {
    .chart-col {
      margin-bottom: 16px;

      .chart-card {
        border-radius: 12px;
        border: 1px solid var(--el-border-color-lighter);
        background: var(--el-bg-color-overlay);
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        height: 100%;

        &:hover {
          box-shadow: 0 10px 24px -4px rgba(0, 0, 0, 0.08);
          border-color: var(--el-color-primary-light-5);
        }

        :deep(.el-card__header) {
          padding: 14px 20px;
          border-bottom: 1px solid var(--el-border-color-extra-light);
          background-color: transparent;
        }

        :deep(.el-card__body) {
          padding: 18px 20px;
        }

        .chart-header {
          display: flex;
          align-items: center;
          justify-content: space-between;

          .header-left {
            display: flex;
            align-items: center;
            gap: 8px;

            .header-icon {
              font-size: 16px;
            }

            .chart-title {
              font-size: 15px;
              font-weight: 600;
              color: var(--el-text-color-primary);
            }
          }

          .header-right {
            display: flex;
            align-items: center;
          }

          .header-tag {
            font-size: 11px;
            border-radius: 4px;
          }

          .more-link {
            font-size: 12px;
            color: var(--el-color-primary);
            text-decoration: none;
            transition: opacity 0.2s ease;

            &:hover {
              opacity: 0.8;
            }
          }
        }
      }
    }
  }
}

.top-posts-card {
  :deep(.el-card__body) {
    height: 310px;
    overflow-y: auto;
  }
}

.top10-card {
  :deep(.el-card__body) {
    padding: 14px 16px;
    height: 480px;
    overflow-y: auto;
  }
}

.top10-list {
  display: flex;
  flex-direction: column;
  gap: 8px;

  .top10-item {
    display: flex;
    align-items: center;
    padding: 9px 12px;
    background: var(--el-fill-color-light);
    border-radius: 8px;
    border: 1px solid var(--el-border-color-extra-light);
    animation: slideIn 0.4s ease forwards;
    opacity: 0;
    transition: all 0.25s ease;

    &:hover {
      transform: translateX(4px);
      background: var(--el-fill-color);
      border-color: var(--el-color-primary-light-5);
    }

    .rank {
      width: 24px;
      height: 24px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 700;
      font-size: 12px;
      margin-right: 12px;
      flex-shrink: 0;

      &.rank-gold {
        background: linear-gradient(135deg, #ffd700, #ffaa00);
        color: #fff;
        box-shadow: 0 2px 8px rgba(255, 170, 0, 0.4);
      }

      &.rank-silver {
        background: linear-gradient(135deg, #b0bec5, #78909c);
        color: #fff;
        box-shadow: 0 2px 8px rgba(120, 144, 156, 0.35);
      }

      &.rank-bronze {
        background: linear-gradient(135deg, #d7ccc8, #a1887f);
        color: #fff;
        box-shadow: 0 2px 8px rgba(161, 136, 127, 0.35);
      }

      &.rank-normal {
        background: var(--el-fill-color-darker);
        color: var(--el-text-color-regular);
      }
    }

    .province-name {
      flex: 1;
      font-size: 13px;
      color: var(--el-text-color-primary);
      font-weight: 500;
    }

    .visitor-count {
      margin-right: 12px;
      text-align: right;

      .count {
        font-size: 14px;
        font-weight: 700;
        color: var(--el-color-primary);
      }

      .unit {
        font-size: 11px;
        color: var(--el-text-color-placeholder);
        margin-left: 2px;
      }
    }

    .progress-bar {
      width: 56px;
      height: 6px;
      background: var(--el-fill-color-darker);
      border-radius: 3px;
      overflow: hidden;

      .progress-fill {
        height: 100%;
        background: linear-gradient(
          90deg,
          var(--el-color-primary-light-3),
          var(--el-color-primary)
        );
        border-radius: 3px;
        transition: width 0.6s ease;
      }
    }
  }

  .no-data {
    text-align: center;
    padding: 40px 0;
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-12px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 12px;
  }
}
</style>
