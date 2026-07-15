import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

export const updateWebsite = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/web/update", {
    data
  });
};

export const getWebsite = () => {
  return http.request<Result>("get", "/api/v1/site/web/detail");
};
