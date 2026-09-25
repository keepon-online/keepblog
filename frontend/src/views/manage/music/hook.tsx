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
import { reactive, ref, onMounted, onUnmounted, h } from "vue";

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
    url: [{ required: true, message: "请输入音乐URL", trigger: "blur" }]
  });

  const dataList = ref([]);
  const title = ref("");
  const dialogFormVisible = ref(false);
  const loading = ref(true);
  const switchLoadMap = ref({});

  // 在线试听播放器状态
  const currentPlaying = ref<any>(null);
  const isPlaying = ref(false);
  let audio: HTMLAudioElement | null = null;

  const initAudio = () => {
    if (!audio && typeof Audio !== "undefined") {
      audio = new Audio();
      audio.addEventListener("ended", () => {
        isPlaying.value = false;
      });
      audio.addEventListener("error", () => {
        isPlaying.value = false;
        message("音频加载失败，请检查直链有效性", { type: "error" });
      });
    }
  };

  const togglePlay = (row: any) => {
    if (!row.url) {
      message("未配置音频播放链接", { type: "warning" });
      return;
    }
    initAudio();
    if (!audio) return;

    if (currentPlaying.value?.id === row.id) {
      if (isPlaying.value) {
        audio.pause();
        isPlaying.value = false;
      } else {
        audio.play().then(() => {
          isPlaying.value = true;
        }).catch(() => {
          message("播放失败，请检查浏览器权限或直链", { type: "error" });
        });
      }
    } else {
      audio.src = row.url;
      currentPlaying.value = row;
      audio.play().then(() => {
        isPlaying.value = true;
      }).catch(() => {
        message("播放失败，请检查音频地址", { type: "error" });
      });
    }
  };

  const stopPlay = () => {
    if (audio) {
      audio.pause();
      audio.currentTime = 0;
    }
    isPlaying.value = false;
    currentPlaying.value = null;
  };

  const copyUrl = (url: string) => {
    navigator.clipboard.writeText(url);
    message("音频地址已复制", { type: "success" });
  };

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
      label: "序号",
      type: "index",
      width: 60
    },
    {
      label: "封面与试听",
      prop: "cover",
      width: 105,
      cellRenderer: ({ row }) => {
        const isCurrent = currentPlaying.value?.id === row.id;
        const active = isCurrent && isPlaying.value;

        return (
          <div
            class="relative cursor-pointer group flex items-center justify-center"
            style="width: 44px; height: 44px; margin: 0 auto;"
            onClick={() => togglePlay(row)}
          >
            {row.cover ? (
              <el-image
                src={row.cover}
                fit="cover"
                loading="lazy"
                style="width: 44px; height: 44px; border-radius: 8px; border: 1px solid var(--el-border-color-lighter);"
              />
            ) : (
              <div
                style="width: 44px; height: 44px; border-radius: 8px; background: linear-gradient(135deg, #9b59b6, #8e44ad); display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px;"
              >
                🎵
              </div>
            )}
            <div
              style={{
                position: "absolute",
                inset: 0,
                borderRadius: "8px",
                backgroundColor: active ? "rgba(0,0,0,0.4)" : "rgba(0,0,0,0.25)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#fff",
                fontSize: "16px",
                transition: "all 0.2s"
              }}
            >
              {active ? "⏸️" : "▶️"}
            </div>
          </div>
        );
      }
    },
    {
      label: "歌曲与歌手",
      prop: "name",
      minWidth: 200,
      align: "left",
      cellRenderer: ({ row }) => (
        <div class="flex flex-col gap-1 py-1">
          <div class="flex items-center gap-2">
            <span style="font-weight: 600; color: var(--el-text-color-primary); font-size: 14px;">
              {row.name}
            </span>
            {currentPlaying.value?.id === row.id && isPlaying.value ? (
              <span class="text-xs px-1.5 py-0.5 rounded bg-green-500/10 text-green-500 animate-pulse">
                正在播放
              </span>
            ) : null}
          </div>
          <div class="flex items-center gap-2 text-xs" style="color: var(--el-text-color-secondary);">
            <span>🎤 {row.artist}</span>
            <el-tooltip content="复制音频直链" placement="top">
              <span
                class="cursor-pointer hover:text-primary"
                onClick={() => copyUrl(row.url)}
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
      label: "音频直链",
      prop: "url",
      minWidth: 180,
      showOverflowTooltip: true,
      cellRenderer: ({ row }) => (
        <a
          href={row.url}
          target="_blank"
          rel="noopener noreferrer"
          style="color: var(--el-color-primary); font-size: 12px; font-family: monospace;"
        >
          {row.url}
        </a>
      )
    },
    {
      label: "歌词状态",
      width: 100,
      cellRenderer: ({ row }) => (
        <el-tag size="small" type={row.lrc ? "success" : "info"} effect="plain">
          {row.lrc ? "有歌词" : "无歌词"}
        </el-tag>
      )
    },
    {
      label: "播放状态",
      width: 110,
      cellRenderer: scope => (
        <el-switch
          size="default"
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
      width: 160,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs.unix(createTime).format("YYYY-MM-DD HH:mm") : "-"
    },
    {
      label: "操作",
      fixed: "right",
      width: 160,
      slot: "operation"
    }
  ];

  function onChange({ row }) {
    ElMessageBox.confirm(
      `确认要<strong>${
        row.state === 0 ? "禁用" : "启用"
      }</strong>【<strong style='color:var(--el-color-primary)'>${
        row.name
      }</strong>】吗?`,
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
            message("已成功修改音乐状态", { type: "success" });
          } else {
            message(r.message, { type: "error" });
          }
        });
      })
      .catch(() => {
        row.state = row.state === 0 ? 1 : 0;
      });
  }

  async function handleUpdate(row) {
    dialogFormVisible.value = true;
    title.value = "编辑音乐信息";
    const { payload } = await getMusic(row.id);
    form.value = payload;
  }

  const handleAdd = async () => {
    resetForm();
    dialogFormVisible.value = true;
    title.value = "添加背景音乐";
  };

  function handleDelete(row) {
    if (currentPlaying.value?.id === row.id) {
      stopPlay();
    }
    deleteMusic(row.id).then(r => {
      if (r.code === 200) {
        message("已成功删除该音乐", { type: "success" });
        onSearch();
      } else {
        message(r.message, { type: "error" });
      }
    });
  }

  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    await formEl.validate((valid, _fields) => {
      if (valid) {
        if (form.value.id !== undefined) {
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
      }
    });
  };

  async function onSearch() {
    loading.value = true;
    try {
      const { payload } = await getMusicList();
      dataList.value = payload || [];
    } finally {
      loading.value = false;
    }
  }

  onMounted(() => {
    onSearch();
  });

  onUnmounted(() => {
    stopPlay();
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
    currentPlaying,
    isPlaying,
    togglePlay,
    stopPlay,
    handleAdd,
    submitForm,
    onSearch,
    handleUpdate,
    handleDelete
  };
}
