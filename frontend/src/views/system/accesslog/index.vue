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
  <div class="main">
    <el-form
      ref="formRef"
      :inline="true"
      :model="form"
      class="bg-bg_color w-[99/100] pl-8 pt-4"
    >
      <el-form-item label="IP:" prop="title">
        <el-input
          v-model="form.ip"
          placeholder="请输入IP地址"
          clearable
          class="!w-[200px]"
        />
      </el-form-item>
      <el-form-item label="状态:" prop="status">
        <el-select
          v-model="form.status"
          placeholder="请选择状态"
          clearable
          class="!w-[180px]"
        >
          <el-option label="成功" :value="200" />
          <el-option label="失败" :value="404" />
        </el-select>
      </el-form-item>
      <el-form-item label="日期" prop="published">
        <el-date-picker
          v-model="daterange"
          placeholder="选择日期"
          type="daterange"
          unlink-panels
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :shortcuts="shortcuts"
          :disabled-date="disabledDate"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          :icon="useRenderIcon(Search)"
          :loading="loading"
          @click="onSearch"
        >
          搜索
        </el-button>
        <el-button :icon="useRenderIcon(Refresh)" @click="resetForm(formRef)">
          重置
        </el-button>
      </el-form-item>
    </el-form>
    <PureTableBar title="访问日志" :columns="columns" @refresh="onSearch">
      <template v-slot="{ size, dynamicColumns }">
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
      </template>
    </PureTableBar>
  </div>
</template>

<style scoped lang="scss">
:deep(.el-dropdown-menu__item i) {
  margin: 0;
}
</style>
