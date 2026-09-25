<template>
  <div class="line-chart-wrapper">
    <div v-if="lineData && lineData.length > 0" ref="lineChartRef" class="line-chart" />
    <div v-else class="empty-chart">
      <el-empty description="暂无访问流量数据" :image-size="80" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, type Ref } from "vue";
import { useAppStoreHook } from "@/store/modules/app";
import { delay, useDark, useECharts, type EchartOptions } from "@pureadmin/utils";
import * as echarts from "echarts/core";

defineOptions({
  name: "TrafficLineChart"
});

const props = defineProps<{
  lineData: Array<{
    name: string;
    pv: number;
    uv: number;
  }>;
}>();

const { isDark } = useDark();

const theme: EchartOptions["theme"] = computed(() => {
  return isDark.value ? "dark" : "default";
});

const lineChartRef = ref<HTMLDivElement | null>(null);
const { setOptions, resize } = useECharts(lineChartRef as Ref<HTMLDivElement>, {
  theme
});

// 监听侧边栏折叠/展开
watch(
  () => useAppStoreHook().getSidebarStatus,
  () => {
    delay(400).then(() => resize());
  }
);

// 监听深色模式切换
watch(isDark, () => {
  if (props.lineData && props.lineData.length > 0) {
    updateChart(props.lineData);
  }
});

const updateChart = (data: Array<{ name: string; pv: number; uv: number }>) => {
  // 不在此处判空 ref：setOptions 内部经 nextTick 懒初始化，元素就绪后会自动补建实例；
  // 提前 return 会在 immediate watch（setup 阶段 ref 尚为 null）时永久丢失首次渲染
  const xData = data.map(e => e.name);
  const pvData = data.map(e => e.pv || 0);
  const uvData = data.map(e => e.uv || 0);

  const primaryColor = "#409eff";
  const successColor = "#67c23a";

  const options: echarts.EChartsCoreOption = {
    backgroundColor: "transparent",
    tooltip: {
      trigger: "axis",
      backgroundColor: isDark.value ? "rgba(24, 24, 28, 0.92)" : "rgba(255, 255, 255, 0.96)",
      borderColor: isDark.value ? "rgba(255, 255, 255, 0.12)" : "rgba(0, 0, 0, 0.08)",
      borderWidth: 1,
      padding: [10, 14],
      textStyle: {
        color: isDark.value ? "#f2f2f2" : "#303133",
        fontSize: 13
      },
      extraCssText: "box-shadow: 0 8px 24px -4px rgba(0,0,0,0.15); border-radius: 8px;",
      formatter: (params: any) => {
        if (!Array.isArray(params) || params.length === 0) return "";
        let str = `<div style="font-weight: 600; margin-bottom: 6px;">${params[0].axisValue} 访问明细</div>`;
        let pv = 0;
        let uv = 0;
        params.forEach(item => {
          if (item.seriesName === "浏览量 (PV)") pv = item.value;
          if (item.seriesName === "独立访客 (UV)") uv = item.value;
          str += `<div style="display: flex; align-items: center; justify-content: space-between; gap: 20px; font-size: 12px; margin-top: 3px;">
            <span>${item.marker} ${item.seriesName}</span>
            <span style="font-weight: 600;">${item.value.toLocaleString()}</span>
          </div>`;
        });
        if (uv > 0) {
          const ratio = (pv / uv).toFixed(1);
          str += `<div style="margin-top: 6px; padding-top: 4px; border-top: 1px dashed ${
            isDark.value ? "#444" : "#eee"
          }; font-size: 11px; color: #909399;">
            人均浏览量: <span style="font-weight: 600; color: ${primaryColor}">${ratio}</span> 页/人
          </div>`;
        }
        return str;
      }
    },
    legend: {
      right: "10px",
      top: "0px",
      icon: "circle",
      itemWidth: 8,
      itemHeight: 8,
      textStyle: {
        color: isDark.value ? "rgba(255, 255, 255, 0.65)" : "#606266",
        fontSize: 12
      },
      data: ["浏览量 (PV)", "独立访客 (UV)"]
    },
    grid: {
      top: "36px",
      left: "12px",
      right: "12px",
      bottom: "8px",
      containLabel: true
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: xData,
      axisLine: {
        lineStyle: {
          color: isDark.value ? "rgba(255, 255, 255, 0.1)" : "rgba(0, 0, 0, 0.08)"
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: isDark.value ? "rgba(255, 255, 255, 0.55)" : "#909399",
        fontSize: 12
      }
    },
    yAxis: {
      type: "value",
      splitLine: {
        lineStyle: {
          type: "dashed",
          color: isDark.value ? "rgba(255, 255, 255, 0.08)" : "rgba(0, 0, 0, 0.05)"
        }
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: isDark.value ? "rgba(255, 255, 255, 0.55)" : "#909399",
        fontSize: 12
      }
    },
    series: [
      {
        name: "浏览量 (PV)",
        type: "line",
        smooth: 0.35,
        symbol: "circle",
        symbolSize: 6,
        showSymbol: false,
        itemStyle: {
          color: primaryColor
        },
        lineStyle: {
          width: 3,
          color: primaryColor,
          shadowColor: "rgba(64, 158, 255, 0.3)",
          shadowBlur: 8,
          shadowOffsetY: 4
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "rgba(64, 158, 255, 0.32)" },
            { offset: 1, color: "rgba(64, 158, 255, 0.01)" }
          ])
        },
        data: pvData
      },
      {
        name: "独立访客 (UV)",
        type: "line",
        smooth: 0.35,
        symbol: "circle",
        symbolSize: 6,
        showSymbol: false,
        itemStyle: {
          color: successColor
        },
        lineStyle: {
          width: 2.5,
          color: successColor,
          shadowColor: "rgba(103, 194, 58, 0.3)",
          shadowBlur: 8,
          shadowOffsetY: 3
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "rgba(103, 194, 58, 0.28)" },
            { offset: 1, color: "rgba(103, 194, 58, 0.01)" }
          ])
        },
        data: uvData
      }
    ]
  };

  setOptions(options as any);
};

watch(
  () => props.lineData,
  newData => {
    if (newData && newData.length > 0) {
      updateChart(newData);
    }
  },
  { immediate: true, deep: true }
);

onMounted(() => {
  window.addEventListener("resize", resize);
});

onUnmounted(() => {
  window.removeEventListener("resize", resize);
});
</script>

<style lang="scss" scoped>
.line-chart-wrapper {
  width: 100%;
  height: 310px;
  position: relative;

  .line-chart {
    width: 100%;
    height: 100%;
  }

  .empty-chart {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
