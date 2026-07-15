import { http } from "@/utils/http";

type DashboardData = {
  panel: {
    postTotal: number;
    categoryTotal: number;
    tagTotal: number;
    visit: number;
    totalWords: number;
    totalReadCount: number;
    todayVisit: number;
    totalMusic: number;
  };
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
};

/** 获取仪表板首页所有数据 */
export const getDashboardData = () => {
  return http.request<{
    code: number;
    message: string;
    payload: DashboardData;
  }>("get", "/api/v1/site/dashboard/data");
};
