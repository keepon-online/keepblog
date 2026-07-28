<script setup lang="ts">
import { ref, computed, watch, type Ref, defineProps } from "vue";
import { useAppStoreHook } from "@/store/modules/app";
import {
  delay,
  useDark,
  useECharts,
  type EchartOptions
} from "@pureadmin/utils";
import * as echarts from "echarts/core";

const props = defineProps<{
  barData: any[];
}>();

const { isDark } = useDark();
const theme: EchartOptions["theme"] = computed(() => {
  return isDark.value ? "dark" : "light";
});

const barChartRef = ref<HTMLDivElement | null>(null);
const { setOptions, resize } = useECharts(barChartRef as Ref<HTMLDivElement>, {
  theme
});

watch(
  () => useAppStoreHook().getSidebarStatus,
  () => {
    delay(600).then(() => resize());
  }
);

// 更新图表的函数
const updateChartFromProps = (data: any[]) => {
  const xData = data.map(e => e.name);
  const yData = data.map(e => e.value);

  setOptions(
    {
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "shadow"
        }
      },
      grid: {
        top: "20px",
        left: "40px",
        right: "20px",
        bottom: "40px"
      },
      legend: {
        //@ts-expect-error
        right: true,
        data: ["star"]
      },
      xAxis: [
        {
          type: "category",
          axisTick: {
            alignWithLabel: true
          },
          axisLabel: {
            interval: 0,
            fontSize: 12
          },
          data: xData,
          triggerEvent: true
        }
      ],
      yAxis: [
        {
          type: "value",
          triggerEvent: true,
          axisLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          splitLine: {
            show: true,
            lineStyle: {
              type: "dashed"
            }
          }
        }
      ],
      series: [
        {
          type: "bar",
          barWidth: "40%",
          itemStyle: {
            borderRadius: [4, 4, 0, 0],
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: "#409EFF"
              },
              {
                offset: 1,
                color: "#53a7ff"
              }
            ])
          },
          data: yData,
          emphasis: {
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                  offset: 0,
                  color: "#53a7ff"
                },
                {
                  offset: 1,
                  color: "#409EFF"
                }
              ])
            }
          }
        }
      ],
      addTooltip: true
    },
    {
      name: "click",
      callback: () => {
        // 图表点击事件处理
      }
    }
  );
};

// 监听数据变化
watch(
  () => props.barData,
  newData => {
    if (newData && newData.length > 0) {
      updateChartFromProps(newData);
    }
  },
  { immediate: true }
);
</script>

<template>
  <div ref="barChartRef" style="width: 100%; height: 300px" />
</template>
