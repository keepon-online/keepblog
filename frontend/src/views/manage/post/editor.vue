<template>
  <div class="editor-container" :class="{ 'focus-mode': focusMode }" @keydown.esc.exact="exitFocusMode">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ isEdit ? '编辑文章' : '新增文章' }}</span>
          <div class="header-actions">
            <el-button
              :type="focusMode ? 'primary' : 'default'"
              :title="focusMode ? '退出专注模式 (Esc)' : '专注模式'"
              @click="toggleFocusMode"
            >
              <el-icon class="mr-2px"><Aim /></el-icon>
              {{ focusMode ? '退出专注' : '专注模式' }}
            </el-button>
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
          <el-col v-bind="mainColProps">
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
                @onSave="handleSaveDraft"
                :toolbars="toolbars"
                language="zh-CN"
                :preview="false"
                :style="editorStyle"
                class="markdown-editor"
              />
              <div class="form-tip">
                支持 Markdown 语法格式
                <span v-if="draftSavedAt" class="draft-indicator">
                  · 草稿已自动保存于 {{ draftSavedAtText }}
                </span>
              </div>
            </el-form-item>
          </el-col>

          <el-col v-if="!focusMode" :xs="24" :sm="24" :md="8" :lg="6">
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
              <el-col :xs="24" :span="24">
                <el-form-item label="文章状态" prop="status">
                  <el-radio-group v-model="ruleForm.status" class="radio-group-block">
                    <el-radio :label="1" border>公开</el-radio>
                    <el-radio :label="0" border>草稿</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :span="24">
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

    <!-- 底部固定操作栏：长文写到结尾无需滚回顶部即可发布 -->
    <div class="editor-footer">
      <span v-if="draftSavedAtText" class="draft-indicator">
        <el-icon><DocumentChecked /></el-icon>
        草稿已保存 {{ draftSavedAtText }}
      </span>
      <span class="footer-spacer" />
      <el-button @click="router.go(-1)">取消</el-button>
      <el-button :loading="draftSaving" @click="handleSaveDraft">
        <el-icon class="mr-2px"><Document /></el-icon>
        存草稿
      </el-button>
      <el-button type="primary" :loading="saveLoading" @click="submitForm(ruleFormRef)">
        {{ isEdit ? '更新' : '发布' }}
      </el-button>
    </div>

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
import { onMounted, onBeforeUnmount, reactive, ref, computed, watch } from "vue";
import type { FormInstance, FormRules, UploadFile } from "element-plus";
import { ElMessageBox } from "element-plus";
import {
  Aim,
  Delete,
  Document,
  DocumentChecked,
  Plus,
  ZoomIn
} from "@element-plus/icons-vue";
import { MdEditor, type ToolbarNames } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { getPost, savePost, updatePost } from "@/api/post";
import { useRouter, useRoute } from "vue-router";
import { upload } from "@/api/common";
import { message } from "@/utils/message";
import { localForage } from "@/utils/localforage";
import { useDebounceFn } from "@vueuse/core";

const router = useRouter();
const route = useRoute();
const ruleFormRef = ref<FormInstance>();
const saveLoading = ref(false);

// 判断是否为编辑模式
const isEdit = computed(() => !!route.params.id);

// 专注模式：隐藏右栏元信息，左栏编辑区撑满
const focusMode = ref(false);

// 左栏响应式 props：专注模式下撑满整行
const mainColProps = computed(() =>
  focusMode.value
    ? { xs: 24, sm: 24, md: 24, lg: 24 }
    : { xs: 24, sm: 24, md: 16, lg: 18 }
);

// 编辑器高度：专注模式下更高（无右栏挤压）
const editorStyle = computed(() => {
  const offset = focusMode.value ? 260 : 320;
  return {
    height: `calc(100vh - ${offset}px)`,
    minHeight: focusMode.value ? "600px" : "500px"
  };
});

// ---------- 自动草稿（按文章 ID 隔离，防抖 2s）----------
interface PostDraft {
  ruleForm: typeof ruleForm.value;
  savedAt: number;
}

const DRAFT_KEY = computed(
  () => `post-draft:${(route.params.id as string) || "new"}`
);
const draftSaving = ref(false); // 手动「存草稿」按钮 loading
const draftSavedAt = ref<number | null>(null); // 自动保存时间戳，用于显示
const draftRestored = ref(false); // 草稿恢复流程完成才开始自动存，避免覆盖刚加载的数据
const draftSavedAtText = computed(() =>
  draftSavedAt.value
    ? new Date(draftSavedAt.value).toLocaleTimeString("zh-CN", {
        hour: "2-digit",
        minute: "2-digit"
      })
    : ""
);

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
const toolbars: ToolbarNames[] = [
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

// ---------- 草稿：写入/读取/清除 ----------
const persistDraftImmediately = async () => {
  // 标题与正文都基本为空时，不写入无意义空草稿
  if (
    !ruleForm.value.title &&
    ruleForm.value.postContent.trim().length < 10
  ) {
    return false;
  }
  draftSaving.value = true;
  try {
    await localForage().setItem<PostDraft>(DRAFT_KEY.value, {
      ruleForm: JSON.parse(JSON.stringify(ruleForm.value)),
      savedAt: Date.now()
    });
    draftSavedAt.value = Date.now();
    return true;
  } catch (e) {
    console.warn("[post-draft] 保存草稿失败", e);
    return false;
  } finally {
    draftSaving.value = false;
  }
};

// 防抖 2s 自动保存（监听 ruleForm 深度变化触发）
const persistDraftDebounced = useDebounceFn(persistDraftImmediately, 2000);

// 手动「存草稿」按钮 + 工具栏 onSave 触发（立即保存并提示）
const handleSaveDraft = async () => {
  const ok = await persistDraftImmediately();
  message(ok ? "草稿已保存到本地" : "内容太少，未保存草稿", {
    type: ok ? "success" : "info"
  });
};

// 检测并提示恢复草稿（编辑模式下仅当草稿比远程数据新才提示）
const tryRestoreDraft = async (remoteUpdateTime?: number) => {
  let draft: PostDraft | null = null;
  try {
    draft = await localForage().getItem<PostDraft>(DRAFT_KEY.value);
  } catch (e) {
    console.warn("[post-draft] 读取草稿失败", e);
  }
  if (!draft || !draft.ruleForm) {
    draftRestored.value = true;
    return;
  }
  // 编辑模式且远程数据更新时间 >= 草稿时间：草稿已过期，静默丢弃
  if (remoteUpdateTime && remoteUpdateTime >= draft.savedAt) {
    await localForage().removeItem(DRAFT_KEY.value);
    draftRestored.value = true;
    return;
  }
  const minutes = Math.max(
    1,
    Math.round((Date.now() - draft.savedAt) / 60000)
  );
  try {
    await ElMessageBox.confirm(
      `检测到 ${minutes} 分钟前的未保存草稿，是否恢复？选择「丢弃」将删除该草稿。`,
      "恢复草稿",
      {
        confirmButtonText: "恢复",
        cancelButtonText: "丢弃",
        type: "info",
        distinguishCancelAndClose: true,
        closeOnClickModal: false
      }
    );
    // 用户选择「恢复」
    ruleForm.value = { ...ruleForm.value, ...draft.ruleForm };
    draftSavedAt.value = draft.savedAt;
    if (draft.ruleForm.coverImage) {
      coverFileList.value = [
        { name: "cover.jpg", url: draft.ruleForm.coverImage } as UploadFile
      ];
    }
    message("已恢复草稿", { type: "success" });
  } catch (action) {
    // 用户选择「丢弃」或关闭
    if (action === "cancel") {
      await localForage().removeItem(DRAFT_KEY.value);
      message("已丢弃草稿", { type: "info" });
    }
    // close（点 X / Esc）时保留草稿不处理
  } finally {
    draftRestored.value = true;
  }
};

// ---------- 专注模式 ----------
const toggleFocusMode = () => {
  focusMode.value = !focusMode.value;
};
const exitFocusMode = () => {
  if (focusMode.value) focusMode.value = false;
};

// 监听表单变化自动存草稿（深度，防抖 2s）
watch(
  ruleForm,
  () => {
    if (draftRestored.value) persistDraftDebounced();
  },
  { deep: true }
);

// 组件卸载前若仍有未提交改动，立刻存一次（防止用户直接关闭路由）
onBeforeUnmount(() => {
  persistDraftImmediately();
});

onMounted(() => {
  // 获取分类和标签列表
  getCategoryList().then(res => {
    categories.value = res.payload;
  });
  getTagList().then(res => {
    tags.value = res.payload;
  });

  // 如果是编辑模式，获取文章详情，加载完后再尝试恢复草稿
  const id = route.params.id;
  if (id) {
    getPost(id)
      .then(res => {
        if (res.code === 200) {
          ruleForm.value = res.payload;
          // 如果有封面图片，设置到文件列表中
          if (res.payload.coverImage) {
            coverFileList.value = [
              {
                name: "cover.jpg",
                url: res.payload.coverImage
              }
            ] as UploadFile[];
          }
        }
        // 编辑模式：草稿需比远程 updateTime 更新才提示恢复
        const remoteUpdateTime = res?.payload?.updateTime
          ? new Date(res.payload.updateTime).getTime()
          : undefined;
        return tryRestoreDraft(remoteUpdateTime);
      })
      .catch(() => tryRestoreDraft());
  } else {
    // 新增模式：直接尝试恢复草稿
    tryRestoreDraft();
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

  // 专注模式下校验失败时，右栏元信息被隐藏，先退出专注模式让用户看到错误字段
  const handleInvalid = () => {
    if (focusMode.value) focusMode.value = false;
    message("表单验证失败，请检查输入内容", { type: "warning" });
  };

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
          // 发布成功：清除本地草稿（草稿使命完成）
          try {
            await localForage().removeItem(DRAFT_KEY.value);
          } catch (e) {
            console.warn("[post-draft] 清除草稿失败", e);
          }
          message(`${isEdit.value ? '更新' : '发布'}文章成功`, { type: "success" });
          router.push({ name: "内容管理" });
        } else {
          message(`${isEdit.value ? '更新' : '发布'}文章失败: ${res.message || "未知错误"}`, { type: "error" });
        }
      } catch (error) {
        const errMsg = error instanceof Error ? error.message : String(error);
        message(`操作过程中发生错误: ${errMsg}`, { type: "error" });
      } finally {
        saveLoading.value = false;
      }
    } else {
      handleInvalid();
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
  const urls = res.map((item: string | { url?: string }) => {
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

    .draft-indicator {
      color: var(--el-color-success);
      margin-left: 4px;
    }
  }

  // 草稿已保存指示（底部 footer 用）
  .draft-indicator {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--el-color-success);
  }

  // 底部固定操作栏
  .editor-footer {
    position: sticky;
    bottom: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 16px;
    padding: 12px 20px;
    background-color: var(--el-bg-color);
    border: 1px solid var(--el-border-color-light);
    border-radius: 8px;
    box-shadow: 0 -2px 12px 0 rgba(0, 0, 0, 0.06);

    .footer-spacer {
      flex: 1;
    }
  }

  // 专注模式：让编辑器区域成为视觉重心
  &.focus-mode {
    .editor-card {
      :deep(.el-card__body) {
        padding: 16px 20px;
      }
    }
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

    // 小屏下 footer 文案隐藏，仅保留按钮
    .editor-footer {
      .draft-indicator {
        display: none;
      }

      .footer-spacer {
        display: none;
      }
    }
  }
}
</style>