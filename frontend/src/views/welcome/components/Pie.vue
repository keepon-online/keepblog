<template>
  <div class="pie-chart-wrapper">
    <div v-if="pieData && pieData.length > 0" ref="pieChartRef" class="pie-chart" />
    <div v-else class="empty-chart">
      <el-empty description="暂无分类数据" :image-size="80" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, type Ref } from "vue";
import { useAppStoreHook } from "@/store/modules/app";
import { delay, useDark, useECharts, type EchartOptions } from "@pureadmin/utils";
import * as echarts from "echarts/core";

defineOptions({
  name: "CategoryPieChart"
});

const props = defineProps<{
  pieData: Array<{
    name: string;
    value: number;
  }>;
}>();

const { isDark } = useDark();

const theme: EchartOptions["theme"] = computed(() => {
  return isDark.value ? "dark" : "default";
});

const pieChartRef = ref<HTMLDivElement | null>(null);
const { setOptions, resize } = useECharts(pieChartRef as Ref<HTMLDivElement>, {
  theme
});

watch(
  () => useAppStoreHook().getSidebarStatus,
  () => {
    delay(400).then(() => resize());
  }
);

watch(isDark, () => {
  if (props.pieData && props.pieData.length > 0) {
    updateChart(props.pieData);
  }
});

const palette = [
  "#409eff",
  "#67c23a",
  "#e6a23c",
  "#9b59b6",
  "#f56c6c",
  "#36a3f7",
  "#34bfa3",
  "#f4516c"
];

const updateChart = (data: Array<{ name: string; value: number }>) => {
  // 不在此处判空 ref：setOptions 内部经 nextTick 懒初始化，元素就绪后会自动补建实例；
  // 提前 return 会在 immediate watch（setup 阶段 ref 尚为 null）时永久丢失首次渲染
  const total = data.reduce((acc, cur) => acc + (cur.value || 0), 0);

  const options: echarts.EChartsCoreOption = {
    backgroundColor: "transparent",
    color: palette,
    tooltip: {
      trigger: "item",
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
        return `<div style="font-weight: 600; margin-bottom: 4px;">${params.name}</div>
          <div style="font-size: 12px; color: ${isDark.value ? '#ccc' : '#606266'}">
            文章篇数: <span style="font-weight: 600; color: #409eff">${params.value}</span> 篇<br/>
            分类占比: <span style="font-weight: 600; color: #67c23a">${params.percent}%</span>
          </div>`;
      }
    },
    legend: {
      orient: "vertical",
      right: "10px",
      top: "center",
      icon: "circle",
      itemWidth: 8,
      itemHeight: 8,
      textStyle: {
        color: isDark.value ? "rgba(255, 255, 255, 0.65)" : "#606266",
        fontSize: 12
      },
      formatter: (name: string) => {
        const item = data.find(d => d.name === name);
        const val = item ? item.value : 0;
        const pct = total > 0 ? ((val / total) * 100).toFixed(0) : "0";
        return `${name}  ${val}篇 (${pct}%)`;
      }
    },
    series: [
      {
        name: "分类占比",
        type: "pie",
        radius: ["50%", "72%"],
        center: ["36%", "50%"],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 6,
          borderColor: isDark.value ? "#1d1e1f" : "#ffffff",
          borderWidth: 2
        },
        label: {
          show: false,
          position: "center"
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: "bold",
            color: isDark.value ? "#fff" : "#303133",
            formatter: "{b}\n{c} 篇 ({d}%)"
          },
          scale: true,
          scaleSize: 6
        },
        labelLine: {
          show: false
        },
        data: data
      }
    ]
  };

  setOptions(options as any);
};

watch(
  () => props.pieData,
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
.pie-chart-wrapper {
  width: 100%;
  height: 310px;
  position: relative;

  .pie-chart {
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
