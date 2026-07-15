import dayjs from "dayjs";
import { getLoginLog } from "@/api/system";
import { Int2Ip } from "@/utils/utils";
import { ref, onMounted, reactive } from "vue";
import type { PaginationProps } from "@pureadmin/table";

export function useLogs() {
  const dataList = ref([]);
  const loading = ref(true);
  const form = reactive({
    title: undefined,
    success: undefined
  });
  const formRef = ref();
  // const success = ref(undefined);
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
      label: "IP",
      prop: "ip",
      minWidth: 100,
      formatter: ({ ip }) => Int2Ip(ip)
    },
    {
      label: "状态",
      prop: "success",
      minWidth: 50,
      cellRenderer: ({ row, props }) => (
        <el-tag
          size={props.size}
          type={row.success === 1 ? "danger" : ""}
          effect="plain"
        >
          {row.success === 1 ? "成功" : "失败"}
        </el-tag>
      )
    },
    {
      label: "描述",
      prop: "note",
      minWidth: 100
    },
    {
      label: "UA",
      prop: "ua"
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
    onSearch();
  };

  async function handleCurrentChange(val) {
    loading.value = true;
    pagination.currentPage = val;
    const params = {
      pageNum: pagination.currentPage,
      pageSize: pagination.pageSize,
      success: form.success
    };
    const { payload } = await getLoginLog(params);
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
      success: form.success
    };
    const { payload } = await getLoginLog(params);
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
    form,
    resetForm,
    onSizeChange,
    handleCurrentChange,
    onSearch
  };
}
