import dayjs from "dayjs";
import { message } from "@/utils/message";
import {
  getMusicList,
  getMusic,
  addMusic,
  updateMusic,
  deleteMusic,
  changeMusicState
} from "@/api/music";
import { ElMessageBox, type FormInstance, type FormRules } from "element-plus";
import { reactive, ref, onMounted } from "vue";

export function useMusic() {
  const ruleFormRef = ref<FormInstance>();
  const form = ref({
    id: undefined as number | undefined,
    name: "",
    artist: "",
    url: "",
    cover: "",
    lrc: "",
    sort: 0,
    state: 1
  });

  const rules = reactive<FormRules>({
    name: [
      { required: true, message: "请输入歌曲名称", trigger: "blur" },
      { min: 1, max: 100, message: "长度在1-100", trigger: "blur" }
    ],
    artist: [
      { required: true, message: "请输入艺术家", trigger: "blur" },
      { min: 1, max: 100, message: "长度在1-100", trigger: "blur" }
    ],
    url: [
      { required: true, message: "请输入音乐URL", trigger: "blur" }
    ]
  });

  const dataList = ref([]);
  const title = ref("");
  const dialogFormVisible = ref(false);
  const loading = ref(true);
  const switchLoadMap = ref({});

  const resetForm = () => {
    form.value.id = undefined;
    form.value.name = "";
    form.value.artist = "";
    form.value.url = "";
    form.value.cover = "";
    form.value.lrc = "";
    form.value.sort = 0;
    form.value.state = 1;
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
      label: "歌曲名称",
      prop: "name",
      minWidth: 120
    },
    {
      label: "艺术家",
      prop: "artist",
      minWidth: 100
    },
    {
      label: "排序",
      prop: "sort",
      width: 80
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
          active-text="启用"
          inactive-text="禁用"
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
      `确认要<strong>${row.state === 0 ? "禁用" : "启用"
      }</strong><strong style='color:var(--el-color-primary)'>${row.name
      }</strong>吗?`,
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
        changeMusicState({ id: row.id, state: row.state }).then(r => {
          if (r.code === 200) {
            message("已成功修改状态", { type: "success" });
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
    title.value = "编辑音乐";
    const { payload } = await getMusic(row.id);
    form.value = payload;
  }

  const handleAdd = async () => {
    resetForm();
    dialogFormVisible.value = true;
    title.value = "添加音乐";
  };

  function handleDelete(row) {
    deleteMusic(row.id).then(r => {
      if (r.code === 200) {
        message("删除成功", { type: "success" });
        onSearch();
      } else {
        message(r.message, { type: "error" });
      }
    });
  }

  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    await formEl.validate((valid, fields) => {
      if (valid) {
        if (form.value.id != undefined) {
          updateMusic(form.value).then(r => {
            if (r.code === 200) {
              message("更新成功", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        } else {
          addMusic(form.value).then(r => {
            if (r.code === 200) {
              message("添加成功", { type: "success" });
            } else {
              message(r.message, { type: "error" });
            }
            dialogFormVisible.value = false;
            onSearch();
          });
        }
      } else {
        console.log("error submit!", fields);
      }
    });
  };

  async function onSearch() {
    loading.value = true;
    const { payload } = await getMusicList();
    dataList.value = payload || [];
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
