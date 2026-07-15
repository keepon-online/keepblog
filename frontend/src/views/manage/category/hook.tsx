import dayjs from "dayjs";
import { message } from "@/utils/message";
import {
  getCategoryList,
  getCategory,
  addCategory,
  updateCategory,
  deleteCategory,
  changeCategoryState
} from "@/api/category";
import { ElMessageBox, type FormInstance, type FormRules } from "element-plus";
import { reactive, ref, onMounted } from "vue";

export function useCategory() {
  const ruleFormRef = ref<FormInstance>();
  const form = ref({
    categoryId: undefined,
    categoryName: "",
    note: ""
  });
  const rules = reactive<FormRules>({
    categoryName: [
      { required: true, message: "请输入分类名称", trigger: "blur" },
      { min: 2, max: 20, message: "长度在2-20", trigger: "blur" }
    ],
    note: [
      { required: true, message: "请输入分类描述", trigger: "blur" },
      { min: 5, max: 200, message: "长度在5-200", trigger: "blur" }
    ]
  });
  const dataList = ref([]);
  const title = ref("");
  const dialogFormVisible = ref(false);

  const loading = ref(true);
  const switchLoadMap = ref({});
  const resetForm = () => {
    form.value.categoryId = undefined;
    form.value.categoryName = "";
    form.value.note = "";
  };
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
      label: "分类名称",
      prop: "categoryName",
      minWidth: 120
    },
    {
      label: "分类描述",
      prop: "note",
      minWidth: 100
    },
    {
      label: "状态",
      minWidth: 130,
      cellRenderer: scope => (
        <el-switch
          size={scope.props.size === "small" ? "small" : "default"}
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.state}
          active-value={1}
          inactive-value={0}
          active-text="已开启"
          inactive-text="已关闭"
          inline-prompt
          onChange={() => onChange(scope as any)}
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

  function onChange({ row }) {
    ElMessageBox.confirm(
      `确认要<strong>${
        row.state === 0 ? "停用" : "启用"
      }</strong><strong style='color:var(--el-color-primary)'>${
        row.categoryName
      }</strong>分类吗?`,
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
        changeCategoryState(row).then(r => {
          if (r.code === 200) {
            message("已成功修改分类状态", { type: "success" });
          } else {
            message(r.message, { type: "error" });
          }
        });
      })
      .catch(() => {
        row.state === 0 ? (row.state = 1) : (row.state = 0);
      });
  }

  async function handleUpdate(row) {
    dialogFormVisible.value = true;
    title.value = "更新分类";
    const { payload } = await getCategory(row.categoryId);
    form.value = payload;
  }

  const handleAdd = async () => {
    resetForm();
    dialogFormVisible.value = true;
    title.value = "添加分类";
  };

  function handleDelete(row) {
    deleteCategory(row.categoryId).then(r => {
      if (r.code === 200) {
        message("success", { type: "success" });
      } else {
        message(r.message, { type: "error" });
      }
    });
    onSearch();
  }

  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    await formEl.validate((valid, fields) => {
      if (valid) {
        if (form.value.categoryId != undefined) {
          updateCategory(form.value).then(r => {
            if (r.code === 200) {
              message("success", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        } else {
          addCategory(form.value).then(r => {
            if (r.code === 200) {
              message("success", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
          });
          dialogFormVisible.value = false;
          onSearch();
        }
      }
    });
  };

  async function onSearch() {
    loading.value = true;
    const { payload } = await getCategoryList();
    dataList.value = payload;
    setTimeout(() => {
      loading.value = false;
    }, 500);
  }

  onMounted(() => {
    onSearch();
  });

  return {
    dialogFormVisible,
    title,
    ruleFormRef,
    rules,
    form,
    loading,
    columns,
    dataList,
    handleAdd,
    submitForm,
    onSearch,
    handleUpdate,
    handleDelete
  };
}
