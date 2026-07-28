<template>
  <div class="dashboard-container">
    <PanelGroup :panel-data="dashboardData.panelGroup" />
    <el-row :gutter="20" class="charts-container">
      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="12"
        :xl="12"
        class="chart-col"
        :initial="{
          opacity: 0,
          y: 100
        }"
        :enter="{
          opacity: 1,
          y: 0,
          transition: {
            delay: 200
          }
        }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">近30天浏览统计</span>
            </div>
          </template>
          <el-skeleton animated :rows="7" :loading="loading">
            <template #default>
              <Line :line-data="dashboardData.lineChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="12"
        :xl="12"
        class="chart-col"
        :initial="{
          opacity: 0,
          y: 100
        }"
        :enter="{
          opacity: 1,
          y: 0,
          transition: {
            delay: 400
          }
        }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">文章分类占比</span>
            </div>
          </template>
          <el-skeleton animated :rows="7" :loading="loading">
            <template #default>
              <Pie :pie-data="dashboardData.pieChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="24"
        :lg="24"
        :xl="24"
        class="chart-col"
        :initial="{
          opacity: 0,
          y: 100
        }"
        :enter="{
          opacity: 1,
          y: 0,
          transition: {
            delay: 600
          }
        }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">近一周统计</span>
            </div>
          </template>
          <el-skeleton animated :rows="7" :loading="loading">
            <template #default>
              <Bar :bar-data="dashboardData.barChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="16"
        :lg="16"
        :xl="16"
        class="chart-col"
        :initial="{
          opacity: 0,
          y: 100
        }"
        :enter="{
          opacity: 1,
          y: 0,
          transition: {
            delay: 800
          }
        }"
      >
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">访客地理分布</span>
            </div>
          </template>
          <el-skeleton animated :rows="10" :loading="loading">
            <template #default>
              <ChinaMap :map-data="dashboardData.mapChart" />
            </template>
          </el-skeleton>
        </el-card>
      </el-col>

      <el-col
        v-motion
        :xs="24"
        :sm="24"
        :md="8"
        :lg="8"
        :xl="8"
        class="chart-col"
        :initial="{
          opacity: 0,
          y: 100
        }"
        :enter="{
          opacity: 1,
          y: 0,
          transition: {
            delay: 900
          }
        }"
      >
        <el-card shadow="hover" class="chart-card top10-card">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">🏆 访客Top10省份</span>
            </div>
          </template>
          <el-skeleton animated :rows="10" :loading="loading">
            <template #default>
              <div class="top10-list">
                <div
                  v-for="(item, index) in top10Provinces"
                  :key="index"
                  class="top10-item"
                  :style="{ animationDelay: `${index * 0.1}s` }"
                >
                  <div class="rank" :class="getRankClass(index)">
                    {{ index + 1 }}
                  </div>
                  <div class="province-name">{{ item.name }}</div>
                  <div class="visitor-count">
                    <span class="count">{{ item.value }}</span>
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
                  暂无数据
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
defineOptions({
  name: "Dashboard"
});
import PanelGroup from "./components/PanelGroup.vue";
import Bar from "./components/Bar.vue";
import Pie from "./components/Pie.vue";
import Line from "./components/Line.vue";
import ChinaMap from "./components/ChinaMap.vue";

import { ref, onMounted, computed } from "vue";
import { getDashboardData } from "@/api/dashboard";

const loading = ref<boolean>(true);

// 仪表板数据
const dashboardData = ref({
  panelGroup: {
    postTotal: 0,
    categoryTotal: 0,
    tagTotal: 0,
    visit: 0,
    totalWords: 0,
    totalReadCount: 0,
    todayVisit: 0,
    totalMusic: 0
  },
  lineChart: [],
  pieChart: [],
  barChart: [],
  mapChart: []
});

// 获取仪表板数据
const fetchDashboardData = async () => {
  try {
    const res = await getDashboardData();
    if (res.code === 200) {
      // 处理获取到的数据，注意API返回的数据结构
      dashboardData.value = {
        panelGroup: {
          postTotal: res.payload.panel.postTotal,
          categoryTotal: res.payload.panel.categoryTotal,
          tagTotal: res.payload.panel.tagTotal,
          visit: res.payload.panel.visit,
          totalWords: res.payload.panel.totalWords || 0,
          totalReadCount: res.payload.panel.totalReadCount || 0,
          todayVisit: res.payload.panel.todayVisit || 0,
          totalMusic: res.payload.panel.totalMusic || 0
        },
        lineChart: res.payload.line,
        pieChart: res.payload.pie,
        barChart: res.payload.bar,
        mapChart: res.payload.map || []
      };
    } else {
      console.error("获取仪表板数据失败:", res.message);
    }
  } catch (error) {
    console.error("获取仪表板数据异常:", error);
  } finally {
    // 数据加载完成后隐藏骨架屏
    loading.value = false;
  }
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

// Top10省份计算属性（只显示国内省份）
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
  return `${(value / maxValue) * 100}%`;
};

// 获取排名样式类
const getRankClass = (index: number) => {
  if (index === 0) return "rank-gold";
  if (index === 1) return "rank-silver";
  if (index === 2) return "rank-bronze";
  return "";
};
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);

  .charts-container {
    margin-top: 20px;

    .chart-col {
      margin-bottom: 20px;

      .chart-card {
        border-radius: 8px;
        overflow: hidden;
        border: 1px solid var(--el-border-color-light);
        transition: all 0.3s ease;

        &:hover {
          box-shadow: var(--el-box-shadow-light);
          transform: translateY(-2px);
        }

        :deep(.el-card__header) {
          padding: 15px 20px;
          background-color: var(--el-fill-color-light);
          border-bottom: 1px solid var(--el-border-color-light);
        }

        .chart-header {
          display: flex;
          align-items: center;

          .chart-title {
            font-size: 16px;
            font-weight: 600;
            color: var(--el-text-color-primary);
          }
        }
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .dashboard-container {
    padding: 12px;
  }
}

// Top10卡片样式
.top10-card {
  height: 100%;

  :deep(.el-card__body) {
    padding: 16px;
    height: calc(100% - 56px);
    overflow-y: auto;
  }
}

.top10-list {
  .top10-item {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    margin-bottom: 8px;
    background: var(--el-fill-color-light);
    border-radius: 8px;
    animation: slideIn 0.5s ease forwards;
    opacity: 0;
    transition: all 0.3s ease;

    &:hover {
      transform: translateX(5px);
      background: var(--el-fill-color);
    }

    .rank {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: bold;
      font-size: 14px;
      background: var(--el-fill-color-darker);
      color: var(--el-text-color-regular);
      margin-right: 12px;
      flex-shrink: 0;

      &.rank-gold {
        background: linear-gradient(135deg, #ffd700, #ffb700);
        color: #fff;
        box-shadow: 0 2px 8px rgba(255, 215, 0, 0.4);
      }

      &.rank-silver {
        background: linear-gradient(135deg, #c0c0c0, #a8a8a8);
        color: #fff;
        box-shadow: 0 2px 8px rgba(192, 192, 192, 0.4);
      }

      &.rank-bronze {
        background: linear-gradient(135deg, #cd7f32, #b8722e);
        color: #fff;
        box-shadow: 0 2px 8px rgba(205, 127, 50, 0.4);
      }
    }

    .province-name {
      flex: 1;
      font-size: 14px;
      color: var(--el-text-color-primary);
      font-weight: 500;
    }

    .visitor-count {
      margin-right: 12px;
      text-align: right;

      .count {
        font-size: 16px;
        font-weight: 600;
        color: var(--el-color-primary);
      }

      .unit {
        font-size: 12px;
        color: var(--el-text-color-secondary);
        margin-left: 2px;
      }
    }

    .progress-bar {
      width: 60px;
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
        transition: width 0.8s ease;
      }
    }
  }

  .no-data {
    text-align: center;
    padding: 40px 0;
    color: var(--el-text-color-secondary);
    font-size: 14px;
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
