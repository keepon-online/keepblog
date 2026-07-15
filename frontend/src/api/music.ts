import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

// 获取音乐列表
export const getMusicList = () => {
  return http.request<Result>("get", "/api/v1/site/music/list");
};

// 获取音乐详情
export const getMusic = (id?: number) => {
  return http.request<Result>("get", "/api/v1/site/music/detail/" + id);
};

// 添加音乐
export const addMusic = (data?: object) => {
  return http.request<Result>("post", "/api/v1/site/music/save", {
    data
  });
};

// 更新音乐
export const updateMusic = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/music/update", {
    data
  });
};

// 删除音乐
export const deleteMusic = (id?: number) => {
  return http.request<Result>("delete", "/api/v1/site/music/delete/" + id);
};

// 更新音乐状态
export const changeMusicState = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/music/update-state", {
    data
  });
};
