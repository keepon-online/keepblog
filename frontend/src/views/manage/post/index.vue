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
  Menu,
  Plus,
  View
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
    <el-card class="search-card">
      <el-form
        ref="formRef"
        :inline="true"
        :model="form"
        class="search-form"
      >
        <el-form-item label="标题:" prop="title">
          <el-input
            v-model="form.title"
            placeholder="请输入标题名称"
            clearable
            class="search-input"
          />
        </el-form-item>
        <el-form-item label="分类：" prop="categoryId">
          <el-select
            v-model="form.categoryId"
            clearable
            filterable
            placeholder="请选择分类"
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
        <el-form-item label="状态：" prop="published">
          <el-select
            v-model="form.published"
            placeholder="请选择状态"
            clearable
            class="search-select"
          >
            <el-option label="已发布" :value="1" />
            <el-option label="未发布" :value="0" />
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
    </el-card>

    <PureTableBar title="文章列表" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="$router.push({ path: '/manage/editor', name: '内容新增' })"
        >
          新增
        </el-button>
        <el-button
          type="primary"
          :icon="useRenderIcon(Refresh)"
          @click="handleUpdateAllCoverImage()"
        >
          更新封面
        </el-button>
      </template>
      <template v-slot="{ size, dynamicColumns }">
        <el-card class="table-card">
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
              <el-popconfirm title="是否确认删除?" @confirm="handleDelete(row)">
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
              <el-dropdown>
                <el-button
                  class="ml-3 mt-[2px]"
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
                        :icon="useRenderIcon(View)"
                        @click="
                          $router.push({
                            params: { id: String(row.postSlug) },
                            name: '内容预览'
                          })
                        "
                      >
                        预览
                      </el-button>
                    </el-dropdown-item>
                    <el-dropdown-item>
                      <el-button
                        :class="buttonClass"
                        link
                        type="primary"
                        :size="size"
                        :icon="useRenderIcon(Folder)"
                        @click="handleUpdateCoverImage(row)"
                      >
                        封面
                      </el-button>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </pure-table>
        </el-card>
      </template>
    </PureTableBar>
  </div>
</template>

<style scoped lang="scss">
.post-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);

  .search-card {
    margin-bottom: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    border: none;

    :deep(.el-card__body) {
      padding: 20px;
    }
  }

  .search-form {
    .search-input {
      width: 200px;
    }

    .search-select {
      width: 180px;
    }
  }

  .table-card {
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    border: none;

    :deep(.el-card__body) {
      padding: 0;
    }
  }
}

// 响应式优化
@media (max-width: 768px) {
  .post-container {
    padding: 12px;

    .search-card {
      :deep(.el-card__body) {
        padding: 15px;
      }
    }

    .search-form {
      .el-form-item {
        display: block;
        margin-right: 0;
        margin-bottom: 15px;
      }

      .search-input,
      .search-select {
        width: 100%;
      }
    }
  }
}
</style>
