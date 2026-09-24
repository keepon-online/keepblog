<template>
  <div
    v-if="visible"
    class="ai-floating-palette-overlay"
    @click.self="$emit('close')"
  >
    <div
      ref="barRef"
      class="ai-floating-bar"
      :style="barStyle"
      @click.stop
    >
      <div class="bar-header">
        <div class="header-left">
          <span class="ai-sparkle">✨</span>
          <span class="palette-title">AI 指令中心</span>
          <span class="context-badge" :class="{ 'has-selection': hasSelection }">
            {{ hasSelection ? `已选中 ${selectionLength} 字符` : "光标位置" }}
          </span>
        </div>
        <div class="header-right">
          <span class="kbd-badge">Esc 退出</span>
        </div>
      </div>

      <div class="bar-input-wrap">
        <el-input
          ref="inputRef"
          v-model="customPrompt"
          placeholder="给 AI 下达指令（如：润色这段话、改写为对比表格、解释代码、续写）..."
          clearable
          size="large"
          @keydown.enter="handleSubmit"
          @keydown.esc="$emit('close')"
        >
          <template #prefix>
            <el-icon class="search-icon"><Promotion /></el-icon>
          </template>
          <template #suffix>
            <el-button
              type="primary"
              size="small"
              :disabled="!customPrompt.trim() && !hasSelection"
              @click="handleSubmit"
            >
              执行 ↵
            </el-button>
          </template>
        </el-input>
      </div>

      <div class="bar-chips-row">
        <span class="chips-label">快捷动作：</span>
        <div class="chips-list">
          <span
            v-for="action in QUICK_ACTIONS"
            :key="action.label"
            class="action-chip"
            @click="handleQuickAction(action)"
          >
            {{ action.label }}
          </span>
        </div>
      </div>

      <div v-if="models && models.length > 1" class="bar-models-row">
        <span class="models-label">切换模型：</span>
        <div class="models-list">
          <span
            v-for="m in models"
            :key="m.name"
            class="model-chip"
            :class="{ active: currentModel === m.name }"
            @click="selectModel(m.name)"
          >
            {{ m.name === 'default' ? m.model : m.name }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from "vue";
import { Promotion } from "@element-plus/icons-vue";

export interface QuickAction {
  label: string;
  task?: string;
  instruction?: string;
}

const props = defineProps<{
  visible: boolean;
  hasSelection: boolean;
  selectionLength?: number;
  models?: Array<{ name: string; model: string }>;
  currentModel?: string;
  position?: { top: number; left: number };
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "update:currentModel", model: string): void;
  (
    e: "submit",
    payload: {
      prompt: string;
      task?: string;
      model?: string;
    }
  ): void;
}>();

const customPrompt = ref("");
const inputRef = ref<any>(null);

const QUICK_ACTIONS: QuickAction[] = [
  { label: "✨ 智能润色", task: "polish", instruction: "语言更优美洗练" },
  { label: "⚡ 精炼要点", task: "polish", instruction: "精炼内容并提炼核心要点" },
  { label: "✍️ 丰富扩写", task: "polish", instruction: "丰富细节与论述" },
  { label: "🔍 全文校对", task: "proofread" },
  { label: "📊 转为表格", task: "polish", instruction: "整理为格式严谨的 Markdown 表格" },
  { label: "💡 解释代码", task: "code-explain" },
  { label: "🌐 译为英文", task: "polish", instruction: "准确翻译为专业地道英文" }
];

const barStyle = computed(() => {
  // 居中或在光标附近展示
  return {};
});

watch(
  () => props.visible,
  val => {
    if (val) {
      customPrompt.value = "";
      nextTick(() => {
        inputRef.value?.focus();
      });
    }
  }
);

const selectModel = (model: string) => {
  emit("update:currentModel", model);
};

const handleQuickAction = (action: QuickAction) => {
  emit("submit", {
    prompt: action.instruction || action.label,
    task: action.task,
    model: props.currentModel
  });
  emit("close");
};

const handleSubmit = () => {
  const prompt = customPrompt.value.trim();
  if (!prompt && !props.hasSelection) return;

  emit("submit", {
    prompt: prompt || "润色优化",
    task: "polish",
    model: props.currentModel
  });
  emit("close");
};
</script>

<style scoped lang="scss">
.ai-floating-palette-overlay {
  position: fixed;
  inset: 0;
  z-index: 2300;
  background: rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(2px);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 15vh;
  animation: overlayFadeIn 0.15s ease-out;
}

.ai-floating-bar {
  width: 620px;
  background: var(--el-bg-color-overlay, #ffffff);
  border-radius: 12px;
  border: 1px solid var(--el-border-color-light);
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.2);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  animation: barScaleIn 0.18s cubic-bezier(0.16, 1, 0.3, 1);

  .bar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
      display: flex;
      align-items: center;
      gap: 8px;

      .ai-sparkle {
        font-size: 16px;
      }

      .palette-title {
        font-size: 13px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .context-badge {
        font-size: 11px;
        background: var(--el-fill-color-light);
        color: var(--el-text-color-secondary);
        padding: 2px 8px;
        border-radius: 4px;

        &.has-selection {
          background: var(--el-color-primary-light-9);
          color: var(--el-color-primary);
          font-weight: 500;
        }
      }
    }

    .header-right {
      .kbd-badge {
        font-size: 11px;
        color: var(--el-text-color-placeholder);
        background: var(--el-fill-color-darker);
        padding: 2px 6px;
        border-radius: 4px;
        font-family: ui-monospace, SFMono-Regular, monospace;
      }
    }
  }

  .bar-input-wrap {
    :deep(.el-input__wrapper) {
      box-shadow: 0 0 0 1px var(--el-border-color) inset;
      border-radius: 8px;
      padding-left: 12px;

      &.is-focus {
        box-shadow: 0 0 0 2px var(--el-color-primary) inset;
      }
    }

    .search-icon {
      font-size: 16px;
      color: var(--el-color-primary);
    }
  }

  .bar-chips-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;

    .chips-label {
      color: var(--el-text-color-secondary);
      flex-shrink: 0;
      font-size: 11px;
    }

    .chips-list {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;

      .action-chip {
        padding: 3px 8px;
        border-radius: 4px;
        background: var(--el-fill-color-light);
        color: var(--el-text-color-regular);
        cursor: pointer;
        transition: all 0.15s ease;
        border: 1px solid var(--el-border-color-extra-light);

        &:hover {
          background: var(--el-color-primary-light-9);
          color: var(--el-color-primary);
          border-color: var(--el-color-primary-light-5);
        }
      }
    }
  }

  .bar-models-row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-top: 6px;
    border-top: 1px solid var(--el-border-color-extra-light);
    font-size: 11px;

    .models-label {
      color: var(--el-text-color-placeholder);
      flex-shrink: 0;
    }

    .models-list {
      display: flex;
      gap: 6px;

      .model-chip {
        padding: 2px 6px;
        border-radius: 3px;
        background: var(--el-fill-color);
        color: var(--el-text-color-secondary);
        cursor: pointer;

        &.active {
          background: var(--el-color-primary-light-8);
          color: var(--el-color-primary);
          font-weight: 600;
        }
      }
    }
  }
}

@keyframes overlayFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes barScaleIn {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>
