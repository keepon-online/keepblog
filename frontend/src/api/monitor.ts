import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 获取分类管理列表 */
export const getServe = () => {
  return http.request<Result>("get", "/api/monitor/server");
};
