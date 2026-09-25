import { http } from "@/utils/http";

export type DashboardPanel = {
  postTotal: number;
  categoryTotal: number;
  tagTotal: number;
  visit: number;
  totalWords: number;
  totalReadCount: number;
  todayVisit: number;
  yesterdayVisit?: number;
  visitGrowth?: number;
  weekPostTotal?: number;
  totalMusic: number;
};

export type TopPostItem = {
  id: number;
  title: string;
  readCount: number;
  categoryName: string;
  createTime: number;
};

export type DashboardData = {
  panel: DashboardPanel;
  line: Array<{
    name: string;
    pv: number;
    uv: number;
  }>;
  pie: Array<{
    name: string;
    value: number;
  }>;
  bar: Array<{
    name: string;
    value: number;
  }>;
  map: Array<{
    name: string;
    value: number;
  }>;
  topPosts?: Array<TopPostItem>;
  days?: number;
};

/** 获取仪表板首页所有数据 */
export const getDashboardData = (params?: { days?: number }) => {
  return http.request<{
    code: number;
    message: string;
    payload: DashboardData;
  }>("get", "/api/v1/site/dashboard/data", { params });
};

