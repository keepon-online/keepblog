<template>
  <div class="about-container">
    <el-card class="about-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">关于我们</span>
          <el-button
            type="primary"
            :loading="saveLoading"
            @click="submitForm(ruleFormRef)"
          >
            保存内容
          </el-button>
        </div>
      </template>
      
      <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        label-position="top"
        class="about-form"
        status-icon
      >
        <el-form-item label="文章标题" prop="title">
          <el-input 
            v-model="ruleForm.title" 
            placeholder="请输入文章标题"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="文章正文" prop="note" class="markdown-field">
          <MdEditor 
            v-model="ruleForm.note" 
            @onUploadImg="onUploadImg" 
            :toolbars="toolbars"
            :preview="false"
            language="zh-CN"
            class="markdown-editor"
          />
          <div class="form-tip">支持 Markdown 语法格式</div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { getAbout, updateAbout } from "@/api/about";
import { message } from "@/utils/message";
import { upload } from "@/api/common";

const ruleFormRef = ref<FormInstance>();
const saveLoading = ref(false);

const ruleForm = ref({
  postId: undefined,
  title: "",
  note: ""
});

// Markdown编辑器工具栏配置
const toolbars = [
  'bold', 'underline', 'italic', '-',
  'title', 'strikeThrough', 'sub', 'sup', '-',
  'quote', 'unorderedList', 'orderedList', '-',
  'codeRow', 'code', 'link', 'image', 'table', '-',
  'revoke', 'next', 'save', '-',
  'preview', 'catalog'
];

onMounted(() => {
  loadData();
});

const loadData = async () => {
  try {
    const res = await getAbout();
    if (res.code === 200) {
      ruleForm.value = res.payload;
    } else {
      message(`加载数据失败: ${res.message || "未知错误"}`, { type: "error" });
    }
  } catch (error) {
    message("加载数据时发生错误", { type: "error" });
  }
};

const rules = reactive<FormRules>({
  title: [{ required: true, message: "请输入标题", trigger: "blur" }],
  note: [{ required: true, message: "请输入内容", trigger: "blur" }]
});

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      saveLoading.value = true;
      try {
        const res = await updateAbout(ruleForm.value);
        if (res.code === 200) {
          message("内容更新成功", { type: "success" });
        } else {
          message(`更新失败: ${res.message || "未知错误"}`, { type: "error" });
        }
      } catch (error) {
        message("保存过程中发生错误", { type: "error" });
      } finally {
        saveLoading.value = false;
      }
    } else {
      message("表单验证失败，请检查输入内容", { type: "warning" });
    }
  });
};

const onUploadImg = async (files: File[], callback: (urls: string[]) => void) => {
  const res = await Promise.all(
    files.map(file => {
      return new Promise((resolve, reject) => {
        const form = new FormData();
        form.append("file", file);
        upload(form)
          .then(res => {
            if (res.code === 200) {
              resolve(res.payload);
            } else {
              reject(new Error(res.message || "上传失败"));
            }
          })
          .catch(error => reject(error));
      });
    })
  );
  
  // 处理上传结果
  const urls = res.map(item => {
    if (typeof item === 'string') {
      return item;
    } else if (item && typeof item === 'object' && item.url) {
      return item.url;
    }
    return '';
  });
  
  callback(urls);
  
  // 显示上传结果提示
  const successCount = urls.filter(url => url).length;
  const failCount = files.length - successCount;
  
  if (failCount > 0) {
    message(`${failCount}个文件上传失败`, { type: "warning" });
  }
};
</script>

<style scoped lang="scss">
.about-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);
  
  .about-card {
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    border: none;
    
    :deep(.el-card__header) {
      border-bottom: 1px solid var(--el-border-color-light);
      padding: 18px 20px;
      
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        
        .card-title {
          font-size: 18px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }
      }
    }
  }
  
  .about-form {
    :deep(.el-form-item) {
      margin-bottom: 24px;
      
      .el-form-item__label {
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin-bottom: 8px;
      }
    }
    
    .markdown-field {
      .markdown-editor {
        border: 1px solid var(--el-border-color);
        border-radius: 4px;
        transition: var(--el-transition-border);
        
        &:hover {
          border-color: var(--el-border-color-hover);
        }
        
        &:focus-within {
          border-color: var(--el-color-primary);
        }
      }
    }
  }
  
  .form-tip {
    font-size: 12px;
    color: var(--el-color-info);
    margin-top: 6px;
    line-height: 1.5;
  }
}

// 响应式优化
@media (max-width: 768px) {
  .about-container {
    padding: 12px;
    
    :deep(.el-card__header) {
      padding: 15px;
    }
    
    .card-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
      
      .card-title {
        font-size: 16px;
      }
      
      .el-button {
        align-self: flex-end;
      }
    }
  }
}
</style>