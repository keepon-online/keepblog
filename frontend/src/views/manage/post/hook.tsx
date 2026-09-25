import dayjs from "dayjs";
import { message } from "@/utils/message";
import {
  getPostList,
  deletePost,
  updatePostTop,
  updatePostPublish,
  updateCover,
  updateAllCover
} from "@/api/post";
import { getCategoryList } from "@/api/category";
import { ElMessageBox } from "element-plus";
import type { PaginationProps } from "@pureadmin/table";
import { reactive, ref, computed, onMounted, h } from "vue";
import { useRouter } from "vue-router";

export function usePost() {
  const router = useRouter();
  const form = reactive({
    title: undefined,
    categoryId: undefined,
    published: undefined
  });
  const dataList = ref([]);
  const categories = ref([]);
  const selectedRows = ref([]);
  const batchLoading = ref(false);
  const loading = ref(true);
  const switchLoadMap = ref({});
  const pagination = reactive<PaginationProps>({
    total: 0,
    pageSize: 10,
    pageSizes: [10, 15, 20, 50, 100],
    currentPage: 1,
    background: true
  });

  const columns: TableColumnList = [
    {
      type: "selection",
      width: 50,
      align: "left"
    },
    {
      label: "序号",
      type: "index",
      width: 60
    },
    {
      label: "文章详情",
      prop: "title",
      minWidth: 280,
      align: "left",
      cellRenderer: ({ row }) => (
        <div class="flex items-center gap-3 py-1">
          <div class="relative flex-shrink-0" style="width: 54px; height: 38px;">
            {row.coverImage ? (
              <el-image
                src={row.coverImage}
                fit="cover"
                loading="lazy"
                preview-teleported
                preview-src-list={[row.coverImage]}
                style="width: 54px; height: 38px; border-radius: 6px; border: 1px solid var(--el-border-color-lighter);"
              />
            ) : (
              <div
                style="width: 54px; height: 38px; border-radius: 6px; background: linear-gradient(135deg, rgba(64,158,255,0.15), rgba(54,163,247,0.3)); display: flex; align-items: center; justify-content: center; color: var(--el-color-primary); font-size: 16px;"
              >
                📄
              </div>
            )}
            {row.top === 1 ? (
              <span
                style="position: absolute; top: -4px; left: -4px; background: #e6a23c; color: #fff; font-size: 10px; font-weight: bold; padding: 1px 4px; border-radius: 4px; box-shadow: 0 1px 4px rgba(0,0,0,0.2);"
              >
                TOP
              </span>
            ) : null}
          </div>

          <div class="flex flex-col min-w-0 flex-1">
            <div
              class="font-semibold text-sm cursor-pointer hover:text-primary transition-colors duration-200 truncate"
              style="color: var(--el-text-color-primary);"
              title={row.title}
              onClick={() => handleUpdate(row)}
            >
              {row.title}
            </div>
            <div class="flex items-center gap-3 text-xs mt-1" style="color: var(--el-text-color-placeholder);">
              {row.series ? (
                <span class="px-1.5 py-0.5 rounded text-[11px]" style="background: rgba(103,194,58,0.1); color: #67c23a;">
                  📚 {row.series}
                </span>
              ) : null}
              <span>👁️ {row.readCount ?? 0} 阅读</span>
              <span>📝 {row.wordCount ?? 0} 字</span>
            </div>
          </div>
        </div>
      )
    },
    {
      label: "所属分类",
      prop: "categoryName",
      width: 130,
      cellRenderer: ({ row }) => (
        <el-tag size="small" effect="plain" type="primary">
          {row.categoryName || "未分类"}
        </el-tag>
      )
    },
    {
      label: "发布状态",
      width: 110,
      cellRenderer: scope => (
        <el-switch
          size="default"
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.published}
          active-value={1}
          inactive-value={0}
          active-text="公开"
          inactive-text="草稿"
          inline-prompt
          onChange={() => onPublishChange(scope as any)}
        />
      )
    },
    {
      label: "置顶状态",
      width: 110,
      cellRenderer: scope => (
        <el-switch
          size="default"
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.top}
          active-value={1}
          inactive-value={0}
          active-text="置顶"
          inactive-text="普通"
          inline-prompt
          onChange={() => onTopChange(scope as any)}
        />
      )
    },
    {
      label: "发布时间",
      width: 160,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs.unix(createTime).format("YYYY-MM-DD HH:mm") : "-"
    },
    {
      label: "操作",
      fixed: "right",
      width: 180,
      slot: "operation"
    }
  ];

  const buttonClass = computed(() => {
    return [
      "!h-[26px]",
      "reset-margin",
      "!text-gray-500",
      "dark:!text-white",
      "dark:hover:!text-primary"
    ];
  });

  function onPublishChange({ row }) {
    const targetStatus = row.published;
    const actionText = targetStatus === 1 ? "公开发布" : "设为草稿";
    ElMessageBox.confirm(
      `确认要将文章【<strong style='color:var(--el-color-primary)'>${row.title}</strong>】${actionText}吗?`,
      "状态变更提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true,
        draggable: true
      }
    )
      .then(() => {
        updatePostPublish({
          postId: row.postId,
          published: targetStatus
        }).then(r => {
          if (r.code === 200) {
            message(`已成功${actionText}`, { type: "success" });
          } else {
            message(r.message, { type: "error" });
          }
        });
      })
      .catch(() => {
        row.published = targetStatus === 1 ? 0 : 1;
      });
  }

  function onTopChange({ row }) {
    const targetTop = row.top;
    const actionText = targetTop === 1 ? "置顶" : "取消置顶";
    ElMessageBox.confirm(
      `确认要将文章【<strong style='color:var(--el-color-primary)'>${row.title}</strong>】${actionText}吗?`,
      "置顶调整提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true,
        draggable: true
      }
    )
      .then(() => {
        updatePostTop({
          postId: row.postId,
          top: targetTop
        }).then(r => {
          if (r.code === 200) {
            message(`已成功${actionText}`, { type: "success" });
          } else {
            message(r.message, { type: "error" });
          }
        });
      })
      .catch(() => {
        row.top = targetTop === 1 ? 0 : 1;
      });
  }

  async function handleUpdate(row) {
    const id = row.postSlug || row.postId;
    if (!id) {
      message("无法获取文章标识", { type: "warning" });
      return;
    }
    await router.push({
      name: "内容编辑",
      params: { id: String(id) }
    });
  }

  async function handleUpdateCoverImage(row) {
    const id = row.postSlug || row.postId;
    if (id) {
      await updateCover(String(id));
      await onSearch();
    }
  }

  async function handleUpdateAllCoverImage() {
    await updateAllCover();
    await onSearch();
  }

  async function handleDelete(row) {
    const id = row.postSlug || row.postId;
    if (id) {
      await deletePost(String(id));
      message("删除成功", { type: "success" });
      await onSearch();
    }
  }

  // 批量设为发布 / 草稿
  const handleBatchPublish = async (status: number) => {
    if (selectedRows.value.length === 0) return;
    const text = status === 1 ? "发布" : "下架为草稿";
    try {
      await ElMessageBox.confirm(
        `确认将选中的 ${selectedRows.value.length} 篇文章批量${text}吗?`,
        "批量操作确认",
        { type: "warning" }
      );
      batchLoading.value = true;
      for (const item of selectedRows.value) {
        await updatePostPublish({ postId: item.postId, published: status });
      }
      message(`已成功批量${text}`, { type: "success" });
      await onSearch();
    } catch {
      // cancel
    } finally {
      batchLoading.value = false;
    }
  };

  // 批量删除
  const handleBatchDelete = async () => {
    if (selectedRows.value.length === 0) return;
    try {
      await ElMessageBox.confirm(
        `确定要永久删除选中的 ${selectedRows.value.length} 篇文章吗？此操作不可逆！`,
        "批量删除警示",
        {
          type: "warning",
          confirmButtonText: "确认删除",
          cancelButtonText: "取消"
        }
      );
      batchLoading.value = true;
      for (const item of selectedRows.value) {
        const id = item.postSlug || item.postId;
        if (id) await deletePost(String(id));
      }
      message("批量删除完成", { type: "success" });
      await onSearch();
    } catch {
      // cancel
    } finally {
      batchLoading.value = false;
    }
  };

  const handleSelectionChange = (rows: any) => {
    selectedRows.value = rows;
  };

  async function handleCurrentChange(val: number) {
    loading.value = true;
    pagination.currentPage = val;
    const params = {
      pageNum: pagination.currentPage,
      pageSize: pagination.pageSize,
      title: form.title,
      categoryId: form.categoryId,
      published: form.published
    };
    try {
      const { payload } = await getPostList(params);
      dataList.value = payload.list || [];
      pagination.total = payload.total || 0;
    } finally {
      loading.value = false;
    }
  }

  function onSizeChange(val: number) {
    pagination.pageSize = val;
    onSearch();
  }

  async function onSearch() {
    loading.value = true;
    const params = {
      pageNum: 1,
      pageSize: pagination.pageSize,
      title: form.title,
      categoryId: form.categoryId,
      published: form.published
    };
    try {
      const { payload } = await getPostList(params);
      dataList.value = payload.list || [];
      pagination.total = payload.total || 0;
    } finally {
      loading.value = false;
    }
  }

  const resetForm = formEl => {
    if (!formEl) return;
    formEl.resetFields();
    onSearch();
  };

  onMounted(() => {
    onSearch();
    getCategoryList().then(res => {
      categories.value = res.payload || [];
    });
  });

  return {
    form,
    loading,
    columns,
    dataList,
    pagination,
    buttonClass,
    categories,
    selectedRows,
    batchLoading,
    handleSelectionChange,
    handleBatchPublish,
    handleBatchDelete,
    onSearch,
    resetForm,
    handleUpdate,
    handleDelete,
    onSizeChange,
    handleUpdateAllCoverImage,
    handleUpdateCoverImage,
    handleCurrentChange
  };
}
