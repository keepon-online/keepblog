import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 获取文章管理列表 */
export const getPostList = (params?: object) => {
  return http.request<Result>("get", "/api/v1/site/post/list", {
    params
  });
};

/** 保存 */
export const savePost = (data?: object) => {
  return http.request<Result>("post", "/api/v1/site/post/save", {
    data
  });
};
/** 更新 */
export const updatePost = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/post/update", {
    data
  });
}; /** 更新 */
export const updatePostPublish = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/post/update-publish", {
    data
  });
};
export const updatePostTop = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/post/update-top", {
    data
  });
};
/** 删除**/
export const deletePost = (postId?: string) => {
  return http.request<Result>("delete", "/api/v1/site/post/delete/" + postId);
};
export const updateCover = (postId?: string) => {
  return http.request<Result>(
    "put",
    "/api/v1/site/post/update-cover/" + postId
  );
};

export const updateAllCover = () => {
  return http.request<Result>("put", "/api/v1/site/post/update-cover");
};

/** 获取文章 */
export const getPost = (postId?: string | string[]) => {
  return http.request<Result>("get", "/api/v1/site/post/detail/" + postId);
};
