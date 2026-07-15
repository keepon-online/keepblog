import { http } from "@/utils/http";

type Result = {
  code: number;
  message: string;
  payload?: string;
};

export const upload = (data?: object) => {
  return http.request<Result>(
    "post",
    "/api/v1/upload/images",
    {
      data
    },
    {
      headers: {
        "Content-Type": "multipart/form-data"
      }
    }
  );
};
