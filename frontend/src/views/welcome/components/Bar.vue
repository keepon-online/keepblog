<template>
  <div class="bar-chart-wrapper">
    <div v-if="barData && barData.length > 0" ref="barChartRef" class="bar-chart" />
    <div v-else class="empty-chart">
      <el-empty description="暂无近期访问统计数据" :image-size="80" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, type Ref } from "vue";
import { useAppStoreHook } from "@/store/modules/app";
import { delay, useDark, useECharts, type EchartOptions } from "@pureadmin/utils";
import * as echarts from "echarts/core";

defineOptions({
  name: "WeekBarChart"
});

const props = defineProps<{
  barData: Array<{
    name: string;
    value: number;
  }>;
}>();

const { isDark } = useDark();

const theme: EchartOptions["theme"] = computed(() => {
  return isDark.value ? "dark" : "default";
});

const barChartRef = ref<HTMLDivElement | null>(null);
const { setOptions, resize } = useECharts(barChartRef as Ref<HTMLDivElement>, {
  theme
});

watch(
  () => useAppStoreHook().getSidebarStatus,
  () => {
    delay(400).then(() => resize());
  }
);

watch(isDark, () => {
  if (props.barData && props.barData.length > 0) {
    updateChart(props.barData);
  }
});

const updateChart = (data: Array<{ name: string; value: number }>) => {
  // 不在此处判空 ref：setOptions 内部经 nextTick 懒初始化，元素就绪后会自动补建实例；
  // 提前 return 会在 immediate watch（setup 阶段 ref 尚为 null）时永久丢失首次渲染
  const xData = data.map(e => e.name);
  const yData = data.map(e => e.value);

  const primaryColor = "#409eff";

  const options: echarts.EChartsCoreOption = {
    backgroundColor: "transparent",
    tooltip: {
      trigger: "axis",
      axisPointer: {
        type: "shadow",
        shadowStyle: {
          color: isDark.value ? "rgba(255, 255, 255, 0.03)" : "rgba(0, 0, 0, 0.03)"
        }
      },
      backgroundColor: isDark.value ? "rgba(24, 24, 28, 0.92)" : "rgba(255, 255, 255, 0.96)",
      borderColor: isDark.value ? "rgba(255, 255, 255, 0.12)" : "rgba(0, 0, 0, 0.08)",
      borderWidth: 1,
      padding: [8, 12],
      textStyle: {
        color: isDark.value ? "#f2f2f2" : "#303133",
        fontSize: 13
      },
      extraCssText: "box-shadow: 0 8px 24px -4px rgba(0,0,0,0.15); border-radius: 8px;",
      formatter: (params: any) => {
        if (!Array.isArray(params) || params.length === 0) return "";
        const item = params[0];
        return `<div style="font-weight: 600; margin-bottom: 4px;">${item.name}</div>
          <div style="font-size: 12px; color: ${isDark.value ? '#ccc' : '#606266'}">
            访问量: <span style="font-weight: 600; color: ${primaryColor}">${item.value.toLocaleString()}</span> 次 PV
          </div>`;
      }
    },
    grid: {
      top: "24px",
      left: "12px",
      right: "12px",
      bottom: "8px",
      containLabel: true
    },
    xAxis: {
      type: "category",
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
        name: "访问量 (PV)",
        type: "bar",
        barWidth: "36%",
        itemStyle: {
          borderRadius: [6, 6, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "#409eff" },
            { offset: 1, color: "rgba(64, 158, 255, 0.45)" }
          ])
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: "#66b1ff" },
              { offset: 1, color: "#409eff" }
            ]),
            shadowBlur: 8,
            shadowColor: "rgba(64, 158, 255, 0.4)"
          }
        },
        data: yData
      }
    ]
  };

  setOptions(options as any);
};

watch(
  () => props.barData,
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
.bar-chart-wrapper {
  width: 100%;
  height: 280px;
  position: relative;

  .bar-chart {
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
