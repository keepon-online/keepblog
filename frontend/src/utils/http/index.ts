import Axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type CustomParamsSerializer
} from "axios";
import type {
  PureHttpError,
  RequestMethods,
  PureHttpResponse,
  PureHttpRequestConfig
} from "./types.d";
import { stringify } from "qs";
import { getToken, formatToken } from "@/utils/auth";
import { useUserStoreHook } from "@/store/modules/user";
// 相关配置请参考：www.axios-js.com/zh-cn/docs/#axios-request-config-1
let Base_Url = "";

async function DynamicBaseUrl() {
  if (import.meta.env.PROD) {
    // 从后端请求获取动态地址
    // const response = await Axios.post(
    //   window.location.protocol + "//" + window.location.host
    // );
    // 设置动态地址作为baseURL
    Base_Url = window.location.protocol + "//" + window.location.host;
  }
  if (import.meta.env.DEV) {
    Base_Url = import.meta.env.VITE_BASE_URL;
  }
  return Base_Url;
}

const defaultConfig: AxiosRequestConfig = {
  baseURL: Base_Url,
  // 请求超时时间
  timeout: 10000,
  headers: {
    Accept: "application/json, text/plain, */*",
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest"
  },
  // 数组格式参数序列化（https://github.com/axios/axios/issues/5142）
  paramsSerializer: {
    serialize: stringify as unknown as CustomParamsSerializer
  }
};

class PureHttp {
  constructor() {
    this.httpInterceptorsRequest();
    this.httpInterceptorsResponse();
  }

  /** `token`过期后，暂存待执行的请求 */
  private static requests = [];

  /** 防止重复刷新`token` */
  private static isRefreshing = false;

  /** 初始化配置对象 */
  private static initConfig: PureHttpRequestConfig = {};

  /** 保存当前`Axios`实例对象 */
  private static axiosInstance: AxiosInstance = Axios.create(defaultConfig);

  /** 重连原始请求 */
  private static retryOriginalRequest(config: PureHttpRequestConfig) {
    return new Promise(resolve => {
      PureHttp.requests.push((token: string) => {
        config.headers["Authorization"] = formatToken(token);
        resolve(config);
      });
    });
  }

  /** 请求拦截 */
  private httpInterceptorsRequest(): void {
    PureHttp.axiosInstance.interceptors.request.use(
      async (config: PureHttpRequestConfig): Promise<any> => {
        if (Base_Url == "") {
          config.baseURL = await DynamicBaseUrl();
        } else {
          config.baseURL = Base_Url;
        }
        // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
        if (typeof config.beforeRequestCallback === "function") {
          config.beforeRequestCallback(config);
          return config;
        }
        if (PureHttp.initConfig.beforeRequestCallback) {
          PureHttp.initConfig.beforeRequestCallback(config);
          return config;
        }
        /** 请求白名单，放置一些不需要`token`的接口（通过设置请求白名单，防止`token`过期后再请求造成的死循环问题） */
        const whiteList = ["/refreshToken", "/login"];
        return whiteList.some(url => config.url.endsWith(url))
          ? config
          : new Promise(resolve => {
            const data = getToken();
            if (data) {
              const now = new Date().getTime();
              const expired = parseInt(data.expires) - now <= 0;
              if (expired) {
                if (!PureHttp.isRefreshing) {
                  PureHttp.isRefreshing = true;
                  // token过期刷新
                  useUserStoreHook()
                    .handRefreshToken({ refreshToken: data.refreshToken })
                    .then(res => {
                      const token = res.payload.accessToken;
                      config.headers["Authorization"] = formatToken(token);
                      PureHttp.requests.forEach(cb => cb(token));
                      PureHttp.requests = [];
                    })
                    .finally(() => {
                      PureHttp.isRefreshing = false;
                    });
                }
                resolve(PureHttp.retryOriginalRequest(config));
              } else {
                config.headers["Authorization"] = formatToken(
                  data.accessToken
                );
                resolve(config);
              }
            } else {
              resolve(config);
            }
          });
      },
      error => {
        return Promise.reject(error);
      }
    );
  }

  /** 响应拦截 */
  private httpInterceptorsResponse(): void {
    const instance = PureHttp.axiosInstance;
    instance.interceptors.response.use(
      (response: PureHttpResponse) => {
        const $config = response.config;
        const data = response.data as { code?: number; message?: string };

        // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
        if (typeof $config.beforeResponseCallback === "function") {
          $config.beforeResponseCallback(response);
          return response.data;
        }
        if (PureHttp.initConfig.beforeResponseCallback) {
          PureHttp.initConfig.beforeResponseCallback(response);
          return response.data;
        }

        // 统一处理业务错误码
        if (data.code !== undefined && data.code !== 200) {
          // 401 未授权，跳转到登录页
          if (data.code === 401) {
            useUserStoreHook().logOut();
          }
          // 其他业务错误不在这里处理，由具体业务代码处理
        }

        return response.data;
      },
      (error: PureHttpError) => {
        const $error = error;
        $error.isCancelRequest = Axios.isCancel($error);

        // 网络错误处理
        if (!$error.isCancelRequest) {
          const status = error.response?.status;
          let message = "网络请求失败";

          switch (status) {
            case 400:
              message = "请求参数错误";
              break;
            case 401:
              message = "登录已过期，请重新登录";
              useUserStoreHook().logOut();
              break;
            case 403:
              message = "没有访问权限";
              break;
            case 404:
              message = "请求的资源不存在";
              break;
            case 500:
              message = "服务器内部错误";
              break;
            case 502:
              message = "网关错误";
              break;
            case 503:
              message = "服务暂不可用";
              break;
            case 504:
              message = "网关超时";
              break;
            default:
              if (!navigator.onLine) {
                message = "网络连接已断开";
              }
          }

          // 将错误信息附加到 error 对象
          $error.message = message;
        }

        return Promise.reject($error);
      }
    );
  }

  /** 通用请求工具函数 */
  public request<T>(
    method: RequestMethods,
    url: string,
    param?: AxiosRequestConfig,
    axiosConfig?: PureHttpRequestConfig
  ): Promise<T> {
    const config = {
      method,
      url,
      ...param,
      ...axiosConfig
    } as PureHttpRequestConfig;

    // 单独处理自定义请求/响应回调
    return new Promise((resolve, reject) => {
      PureHttp.axiosInstance
        .request(config)
        .then((response: undefined) => {
          resolve(response);
        })
        .catch(error => {
          reject(error);
        });
    });
  }

  /** 单独抽离的`post`工具函数 */
  public post<T, P>(
    url: string,
    params?: AxiosRequestConfig<P>,
    config?: PureHttpRequestConfig
  ): Promise<T> {
    return this.request<T>("post", url, params, config);
  }

  /** 单独抽离的`get`工具函数 */
  public get<T, P>(
    url: string,
    params?: AxiosRequestConfig<P>,
    config?: PureHttpRequestConfig
  ): Promise<T> {
    return this.request<T>("get", url, params, config);
  }
}

export const http = new PureHttp();
