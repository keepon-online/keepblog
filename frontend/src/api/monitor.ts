import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 获取系统监控完整信息 */
export const getServe = () => {
  return http.request<Result>("get", "/api/monitor/server");
};

/** 获取实时监控统计（轻量） */
export const getRealtime = () => {
  return http.request<Result>("get", "/api/monitor/realtime");
};
