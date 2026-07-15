import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 用户信息响应类型 */
export interface UserInfoResponse {
  userId: number;
  username: string;
  nickName: string;
  email: string;
  phonenumber: string;
  sex: number;
  avatar: string;
  role: string;
  loginIp: string;
  loginTime: string;
  loginCount: number;
  registerTime: string;
}

/** 更新资料请求类型 */
export interface UpdateProfileRequest {
  nickName?: string;
  email?: string;
  phonenumber?: string;
  sex?: number;
  avatar?: string;
}

/** 获取用户信息 */
export const getUserInfo = () => {
  return http.request<Result>("get", "/api/getInfo");
};

/** 更新用户资料 */
export const updateProfile = (data: UpdateProfileRequest) => {
  return http.request<Result>("put", "/api/updateProfile", { data });
};

/** 修改密码 */
export const changePassword = (data?: object) => {
  return http.request<Result>("post", "/api/change-password", {
    data
  });
};
