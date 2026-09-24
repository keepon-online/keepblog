<template>
  <el-dialog
    :model-value="visible"
    title="🎨 AI 极客技术封面工作台 (Canvas 16:9)"
    width="1080px"
    append-to-body
    class="ai-tech-cover-dialog"
    @close="$emit('update:visible', false)"
  >
    <div class="tech-cover-workbench">
      <!-- 左栏：Canvas 高清预览区 -->
      <div class="workbench-preview-pane">
        <div class="canvas-wrapper">
          <canvas
            ref="canvasRef"
            class="tech-cover-canvas"
            width="1280"
            height="720"
          ></canvas>
        </div>
        <div class="preview-meta-bar">
          <span class="ratio-badge">16:9 · 1280 × 720 高清输出</span>
          <span class="preview-hint">支持鼠标右键直接复制图片</span>
        </div>
      </div>

      <!-- 右栏：微调与参数定制 -->
      <div class="workbench-control-pane">
        <div class="control-header">
          <span class="control-title">封面属性与样式定制</span>
          <el-button
            size="small"
            type="primary"
            plain
            @click="syncFromPost"
          >
            <el-icon class="mr-2px"><DocumentCopy /></el-icon>
            一键同步文章信息
          </el-button>
        </div>

        <!-- 视觉主题预设选择器 -->
        <div class="control-section">
          <label class="control-label">视觉主题预设</label>
          <div class="theme-grid">
            <div
              v-for="t in TECH_COVER_THEMES"
              :key="t.id"
              class="theme-card"
              :class="{ active: currentThemeId === t.id }"
              @click="selectTheme(t.id)"
            >
              <div
                class="theme-swatch"
                :style="{ background: `linear-gradient(135deg, ${t.bgStart}, ${t.accent})` }"
              ></div>
              <div class="theme-info">
                <span class="theme-name">{{ t.name }}</span>
                <span class="theme-desc">{{ t.desc }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 标题与副标题 -->
        <div class="control-section">
          <label class="control-label">封面主标题</label>
          <el-input
            v-model="coverTitle"
            placeholder="请输入文章主标题"
            clearable
            maxlength="80"
            show-word-limit
            @input="renderCanvas"
          />
        </div>

        <div class="control-section">
          <div class="row-inputs">
            <div class="input-half">
              <label class="control-label">署名作者 / 团队</label>
              <el-input
                v-model="coverAuthor"
                placeholder="例如：KeepBlog Tech"
                clearable
                @input="renderCanvas"
              />
            </div>
            <div class="input-half">
              <label class="control-label">站点徽标文字</label>
              <el-input
                v-model="siteName"
                placeholder="例如：KEEPBLOG TECH"
                clearable
                @input="renderCanvas"
              />
            </div>
          </div>
        </div>

        <!-- 技术栈徽章推荐快捷选择 -->
        <div class="control-section">
          <div class="label-with-tip">
            <label class="control-label">技术栈徽章 (点击快速增删)</label>
            <span class="tip-sub">最多显示 4 个</span>
          </div>
          <div class="badge-picker-row">
            <el-check-tag
              v-for="badge in PRESET_TECH_BADGES"
              :key="badge"
              :checked="selectedBadges.includes(badge)"
              class="tech-badge-item"
              @change="toggleBadge(badge)"
            >
              {{ badge }}
            </el-check-tag>
          </div>
          <el-input
            v-model="customTagStr"
            placeholder="自定义标签（逗号隔开，如：K8s, 微服务, eBPF）"
            size="small"
            style="margin-top: 8px;"
            clearable
            @input="onCustomTagsInput"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <div class="workbench-footer">
        <el-button @click="$emit('update:visible', false)">取消</el-button>
        <div class="footer-actions">
          <el-button @click="downloadImage">
            <el-icon class="mr-4px"><Download /></el-icon>
            下载封面 PNG
          </el-button>
          <el-button
            type="primary"
            :loading="applying"
            @click="applyToPost"
          >
            <el-icon class="mr-4px"><Check /></el-icon>
            一键设为文章封面
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from "vue";
import { Download, Check, DocumentCopy } from "@element-plus/icons-vue";
import {
  renderTechCover,
  TECH_COVER_THEMES,
  canvasToBlob,
  downloadCanvas
} from "@/utils/techCover";
import { upload } from "@/api/common";
import { message } from "@/utils/message";

const props = defineProps<{
  visible: boolean;
  postTitle: string;
  postTags: string[];
  postCategoryName?: string;
}>();

const emit = defineEmits<{
  (e: "update:visible", val: boolean): void;
  (e: "applied", url: string): void;
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
const currentThemeId = ref("galaxy");
const coverTitle = ref("");
const coverAuthor = ref("KeepBlog Author");
const siteName = ref("KEEPBLOG TECH");
const selectedBadges = ref<string[]>(["Go", "Vue 3", "Architecture"]);
const customTagStr = ref("");
const applying = ref(false);

const PRESET_TECH_BADGES = [
  "Go",
  "Vue 3",
  "Docker",
  "Kubernetes",
  "Rust",
  "Python",
  "Linux",
  "TypeScript",
  "MySQL",
  "Redis",
  "AI / LLM",
  "Microservices",
  "DevOps",
  "Security"
];

const syncFromPost = () => {
  if (props.postTitle) {
    coverTitle.value = props.postTitle;
  }
  if (props.postTags && props.postTags.length > 0) {
    selectedBadges.value = props.postTags.slice(0, 4);
    customTagStr.value = props.postTags.join(", ");
  }
  renderCanvas();
};

const selectTheme = (themeId: string) => {
  currentThemeId.value = themeId;
  renderCanvas();
};

const toggleBadge = (badge: string) => {
  const idx = selectedBadges.value.indexOf(badge);
  if (idx > -1) {
    selectedBadges.value.splice(idx, 1);
  } else {
    if (selectedBadges.value.length >= 4) {
      selectedBadges.value.shift();
    }
    selectedBadges.value.push(badge);
  }
  customTagStr.value = selectedBadges.value.join(", ");
  renderCanvas();
};

const onCustomTagsInput = () => {
  const tags = customTagStr.value
    .split(/[,，]/)
    .map(s => s.trim())
    .filter(Boolean);
  selectedBadges.value = tags.slice(0, 4);
  renderCanvas();
};

const renderCanvas = () => {
  if (!canvasRef.value) return;
  const tags = selectedBadges.value.length > 0 ? selectedBadges.value : ["Architecture", "Code"];
  renderTechCover(canvasRef.value, {
    title: coverTitle.value || "极简现代技术博文封面",
    themeId: currentThemeId.value,
    tags,
    author: coverAuthor.value || "KeepBlog Author",
    siteName: siteName.value || "KEEPBLOG TECH",
    width: 1280,
    height: 720
  });
};

watch(
  () => props.visible,
  val => {
    if (val) {
      if (!coverTitle.value && props.postTitle) {
        coverTitle.value = props.postTitle;
      }
      if (props.postTags && props.postTags.length > 0 && selectedBadges.value.length === 0) {
        selectedBadges.value = props.postTags.slice(0, 4);
        customTagStr.value = props.postTags.join(", ");
      }
      nextTick(() => {
        renderCanvas();
      });
    }
  }
);

const downloadImage = () => {
  if (!canvasRef.value) return;
  downloadCanvas(
    canvasRef.value,
    `cover_${(coverTitle.value || "tech").slice(0, 15)}_${Date.now()}.png`
  );
  message("封面已开始下载", { type: "success" });
};

const applyToPost = async () => {
  if (!canvasRef.value) return;
  applying.value = true;
  try {
    const blob = await canvasToBlob(canvasRef.value);
    const file = new File([blob], `tech_cover_${Date.now()}.png`, {
      type: "image/png"
    });
    const formData = new FormData();
    formData.append("file", file);

    const res = await upload(formData);
    if (res?.code === 200 && res.payload) {
      emit("applied", res.payload);
      message("极客技术封面已成功上传并设为文章封面！", { type: "success" });
      emit("update:visible", false);
    } else {
      const dataUrl = canvasRef.value.toDataURL("image/png");
      emit("applied", dataUrl);
      message("已使用本地数据设为文章封面", { type: "success" });
      emit("update:visible", false);
    }
  } catch {
    const dataUrl = canvasRef.value.toDataURL("image/png");
    emit("applied", dataUrl);
    message("已设为文章封面", { type: "success" });
    emit("update:visible", false);
  } finally {
    applying.value = false;
  }
};
</script>

<style scoped lang="scss">
.tech-cover-workbench {
  display: flex;
  gap: 20px;
  height: 560px;
}

.workbench-preview-pane {
  flex: 1.1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: var(--el-fill-color-darker);
  border-radius: 8px;
  padding: 16px;
  border: 1px solid var(--el-border-color-lighter);

  .canvas-wrapper {
    width: 100%;
    aspect-ratio: 16 / 9;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.35);
    background: #000;

    .tech-cover-canvas {
      width: 100%;
      height: 100%;
      display: block;
      object-fit: contain;
    }
  }

  .preview-meta-bar {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 14px;
    font-size: 12px;
    color: var(--el-text-color-secondary);

    .ratio-badge {
      background: var(--el-fill-color);
      padding: 3px 8px;
      border-radius: 4px;
      font-weight: 500;
    }

    .preview-hint {
      color: var(--el-text-color-placeholder);
    }
  }
}

.workbench-control-pane {
  flex: 0.9;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding-right: 6px;
  gap: 16px;

  .control-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--el-border-color-lighter);

    .control-title {
      font-size: 14px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }
  }

  .control-section {
    display: flex;
    flex-direction: column;
    gap: 6px;

    .control-label {
      font-size: 12px;
      font-weight: 500;
      color: var(--el-text-color-regular);
    }

    .label-with-tip {
      display: flex;
      justify-content: space-between;
      align-items: baseline;

      .tip-sub {
        font-size: 11px;
        color: var(--el-text-color-placeholder);
      }
    }
  }

  .theme-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;

    .theme-card {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 8px 10px;
      border-radius: 6px;
      border: 1px solid var(--el-border-color-lighter);
      background: var(--el-fill-color-light);
      cursor: pointer;
      transition: all 0.2s ease;

      &:hover {
        border-color: var(--el-color-primary-light-5);
        background: var(--el-fill-color);
      }

      &.active {
        border-color: var(--el-color-primary);
        background: var(--el-color-primary-light-9);
      }

      .theme-swatch {
        width: 30px;
        height: 30px;
        border-radius: 6px;
        flex-shrink: 0;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
      }

      .theme-info {
        display: flex;
        flex-direction: column;
        overflow: hidden;

        .theme-name {
          font-size: 12px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .theme-desc {
          font-size: 10px;
          color: var(--el-text-color-secondary);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
  }

  .row-inputs {
    display: flex;
    gap: 12px;

    .input-half {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 6px;
    }
  }

  .badge-picker-row {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;

    .tech-badge-item {
      font-size: 12px;
      padding: 2px 8px;
      border-radius: 4px;
    }
  }
}

.workbench-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .footer-actions {
    display: flex;
    gap: 10px;
  }
}
</style>
