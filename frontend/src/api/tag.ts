import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 获取专栏管理列表 */
export const getTagList = () => {
  return http.request<Result>("get", "/api/v1/site/tags/list");
};
