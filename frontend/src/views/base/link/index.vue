<script setup lang="ts">
import { useLink } from "./hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { Delete, Edit, Plus, CircleCheck } from "@element-plus/icons-vue";

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
  checkingAll,
  handleAdd,
  handleCheckAll,
  submitForm,
  onSearch,
  handleUpdate,
  handleDelete
} = useLink();
</script>

<template>
  <div class="link-manage-container">
    <PureTableBar title="友情链接与外链管理" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="handleAdd"
        >
          新增友链
        </el-button>
        <el-button
          type="success"
          plain
          :icon="useRenderIcon(CircleCheck)"
          :loading="checkingAll"
          @click="handleCheckAll"
        >
          一键存活体检
        </el-button>
      </template>

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
              <el-popconfirm
                :title="`是否确认删除【${row.title}】?`"
                @confirm="handleDelete(row)"
              >
                <template #reference>
                  <el-button
                    class="reset-margin"
                    link
                    type="danger"
                    :size="size"
                    :icon="useRenderIcon(Delete)"
                  >
                    删除
                  </el-button>
                </template>
              </el-popconfirm>
            </template>
          </pure-table>
        </el-card>
      </template>
    </PureTableBar>

    <el-dialog
      v-model="dialogFormVisible"
      :title="title"
      width="580px"
      destroy-on-close
      class="link-dialog"
    >
      <el-form
        ref="ruleFormRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-row :gutter="16">
          <el-col :span="14">
            <el-form-item label="网站名称" prop="title">
              <el-input
                v-model="form.title"
                placeholder="例如: 张三的技术博客"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="10">
            <el-form-item label="类型" prop="type">
              <el-select v-model="form.type" placeholder="选择类型" class="w-full">
                <el-option label="网站" :value="0" />
                <el-option label="友链" :value="1" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="网站链接" prop="linkUrl">
          <el-input
            v-model="form.linkUrl"
            placeholder="例如: https://example.com"
            clearable
          />
        </el-form-item>

        <el-form-item label="站点 Logo / Favicon 图标" prop="linkIcon">
          <el-input
            v-model="form.linkIcon"
            placeholder="图片直链地址，如 https://example.com/favicon.png"
            clearable
          />
        </el-form-item>

        <el-form-item label="网站简介描述" prop="linkDesc">
          <el-input
            v-model="form.linkDesc"
            type="textarea"
            :rows="3"
            placeholder="简要描述站点特色与博主介绍..."
          />
        </el-form-item>

        <!-- 实时卡片预览 -->
        <div v-if="form.title || form.linkUrl" class="preview-box">
          <div class="preview-title">实时卡片效果预览</div>
          <div class="link-card-preview">
            <div class="preview-avatar">
              <el-image
                v-if="form.linkIcon"
                :src="form.linkIcon"
                fit="cover"
                style="width: 44px; height: 44px; border-radius: 8px;"
              >
                <template #error>
                  <div class="avatar-fallback">
                    {{ (form.title || "L").slice(0, 1).toUpperCase() }}
                  </div>
                </template>
              </el-image>
              <div v-else class="avatar-fallback">
                {{ (form.title || "L").slice(0, 1).toUpperCase() }}
              </div>
            </div>
            <div class="preview-content">
              <div class="preview-name">{{ form.title || "网站名称" }}</div>
              <div class="preview-desc">{{ form.linkDesc || "暂无简介" }}</div>
              <div class="preview-url">{{ form.linkUrl || "https://..." }}</div>
            </div>
          </div>
        </div>
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
.link-manage-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 120px);

  .table-card {
    border-radius: 10px;
    border: 1px solid var(--el-border-color-lighter);

    :deep(.el-card__body) {
      padding: 0;
    }
  }
}

.preview-box {
  margin-top: 14px;
  padding: 12px;
  background: var(--el-fill-color-light);
  border-radius: 8px;
  border: 1px dashed var(--el-border-color);

  .preview-title {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-bottom: 8px;
    font-weight: 500;
  }

  .link-card-preview {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: var(--el-bg-color-overlay);
    border-radius: 8px;
    border: 1px solid var(--el-border-color-light);

    .preview-avatar {
      flex-shrink: 0;

      .avatar-fallback {
        width: 44px;
        height: 44px;
        border-radius: 8px;
        background: linear-gradient(135deg, #409eff, #36a3f7);
        color: #fff;
        font-size: 18px;
        font-weight: bold;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .preview-content {
      flex: 1;
      min-width: 0;

      .preview-name {
        font-size: 14px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .preview-desc {
        font-size: 12px;
        color: var(--el-text-color-secondary);
        margin-top: 2px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .preview-url {
        font-size: 11px;
        color: var(--el-color-primary);
        margin-top: 2px;
      }
    }
  }
}
</style>
