<template>
  <el-dialog
    :model-value="visible"
    title="⏳ 草稿时光机 · 本地修订历史快照"
    width="900px"
    append-to-body
    class="draft-history-dialog"
    @close="$emit('update:visible', false)"
  >
    <div v-if="revisions.length === 0" class="history-empty">
      <el-empty description="暂无历史快照。编辑过程中或点击「保存草稿」时将自动记录版本快照。">
        <template #image>
          <div class="empty-icon-wrap">
            <el-icon class="text-4xl text-gray-300"><DocumentCopy /></el-icon>
          </div>
        </template>
      </el-empty>
    </div>

    <div v-else class="history-layout">
      <!-- 左侧版本时间轴列表 -->
      <div class="history-sidebar">
        <div class="sidebar-header">
          <span class="count-badge">共 {{ revisions.length }} 个版本</span>
          <el-button link type="danger" size="small" @click="$emit('clear-all')">
            清空所有快照
          </el-button>
        </div>
        <div class="history-list">
          <div
            v-for="(rev, idx) in revisions"
            :key="rev.id"
            class="history-item"
            :class="{ active: selectedId === rev.id }"
            @click="selectRevision(rev)"
          >
            <div class="item-top">
              <span class="item-time">{{ formatTime(rev.timestamp) }}</span>
              <span class="item-tag" :class="idx === 0 ? 'latest' : ''">
                {{ idx === 0 ? '最新' : `#${revisions.length - idx}` }}
              </span>
            </div>
            <div class="item-title" :title="rev.title">
              {{ rev.title || "无标题草稿" }}
            </div>
            <div class="item-meta">
              <span>{{ rev.wordCount }} 字</span>
              <span v-if="calcDiffLen(rev) !== 0" :class="calcDiffLen(rev) > 0 ? 'len-up' : 'len-down'">
                {{ calcDiffLen(rev) > 0 ? `+${calcDiffLen(rev)}` : calcDiffLen(rev) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧对比/详情展示区 -->
      <div v-if="currentSelected" class="history-preview-panel">
        <div class="preview-toolbar">
          <div class="view-switch">
            <el-radio-group v-model="compareMode" size="small">
              <el-radio-button label="diff">与当前正文 Diff</el-radio-button>
              <el-radio-button label="raw">快照正文预览</el-radio-button>
            </el-radio-group>
          </div>
          <div class="action-buttons">
            <el-button
              type="danger"
              plain
              size="small"
              @click="$emit('delete', currentSelected.id)"
            >
              删除该快照
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="confirmRestore(currentSelected)"
            >
              恢复此版本
            </el-button>
          </div>
        </div>

        <div class="preview-content-box">
          <div v-if="compareMode === 'diff'" class="diff-viewer">
            <div class="diff-legend-bar">
              <span class="legend-badge legend-del">红色删除 (该快照有但当前已删)</span>
              <span class="legend-badge legend-ins">绿色新增 (当前正文相比新增)</span>
            </div>
            <div class="diff-render-body">
              <template v-for="(chunk, idx) in diffChunks" :key="idx">
                <del v-if="chunk.type === 'removed'" class="diff-chunk diff-del">{{ chunk.value }}</del>
                <ins v-else-if="chunk.type === 'added'" class="diff-chunk diff-ins">{{ chunk.value }}</ins>
                <span v-else class="diff-chunk diff-common">{{ chunk.value }}</span>
              </template>
            </div>
          </div>

          <div v-else class="raw-viewer">
            <pre class="raw-code-block">{{ currentSelected.content || '正文为空' }}</pre>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="history-dialog-footer">
        <span class="footer-tip">
          💡 历史快照仅保存在当前浏览器的离线存储中，清理浏览器缓存时会被重置。
        </span>
        <el-button @click="$emit('update:visible', false)">关闭</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { DocumentCopy } from "@element-plus/icons-vue";
import { computeDiff, type DiffChunk } from "@/utils/diff";
import { ElMessageBox } from "element-plus";
import type { DraftRevision } from "../hooks/useDraftHistory";

const props = defineProps<{
  visible: boolean;
  revisions: DraftRevision[];
  currentContent: string;
}>();

const emit = defineEmits<{
  (e: "update:visible", val: boolean): void;
  (e: "restore", rev: DraftRevision): void;
  (e: "delete", id: string): void;
  (e: "clear-all"): void;
}>();

const selectedId = ref<string>("");
const compareMode = ref<"diff" | "raw">("diff");

const currentSelected = computed(() => {
  return props.revisions.find(r => r.id === selectedId.value) || props.revisions[0] || null;
});

watch(
  () => props.revisions,
  newRevs => {
    if (newRevs.length > 0 && (!selectedId.value || !newRevs.some(r => r.id === selectedId.value))) {
      selectedId.value = newRevs[0].id;
    }
  },
  { immediate: true }
);

const selectRevision = (rev: DraftRevision) => {
  selectedId.value = rev.id;
};

const formatTime = (ts: number) => {
  const diff = Date.now() - ts;
  if (diff < 60000) return "刚刚";
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`;
  const d = new Date(ts);
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`;
};

const calcDiffLen = (rev: DraftRevision) => {
  return (rev.wordCount || 0) - (props.currentContent?.length || 0);
};

const diffChunks = computed<DiffChunk[]>(() => {
  if (!currentSelected.value) return [];
  // old is snapshot revision content, new is current editor content
  return computeDiff(currentSelected.value.content || "", props.currentContent || "");
});

const confirmRestore = (rev: DraftRevision) => {
  ElMessageBox.confirm(
    `确定要将当前编辑器内容恢复为快照版本「${rev.title}」吗？未保存的临时改动将被覆盖。`,
    "恢复快照确认",
    {
      confirmButtonText: "确认恢复",
      cancelButtonText: "取消",
      type: "warning"
    }
  ).then(() => {
    emit("restore", rev);
  }).catch(() => {});
};
</script>

<style scoped lang="scss">
.history-layout {
  display: flex;
  height: 520px;
  gap: 16px;
}

.history-sidebar {
  width: 260px;
  border-right: 1px solid var(--el-border-color-lighter);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 8px;
  margin-right: 12px;
  border-bottom: 1px solid var(--el-border-color-extra-light);

  .count-badge {
    font-size: 12px;
    font-weight: 600;
    color: var(--el-text-color-secondary);
  }
}

.history-list {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
  margin-top: 8px;
}

.history-item {
  padding: 10px 12px;
  border-radius: 6px;
  background: var(--el-fill-color-light);
  margin-bottom: 8px;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--el-color-primary-light-5);
    background: var(--el-fill-color);
  }

  &.active {
    border-color: var(--el-color-primary);
    background: var(--el-color-primary-light-9);
  }

  .item-top {
    display: flex;
    justify-content: space-between;
    font-size: 11px;
    color: var(--el-text-color-placeholder);
    margin-bottom: 4px;

    .item-tag {
      padding: 0 4px;
      border-radius: 3px;
      background: var(--el-fill-color-darker);
      color: var(--el-text-color-secondary);

      &.latest {
        background: var(--el-color-success-light-8);
        color: var(--el-color-success);
        font-weight: 600;
      }
    }
  }

  .item-title {
    font-size: 13px;
    font-weight: 500;
    color: var(--el-text-color-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-bottom: 4px;
  }

  .item-meta {
    display: flex;
    justify-content: space-between;
    font-size: 11px;
    color: var(--el-text-color-secondary);

    .len-up {
      color: var(--el-color-success);
      font-weight: 600;
    }

    .len-down {
      color: var(--el-color-danger);
      font-weight: 600;
    }
  }
}

.history-preview-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.preview-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.preview-content-box {
  flex: 1;
  overflow-y: auto;
  padding-top: 10px;
}

.diff-legend-bar {
  display: flex;
  gap: 12px;
  font-size: 12px;
  margin-bottom: 10px;

  .legend-badge {
    padding: 2px 8px;
    border-radius: 4px;

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

.diff-render-body {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  background: var(--el-fill-color-blank);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  padding: 12px;
  max-height: 420px;
  overflow-y: auto;

  .diff-chunk {
    &.diff-del {
      background: rgba(239, 68, 68, 0.18);
      color: #b91c1c;
      text-decoration: line-through;
    }

    &.diff-ins {
      background: rgba(34, 197, 94, 0.18);
      color: #15803d;
      text-decoration: none;
    }

    &.diff-common {
      color: var(--el-text-color-regular);
    }
  }
}

.raw-code-block {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  background: var(--el-fill-color-blank);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  padding: 12px;
  max-height: 440px;
  overflow-y: auto;
  margin: 0;
  color: var(--el-text-color-primary);
}

.history-empty {
  padding: 40px 0;
}

.history-dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .footer-tip {
    font-size: 12px;
    color: var(--el-text-color-placeholder);
  }
}
</style>
