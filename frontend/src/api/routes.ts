import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload: Array<any>;
};

export const getAsyncRoutes = () => {
  return http.request<Result>("get", "/api/getAsyncRoutes");
};
