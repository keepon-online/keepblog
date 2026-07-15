<script setup lang="ts">
defineOptions({
  name: "LoginLog"
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
  form,
  formRef,
  resetForm,
  onSizeChange,
  handleCurrentChange,
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
      <el-form-item label="标题:" prop="title">
        <el-input
          v-model="form.title"
          placeholder="请输入标题名称"
          clearable
          class="!w-[200px]"
        />
      </el-form-item>
      <el-form-item label="状态：" prop="published">
        <el-select
          v-model="form.success"
          placeholder="请选择状态"
          clearable
          class="!w-[180px]"
        >
          <el-option label="成功" :value="1" />
          <el-option label="失败" :value="2" />
        </el-select>
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
    <PureTableBar title="登录日志" :columns="columns" @refresh="onSearch">
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
