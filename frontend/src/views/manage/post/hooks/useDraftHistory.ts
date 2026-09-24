import { ref, computed } from "vue";
import type { Ref } from "vue";
import { localForage } from "@/utils/localforage";
import { message } from "@/utils/message";
import { ElMessageBox } from "element-plus";
import { useDebounceFn } from "@vueuse/core";

export interface PostDraft {
  ruleForm: any;
  savedAt: number;
}

export interface DraftRevision {
  id: string;
  timestamp: number;
  title: string;
  wordCount: number;
  summary: string;
  content: string;
  ruleForm: any;
}

export function useDraftHistory(
  postId: Ref<string | number | undefined>,
  ruleForm: Ref<any>,
  isEdit: Ref<boolean>
) {
  const DRAFT_KEY = computed(
    () => `post-draft:${postId.value ? String(postId.value) : "new"}`
  );
  const REVISIONS_KEY = computed(
    () => `post-revisions:${postId.value ? String(postId.value) : "new"}`
  );

  const draftSaving = ref(false);
  const draftSavedAt = ref<number | null>(null);
  const draftRestored = ref(false);
  const historyModalVisible = ref(false);
  const revisionHistory = ref<DraftRevision[]>([]);

  const draftSavedAtText = computed(() =>
    draftSavedAt.value
      ? new Date(draftSavedAt.value).toLocaleTimeString("zh-CN", {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit"
        })
      : ""
  );

  // 加载修订历史快照列表
  const loadRevisions = async () => {
    try {
      const list = await localForage().getItem<DraftRevision[]>(REVISIONS_KEY.value);
      revisionHistory.value = Array.isArray(list) ? list : [];
    } catch (e) {
      console.warn("[draft-history] 加载修订历史失败", e);
      revisionHistory.value = [];
    }
  };

  // 添加快照到修订历史列表（最多保留 25 条）
  const pushRevision = async (form: any) => {
    try {
      const content = form.postContent || "";
      const rev: DraftRevision = {
        id: `rev-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
        timestamp: Date.now(),
        title: form.title || "无标题草稿",
        wordCount: content.length,
        summary: form.summary || (content.slice(0, 100) + (content.length > 100 ? "..." : "")),
        content,
        ruleForm: JSON.parse(JSON.stringify(form))
      };

      const list = await localForage().getItem<DraftRevision[]>(REVISIONS_KEY.value) || [];
      // 避免 30 秒内重复无差异生成快照
      if (list.length > 0) {
        const last = list[0];
        if (last.content === rev.content && last.title === rev.title) {
          return;
        }
      }

      const updated = [rev, ...list].slice(0, 25);
      await localForage().setItem(REVISIONS_KEY.value, updated);
      revisionHistory.value = updated;
    } catch (e) {
      console.warn("[draft-history] 存储修订历史快照失败", e);
    }
  };

  // 立即保存草稿
  const persistDraftImmediately = async (saveRevisionSnapshot = false) => {
    if (!ruleForm.value.title && (!ruleForm.value.postContent || ruleForm.value.postContent.trim().length < 10)) {
      return false;
    }
    draftSaving.value = true;
    try {
      const snapshot = JSON.parse(JSON.stringify(ruleForm.value));
      await localForage().setItem<PostDraft>(DRAFT_KEY.value, {
        ruleForm: snapshot,
        savedAt: Date.now()
      });
      draftSavedAt.value = Date.now();

      if (saveRevisionSnapshot) {
        await pushRevision(snapshot);
      }
      return true;
    } catch (e) {
      console.warn("[post-draft] 保存草稿失败", e);
      return false;
    } finally {
      draftSaving.value = false;
    }
  };

  // 防抖 2s 自动保存
  const persistDraftDebounced = useDebounceFn(() => {
    if (!draftRestored.value) return;
    persistDraftImmediately(false);
  }, 2000);

  // 手动保存草稿（包含生成历史快照）
  const handleSaveDraft = async () => {
    const ok = await persistDraftImmediately(true);
    message(ok ? "草稿与快照已保存到本地" : "内容过少，未保存草稿", {
      type: ok ? "success" : "info"
    });
  };

  // 尝试恢复草稿提示
  const tryRestoreDraft = async (remoteUpdateTime?: number) => {
    let draft: PostDraft | null = null;
    try {
      draft = await localForage().getItem<PostDraft>(DRAFT_KEY.value);
    } catch (e) {
      console.warn("[post-draft] 读取草稿失败", e);
    }

    await loadRevisions();

    if (!draft || !draft.ruleForm) {
      draftRestored.value = true;
      return;
    }

    if (remoteUpdateTime && remoteUpdateTime >= draft.savedAt) {
      await localForage().removeItem(DRAFT_KEY.value);
      draftRestored.value = true;
      return;
    }

    const minutes = Math.max(1, Math.round((Date.now() - draft.savedAt) / 60000));
    try {
      await ElMessageBox.confirm(
        `检测到 ${minutes} 分钟前的未保存本地草稿，是否恢复？选择「丢弃」将清除该草稿。`,
        "恢复草稿",
        {
          confirmButtonText: "恢复草稿",
          cancelButtonText: "丢弃草稿",
          type: "info"
        }
      );

      const d = draft.ruleForm;
      ruleForm.value.title = d.title ?? ruleForm.value.title;
      ruleForm.value.categoryId = d.categoryId ?? ruleForm.value.categoryId;
      ruleForm.value.tags = Array.isArray(d.tags) ? d.tags : ruleForm.value.tags;
      ruleForm.value.summary = d.summary ?? ruleForm.value.summary;
      ruleForm.value.status = d.status ?? ruleForm.value.status;
      ruleForm.value.type = d.type ?? ruleForm.value.type;
      ruleForm.value.postContent = d.postContent ?? ruleForm.value.postContent;
      ruleForm.value.coverImage = d.coverImage ?? ruleForm.value.coverImage;
      ruleForm.value.series = d.series ?? ruleForm.value.series;
      draftSavedAt.value = draft.savedAt;

      message("已成功恢复上次草稿", { type: "success" });
    } catch (action) {
      if (action === "cancel") {
        await localForage().removeItem(DRAFT_KEY.value);
        message("已丢弃草稿", { type: "info" });
      }
    } finally {
      draftRestored.value = true;
    }
  };

  // 清除草稿（发布成功时调用）
  const clearDraft = async () => {
    try {
      await localForage().removeItem(DRAFT_KEY.value);
    } catch (e) {
      console.warn("[post-draft] 清除草稿失败", e);
    }
  };

  const openHistoryModal = async () => {
    await loadRevisions();
    historyModalVisible.value = true;
  };

  const closeHistoryModal = () => {
    historyModalVisible.value = false;
  };

  const restoreRevision = (rev: DraftRevision) => {
    if (!rev?.ruleForm) return;
    const d = rev.ruleForm;
    ruleForm.value.title = d.title ?? ruleForm.value.title;
    ruleForm.value.categoryId = d.categoryId ?? ruleForm.value.categoryId;
    ruleForm.value.tags = Array.isArray(d.tags) ? d.tags : ruleForm.value.tags;
    ruleForm.value.summary = d.summary ?? ruleForm.value.summary;
    ruleForm.value.status = d.status ?? ruleForm.value.status;
    ruleForm.value.type = d.type ?? ruleForm.value.type;
    ruleForm.value.postContent = d.postContent ?? ruleForm.value.postContent;
    ruleForm.value.coverImage = d.coverImage ?? ruleForm.value.coverImage;
    ruleForm.value.series = d.series ?? ruleForm.value.series;
    message(`已恢复至快照版本：${rev.title}`, { type: "success" });
    closeHistoryModal();
  };

  const deleteRevision = async (id: string) => {
    const list = revisionHistory.value.filter(r => r.id !== id);
    revisionHistory.value = list;
    await localForage().setItem(REVISIONS_KEY.value, list);
    message("已删除该版本快照", { type: "info" });
  };

  const clearAllRevisions = async () => {
    try {
      await ElMessageBox.confirm("确定要清空该文章所有的本地历史快照吗？", "清空历史快照", {
        type: "warning"
      });
      revisionHistory.value = [];
      await localForage().removeItem(REVISIONS_KEY.value);
      message("已清空历史快照", { type: "success" });
    } catch {
      // cancel
    }
  };

  return {
    draftSaving,
    draftSavedAt,
    draftSavedAtText,
    draftRestored,
    historyModalVisible,
    revisionHistory,
    persistDraftImmediately,
    persistDraftDebounced,
    handleSaveDraft,
    tryRestoreDraft,
    clearDraft,
    openHistoryModal,
    closeHistoryModal,
    restoreRevision,
    deleteRevision,
    clearAllRevisions,
    pushRevision
  };
}
