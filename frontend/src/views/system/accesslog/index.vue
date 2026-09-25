<script setup lang="ts">
defineOptions({
  name: "AccessLog"
});

import { useLogs } from "./hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { Search, Refresh } from "@element-plus/icons-vue";

const {
  loading,
  columns,
  dataList,
  pagination,
  shortcuts,
  daterange,
  form,
  formRef,
  disabledDate,
  resetForm,
  handleCurrentChange,
  onSizeChange,
  onSearch
} = useLogs();
</script>

<template>
  <div class="access-log-container">
    <el-card shadow="never" class="search-card">
      <el-form
        ref="formRef"
        :inline="true"
        :model="form"
        class="search-form"
      >
        <el-form-item label="访客 IP:" prop="ip">
          <el-input
            v-model="form.ip"
            placeholder="例如: 127.0.0.1"
            clearable
            class="search-input"
          />
        </el-form-item>
        <el-form-item label="HTTP 状态:" prop="status">
          <el-select
            v-model="form.status"
            placeholder="全部状态"
            clearable
            class="search-select"
          >
            <el-option label="200 (成功)" :value="200" />
            <el-option label="301/302 (跳转)" :value="301" />
            <el-option label="404 (未找到)" :value="404" />
            <el-option label="500 (服务异常)" :value="500" />
          </el-select>
        </el-form-item>
        <el-form-item label="访问日期:">
          <el-date-picker
            v-model="daterange"
            type="daterange"
            unlink-panels
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :shortcuts="shortcuts"
            :disabled-date="disabledDate"
            class="date-picker"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :icon="useRenderIcon(Search)"
            :loading="loading"
            @click="onSearch"
          >
            查询
          </el-button>
          <el-button :icon="useRenderIcon(Refresh)" @click="resetForm(formRef)">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <PureTableBar title="网站访问安全审计日志" :columns="columns" @refresh="onSearch">
      <template v-slot="{ size, dynamicColumns }">
        <el-card shadow="never" class="table-card">
          <pure-table
            border
            align-whole="center"
            showOverflowTooltip
            table-layout="auto"
            :loading="loading"
            :size="size"
            :data="dataList"
            :columns="dynamicColumns"
            :pagination="pagination"
            :paginationSmall="size === 'small' ? true : false"
            :header-cell-style="{
              background: 'var(--el-table-row-hover-bg-color)',
              color: 'var(--el-text-color-primary)'
            }"
            @page-current-change="handleCurrentChange"
            @page-size-change="onSizeChange"
          />
        </el-card>
      </template>
    </PureTableBar>
  </div>
</template>

<style scoped lang="scss">
.access-log-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 120px);

  .search-card {
    margin-bottom: 16px;
    border-radius: 10px;
    border: 1px solid var(--el-border-color-lighter);

    :deep(.el-card__body) {
      padding: 18px 20px 2px;
    }

    .search-input {
      width: 180px;
    }

    .search-select {
      width: 160px;
    }

    .date-picker {
      width: 260px;
    }
  }

  .table-card {
    border-radius: 10px;
    border: 1px solid var(--el-border-color-lighter);

    :deep(.el-card__body) {
      padding: 0;
    }
  }
}
</style>
