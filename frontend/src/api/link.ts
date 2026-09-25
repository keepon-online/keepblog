import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

export const getLinkList = () => {
  return http.request<Result>("get", "/api/v1/site/link/list");
};

export const updateLink = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/link/update", {
    data
  });
};

export const addLink = (data?: object) => {
  return http.request<Result>("post", "/api/v1/site/link/save", {
    data
  });
};

export const deleteLink = (id?: number) => {
  return http.request<Result>("delete", "/api/v1/site/link/delete/" + id);
};

export const changeLinkState = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/link/update-state", {
    data
  });
};

export const getLink = (id?: number) => {
  return http.request<Result>("get", "/api/v1/site/link/detail/" + id);
};

export const checkLink = (url: string) => {
  return http.request<{
    code: number;
    message: string;
    payload: { status: number; msg: string };
  }>("get", "/api/v1/site/link/check", { params: { url } });
};

