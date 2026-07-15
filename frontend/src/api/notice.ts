import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: any;
};

/** 通知项 */
export interface NoticeItem {
  id: number;
  avatar: string;
  title: string;
  description: string;
  datetime: string;
  type: string;
  extra?: string;
  status?: "primary" | "success" | "warning" | "info" | "danger";
}

/** 通知分类 */
export interface NoticeTab {
  key: string;
  name: string;
  list: NoticeItem[];
  emptyText: string;
}

/** 通知列表响应 */
export interface NoticeListResponse {
  notices: NoticeTab[];
  unreadCount: number;
}

/** 获取通知列表 */
export const getNoticeList = () => {
  return http.request<Result>("get", "/api/v1/notice/list");
};

/** 标记通知为已读 */
export const markNoticeRead = (id: number) => {
  return http.request<Result>("put", `/api/v1/notice/read/${id}`);
};

/** 标记所有通知为已读 */
export const markAllNoticeRead = (type?: number) => {
  return http.request<Result>("put", "/api/v1/notice/read-all", {
    data: type !== undefined ? { type } : {}
  });
};

/** 删除通知 */
export const deleteNotice = (id: number) => {
  return http.request<Result>("delete", `/api/v1/notice/delete/${id}`);
};
