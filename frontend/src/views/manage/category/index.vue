<script setup lang="ts">
defineOptions({
  name: "Category"
});

import { useCategory } from "./hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { Delete, Edit, Plus } from "@element-plus/icons-vue";

const {
  dialogFormVisible,
  loading,
  columns,
  dataList,
  form,
  rules,
  ruleFormRef,
  title,
  submitForm,
  handleAdd,
  onSearch,
  handleUpdate,
  handleDelete
} = useCategory();
</script>

<template>
  <div class="main">
    <PureTableBar title="分类列表" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="handleAdd"
        >
          新增分类
        </el-button>
      </template>
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
          :paginationSmall="size === 'small' ? true : false"
          :header-cell-style="{
            background: 'var(--el-table-row-hover-bg-color)',
            color: 'var(--el-text-color-primary)'
          }"
        >
          <template #operation="{ row }">
            <el-button
              class="reset-margin"
              link
              type="primary"
              :size="size"
              :icon="useRenderIcon(Edit)"
              @click="handleUpdate(row)"
            >
              修改
            </el-button>
            <el-popconfirm title="是否确认删除?" @confirm="handleDelete(row)">
              <template #reference>
                <el-button
                  class="reset-margin"
                  link
                  type="primary"
                  :size="size"
                  :icon="useRenderIcon(Delete)"
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </pure-table>
      </template>
    </PureTableBar>
    <el-dialog v-model="dialogFormVisible" :title="title">
      <el-form ref="ruleFormRef" :model="form" :rules="rules">
        <el-form-item label="分类名称" label-width="140px" prop="categoryName">
          <el-input v-model="form.categoryName" autocomplete="off" />
        </el-form-item>
        <el-form-item label="分类描述" label-width="140px" prop="note">
          <el-input v-model="form.note" type="textarea" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogFormVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm(ruleFormRef)">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
:deep(.el-dropdown-menu__item i) {
  margin: 0;
}
</style>
