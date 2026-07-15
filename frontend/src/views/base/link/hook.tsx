import dayjs from "dayjs";
import { message } from "@/utils/message";
import {
  getLinkList,
  getLink,
  addLink,
  updateLink,
  deleteLink,
  changeLinkState
} from "@/api/link";
import { ElMessageBox, type FormInstance, type FormRules } from "element-plus";
import { reactive, ref, onMounted } from "vue";

export function useLink() {
  const ruleFormRef = ref<FormInstance>();
  const form = ref({
    id: undefined,
    title: "",
    linkUrl: "",
    linkIcon: "",
    linkDesc: ""
  });
  const rules = reactive<FormRules>({
    title: [
      { required: true, message: "请输入名称", trigger: "blur" },
      { min: 2, max: 20, message: "长度在2-20", trigger: "blur" }
    ],
    linkUrl: [
      { required: true, message: "请输入链接", trigger: "blur" },
      { min: 2, max: 120, message: "长度在2-120", trigger: "blur" }
    ],
    linkDesc: [
      { required: true, message: "请输入描述", trigger: "blur" },
      { min: 5, max: 200, message: "长度在5-200", trigger: "blur" }
    ]
  });
  const dataList = ref([]);
  const title = ref("");
  const dialogFormVisible = ref(false);

  const loading = ref(true);
  const switchLoadMap = ref({});
  const resetForm = () => {
    form.value.id = undefined;
    form.value.title = "";
    form.value.linkUrl = "";
    form.value.linkIcon = "";
    form.value.linkDesc = "";
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
      label: "LOGO",
      prop: "linkIcon",
      minWidth: 120,
      cellRenderer: ({ row }) => (
        <el-image
          style="width: 100px; height: 100px"
          src={row.linkIcon}
          fit={"fit"}
        />
      )
    },
    {
      label: "名称",
      prop: "title",
      minWidth: 120
    },
    {
      label: "链接",
      prop: "linkUrl",
      minWidth: 120,
      cellRenderer: ({ row }) => (
        <a href={row.linkUrl} target="_black">
          {row.linkUrl}{" "}
        </a>
      )
    },
    {
      label: "类型",
      prop: "type",
      minWidth: 120,
      cellRenderer: ({ row, props }) => (
        <el-tag
          size={props.size}
          type={row.type === 1 ? "danger" : ""}
          effect="plain"
        >
          {row.type === 0 ? "网站" : "友链"}
        </el-tag>
      )
    },
    {
      label: "描述",
      prop: "linkDesc",
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
        row.title
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
        changeLinkState(row).then(r => {
          if (r.code === 200) {
            message("已成功修改链接状态", { type: "success" });
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
    title.value = "更新友链";
    const { payload } = await getLink(row.id);
    form.value = payload;
  }

  const handleAdd = async () => {
    resetForm();
    dialogFormVisible.value = true;
    title.value = "添加友链";
  };

  function handleDelete(row) {
    deleteLink(row.id).then(r => {
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
        if (form.value.id != undefined) {
          updateLink(form.value).then(r => {
            if (r.code === 200) {
              message("success", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        } else {
          addLink(form.value).then(r => {
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
    const { payload } = await getLinkList();
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
