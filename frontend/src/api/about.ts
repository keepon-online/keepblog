import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

export const updateAbout = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/about/update", {
    data
  });
};

export const getAbout = () => {
  return http.request<Result>("get", "/api/v1/site/about/detail");
};
