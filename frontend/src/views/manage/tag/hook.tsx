import dayjs from "dayjs";
import { message } from "@/utils/message";
import {
  getTagManageList,
  getTag,
  addTag,
  updateTag,
  deleteTag
} from "@/api/tag";
import type { FormInstance, FormRules } from "element-plus";
import { reactive, ref, onMounted } from "vue";

export function useTag() {
  const ruleFormRef = ref<FormInstance>();
  const form = ref({
    tagId: undefined,
    tagName: ""
  });
  const rules = reactive<FormRules>({
    tagName: [
      { required: true, message: "请输入标签名称", trigger: "blur" },
      { min: 1, max: 30, message: "长度在1-30", trigger: "blur" }
    ]
  });
  const dataList = ref([]);
  const title = ref("");
  const dialogFormVisible = ref(false);

  const loading = ref(true);
  const resetForm = () => {
    form.value.tagId = undefined;
    form.value.tagName = "";
  };
  const tagTypes = ["primary", "success", "warning", "danger", "info"] as const;
  const getTagType = (name: string) => {
    let hash = 0;
    for (let i = 0; i < (name || "").length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    return tagTypes[Math.abs(hash) % tagTypes.length];
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
      label: "标签名称",
      prop: "tagName",
      minWidth: 140,
      cellRenderer: ({ row }) => (
        <el-tag size="default" effect="light" type={getTagType(row.tagName)}>
          #{row.tagName}
        </el-tag>
      )
    },
    {
      label: "关联文章数",
      prop: "postCount",
      width: 130,
      cellRenderer: ({ row }) => (
        <span style="font-weight: 600; color: var(--el-text-color-regular);">
          {row.postCount ?? 0} 篇
        </span>
      )
    },
    {
      label: "创建时间",
      minWidth: 180,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs.unix(createTime).format("YYYY-MM-DD HH:mm:ss") : "-"
    },
    {
      label: "操作",
      fixed: "right",
      width: 180,
      slot: "operation"
    }
  ];

  async function handleUpdate(row) {
    dialogFormVisible.value = true;
    title.value = "更新标签";
    const { payload } = await getTag(row.tagId);
    form.value = payload;
  }

  const handleAdd = async () => {
    resetForm();
    dialogFormVisible.value = true;
    title.value = "添加标签";
  };

  function handleDelete(row) {
    deleteTag(row.tagId).then(r => {
      if (r.code === 200) {
        message("删除成功", { type: "success" });
      } else {
        message(r.message, { type: "error" });
      }
    });
    onSearch();
  }

  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    await formEl.validate((valid, _fields) => {
      if (valid) {
        if (form.value.tagId != undefined) {
          updateTag(form.value).then(r => {
            if (r.code === 200) {
              message("更新成功", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        } else {
          addTag(form.value).then(r => {
            if (r.code === 200) {
              message("新增成功", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        }
      }
    });
  };

  async function onSearch() {
    loading.value = true;
    const { payload } = await getTagManageList();
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
