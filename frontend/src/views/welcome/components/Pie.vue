<script setup lang="ts">
import { ref, computed, watch, type Ref, defineProps } from "vue";
import { useAppStoreHook } from "@/store/modules/app";
import {
  delay,
  useDark,
  useECharts,
  type EchartOptions
} from "@pureadmin/utils";

const props = defineProps<{
  pieData: any[];
}>();

const { isDark } = useDark();

const theme: EchartOptions["theme"] = computed(() => {
  return isDark.value ? "dark" : "light";
});

const pieChartRef = ref<HTMLDivElement | null>(null);
const { setOptions, resize } = useECharts(pieChartRef as Ref<HTMLDivElement>, {
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
  setOptions(
    {
      tooltip: {
        trigger: "item",
        formatter: "{a} <br/>{b}: {c} ({d}%)"
      },
      legend: {
        icon: "circle",
        orient: "vertical",
        right: "20px",
        top: "middle",
        textStyle: {
          color: isDark.value
            ? "rgba(255, 255, 255, 0.7)"
            : "rgba(0, 0, 0, 0.65)"
        }
      },
      series: [
        {
          name: "分类",
          type: "pie",
          top: "0%",
          radius: ["40%", "70%"],
          center: ["35%", "50%"],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 4,
            borderColor: isDark.value ? "#141414" : "#fff",
            borderWidth: 2
          },
          label: {
            show: false,
            position: "center"
          },
          emphasis: {
            label: {
              show: true,
              fontSize: "14",
              fontWeight: "bold",
              formatter: "{b}\n\n{d}%"
            }
          },
          labelLine: {
            show: false
          },
          data: data
        }
      ]
    },
    {
      name: "click",
      callback: () => {
        // 图表点击事件处理
      }
    },
    // 点击空白处
    {
      type: "zrender",
      name: "click",
      callback: () => {
        // 空白区域点击处理
      }
    }
  );
};

// 监听数据变化
watch(
  () => props.pieData,
  newData => {
    if (newData && newData.length > 0) {
      updateChartFromProps(newData);
    }
  },
  { immediate: true }
);
</script>

<template>
  <div ref="pieChartRef" style="width: 100%; height: 300px" />
</template>
