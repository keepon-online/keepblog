<template>
  <el-drawer
    :model-value="visible"
    :size="lastRequest?.task === 'proofread' ? '560px' : '500px'"
    :modal="false"
    :lock-scroll="false"
    class="ai-docked-drawer"
    :close-on-click-modal="false"
    @close="$emit('close')"
    @keydown.ctrl.enter.prevent="handlePrimaryAction"
    @keydown.meta.enter.prevent="handlePrimaryAction"
  >
    <template #header>
      <div class="ai-drawer-header">
        <span class="ai-drawer-title">{{ taskLabel }}</span>
        <!-- 历史版本药丸导航 -->
        <div v-if="history.length > 1" class="ai-history-pills">
          <span
            v-for="(h, idx) in history"
            :key="idx"
            class="ai-history-pill"
            :class="{ active: historyIdx === idx }"
            :title="h.content"
            @click="$emit('select-history', idx)"
          >
            {{ h.label }}
          </span>
        </div>
      </div>
    </template>

    <div ref="bodyRef" class="ai-result-body" @scroll="onBodyScroll">
      <el-alert
        v-if="error"
        :title="error"
        type="error"
        show-icon
        :closable="false"
        class="ai-error-banner"
      />

      <!-- 任务为全文校对时的专属交互界面 -->
      <template v-if="lastRequest?.task === 'proofread'">
        <!-- 校对控制条 -->
        <div class="proofread-header-bar">
          <div class="proofread-stats">
            <span class="proofread-total-badge">
              发现 {{ proofreadItems.length }} 处问题
            </span>
            <span v-if="pendingProofreadCount > 0" class="text-orange-500 text-xs">
              ({{ pendingProofreadCount }} 待处理)
            </span>
            <span v-if="appliedProofreadCount > 0" class="text-green-600 text-xs">
              ({{ appliedProofreadCount }} 已采纳)
            </span>
          </div>

          <div class="proofread-header-actions">
            <el-button
              v-if="pendingProofreadCount > 0 && !streaming"
              size="small"
              type="primary"
              plain
              @click="$emit('apply-all-proofread')"
            >
              一键采纳全部 ({{ pendingProofreadCount }})
            </el-button>
            <el-radio-group v-model="proofreadViewMode" size="small">
              <el-radio-button label="cards">卡片视图</el-radio-button>
              <el-radio-button label="raw">原始输出</el-radio-button>
            </el-radio-group>
          </div>
        </div>

        <!-- 卡片视图下的筛选与列表 -->
        <div v-if="proofreadViewMode === 'cards'" class="proofread-cards-container">
          <div v-if="proofreadItems.length > 0" class="proofread-filter-row">
            <el-radio-group v-model="proofreadFilter" size="small">
              <el-radio-button label="all">全部 ({{ proofreadItems.length }})</el-radio-button>
              <el-radio-button label="pending">待处理 ({{ pendingProofreadCount }})</el-radio-button>
              <el-radio-button label="applied">已采纳 ({{ appliedProofreadCount }})</el-radio-button>
            </el-radio-group>
          </div>

          <div v-if="proofreadItems.length === 0 && !streaming" class="proofread-empty">
            <el-icon class="text-green-500 text-3xl mb-2"><CircleCheck /></el-icon>
            <div>未发现明显错别字或语病，全文表述良好！</div>
          </div>

          <div class="proofread-cards-list">
            <div
              v-for="item in filteredProofreadItems"
              :key="item.id"
              class="proofread-card"
              :class="{
                'is-applied': item.status === 'applied',
                'is-ignored': item.status === 'ignored',
                'is-not-found': item.status === 'not_found'
              }"
            >
              <div class="proofread-card-header">
                <el-tag
                  size="small"
                  :type="getProofreadTagType(item.type)"
                  effect="light"
                >
                  {{ item.type }}
                </el-tag>
                <span v-if="item.reason" class="proofread-reason" :title="item.reason">
                  {{ item.reason }}
                </span>
                <span class="proofread-status-indicator">
                  <span v-if="item.status === 'applied'" class="text-green-600 font-medium">✓ 已采纳</span>
                  <span v-else-if="item.status === 'ignored'" class="text-gray-400">已忽略</span>
                  <span v-else-if="item.status === 'not_found'" class="text-orange-400">未在正文中匹配</span>
                </span>
              </div>

              <div class="proofread-card-diff">
                <div class="diff-line">
                  <span class="diff-tag del">原</span>
                  <span class="diff-text del">{{ item.original }}</span>
                </div>
                <div class="diff-line">
                  <span class="diff-tag ins">改</span>
                  <span class="diff-text ins">{{ item.suggestion }}</span>
                </div>
              </div>

              <div class="proofread-card-footer">
                <el-button
                  link
                  size="small"
                  type="primary"
                  @click="$emit('locate-proofread', item)"
                >
                  <el-icon class="mr-2px"><Search /></el-icon>
                  定位原文
                </el-button>
                <template v-if="item.status === 'pending' || !item.status">
                  <el-button
                    size="small"
                    type="primary"
                    @click="$emit('apply-proofread', item)"
                  >
                    采纳修改
                  </el-button>
                  <el-button
                    link
                    size="small"
                    @click="$emit('ignore-proofread', item)"
                  >
                    忽略
                  </el-button>
                </template>
              </div>
            </div>
          </div>
        </div>

        <!-- 原始输出视图 -->
        <MdPreview
          v-else
          :model-value="resultText || (streaming ? '正在校对全文…' : '')"
          preview-theme="github"
          code-theme="atom"
          class="ai-result-preview"
        />
      </template>

      <!-- 任务为润色/改写等其他任务时的常规视图 -->
      <template v-else>
        <!-- 差异对比 / 最终效果 视图切换栏（存在原始选区时展示） -->
        <div v-if="originalSelection" class="ai-view-switch-row">
          <el-radio-group v-if="!editMode" v-model="viewMode" size="small">
            <el-radio-button label="diff">
              <el-icon class="mr-2px"><DocumentCopy /></el-icon>
              差异对比
            </el-radio-button>
            <el-radio-button label="preview">
              <el-icon class="mr-2px"><View /></el-icon>
              最终效果
            </el-radio-button>
          </el-radio-group>
          <span v-if="viewMode === 'diff' && !editMode" class="diff-legend">
            <span class="legend-badge legend-del">红色删除</span>
            <span class="legend-badge legend-ins">绿色新增</span>
          </span>
          <el-button
            size="small"
            :type="editMode ? 'primary' : 'default'"
            :disabled="streaming || !resultText"
            class="ml-auto"
            @click="editMode = !editMode"
          >
            <el-icon class="mr-2px"><EditPen /></el-icon>
            {{ editMode ? "完成微调" : "就地微调" }}
          </el-button>
        </div>

        <!-- 思维链（Reasoning）智能胶囊折叠 -->
        <div
          v-if="reasoningText || (streaming && !resultText)"
          class="ai-reasoning-capsule"
          :class="{ 'is-collapsed': reasoningCollapsed }"
        >
          <div
            class="ai-capsule-header"
            @click="reasoningCollapsed = !reasoningCollapsed"
          >
            <span
              class="ai-capsule-dot"
              :class="{ pulse: streaming && !resultText }"
            ></span>
            <span class="ai-capsule-title">
              {{
                streaming && !resultText
                  ? "正在深度思考…"
                  : `思考过程 (${reasoningText.length} 字${
                      meta?.elapsed ? " · " + meta.elapsed : ""
                    })`
              }}
            </span>
            <span class="ai-capsule-toggle-tip">
              {{ reasoningCollapsed ? "展开回看" : "收起" }}
            </span>
            <el-icon
              class="ai-capsule-arrow"
              :class="{ 'is-open': !reasoningCollapsed }"
            >
              <ArrowDown />
            </el-icon>
          </div>
          <div v-show="!reasoningCollapsed" class="ai-reasoning-content">
            {{ reasoningText }}
          </div>
        </div>

        <!-- 就地直接编辑模式 -->
        <div v-if="editMode" class="ai-direct-edit-box">
          <div class="ai-edit-tip text-xs text-gray-400 mb-1">
            可在下方直接编辑文字，点击「替换选区」将应用修改后的最终内容：
          </div>
          <el-input
            :model-value="resultText"
            type="textarea"
            :rows="14"
            class="ai-direct-edit-textarea"
            @update:model-value="$emit('update:resultText', $event)"
          />
        </div>

        <!-- 差异对比视图 -->
        <div
          v-else-if="originalSelection && viewMode === 'diff'"
          class="ai-diff-container"
        >
          <div v-if="!resultText && streaming" class="ai-diff-placeholder">
            正在生成并计算差异…
          </div>
          <div v-else class="ai-diff-content">
            <template v-for="(chunk, idx) in diffChunks" :key="idx">
              <del v-if="chunk.type === 'removed'" class="diff-chunk diff-del">{{ chunk.value }}</del>
              <ins v-else-if="chunk.type === 'added'" class="diff-chunk diff-ins">{{ chunk.value }}</ins>
              <span v-else class="diff-chunk diff-common">{{ chunk.value }}</span>
            </template>
          </div>
        </div>

        <!-- 原 Markdown 预览视图 -->
        <MdPreview
          v-else
          :model-value="resultText || (streaming && !reasoningText ? '生成中…' : '')"
          preview-theme="github"
          code-theme="atom"
          class="ai-result-preview"
        />
      </template>
    </div>

    <template #footer>
      <div class="ai-result-foot">
        <!-- 追加指令快捷芯片（Quick Refine Chips） -->
        <div
          v-if="!streaming && resultText && !error && lastRequest?.task !== 'proofread'"
          class="ai-quick-refine-chips"
        >
          <span class="chips-label">快捷调整：</span>
          <span
            v-for="chip in quickRefineChips"
            :key="chip"
            class="ai-refine-chip"
            @click="applyRefineChip(chip)"
          >
            {{ chip }}
          </span>
        </div>

        <!-- 追加指令：以上一版为基础继续调整 -->
        <div
          v-if="!streaming && resultText && !error && lastRequest?.task !== 'proofread'"
          class="ai-refine-row"
        >
          <el-input
            v-model="refineInput"
            size="small"
            placeholder="输入追加要求，或点击上方快捷芯片"
            clearable
            @keydown.enter="submitRefine"
          />
          <el-button
            size="small"
            type="primary"
            :disabled="!refineInput.trim()"
            @click="submitRefine"
          >
            继续调整
          </el-button>
        </div>

        <div class="ai-result-meta">
          <span v-if="history.length > 1" class="ai-version-nav">
            <el-button
              link
              size="small"
              :disabled="historyIdx <= 0"
              @click="$emit('nav-history', -1)"
            >
              ‹
            </el-button>
            {{ historyIdx + 1 }}/{{ history.length }}
            <el-button
              link
              size="small"
              :disabled="historyIdx >= history.length - 1"
              @click="$emit('nav-history', 1)"
            >
              ›
            </el-button>
          </span>
          <span v-if="meta && !error">
            {{ meta.elapsed }}{{ meta.tokens ? ` · ${meta.tokens} tokens` : "" }}
          </span>
          <span v-if="quotaText" class="ai-quota-text">
            {{ quotaText }}
          </span>
        </div>

        <div class="ai-result-actions">
          <!-- 全文校对任务的批量采纳按钮 -->
          <el-button
            v-if="lastRequest?.task === 'proofread' && pendingProofreadCount > 0"
            type="primary"
            :disabled="streaming"
            title="按 Ctrl/Cmd + Enter 快速采纳"
            @click="$emit('apply-all-proofread')"
          >
            一键采纳全部 ({{ pendingProofreadCount }}) <span class="kbd-hint">⌘⏎</span>
          </el-button>
          <el-button
            v-if="applyType === 'replace-selection'"
            type="primary"
            :disabled="streaming || !resultText"
            title="按 Ctrl/Cmd + Enter 快速替换"
            @click="$emit('apply-result')"
          >
            替换选区 <span class="kbd-hint">⌘⏎</span>
          </el-button>
          <el-button
            v-if="applyType === 'replace-selection'"
            :disabled="streaming || !resultText"
            title="保留原选区文本，在后方追加换行并插入本次生成内容"
            @click="$emit('insert-below')"
          >
            在选区后插入
          </el-button>
          <el-button
            :disabled="streaming || !resultText"
            @click="$emit('copy-result')"
          >
            复制
          </el-button>
          <el-button
            v-if="!streaming && !error"
            :disabled="!resultText"
            @click="$emit('regen')"
          >
            重新生成
          </el-button>
          <el-button
            v-if="error && !streaming"
            type="primary"
            plain
            @click="$emit('regen')"
          >
            重试
          </el-button>
          <el-button
            v-if="streaming"
            type="danger"
            plain
            @click="$emit('stop')"
          >
            停止
          </el-button>
        </div>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue";
import { MdPreview } from "md-editor-v3";
import {
  DocumentCopy,
  View,
  EditPen,
  ArrowDown,
  CircleCheck,
  Search
} from "@element-plus/icons-vue";
import { computeDiff, type DiffChunk } from "@/utils/diff";
import type { ProofreadItem } from "@/utils/proofread";
import type { AIEditRequest } from "@/api/ai";

const props = defineProps<{
  visible: boolean;
  lastRequest: AIEditRequest | null;
  taskLabel: string;
  streaming: boolean;
  error: string;
  resultText: string;
  reasoningText: string;
  meta: { elapsed?: string; tokens?: number } | null;
  quotaText: string;
  history: Array<{ label: string; content: string }>;
  historyIdx: number;
  originalSelection: string;
  applyType: "replace-selection" | "insert-cursor";
  proofreadItems: ProofreadItem[];
  pendingProofreadCount: number;
  appliedProofreadCount: number;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "update:resultText", val: string): void;
  (e: "apply-result"): void;
  (e: "insert-below"): void;
  (e: "copy-result"): void;
  (e: "regen"): void;
  (e: "stop"): void;
  (e: "refine", prompt: string): void;
  (e: "nav-history", delta: number): void;
  (e: "select-history", idx: number): void;
  (e: "apply-proofread", item: ProofreadItem): void;
  (e: "ignore-proofread", item: ProofreadItem): void;
  (e: "locate-proofread", item: ProofreadItem): void;
  (e: "apply-all-proofread"): void;
}>();

const bodyRef = ref<HTMLDivElement | null>(null);
const userScrolled = ref(false);
const viewMode = ref<"diff" | "preview">("preview");
const editMode = ref(false);
const reasoningCollapsed = ref(false);
const proofreadViewMode = ref<"cards" | "raw">("cards");
const proofreadFilter = ref<"all" | "pending" | "applied">("all");
const refineInput = ref("");

const quickRefineChips = [
  "更口语自然一点",
  "精简压缩篇幅",
  "转为清晰列表",
  "学术技术风",
  "语气更温和"
];

const filteredProofreadItems = computed(() => {
  if (proofreadFilter.value === "pending") {
    return props.proofreadItems.filter(i => !i.status || i.status === "pending");
  }
  if (proofreadFilter.value === "applied") {
    return props.proofreadItems.filter(i => i.status === "applied");
  }
  return props.proofreadItems;
});

const diffChunks = computed<DiffChunk[]>(() => {
  if (!props.originalSelection || !props.resultText) return [];
  return computeDiff(props.originalSelection, props.resultText);
});

// 流式生成时自动滚动到底部
watch(
  () => props.resultText,
  () => {
    if (props.streaming && !userScrolled.value && bodyRef.value) {
      nextTick(() => {
        if (bodyRef.value) {
          bodyRef.value.scrollTop = bodyRef.value.scrollHeight;
        }
      });
    }
  }
);

watch(
  () => props.streaming,
  streaming => {
    if (streaming) {
      userScrolled.value = false;
      reasoningCollapsed.value = false;
    }
  }
);

const onBodyScroll = () => {
  if (!bodyRef.value) return;
  const isNearBottom =
    bodyRef.value.scrollHeight - bodyRef.value.scrollTop - bodyRef.value.clientHeight < 40;
  if (!isNearBottom && props.streaming) {
    userScrolled.value = true;
  }
};

const getProofreadTagType = (type: string) => {
  if (type.includes("错") || type.includes("字")) return "danger";
  if (type.includes("法") || type.includes("语") || type.includes("病")) return "warning";
  if (type.includes("风") || type.includes("达") || type.includes("词")) return "primary";
  return "info";
};

const applyRefineChip = (chip: string) => {
  emit("refine", chip);
};

const submitRefine = () => {
  const p = refineInput.value.trim();
  if (!p) return;
  emit("refine", p);
  refineInput.value = "";
};

const handlePrimaryAction = () => {
  if (props.lastRequest?.task === "proofread") {
    if (props.pendingProofreadCount > 0 && !props.streaming) {
      emit("apply-all-proofread");
    }
  } else if (props.applyType === "replace-selection") {
    if (!props.streaming && props.resultText) {
      emit("apply-result");
    }
  }
};
</script>

<style scoped lang="scss">
.ai-drawer-header {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;

  .ai-drawer-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--el-text-color-primary);
  }

  .ai-history-pills {
    display: flex;
    gap: 4px;

    .ai-history-pill {
      font-size: 11px;
      padding: 1px 6px;
      border-radius: 4px;
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

.ai-result-body {
  height: 100%;
  overflow-y: auto;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ai-error-banner {
  margin-bottom: 8px;
}

.proofread-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--el-border-color-lighter);

  .proofread-stats {
    display: flex;
    align-items: baseline;
    gap: 6px;

    .proofread-total-badge {
      font-weight: 600;
      font-size: 13px;
      color: var(--el-text-color-primary);
    }
  }

  .proofread-header-actions {
    display: flex;
    align-items: center;
    gap: 10px;
  }
}

.proofread-filter-row {
  margin-bottom: 10px;
}

.proofread-empty {
  padding: 40px 0;
  text-align: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.proofread-cards-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.proofread-card {
  padding: 12px;
  border-radius: 8px;
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-lighter);
  display: flex;
  flex-direction: column;
  gap: 8px;

  &.is-applied {
    opacity: 0.7;
    background: var(--el-color-success-light-9);
    border-color: var(--el-color-success-light-7);
  }

  &.is-ignored {
    opacity: 0.5;
  }

  .proofread-card-header {
    display: flex;
    align-items: center;
    gap: 8px;

    .proofread-reason {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      flex: 1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .proofread-status-indicator {
      font-size: 11px;
    }
  }

  .proofread-card-diff {
    display: flex;
    flex-direction: column;
    gap: 4px;
    background: var(--el-fill-color-blank);
    padding: 8px 10px;
    border-radius: 6px;
    font-family: ui-monospace, SFMono-Regular, monospace;
    font-size: 13px;

    .diff-line {
      display: flex;
      align-items: baseline;
      gap: 6px;

      .diff-tag {
        font-size: 10px;
        padding: 0 4px;
        border-radius: 3px;

        &.del {
          background: rgba(239, 68, 68, 0.15);
          color: #dc2626;
        }

        &.ins {
          background: rgba(34, 197, 94, 0.15);
          color: #16a34a;
        }
      }

      .diff-text {
        &.del {
          color: #b91c1c;
          text-decoration: line-through;
        }

        &.ins {
          color: #15803d;
          font-weight: 600;
        }
      }
    }
  }

  .proofread-card-footer {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 8px;
  }
}

.ai-view-switch-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);

  .diff-legend {
    font-size: 11px;
    display: flex;
    gap: 6px;

    .legend-badge {
      padding: 1px 6px;
      border-radius: 3px;

      &.legend-del {
        background: rgba(239, 68, 68, 0.15);
        color: #dc2626;
      }

      &.legend-ins {
        background: rgba(34, 197, 94, 0.15);
        color: #16a34a;
      }
    }
  }
}

.ai-reasoning-capsule {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  background: var(--el-fill-color-light);
  overflow: hidden;

  .ai-capsule-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    cursor: pointer;

    .ai-capsule-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--el-color-primary);

      &.pulse {
        animation: pulse 1s infinite;
      }
    }

    .ai-capsule-title {
      font-size: 12px;
      font-weight: 500;
      color: var(--el-text-color-regular);
      flex: 1;
    }

    .ai-capsule-toggle-tip {
      font-size: 11px;
      color: var(--el-text-color-placeholder);
    }

    .ai-capsule-arrow {
      transition: transform 0.2s ease;
      &.is-open {
        transform: rotate(180deg);
      }
    }
  }

  .ai-reasoning-content {
    padding: 10px 12px;
    background: var(--el-fill-color-blank);
    font-size: 12px;
    color: var(--el-text-color-secondary);
    line-height: 1.6;
    white-space: pre-wrap;
    border-top: 1px solid var(--el-border-color-extra-light);
  }
}

.ai-diff-container {
  background: var(--el-fill-color-blank);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  padding: 12px;
  font-family: ui-monospace, SFMono-Regular, monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;

  .diff-chunk {
    &.diff-del {
      background: rgba(239, 68, 68, 0.18);
      color: #b91c1c;
      text-decoration: line-through;
    }

    &.diff-ins {
      background: rgba(34, 197, 94, 0.18);
      color: #15803d;
    }

    &.diff-common {
      color: var(--el-text-color-regular);
    }
  }
}

.ai-result-foot {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-top: 6px;

  .ai-quick-refine-chips {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 6px;
    font-size: 11px;

    .chips-label {
      color: var(--el-text-color-placeholder);
    }

    .ai-refine-chip {
      padding: 2px 8px;
      border-radius: 4px;
      background: var(--el-fill-color-light);
      color: var(--el-text-color-regular);
      cursor: pointer;
      border: 1px solid var(--el-border-color-extra-light);

      &:hover {
        background: var(--el-color-primary-light-9);
        color: var(--el-color-primary);
        border-color: var(--el-color-primary-light-5);
      }
    }
  }

  .ai-refine-row {
    display: flex;
    gap: 8px;
  }

  .ai-result-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
    color: var(--el-text-color-secondary);

    .ai-version-nav {
      display: flex;
      align-items: center;
      gap: 2px;
    }

    .ai-quota-text {
      color: var(--el-text-color-placeholder);
    }
  }

  .ai-result-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;

    .kbd-hint {
      margin-left: 4px;
      opacity: 0.7;
      font-size: 11px;
    }
  }
}

@keyframes pulse {
  0% {
    opacity: 0.4;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.4;
  }
}
</style>
