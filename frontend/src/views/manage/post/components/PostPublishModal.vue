<template>
  <el-dialog
    :model-value="visible"
    width="540px"
    append-to-body
    :show-close="false"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    class="post-celebrate-dialog"
  >
    <div class="celebrate-body">
      <div class="celebrate-badge">
        <el-icon class="badge-icon"><CircleCheckFilled /></el-icon>
      </div>

      <h3 class="celebrate-title">
        {{ isEdit ? "文章更新成功！" : "🎉 恭喜，文章发布成功！" }}
      </h3>
      <p class="celebrate-subtitle">
        您的博文现已就绪并同步至生产环境。
      </p>

      <div class="post-info-card">
        <div class="info-row title-row">
          <span class="info-label">文章标题：</span>
          <span class="info-value font-semibold">{{ postTitle || "无标题" }}</span>
        </div>
        <div class="info-row link-row">
          <span class="info-label">访问链接：</span>
          <el-link
            type="primary"
            :href="publicUrl"
            target="_blank"
            class="info-link"
          >
            {{ publicUrl }}
          </el-link>
        </div>
      </div>

      <div class="share-actions-row">
        <el-button
          type="primary"
          @click="openPublicPage"
        >
          <el-icon class="mr-4px"><View /></el-icon>
          在新窗口打开文章
        </el-button>
        <el-button
          @click="copyShareLink"
        >
          <el-icon class="mr-4px"><DocumentCopy /></el-icon>
          复制文章链接
        </el-button>
      </div>

      <div class="dialog-foot-nav">
        <el-button text @click="$emit('continue-edit')">
          继续编辑
        </el-button>
        <span class="nav-divider">·</span>
        <el-button text @click="$emit('write-another')">
          再写一篇
        </el-button>
        <span class="nav-divider">·</span>
        <el-button text type="primary" @click="$emit('back-to-list')">
          返回文章列表
        </el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { CircleCheckFilled, View, DocumentCopy } from "@element-plus/icons-vue";
import { message } from "@/utils/message";

const props = defineProps<{
  visible: boolean;
  isEdit: boolean;
  postId?: string | number;
  postTitle: string;
}>();

defineEmits<{
  (e: "continue-edit"): void;
  (e: "write-another"): void;
  (e: "back-to-list"): void;
}>();

const publicUrl = computed(() => {
  const origin = window.location.origin;
  return `${origin}/post/${props.postId || ""}`;
});

const openPublicPage = () => {
  window.open(publicUrl.value, "_blank");
};

const copyShareLink = () => {
  navigator.clipboard.writeText(publicUrl.value);
  message("文章公开链接已复制到剪贴板", { type: "success" });
};
</script>

<style scoped lang="scss">
.post-celebrate-dialog {
  :deep(.el-dialog__header) {
    display: none;
  }
}

.celebrate-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 20px 10px 10px;

  .celebrate-badge {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background: var(--el-color-success-light-9);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 16px;

    .badge-icon {
      font-size: 38px;
      color: var(--el-color-success);
    }
  }

  .celebrate-title {
    font-size: 18px;
    font-weight: 700;
    color: var(--el-text-color-primary);
    margin: 0 0 6px;
  }

  .celebrate-subtitle {
    font-size: 13px;
    color: var(--el-text-color-secondary);
    margin: 0 0 20px;
  }

  .post-info-card {
    width: 100%;
    background: var(--el-fill-color-light);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    padding: 12px 16px;
    margin-bottom: 20px;
    text-align: left;
    display: flex;
    flex-direction: column;
    gap: 8px;

    .info-row {
      display: flex;
      align-items: center;
      font-size: 13px;

      .info-label {
        color: var(--el-text-color-secondary);
        width: 72px;
        flex-shrink: 0;
      }

      .info-value {
        color: var(--el-text-color-primary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .info-link {
        font-size: 12px;
        word-break: break-all;
      }
    }
  }

  .share-actions-row {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
  }

  .dialog-foot-nav {
    display: flex;
    align-items: center;
    gap: 8px;
    border-top: 1px solid var(--el-border-color-extra-light);
    width: 100%;
    padding-top: 14px;
    justify-content: center;

    .nav-divider {
      color: var(--el-border-color);
      font-size: 12px;
    }
  }
}
</style>
