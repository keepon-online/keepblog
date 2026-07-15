import dayjs from "dayjs";
import { getAccessLog } from "@/api/system";
import { Int2Ip } from "@/utils/utils";
import { ref, onMounted, reactive } from "vue";
import type { PaginationProps } from "@pureadmin/table";

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
      text: "上周",
      value: () => {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
        return [start, end];
      }
    },
    {
      text: "上月",
      value: () => {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
        return [start, end];
      }
    },
    {
      text: "三月前",
      value: () => {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
        return [start, end];
      }
    }
  ];
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
      label: "IP",
      prop: "ip",
      minWidth: 100,
      formatter: ({ ip }) => Int2Ip(ip)
    },
    {
      label: "状态",
      prop: "status"
    },
    {
      label: "地区",
      prop: "area",
      minWidth: 100
    },
    {
      label: "UA",
      prop: "ua"
    },
    {
      label: "Referer",
      prop: "referer"
    },
    {
      label: "URL",
      prop: "url"
    },
    {
      label: "创建时间",
      minWidth: 80,
      prop: "createTime",
      formatter: ({ createTime }) =>
        dayjs.unix(createTime).format("YYYY-MM-DD HH:mm:ss")
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

  async function handleCurrentChange(val) {
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
    const { payload } = await getAccessLog(params);
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
    const [startDate, endDate] = daterange.value || [];
    const params = {
      pageNum: 1,
      pageSize: pagination.pageSize,
      ip: form.ip,
      status: form.status,
      start: startDate ? dayjs(startDate).format("YYYY-MM-DD") : "",
      end: endDate ? dayjs(endDate).format("YYYY-MM-DD") : ""
    };
    const { payload } = await getAccessLog(params);
    dataList.value = payload.list;
    pagination.total = payload.total;
    setTimeout(() => {
      loading.value = false;
    }, 500);
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
