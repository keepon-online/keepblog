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
import { reactive, ref, computed, onMounted } from "vue";
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
      label: "多选",
      width: 55,
      align: "left"
    },
    {
      label: "序号",
      type: "index",
      width: 70
    },
    {
      label: "标题",
      prop: "title",
      minWidth: 100
    },
    {
      label: "分类名称",
      prop: "categoryName",
      minWidth: 120
    },
    {
      label: "状态",
      minWidth: 130,
      cellRenderer: scope => (
        <el-switch
          size={scope.props.size === "small" ? "small" : "default"}
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.published}
          active-value={1}
          inactive-value={0}
          active-text="已发布"
          inactive-text="未发布"
          inline-prompt
          onChange={() => onPublishChange(scope as any)}
        />
      )
    },
    {
      label: "置顶",
      minWidth: 130,
      cellRenderer: scope => (
        <el-switch
          size={scope.props.size === "small" ? "small" : "default"}
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.top}
          active-value={1}
          inactive-value={0}
          active-text="置顶"
          inactive-text="正常"
          inline-prompt
          onChange={() => onTopChange(scope as any)}
        />
      )
    },
    {
      label: "创建时间",
      minWidth: 180,
      prop: "createTime",
      formatter: ({ createTime }) =>
        dayjs.unix(createTime).format("YYYY-MM-DD HH:mm:ss")
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
      "!h-[20px]",
      "reset-margin",
      "!text-gray-500",
      "dark:!text-white",
      "dark:hover:!text-primary"
    ];
  });

  function onPublishChange({ row }) {
    ElMessageBox.confirm(
      `确认要<strong>${
        row.published === 0 ? "停用" : "启用"
      }</strong><strong style='color:var(--el-color-primary)'>${
        row.title
      }</strong>文章吗?`,
      "系统提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true,
        draggable: true
      }
    )
      .then(() => {
        setTimeout(() => {
          const obj = {
            postId: row.postId,
            published: row.published
          };
          updatePostPublish(obj).then(r => {
            if (r.code === 200) {
              message(`已成功修改文章状态`, {
                type: "success"
              });
            }
          });
        }, 300);
      })
      .catch(() => {
        row.published === 0 ? (row.published = 1) : (row.published = 0);
      });
  }
  function onTopChange({ row }) {
    ElMessageBox.confirm(
      `确认要<strong>${
        row.top === 1 ? "置顶" : "启用"
      }</strong><strong style='color:var(--el-color-primary)'>${
        row.title
      }</strong>文章吗?`,
      "系统提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true,
        draggable: true
      }
    )
      .then(() => {
        setTimeout(() => {
          const obj = {
            postId: row.postId,
            top: row.top
          };
          updatePostTop(obj).then(r => {
            if (r.code === 200) {
              message(`已成功置顶文章`, {
                type: "success"
              });
            }
          });
        }, 300);
      })
      .catch(() => {
        row.top === 0 ? (row.top = 1) : (row.top = 0);
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
    if (id) await updateCover(String(id));
  }
  async function handleUpdateAllCoverImage() {
    await updateAllCover();
  }

  async function handleDelete(row) {
    const id = row.postSlug || row.postId;
    if (id) {
      await deletePost(String(id));
      await onSearch();
    }
  }

  async function handleCurrentChange(val) {
    loading.value = true;
    pagination.currentPage = val;
    const params = {
      pageNum: pagination.currentPage,
      pageSize: pagination.pageSize
    };
    const { payload } = await getPostList(params);
    dataList.value = payload.list;
    pagination.total = payload.total;
    setTimeout(() => {
      loading.value = false;
    }, 500);
  }
  function onSizeChange(val) {
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
    const { payload } = await getPostList(params);
    dataList.value = payload.list;
    pagination.total = payload.total;
    setTimeout(() => {
      loading.value = false;
    }, 500);
  }

  const resetForm = formEl => {
    if (!formEl) return;
    formEl.resetFields();
    onSearch();
  };

  onMounted(() => {
    onSearch();
    getCategoryList().then(res => {
      categories.value = res.payload;
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
