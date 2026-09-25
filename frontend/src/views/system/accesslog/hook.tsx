import dayjs from "dayjs";
import { getAccessLog } from "@/api/system";
import { Int2Ip } from "@/utils/utils";
import { ref, onMounted, reactive, h } from "vue";
import type { PaginationProps } from "@pureadmin/table";
import { message } from "@/utils/message";

export function useLogs() {
  const dataList = ref([]);
  const loading = ref(true);
  const daterange = ref([]);
  const form = reactive({
    title: undefined,
    status: undefined,
    ip: undefined
  });
  const formRef = ref();
  const pagination = reactive<PaginationProps>({
    total: 0,
    pageSize: 10,
    pageSizes: [10, 15, 20, 50, 100],
    currentPage: 1,
    background: true
  });

  const disabledDate = (time: Date) => {
    return time.getTime() > Date.now();
  };

  const copyText = (text: string, label = "文本") => {
    navigator.clipboard.writeText(text);
    message(`${label}已复制`, { type: "success" });
  };

  const parseUserAgent = (ua: string) => {
    if (!ua) return { os: "未知系统", browser: "未知浏览器" };
    let os = "其他系统";
    let browser = "其他浏览器";

    if (/windows nt 10/i.test(ua)) os = "Windows 10/11";
    else if (/windows nt 6.3/i.test(ua)) os = "Windows 8.1";
    else if (/windows nt 6.1/i.test(ua)) os = "Windows 7";
    else if (/windows/i.test(ua)) os = "Windows";
    else if (/macintosh|mac os x/i.test(ua)) os = "macOS";
    else if (/iphone/i.test(ua)) os = "iPhone";
    else if (/ipad/i.test(ua)) os = "iPad";
    else if (/android/i.test(ua)) os = "Android";
    else if (/linux/i.test(ua)) os = "Linux";

    if (/micromessenger/i.test(ua)) browser = "微信内置";
    else if (/edg/i.test(ua)) browser = "Edge";
    else if (/chrome/i.test(ua) && !/edg/i.test(ua)) browser = "Chrome";
    else if (/safari/i.test(ua) && !/chrome/i.test(ua)) browser = "Safari";
    else if (/firefox/i.test(ua)) browser = "Firefox";
    else if (/opera|opr/i.test(ua)) browser = "Opera";

    return { os, browser };
  };

  const formatArea = (area: string) => {
    if (!area) return "-";
    // 过滤掉 ip2region 中的 "0|"
    return area
      .split("|")
      .filter(p => p && p !== "0")
      .join(" ");
  };

  const shortcuts = [
    {
      text: "今天",
      value: new Date()
    },
    {
      text: "昨天",
      value: () => {
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24);
        return [start, start];
      }
    },
    {
      text: "近7天",
      value: () => {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
        return [start, end];
      }
    },
    {
      text: "近30天",
      value: () => {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
        return [start, end];
      }
    }
  ];

  const columns: TableColumnList = [
    {
      label: "序号",
      type: "index",
      width: 60
    },
    {
      label: "访客 IP / 地区",
      prop: "ip",
      minWidth: 180,
      align: "left",
      cellRenderer: ({ row }) => {
        const ipStr = Int2Ip(row.ip);
        const areaStr = formatArea(row.area);
        return (
          <div class="flex flex-col gap-1 py-1">
            <div class="flex items-center gap-2">
              <span style="font-family: monospace; font-weight: 600; color: var(--el-text-color-primary);">
                {ipStr}
              </span>
              <span
                class="cursor-pointer text-xs"
                title="复制 IP"
                onClick={() => copyText(ipStr, "IP")}
                style="color: var(--el-text-color-placeholder);"
              >
                📋
              </span>
            </div>
            <div style="font-size: 11px; color: var(--el-text-color-secondary);">
              📍 {areaStr}
            </div>
          </div>
        );
      }
    },
    {
      label: "HTTP 状态",
      prop: "status",
      width: 110,
      cellRenderer: ({ row }) => {
        const status = row.status || 200;
        let tagType: "success" | "warning" | "danger" | "info" = "info";
        if (status >= 200 && status < 300) tagType = "success";
        else if (status >= 300 && status < 400) tagType = "warning";
        else if (status >= 400) tagType = "danger";

        return (
          <el-tag size="small" type={tagType} effect="light">
            {status} {status === 200 ? "OK" : ""}
          </el-tag>
        );
      }
    },
    {
      label: "访问目标 URL",
      prop: "url",
      minWidth: 200,
      align: "left",
      cellRenderer: ({ row }) => (
        <div class="flex flex-col gap-1 py-1">
          <span
            style="color: var(--el-color-primary); font-family: monospace; font-size: 13px; max-width: 320px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
            title={row.url}
          >
            {row.url || "/"}
          </span>
          {row.referer ? (
            <span
              style="font-size: 11px; color: var(--el-text-color-placeholder); max-width: 320px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
              title={`来源: ${row.referer}`}
            >
              来源: {row.referer}
            </span>
          ) : null}
        </div>
      )
    },
    {
      label: "客户端环境",
      prop: "ua",
      minWidth: 180,
      cellRenderer: ({ row }) => {
        const { os, browser } = parseUserAgent(row.ua);
        return (
          <el-tooltip content={row.ua || "无 User-Agent"} placement="top">
            <div class="flex items-center gap-1.5 cursor-pointer">
              <el-tag size="small" effect="plain" type="info">
                {os}
              </el-tag>
              <el-tag size="small" effect="plain" type="primary">
                {browser}
              </el-tag>
            </div>
          </el-tooltip>
        );
      }
    },
    {
      label: "访问时间",
      minWidth: 160,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs.unix(createTime).format("YYYY-MM-DD HH:mm:ss") : "-"
    }
  ];

  const resetForm = formEl => {
    if (!formEl) return;
    formEl.resetFields();
    daterange.value = [];
    form.ip = undefined;
    form.status = undefined;
    onSearch();
  };

  async function handleCurrentChange(val: number) {
    loading.value = true;
    pagination.currentPage = val;
    const [startDate, endDate] = daterange.value || [];
    const params = {
      pageNum: pagination.currentPage,
      pageSize: pagination.pageSize,
      ip: form.ip,
      status: form.status,
      start: startDate ? dayjs(startDate).format("YYYY-MM-DD") : "",
      end: endDate ? dayjs(endDate).format("YYYY-MM-DD") : ""
    };
    try {
      const { payload } = await getAccessLog(params);
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
    const [startDate, endDate] = daterange.value || [];
    const params = {
      pageNum: 1,
      pageSize: pagination.pageSize,
      ip: form.ip,
      status: form.status,
      start: startDate ? dayjs(startDate).format("YYYY-MM-DD") : "",
      end: endDate ? dayjs(endDate).format("YYYY-MM-DD") : ""
    };
    try {
      const { payload } = await getAccessLog(params);
      dataList.value = payload.list || [];
      pagination.total = payload.total || 0;
    } finally {
      loading.value = false;
    }
  }

  onMounted(() => {
    onSearch();
  });

  return {
    loading,
    formRef,
    columns,
    dataList,
    pagination,
    shortcuts,
    daterange,
    form,
    disabledDate,
    resetForm,
    onSizeChange,
    handleCurrentChange,
    onSearch
  };
}
