/**
 * ECharts 按需引入配置
 * 按需引入可显著减少打包体积
 */

// 核心模块
import * as echarts from "echarts/core";

// 图表类型
import {
  BarChart,
  LineChart,
  PieChart,
  GaugeChart,
  MapChart,
  ScatterChart
} from "echarts/charts";

// 组件
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  LegendComponent,
  DataZoomComponent,
  ToolboxComponent,
  VisualMapComponent,
  GeoComponent,
  MarkPointComponent,
  MarkLineComponent
} from "echarts/components";

// 标签自动布局、全局过渡动画等特性
import { LabelLayout, UniversalTransition } from "echarts/features";

// 引入 Canvas 渲染器（也可使用 SVGRenderer）
import { CanvasRenderer } from "echarts/renderers";

// 注册必要的组件
echarts.use([
  // 图表
  BarChart,
  LineChart,
  PieChart,
  GaugeChart,
  MapChart,
  ScatterChart,
  // 组件
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  LegendComponent,
  DataZoomComponent,
  ToolboxComponent,
  VisualMapComponent,
  GeoComponent,
  MarkPointComponent,
  MarkLineComponent,
  // 特性
  LabelLayout,
  UniversalTransition,
  // 渲染器
  CanvasRenderer
]);

export default echarts;

// 导出类型
export type EChartsOption = echarts.ComposeOption<
  | import("echarts/charts").BarSeriesOption
  | import("echarts/charts").LineSeriesOption
  | import("echarts/charts").PieSeriesOption
  | import("echarts/charts").GaugeSeriesOption
  | import("echarts/charts").MapSeriesOption
  | import("echarts/charts").ScatterSeriesOption
  | import("echarts/components").TitleComponentOption
  | import("echarts/components").TooltipComponentOption
  | import("echarts/components").GridComponentOption
  | import("echarts/components").DatasetComponentOption
  | import("echarts/components").LegendComponentOption
  | import("echarts/components").DataZoomComponentOption
  | import("echarts/components").ToolboxComponentOption
  | import("echarts/components").VisualMapComponentOption
  | import("echarts/components").GeoComponentOption
>;
