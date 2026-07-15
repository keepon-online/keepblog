<template>
  <div class="notice-manage">
    <el-card class="notice-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">通知管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleSendNotice">
              <IconifyIconOffline :icon="PlusIcon" class="mr-1" />
              发送通知
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选区 -->
      <div class="filter-area">
        <el-form :inline="true" :model="queryParams">
          <el-form-item label="通知类型">
            <el-select v-model="queryParams.type" placeholder="全部" clearable style="width: 120px">
              <el-option label="系统通知" :value="1" />
              <el-option label="消息" :value="2" />
              <el-option label="待办" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="queryParams.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="未读" :value="0" />
              <el-option label="已读" :value="1" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchData">查询</el-button>
            <el-button @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="noticeList"
        stripe
        style="width: 100%"
      >
        <el-table-column type="index" label="#" width="60" />
        <el-table-column prop="title" label="标题" min-width="180">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="danger" size="small" class="mr-2">未读</el-tag>
            {{ row.title }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">{{ getTypeName(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="内容" min-width="250" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">查看</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-area">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 发送通知对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="发送通知"
      width="550px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="通知类型" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio :label="1">系统通知</el-radio>
            <el-radio :label="2">消息</el-radio>
            <el-radio :label="3">待办</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入标题" maxlength="100" />
        </el-form-item>
        <el-form-item label="内容" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入内容"
            maxlength="500"
          />
        </el-form-item>
        <el-form-item label="接收用户" prop="userId">
          <el-select v-model="formData.userId" placeholder="全部用户（广播）" clearable style="width: 100%">
            <el-option label="全部用户" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { getNoticeList, deleteNotice } from "@/api/notice";
import { http } from "@/utils/http";
import { message } from "@/utils/message";
import { ElMessageBox } from "element-plus";
import PlusIcon from "~icons/ep/plus";

defineOptions({ name: "NoticeManage" });

const loading = ref(false);
const noticeList = ref<any[]>([]);
const total = ref(0);
const dialogVisible = ref(false);
const submitting = ref(false);
const formRef = ref();

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  type: undefined as number | undefined,
  status: undefined as number | undefined
});

const formData = reactive({
  type: 1,
  title: "",
  description: "",
  userId: 0
});

const formRules = {
  type: [{ required: true, message: "请选择通知类型", trigger: "change" }],
  title: [{ required: true, message: "请输入标题", trigger: "blur" }],
  description: [{ required: true, message: "请输入内容", trigger: "blur" }]
};

const getTypeName = (type: number) => {
  const map: Record<number, string> = { 1: "系统通知", 2: "消息", 3: "待办" };
  return map[type] || "未知";
};

const getTypeTag = (type: number) => {
  const map: Record<number, string> = { 1: "danger", 2: "primary", 3: "warning" };
  return map[type] || "info";
};

const fetchData = async () => {
  loading.value = true;
  try {
    const res = await getNoticeList();
    if (res.code === 200 && res.payload) {
      // 展平所有分类的通知
      const allNotices: any[] = [];
      res.payload.notices?.forEach((tab: any) => {
        tab.list?.forEach((item: any) => {
          allNotices.push({
            ...item,
            type: parseInt(item.type) || parseInt(tab.key)
          });
        });
      });
      noticeList.value = allNotices;
      total.value = allNotices.length;
    }
  } catch (error) {
    message("获取通知列表失败", { type: "error" });
  } finally {
    loading.value = false;
  }
};

const resetQuery = () => {
  queryParams.type = undefined;
  queryParams.status = undefined;
  queryParams.pageNum = 1;
  fetchData();
};

const handleSendNotice = () => {
  formData.type = 1;
  formData.title = "";
  formData.description = "";
  formData.userId = 0;
  dialogVisible.value = true;
};

const handleSubmit = async () => {
  const valid = await formRef.value?.validate();
  if (!valid) return;

  submitting.value = true;
  try {
    const res = await http.request("post", "/api/v1/notice/send", { data: formData });
    if (res.code === 200) {
      message("发送成功", { type: "success" });
      dialogVisible.value = false;
      fetchData();
    } else {
      message(res.message || "发送失败", { type: "error" });
    }
  } catch (error) {
    message("发送失败", { type: "error" });
  } finally {
    submitting.value = false;
  }
};

const handleView = (row: any) => {
  ElMessageBox.alert(row.description, row.title, {
    confirmButtonText: "确定"
  });
};

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm("确定要删除该通知吗？", "提示", {
      type: "warning"
    });
    await deleteNotice(row.id);
    message("删除成功", { type: "success" });
    fetchData();
  } catch (error) {
    // 取消删除
  }
};

onMounted(() => {
  fetchData();
});
</script>

<style lang="scss" scoped>
.notice-manage {
  padding: 20px;

  .notice-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .card-title {
        font-size: 18px;
        font-weight: 600;
      }
    }
  }

  .filter-area {
    margin-bottom: 20px;
  }

  .pagination-area {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
