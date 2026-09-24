import { ref, computed } from "vue";
import {
  getAIStatus,
  streamAIEdit,
  type AIEditRequest,
  type AIStatus
} from "@/api/ai";
import { message } from "@/utils/message";

export interface AIHistoryEntry {
  content: string;
  label: string;
}

export function useAiStreaming() {
  const aiStatus = ref<AIStatus | null>(null);
  const aiEnabled = computed(() => !!aiStatus.value?.enabled);

  const aiDrawerVisible = ref(false);
  const aiStreaming = ref(false);
  const aiResultText = ref("");
  const aiReasoningText = ref("");
  const aiError = ref("");
  const aiMeta = ref<{ elapsed: string; tokens: number } | null>(null);

  const aiTaskLabel = ref("AI 助手");
  const aiApplyType = ref<"replace-selection" | "insert-cursor">("replace-selection");
  const aiOriginalSelection = ref("");
  const aiSelRange = ref<{ from: number; to: number } | null>(null);

  let aiResultBuf = "";
  let aiReasoningBuf = "";
  let aiFlushTimer: ReturnType<typeof setTimeout> | null = null;
  let aiAbort: (() => void) | null = null;
  let aiLastRequest: AIEditRequest | null = null;
  let aiStartedAt = 0;

  const aiHistory = ref<AIHistoryEntry[]>([]);
  const aiHistoryIdx = ref(-1);

  // 配额提示
  const aiQuotaText = computed(() => {
    const s = aiStatus.value;
    if (!s || !s.enabled) return "";
    if (s.dailyQuota <= 0) return "无限配额";
    const remaining = Math.max(0, s.dailyQuota - (s.todayUsed || 0));
    return `今日剩余 ${remaining}/${s.dailyQuota} 次`;
  });

  const aiQuotaExhausted = computed(() => {
    const s = aiStatus.value;
    if (!s || !s.enabled) return true;
    return s.dailyQuota > 0 && s.todayUsed >= s.dailyQuota;
  });

  const aiModelOptions = computed(
    () => aiStatus.value?.models ?? []
  );

  const aiModelChoice = ref<string>(
    localStorage.getItem("post_editor_ai_model") || "default"
  );

  const setAIModel = (name: string) => {
    aiModelChoice.value = name;
    localStorage.setItem("post_editor_ai_model", name);
  };

  const refreshAIStatus = async () => {
    try {
      const res = await getAIStatus();
      if (res.code === 200) {
        aiStatus.value = res.payload ?? null;
      }
    } catch {
      // 接口失败降级为未配置
    }
  };

  const flushAI = () => {
    aiResultText.value = aiResultBuf;
    aiReasoningText.value = aiReasoningBuf;
  };

  const scheduleFlush = () => {
    if (aiFlushTimer) return;
    aiFlushTimer = setTimeout(() => {
      aiFlushTimer = null;
      flushAI();
    }, 50);
  };

  const pushHistory = (content: string, label = "生成") => {
    if (!content) return;
    if (aiHistoryIdx.value < aiHistory.value.length - 1) {
      aiHistory.value = aiHistory.value.slice(0, aiHistoryIdx.value + 1);
    }
    aiHistory.value.push({ content, label });
    aiHistoryIdx.value = aiHistory.value.length - 1;
  };

  const aiNavHistory = (delta: number) => {
    const next = aiHistoryIdx.value + delta;
    if (next < 0 || next >= aiHistory.value.length) return;
    aiHistoryIdx.value = next;
    aiResultText.value = aiHistory.value[next].content;
  };

  const selectHistoryIndex = (idx: number) => {
    if (idx < 0 || idx >= aiHistory.value.length) return;
    aiHistoryIdx.value = idx;
    aiResultText.value = aiHistory.value[idx].content;
  };

  const stopAI = () => {
    if (aiAbort) {
      aiAbort();
      aiAbort = null;
    }
    if (aiFlushTimer) {
      clearTimeout(aiFlushTimer);
      aiFlushTimer = null;
    }
    flushAI();
    aiStreaming.value = false;
  };

  const startAiStream = (
    req: AIEditRequest,
    options?: {
      label?: string;
      applyType?: "replace-selection" | "insert-cursor";
      selectionText?: string;
      selectionRange?: { from: number; to: number };
      onDelta?: (text: string) => void;
      onDone?: (fullText: string) => void;
    }
  ) => {
    stopAI();

    if (aiQuotaExhausted.value) {
      message("今日 AI 调用配额已用完，明天再来吧", { type: "warning" });
      return;
    }

    aiLastRequest = { ...req, model: aiModelChoice.value === "default" ? undefined : aiModelChoice.value };
    aiTaskLabel.value = options?.label || "AI 助手";
    aiApplyType.value = options?.applyType || "replace-selection";
    aiOriginalSelection.value = options?.selectionText || "";
    aiSelRange.value = options?.selectionRange || null;

    aiDrawerVisible.value = true;
    aiStreaming.value = true;
    aiError.value = "";
    aiMeta.value = null;
    aiResultBuf = "";
    aiReasoningBuf = "";
    aiResultText.value = "";
    aiReasoningText.value = "";
    aiStartedAt = Date.now();

    aiAbort = streamAIEdit(aiLastRequest, {
      onDelta: delta => {
        aiResultBuf += delta;
        scheduleFlush();
        options?.onDelta?.(delta);
      },
      onReasoning: r => {
        aiReasoningBuf += r;
        scheduleFlush();
      },
      onDone: () => {
        stopAI();
        const elapsed = `${((Date.now() - aiStartedAt) / 1000).toFixed(1)}s`;
        aiMeta.value = { elapsed, tokens: aiResultBuf.length };
        pushHistory(aiResultText.value, options?.label || "生成");
        refreshAIStatus();
        options?.onDone?.(aiResultText.value);
      },
      onError: msg => {
        stopAI();
        aiError.value = msg || "AI 生成中断或失败";
        message(aiError.value, { type: "error" });
      }
    });
  };

  const regenAI = () => {
    if (!aiLastRequest) return;
    startAiStream(aiLastRequest, {
      label: "重新生成",
      applyType: aiApplyType.value,
      selectionText: aiOriginalSelection.value,
      selectionRange: aiSelRange.value || undefined
    });
  };

  const refineAI = (instruction: string) => {
    if (!aiLastRequest || !aiResultText.value) return;
    const req: AIEditRequest = {
      task: "refine",
      previous: aiResultText.value,
      instruction,
      model: aiModelChoice.value === "default" ? undefined : aiModelChoice.value
    };
    startAiStream(req, {
      label: instruction.slice(0, 8),
      applyType: aiApplyType.value,
      selectionText: aiOriginalSelection.value,
      selectionRange: aiSelRange.value || undefined
    });
  };

  return {
    aiStatus,
    aiEnabled,
    aiQuotaText,
    aiQuotaExhausted,
    aiModelOptions,
    aiModelChoice,
    setAIModel,
    refreshAIStatus,
    aiDrawerVisible,
    aiStreaming,
    aiResultText,
    aiReasoningText,
    aiError,
    aiMeta,
    aiTaskLabel,
    aiApplyType,
    aiOriginalSelection,
    aiSelRange,
    aiLastRequest,
    aiHistory,
    aiHistoryIdx,
    pushHistory,
    aiNavHistory,
    selectHistoryIndex,
    startAiStream,
    stopAI,
    regenAI,
    refineAI
  };
}
