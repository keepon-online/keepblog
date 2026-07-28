<script setup lang="ts">
defineOptions({
  name: "Monitor"
});
import { ref, onMounted, computed, onBeforeUnmount, watch } from "vue";
import { getServe } from "@/api/monitor";
import { ElMessage } from "element-plus";
import echarts from "@/utils/echarts"; // 按需引入
import {
  Monitor,
  Refresh,
  Cpu,
  Memo,
  TrendCharts,
  Platform,
  FolderOpened,
  Connection,
  Timer
} from "@element-plus/icons-vue";

interface ServerInfo {
  general?: {
    hostname?: string;
    os?: string;
    arch?: string;
    kernel?: string;
    uptime?: number;
    uptimeFormat?: string;
    cpuCores?: number;
    goVersion?: string;
  };
  cpu?: {
    usagePercent?: number;
    cores?: number;
    modelName?: string;
    frequency?: number;
    coreDetails?: number[];
  };
  memory?: {
    total?: number;
    used?: number;
    free?: number;
    available?: number;
    usedPercent?: number;
    totalFormat?: string;
    usedFormat?: string;
    freeFormat?: string;
    swapTotal?: number;
    swapUsed?: number;
    swapFree?: number;
    swapUsedPercent?: number;
    swapTotalFormat?: string;
    swapUsedFormat?: string;
  };
  disk?: Array<{
    device: string;
    mountpoint: string;
    fstype: string;
    total: number;
    used: number;
    free: number;
    usedPercent: number;
    totalFormat: string;
    usedFormat: string;
    freeFormat: string;
  }>;
  network?: Array<{
    name: string;
    bytesRecv: number;
    bytesSent: number;
    packetsRecv: number;
    packetsSent: number;
    recvFormat: string;
    sentFormat: string;
    isUp: boolean;
  }>;
  load?: {
    load1?: number;
    load5?: number;
    load15?: number;
  };
  process?: {
    total?: number;
    running?: number;
    sleeping?: number;
    stopped?: number;
    zombie?: number;
  };
  timestamp?: number;
}

const server = ref<ServerInfo>({});
const loading = ref(true);
const refreshInterval = ref<NodeJS.Timeout | null>(null);

// 历史数据存储
const historyData = ref<{
  timestamps: string[];
  cpuData: number[];
  memoryData: number[];
}>({
  timestamps: [],
  cpuData: [],
  memoryData: []
});

// 图表引用
const cpuGaugeRef = ref<HTMLElement | null>(null);
const memGaugeRef = ref<HTMLElement | null>(null);
const historyChartRef = ref<HTMLElement | null>(null);
let cpuGaugeChart: echarts.ECharts | null = null;
let memGaugeChart: echarts.ECharts | null = null;
let historyChart: echarts.ECharts | null = null;

// CPU核心显示控制
const showAllCores = ref(false);
const displayedCores = computed(() => {
  const cores = server.value.cpu?.coreDetails || [];
  if (showAllCores.value || cores.length <= 8) {
    return cores;
  }
  return cores.slice(0, 8);
});

// 计算CPU使用状态颜色
const cpuUsageClass = computed(() => {
  const usage = server.value.cpu?.usagePercent || 0;
  if (usage > 80) return "text-red-500";
  if (usage > 60) return "text-orange-500";
  return "text-green-500";
});

// 计算内存使用状态颜色
const memUsageClass = computed(() => {
  const usage = server.value.memory?.usedPercent || 0;
  if (usage > 80) return "text-red-500";
  if (usage > 60) return "text-orange-500";
  return "text-green-500";
});

// 获取使用率颜色
const getUsageColor = (percent: number) => {
  if (percent > 80) return "#f56c6c";
  if (percent > 60) return "#e6a23c";
  return "#67c23a";
};

// 初始化仪表盘图表
const initGaugeCharts = () => {
  if (cpuGaugeRef.value) {
    cpuGaugeChart = echarts.init(cpuGaugeRef.value);
    updateCpuGauge();
  }
  if (memGaugeRef.value) {
    memGaugeChart = echarts.init(memGaugeRef.value);
    updateMemGauge();
  }
  if (historyChartRef.value) {
    historyChart = echarts.init(historyChartRef.value);
    updateHistoryChart();
  }
  window.addEventListener("resize", handleResize);
};

const handleResize = () => {
  cpuGaugeChart?.resize();
  memGaugeChart?.resize();
  historyChart?.resize();
};

// 更新CPU仪表盘
const updateCpuGauge = () => {
  if (!cpuGaugeChart) return;
  const value = server.value.cpu?.usagePercent || 0;
  const option = {
    series: [
      {
        type: "gauge",
        startAngle: 200,
        endAngle: -20,
        min: 0,
        max: 100,
        splitNumber: 10,
        itemStyle: {
          color: getUsageColor(value)
        },
        progress: {
          show: true,
          roundCap: true,
          width: 12
        },
        pointer: {
          show: false
        },
        axisLine: {
          roundCap: true,
          lineStyle: {
            width: 12,
            color: [[1, "#e6e8f0"]]
          }
        },
        axisTick: { show: false },
        splitLine: { show: false },
        axisLabel: { show: false },
        title: {
          show: true,
          offsetCenter: [0, "30%"],
          fontSize: 14,
          color: "#666"
        },
        detail: {
          valueAnimation: true,
          offsetCenter: [0, "-10%"],
          fontSize: 28,
          fontWeight: "bold",
          formatter: "{value}%",
          color: getUsageColor(value)
        },
        data: [{ value: value.toFixed(1), name: "CPU" }]
      }
    ]
  };
  cpuGaugeChart.setOption(option);
};

// 更新内存仪表盘
const updateMemGauge = () => {
  if (!memGaugeChart) return;
  const value = server.value.memory?.usedPercent || 0;
  const option = {
    series: [
      {
        type: "gauge",
        startAngle: 200,
        endAngle: -20,
        min: 0,
        max: 100,
        splitNumber: 10,
        itemStyle: {
          color: getUsageColor(value)
        },
        progress: {
          show: true,
          roundCap: true,
          width: 12
        },
        pointer: {
          show: false
        },
        axisLine: {
          roundCap: true,
          lineStyle: {
            width: 12,
            color: [[1, "#e6e8f0"]]
          }
        },
        axisTick: { show: false },
        splitLine: { show: false },
        axisLabel: { show: false },
        title: {
          show: true,
          offsetCenter: [0, "30%"],
          fontSize: 14,
          color: "#666"
        },
        detail: {
          valueAnimation: true,
          offsetCenter: [0, "-10%"],
          fontSize: 28,
          fontWeight: "bold",
          formatter: "{value}%",
          color: getUsageColor(value)
        },
        data: [{ value: value.toFixed(1), name: "内存" }]
      }
    ]
  };
  memGaugeChart.setOption(option);
};

// 更新历史曲线图
const updateHistoryChart = () => {
  if (!historyChart) return;
  const option = {
    tooltip: {
      trigger: "axis",
      axisPointer: { type: "cross" }
    },
    legend: {
      data: ["CPU使用率", "内存使用率"],
      bottom: 0
    },
    grid: {
      left: "3%",
      right: "4%",
      bottom: "15%",
      top: "10%",
      containLabel: true
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: historyData.value.timestamps,
      axisLabel: { fontSize: 10 }
    },
    yAxis: {
      type: "value",
      min: 0,
      max: 100,
      axisLabel: { formatter: "{value}%" }
    },
    series: [
      {
        name: "CPU使用率",
        type: "line",
        smooth: true,
        symbol: "none",
        areaStyle: {
          opacity: 0.3,
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "rgba(59, 130, 246, 0.5)" },
            { offset: 1, color: "rgba(59, 130, 246, 0.05)" }
          ])
        },
        lineStyle: { width: 2, color: "#3b82f6" },
        data: historyData.value.cpuData
      },
      {
        name: "内存使用率",
        type: "line",
        smooth: true,
        symbol: "none",
        areaStyle: {
          opacity: 0.3,
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "rgba(16, 185, 129, 0.5)" },
            { offset: 1, color: "rgba(16, 185, 129, 0.05)" }
          ])
        },
        lineStyle: { width: 2, color: "#10b981" },
        data: historyData.value.memoryData
      }
    ]
  };
  historyChart.setOption(option);
};

// 更新历史数据
const updateHistoryData = () => {
  const now = new Date();
  const timeStr = now.toLocaleTimeString("zh-CN", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit"
  });

  historyData.value.timestamps.push(timeStr);
  historyData.value.cpuData.push(server.value.cpu?.usagePercent || 0);
  historyData.value.memoryData.push(server.value.memory?.usedPercent || 0);

  // 只保留最近30个数据点
  if (historyData.value.timestamps.length > 30) {
    historyData.value.timestamps.shift();
    historyData.value.cpuData.shift();
    historyData.value.memoryData.shift();
  }
};

// 获取服务器信息
const fetchServerInfo = async () => {
  try {
    loading.value = true;
    const rep = await getServe();
    if (rep && rep.code === 200) {
      server.value = rep.payload;
      updateHistoryData();
      updateCpuGauge();
      updateMemGauge();
      updateHistoryChart();
    } else {
      ElMessage.error("获取服务器信息失败");
    }
  } catch (error) {
    console.error("获取服务器信息失败:", error);
  } finally {
    loading.value = false;
  }
};

// 刷新数据
const refreshData = () => {
  fetchServerInfo();
  ElMessage.success("数据刷新成功");
};

// 自动刷新
const startAutoRefresh = () => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value);
  }
  refreshInterval.value = setInterval(fetchServerInfo, 5000); // 5秒刷新一次
};

const stopAutoRefresh = () => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value);
    refreshInterval.value = null;
  }
};

onMounted(() => {
  fetchServerInfo();
  setTimeout(initGaugeCharts, 100);
  startAutoRefresh();
});

// 组件卸载时清除定时器
onBeforeUnmount(() => {
  stopAutoRefresh();
  window.removeEventListener("resize", handleResize);
  cpuGaugeChart?.dispose();
  memGaugeChart?.dispose();
  historyChart?.dispose();
});
</script>

<template>
  <div class="monitor-page">
    <!-- 资源使用概览 -->
    <el-row :gutter="16" class="overview-section">
      <!-- CPU仪表盘 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="gauge-card">
          <div class="gauge-header">
            <el-icon class="gauge-icon cpu"><Cpu /></el-icon>
            <span>CPU 使用率</span>
          </div>
          <div ref="cpuGaugeRef" class="gauge-chart" />
          <div class="gauge-info">
            <div class="info-row">
              <span>核心数</span>
              <span class="value">{{ server.cpu?.cores || 0 }} 核</span>
            </div>
            <div class="info-row">
              <span>负载</span>
              <span class="value">{{
                server.load?.load1?.toFixed(2) || 0
              }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 内存仪表盘 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="gauge-card">
          <div class="gauge-header">
            <el-icon class="gauge-icon memory"><Memo /></el-icon>
            <span>内存使用率</span>
          </div>
          <div ref="memGaugeRef" class="gauge-chart" />
          <div class="gauge-info">
            <div class="info-row">
              <span>已用</span>
              <span class="value">{{ server.memory?.usedFormat || "-" }}</span>
            </div>
            <div class="info-row">
              <span>总计</span>
              <span class="value">{{ server.memory?.totalFormat || "-" }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- Swap内存 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="gauge-card">
          <div class="gauge-header">
            <el-icon class="gauge-icon swap"><TrendCharts /></el-icon>
            <span>Swap 交换内存</span>
          </div>
          <div class="swap-chart-area">
            <el-progress
              type="dashboard"
              :percentage="server.memory?.swapUsedPercent || 0"
              :color="getUsageColor(server.memory?.swapUsedPercent || 0)"
              :stroke-width="12"
              :width="140"
            >
              <template #default="{ percentage }">
                <div class="swap-center">
                  <span class="swap-percent">{{ percentage.toFixed(1) }}%</span>
                  <span class="swap-label">Swap</span>
                </div>
              </template>
            </el-progress>
          </div>
          <div class="gauge-info">
            <div class="info-row">
              <span>已用</span>
              <span class="value orange">{{
                server.memory?.swapUsedFormat || "0 B"
              }}</span>
            </div>
            <div class="info-row">
              <span>总计</span>
              <span class="value">{{
                server.memory?.swapTotalFormat || "0 B"
              }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 运行时间 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="gauge-card">
          <div class="gauge-header">
            <el-icon class="gauge-icon uptime"><Timer /></el-icon>
            <span>系统运行时间</span>
          </div>
          <div class="uptime-chart-area">
            <div class="uptime-value">
              {{ server.general?.uptimeFormat || "-" }}
            </div>
          </div>
          <div class="gauge-info">
            <div class="info-row">
              <span>主机名</span>
              <span class="value">{{ server.general?.hostname || "-" }}</span>
            </div>
            <div class="info-row">
              <span>系统</span>
              <span class="value">{{
                server.general?.os?.split(" ").slice(0, 2).join(" ") || "-"
              }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 实时曲线图 -->
    <el-row :gutter="16" class="chart-section">
      <el-col :span="24">
        <el-card shadow="hover" class="history-card">
          <template #header>
            <div class="chart-header">
              <el-icon class="chart-icon"><TrendCharts /></el-icon>
              <span>实时资源监控曲线</span>
              <el-tag size="small" type="info">5秒刷新</el-tag>
            </div>
          </template>
          <div ref="historyChartRef" class="history-chart" />
        </el-card>
      </el-col>
    </el-row>

    <!-- CPU核心详情 -->
    <el-row
      v-if="server.cpu?.coreDetails?.length"
      :gutter="16"
      class="cores-section"
    >
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><Cpu /></el-icon>
              <span>CPU 核心使用率</span>
              <el-tag size="small" type="info" class="core-count"
                >{{ server.cpu?.coreDetails?.length }} 核心</el-tag
              >
            </div>
          </template>
          <div class="cores-grid">
            <div
              v-for="(usage, index) in server.cpu?.coreDetails"
              :key="index"
              class="core-card"
            >
              <div class="core-header">
                <span class="core-name">核心 {{ index }}</span>
                <span
                  class="core-percent"
                  :style="{ color: getUsageColor(usage) }"
                  >{{ usage?.toFixed(0) }}%</span
                >
              </div>
              <div class="core-progress-bg">
                <div
                  class="core-progress-fill"
                  :style="{
                    width: `${usage}%`,
                    background: `linear-gradient(90deg, ${getUsageColor(usage)}88, ${getUsageColor(usage)})`
                  }"
                />
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 磁盘状态 -->
    <el-row :gutter="16" class="disk-section">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon disk"><FolderOpened /></el-icon>
              <span>磁盘存储状态</span>
            </div>
          </template>
          <div class="disk-grid">
            <div
              v-for="(disk, index) in server.disk"
              :key="index"
              class="disk-item"
            >
              <div class="disk-header">
                <el-tag type="info" size="small">{{ disk.device }}</el-tag>
                <span class="mount">{{ disk.mountpoint }}</span>
              </div>
              <el-progress
                :percentage="disk.usedPercent"
                :stroke-width="12"
                :color="getUsageColor(disk.usedPercent)"
                :format="() => `${disk.usedPercent?.toFixed(1)}%`"
              />
              <div class="disk-details">
                <span class="used">已用: {{ disk.usedFormat }}</span>
                <span class="total">总计: {{ disk.totalFormat }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 网络状态 -->
    <el-row :gutter="16" class="network-section">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon network"><Connection /></el-icon>
              <span>网络接口状态</span>
            </div>
          </template>
          <el-table
            :data="server.network?.filter(item => item.isUp) || []"
            stripe
            :header-cell-style="{ background: '#f5f7fa', color: '#606266' }"
          >
            <el-table-column
              prop="name"
              label="接口"
              width="150"
              align="center"
            >
              <template #default="{ row }">
                <el-tag :type="row.isUp ? 'success' : 'danger'" size="small">
                  {{ row.name }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="接收" align="center">
              <template #default="{ row }">
                <span class="text-green">↓ {{ row.recvFormat }}</span>
              </template>
            </el-table-column>
            <el-table-column label="发送" align="center">
              <template #default="{ row }">
                <span class="text-blue">↑ {{ row.sentFormat }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="packetsRecv" label="收包数" align="center" />
            <el-table-column prop="packetsSent" label="发包数" align="center" />
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag
                  :type="row.isUp ? 'success' : 'danger'"
                  size="small"
                  effect="light"
                >
                  {{ row.isUp ? "● 活跃" : "○ 断开" }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 服务器详细信息 -->
    <el-row :gutter="16" class="server-section">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon server"><Platform /></el-icon>
              <span>服务器详细信息</span>
            </div>
          </template>
          <el-descriptions :column="4" border>
            <el-descriptions-item label="主机名">{{
              server.general?.hostname || "-"
            }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{
              server.general?.os || "-"
            }}</el-descriptions-item>
            <el-descriptions-item label="系统架构">{{
              server.general?.arch || "-"
            }}</el-descriptions-item>
            <el-descriptions-item label="内核版本">{{
              server.general?.kernel || "-"
            }}</el-descriptions-item>
            <el-descriptions-item label="Go版本">{{
              server.general?.goVersion || "-"
            }}</el-descriptions-item>
            <el-descriptions-item label="CPU型号">{{
              server.cpu?.modelName || "-"
            }}</el-descriptions-item>
            <el-descriptions-item label="CPU频率"
              >{{
                server.cpu?.frequency?.toFixed(0) || 0
              }}
              MHz</el-descriptions-item
            >
            <el-descriptions-item label="系统负载">
              {{ server.load?.load1?.toFixed(2) || 0 }} /
              {{ server.load?.load5?.toFixed(2) || 0 }} /
              {{ server.load?.load15?.toFixed(2) || 0 }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.monitor-page {
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  min-height: calc(100vh - 100px);
}

/* 顶部栏 */
.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  margin-bottom: 20px;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  font-size: 28px;
}

.header-title {
  font-size: 20px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 6px;
}

.pulse {
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 仪表盘卡片 */
.gauge-card {
  border-radius: 12px;
  border: none;
  margin-bottom: 16px;
}

.gauge-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-weight: 600;
  color: #333;
}

.gauge-icon {
  font-size: 20px;
  padding: 6px;
  border-radius: 8px;
  color: white;
}

.gauge-icon.cpu {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
}

.gauge-icon.memory {
  background: linear-gradient(135deg, #10b981, #059669);
}

.gauge-icon.swap {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.gauge-icon.uptime {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
}

.gauge-chart {
  height: 180px;
}

.swap-chart-area {
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.swap-center {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.swap-percent {
  font-size: 22px;
  font-weight: bold;
  color: #f59e0b;
}

.swap-label {
  font-size: 12px;
  color: #999;
}

.uptime-chart-area {
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.uptime-value {
  font-size: 22px;
  font-weight: bold;
  color: #8b5cf6;
  text-align: center;
}

.info-row .value.orange {
  color: #f59e0b;
}

.gauge-info {
  border-top: 1px solid #eee;
  padding-top: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 13px;
}

.info-row .value {
  font-weight: 600;
  color: #333;
}

/* 统计卡片 */
.stat-card {
  border-radius: 12px;
  border: none;
  margin-bottom: 16px;
  height: 100%;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-weight: 600;
  color: #333;
}

.stat-icon {
  font-size: 20px;
  padding: 6px;
  border-radius: 8px;
  color: white;
}

.stat-icon.swap {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.stat-icon.uptime {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
}

/* Swap卡片 */
.swap-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.swap-percent {
  font-size: 16px;
  font-weight: bold;
}

.swap-info {
  flex: 1;
}

.swap-info .info-item {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px dashed #eee;
}

.swap-info .info-item:last-child {
  border-bottom: none;
}

.swap-info .label {
  color: #666;
  font-size: 13px;
}

.swap-info .value {
  font-weight: 600;
  color: #333;
}

.swap-info .value.orange {
  color: #f59e0b;
}

/* 运行时间卡片 */
.uptime-content {
  text-align: center;
}

.uptime-value {
  font-size: 24px;
  font-weight: bold;
  color: #8b5cf6;
  margin-bottom: 16px;
}

.uptime-details {
  text-align: left;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 12px;
}

.detail-item .label {
  color: #999;
}

.detail-item .value {
  color: #333;
  font-weight: 500;
}

/* 历史曲线 */
.history-card {
  border-radius: 12px;
  border: none;
  margin-bottom: 16px;
}

.chart-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
}

.chart-icon {
  font-size: 18px;
  color: #667eea;
}

.history-chart {
  height: 300px;
}

/* 通用卡片头 */
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.card-header .header-icon {
  font-size: 18px;
  color: #667eea;
}

.card-header .header-icon.disk {
  color: #ef4444;
}

.card-header .header-icon.network {
  color: #3b82f6;
}

.card-header .header-icon.server {
  color: #f59e0b;
}

/* CPU核心显示 */
.core-count {
  margin-left: auto;
}

.cores-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.core-card {
  padding: 12px;
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;
}

.core-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: #cbd5e1;
}

.core-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.core-name {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.core-percent {
  font-size: 14px;
  font-weight: 700;
}

.core-progress-bg {
  height: 6px;
  background: #e2e8f0;
  border-radius: 3px;
  overflow: hidden;
}

.core-progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.5s ease;
}

/* 磁盘网格 */
.disk-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.disk-item {
  padding: 16px;
  background: #f9fafb;
  border-radius: 10px;
  transition: all 0.3s ease;
}

.disk-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.disk-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.disk-header .mount {
  font-size: 13px;
  color: #666;
}

.disk-details {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 12px;
  color: #666;
}

.disk-details .used {
  color: #f59e0b;
}

/* 网络表格 */
.text-green {
  color: #10b981;
  font-weight: 500;
}

.text-blue {
  color: #3b82f6;
  font-weight: 500;
}

/* 各区块间距 */
.overview-section,
.chart-section,
.cores-section,
.disk-section,
.network-section,
.server-section {
  margin-bottom: 16px;
}

/* 响应式 */
@media (max-width: 768px) {
  .monitor-page {
    padding: 12px;
  }

  .header-bar {
    flex-direction: column;
    gap: 12px;
  }

  .swap-content {
    flex-direction: column;
  }

  .cores-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .disk-grid {
    grid-template-columns: 1fr;
  }
}

/* 卡片动画 */
.el-card {
  transition: all 0.3s ease;
}

.el-card:hover {
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}
</style>
