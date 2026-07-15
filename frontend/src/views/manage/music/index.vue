<script setup lang="ts">
defineOptions({
  name: "Music"
});

import { useMusic } from "./hook";
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
} = useMusic();
</script>

<template>
  <div class="main">
    <PureTableBar title="音乐列表" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="handleAdd"
        >
          添加音乐
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
              编辑
            </el-button>
            <el-popconfirm title="确认删除该音乐?" @confirm="handleDelete(row)">
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

    <el-dialog v-model="dialogFormVisible" :title="title" width="600px">
      <el-form ref="ruleFormRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="歌曲名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入歌曲名称" />
        </el-form-item>
        <el-form-item label="艺术家" prop="artist">
          <el-input v-model="form.artist" placeholder="请输入艺术家/歌手" />
        </el-form-item>
        <el-form-item label="音乐URL" prop="url">
          <el-input v-model="form.url" placeholder="请输入音乐文件URL" />
        </el-form-item>
        <el-form-item label="封面URL" prop="cover">
          <el-input v-model="form.cover" placeholder="请输入封面图片URL (可选)" />
        </el-form-item>
        <el-form-item label="歌词URL" prop="lrc">
          <el-input v-model="form.lrc" placeholder="请输入歌词URL (可选)" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="9999" />
          <span class="ml-2 text-gray-400 text-sm">值越小越靠前</span>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-switch
            v-model="form.state"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
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
