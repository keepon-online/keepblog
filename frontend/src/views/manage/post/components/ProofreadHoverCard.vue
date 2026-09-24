<template>
  <div
    v-if="visible && item"
    class="proofread-hover-card"
    :style="{ top: `${position.top}px`, left: `${position.left}px` }"
    @click.stop
  >
    <div class="card-header">
      <el-tag size="small" :type="tagType" effect="light" class="type-tag">
        {{ item.type || "错别字" }}
      </el-tag>
      <span class="card-title">校对建议</span>
      <el-icon class="close-btn" @click="$emit('close')"><Close /></el-icon>
    </div>

    <div class="card-body">
      <div class="diff-preview">
        <span class="diff-original">{{ item.original }}</span>
        <el-icon class="diff-arrow"><Right /></el-icon>
        <span class="diff-suggestion">{{ item.suggestion }}</span>
      </div>
      <div v-if="item.reason" class="diff-reason" :title="item.reason">
        💡 {{ item.reason }}
      </div>
    </div>

    <div class="card-footer">
      <el-button
        size="small"
        type="primary"
        @click="$emit('apply', item)"
      >
        <el-icon class="mr-2px"><Check /></el-icon>
        采纳修改
      </el-button>
      <el-button
        size="small"
        text
        @click="$emit('ignore', item)"
      >
        忽略
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { Right, Check, Close } from "@element-plus/icons-vue";
import type { ProofreadItem } from "@/utils/proofread";

const props = defineProps<{
  visible: boolean;
  item: ProofreadItem | null;
  position: { top: number; left: number };
}>();

defineEmits<{
  (e: "apply", item: ProofreadItem): void;
  (e: "ignore", item: ProofreadItem): void;
  (e: "close"): void;
}>();

const tagType = computed(() => {
  const t = props.item?.type || "";
  if (t.includes("错") || t.includes("字")) return "danger";
  if (t.includes("法") || t.includes("语") || t.includes("病")) return "warning";
  if (t.includes("风") || t.includes("达") || t.includes("词")) return "primary";
  return "info";
});
</script>

<style scoped lang="scss">
.proofread-hover-card {
  position: fixed;
  z-index: 2500;
  width: 290px;
  background: var(--el-bg-color-overlay, #ffffff);
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  animation: fadeIn 0.15s ease-out;

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .type-tag {
      font-weight: 500;
    }

    .card-title {
      font-size: 12px;
      font-weight: 600;
      color: var(--el-text-color-regular);
      flex: 1;
      margin-left: 8px;
    }

    .close-btn {
      font-size: 14px;
      color: var(--el-text-color-placeholder);
      cursor: pointer;

      &:hover {
        color: var(--el-text-color-primary);
      }
    }
  }

  .card-body {
    display: flex;
    flex-direction: column;
    gap: 6px;

    .diff-preview {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;

      .diff-original {
        color: #dc2626;
        background: rgba(239, 68, 68, 0.12);
        padding: 1px 6px;
        border-radius: 4px;
        text-decoration: line-through;
      }

      .diff-arrow {
        color: var(--el-text-color-placeholder);
        font-size: 12px;
      }

      .diff-suggestion {
        color: #16a34a;
        background: rgba(34, 197, 94, 0.12);
        padding: 1px 6px;
        border-radius: 4px;
        font-weight: 600;
      }
    }

    .diff-reason {
      font-size: 11px;
      color: var(--el-text-color-secondary);
      line-height: 1.4;
      background: var(--el-fill-color-light);
      padding: 4px 8px;
      border-radius: 4px;
    }
  }

  .card-footer {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 6px;
    padding-top: 6px;
    border-top: 1px solid var(--el-border-color-extra-light);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
