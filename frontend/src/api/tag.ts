import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 获取标签列表（前台 tag cloud，供文章编辑器候选标签使用） */
export const getTagList = () => {
  return http.request<Result>("get", "/api/v1/site/tags/list");
};

/** 获取标签管理列表（后台管理，返回全部标签） */
export const getTagManageList = () => {
  return http.request<Result>("get", "/api/v1/site/tags/manage-list");
};

/** 获取单个标签 */
export const getTag = (tagId?: number) => {
  return http.request<Result>("get", "/api/v1/site/tags/detail/" + tagId);
};

/** 新增标签 */
export const addTag = (data?: object) => {
  return http.request<Result>("post", "/api/v1/site/tags/save", { data });
};

/** 更新标签 */
export const updateTag = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/tags/update", { data });
};

/** 删除标签 */
export const deleteTag = (tagId?: number) => {
  return http.request<Result>("delete", "/api/v1/site/tags/delete/" + tagId);
};
