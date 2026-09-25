<script setup lang="ts">
defineOptions({
  name: "Post"
});
import { ref } from "vue";
import { usePost } from "./hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";

import {
  Folder,
  MoreFilled,
  Delete,
  Edit,
  Search,
  Refresh,
  Plus,
  View,
  Check,
  Close
} from "@element-plus/icons-vue";

const formRef = ref();
const {
  form,
  loading,
  columns,
  dataList,
  pagination,
  buttonClass,
  categories,
  selectedRows,
  batchLoading,
  handleSelectionChange,
  handleBatchPublish,
  handleBatchDelete,
  onSearch,
  resetForm,
  handleUpdate,
  handleDelete,
  onSizeChange,
  handleUpdateCoverImage,
  handleUpdateAllCoverImage,
  handleCurrentChange
} = usePost();
</script>

<template>
  <div class="post-container">
    <el-card shadow="never" class="search-card">
      <el-form ref="formRef" :inline="true" :model="form" class="search-form">
        <el-form-item label="文章标题:" prop="title">
          <el-input
            v-model="form.title"
            placeholder="支持关键词模糊匹配"
            clearable
            class="search-input"
          />
        </el-form-item>
        <el-form-item label="所属分类:" prop="categoryId">
          <el-select
            v-model="form.categoryId"
            clearable
            filterable
            placeholder="全部分类"
            class="search-select"
          >
            <el-option
              v-for="item in categories"
              :key="item.categoryId"
              :label="item.categoryName"
              :value="item.categoryId"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发布状态:" prop="published">
          <el-select
            v-model="form.published"
            placeholder="全部状态"
            clearable
            class="search-select"
          >
            <el-option label="已发布" :value="1" />
            <el-option label="未发布 (草稿)" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :icon="useRenderIcon(Search)"
            :loading="loading"
            @click="onSearch"
          >
            筛选
          </el-button>
          <el-button :icon="useRenderIcon(Refresh)" @click="resetForm(formRef)">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <PureTableBar title="博文内容资产列表" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="$router.push({ path: '/manage/editor' })"
        >
          撰写文章
        </el-button>
        <el-button
          type="primary"
          plain
          :icon="useRenderIcon(Refresh)"
          @click="handleUpdateAllCoverImage()"
        >
          全量刷新封面
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
            :pagination="pagination"
            :paginationSmall="size === 'small' ? true : false"
            :header-cell-style="{
              background: 'var(--el-table-row-hover-bg-color)',
              color: 'var(--el-text-color-primary)'
            }"
            @selection-change="handleSelectionChange"
            @page-current-change="handleCurrentChange"
            @page-size-change="onSizeChange"
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
              <el-button
                class="reset-margin"
                link
                type="primary"
                :size="size"
                :icon="useRenderIcon(View)"
                @click="
                  $router.push({
                    params: { id: String(row.postSlug || row.postId) },
                    name: '内容预览'
                  })
                "
              >
                预览
              </el-button>
              <el-dropdown>
                <el-button
                  class="ml-2 mt-[2px]"
                  link
                  type="primary"
                  :size="size"
                  :icon="useRenderIcon(MoreFilled)"
                />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item>
                      <el-button
                        :class="buttonClass"
                        link
                        type="primary"
                        :size="size"
                        :icon="useRenderIcon(Folder)"
                        @click="handleUpdateCoverImage(row)"
                      >
                        随机更换封面
                      </el-button>
                    </el-dropdown-item>
                    <el-dropdown-item divided>
                      <el-popconfirm
                        :title="`确认永久删除【${row.title}】?`"
                        @confirm="handleDelete(row)"
                      >
                        <template #reference>
                          <el-button
                            class="reset-margin !text-red-500"
                            link
                            type="danger"
                            :size="size"
                            :icon="useRenderIcon(Delete)"
                          >
                            彻底删除
                          </el-button>
                        </template>
                      </el-popconfirm>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </pure-table>
        </el-card>
      </template>
    </PureTableBar>

    <!-- 批量操作悬浮底栏 -->
    <transition name="el-zoom-in-bottom">
      <div v-if="selectedRows.length > 0" class="floating-batch-bar">
        <div class="batch-info">
          已选择 <span class="count-num">{{ selectedRows.length }}</span> 篇文章
        </div>
        <div class="batch-actions">
          <el-button
            type="success"
            size="small"
            :loading="batchLoading"
            :icon="useRenderIcon(Check)"
            @click="handleBatchPublish(1)"
          >
            批量发布
          </el-button>
          <el-button
            type="warning"
            size="small"
            :loading="batchLoading"
            :icon="useRenderIcon(Close)"
            @click="handleBatchPublish(0)"
          >
            批量设为草稿
          </el-button>
          <el-button
            type="danger"
            size="small"
            :loading="batchLoading"
            :icon="useRenderIcon(Delete)"
            @click="handleBatchDelete"
          >
            批量删除
          </el-button>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped lang="scss">
.post-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 120px);
  position: relative;

  .search-card {
    margin-bottom: 16px;
    border-radius: 10px;
    border: 1px solid var(--el-border-color-lighter);

    :deep(.el-card__body) {
      padding: 18px 20px 2px;
    }

    .search-input {
      width: 200px;
    }

    .search-select {
      width: 160px;
    }
  }

  .table-card {
    border-radius: 10px;
    border: 1px solid var(--el-border-color-lighter);

    :deep(.el-card__body) {
      padding: 0;
    }
  }

  .floating-batch-bar {
    position: fixed;
    bottom: 28px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 999;
    display: flex;
    align-items: center;
    gap: 18px;
    padding: 10px 22px;
    background: var(--el-bg-color-overlay);
    border: 1px solid var(--el-border-color-light);
    border-radius: 30px;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
    backdrop-filter: blur(10px);

    .batch-info {
      font-size: 13px;
      font-weight: 500;
      color: var(--el-text-color-regular);

      .count-num {
        color: var(--el-color-primary);
        font-weight: 700;
        font-size: 15px;
      }
    }

    .batch-actions {
      display: flex;
      align-items: center;
      gap: 10px;
    }
  }
}

@media (max-width: 768px) {
  .post-container {
    padding: 12px;

    .floating-batch-bar {
      width: 90%;
      flex-direction: column;
      gap: 10px;
      border-radius: 14px;
      bottom: 16px;
    }
  }
}
</style>
