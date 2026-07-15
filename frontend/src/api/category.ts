import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 获取分类管理列表 */
export const getCategoryList = () => {
  return http.request<Result>("get", "/api/v1/site/category/list");
};
/**更新分类**/
export const updateCategory = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/category/update", {
    data
  });
};
/**更新分类**/
export const addCategory = (data?: object) => {
  return http.request<Result>("post", "/api/v1/site/category/save", {
    data
  });
};

/**删除分类**/
export const deleteCategory = (categoryId?: number) => {
  return http.request<Result>(
    "delete",
    "/api/v1/site/category/delete/" + categoryId
  );
};

/**更新分类状态**/
export const changeCategoryState = (data?: object) => {
  return http.request<Result>("put", "/api/v1/site/category/update-state", {
    data
  });
};
/**获取分类**/
export const getCategory = (categoryId?: number) => {
  return http.request<Result>(
    "get",
    "/api/v1/site/category/detail/" + categoryId
  );
};
