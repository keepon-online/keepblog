<template>
  <div class="cover-manager-card">
    <!-- 封面预览区域 -->
    <div class="cover-preview-box">
      <template v-if="modelValue">
        <el-image
          :src="modelValue"
          fit="cover"
          class="cover-image-display"
        />
        <div class="cover-overlay">
          <el-tooltip content="预览大图" placement="top">
            <span class="overlay-btn" @click="handlePreviewCover">
              <el-icon><ZoomIn /></el-icon>
            </span>
          </el-tooltip>
          <el-tooltip
            content="设为无封面（前台展示星空流星效果）"
            placement="top"
          >
            <span
              class="overlay-btn danger"
              @click="handleClearCover"
            >
              <el-icon><Delete /></el-icon>
            </span>
          </el-tooltip>
        </div>
      </template>
      <template v-else>
        <div class="cover-empty-state">
          <div class="starry-badge">
            <el-icon><Picture /></el-icon>
            <span>无封面状态</span>
          </div>
          <span class="starry-hint">前台将自动以动态星空流星效果展示</span>
        </div>
      </template>
    </div>

    <!-- 封面控制操作区：支持5种模式 -->
    <div class="cover-toolbar">
      <!-- 方式1：本地/OSS 上传 -->
      <el-upload
        :action="uploadAction"
        :show-file-list="false"
        :auto-upload="true"
        accept="image/*"
        :before-upload="beforeUpload"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        class="upload-trigger"
      >
        <el-button size="small" type="primary" plain>
          <el-icon class="mr-2px"><Upload /></el-icon>
          上传图片
        </el-button>
      </el-upload>

      <!-- 方式2：输入网络图片 URL -->
      <el-popover
        v-model:visible="urlPopoverVisible"
        placement="bottom"
        :width="280"
        trigger="click"
      >
        <template #reference>
          <el-button size="small" plain>
            <el-icon class="mr-2px"><Link /></el-icon>
            外链地址
          </el-button>
        </template>
        <div class="url-input-popover">
          <el-input
            v-model="customImageUrl"
            placeholder="请输入图片 URL (https://...)"
            size="small"
            clearable
            @keydown.enter="applyCustomImageUrl"
          />
          <div class="popover-actions">
            <el-button
              size="small"
              @click="urlPopoverVisible = false"
            >
              取消
            </el-button>
            <el-button
              size="small"
              type="primary"
              @click="applyCustomImageUrl"
            >
              确定
            </el-button>
          </div>
        </div>
      </el-popover>

      <!-- 方式3：随机 Pixabay 封面 -->
      <el-button
        size="small"
        plain
        :loading="randomCoverLoading"
        title="从 Pixabay 获取一张高清壁纸"
        @click="handleRandomCover"
      >
        <el-icon class="mr-2px"><MagicStick /></el-icon>
        随机壁纸
      </el-button>

      <!-- 方式4：AI 极客技术封面定制 -->
      <el-button
        size="small"
        type="primary"
        plain
        title="打开视觉封面工作台定制极客卡片封面"
        @click="$emit('open-tech-cover')"
      >
        <el-icon class="mr-2px"><PictureFilled /></el-icon>
        技术封面
      </el-button>

      <!-- 方式5：设为无封面（星空流星模式） -->
      <el-button
        v-if="modelValue"
        size="small"
        type="danger"
        link
        @click="handleClearCover"
      >
        设为无封面
      </el-button>
    </div>

    <!-- 图片预览大图弹窗 -->
    <el-dialog v-model="previewVisible" title="封面大图预览" width="700px" append-to-body>
      <div style="text-align: center; max-height: 550px; overflow: hidden; display: flex; justify-content: center; align-items: center;">
        <img
          :src="previewImageUrl"
          alt="封面预览"
          style="max-width: 100%; max-height: 500px; border-radius: 8px; object-fit: contain;"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { UploadFile } from "element-plus";
import {
  Upload,
  Link,
  MagicStick,
  PictureFilled,
  ZoomIn,
  Delete,
  Picture
} from "@element-plus/icons-vue";
import { getRandomCover } from "@/api/post";
import { message } from "@/utils/message";

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", val: string): void;
  (e: "open-tech-cover"): void;
}>();

const uploadAction = import.meta.env.VITE_BASE_URL + "/api/upload/images";
const previewVisible = ref(false);
const previewImageUrl = ref("");
const randomCoverLoading = ref(false);
const urlPopoverVisible = ref(false);
const customImageUrl = ref("");

const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith("image/");
  const isLt5M = file.size / 1024 / 1024 < 5;

  if (!isImage) {
    message("只能上传图片文件", { type: "error" });
    return false;
  }
  if (!isLt5M) {
    message("图片大小不能超过 5MB", { type: "error" });
    return false;
  }
  return true;
};

const handleUploadSuccess = (res: any) => {
  if (res.code === 200) {
    emit("update:modelValue", res.data?.url || res.data);
    message("封面图片上传成功", { type: "success" });
  } else {
    message(res.message || "上传失败", { type: "error" });
  }
};

const handleUploadError = () => {
  message("上传图片失败，请检查网络或后端接口", { type: "error" });
};

const applyCustomImageUrl = () => {
  const url = customImageUrl.value.trim();
  if (!url) {
    message("请输入有效的图片 URL", { type: "warning" });
    return;
  }
  if (!/^https?:\/\/.+/i.test(url)) {
    message("图片 URL 需以 http:// 或 https:// 开头", { type: "warning" });
    return;
  }
  emit("update:modelValue", url);
  customImageUrl.value = "";
  urlPopoverVisible.value = false;
  message("已应用外链封面", { type: "success" });
};

const handleRandomCover = async () => {
  randomCoverLoading.value = true;
  try {
    const res = await getRandomCover();
    if (res.code === 200 && res.payload) {
      emit("update:modelValue", res.payload);
      message("已从壁纸库获取新封面", { type: "success" });
    } else {
      message(res.message || "获取随机封面失败", { type: "warning" });
    }
  } catch (e) {
    message("获取随机封面失败，请稍后重试", { type: "error" });
  } finally {
    randomCoverLoading.value = false;
  }
};

const handleClearCover = () => {
  emit("update:modelValue", "");
  message("已清除封面，前台将以动态星空流星效果展示", { type: "info" });
};

const handlePreviewCover = () => {
  previewImageUrl.value = props.modelValue;
  previewVisible.value = true;
};
</script>

<style scoped lang="scss">
.cover-manager-card {
  display: flex;
  flex-direction: column;
  gap: 10px;

  .cover-preview-box {
    position: relative;
    width: 100%;
    height: 140px;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid var(--el-border-color);
    background-color: var(--el-fill-color-lighter);

    .cover-image-display {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: block;
    }

    .cover-overlay {
      position: absolute;
      inset: 0;
      background: rgba(0, 0, 0, 0.55);
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 16px;
      opacity: 0;
      transition: opacity 0.25s ease;

      .overlay-btn {
        width: 34px;
        height: 34px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.2);
        color: #fff;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 18px;
        cursor: pointer;
        backdrop-filter: blur(4px);
        transition: all 0.2s ease;

        &:hover {
          background: rgba(255, 255, 255, 0.4);
          transform: scale(1.1);
        }

        &.danger:hover {
          background: rgba(239, 68, 68, 0.8);
        }
      }
    }

    &:hover .cover-overlay {
      opacity: 1;
    }

    .cover-empty-state {
      height: 100%;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 8px;
      background: radial-gradient(ellipse at bottom, #1b2735 0%, #090a0f 100%);
      color: #94a3b8;

      .starry-badge {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        font-weight: 500;
        color: #e2e8f0;
      }

      .starry-hint {
        font-size: 11px;
        color: #64748b;
      }
    }
  }

  .cover-toolbar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;

    .upload-trigger {
      display: inline-flex;
    }
  }
}

.url-input-popover {
  display: flex;
  flex-direction: column;
  gap: 10px;

  .popover-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}
</style>
