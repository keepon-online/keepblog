import dayjs from "dayjs";
import { message } from "@/utils/message";
import {
  getLinkList,
  getLink,
  addLink,
  updateLink,
  deleteLink,
  changeLinkState,
  checkLink
} from "@/api/link";
import { ElMessageBox, type FormInstance, type FormRules } from "element-plus";
import { reactive, ref, onMounted, h } from "vue";

export function useLink() {
  const ruleFormRef = ref<FormInstance>();
  const form = ref({
    id: undefined,
    title: "",
    linkUrl: "",
    linkIcon: "",
    linkDesc: "",
    type: 0
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
      { min: 2, max: 200, message: "长度在2-200", trigger: "blur" }
    ]
  });

  const dataList = ref([]);
  const title = ref("");
  const dialogFormVisible = ref(false);

  const loading = ref(true);
  const switchLoadMap = ref({});
  const checkLoadingMap = ref<Record<number, boolean>>({});
  const healthStatusMap = ref<Record<number, { status: number; msg: string }>>({});
  const checkingAll = ref(false);

  const resetForm = () => {
    form.value.id = undefined;
    form.value.title = "";
    form.value.linkUrl = "";
    form.value.linkIcon = "";
    form.value.linkDesc = "";
    form.value.type = 0;
  };

  const copyUrl = (url: string) => {
    navigator.clipboard.writeText(url);
    message("链接已复制到剪贴板", { type: "success" });
  };

  const handleCheckLink = async (row: any) => {
    if (!row.linkUrl) return;
    checkLoadingMap.value[row.id] = true;
    try {
      const res = await checkLink(row.linkUrl);
      if (res.code === 200 && res.payload) {
        healthStatusMap.value[row.id] = res.payload;
        if (res.payload.status === 200) {
          message(`【${row.title}】连接正常 (200 OK)`, { type: "success" });
        } else {
          message(`【${row.title}】检测状态: ${res.payload.msg}`, { type: "warning" });
        }
      }
    } catch {
      healthStatusMap.value[row.id] = { status: 0, msg: "网络错误" };
    } finally {
      checkLoadingMap.value[row.id] = false;
    }
  };

  const handleCheckAll = async () => {
    if (dataList.value.length === 0) return;
    checkingAll.value = true;
    message("正在检测所有友链连通性...", { type: "info" });
    for (const item of dataList.value) {
      await handleCheckLink(item);
    }
    checkingAll.value = false;
    message("所有友链体检完成", { type: "success" });
  };

  const columns: TableColumnList = [
    {
      label: "序号",
      type: "index",
      width: 65
    },
    {
      label: "站点标识",
      prop: "linkIcon",
      width: 100,
      cellRenderer: ({ row }) => {
        return (
          <div class="flex items-center justify-center">
            {row.linkIcon ? (
              <el-image
                style="width: 38px; height: 38px; border-radius: 8px; border: 1px solid var(--el-border-color-lighter);"
                src={row.linkIcon}
                fit="cover"
                loading="lazy"
              >
                {{
                  error: () => (
                    <div style="width: 38px; height: 38px; border-radius: 8px; background: linear-gradient(135deg, #409eff, #36a3f7); color: #fff; font-weight: bold; display: flex; align-items: center; justify-content: center; font-size: 14px;">
                      {(row.title || "L").slice(0, 1).toUpperCase()}
                    </div>
                  )
                }}
              </el-image>
            ) : (
              <div style="width: 38px; height: 38px; border-radius: 8px; background: linear-gradient(135deg, #409eff, #36a3f7); color: #fff; font-weight: bold; display: flex; align-items: center; justify-content: center; font-size: 14px;">
                {(row.title || "L").slice(0, 1).toUpperCase()}
              </div>
            )}
          </div>
        );
      }
    },
    {
      label: "网站名称与链接",
      prop: "title",
      minWidth: 220,
      align: "left",
      cellRenderer: ({ row }) => (
        <div class="flex flex-col gap-1 py-1">
          <div class="flex items-center gap-2">
            <span style="font-weight: 600; color: var(--el-text-color-primary); font-size: 14px;">
              {row.title}
            </span>
            <el-tag
              size="small"
              type={row.type === 1 ? "warning" : "primary"}
              effect="plain"
            >
              {row.type === 0 ? "网站" : "友链"}
            </el-tag>
          </div>
          <div class="flex items-center gap-2 text-xs" style="color: var(--el-text-color-secondary);">
            <a
              href={row.linkUrl}
              target="_blank"
              rel="noopener noreferrer"
              style="color: var(--el-color-primary); text-decoration: none; max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
              class="hover:underline"
            >
              {row.linkUrl}
            </a>
            <el-tooltip content="复制网址" placement="top">
              <span
                class="cursor-pointer hover:text-primary"
                onClick={() => copyUrl(row.linkUrl)}
                style="color: var(--el-text-color-placeholder); font-size: 11px;"
              >
                📋
              </span>
            </el-tooltip>
          </div>
        </div>
      )
    },
    {
      label: "网站描述",
      prop: "linkDesc",
      minWidth: 160,
      showOverflowTooltip: true,
      cellRenderer: ({ row }) => (
        <span style="color: var(--el-text-color-regular); font-size: 13px;">
          {row.linkDesc || "-"}
        </span>
      )
    },
    {
      label: "存活状态",
      width: 140,
      cellRenderer: ({ row }) => {
        const isChecking = checkLoadingMap.value[row.id];
        const health = healthStatusMap.value[row.id];

        if (isChecking) {
          return <el-tag size="small" type="info">⏳ 探测中...</el-tag>;
        }

        if (health) {
          if (health.status === 200) {
            return (
              <el-tag size="small" type="success" effect="light">
                🟢 正常 (200)
              </el-tag>
            );
          }
          if (health.status >= 300 && health.status < 400) {
            return (
              <el-tag size="small" type="warning" effect="light">
                🟡 重定向 ({health.status})
              </el-tag>
            );
          }
          return (
            <el-tag size="small" type="danger" effect="light" title={health.msg}>
              🔴 异常 ({health.status || "超时"})
            </el-tag>
          );
        }

        return (
          <el-button
            link
            type="primary"
            size="small"
            onClick={() => handleCheckLink(row)}
          >
            🔍 检测存活
          </el-button>
        );
      }
    },
    {
      label: "展示状态",
      width: 110,
      cellRenderer: scope => (
        <el-switch
          size="default"
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.state}
          active-value={1}
          inactive-value={0}
          active-text="展示"
          inactive-text="隐藏"
          inline-prompt
          onChange={() => onChange(scope as any)}
        />
      )
    },
    {
      label: "创建时间",
      width: 160,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs.unix(createTime).format("YYYY-MM-DD HH:mm") : "-"
    },
    {
      label: "操作",
      fixed: "right",
      width: 170,
      slot: "operation"
    }
  ];

  function onChange({ row }) {
    ElMessageBox.confirm(
      `确认要<strong>${
        row.state === 0 ? "隐藏" : "展示"
      }</strong>【<strong style='color:var(--el-color-primary)'>${
        row.title
      }</strong>】友链吗?`,
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
            message("已成功修改友链状态", { type: "success" });
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
    title.value = "编辑友链信息";
    const { payload } = await getLink(row.id);
    form.value = payload;
  }

  const handleAdd = async () => {
    resetForm();
    dialogFormVisible.value = true;
    title.value = "添加新友链";
  };

  function handleDelete(row) {
    deleteLink(row.id).then(r => {
      if (r.code === 200) {
        message("已成功删除该友链", { type: "success" });
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
        if (form.value.id !== undefined) {
          updateLink(form.value).then(r => {
            if (r.code === 200) {
              message("更新成功", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        } else {
          addLink(form.value).then(r => {
            if (r.code === 200) {
              message("添加成功", { type: "success" });
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
    const { payload } = await getLinkList();
    dataList.value = payload || [];
    setTimeout(() => {
      loading.value = false;
    }, 300);
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
    checkingAll,
    handleAdd,
    handleCheckLink,
    handleCheckAll,
    submitForm,
    onSearch,
    handleUpdate,
    handleDelete
  };
}
