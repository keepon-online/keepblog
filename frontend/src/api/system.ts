import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

export const getLoginLog = (params?: object) => {
  return http.request<Result>("get", "/api/v1/site/logs/logon", {
    params
  });
};
export const getAccessLog = (params?: object) => {
  return http.request<Result>("get", "/api/v1/site/logs/access", {
    params
  });
};
