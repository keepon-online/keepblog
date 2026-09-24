<template>
  <el-dialog
    :model-value="visible"
    title="✨ AI 主题大纲生成器"
    width="780px"
    append-to-body
    class="ai-outline-dialog"
    :before-close="closeDialog"
  >
    <div class="outline-modal-body">
      <el-alert
        type="info"
        :closable="false"
        show-icon
        title="输入文章构思或主题，选择目标受众深度与结构规模，AI 将流式构建严密的多级 Markdown 大纲与要点提示。"
        style="margin-bottom: 16px"
      />

      <div class="outline-config-grid">
        <div class="config-item topic-item">
          <label class="config-label">文章主题 / 核心构思</label>
          <el-input
            v-model="outlineTopic"
            placeholder="例如：Go 语言并发原语与微服务实战陷阱"
            clearable
            maxlength="150"
            show-word-limit
          />
        </div>

        <div class="config-row">
          <div class="config-item">
            <label class="config-label">读者定位与深度偏好</label>
            <el-radio-group v-model="outlineDepth" size="default">
              <el-radio-button label="beginner">🟢 入门教程</el-radio-button>
              <el-radio-button label="advanced">🔵 进阶实战</el-radio-button>
              <el-radio-button label="architecture">🟣 架构选型</el-radio-button>
            </el-radio-group>
          </div>

          <div class="config-item">
            <label class="config-label">篇幅规模</label>
            <el-radio-group v-model="outlineScale" size="default">
              <el-radio-button label="concise">⚡ 精简核心 (3~4章)</el-radio-button>
              <el-radio-button label="detailed">📖 深入详尽 (5~7章)</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </div>

      <div class="outline-action-bar">
        <el-button
          type="primary"
          :loading="outlineLoading"
          :disabled="!outlineTopic.trim() || aiQuotaExhausted"
          @click="startOutlineGeneration"
        >
          <el-icon class="mr-4px"><Promotion /></el-icon>
          {{ outlineResult ? "重新生成大纲" : "开始生成大纲" }}
        </el-button>
        <el-button
          v-if="outlineLoading"
          type="danger"
          plain
          size="small"
          @click="stopOutlineGeneration"
        >
          停止生成
        </el-button>
      </div>

      <!-- 大纲结果展示区（带 Markdown 预览） -->
      <div v-if="outlineResult || outlineLoading" class="outline-preview-container">
        <div class="preview-header">
          <span class="preview-title">
            大纲预览
            <span v-if="outlineLoading" class="outline-streaming-tag">生成中…</span>
          </span>
          <el-button
            v-if="outlineResult"
            size="small"
            text
            @click="copyOutlineContent"
          >
            复制大纲
          </el-button>
        </div>
        <div class="outline-preview-scroll">
          <MdPreview
            :model-value="outlineResult"
            language="zh-CN"
            code-theme="atom"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <div class="outline-foot">
        <span class="outline-foot-tip text-xs text-gray-400">
          {{ outlineResult ? `已生成约 ${outlineResult.length} 字符` : "配置后点击开始生成" }}
        </span>
        <div class="outline-foot-actions">
          <el-button @click="closeDialog">关闭</el-button>
          <el-dropdown
            v-if="outlineResult && !outlineLoading"
            trigger="click"
            @command="handleApplyOutlineCommand"
          >
            <el-button type="primary">
              应用到编辑器 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="replace">
                  覆盖当前正文
                </el-dropdown-item>
                <el-dropdown-item command="append">
                  追加到正文末尾
                </el-dropdown-item>
                <el-dropdown-item command="sync-title" divided>
                  同时提取并填充主标题与摘要
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { MdPreview } from "md-editor-v3";
import { Promotion, ArrowDown } from "@element-plus/icons-vue";
import { streamAIEdit } from "@/api/ai";
import { message } from "@/utils/message";

const props = defineProps<{
  visible: boolean;
  initialTopic?: string;
  aiModelChoice?: string;
  aiQuotaExhausted?: boolean;
}>();

const emit = defineEmits<{
  (e: "update:visible", val: boolean): void;
  (
    e: "apply",
    payload: {
      content: string;
      mode: "replace" | "append" | "sync-title";
      title?: string;
      summary?: string;
    }
  ): void;
}>();

const outlineTopic = ref("");
const outlineDepth = ref<"beginner" | "advanced" | "architecture">("advanced");
const outlineScale = ref<"concise" | "detailed">("detailed");
const outlineResult = ref("");
const outlineLoading = ref(false);
let outlineAbort: (() => void) | null = null;

watch(
  () => props.visible,
  val => {
    if (val) {
      if (props.initialTopic && !outlineTopic.value) {
        outlineTopic.value = props.initialTopic;
      }
    } else {
      stopOutlineGeneration();
    }
  }
);

const closeDialog = () => {
  stopOutlineGeneration();
  emit("update:visible", false);
};

const stopOutlineGeneration = () => {
  if (outlineAbort) {
    outlineAbort();
    outlineAbort = null;
  }
  outlineLoading.value = false;
};

const startOutlineGeneration = () => {
  const topic = outlineTopic.value.trim();
  if (!topic) {
    message("请输入文章主题或构思", { type: "warning" });
    return;
  }

  stopOutlineGeneration();
  outlineLoading.value = true;
  outlineResult.value = "";

  const depthMap = {
    beginner: "定位为面向初学者的平易近人教程，重基础概念与循序渐进示例",
    advanced: "定位为面向中高级工程师的深度实战，重底层原理解析与生产避坑经验",
    architecture: "定位为技术架构师视角，重方案选型权衡、扩展性与多方案横向对比"
  };

  const scaleMap = {
    concise: "精简紧凑版，提炼3~4个核心章节，快速落地",
    detailed: "深入详尽版，包含5~7个多级深入章节，涵盖演进过程、实战代码架构与最佳实践"
  };

  const instruction = `${depthMap[outlineDepth.value]}；${scaleMap[outlineScale.value]}`;

  outlineAbort = streamAIEdit(
    {
      task: "outline",
      title: topic,
      instruction,
      // "default" 是前端"默认模型"选项的哨兵值，后端按空 name 走顶层默认
      model: props.aiModelChoice === "default" ? undefined : props.aiModelChoice
    },
    {
      onDelta: delta => {
        outlineResult.value += delta;
      },
      onDone: () => {
        outlineLoading.value = false;
        outlineAbort = null;
      },
      onError: err => {
        message(err || "生成大纲失败", { type: "error" });
        outlineLoading.value = false;
        outlineAbort = null;
      }
    }
  );
};

const copyOutlineContent = () => {
  if (!outlineResult.value) return;
  navigator.clipboard.writeText(outlineResult.value);
  message("大纲内容已复制到剪贴板", { type: "success" });
};

const handleApplyOutlineCommand = (command: "replace" | "append" | "sync-title") => {
  if (!outlineResult.value) return;

  let title: string | undefined;
  let summary: string | undefined;

  const titleMatch = outlineResult.value.match(/^#\s+(.+)$/m);
  if (titleMatch) {
    title = titleMatch[1].trim();
  }

  const summaryMatch = outlineResult.value.match(/> 💡 写作要点：\s*(.+)/);
  if (summaryMatch) {
    summary = summaryMatch[1].trim();
  }

  emit("apply", {
    content: outlineResult.value,
    mode: command,
    title,
    summary
  });

  closeDialog();
};
</script>

<style scoped lang="scss">
.ai-outline-dialog {
  .outline-modal-body {
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .outline-config-grid {
    display: flex;
    flex-direction: column;
    gap: 12px;
    background: var(--el-fill-color-lighter, #fafafa);
    padding: 14px;
    border-radius: 8px;
    border: 1px solid var(--el-border-color-lighter, #ebeef5);

    .config-label {
      display: block;
      font-size: 12px;
      font-weight: 500;
      color: var(--el-text-color-regular, #606266);
      margin-bottom: 6px;
    }

    .config-row {
      display: flex;
      gap: 20px;
      flex-wrap: wrap;
    }
  }

  .outline-action-bar {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .outline-preview-container {
    border-radius: 8px;
    border: 1px solid var(--el-border-color-lighter, #ebeef5);
    background: var(--el-bg-color, #fff);
    overflow: hidden;

    .preview-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 14px;
      background: var(--el-fill-color-light, #f5f7fa);
      border-bottom: 1px solid var(--el-border-color-lighter, #ebeef5);

      .preview-title {
        font-size: 13px;
        font-weight: 600;
        color: var(--el-text-color-primary, #303133);
        display: flex;
        align-items: center;
        gap: 8px;
      }

      .outline-streaming-tag {
        font-size: 11px;
        color: var(--el-color-primary, #409eff);
      }
    }

    .outline-preview-scroll {
      max-height: 380px;
      overflow-y: auto;
      padding: 14px;
    }
  }

  .outline-foot {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }
}
</style>
