<script setup lang="ts">
import { useIntervalFn } from "@vueuse/core";
import { ref, computed, watch, type Ref, defineProps } from "vue";
import { useAppStoreHook } from "@/store/modules/app";
import {
  delay,
  useDark,
  useECharts,
  type EchartOptions
} from "@pureadmin/utils";

const props = defineProps<{
  lineData: any[];
}>();

const { isDark } = useDark();

const theme: EchartOptions["theme"] = computed(() => {
  return isDark.value ? "dark" : "default";
});

const lineChartRef = ref<HTMLDivElement | null>(null);
const { setOptions, getInstance, resize } = useECharts(
  lineChartRef as Ref<HTMLDivElement>,
  { theme }
);

const xData = ref([]);
const pv = ref([]);
const uv = ref([]);
let a = 1;

useIntervalFn(() => {
  if (xData.value.length > 0) {
    if (a == xData.value.length - 24) {
      a = 0;
    }
    const instance = getInstance();
    if (instance) {
      instance.dispatchAction({
        type: "dataZoom",
        startValue: a,
        endValue: a + 24
      });
    }
    a++;
  }
}, 2000);

watch(
  () => useAppStoreHook().getSidebarStatus,
  () => {
    delay(600).then(() => resize());
  }
);

// 更新图表的函数
const updateChartFromProps = (data: any[]) => {
  xData.value = data.map(e => e.name);
  pv.value = data.map(e => e.pv);
  uv.value = data.map(e => e.uv);

  setOptions(
    {
      tooltip: {
        trigger: "axis",
        backgroundColor: isDark.value
          ? "rgba(0, 0, 0, 0.8)"
          : "rgba(255, 255, 255, 0.9)",
        borderColor: isDark.value ? "#303133" : "#e4e7ed",
        textStyle: {
          color: isDark.value
            ? "rgba(255, 255, 255, 0.9)"
            : "rgba(0, 0, 0, 0.9)"
        }
      },
      grid: {
        top: "30px",
        left: "40px",
        right: "20px",
        bottom: "40px"
      },
      legend: {
        //@ts-expect-error
        right: true,
        data: ["PV", "UV"],
        textStyle: {
          color: isDark.value
            ? "rgba(255, 255, 255, 0.7)"
            : "rgba(0, 0, 0, 0.65)"
        }
      },
      calculable: true,
      xAxis: [
        {
          triggerEvent: true,
          type: "category",
          splitLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          axisLine: {
            lineStyle: {
              color: isDark.value ? "#303133" : "#e4e7ed"
            }
          },
          axisLabel: {
            color: isDark.value
              ? "rgba(255, 255, 255, 0.7)"
              : "rgba(0, 0, 0, 0.65)"
          },
          data: xData.value
        }
      ],
      yAxis: [
        {
          triggerEvent: true,
          type: "value",
          splitLine: {
            show: true,
            lineStyle: {
              type: "dashed"
            }
          },
          axisLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          axisLabel: {
            color: isDark.value
              ? "rgba(255, 255, 255, 0.7)"
              : "rgba(0, 0, 0, 0.65)"
          }
        }
      ],
      dataZoom: [
        {
          type: "slider",
          show: false,
          realtime: true,
          startValue: 0,
          endValue: 24
        }
      ],
      series: [
        {
          name: "PV",
          type: "line",
          smooth: true,
          symbolSize: 6,
          symbol: "circle",
          color: "#f56c6c",
          lineStyle: {
            width: 3
          },
          markPoint: {
            label: {
              color: "#fff"
            },
            data: [
              {
                type: "max",
                name: "最大值"
              },
              {
                type: "min",
                name: "最小值"
              }
            ]
          },
          data: pv.value
        },
        {
          name: "UV",
          type: "line",
          smooth: true,
          symbolSize: 6,
          symbol: "circle",
          color: "#409EFF",
          lineStyle: {
            width: 3
          },
          markPoint: {
            label: {
              color: "#fff"
            },
            data: [
              {
                type: "max",
                name: "最大值"
              },
              {
                type: "min",
                name: "最小值"
              }
            ]
          },
          data: uv.value
        }
      ],
      addTooltip: true
    },
    {
      name: "click",
      callback: () => {
        // 图表点击事件处理
      }
    },
    {
      name: "contextmenu",
      callback: () => {
        // 右键菜单事件处理
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
  () => props.lineData,
  newData => {
    if (newData && newData.length > 0) {
      updateChartFromProps(newData);
    }
  },
  { immediate: true }
);
</script>

<template>
  <div ref="lineChartRef" style="width: 100%; height: 300px" />
</template>
