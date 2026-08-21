<template>
  <div class="map-container">
    <div ref="chartRef" class="map-chart" />
    <div v-if="overseasData.length > 0" class="overseas-section">
      <div class="overseas-header">
        <span class="overseas-title">🌍 海外访客分布</span>
        <span class="overseas-total">共 {{ overseasTotal }} 人次</span>
      </div>
      <div ref="overseasChartRef" class="overseas-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "ChinaMap"
});

import { ref, onMounted, watch, onUnmounted, computed } from "vue";
// 按需引入：全量 import "echarts" 会把整个 echarts 打进本组件的懒加载 chunk
import echarts from "@/utils/echarts";
import chinaJson from "./china.json";

const props = defineProps<{
  mapData: Array<{ name: string; value: number }>;
}>();

const chartRef = ref<HTMLElement | null>(null);
let chartInstance: echarts.ECharts | null = null;

// 注册中国地图
echarts.registerMap("china", chinaJson as any);

// 省份名称映射表：简称 -> 地图全称
const provinceNameMap: Record<string, string> = {
  北京: "北京市",
  天津: "天津市",
  上海: "上海市",
  重庆: "重庆市",
  河北: "河北省",
  山西: "山西省",
  辽宁: "辽宁省",
  吉林: "吉林省",
  黑龙江: "黑龙江省",
  江苏: "江苏省",
  浙江: "浙江省",
  安徽: "安徽省",
  福建: "福建省",
  江西: "江西省",
  山东: "山东省",
  河南: "河南省",
  湖北: "湖北省",
  湖南: "湖南省",
  广东: "广东省",
  海南: "海南省",
  四川: "四川省",
  贵州: "贵州省",
  云南: "云南省",
  陕西: "陕西省",
  甘肃: "甘肃省",
  青海: "青海省",
  台湾: "台湾省",
  内蒙古: "内蒙古自治区",
  广西: "广西壮族自治区",
  西藏: "西藏自治区",
  宁夏: "宁夏回族自治区",
  新疆: "新疆维吾尔自治区",
  香港: "香港特别行政区",
  澳门: "澳门特别行政区"
};

// 转换数据：将简称映射为地图所需的全称
const convertMapData = (data: Array<{ name: string; value: number }>) => {
  if (!data) return [];
  return data.map(item => ({
    name: provinceNameMap[item.name] || item.name,
    value: item.value
  }));
};

const initChart = () => {
  if (!chartRef.value) return;

  chartInstance = echarts.init(chartRef.value);
  updateChart();

  // 响应式
  window.addEventListener("resize", handleResize);
};

const handleResize = () => {
  chartInstance?.resize();
};

const updateChart = () => {
  if (!chartInstance) return;

  const convertedData = convertMapData(props.mapData);
  const maxValue = Math.max(
    ...(convertedData?.map(item => item.value) || [100]),
    100
  );

  const option = {
    title: {
      text: "访客地理分布",
      left: "center",
      textStyle: {
        color: "#333",
        fontSize: 16
      }
    },
    tooltip: {
      trigger: "item",
      formatter: (params: any) => {
        if (params.value) {
          return `${params.name}<br/>访客数: ${params.value}`;
        }
        return `${params.name}<br/>访客数: 0`;
      }
    },
    visualMap: {
      min: 0,
      max: maxValue,
      left: "left",
      top: "bottom",
      text: ["高", "低"],
      calculable: true,
      inRange: {
        color: ["#e0f3f8", "#abd9e9", "#74add1", "#4575b4", "#313695"]
      }
    },
    series: [
      {
        name: "访客数",
        type: "map",
        map: "china",
        roam: true,
        scaleLimit: {
          min: 0.5,
          max: 3
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            color: "#fff"
          },
          itemStyle: {
            areaColor: "#2a8fd8"
          }
        },
        data: convertedData
      }
    ]
  };

  chartInstance.setOption(option);
};

// 过滤出国外访客数据（非中国省份的数据）
const overseasData = computed(() => {
  if (!props.mapData) return [];
  const chineseProvinces = Object.keys(provinceNameMap);
  return props.mapData
    .filter(item => !chineseProvinces.includes(item.name))
    .sort((a, b) => b.value - a.value);
});

// 国家名称到旗帜emoji的映射
const countryFlagMap: Record<string, string> = {
  美国: "🇺🇸",
  日本: "🇯🇵",
  韩国: "🇰🇷",
  英国: "🇬🇧",
  德国: "🇩🇪",
  法国: "🇫🇷",
  俄罗斯: "🇷🇺",
  加拿大: "🇨🇦",
  澳大利亚: "🇦🇺",
  新加坡: "🇸🇬",
  马来西亚: "🇲🇾",
  泰国: "🇹🇭",
  越南: "🇻🇳",
  印度: "🇮🇳",
  印度尼西亚: "🇮🇩",
  菲律宾: "🇵🇭",
  巴西: "🇧🇷",
  墨西哥: "🇲🇽",
  意大利: "🇮🇹",
  西班牙: "🇪🇸",
  荷兰: "🇳🇱",
  瑞士: "🇨🇭",
  瑞典: "🇸🇪",
  挪威: "🇳🇴",
  丹麦: "🇩🇰",
  芬兰: "🇫🇮",
  波兰: "🇵🇱",
  土耳其: "🇹🇷",
  阿联酋: "🇦🇪",
  沙特阿拉伯: "🇸🇦",
  以色列: "🇮🇱",
  南非: "🇿🇦",
  埃及: "🇪🇬",
  新西兰: "🇳🇿",
  爱尔兰: "🇮🇪",
  奥地利: "🇦🇹",
  比利时: "🇧🇪",
  捷克: "🇨🇿",
  希腊: "🇬🇷",
  葡萄牙: "🇵🇹",
  阿根廷: "🇦🇷",
  智利: "🇨🇱",
  哥伦比亚: "🇨🇴",
  乌克兰: "🇺🇦"
};

const getCountryFlag = (name: string) => {
  return countryFlagMap[name] || "🌐";
};

// 海外访客总数
const overseasTotal = computed(() => {
  return overseasData.value.reduce((sum, item) => sum + item.value, 0);
});

// 海外图表引用和实例
const overseasChartRef = ref<HTMLElement | null>(null);
let overseasChartInstance: echarts.ECharts | null = null;

// 初始化海外访客图表
const initOverseasChart = () => {
  if (!overseasChartRef.value || overseasData.value.length === 0) return;

  overseasChartInstance = echarts.init(overseasChartRef.value);
  updateOverseasChart();
};

// 更新海外访客图表
const updateOverseasChart = () => {
  if (!overseasChartInstance) return;

  const data = overseasData.value.slice(0, 10); // 最多显示10个
  const countries = data.map(
    item => `${getCountryFlag(item.name)} ${item.name}`
  );
  const values = data.map(item => item.value);

  const option = {
    grid: {
      left: "3%",
      right: "10%",
      top: "5%",
      bottom: "5%",
      containLabel: true
    },
    xAxis: {
      type: "value",
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { show: false },
      splitLine: { show: false }
    },
    yAxis: {
      type: "category",
      data: countries.reverse(),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: {
        fontSize: 12,
        color: "#666"
      }
    },
    series: [
      {
        type: "bar",
        data: values.reverse(),
        barWidth: 16,
        itemStyle: {
          borderRadius: [0, 8, 8, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: "#36a3f7" },
            { offset: 1, color: "#6dd5ed" }
          ])
        },
        label: {
          show: true,
          position: "right",
          formatter: "{c}",
          fontSize: 12,
          color: "#36a3f7",
          fontWeight: "bold"
        }
      }
    ]
  };

  overseasChartInstance.setOption(option);
};

watch(
  () => props.mapData,
  () => {
    updateChart();
    if (overseasChartInstance) {
      updateOverseasChart();
    } else {
      setTimeout(initOverseasChart, 100);
    }
  },
  { deep: true }
);

watch(
  () => overseasData.value,
  () => {
    if (overseasData.value.length > 0 && !overseasChartInstance) {
      setTimeout(initOverseasChart, 100);
    }
  },
  { immediate: true }
);

onMounted(() => {
  initChart();
  setTimeout(initOverseasChart, 200);
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
  chartInstance?.dispose();
  overseasChartInstance?.dispose();
});
</script>

<style scoped>
.map-container {
  width: 100%;
}

.map-chart {
  width: 100%;
  height: 400px;
}

.overseas-section {
  margin-top: 16px;
  padding: 16px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  border-radius: 12px;
  border: 1px solid #e0e6ed;
}

.overseas-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.overseas-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.overseas-total {
  font-size: 13px;
  color: #36a3f7;
  font-weight: 500;
}

.overseas-chart {
  width: 100%;
  height: 200px;
}
</style>
