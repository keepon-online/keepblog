import { http } from "@/utils/http";
import { getToken, formatToken } from "@/utils/auth";

type Result = {
  code: number;
  message?: string;
  payload?: any;
};

export type AIStatus = {
  enabled: boolean;
  model: string;
  todayUsed: number;
  dailyQuota: number;
};

export type AIEditRequest = {
  task:
    | "polish"
    | "continue"
    | "title"
    | "summary"
    | "refine"
    | "tags"
    | "proofread";
  mode?: "polish" | "expand" | "shorten" | "translate";
  selection?: string;
  before?: string;
  after?: string;
  title?: string;
  series?: string;
  digest?: string;
  /** refine 任务：上一版结果 */
  previous?: string;
  /** refine 任务：追加修饰要求 */
  instruction?: string;
};

export type AIStreamHandlers = {
  onDelta?: (text: string) => void;
  /** 推理模型正文前的思考增量（用于"思考中"展示） */
  onReasoning?: (text: string) => void;
  onDone?: (usage?: any) => void;
  onError?: (message: string) => void;
};

/** 探测 AI 是否可用（未配置 apiKey 时前端隐藏所有 AI 入口） */
export const getAIStatus = () => {
  return http.request<Result>("get", "/api/v1/ai/status");
};

const resolveBaseURL = () => {
  const base =
    (import.meta.env.VITE_BASE_URL as string) ||
    (typeof window !== "undefined" ? window.location.origin : "");
  return base.replace(/\/$/, "");
};

/**
 * 取可用 accessToken：过期时先经 /api/refreshToken 无感续期（与 axios
 * 拦截器同款语义）。原生 fetch 不走拦截器，token 过期会以 401 失败，
 * 这里在发流式请求前对齐一次。
 */
async function ensureAccessToken(): Promise<string> {
  const data: any = getToken();
  if (!data?.accessToken) return "";
  const expired = parseInt(data.expires) - Date.now() <= 0;
  if (!expired) return data.accessToken;
  // 裸 fetch 刷新：不走 axios 拦截器——后者在刷新失败（401）时会触发
  // 全局 logOut() 把用户踢回登录页，对"点一下 AI 按钮"来说副作用过重。
  // 这里失败就静默返回旧 token，由后续请求自行暴露真实错误。
  try {
    const resp = await fetch(resolveBaseURL() + "/api/refreshToken", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refreshToken: data.refreshToken })
    });
    if (resp.ok) {
      const res: any = await resp.json();
      if (res?.code === 200 && res?.payload?.accessToken) {
        return res.payload.accessToken;
      }
    }
    return getToken()?.accessToken ?? data.accessToken;
  } catch {
    return data.accessToken;
  }
}

/**
 * SSE 流式 AI 编辑。
 *
 * 不走 axios 封装：纯 http 客户端会缓冲整个响应体，流式语义必须用原生
 * fetch + ReadableStream 手工解析 `data:` 行。鉴权头与 axios 侧保持一致。
 *
 * @returns abort：中止生成（同时中断上游请求）
 */
export const streamAIEdit = (
  body: AIEditRequest,
  handlers: AIStreamHandlers
): (() => void) => {
  const controller = new AbortController();
  const base = resolveBaseURL();

  (async () => {
    try {
      const accessToken = await ensureAccessToken();
      const resp = await fetch(base + "/api/v1/ai/edit", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(accessToken
            ? { Authorization: formatToken(accessToken) }
            : {})
        },
        body: JSON.stringify(body),
        signal: controller.signal
      });

      if (!resp.ok || !resp.body) {
        // 非流式错误信封（参数错/未配置/超配额），读完整体再回调
        const text = await resp.text();
        let message = `请求失败（${resp.status}）`;
        try {
          const data = JSON.parse(text);
          message = data.message || message;
        } catch {
          /* 非 JSON 保持默认消息 */
        }
        handlers.onError?.(message);
        return;
      }

      const reader = resp.body.getReader();
      const decoder = new TextDecoder("utf-8");
      let buffer = "";

      for (;;) {
        const { done, value } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });

        // SSE 事件以空行分隔；按块切出完整事件再解析，半包留在 buffer
        const blocks = buffer.split("\n\n");
        buffer = blocks.pop() ?? "";
        for (const block of blocks) {
          for (const line of block.split("\n")) {
            if (!line.startsWith("data:")) continue;
            let evt: any;
            try {
              evt = JSON.parse(line.slice(5).trim());
            } catch {
              continue;
            }
            if (evt.type === "delta" && typeof evt.content === "string") {
              handlers.onDelta?.(evt.content);
            } else if (evt.type === "reasoning" && typeof evt.content === "string") {
              handlers.onReasoning?.(evt.content);
            } else if (evt.type === "done") {
              handlers.onDone?.(evt.usage);
            } else if (evt.type === "error") {
              handlers.onError?.(evt.message || "生成失败");
            }
          }
        }
      }
      handlers.onDone?.();
    } catch (e: any) {
      if (e?.name === "AbortError") return; // 主动停止不算错误
      const detail = e?.message ? `${e.message}` : "网络错误（无法连接服务）";
      handlers.onError?.(detail);
    }
  })();

  return () => controller.abort();
};
