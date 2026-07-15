<script setup lang="ts">
import { useLink } from "./hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { Delete, Edit, Plus } from "@element-plus/icons-vue";

defineOptions({
  name: "Link"
});

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
} = useLink();
</script>

<template>
  <div class="main">
    <PureTableBar title="友情链接" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="handleAdd"
        >
          新增友链
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
          <template #image="{ row, index }">
            <el-image
              preview-teleported
              loading="lazy"
              :src="row.linkIcon"
              :initial-index="index"
              fit="cover"
              class="w-[100px] h-[100px]"
            />
          </template>
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
        <el-form-item label="名称" label-width="140px" prop="title">
          <el-input v-model="form.title" autocomplete="off" />
        </el-form-item>
        <el-form-item label="链接" label-width="140px" prop="linkUrl">
          <el-input v-model="form.linkUrl" autocomplete="off" />
        </el-form-item>
        <el-form-item label="LOGO" label-width="140px" prop="linkIcon">
          <el-input v-model="form.linkIcon" autocomplete="off" />
        </el-form-item>
        <el-form-item label="LOGO" label-width="140px" prop="type">
          <el-select
            v-model="form.type"
            placeholder="请选择类型"
            clearable
            style="width: 100%"
          >
            <el-option label="友链" :value="1" />
            <el-option label="网站" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" label-width="140px" prop="linkDesc">
          <el-input
            v-model="form.linkDesc"
            type="textarea"
            autocomplete="off"
          />
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
