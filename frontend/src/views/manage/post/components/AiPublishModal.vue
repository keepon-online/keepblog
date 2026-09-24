<template>
  <el-dialog
    :model-value="visible"
    title="✨ AI 发文助手与质量体检中心"
    width="740px"
    append-to-body
    class="ai-publish-dialog"
    :before-close="handleClose"
  >
    <el-tabs v-model="activeTab" class="publish-tabs">
      <!-- 选项卡 1：AI 元数据智能提炼 -->
      <el-tab-pane label="AI 智能提炼 (Copilot)" name="copilot">
        <div class="copilot-tab-body">
          <el-alert
            type="info"
            :closable="false"
            show-icon
            title="基于全文深度提炼推荐标题、精炼摘要与精准匹配标签，确认后可一键注入文章表单。"
            style="margin-bottom: 16px"
          />

          <!-- 1. 推荐标题 -->
          <div class="copilot-section">
            <div class="copilot-section-header">
              <span class="copilot-section-title">
                <el-icon class="mr-4px text-primary"><CollectionTag /></el-icon>
                推荐候选标题
              </span>
              <el-button
                link
                type="primary"
                size="small"
                :loading="copilotTitleLoading"
                @click="regenTitle"
              >
                换一批
              </el-button>
            </div>
            <div class="copilot-title-list">
              <el-radio-group v-model="chosenTitle" class="copilot-radio-group">
                <el-radio
                  v-for="(t, idx) in titleCandidates"
                  :key="idx"
                  :label="t"
                  class="copilot-title-radio"
                >
                  {{ t }}
                </el-radio>
                <el-radio
                  v-if="currentTitle"
                  :label="currentTitle"
                  class="copilot-title-radio text-gray-500"
                >
                  保持当前标题：{{ currentTitle }}
                </el-radio>
              </el-radio-group>
            </div>
          </div>

          <!-- 2. 文章摘要 -->
          <div class="copilot-section">
            <div class="copilot-section-header">
              <span class="copilot-section-title">
                <el-icon class="mr-4px text-primary"><Tickets /></el-icon>
                精炼核心摘要
              </span>
              <el-button
                link
                type="primary"
                size="small"
                :loading="copilotSummaryLoading"
                @click="regenSummary"
              >
                重新提炼
              </el-button>
            </div>
            <el-input
              v-model="summaryText"
              type="textarea"
              :rows="3"
              maxlength="500"
              show-word-limit
              placeholder="AI 正在生成摘要…"
            />
          </div>

          <!-- 3. 标签匹配与建议 -->
          <div class="copilot-section">
            <div class="copilot-section-header">
              <span class="copilot-section-title">
                <el-icon class="mr-4px text-primary"><PriceTag /></el-icon>
                智能标签匹配 (比对系统已有标签库)
              </span>
              <el-button
                link
                type="primary"
                size="small"
                :loading="copilotTagsLoading"
                @click="regenTags"
              >
                重新建议
              </el-button>
            </div>

            <!-- 命中已有标签 -->
            <div v-if="matchedTags.length > 0" class="copilot-tag-subgroup">
              <div class="subgroup-label">命中系统已有标签（建议勾选复用）：</div>
              <div class="copilot-tags-wrap">
                <el-check-tag
                  v-for="t in matchedTags"
                  :key="t"
                  :checked="selectedTags.includes(t)"
                  class="copilot-tag-item matched"
                  @change="toggleTag(t)"
                >
                  <el-icon class="mr-2px"><Check /></el-icon>
                  {{ t }}
                </el-check-tag>
              </div>
            </div>

            <!-- 推荐新建的标签 -->
            <div v-if="newTags.length > 0" class="copilot-tag-subgroup">
              <div class="subgroup-label">推荐新建标签：</div>
              <div class="copilot-tags-wrap">
                <el-check-tag
                  v-for="t in newTags"
                  :key="t"
                  :checked="selectedTags.includes(t)"
                  class="copilot-tag-item new-tag"
                  @change="toggleTag(t)"
                >
                  + {{ t }}
                </el-check-tag>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 选项卡 2：SEO 与发文质量体检 -->
      <el-tab-pane label="发文质量与 SEO 体检" name="checklist">
        <div class="checklist-tab-body">
          <!-- 总体健康度卡片 -->
          <div class="score-card" :class="scoreClass">
            <div class="score-left">
              <div class="score-num">{{ healthScore }}</div>
              <div class="score-desc">
                <div class="score-title">发文综合健康度</div>
                <div class="score-subtitle">{{ scoreLabel }}</div>
              </div>
            </div>
            <div class="score-meta">
              <span>全文共约 {{ wordCount }} 字</span>
              <span>· 预估阅读时间约 {{ readMinutes }} 分钟</span>
            </div>
          </div>

          <!-- 各项指标列表 -->
          <div class="check-items-list">
            <!-- 标题体检 -->
            <div class="check-item">
              <el-icon :class="titleCheck.iconClass"><component :is="titleCheck.icon" /></el-icon>
              <div class="check-content">
                <div class="check-title">文章标题 ({{ currentTitle.length }} 字符)</div>
                <div class="check-hint">{{ titleCheck.hint }}</div>
              </div>
            </div>

            <!-- 摘要体检 -->
            <div class="check-item">
              <el-icon :class="summaryCheck.iconClass"><component :is="summaryCheck.icon" /></el-icon>
              <div class="check-content">
                <div class="check-title">文章摘要 ({{ (summaryText || currentSummary).length }} 字符)</div>
                <div class="check-hint">{{ summaryCheck.hint }}</div>
              </div>
            </div>

            <!-- 分类与标签体检 -->
            <div class="check-item">
              <el-icon :class="taxonomyCheck.iconClass"><component :is="taxonomyCheck.icon" /></el-icon>
              <div class="check-content">
                <div class="check-title">分类与标签</div>
                <div class="check-hint">{{ taxonomyCheck.hint }}</div>
              </div>
            </div>

            <!-- 封面体检 -->
            <div class="check-item">
              <el-icon :class="coverCheck.iconClass"><component :is="coverCheck.icon" /></el-icon>
              <div class="check-content">
                <div class="check-title">文章封面配置</div>
                <div class="check-hint">{{ coverCheck.hint }}</div>
              </div>
            </div>

            <!-- 代码块语法高亮体检 -->
            <div class="check-item">
              <el-icon :class="codeBlockCheck.iconClass"><component :is="codeBlockCheck.icon" /></el-icon>
              <div class="check-content">
                <div class="check-title">代码块语法高亮检查</div>
                <div class="check-hint">{{ codeBlockCheck.hint }}</div>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <div class="publish-dialog-foot">
        <div class="foot-tip text-xs text-gray-400">
          已选：标题 {{ chosenTitle ? '✓' : '-' }} · 摘要 {{ summaryText ? '✓' : '-' }} · 标签 ({{ selectedTags.length }})
        </div>
        <div class="foot-actions">
          <el-button @click="handleClose">取消</el-button>
          <el-button
            type="primary"
            plain
            :disabled="!chosenTitle && !summaryText && selectedTags.length === 0"
            @click="applyToFormOnly"
          >
            仅应用到表单
          </el-button>
          <el-button
            type="primary"
            :loading="submitLoading"
            @click="applyAndPublish"
          >
            应用并立即提交
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  CollectionTag,
  Tickets,
  PriceTag,
  Check,
  CircleCheckFilled,
  WarningFilled,
  CircleCloseFilled
} from "@element-plus/icons-vue";
import { streamAIEdit } from "@/api/ai";

const props = defineProps<{
  visible: boolean;
  currentTitle: string;
  currentSummary: string;
  currentTags: string[];
  currentCategoryId?: number | string;
  currentCover?: string;
  postContent: string;
  systemTags: any[];
  submitLoading?: boolean;
}>();

const emit = defineEmits<{
  (e: "update:visible", val: boolean): void;
  (
    e: "apply",
    payload: {
      title: string;
      summary: string;
      tags: string[];
      publishImmediately: boolean;
    }
  ): void;
}>();

const activeTab = ref("copilot");
const chosenTitle = ref("");
const summaryText = ref("");
const selectedTags = ref<string[]>([]);
const matchedTags = ref<string[]>([]);
const newTags = ref<string[]>([]);
const titleCandidates = ref<string[]>([]);

const copilotTitleLoading = ref(false);
const copilotSummaryLoading = ref(false);
const copilotTagsLoading = ref(false);

let abortTitle: (() => void) | null = null;
let abortSummary: (() => void) | null = null;
let abortTags: (() => void) | null = null;

watch(
  () => props.visible,
  val => {
    if (val) {
      chosenTitle.value = props.currentTitle || "";
      summaryText.value = props.currentSummary || "";
      selectedTags.value = [...(props.currentTags || [])];
      matchedTags.value = [];
      newTags.value = [];
      titleCandidates.value = [];

      regenTitle();
      regenSummary();
      regenTags();
    } else {
      stopAll();
    }
  }
);

const stopAll = () => {
  abortTitle?.();
  abortSummary?.();
  abortTags?.();
  copilotTitleLoading.value = false;
  copilotSummaryLoading.value = false;
  copilotTagsLoading.value = false;
};

const handleClose = () => {
  stopAll();
  emit("update:visible", false);
};

const regenTitle = () => {
  if (!props.postContent) return;
  copilotTitleLoading.value = true;
  titleCandidates.value = [];
  let acc = "";
  abortTitle = streamAIEdit(
    { task: "title", digest: props.postContent },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        copilotTitleLoading.value = false;
        const candidates = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 4);
        titleCandidates.value = candidates;
        if (!chosenTitle.value && candidates.length > 0) {
          chosenTitle.value = candidates[0];
        }
      },
      onError: () => {
        copilotTitleLoading.value = false;
      }
    }
  );
};

const regenSummary = () => {
  if (!props.postContent) return;
  copilotSummaryLoading.value = true;
  let acc = "";
  abortSummary = streamAIEdit(
    { task: "summary", digest: props.postContent },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        copilotSummaryLoading.value = false;
        if (acc.trim()) {
          summaryText.value = acc.trim();
        }
      },
      onError: () => {
        copilotSummaryLoading.value = false;
      }
    }
  );
};

const regenTags = () => {
  if (!props.postContent) return;
  copilotTagsLoading.value = true;
  matchedTags.value = [];
  newTags.value = [];
  let acc = "";
  abortTags = streamAIEdit(
    {
      task: "tags",
      digest: props.postContent,
      title: chosenTitle.value || props.currentTitle
    },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        copilotTagsLoading.value = false;
        const recTags = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 8);

        const systemTagNames: string[] = (props.systemTags || []).map(
          (t: any) => t.tagName || ""
        );

        const matched: string[] = [];
        const created: string[] = [];

        recTags.forEach(tag => {
          const found = systemTagNames.find(
            st => st.toLowerCase() === tag.toLowerCase()
          );
          if (found) {
            if (!matched.includes(found)) matched.push(found);
            if (!selectedTags.value.includes(found)) {
              selectedTags.value.push(found);
            }
          } else {
            if (!created.includes(tag)) created.push(tag);
          }
        });

        matchedTags.value = matched;
        newTags.value = created;
      },
      onError: () => {
        copilotTagsLoading.value = false;
      }
    }
  );
};

const toggleTag = (t: string) => {
  const idx = selectedTags.value.indexOf(t);
  if (idx > -1) {
    selectedTags.value.splice(idx, 1);
  } else {
    selectedTags.value.push(t);
  }
};

// ==================== SEO & 质量体检计算 ====================
const wordCount = computed(() => (props.postContent || "").replace(/\s+/g, "").length);
const readMinutes = computed(() => Math.max(1, Math.ceil(wordCount.value / 350)));

const unlabeledCodeBlockCount = computed(() => {
  const matches = (props.postContent || "").match(/```[ \t]*\r?\n/g);
  return matches ? matches.length : 0;
});

const titleCheck = computed(() => {
  const len = props.currentTitle?.length || 0;
  if (len === 0) {
    return { icon: CircleCloseFilled, iconClass: "text-red-500", hint: "缺少文章标题，必须填写" };
  }
  if (len < 6) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "标题偏短（少于6字），不利于搜索引擎收录" };
  }
  if (len > 50) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "标题略长（超过50字），展示时可能会被省略截断" };
  }
  return { icon: CircleCheckFilled, iconClass: "text-green-500", hint: "标题长度适中，符合 SEO 最佳实践" };
});

const summaryCheck = computed(() => {
  const text = summaryText.value || props.currentSummary || "";
  if (!text) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "未配置文章摘要，前台将自动截取正文首段" };
  }
  if (text.length < 20) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "摘要过短，建议提炼更完整的核心观点" };
  }
  return { icon: CircleCheckFilled, iconClass: "text-green-500", hint: "摘要提炼完备，有利于社交分享预览与卡片展示" };
});

const taxonomyCheck = computed(() => {
  const hasCat = !!props.currentCategoryId;
  const tagLen = (selectedTags.value.length || props.currentTags?.length || 0);
  if (!hasCat && tagLen === 0) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "未设置分类与标签，文章将难以被分类聚合归档" };
  }
  if (!hasCat) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "未设置文章分类，建议为其指定一个主分类" };
  }
  if (tagLen === 0) {
    return { icon: WarningFilled, iconClass: "text-orange-500", hint: "未设置标签，建议选择 1~4 个技术关键词" };
  }
  return { icon: CircleCheckFilled, iconClass: "text-green-500", hint: "分类与标签完备，文章归档结构清晰" };
});

const coverCheck = computed(() => {
  if (props.currentCover) {
    return { icon: CircleCheckFilled, iconClass: "text-green-500", hint: "已配置高清特色封面图片" };
  }
  return { icon: CircleCheckFilled, iconClass: "text-blue-500", hint: "未配置图片封面，前台将以动态星空流星效果展示（符合预期）" };
});

const codeBlockCheck = computed(() => {
  const count = unlabeledCodeBlockCount.value;
  if (count > 0) {
    return {
      icon: WarningFilled,
      iconClass: "text-orange-500",
      hint: "检测到 " + count + " 处未声明语言的代码块 (```)，建议指定语言 (如 ```go) 以获得精准语法高亮"
    };
  }
  return { icon: CircleCheckFilled, iconClass: "text-green-500", hint: "所有代码块均已明确声明语法高亮语言" };
});

const healthScore = computed(() => {
  let score = 0;
  if (titleCheck.value.iconClass.includes("text-green-500")) score += 25;
  else if (titleCheck.value.iconClass.includes("text-orange-500")) score += 15;

  if (summaryCheck.value.iconClass.includes("text-green-500")) score += 20;
  else if (summaryCheck.value.iconClass.includes("text-orange-500")) score += 10;

  if (taxonomyCheck.value.iconClass.includes("text-green-500")) score += 20;
  else if (taxonomyCheck.value.iconClass.includes("text-orange-500")) score += 10;

  score += 15; // cover is either custom or starry
  if (codeBlockCheck.value.iconClass.includes("text-green-500")) score += 10;
  if (wordCount.value >= 100) score += 10;

  return Math.min(100, score);
});

const scoreClass = computed(() => {
  if (healthScore.value >= 90) return "score-perfect";
  if (healthScore.value >= 75) return "score-good";
  return "score-warning";
});

const scoreLabel = computed(() => {
  if (healthScore.value >= 90) return "极佳 · 符合高标准技术出版规范";
  if (healthScore.value >= 75) return "良好 · 建议参考体检建议优化";
  return "待完善 · 建议补充标题或摘要";
});

const applyToFormOnly = () => {
  emit("apply", {
    title: chosenTitle.value || props.currentTitle,
    summary: summaryText.value || props.currentSummary,
    tags: Array.from(new Set([...props.currentTags, ...selectedTags.value])),
    publishImmediately: false
  });
  handleClose();
};

const applyAndPublish = () => {
  emit("apply", {
    title: chosenTitle.value || props.currentTitle,
    summary: summaryText.value || props.currentSummary,
    tags: Array.from(new Set([...props.currentTags, ...selectedTags.value])),
    publishImmediately: true
  });
  handleClose();
};
</script>

<style scoped lang="scss">
.publish-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 16px;
  }
}

.copilot-tab-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 480px;
  overflow-y: auto;
  padding-right: 6px;
}

.copilot-section {
  background: var(--el-fill-color-lighter);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 12px 14px;

  .copilot-section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;

    .copilot-section-title {
      font-size: 13px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      display: flex;
      align-items: center;
    }
  }

  .copilot-radio-group {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    width: 100%;

    .copilot-title-radio {
      margin-right: 0;
      white-space: normal;
      height: auto;
      line-height: 1.4;
      padding: 6px 8px;
      border-radius: 6px;
      width: 100%;

      &:hover {
        background: var(--el-fill-color);
      }
    }
  }

  .copilot-tag-subgroup {
    margin-top: 10px;

    .subgroup-label {
      font-size: 11px;
      color: var(--el-text-color-secondary);
      margin-bottom: 6px;
    }

    .copilot-tags-wrap {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;

      .copilot-tag-item {
        font-size: 12px;
        padding: 3px 10px;
        border-radius: 4px;
      }
    }
  }
}

.checklist-tab-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 480px;
  overflow-y: auto;
  padding-right: 6px;

  .score-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-radius: 8px;
    background: var(--el-fill-color-light);
    border: 1px solid var(--el-border-color-lighter);

    &.score-perfect {
      background: linear-gradient(135deg, rgba(34, 197, 94, 0.1), rgba(16, 185, 129, 0.05));
      border-color: rgba(34, 197, 94, 0.3);

      .score-num {
        color: #16a34a;
      }
    }

    &.score-good {
      background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(37, 99, 235, 0.05));
      border-color: rgba(59, 130, 246, 0.3);

      .score-num {
        color: #2563eb;
      }
    }

    &.score-warning {
      background: linear-gradient(135deg, rgba(245, 158, 11, 0.1), rgba(217, 119, 6, 0.05));
      border-color: rgba(245, 158, 11, 0.3);

      .score-num {
        color: #d97706;
      }
    }

    .score-left {
      display: flex;
      align-items: center;
      gap: 16px;

      .score-num {
        font-size: 38px;
        font-weight: 800;
        line-height: 1;
      }

      .score-title {
        font-size: 14px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .score-subtitle {
        font-size: 12px;
        color: var(--el-text-color-secondary);
        margin-top: 2px;
      }
    }

    .score-meta {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }

  .check-items-list {
    display: flex;
    flex-direction: column;
    gap: 10px;

    .check-item {
      display: flex;
      align-items: flex-start;
      gap: 12px;
      padding: 12px 14px;
      border-radius: 6px;
      background: var(--el-fill-color-blank);
      border: 1px solid var(--el-border-color-lighter);

      .el-icon {
        font-size: 18px;
        margin-top: 2px;
      }

      .check-title {
        font-size: 13px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .check-hint {
        font-size: 12px;
        color: var(--el-text-color-secondary);
        margin-top: 2px;
      }
    }
  }
}

.publish-dialog-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .foot-actions {
    display: flex;
    gap: 10px;
  }
}
</style>
