<template>
  <div class="preview-container">
    <el-card class="preview-card">
      <template #header>
        <div class="card-header">
          <div class="header-info">
            <h2 class="article-title">{{ ruleForm.title }}</h2>
            <div class="article-meta">
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                {{ formattedCreateTime }}
              </span>
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ ruleForm.viewCount || 0 }} 次阅读
              </span>
              <span class="meta-item">
                <el-icon><FolderOpened /></el-icon>
                {{ ruleForm.categoryName }}
              </span>
            </div>
          </div>
          <div class="header-actions">
            <el-button @click="router.go(-1)">返回</el-button>
            <el-button 
              type="primary" 
              @click="handleEdit"
            >
              编辑
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="article-content">
        <div class="article-cover" v-if="ruleForm.coverImage">
          <el-image 
            :src="ruleForm.coverImage" 
            fit="cover" 
            class="cover-image"
          />
        </div>
        
        <div class="article-summary" v-if="ruleForm.summary">
          <blockquote class="summary-content">{{ ruleForm.summary }}</blockquote>
        </div>
        
        <MdEditor 
          v-model="ruleForm.postContent" 
          previewOnly 
          editorId="postId" 
          :preview-theme="theme"
          class="markdown-preview"
        />
      </div>
      
      <div class="article-tags" v-if="ruleForm.tags && ruleForm.tags.length > 0">
        <el-tag 
          v-for="(tag, index) in ruleForm.tags" 
          :key="index" 
          type="info" 
          class="tag-item"
        >
          {{ tag }}
        </el-tag>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
defineOptions({
  name: "Preview"
});

import { onMounted, ref, computed } from "vue";
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { getPost } from "@/api/post";
import { useRoute, useRouter } from "vue-router";
import { Calendar, View, FolderOpened } from "@element-plus/icons-vue";

const route = useRoute();
const router = useRouter();

const ruleForm = ref({
  postId: undefined,
  title: "",
  categoryName: "",
  tags: [],
  summary: "",
  postContent: "",
  coverImage: "",
  viewCount: 0,
  createTime: ""
});

const theme = ref("default");

// 修复时间格式化计算属性
const formattedCreateTime = computed(() => {
  if (!ruleForm.value.createTime) return "";
  
  // 尝试解析时间戳
  const date = new Date(ruleForm.value.createTime);
  if (isNaN(date.getTime())) {
    // 如果不是有效日期，尝试作为时间戳解析
    const timestampNum = Number(ruleForm.value.createTime);
    if (!isNaN(timestampNum)) {
      return new Date(timestampNum).toLocaleDateString('zh-CN');
    }
    return ruleForm.value.createTime;
  }
  
  // 格式化为 YYYY-MM-DD HH:mm:ss
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
});

const handleEdit = () => {
  if (ruleForm.value.postId) {
    router.push({ 
      name: "内容编辑", 
      params: { id: String(ruleForm.value.postId) } 
    });
  }
};

onMounted(() => {
  const id = route.params.id;
  if (id) {
    getPost(id).then(res => {
      if (res.code === 200) {
        ruleForm.value = res.payload;
      }
    });
  }
});
</script>

<style scoped lang="scss">
.preview-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);
  
  .preview-card {
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    border: none;
    
    :deep(.el-card__header) {
      border-bottom: 1px solid var(--el-border-color-light);
      padding: 20px;
      
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 20px;
        
        .header-info {
          flex: 1;
          
          .article-title {
            margin: 0 0 15px 0;
            font-size: 24px;
            font-weight: 600;
            color: var(--el-text-color-primary);
            line-height: 1.4;
          }
          
          .article-meta {
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
            
            .meta-item {
              display: flex;
              align-items: center;
              gap: 5px;
              font-size: 14px;
              color: var(--el-text-color-secondary);
              
              .el-icon {
                font-size: 16px;
              }
            }
          }
        }
        
        .header-actions {
          display: flex;
          gap: 12px;
          flex-shrink: 0;
        }
      }
    }
    
    :deep(.el-card__body) {
      padding: 30px 20px;
    }
  }
  
  .article-content {
    .article-cover {
      margin-bottom: 30px;
      text-align: center;
      
      .cover-image {
        max-width: 100%;
        max-height: 400px;
        border-radius: 8px;
        box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      }
    }
    
    .article-summary {
      margin-bottom: 30px;
      
      .summary-content {
        margin: 0;
        padding: 15px 20px;
        border-left: 4px solid var(--el-color-primary);
        background-color: var(--el-fill-color-light);
        color: var(--el-text-color-secondary);
        font-size: 15px;
        line-height: 1.6;
        border-radius: 0 4px 4px 0;
      }
    }
    
    .markdown-preview {
      :deep(.md-editor-preview) {
        padding: 0;
      }
      
      :deep(.md-editor) {
        box-shadow: none;
      }
    }
  }
  
  .article-tags {
    margin-top: 30px;
    padding-top: 20px;
    border-top: 1px solid var(--el-border-color-light);
    
    .tag-item {
      margin-right: 10px;
      margin-bottom: 10px;
    }
  }
}

// 响应式优化
@media (max-width: 768px) {
  .preview-container {
    padding: 12px;
    
    :deep(.el-card__header) {
      padding: 15px;
    }
    
    .card-header {
      flex-direction: column;
      align-items: stretch;
      gap: 15px;
      
      .header-actions {
        align-self: flex-end;
      }
    }
    
    .article-title {
      font-size: 20px !important;
    }
    
    .article-meta {
      gap: 12px !important;
      
      .meta-item {
        font-size: 13px !important;
      }
    }
  }
}
</style>