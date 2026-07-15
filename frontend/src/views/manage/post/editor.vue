<template>
  <div class="editor-container">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ isEdit ? '编辑文章' : '新增文章' }}</span>
          <div class="header-actions">
            <el-button @click="router.go(-1)">取消</el-button>
            <el-button type="primary" :loading="saveLoading" @click="submitForm(ruleFormRef)">
              {{ isEdit ? '更新' : '发布' }}
            </el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        label-position="top"
        class="editor-form"
        status-icon
      >
        <el-row :gutter="20">
          <el-col :xs="24" :sm="24" :md="16" :lg="18">
            <el-form-item label="文章标题" prop="title">
              <el-input 
                v-model="ruleForm.title" 
                placeholder="请输入文章标题"
                clearable
                maxlength="200"
                show-word-limit
              />
            </el-form-item>

            <el-form-item label="文章正文" prop="postContent">
              <MdEditor
                v-model="ruleForm.postContent"
                @onHtmlChanged="onHtmlChanged"
                @onGetCatalog="onGetCatalog"
                @onUploadImg="onUploadImg"
                :toolbars="toolbars"
                language="zh-CN"
                :preview="false"
                class="markdown-editor"
              />
              <div class="form-tip">支持 Markdown 语法格式</div>
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="24" :md="8" :lg="6">
            <el-form-item label="文章分类" prop="categoryId">
              <el-select
                v-model="ruleForm.categoryId"
                clearable
                filterable
                placeholder="请选择分类"
                style="width: 100%"
              >
                <el-option
                  v-for="item in categories"
                  :key="item.categoryId"
                  :label="item.categoryName"
                  :value="item.categoryId"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="文章标签" prop="tags">
              <el-select
                v-model="ruleForm.tags"
                style="width: 100%"
                multiple
                filterable
                allow-create
                default-first-option
                :reserve-keyword="false"
                placeholder="选择或创建文章标签"
              >
                <el-option
                  v-for="item in tags"
                  :key="item.tagId"
                  :label="item.tagName"
                  :value="item.tagName"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="文章摘要" prop="summary">
              <el-input
                v-model="ruleForm.summary"
                type="textarea"
                :rows="4"
                placeholder="请输入文章摘要"
                maxlength="500"
                show-word-limit
              />
              <div class="form-tip">如果不填写，系统将自动从正文中提取</div>
            </el-form-item>

            <el-form-item label="文章封面">
              <el-upload
                :action="uploadAction"
                :on-success="handleUploadSuccess"
                :on-error="handleUploadError"
                list-type="picture-card"
                :auto-upload="true"
                :limit="1"
                :file-list="coverFileList"
                accept="image/*"
                :before-upload="beforeUpload"
              >
                <el-icon><Plus /></el-icon>
                <template #file="{ file }">
                  <div class="image-preview">
                    <el-image 
                      :src="file.url" 
                      fit="cover" 
                      class="cover-image"
                    />
                    <span class="image-actions">
                      <span
                        class="image-action-item"
                        @click="handlePictureCardPreview(file)"
                      >
                        <el-icon><zoom-in /></el-icon>
                      </span>
                      <span
                        class="image-action-item"
                        @click="handleRemoveCover(file)"
                      >
                        <el-icon><Delete /></el-icon>
                      </span>
                    </span>
                  </div>
                </template>
              </el-upload>
              <div class="form-tip">建议尺寸：800x600像素</div>
            </el-form-item>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="文章状态" prop="status">
                  <el-radio-group v-model="ruleForm.status" class="radio-group-block">
                    <el-radio :label="1" border>公开</el-radio>
                    <el-radio :label="0" border>草稿</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="文章类型" prop="type">
                  <el-radio-group v-model="ruleForm.type" class="radio-group-block">
                    <el-radio :label="1" border>原创</el-radio>
                    <el-radio :label="0" border>转载</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 图片预览 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="600px">
      <el-image
        :src="previewImageUrl"
        fit="contain"
        style="width: 100%; height: 400px"
      />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
defineOptions({
  name: "Editor"
});

import { getCategoryList } from "@/api/category";
import { getTagList } from "@/api/tag";
import { onMounted, reactive, ref, computed } from "vue";
import type { FormInstance, FormRules, UploadFile } from "element-plus";
import { Delete, Plus, ZoomIn } from "@element-plus/icons-vue";
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { getPost, savePost, updatePost } from "@/api/post";
import { useRouter, useRoute } from "vue-router";
import { upload } from "@/api/common";
import { ElMessage } from "element-plus";
import { message } from "@/utils/message";

const router = useRouter();
const route = useRoute();
const ruleFormRef = ref<FormInstance>();
const saveLoading = ref(false);

// 判断是否为编辑模式
const isEdit = computed(() => !!route.params.id);

const tags = ref([]);
const categories = ref([]);

const ruleForm = ref({
  postId: undefined,
  title: "",
  categoryId: undefined,
  tags: [],
  summary: "",
  status: 1,
  type: 1,
  postContent: "",
  postContentHtml: "",
  coverImage: ""
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

// 上传相关
const uploadAction = import.meta.env.VITE_BASE_URL + "/api/upload/images";
const coverFileList = ref<UploadFile[]>([]);
const previewVisible = ref(false);
const previewImageUrl = ref("");

onMounted(() => {
  // 获取分类和标签列表
  getCategoryList().then(res => {
    categories.value = res.payload;
  });
  getTagList().then(res => {
    tags.value = res.payload;
  });
  
  // 如果是编辑模式，获取文章详情
  const id = route.params.id;
  if (id) {
    getPost(id).then(res => {
      if (res.code === 200) {
        ruleForm.value = res.payload;
        // 如果有封面图片，设置到文件列表中
        if (res.payload.coverImage) {
          coverFileList.value = [{
            name: 'cover.jpg',
            url: res.payload.coverImage
          }] as UploadFile[];
        }
      }
    });
  }
});

const rules = reactive<FormRules>({
  title: [
    { required: true, message: "请输入标题", trigger: "blur" },
    { min: 5, max: 200, message: "标题长度应在5-200字符之间", trigger: "blur" }
  ],
  categoryId: [{ required: true, message: "请选择分类", trigger: "change" }],
  tags: [{ required: true, message: "请至少选择一个标签", trigger: "change" }],
  status: [{ required: true, message: "请选择文章状态", trigger: "change" }],
  type: [{ required: true, message: "请选择文章类型", trigger: "change" }],
  postContent: [{ required: true, message: "请输入文章正文", trigger: "blur" }]
});

const onGetCatalog = (_e: any) => {
  // 目录获取回调（如需处理目录可在此添加逻辑）
};

const onHtmlChanged = (html: string) => {
  ruleForm.value.postContentHtml = html;
};

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      saveLoading.value = true;
      try {
        let res;
        if (isEdit.value) {
          res = await updatePost(ruleForm.value);
        } else {
          res = await savePost(ruleForm.value);
        }
        
        if (res.code === 200) {
          message(`${isEdit.value ? '更新' : '发布'}文章成功`, { type: "success" });
          router.push({ name: "内容管理" });
        } else {
          message(`${isEdit.value ? '更新' : '发布'}文章失败: ${res.message || "未知错误"}`, { type: "error" });
        }
      } catch (error) {
        message(`操作过程中发生错误: ${error.message || "未知错误"}`, { type: "error" });
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

// 上传前检查
const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/');
  const isLt2M = file.size / 1024 / 1024 < 2;
  
  if (!isImage) {
    message('只能上传图片文件!', { type: "error" });
    return false;
  }
  
  if (!isLt2M) {
    message('图片大小不能超过 2MB!', { type: "error" });
    return false;
  }
  
  return true;
};

// 上传成功处理
const handleUploadSuccess = (response: any, file: UploadFile) => {
  if (response.code === 200) {
    ruleForm.value.coverImage = response.payload;
    message("封面上传成功", { type: "success" });
  } else {
    message(`上传失败: ${response.message || "未知错误"}`, { type: "error" });
  }
};

// 上传失败处理
const handleUploadError = (error: any, file: UploadFile) => {
  message(`上传失败: ${error.message || "网络错误"}`, { type: "error" });
};

// 移除封面
const handleRemoveCover = (file: UploadFile) => {
  ruleForm.value.coverImage = "";
  coverFileList.value = [];
  message("已移除封面图片", { type: "success" });
};

// 预览图片
const handlePictureCardPreview = (file: UploadFile) => {
  previewImageUrl.value = file.url!;
  previewVisible.value = true;
};

// 重置表单
const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields();
};
</script>

<style scoped lang="scss">
.editor-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);
  
  .editor-card {
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
        
        .header-actions {
          display: flex;
          gap: 12px;
        }
      }
    }
  }
  
  .editor-form {
    :deep(.el-form-item) {
      margin-bottom: 24px;
      
      .el-form-item__label {
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin-bottom: 8px;
      }
    }
    
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
    
    .radio-group-block {
      display: flex;
      flex-direction: column;
      gap: 10px;
      
      .el-radio {
        margin-right: 0;
      }
    }
  }
  
  .image-preview {
    position: relative;
    width: 100%;
    height: 100%;
    
    .cover-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    
    .image-actions {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.5);
      display: flex;
      justify-content: center;
      align-items: center;
      gap: 20px;
      opacity: 0;
      transition: opacity 0.3s;
      
      .image-action-item {
        font-size: 20px;
        color: #fff;
        cursor: pointer;
      }
    }
    
    &:hover {
      .image-actions {
        opacity: 1;
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
  .editor-container {
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
      
      .header-actions {
        align-self: flex-end;
      }
    }
  }
}
</style>