<script setup lang="ts">
defineOptions({
  name: "Monitor"
});

import {
  ref,
  onMounted,
  computed,
  onBeforeUnmount,
  watch,
  nextTick
} from "vue";
import dayjs from "dayjs";
import { ElMessage } from "element-plus";
import { getServe } from "@/api/monitor";
import echarts from "@/utils/echarts";
import {
  Monitor,
  Refresh,
  Cpu,
  Memo,
  TrendCharts,
  Platform,
  FolderOpened,
  Connection,
  Timer,
  Clock,
  ArrowDown,
  ArrowUp,
  Search,
  Warning
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
    inodesTotal?: number;
    inodesUsed?: number;
    inodesFree?: number;
    inodesUsedPercent?: number;
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

// 核心数据状态
const server = ref<ServerInfo>({});
const loading = ref(true); // 首次加载占位
const refreshing = ref(false); // 手动刷新按钮加载动画
const autoRefresh = ref(true);
const refreshIntervalSec = ref(5);
let timer: ReturnType<typeof setInterval> | null = null;
const lastUpdatedTime = ref("");

// 图表维度切换
const activeChartMetric = ref<"resource" | "network">("resource");
const maxDataPoints = ref<number>(30);

// 历史数据存储
const historyData = ref<{
  timestamps: string[];
  cpuData: number[];
  memoryData: number[];
  netRecvData: number[]; // KB/s
  netSentData: number[]; // KB/s
}>({
  timestamps: [],
  cpuData: [],
  memoryData: [],
  netRecvData: [],
  netSentData: []
});

// 网络吞吐速率计算快照
const prevNetworkSnapshot = ref<{
  timestamp: number;
  rates: Record<string, { bytesRecv: number; bytesSent: number }>;
} | null>(null);

const networkRates = ref<
  Record<
    string,
    {
      recvRate: number;
      sentRate: number;
      recvRateFormat: string;
      sentRateFormat: string;
    }
  >
>({});

const totalNetworkRate = ref({
  recvRate: 0,
  sentRate: 0,
  recvRateFormat: "0 B/s",
  sentRateFormat: "0 B/s"
});

// ================= CPU 核心自适应与热力矩阵相关逻辑 =================
const coreViewMode = ref<"auto" | "heatmap" | "card">("auto");
const showAllCores = ref(false);

const coreCount = computed(() => server.value.cpu?.coreDetails?.length || 0);

const activeCoreView = computed(() => {
  if (coreViewMode.value !== "auto") return coreViewMode.value;
  return coreCount.value > 8 ? "heatmap" : "card";
});

const displayedCores = computed(() => {
  const cores = server.value.cpu?.coreDetails || [];
  if (showAllCores.value || cores.length <= 8) {
    return cores;
  }
  return cores.slice(0, 8);
});

const coreDiagnostics = computed(() => {
  const cores = server.value.cpu?.coreDetails || [];
  if (cores.length === 0) {
    return {
      maxIndex: 0,
      maxVal: 0,
      minIndex: 0,
      minVal: 0,
      avgVal: 0,
      skew: 0,
      isSkewed: false,
      statusText: "无核心数据",
      statusType: "info" as const
    };
  }

  let maxVal = cores[0];
  let maxIndex = 0;
  let minVal = cores[0];
  let minIndex = 0;
  let sum = 0;

  cores.forEach((val, idx) => {
    sum += val;
    if (val > maxVal) {
      maxVal = val;
      maxIndex = idx;
    }
    if (val < minVal) {
      minVal = val;
      minIndex = idx;
    }
  });

  const avgVal = sum / cores.length;
  const skew = maxVal - minVal;
  const isSkewed = skew >= 50 && maxVal >= 70;

  let statusText = "负载均衡";
  let statusType: "success" | "warning" | "danger" = "success";
  if (isSkewed) {
    statusText = "单核高负荷倾斜";
    statusType = "danger";
  } else if (skew >= 35) {
    statusText = "轻度波动";
    statusType = "warning";
  }

  return {
    maxIndex,
    maxVal,
    minIndex,
    minVal,
    avgVal,
    skew,
    isSkewed,
    statusText,
    statusType
  };
});

const getHeatmapColor = (usage: number) => {
  const val = Math.min(100, Math.max(0, usage || 0));
  if (val >= 90) return "#ef4444";
  if (val >= 75) return "#f97316";
  if (val >= 50) return "#eab308";
  if (val >= 25) return "#10b981";
  if (val >= 10) return "#34d399";
  return "var(--el-fill-color, #e2e8f0)";
};

const getHeatmapTextColor = (usage: number) => {
  const val = Math.min(100, Math.max(0, usage || 0));
  if (val >= 25) return "#ffffff";
  return "var(--el-text-color-regular, #475569)";
};

const getCoreStatusText = (usage: number) => {
  const val = Math.min(100, Math.max(0, usage || 0));
  if (val >= 85) return "满载高负荷";
  if (val >= 60) return "计算繁忙";
  if (val >= 30) return "正常计算";
  return "平稳空闲";
};
// ===================================================================

// ================= 磁盘存储全局池与过滤排序逻辑 =================
// 磁盘展示模式：auto（智能推荐） | card（卡片） | table（紧凑表格）
const diskViewMode = ref<"auto" | "card" | "table">("auto");
const diskSearchQuery = ref("");
const diskOnlyAlert = ref(false); // 仅查看告警盘 (> 80% 或 Inode > 85%)
const diskSortBy = ref<"usage_desc" | "total_desc" | "mountpoint">("usage_desc");

// 实际生效的视图模式（盘数 > 6 时默认采用表格，<= 6 时采用卡片）
const activeDiskView = computed(() => {
  if (diskViewMode.value !== "auto") return diskViewMode.value;
  return (server.value.disk?.length || 0) > 6 ? "table" : "card";
});

// 格式化字节数
const formatBytes = (bytes: number): string => {
  if (bytes <= 0 || isNaN(bytes)) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  let i = 0;
  let val = bytes;
  while (val >= 1024 && i < units.length - 1) {
    val /= 1024;
    i++;
  }
  return `${val < 10 ? val.toFixed(2) : val.toFixed(1)} ${units[i]}`;
};

// 全局存储池汇总指标 (Total Storage Pool)
const storagePool = computed(() => {
  const disks = server.value.disk || [];
  if (disks.length === 0) {
    return {
      totalBytes: 0,
      usedBytes: 0,
      freeBytes: 0,
      usedPercent: 0,
      totalFormat: "0 B",
      usedFormat: "0 B",
      freeFormat: "0 B",
      diskCount: 0,
      maxUsedDisk: null as any,
      alertCount: 0,
      hasAlert: false
    };
  }

  let totalBytes = 0;
  let usedBytes = 0;
  let freeBytes = 0;
  let maxUsedDisk = disks[0];
  let alertCount = 0;

  for (const d of disks) {
    totalBytes += d.total || 0;
    usedBytes += d.used || 0;
    freeBytes += d.free || 0;
    if ((d.usedPercent || 0) > (maxUsedDisk.usedPercent || 0)) {
      maxUsedDisk = d;
    }
    if ((d.usedPercent || 0) >= 85 || (d.inodesUsedPercent || 0) >= 85) {
      alertCount++;
    }
  }

  const usedPercent = totalBytes > 0 ? (usedBytes / totalBytes) * 100 : 0;
  const hasAlert = (maxUsedDisk?.usedPercent || 0) >= 85;

  return {
    totalBytes,
    usedBytes,
    freeBytes,
    usedPercent,
    totalFormat: formatBytes(totalBytes),
    usedFormat: formatBytes(usedBytes),
    freeFormat: formatBytes(freeBytes),
    diskCount: disks.length,
    maxUsedDisk,
    alertCount,
    hasAlert
  };
});

// 经过搜索、告警筛选与排序后的磁盘挂载列表
const filteredDisks = computed(() => {
  let list = [...(server.value.disk || [])];

  // 1. 告警筛选
  if (diskOnlyAlert.value) {
    list = list.filter(
      d =>
        (d.usedPercent || 0) >= 80 ||
        (d.inodesUsedPercent && d.inodesUsedPercent >= 85)
    );
  }

  // 2. 关键字搜索 (挂载点、设备名、文件系统类型)
  const q = diskSearchQuery.value.trim().toLowerCase();
  if (q) {
    list = list.filter(
      d =>
        d.mountpoint.toLowerCase().includes(q) ||
        d.device.toLowerCase().includes(q) ||
        (d.fstype && d.fstype.toLowerCase().includes(q))
    );
  }

  // 3. 排序
  if (diskSortBy.value === "usage_desc") {
    list.sort((a, b) => (b.usedPercent || 0) - (a.usedPercent || 0));
  } else if (diskSortBy.value === "total_desc") {
    list.sort((a, b) => (b.total || 0) - (a.total || 0));
  } else if (diskSortBy.value === "mountpoint") {
    list.sort((a, b) => a.mountpoint.localeCompare(b.mountpoint));
  }

  return list;
});
// ===============================================================

// 系统综合健康评估
const systemHealth = computed(() => {
  const cpuVal = server.value.cpu?.usagePercent || 0;
  const memVal = server.value.memory?.usedPercent || 0;
  const maxDisk = storagePool.value.maxUsedDisk?.usedPercent || 0;

  if (cpuVal > 85 || memVal > 90 || maxDisk > 92) {
    return {
      status: "critical",
      text: "高负载告警",
      tagType: "danger" as const,
      color: "#f56c6c"
    };
  }
  if (cpuVal > 70 || memVal > 75 || maxDisk > 80) {
    return {
      status: "warning",
      text: "负载偏高",
      tagType: "warning" as const,
      color: "#e6a23c"
    };
  }
  return {
    status: "healthy",
    text: "运行平稳",
    tagType: "success" as const,
    color: "#67c23a"
  };
});

// 格式化网络速率
const formatRate = (bytesPerSec: number): string => {
  if (bytesPerSec <= 0 || isNaN(bytesPerSec)) return "0 B/s";
  const units = ["B/s", "KB/s", "MB/s", "GB/s", "TB/s"];
  let i = 0;
  let val = bytesPerSec;
  while (val >= 1024 && i < units.length - 1) {
    val /= 1024;
    i++;
  }
  return `${val < 10 ? val.toFixed(2) : val.toFixed(1)} ${units[i]}`;
};

// 状态色阶
const getUsageColor = (percent: number) => {
  if (percent >= 85) return "#f56c6c";
  if (percent >= 70) return "#e6a23c";
  return "#67c23a";
};

const getCoreGradient = (percent: number) => {
  const color = getUsageColor(percent);
  return `linear-gradient(90deg, ${color}66 0%, ${color} 100%)`;
};

const getLoadClass = (val?: number) => {
  if (!val) return "";
  const cores = server.value.cpu?.cores || 1;
  if (val > cores * 1.2) return "text-red font-bold";
  if (val > cores * 0.8) return "text-orange font-medium";
  return "text-green";
};

// 图表 DOM 引用与实例
const cpuGaugeRef = ref<HTMLElement | null>(null);
const memGaugeRef = ref<HTMLElement | null>(null);
const historyChartRef = ref<HTMLElement | null>(null);
let cpuGaugeChart: echarts.ECharts | null = null;
let memGaugeChart: echarts.ECharts | null = null;
let historyChart: echarts.ECharts | null = null;
let themeObserver: MutationObserver | null = null;

const isDark = () => document.documentElement.classList.contains("dark");

// 初始化与更新图表
const initCharts = () => {
  if (cpuGaugeRef.value && !cpuGaugeChart) {
    cpuGaugeChart = echarts.init(cpuGaugeRef.value);
  }
  if (memGaugeRef.value && !memGaugeChart) {
    memGaugeChart = echarts.init(memGaugeRef.value);
  }
  if (historyChartRef.value && !historyChart) {
    historyChart = echarts.init(historyChartRef.value);
  }
  renderAllCharts();
  window.addEventListener("resize", handleResize);

  themeObserver = new MutationObserver(() => {
    renderAllCharts();
  });
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"]
  });
};

const handleResize = () => {
  cpuGaugeChart?.resize();
  memGaugeChart?.resize();
  historyChart?.resize();
};

const updateCpuGauge = () => {
  if (!cpuGaugeChart) return;
  const dark = isDark();
  const value = server.value.cpu?.usagePercent || 0;
  const color = getUsageColor(value);
  const option = {
    series: [
      {
        type: "gauge",
        startAngle: 210,
        endAngle: -30,
        min: 0,
        max: 100,
        splitNumber: 10,
        itemStyle: { color },
        progress: {
          show: true,
          roundCap: true,
          width: 12
        },
        pointer: { show: false },
        axisLine: {
          roundCap: true,
          lineStyle: {
            width: 12,
            color: [[1, dark ? "rgba(255,255,255,0.08)" : "#f1f5f9"]]
          }
        },
        axisTick: { show: false },
        splitLine: { show: false },
        axisLabel: { show: false },
        title: {
          show: true,
          offsetCenter: [0, "32%"],
          fontSize: 13,
          color: dark ? "#94a3b8" : "#64748b"
        },
        detail: {
          valueAnimation: true,
          offsetCenter: [0, "-8%"],
          fontSize: 26,
          fontWeight: "bold",
          formatter: "{value}%",
          color
        },
        data: [{ value: Number(value.toFixed(1)), name: "CPU 使用率" }]
      }
    ]
  };
  cpuGaugeChart.setOption(option);
};

const updateMemGauge = () => {
  if (!memGaugeChart) return;
  const dark = isDark();
  const value = server.value.memory?.usedPercent || 0;
  const color = getUsageColor(value);
  const option = {
    series: [
      {
        type: "gauge",
        startAngle: 210,
        endAngle: -30,
        min: 0,
        max: 100,
        splitNumber: 10,
        itemStyle: { color },
        progress: {
          show: true,
          roundCap: true,
          width: 12
        },
        pointer: { show: false },
        axisLine: {
          roundCap: true,
          lineStyle: {
            width: 12,
            color: [[1, dark ? "rgba(255,255,255,0.08)" : "#f1f5f9"]]
          }
        },
        axisTick: { show: false },
        splitLine: { show: false },
        axisLabel: { show: false },
        title: {
          show: true,
          offsetCenter: [0, "32%"],
          fontSize: 13,
          color: dark ? "#94a3b8" : "#64748b"
        },
        detail: {
          valueAnimation: true,
          offsetCenter: [0, "-8%"],
          fontSize: 26,
          fontWeight: "bold",
          formatter: "{value}%",
          color
        },
        data: [{ value: Number(value.toFixed(1)), name: "内存使用率" }]
      }
    ]
  };
  memGaugeChart.setOption(option);
};

const updateHistoryChart = () => {
  if (!historyChart) return;
  const dark = isDark();
  const textColor = dark ? "#94a3b8" : "#64748b";
  const splitLineColor = dark
    ? "rgba(255, 255, 255, 0.06)"
    : "rgba(0, 0, 0, 0.05)";

  if (activeChartMetric.value === "resource") {
    const option = {
      backgroundColor: "transparent",
      tooltip: {
        trigger: "axis",
        axisPointer: { type: "cross" },
        valueFormatter: (value: any) => `${Number(value || 0).toFixed(1)}%`
      },
      legend: {
        data: ["CPU 使用率", "内存使用率"],
        bottom: 0,
        textStyle: { color: textColor }
      },
      grid: {
        left: "2%",
        right: "3%",
        bottom: "12%",
        top: "10%",
        containLabel: true
      },
      xAxis: {
        type: "category",
        boundaryGap: false,
        data: historyData.value.timestamps,
        axisLabel: { color: textColor, fontSize: 11 },
        axisLine: { lineStyle: { color: dark ? "#334155" : "#e2e8f0" } }
      },
      yAxis: {
        type: "value",
        min: 0,
        max: 100,
        axisLabel: {
          formatter: "{value}%",
          color: textColor,
          fontSize: 11
        },
        splitLine: { lineStyle: { color: splitLineColor } }
      },
      series: [
        {
          name: "CPU 使用率",
          type: "line",
          smooth: true,
          symbol: "circle",
          symbolSize: 4,
          areaStyle: {
            opacity: 0.2,
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: "rgba(59, 130, 246, 0.6)" },
              { offset: 1, color: "rgba(59, 130, 246, 0.02)" }
            ])
          },
          lineStyle: { width: 2.2, color: "#3b82f6" },
          itemStyle: { color: "#3b82f6" },
          data: historyData.value.cpuData
        },
        {
          name: "内存使用率",
          type: "line",
          smooth: true,
          symbol: "circle",
          symbolSize: 4,
          areaStyle: {
            opacity: 0.2,
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: "rgba(16, 185, 129, 0.6)" },
              { offset: 1, color: "rgba(16, 185, 129, 0.02)" }
            ])
          },
          lineStyle: { width: 2.2, color: "#10b981" },
          itemStyle: { color: "#10b981" },
          data: historyData.value.memoryData
        }
      ]
    };
    historyChart.setOption(option, true);
  } else {
    // 网络流量吞吐曲线 (单位: KB/s)
    const option = {
      backgroundColor: "transparent",
      tooltip: {
        trigger: "axis",
        axisPointer: { type: "cross" },
        valueFormatter: (value: any) =>
          formatRate(Number(value || 0) * 1024)
      },
      legend: {
        data: ["实时下行速率 (入网)", "实时上行速率 (出网)"],
        bottom: 0,
        textStyle: { color: textColor }
      },
      grid: {
        left: "2%",
        right: "3%",
        bottom: "12%",
        top: "10%",
        containLabel: true
      },
      xAxis: {
        type: "category",
        boundaryGap: false,
        data: historyData.value.timestamps,
        axisLabel: { color: textColor, fontSize: 11 },
        axisLine: { lineStyle: { color: dark ? "#334155" : "#e2e8f0" } }
      },
      yAxis: {
        type: "value",
        axisLabel: {
          formatter: (val: number) => `${val.toFixed(0)} KB/s`,
          color: textColor,
          fontSize: 11
        },
        splitLine: { lineStyle: { color: splitLineColor } }
      },
      series: [
        {
          name: "实时下行速率 (入网)",
          type: "line",
          smooth: true,
          symbol: "circle",
          symbolSize: 4,
          areaStyle: {
            opacity: 0.2,
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: "rgba(16, 185, 129, 0.6)" },
              { offset: 1, color: "rgba(16, 185, 129, 0.02)" }
            ])
          },
          lineStyle: { width: 2.2, color: "#10b981" },
          itemStyle: { color: "#10b981" },
          data: historyData.value.netRecvData
        },
        {
          name: "实时上行速率 (出网)",
          type: "line",
          smooth: true,
          symbol: "circle",
          symbolSize: 4,
          areaStyle: {
            opacity: 0.2,
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: "rgba(139, 92, 246, 0.6)" },
              { offset: 1, color: "rgba(139, 92, 246, 0.02)" }
            ])
          },
          lineStyle: { width: 2.2, color: "#8b5cf6" },
          itemStyle: { color: "#8b5cf6" },
          data: historyData.value.netSentData
        }
      ]
    };
    historyChart.setOption(option, true);
  }
};

const renderAllCharts = () => {
  updateCpuGauge();
  updateMemGauge();
  updateHistoryChart();
};

// 计算网络接口实时传输速率
const updateNetworkRates = () => {
  const now = Date.now();
  if (prevNetworkSnapshot.value && server.value.network) {
    const timeDeltaSec = (now - prevNetworkSnapshot.value.timestamp) / 1000;
    if (timeDeltaSec > 0.4) {
      let sumRecvRate = 0;
      let sumSentRate = 0;
      const rates: Record<string, any> = {};

      for (const iface of server.value.network) {
        const prev = prevNetworkSnapshot.value.rates[iface.name];
        if (
          prev &&
          iface.bytesRecv >= prev.bytesRecv &&
          iface.bytesSent >= prev.bytesSent
        ) {
          const recvRate = (iface.bytesRecv - prev.bytesRecv) / timeDeltaSec;
          const sentRate = (iface.bytesSent - prev.bytesSent) / timeDeltaSec;
          rates[iface.name] = {
            recvRate,
            sentRate,
            recvRateFormat: formatRate(recvRate),
            sentRateFormat: formatRate(sentRate)
          };
          if (iface.name !== "lo" && iface.isUp) {
            sumRecvRate += recvRate;
            sumSentRate += sentRate;
          }
        } else {
          rates[iface.name] = {
            recvRate: 0,
            sentRate: 0,
            recvRateFormat: "0 B/s",
            sentRateFormat: "0 B/s"
          };
        }
      }
      networkRates.value = rates;
      totalNetworkRate.value = {
        recvRate: sumRecvRate,
        sentRate: sumSentRate,
        recvRateFormat: formatRate(sumRecvRate),
        sentRateFormat: formatRate(sumSentRate)
      };
    }
  }

  // 记录本次快照
  const snapshotMap: Record<string, { bytesRecv: number; bytesSent: number }> =
    {};
  if (server.value.network) {
    for (const iface of server.value.network) {
      snapshotMap[iface.name] = {
        bytesRecv: iface.bytesRecv,
        bytesSent: iface.bytesSent
      };
    }
  }
  prevNetworkSnapshot.value = {
    timestamp: now,
    rates: snapshotMap
  };
};

// 记录历史序列
const updateHistoryData = () => {
  const timeStr = dayjs().format("HH:mm:ss");

  historyData.value.timestamps.push(timeStr);
  historyData.value.cpuData.push(
    Number((server.value.cpu?.usagePercent || 0).toFixed(1))
  );
  historyData.value.memoryData.push(
    Number((server.value.memory?.usedPercent || 0).toFixed(1))
  );
  historyData.value.netRecvData.push(
    Number((totalNetworkRate.value.recvRate / 1024).toFixed(1))
  );
  historyData.value.netSentData.push(
    Number((totalNetworkRate.value.sentRate / 1024).toFixed(1))
  );

  const limit = maxDataPoints.value;
  while (historyData.value.timestamps.length > limit) {
    historyData.value.timestamps.shift();
    historyData.value.cpuData.shift();
    historyData.value.memoryData.shift();
    historyData.value.netRecvData.shift();
    historyData.value.netSentData.shift();
  }
};

// 获取服务器监控完整数据
const fetchServerInfo = async (isManual = false) => {
  if (isManual) {
    refreshing.value = true;
  }
  try {
    const rep = await getServe();
    if (rep && rep.code === 200 && rep.payload) {
      server.value = rep.payload;
      lastUpdatedTime.value = dayjs().format("HH:mm:ss");
      updateNetworkRates();
      updateHistoryData();
      renderAllCharts();
      if (isManual) {
        ElMessage.success("监控数据刷新成功");
      }
    } else {
      if (isManual) {
        ElMessage.error(rep?.message || "获取系统监控数据失败");
      }
    }
  } catch (error: any) {
    console.error("fetch server info error:", error);
    if (isManual) {
      ElMessage.error(error?.message || "网络请求异常");
    }
  } finally {
    loading.value = false;
    refreshing.value = false;
  }
};

// 控制器交互逻辑
const manualRefresh = () => {
  fetchServerInfo(true);
};

const handleAutoRefreshChange = (val: boolean) => {
  if (val) {
    startAutoRefresh();
    ElMessage.success(`已启用自动刷新 (${refreshIntervalSec.value}秒)`);
  } else {
    stopAutoRefresh();
    ElMessage.info("已暂停自动刷新");
  }
};

const handleIntervalChange = () => {
  if (autoRefresh.value) {
    startAutoRefresh();
  }
};

const handleMetricTabChange = () => {
  nextTick(() => {
    updateHistoryChart();
  });
};

const handleDataPointsChange = () => {
  const limit = maxDataPoints.value;
  while (historyData.value.timestamps.length > limit) {
    historyData.value.timestamps.shift();
    historyData.value.cpuData.shift();
    historyData.value.memoryData.shift();
    historyData.value.netRecvData.shift();
    historyData.value.netSentData.shift();
  }
  updateHistoryChart();
};

const startAutoRefresh = () => {
  stopAutoRefresh();
  timer = setInterval(() => {
    fetchServerInfo(false);
  }, refreshIntervalSec.value * 1000);
};

const stopAutoRefresh = () => {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
};

onMounted(async () => {
  await fetchServerInfo();
  setTimeout(initCharts, 80);
  if (autoRefresh.value) {
    startAutoRefresh();
  }
});

onBeforeUnmount(() => {
  stopAutoRefresh();
  window.removeEventListener("resize", handleResize);
  themeObserver?.disconnect();
  cpuGaugeChart?.dispose();
  memGaugeChart?.dispose();
  historyChart?.dispose();
  cpuGaugeChart = null;
  memGaugeChart = null;
  historyChart = null;
});
</script>

<template>
  <div class="system-monitor-container">
    <!-- 顶部操作与状态栏 -->
    <div class="monitor-control-bar">
      <div class="bar-left">
        <div class="title-with-icon">
          <div class="title-icon-box">
            <el-icon><Monitor /></el-icon>
          </div>
          <div class="title-text-group">
            <div class="bar-title">系统服务器监控</div>
            <div class="bar-subtitle">
              实时监控 CPU 计算资源、内存占用、磁盘存储与网卡流量吞吐
            </div>
          </div>
        </div>
        <div class="status-indicator">
          <el-tag
            :type="systemHealth.tagType"
            effect="dark"
            class="status-pulse-tag"
          >
            <span
              class="pulse-dot"
              :style="{ backgroundColor: systemHealth.color }"
            />
            {{ systemHealth.text }}
          </el-tag>
          <span v-if="lastUpdatedTime" class="update-time-text">
            <el-icon><Clock /></el-icon>
            最后更新：{{ lastUpdatedTime }}
          </span>
        </div>
      </div>

      <div class="bar-right">
        <div class="control-item">
          <span class="control-label">自动刷新</span>
          <el-switch
            v-model="autoRefresh"
            inline-prompt
            active-text="开"
            inactive-text="关"
            @change="handleAutoRefreshChange"
          />
        </div>

        <div v-if="autoRefresh" class="control-item">
          <span class="control-label">频率</span>
          <el-select
            v-model="refreshIntervalSec"
            size="small"
            class="interval-select"
            @change="handleIntervalChange"
          >
            <el-option :value="3" label="3 秒" />
            <el-option :value="5" label="5 秒" />
            <el-option :value="10" label="10 秒" />
            <el-option :value="30" label="30 秒" />
          </el-select>
        </div>

        <el-button
          type="primary"
          :loading="refreshing"
          :icon="Refresh"
          @click="manualRefresh"
        >
          刷新数据
        </el-button>
      </div>
    </div>

    <!-- 顶部核心资源四大卡片 -->
    <el-row :gutter="16" class="top-metrics-row">
      <!-- 1. CPU 卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="metric-card">
          <div class="card-top-header">
            <div class="header-badge cpu">
              <el-icon><Cpu /></el-icon>
            </div>
            <div class="header-titles">
              <span class="card-main-title">CPU 处理器</span>
              <span class="card-sub-title">计算算力与总体负载</span>
            </div>
            <el-tag
              size="small"
              :type="getUsageColor(server.cpu?.usagePercent || 0) === '#f56c6c' ? 'danger' : 'info'"
              class="header-tag"
            >
              {{ server.cpu?.cores || 0 }} 核心
            </el-tag>
          </div>

          <div ref="cpuGaugeRef" class="gauge-chart-container" />

          <div class="card-info-footer">
            <div class="info-item">
              <span class="label">主频 / 架构</span>
              <span class="val">
                {{ server.cpu?.frequency ? `${server.cpu.frequency.toFixed(0)} MHz` : '-' }}
                ({{ server.general?.arch || '-' }})
              </span>
            </div>
            <div class="info-item">
              <span class="label">平均负载 (1/5/15m)</span>
              <span class="val load-tag-wrap">
                <span :class="getLoadClass(server.load?.load1)">
                  {{ server.load?.load1?.toFixed(2) || '0.00' }}
                </span>
                <span class="sep">/</span>
                <span>{{ server.load?.load5?.toFixed(2) || '0.00' }}</span>
                <span class="sep">/</span>
                <span>{{ server.load?.load15?.toFixed(2) || '0.00' }}</span>
              </span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 2. 内存 卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="metric-card">
          <div class="card-top-header">
            <div class="header-badge memory">
              <el-icon><Memo /></el-icon>
            </div>
            <div class="header-titles">
              <span class="card-main-title">运行内存 (RAM)</span>
              <span class="card-sub-title">物理内存与 Swap</span>
            </div>
            <el-tag
              size="small"
              :type="getUsageColor(server.memory?.usedPercent || 0) === '#f56c6c' ? 'danger' : 'success'"
              class="header-tag"
            >
              {{ server.memory?.totalFormat || '0 B' }}
            </el-tag>
          </div>

          <div ref="memGaugeRef" class="gauge-chart-container" />

          <div class="card-info-footer">
            <div class="info-item">
              <span class="label">已用 / 可用</span>
              <span class="val">
                {{ server.memory?.usedFormat || '-' }} /
                <span class="text-green">{{ server.memory?.freeFormat || '-' }}</span>
              </span>
            </div>
            <div class="info-item">
              <span class="label">Swap 交换区</span>
              <span
                v-if="server.memory?.swapTotal && server.memory.swapTotal > 0"
                class="val"
              >
                {{ server.memory?.swapUsedFormat }} / {{ server.memory?.swapTotalFormat }}
                <span class="text-orange">({{ server.memory?.swapUsedPercent?.toFixed(1) }}%)</span>
              </span>
              <span v-else class="val text-muted">未启用 (推荐容器配置)</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 3. 磁盘全局存储池卡片 (全新升级：全局池聚合 + 最危险分区预警) -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="metric-card">
          <div class="card-top-header">
            <div class="header-badge disk">
              <el-icon><FolderOpened /></el-icon>
            </div>
            <div class="header-titles">
              <span class="card-main-title">存储空间池</span>
              <span class="card-sub-title">
                {{ storagePool.maxUsedDisk ? `最高负荷: ${storagePool.maxUsedDisk.mountpoint}` : '物理与数据卷' }}
              </span>
            </div>
            <el-tag
              size="small"
              :type="storagePool.hasAlert ? 'danger' : 'info'"
              class="header-tag"
            >
              {{ storagePool.hasAlert ? `⚠️ 最满 ${storagePool.maxUsedDisk?.usedPercent?.toFixed(0)}%` : `共 ${storagePool.diskCount} 个分区` }}
            </el-tag>
          </div>

          <div class="disk-progress-container">
            <el-progress
              type="dashboard"
              :percentage="Number(storagePool.usedPercent.toFixed(1))"
              :color="getUsageColor(storagePool.usedPercent)"
              :stroke-width="12"
              :width="150"
            >
              <template #default="{ percentage }">
                <div class="gauge-center-content">
                  <span
                    class="gauge-center-val"
                    :style="{ color: getUsageColor(percentage) }"
                  >
                    {{ percentage.toFixed(1) }}%
                  </span>
                  <span class="gauge-center-label">总空间已用</span>
                </div>
              </template>
            </el-progress>
          </div>

          <div class="card-info-footer">
            <div class="info-item">
              <span class="label">总空间 (已用/总量)</span>
              <span class="val">
                {{ storagePool.usedFormat }} / {{ storagePool.totalFormat }}
              </span>
            </div>
            <div class="info-item">
              <span class="label">总剩余可用空间</span>
              <span class="val text-green">{{ storagePool.freeFormat }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 4. 系统与运行时 卡片 -->
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="hover" class="metric-card">
          <div class="card-top-header">
            <div class="header-badge system">
              <el-icon><Platform /></el-icon>
            </div>
            <div class="header-titles">
              <span class="card-main-title">系统与运行时</span>
              <span class="card-sub-title">
                {{ server.general?.os?.split(" ").slice(0, 2).join(" ") || 'Linux' }}
              </span>
            </div>
            <el-tag size="small" type="success" effect="light" class="header-tag">
              {{ server.general?.arch || 'amd64' }}
            </el-tag>
          </div>

          <div class="uptime-display-wrap">
            <div class="uptime-icon-badge">
              <el-icon><Timer /></el-icon>
            </div>
            <div class="uptime-text-block">
              <span class="uptime-title">系统已连续运行</span>
              <span class="uptime-val">
                {{ server.general?.uptimeFormat || '正在获取...' }}
              </span>
            </div>
          </div>

          <div class="card-info-footer">
            <div class="info-item">
              <span class="label">主机名</span>
              <span class="val truncate" :title="server.general?.hostname">
                {{ server.general?.hostname || '-' }}
              </span>
            </div>
            <div class="info-item">
              <span class="label">Go 运行时 / 协程</span>
              <span class="val">
                <span class="text-blue">{{ server.general?.goVersion || '-' }}</span>
                <span class="text-muted ml-1">({{ server.process?.running || 0 }} 协程)</span>
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 实时资源趋势图表 -->
    <el-row :gutter="16" class="chart-section-row">
      <el-col :span="24">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="section-card-header chart-custom-header">
              <div class="header-left">
                <el-icon class="section-icon chart"><TrendCharts /></el-icon>
                <span class="section-title">实时监控动态趋势</span>
                <el-radio-group
                  v-model="activeChartMetric"
                  size="small"
                  class="ml-3"
                  @change="handleMetricTabChange"
                >
                  <el-radio-button value="resource">
                    CPU & 内存趋势
                  </el-radio-button>
                  <el-radio-button value="network">
                    网络吞吐速率 (I/O)
                  </el-radio-button>
                </el-radio-group>
              </div>

              <div class="header-right">
                <span class="points-label">采样点数：</span>
                <el-radio-group
                  v-model="maxDataPoints"
                  size="small"
                  @change="handleDataPointsChange"
                >
                  <el-radio-button :value="30">近 30 次</el-radio-button>
                  <el-radio-button :value="60">近 60 次</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>

          <div ref="historyChartRef" class="history-chart-dom" />
        </el-card>
      </el-col>
    </el-row>

    <!-- ================= CPU 核心负载分布矩阵 (自适应双模：热力矩阵 vs 详细卡片) ================= -->
    <el-row
      v-if="server.cpu?.coreDetails?.length"
      :gutter="16"
      class="cores-section-row"
    >
      <el-col :span="24">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="section-card-header">
              <div class="header-left">
                <el-icon class="section-icon cpu"><Cpu /></el-icon>
                <span class="section-title">CPU 核心负载分布矩阵</span>
                <el-tag size="small" type="info" class="ml-2">
                  共 {{ coreCount }} 核心 / 线程
                </el-tag>
                <el-tag
                  size="small"
                  :type="coreDiagnostics.statusType"
                  effect="light"
                  class="ml-2"
                >
                  {{ coreDiagnostics.statusText }} (极差 {{ coreDiagnostics.skew.toFixed(0) }}%)
                </el-tag>
              </div>

              <div class="header-right">
                <el-radio-group v-model="coreViewMode" size="small">
                  <el-radio-button value="auto">
                    自适应 ({{ activeCoreView === 'heatmap' ? '热力' : '卡片' }})
                  </el-radio-button>
                  <el-radio-button value="heatmap">热力方阵</el-radio-button>
                  <el-radio-button value="card">详细卡片</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>

          <!-- 顶部核心极值与偏斜诊断指示栏 -->
          <div class="core-summary-strip">
            <div class="summary-pill max-pill">
              <span class="pill-icon">🔥</span>
              <span class="pill-label">最高负载核心：</span>
              <span class="pill-value text-red font-bold">
                Core #{{ coreDiagnostics.maxIndex }} ({{ coreDiagnostics.maxVal.toFixed(0) }}%)
              </span>
            </div>

            <div class="summary-pill min-pill">
              <span class="pill-icon">❄️</span>
              <span class="pill-label">最闲核心：</span>
              <span class="pill-value text-green font-semibold">
                Core #{{ coreDiagnostics.minIndex }} ({{ coreDiagnostics.minVal.toFixed(0) }}%)
              </span>
            </div>

            <div class="summary-pill avg-pill">
              <span class="pill-icon">⚖️</span>
              <span class="pill-label">核心均值：</span>
              <span class="pill-value text-blue font-semibold">
                {{ coreDiagnostics.avgVal.toFixed(1) }}%
              </span>
            </div>

            <div v-if="coreDiagnostics.isSkewed" class="summary-pill desc-pill">
              <span class="pill-icon">⚠️</span>
              <span class="pill-label text-red">
                检测到单核高负荷倾斜，可能存在单线程计算瓶颈或密集任务锁定
              </span>
            </div>
          </div>

          <!-- 视图 A：热力方块矩阵 (Heatmap Matrix) -->
          <div v-if="activeCoreView === 'heatmap'" class="cores-heatmap-container">
            <div class="heatmap-tiles-grid">
              <el-tooltip
                v-for="(usage, index) in server.cpu?.coreDetails || []"
                :key="index"
                placement="top"
                :show-after="80"
              >
                <template #content>
                  <div class="heatmap-tip-content">
                    <div class="font-bold">CPU 核心 #{{ index }}</div>
                    <div>
                      当前负载：<span :style="{ color: getUsageColor(usage || 0) }">{{ (usage || 0).toFixed(1) }}%</span>
                    </div>
                    <div>状态研判：{{ getCoreStatusText(usage || 0) }}</div>
                    <div v-if="index === coreDiagnostics.maxIndex" class="text-orange mt-1">
                      🔥 当前全机负载最高核心
                    </div>
                  </div>
                </template>

                <div
                  class="heatmap-tile"
                  :class="{
                    'is-max-core': index === coreDiagnostics.maxIndex && coreDiagnostics.maxVal >= 60,
                    'is-hot': (usage || 0) >= 80
                  }"
                  :style="{
                    backgroundColor: getHeatmapColor(usage || 0),
                    color: getHeatmapTextColor(usage || 0)
                  }"
                >
                  <span class="tile-id">#{{ index }}</span>
                  <span class="tile-val">{{ (usage || 0).toFixed(0) }}%</span>
                  <span
                    v-if="index === coreDiagnostics.maxIndex && coreDiagnostics.maxVal >= 70"
                    class="tile-crown-dot"
                  />
                </div>
              </el-tooltip>
            </div>

            <!-- 热力色阶图例 -->
            <div class="heatmap-legend-bar">
              <span class="legend-text">负载热度：</span>
              <div class="legend-scale">
                <span class="legend-block block-idle">0-10% 空闲</span>
                <span class="legend-block block-light">10-25%</span>
                <span class="legend-block block-normal">25-50% 正常</span>
                <span class="legend-block block-mid">50-75% 中等</span>
                <span class="legend-block block-busy">75-90% 繁忙</span>
                <span class="legend-block block-full">90%+ 满载</span>
              </div>
            </div>
          </div>

          <!-- 视图 B：详细卡片模式 (Card Mode) -->
          <div
            v-else
            class="cores-cards-container"
            :class="{ 'few-cores': coreCount <= 4 }"
          >
            <div
              v-for="(usage, index) in displayedCores"
              :key="index"
              class="core-detail-card"
              :class="{
                'is-extreme': index === coreDiagnostics.maxIndex && coreDiagnostics.maxVal >= 70
              }"
            >
              <div class="card-header-line">
                <div class="core-tag-box">
                  <span class="core-number">Core #{{ index }}</span>
                  <el-tag
                    size="small"
                    :type="getUsageColor(usage || 0) === '#f56c6c' ? 'danger' : ((usage || 0) >= 50 ? 'warning' : 'success')"
                    effect="plain"
                    class="status-tag"
                  >
                    {{ getCoreStatusText(usage || 0) }}
                  </el-tag>
                </div>
                <span
                  class="core-exact-percent"
                  :style="{ color: getUsageColor(usage || 0) }"
                >
                  {{ (usage || 0).toFixed(1) }}%
                </span>
              </div>

              <div class="core-track">
                <div
                  class="core-fill"
                  :style="{
                    width: `${Math.min(100, Math.max(0, usage || 0))}%`,
                    background: getCoreGradient(usage || 0)
                  }"
                />
              </div>

              <div v-if="coreCount <= 8" class="card-footer-meta">
                <span class="meta-label">相对均值偏离：</span>
                <span
                  class="meta-val font-semibold"
                  :class="(usage || 0) >= coreDiagnostics.avgVal ? 'text-orange' : 'text-green'"
                >
                  {{ ((usage || 0) - coreDiagnostics.avgVal) >= 0 ? '+' : '' }}{{ ((usage || 0) - coreDiagnostics.avgVal).toFixed(1) }}%
                </span>
              </div>
            </div>
          </div>

          <!-- 展开更多核心折叠按钮 -->
          <div
            v-if="activeCoreView === 'card' && coreCount > 8"
            class="expand-toggle-bar"
          >
            <el-button
              type="primary"
              plain
              size="small"
              @click="showAllCores = !showAllCores"
            >
              {{ showAllCores ? '收起部分核心 (仅显示前 8 核)' : `展开查看全部 (${coreCount} 核)` }}
              <el-icon class="el-icon--right">
                <ArrowUp v-if="showAllCores" />
                <ArrowDown v-else />
              </el-icon>
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- ================= 磁盘存储与挂载分区 (全新升级：少盘饱满卡片 + 多盘紧凑表格 + Inode监控) ================= -->
    <el-row :gutter="16" class="disk-section-row">
      <el-col :span="24">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="section-card-header flex-wrap gap-3">
              <div class="header-left">
                <el-icon class="section-icon disk"><FolderOpened /></el-icon>
                <span class="section-title">磁盘存储与挂载分区</span>
                <el-tag size="small" type="info" class="ml-2">
                  有效物理/数据分区 {{ storagePool.diskCount }} 个
                </el-tag>
                <el-tag
                  v-if="storagePool.alertCount > 0"
                  size="small"
                  type="danger"
                  effect="dark"
                  class="ml-2"
                >
                  ⚠️ {{ storagePool.alertCount }} 个分区高负荷预警
                </el-tag>
              </div>

              <!-- 右侧控制栏：告警筛选、排序、搜索、视图切换 -->
              <div class="header-right flex-wrap gap-2">
                <el-input
                  v-if="storagePool.diskCount > 3"
                  v-model="diskSearchQuery"
                  placeholder="搜索挂载点 / 设备"
                  size="small"
                  clearable
                  :prefix-icon="Search"
                  style="width: 170px"
                />

                <el-checkbox
                  v-if="storagePool.diskCount > 2"
                  v-model="diskOnlyAlert"
                  label="仅看告警 (≥80%)"
                  size="small"
                  class="mr-2"
                />

                <el-select
                  v-if="storagePool.diskCount > 2"
                  v-model="diskSortBy"
                  size="small"
                  style="width: 130px"
                >
                  <el-option value="usage_desc" label="使用率从高到低" />
                  <el-option value="total_desc" label="总容量从大到小" />
                  <el-option value="mountpoint" label="挂载路径排序" />
                </el-select>

                <el-radio-group v-model="diskViewMode" size="small">
                  <el-radio-button value="auto">
                    自适应 ({{ activeDiskView === 'table' ? '表格' : '卡片' }})
                  </el-radio-button>
                  <el-radio-button value="card">卡片网格</el-radio-button>
                  <el-radio-button value="table">紧凑表格</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>

          <!-- 模式一：详细表格视图 (Table Mode) - 适用于盘数较多或需要精确横向对比 -->
          <div v-if="activeDiskView === 'table'" class="disk-table-container">
            <el-table
              :data="filteredDisks"
              stripe
              style="width: 100%"
              class="custom-disk-table"
            >
              <el-table-column prop="mountpoint" label="挂载路径" min-width="150">
                <template #default="{ row }">
                  <div class="disk-table-mount">
                    <span class="mount-name font-bold">{{ row.mountpoint }}</span>
                    <span class="device-tag text-muted text-xs">({{ row.device }})</span>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="fstype" label="文件系统" width="100">
                <template #default="{ row }">
                  <el-tag size="small" effect="plain" type="info">
                    {{ row.fstype || 'ext4' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column label="空间使用率" min-width="180">
                <template #default="{ row }">
                  <el-progress
                    :percentage="Number((row.usedPercent || 0).toFixed(1))"
                    :stroke-width="10"
                    :color="getUsageColor(row.usedPercent || 0)"
                    :format="() => `${(row.usedPercent || 0).toFixed(1)}%`"
                  />
                </template>
              </el-table-column>

              <el-table-column label="已用 / 总量" min-width="160">
                <template #default="{ row }">
                  <span>{{ row.usedFormat }} / {{ row.totalFormat }}</span>
                </template>
              </el-table-column>

              <el-table-column label="剩余可用空间" min-width="130">
                <template #default="{ row }">
                  <span class="text-green font-semibold">{{ row.freeFormat }}</span>
                </template>
              </el-table-column>

              <el-table-column label="Inode 占用率" min-width="140">
                <template #default="{ row }">
                  <div v-if="row.inodesTotal && row.inodesTotal > 0" class="flex items-center gap-2">
                    <el-progress
                      :percentage="Number((row.inodesUsedPercent || 0).toFixed(1))"
                      :stroke-width="6"
                      :color="getUsageColor(row.inodesUsedPercent || 0)"
                      class="flex-1"
                      :show-text="false"
                    />
                    <span class="text-xs" :style="{ color: getUsageColor(row.inodesUsedPercent || 0) }">
                      {{ (row.inodesUsedPercent || 0).toFixed(0) }}%
                    </span>
                  </div>
                  <span v-else class="text-muted text-xs">-</span>
                </template>
              </el-table-column>

              <el-table-column label="健康研判" width="100" align="center">
                <template #default="{ row }">
                  <el-tag
                    :type="(row.usedPercent || 0) >= 85 ? 'danger' : ((row.usedPercent || 0) >= 70 ? 'warning' : 'success')"
                    size="small"
                    effect="light"
                  >
                    {{ (row.usedPercent || 0) >= 85 ? '高危' : ((row.usedPercent || 0) >= 70 ? '偏紧' : '充裕') }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- 模式二：卡片视图 (Card Mode) - 少盘 (1~2盘) 宽幅饱满展示，多盘网格化 -->
          <div
            v-else
            class="disk-cards-grid"
            :class="{ 'few-disks-grid': filteredDisks.length <= 2 }"
          >
            <!-- 极少盘场景（1~2 盘）：宽幅多维度大卡片，彻底避免单薄留白 -->
            <template v-if="filteredDisks.length <= 2">
              <div
                v-for="(item, index) in filteredDisks"
                :key="index"
                class="wide-disk-card"
                :class="{ 'is-high-usage': (item.usedPercent || 0) >= 85 }"
              >
                <!-- 左侧仪表大环 -->
                <div class="wide-disk-left">
                  <el-progress
                    type="dashboard"
                    :percentage="Number((item.usedPercent || 0).toFixed(1))"
                    :color="getUsageColor(item.usedPercent || 0)"
                    :stroke-width="12"
                    :width="140"
                  >
                    <template #default="{ percentage }">
                      <div class="gauge-center-content">
                        <span
                          class="gauge-center-val"
                          :style="{ color: getUsageColor(percentage) }"
                        >
                          {{ percentage.toFixed(1) }}%
                        </span>
                        <span class="gauge-center-label">空间已用</span>
                      </div>
                    </template>
                  </el-progress>
                </div>

                <!-- 右侧详细指标区 -->
                <div class="wide-disk-right">
                  <div class="disk-meta-header">
                    <div class="disk-title-group">
                      <span class="disk-mount-title">{{ item.mountpoint }}</span>
                      <span class="disk-device-badge">{{ item.device }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <el-tag size="small" effect="plain" type="info">
                        {{ item.fstype || 'ext4' }}
                      </el-tag>
                      <el-tag
                        size="small"
                        :type="(item.usedPercent || 0) >= 85 ? 'danger' : 'success'"
                        effect="light"
                      >
                        {{ (item.usedPercent || 0) >= 85 ? '空间紧张' : '读写正常 (rw)' }}
                      </el-tag>
                    </div>
                  </div>

                  <!-- 容量数值进度 -->
                  <div class="wide-progress-box">
                    <div class="progress-labels">
                      <span class="text-xs text-muted">存储空间容量</span>
                      <span class="text-xs font-semibold">
                        已用 {{ item.usedFormat }} / 剩余 <span class="text-green">{{ item.freeFormat }}</span> (共 {{ item.totalFormat }})
                      </span>
                    </div>
                    <el-progress
                      :percentage="Number((item.usedPercent || 0).toFixed(1))"
                      :stroke-width="10"
                      :color="getUsageColor(item.usedPercent || 0)"
                      :show-text="false"
                    />
                  </div>

                  <!-- Inode 状态条 (业界防故障关键指标) -->
                  <div class="wide-inode-box">
                    <div class="progress-labels">
                      <span class="text-xs text-muted">Inode 索引节点</span>
                      <span class="text-xs font-semibold">
                        {{ item.inodesTotal ? `已用 ${item.inodesUsed} / 总量 ${item.inodesTotal} (${(item.inodesUsedPercent || 0).toFixed(1)}%)` : '不适用或无需统计' }}
                      </span>
                    </div>
                    <el-progress
                      v-if="item.inodesTotal"
                      :percentage="Number((item.inodesUsedPercent || 0).toFixed(1))"
                      :stroke-width="6"
                      :color="getUsageColor(item.inodesUsedPercent || 0)"
                      :show-text="false"
                    />
                  </div>
                </div>
              </div>
            </template>

            <!-- 适中盘数场景（3~6 盘）：精简卡片网格 -->
            <template v-else>
              <div
                v-for="(item, index) in filteredDisks"
                :key="index"
                class="disk-partition-item"
                :class="{ 'is-high-usage': (item.usedPercent || 0) >= 85 }"
              >
                <div class="disk-part-top">
                  <div class="disk-mount-wrap">
                    <span class="disk-mount-point">{{ item.mountpoint }}</span>
                    <span class="disk-device-badge">{{ item.device }}</span>
                  </div>
                  <el-tag size="small" effect="plain" type="info">
                    {{ item.fstype || 'ext4' }}
                  </el-tag>
                </div>

                <div class="disk-progress-row">
                  <el-progress
                    :percentage="Number((item.usedPercent || 0).toFixed(1))"
                    :stroke-width="10"
                    :color="getUsageColor(item.usedPercent || 0)"
                    :format="() => `${(item.usedPercent || 0).toFixed(1)}%`"
                  />
                </div>

                <div class="disk-part-stats">
                  <div class="stat-col">
                    <span class="stat-k">已用容量</span>
                    <span
                      class="stat-v"
                      :style="{ color: getUsageColor(item.usedPercent || 0) }"
                    >
                      {{ item.usedFormat }}
                    </span>
                  </div>
                  <div class="stat-col">
                    <span class="stat-k">可用空间</span>
                    <span class="stat-v text-green">{{ item.freeFormat }}</span>
                  </div>
                  <div class="stat-col">
                    <span class="stat-k">总容量</span>
                    <span class="stat-v">{{ item.totalFormat }}</span>
                  </div>
                </div>

                <div v-if="item.inodesTotal" class="disk-card-inode-meta">
                  <span class="text-xs text-muted">Inode 占用：</span>
                  <span class="text-xs font-semibold" :style="{ color: getUsageColor(item.inodesUsedPercent || 0) }">
                    {{ (item.inodesUsedPercent || 0).toFixed(1) }}%
                  </span>
                </div>
              </div>
            </template>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 网络接口与流量吞吐 -->
    <el-row :gutter="16" class="network-section-row">
      <el-col :span="24">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="section-card-header">
              <div class="header-left">
                <el-icon class="section-icon network"><Connection /></el-icon>
                <span class="section-title">网络接口与实时吞吐监控</span>
                <el-tag size="small" type="success" effect="light" class="ml-2">
                  实时下行：↓ {{ totalNetworkRate.recvRateFormat }} / 实时上行：↑ {{ totalNetworkRate.sentRateFormat }}
                </el-tag>
              </div>
            </div>
          </template>

          <el-table
            :data="server.network || []"
            stripe
            style="width: 100%"
            class="custom-network-table"
          >
            <el-table-column prop="name" label="网卡接口" min-width="130">
              <template #default="{ row }">
                <div class="iface-name-col">
                  <span class="iface-name">{{ row.name }}</span>
                  <el-tag
                    :type="row.isUp ? 'success' : 'info'"
                    size="small"
                    effect="plain"
                  >
                    {{ row.isUp ? 'UP' : 'DOWN' }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>

            <el-table-column label="实时下行速率 (入网)" min-width="160">
              <template #default="{ row }">
                <div class="rate-cell text-green">
                  <span class="arrow font-bold mr-1">↓</span>
                  <span class="rate-val font-semibold">
                    {{ networkRates[row.name]?.recvRateFormat || '0 B/s' }}
                  </span>
                </div>
              </template>
            </el-table-column>

            <el-table-column label="实时上行速率 (出网)" min-width="160">
              <template #default="{ row }">
                <div class="rate-cell text-blue">
                  <span class="arrow font-bold mr-1">↑</span>
                  <span class="rate-val font-semibold">
                    {{ networkRates[row.name]?.sentRateFormat || '0 B/s' }}
                  </span>
                </div>
              </template>
            </el-table-column>

            <el-table-column prop="recvFormat" label="累计接收流量" min-width="140" />
            <el-table-column prop="sentFormat" label="累计发送流量" min-width="140" />
            <el-table-column prop="packetsRecv" label="接收包数" min-width="110" />
            <el-table-column prop="packetsSent" label="发送包数" min-width="110" />

            <el-table-column label="连接状态" width="110" align="center">
              <template #default="{ row }">
                <el-tag
                  :type="row.isUp ? 'success' : 'danger'"
                  size="small"
                  effect="light"
                >
                  {{ row.isUp ? '● 活跃' : '○ 断开' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 服务器详细参数规格 -->
    <el-row :gutter="16" class="specs-section-row">
      <el-col :span="24">
        <el-card shadow="hover" class="section-card">
          <template #header>
            <div class="section-card-header">
              <div class="header-left">
                <el-icon class="section-icon server"><Platform /></el-icon>
                <span class="section-title">服务器环境参数与规格</span>
              </div>
            </div>
          </template>

          <el-descriptions :column="4" border class="server-specs-desc">
            <el-descriptions-item label="主机名称">
              {{ server.general?.hostname || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="操作系统">
              {{ server.general?.os || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="系统架构">
              {{ server.general?.arch || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="内核版本">
              {{ server.general?.kernel || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="Go 运行环境">
              {{ server.general?.goVersion || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="CPU 规格型号">
              {{ server.cpu?.modelName || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="CPU 标称频率">
              {{ server.cpu?.frequency ? `${server.cpu.frequency.toFixed(0)} MHz` : '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="系统负载 (1/5/15m)">
              {{ server.load?.load1?.toFixed(2) || '0.00' }} /
              {{ server.load?.load5?.toFixed(2) || '0.00' }} /
              {{ server.load?.load15?.toFixed(2) || '0.00' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.system-monitor-container {
  padding: 16px 20px 30px;
  background-color: var(--el-bg-color-page, #f8fafc);
  min-height: calc(100vh - 90px);
}

/* 顶部操作控制条 */
.monitor-control-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  padding: 18px 24px;
  background: var(--el-bg-color-overlay, #ffffff);
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.bar-left {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}

.title-with-icon {
  display: flex;
  align-items: center;
  gap: 14px;
}

.title-icon-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: #ffffff;
  font-size: 22px;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.title-text-group .bar-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--el-text-color-primary, #0f172a);
  line-height: 1.3;
}

.title-text-group .bar-subtitle {
  font-size: 12px;
  color: var(--el-text-color-secondary, #64748b);
  margin-top: 2px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 14px;
}

.status-pulse-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 20px;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  box-shadow: 0 0 8px currentColor;
  animation: pulse-glow 1.8s infinite;
}

@keyframes pulse-glow {
  0%,
  100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.3);
    opacity: 0.6;
  }
}

.update-time-text {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--el-text-color-secondary, #64748b);
}

.bar-right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.control-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.control-label {
  font-size: 13px;
  color: var(--el-text-color-regular, #475569);
  font-weight: 500;
}

.interval-select {
  width: 90px;
}

/* 顶部四大卡片 */
.top-metrics-row {
  margin-bottom: 20px;
}

.metric-card {
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  background: var(--el-bg-color-overlay, #ffffff);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  margin-bottom: 16px;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.08);
}

.card-top-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.header-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  color: #ffffff;
  font-size: 18px;
  flex-shrink: 0;
}

.header-badge.cpu {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
}

.header-badge.memory {
  background: linear-gradient(135deg, #10b981, #047857);
}

.header-badge.disk {
  background: linear-gradient(135deg, #f59e0b, #b45309);
}

.header-badge.system {
  background: linear-gradient(135deg, #8b5cf6, #6d28d9);
}

.header-titles {
  flex: 1;
  min-width: 0;
}

.card-main-title {
  display: block;
  font-size: 15px;
  font-weight: 700;
  color: var(--el-text-color-primary, #0f172a);
}

.card-sub-title {
  display: block;
  font-size: 11px;
  color: var(--el-text-color-secondary, #64748b);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.header-tag {
  flex-shrink: 0;
}

.gauge-chart-container {
  height: 165px;
  width: 100%;
}

.disk-progress-container {
  height: 165px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gauge-center-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.gauge-center-val {
  font-size: 22px;
  font-weight: 800;
  line-height: 1.2;
}

.gauge-center-label {
  font-size: 11px;
  color: var(--el-text-color-secondary, #64748b);
  margin-top: 2px;
}

/* 运行时间卡片展示区 */
.uptime-display-wrap {
  height: 165px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  gap: 12px;
}

.uptime-icon-badge {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(139, 92, 246, 0.12);
  color: #8b5cf6;
  font-size: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.uptime-text-block .uptime-title {
  display: block;
  font-size: 12px;
  color: var(--el-text-color-secondary, #64748b);
  margin-bottom: 4px;
}

.uptime-text-block .uptime-val {
  display: block;
  font-size: 18px;
  font-weight: 700;
  color: #8b5cf6;
}

/* 卡片底栏 */
.card-info-footer {
  border-top: 1px dashed var(--el-border-color-lighter, #e2e8f0);
  padding-top: 10px;
  margin-top: 4px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  padding: 3px 0;
}

.info-item .label {
  color: var(--el-text-color-secondary, #64748b);
}

.info-item .val {
  font-weight: 600;
  color: var(--el-text-color-primary, #0f172a);
}

.load-tag-wrap {
  display: flex;
  align-items: center;
  gap: 4px;
}

.load-tag-wrap .sep {
  color: var(--el-text-color-placeholder, #94a3b8);
  font-weight: normal;
}

/* 通用板块卡片 */
.section-card {
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  background: var(--el-bg-color-overlay, #ffffff);
  margin-bottom: 20px;
}

.section-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.section-card-header .header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.section-icon {
  font-size: 20px;
  display: flex;
  align-items: center;
}

.section-icon.chart {
  color: #3b82f6;
}

.section-icon.cpu {
  color: #3b82f6;
}

.section-icon.disk {
  color: #f59e0b;
}

.section-icon.network {
  color: #10b981;
}

.section-icon.server {
  color: #8b5cf6;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--el-text-color-primary, #0f172a);
}

.chart-custom-header .points-label {
  font-size: 12px;
  color: var(--el-text-color-secondary, #64748b);
}

.history-chart-dom {
  height: 320px;
  width: 100%;
}

/* ================= 多核矩阵与热力图样式 ================= */
.core-summary-strip {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  padding: 10px 14px;
  background: var(--el-fill-color-light, #f8fafc);
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 12px;
}

.summary-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.summary-pill .pill-icon {
  font-size: 14px;
}

.summary-pill .pill-label {
  color: var(--el-text-color-secondary, #64748b);
}

.summary-pill .pill-value {
  font-size: 13px;
}

/* 热力方块容器 */
.cores-heatmap-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.heatmap-tiles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(56px, 1fr));
  gap: 8px;
}

.heatmap-tile {
  aspect-ratio: 1;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;
  position: relative;
  border: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.heatmap-tile:hover {
  transform: translateY(-2px) scale(1.05);
  z-index: 5;
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.15);
}

.heatmap-tile.is-max-core {
  outline: 2px solid #ef4444;
  outline-offset: 1px;
}

.tile-id {
  font-size: 10px;
  opacity: 0.85;
  line-height: 1;
}

.tile-val {
  font-weight: 700;
  font-size: 12px;
  margin-top: 3px;
  line-height: 1.1;
}

.tile-crown-dot {
  position: absolute;
  top: 3px;
  right: 3px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background-color: #ffffff;
  box-shadow: 0 0 4px #000000;
}

.heatmap-tip-content {
  line-height: 1.6;
  font-size: 12px;
}

.heatmap-legend-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  padding-top: 6px;
  border-top: 1px dashed var(--el-border-color-lighter, #e2e8f0);
}

.legend-text {
  font-size: 12px;
  color: var(--el-text-color-secondary, #64748b);
}

.legend-scale {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.legend-block {
  font-size: 11px;
  padding: 2px 7px;
  border-radius: 4px;
  font-weight: 500;
}

.legend-block.block-idle {
  background: var(--el-fill-color, #e2e8f0);
  color: var(--el-text-color-secondary, #64748b);
}

.legend-block.block-light {
  background: #34d399;
  color: #ffffff;
}

.legend-block.block-normal {
  background: #10b981;
  color: #ffffff;
}

.legend-block.block-mid {
  background: #eab308;
  color: #ffffff;
}

.legend-block.block-busy {
  background: #f97316;
  color: #ffffff;
}

.legend-block.block-full {
  background: #ef4444;
  color: #ffffff;
}

/* 详细卡片视图 */
.cores-cards-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.cores-cards-container.few-cores {
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.core-detail-card {
  background: var(--el-fill-color-light, #f8fafc);
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  border-radius: 8px;
  padding: 12px 14px;
  transition: all 0.25s ease;
}

.core-detail-card:hover {
  transform: translateY(-2px);
  border-color: var(--el-color-primary, #3b82f6);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.core-detail-card.is-extreme {
  border-left: 3px solid #ef4444;
}

.card-header-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.core-tag-box {
  display: flex;
  align-items: center;
  gap: 6px;
}

.core-number {
  font-weight: 600;
  font-size: 13px;
  color: var(--el-text-color-primary, #0f172a);
}

.core-exact-percent {
  font-weight: 800;
  font-size: 14px;
}

.core-track {
  width: 100%;
  height: 8px;
  background: var(--el-border-color-lighter, #e2e8f0);
  border-radius: 4px;
  overflow: hidden;
}

.core-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.card-footer-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  margin-top: 8px;
  padding-top: 6px;
  border-top: 1px dashed var(--el-border-color-lighter, #e2e8f0);
}

.card-footer-meta .meta-label {
  color: var(--el-text-color-secondary, #64748b);
}

.expand-toggle-bar {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

/* ================= 磁盘存储与挂载分区全新样式 ================= */
.disk-table-container {
  width: 100%;
  overflow-x: auto;
}

.custom-disk-table .disk-table-mount {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.disk-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

/* 极少盘场景（1~2 盘）大号宽幅自适应布局 */
.disk-cards-grid.few-disks-grid {
  grid-template-columns: repeat(auto-fit, minmax(420px, 1fr));
  gap: 18px;
}

.wide-disk-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 18px 22px;
  background: var(--el-fill-color-light, #f8fafc);
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.wide-disk-card:hover {
  transform: translateY(-2px);
  border-color: var(--el-color-primary, #3b82f6);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.06);
}

.wide-disk-card.is-high-usage {
  border-left: 4px solid #ef4444;
}

.wide-disk-left {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wide-disk-right {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.disk-meta-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.disk-title-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.disk-mount-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--el-text-color-primary, #0f172a);
}

.wide-progress-box,
.wide-inode-box {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.progress-labels {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 多盘常规卡片 */
.disk-partition-item {
  padding: 16px;
  background: var(--el-fill-color-light, #f8fafc);
  border: 1px solid var(--el-border-color-lighter, #e2e8f0);
  border-radius: 10px;
  transition: all 0.2s ease;
}

.disk-partition-item:hover {
  transform: translateY(-2px);
  border-color: var(--el-border-color, #cbd5e1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
}

.disk-partition-item.is-high-usage {
  border-left: 3px solid #ef4444;
}

.disk-part-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.disk-mount-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.disk-mount-point {
  font-size: 15px;
  font-weight: 700;
  color: var(--el-text-color-primary, #0f172a);
}

.disk-device-badge {
  font-size: 11px;
  color: var(--el-text-color-secondary, #64748b);
  background: var(--el-fill-color, #e2e8f0);
  padding: 2px 6px;
  border-radius: 4px;
}

.disk-progress-row {
  margin-bottom: 12px;
}

.disk-part-stats {
  display: flex;
  justify-content: space-between;
  border-top: 1px dashed var(--el-border-color-lighter, #e2e8f0);
  padding-top: 8px;
}

.stat-col {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-k {
  font-size: 11px;
  color: var(--el-text-color-secondary, #64748b);
}

.stat-v {
  font-size: 12px;
  font-weight: 600;
  color: var(--el-text-color-primary, #0f172a);
}

.disk-card-inode-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 6px;
  padding-top: 6px;
  border-top: 1px dashed var(--el-border-color-lighter, #e2e8f0);
}

/* 网卡列表 */
.iface-name-col {
  display: flex;
  align-items: center;
  gap: 8px;
}

.iface-name {
  font-weight: 600;
  color: var(--el-text-color-primary, #0f172a);
}

.rate-cell {
  display: inline-flex;
  align-items: center;
}

/* 颜色辅助类 */
.text-green {
  color: #10b981 !important;
}

.text-blue {
  color: #3b82f6 !important;
}

.text-orange {
  color: #f59e0b !important;
}

.text-red {
  color: #f56c6c !important;
}

.text-muted {
  color: var(--el-text-color-placeholder, #94a3b8) !important;
}

.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .system-monitor-container {
    padding: 12px 12px 24px;
  }

  .monitor-control-bar {
    padding: 14px 16px;
  }

  .bar-left,
  .bar-right {
    width: 100%;
    justify-content: space-between;
  }

  .chart-custom-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .chart-custom-header .header-right {
    margin-left: 0;
  }

  .cores-cards-container {
    grid-template-columns: repeat(2, 1fr);
  }

  .cores-cards-container.few-cores {
    grid-template-columns: 1fr;
  }

  .disk-cards-grid {
    grid-template-columns: 1fr;
  }

  .disk-cards-grid.few-disks-grid {
    grid-template-columns: 1fr;
  }

  .wide-disk-card {
    flex-direction: column;
    align-items: flex-start;
  }

  .wide-disk-left {
    width: 100%;
  }

  .heatmap-tiles-grid {
    grid-template-columns: repeat(auto-fill, minmax(46px, 1fr));
  }
}
</style>
