<template>
  <div
    class="editor-container"
    :class="{ 'focus-mode': focusMode }"
    @keydown.esc.exact="exitFocusMode"
  >
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ isEdit ? "编辑文章" : "新增文章" }}</span>
          <div class="header-actions">
            <el-button
              :type="focusMode ? 'primary' : 'default'"
              :title="focusMode ? '退出专注模式 (Esc)' : '专注模式'"
              @click="toggleFocusMode"
            >
              <el-icon class="mr-2px"><Aim /></el-icon>
              {{ focusMode ? "退出专注" : "专注模式" }}
            </el-button>
            <el-button
              title="从本地 .md 文件导入（支持 FrontMatter 元数据解析）"
              @click="triggerImportMd"
            >
              <el-icon class="mr-2px"><Upload /></el-icon>
              导入 .md
            </el-button>
            <el-button
              title="将当前文章与元数据导出为标准 Markdown 文件"
              @click="handleExportMd"
            >
              <el-icon class="mr-2px"><Download /></el-icon>
              导出 .md
            </el-button>
            <el-button @click="router.go(-1)">取消</el-button>
            <el-button
              type="primary"
              :loading="saveLoading"
              @click="submitForm(ruleFormRef)"
            >
              {{ isEdit ? "更新" : "发布" }}
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
                :toolbars="toolbars"
                language="zh-CN"
                :preview="true"
                :sync-scroll="true"
                code-theme="atom"
                :style="editorStyle"
                class="markdown-editor"
                @onHtmlChanged="onHtmlChanged"
                @onGetCatalog="onGetCatalog"
                @onUploadImg="onUploadImg"
                @onSave="handleSaveDraft"
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

            <el-form-item label="文章封面" class="cover-form-item">
              <div class="cover-manager-card">
                <!-- 封面预览区域 -->
                <div class="cover-preview-box">
                  <template v-if="ruleForm.coverImage">
                    <el-image
                      :src="ruleForm.coverImage"
                      fit="cover"
                      class="cover-image-display"
                    />
                    <div class="cover-overlay">
                      <el-tooltip content="预览大图" placement="top">
                        <span class="overlay-btn" @click="handlePreviewCover">
                          <el-icon><ZoomIn /></el-icon>
                        </span>
                      </el-tooltip>
                      <el-tooltip
                        content="设为无封面（前台展示星空流星效果）"
                        placement="top"
                      >
                        <span
                          class="overlay-btn danger"
                          @click="handleClearCover"
                        >
                          <el-icon><Delete /></el-icon>
                        </span>
                      </el-tooltip>
                    </div>
                  </template>
                  <template v-else>
                    <div class="cover-empty-state">
                      <div class="starry-badge">
                        <el-icon><Picture /></el-icon>
                        <span>无封面状态</span>
                      </div>
                      <span class="starry-hint"
                        >前台将自动以动态星空流星效果展示</span
                      >
                    </div>
                  </template>
                </div>

                <!-- 封面控制操作区：支持4种模式 -->
                <div class="cover-toolbar">
                  <!-- 方式1：本地/OSS 上传 -->
                  <el-upload
                    :action="uploadAction"
                    :show-file-list="false"
                    :auto-upload="true"
                    accept="image/*"
                    :before-upload="beforeUpload"
                    :on-success="handleUploadSuccess"
                    :on-error="handleUploadError"
                    class="upload-trigger"
                  >
                    <el-button size="small" type="primary" plain>
                      <el-icon class="mr-2px"><Upload /></el-icon>
                      上传图片
                    </el-button>
                  </el-upload>

                  <!-- 方式2：输入网络图片 URL -->
                  <el-popover
                    v-model:visible="urlPopoverVisible"
                    placement="bottom"
                    :width="280"
                    trigger="click"
                  >
                    <template #reference>
                      <el-button size="small" plain>
                        <el-icon class="mr-2px"><Link /></el-icon>
                        外链地址
                      </el-button>
                    </template>
                    <div class="url-input-popover">
                      <el-input
                        v-model="customImageUrl"
                        placeholder="请输入图片 URL (https://...)"
                        size="small"
                        clearable
                        @keydown.enter="applyCustomImageUrl"
                      />
                      <div class="popover-actions">
                        <el-button
                          size="small"
                          @click="urlPopoverVisible = false"
                          >取消</el-button
                        >
                        <el-button
                          size="small"
                          type="primary"
                          @click="applyCustomImageUrl"
                          >确定</el-button
                        >
                      </div>
                    </div>
                  </el-popover>

                  <!-- 方式3：随机 Pixabay 封面 -->
                  <el-button
                    size="small"
                    plain
                    :loading="randomCoverLoading"
                    title="从 Pixabay 获取一张高清壁纸"
                    @click="handleRandomCover"
                  >
                    <el-icon class="mr-2px"><MagicStick /></el-icon>
                    随机壁纸
                  </el-button>

                  <!-- 方式4：设为无封面（星空流星模式） -->
                  <el-button
                    v-if="ruleForm.coverImage"
                    size="small"
                    type="danger"
                    link
                    @click="handleClearCover"
                  >
                    设为无封面
                  </el-button>
                </div>
              </div>
              <div class="form-tip">
                支持上传图片、输入外链、一键随机配图或留空以启用前台星空流星效果
              </div>
            </el-form-item>

            <el-row :gutter="20">
              <el-col :xs="24" :span="24">
                <el-form-item label="文章状态" prop="status">
                  <el-radio-group
                    v-model="ruleForm.status"
                    class="radio-group-block"
                  >
                    <el-radio :label="1" border>公开</el-radio>
                    <el-radio :label="0" border>草稿</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :span="24">
                <el-form-item label="文章类型" prop="type">
                  <el-radio-group
                    v-model="ruleForm.type"
                    class="radio-group-block"
                  >
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
      <el-button
        type="primary"
        :loading="saveLoading"
        @click="submitForm(ruleFormRef)"
      >
        {{ isEdit ? "更新" : "发布" }}
      </el-button>
    </div>

    <!-- 隐藏的 Markdown 文件导入 input -->
    <input
      ref="mdFileInputRef"
      type="file"
      accept=".md,.markdown"
      style="display: none"
      @change="onMdFileSelected"
    />

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
import {
  onMounted,
  onBeforeUnmount,
  reactive,
  ref,
  computed,
  watch
} from "vue";
import type { FormInstance, FormRules, UploadFile } from "element-plus";
import { ElMessageBox } from "element-plus";
import {
  Aim,
  Delete,
  Document,
  DocumentChecked,
  Download,
  Link,
  MagicStick,
  Picture,
  Plus,
  Upload,
  ZoomIn
} from "@element-plus/icons-vue";
import { MdEditor, type ToolbarNames } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { getPost, getRandomCover, savePost, updatePost } from "@/api/post";
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

const tags = ref<any[]>([]);
const categories = ref<any[]>([]);

const ruleForm = ref({
  postId: undefined,
  title: "",
  categoryId: undefined,
  tags: [] as string[],
  summary: "",
  status: 1,
  type: 1,
  postContent: "",
  postContentHtml: "",
  coverImage: ""
});

// Markdown 文件导入与导出
const mdFileInputRef = ref<HTMLInputElement | null>(null);

// 封面管理相关状态（4种模式：本地上传/网络外链/随机壁纸/清空星空）
const uploadAction = import.meta.env.VITE_BASE_URL + "/api/upload/images";
const coverFileList = ref<UploadFile[]>([]);
const previewVisible = ref(false);
const previewImageUrl = ref("");
const randomCoverLoading = ref(false);
const urlPopoverVisible = ref(false);
const customImageUrl = ref("");

// Markdown 编辑器工具栏配置（支持分屏实时预览、全屏、HTML预览、大纲目录）
const toolbars: ToolbarNames[] = [
  "bold",
  "underline",
  "italic",
  "-",
  "title",
  "strikeThrough",
  "sub",
  "sup",
  "-",
  "quote",
  "unorderedList",
  "orderedList",
  "task",
  "-",
  "codeRow",
  "code",
  "link",
  "image",
  "table",
  "mermaid",
  "katex",
  "-",
  "revoke",
  "next",
  "save",
  "=",
  "pageFullscreen",
  "fullscreen",
  "preview",
  "htmlPreview",
  "catalog"
];

// ---------- 草稿：写入/读取/清除 ----------
const persistDraftImmediately = async () => {
  // 标题与正文都基本为空时，不写入无意义空草稿
  if (!ruleForm.value.title && ruleForm.value.postContent.trim().length < 10) {
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
  const minutes = Math.max(1, Math.round((Date.now() - draft.savedAt) / 60000));
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
          message(`${isEdit.value ? "更新" : "发布"}文章成功`, {
            type: "success"
          });
          router.push({ name: "内容管理" });
        } else {
          message(
            `${isEdit.value ? "更新" : "发布"}文章失败: ${res.message || "未知错误"}`,
            { type: "error" }
          );
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

const onUploadImg = async (
  files: File[],
  callback: (urls: string[]) => void
) => {
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
    if (typeof item === "string") {
      return item;
    } else if (item && typeof item === "object" && item.url) {
      return item.url;
    }
    return "";
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
  const isImage = file.type.startsWith("image/");
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isImage) {
    message("只能上传图片文件!", { type: "error" });
    return false;
  }

  if (!isLt2M) {
    message("图片大小不能超过 2MB!", { type: "error" });
    return false;
  }

  return true;
};

// 封面图片上传成功处理
const handleUploadSuccess = (response: any) => {
  if (response.code === 200) {
    ruleForm.value.coverImage = response.payload;
    coverFileList.value = [
      { name: "cover.jpg", url: response.payload } as UploadFile
    ];
    message("封面上传成功", { type: "success" });
  } else {
    message(`上传失败: ${response.message || "未知错误"}`, { type: "error" });
  }
};

// 上传失败处理
const handleUploadError = (error: any) => {
  message(`上传失败: ${error.message || "网络错误"}`, { type: "error" });
};

// 预览当前封面大图
const handlePreviewCover = () => {
  if (!ruleForm.value.coverImage) return;
  previewImageUrl.value = ruleForm.value.coverImage;
  previewVisible.value = true;
};

// 清除封面（设为无封面，触发前台动态星空流星效果）
const handleClearCover = () => {
  ruleForm.value.coverImage = "";
  coverFileList.value = [];
  message("已设为无封面（前台将展示动态星空流星效果）", { type: "info" });
};

// 应用用户输入的网络图片外链
const applyCustomImageUrl = () => {
  const url = customImageUrl.value.trim();
  if (!url) {
    message("请输入有效的图片 URL", { type: "warning" });
    return;
  }
  ruleForm.value.coverImage = url;
  coverFileList.value = [{ name: "cover.jpg", url } as UploadFile];
  customImageUrl.value = "";
  urlPopoverVisible.value = false;
  message("已应用封面图片外链", { type: "success" });
};

// 随机获取一张 Pixabay 高清壁纸封面
const handleRandomCover = async () => {
  randomCoverLoading.value = true;
  try {
    const res = await getRandomCover();
    if (res.code === 200 && res.payload) {
      ruleForm.value.coverImage = res.payload;
      coverFileList.value = [
        { name: "cover.jpg", url: res.payload } as UploadFile
      ];
      message("已成功匹配一张高清随机封面壁纸", { type: "success" });
    } else {
      message(`获取随机封面失败: ${res.message || "服务异常"}`, {
        type: "warning"
      });
    }
  } catch (err: any) {
    message(`获取随机封面异常: ${err?.message || "网络错误"}`, {
      type: "error"
    });
  } finally {
    randomCoverLoading.value = false;
  }
};

// ---------- Markdown 导入与导出 ----------
const triggerImportMd = () => {
  mdFileInputRef.value?.click();
};

// 解析 Markdown YAML FrontMatter
const parseFrontMatter = (content: string) => {
  const match = content.match(/^---\r?\n([\s\S]*?)\r?\n---\r?\n?([\s\S]*)$/);
  if (!match) {
    return { meta: {} as Record<string, any>, body: content };
  }
  const yamlText = match[1];
  const body = match[2];
  const meta: Record<string, any> = {};

  const lines = yamlText.split(/\r?\n/);
  let currentKey = "";
  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) continue;

    // 列表项: - item
    const listMatch = line.match(/^\s*-\s+(.+)$/);
    if (listMatch) {
      const val = listMatch[1].trim().replace(/^['"]|['"]$/g, "");
      if (currentKey && Array.isArray(meta[currentKey])) {
        meta[currentKey].push(val);
      }
      continue;
    }

    // 键值对: key: value
    const kvMatch = line.match(/^([a-zA-Z0-9_-]+):\s*(.*)$/);
    if (kvMatch) {
      currentKey = kvMatch[1].trim();
      const rawVal = kvMatch[2].trim();

      if (rawVal === "") {
        meta[currentKey] = [];
      } else if (rawVal.startsWith("[") && rawVal.endsWith("]")) {
        meta[currentKey] = rawVal
          .slice(1, -1)
          .split(",")
          .map(s => s.trim().replace(/^['"]|['"]$/g, ""))
          .filter(Boolean);
      } else {
        meta[currentKey] = rawVal.replace(/^['"]|['"]$/g, "");
      }
    }
  }
  return { meta, body };
};

// 选择并解析本地 .md 文件
const onMdFileSelected = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = e => {
    const text = e.target?.result as string;
    if (!text) return;

    const { meta, body } = parseFrontMatter(text);

    // 1. 文章标题
    if (meta.title) {
      ruleForm.value.title = String(meta.title).trim();
    } else {
      const h1Match = text.match(/^#\s+(.+)$/m);
      if (h1Match) {
        ruleForm.value.title = h1Match[1].trim();
      } else {
        ruleForm.value.title = file.name.replace(/\.(md|markdown)$/i, "");
      }
    }

    // 2. 正文内容
    ruleForm.value.postContent = body || text;

    // 3. 文章标签
    if (Array.isArray(meta.tags)) {
      ruleForm.value.tags = meta.tags.map(String);
    } else if (typeof meta.tags === "string" && meta.tags.trim()) {
      ruleForm.value.tags = meta.tags
        .split(/[,，;；]/)
        .map((s: string) => s.trim())
        .filter(Boolean);
    }

    // 4. 分类自动匹配
    const catVal =
      meta.category ||
      (Array.isArray(meta.categories) ? meta.categories[0] : meta.categories);
    if (catVal && categories.value.length > 0) {
      const found = categories.value.find(
        (c: any) =>
          c.categoryName.toLowerCase() === String(catVal).trim().toLowerCase()
      );
      if (found) {
        ruleForm.value.categoryId = found.categoryId;
      }
    }

    // 5. 文章摘要
    if (meta.summary || meta.description) {
      ruleForm.value.summary = String(meta.summary || meta.description).trim();
    }

    // 6. 封面图
    const cover = meta.coverImage || meta.cover || meta.image;
    if (cover) {
      ruleForm.value.coverImage = String(cover).trim();
      coverFileList.value = [
        { name: "cover.jpg", url: ruleForm.value.coverImage } as UploadFile
      ];
    }

    // 7. 文章状态与类型
    if (meta.draft === true || meta.status === 0 || meta.status === "draft") {
      ruleForm.value.status = 0;
    } else if (meta.status === 1 || meta.status === "published") {
      ruleForm.value.status = 1;
    }

    if (
      meta.type === "reproduced" ||
      meta.type === "转载" ||
      meta.type === 0 ||
      meta.type === "0"
    ) {
      ruleForm.value.type = 0;
    } else if (
      meta.type === "original" ||
      meta.type === "原创" ||
      meta.type === 1 ||
      meta.type === "1"
    ) {
      ruleForm.value.type = 1;
    }

    message(`成功导入「${file.name}」，已自动解析元数据与正文`, {
      type: "success"
    });
    // 清空 input 允许重复选择相同文件名
    target.value = "";
  };
  reader.onerror = () => {
    message("读取 Markdown 文件失败", { type: "error" });
    target.value = "";
  };
  reader.readAsText(file, "UTF-8");
};

// 导出当前文章与元数据为 Markdown 文件
const handleExportMd = () => {
  const form = ruleForm.value;
  const currentCat = categories.value.find(
    (c: any) => c.categoryId === form.categoryId
  );
  const catName = currentCat ? currentCat.categoryName : "";

  let yaml = "---\n";
  yaml += `title: "${(form.title || "未命名文章").replace(/"/g, '\\"')}"\n`;
  yaml += `date: ${new Date().toISOString()}\n`;
  if (catName) {
    yaml += `category: "${catName}"\n`;
  }
  if (form.tags && form.tags.length > 0) {
    yaml += `tags:\n${form.tags.map((t: string) => `  - "${t}"`).join("\n")}\n`;
  }
  if (form.summary) {
    yaml += `summary: "${form.summary.replace(/"/g, '\\"')}"\n`;
  }
  if (form.coverImage) {
    yaml += `coverImage: "${form.coverImage}"\n`;
  }
  yaml += `type: ${form.type === 1 ? "original" : "reproduced"}\n`;
  yaml += `status: ${form.status === 1 ? "published" : "draft"}\n`;
  yaml += "---\n\n";

  const fullMd = yaml + (form.postContent || "");
  const blob = new Blob([fullMd], { type: "text/markdown;charset=utf-8" });
  const downloadUrl = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = downloadUrl;
  const safeTitle = (form.title || "post")
    .trim()
    .replace(/[/\\?%*:|"<>]/g, "_");
  link.download = `${safeTitle}.md`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(downloadUrl);
  message("文章已成功导出为 Markdown 文件", { type: "success" });
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
      border-radius: 6px;
      transition: var(--el-transition-border);

      &:hover {
        border-color: var(--el-border-color-hover);
      }

      &:focus-within {
        border-color: var(--el-color-primary);
      }

      // 预览区保真排版与高亮样式定制（对齐前台博客视觉规范）
      :deep(.md-editor-preview) {
        font-family:
          -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC",
          "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
        line-height: 1.8;
        color: var(--el-text-color-primary);

        // 标题层次与锚点感
        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
          font-weight: 600;
          color: var(--el-text-color-primary);
          margin-top: 24px;
          margin-bottom: 12px;
          line-height: 1.4;
        }

        h1 {
          font-size: 24px;
          border-bottom: 2px solid var(--el-border-color-light);
          padding-bottom: 8px;
        }
        h2 {
          font-size: 20px;
          border-bottom: 1px solid var(--el-border-color-lighter);
          padding-bottom: 6px;
        }
        h3 {
          font-size: 17px;
        }
        h4 {
          font-size: 15px;
        }

        // 引用块（与前台一致的左侧品牌色强调边框与浅色底）
        blockquote {
          margin: 16px 0;
          padding: 12px 18px;
          border-left: 4px solid var(--el-color-primary);
          background-color: var(--el-color-primary-light-9);
          border-radius: 0 6px 6px 0;
          color: var(--el-text-color-regular);

          p {
            margin: 0;
          }
        }

        // 行内代码
        :not(pre) > code {
          background-color: var(--el-fill-color);
          color: var(--el-color-primary);
          padding: 2px 6px;
          border-radius: 4px;
          font-size: 0.9em;
          font-family:
            "JetBrains Mono", "Fira Code", Consolas, Monaco, monospace;
        }

        // 代码块：Monokai 暗黑基调 + macOS 红黄绿红绿灯控制点
        pre {
          position: relative;
          background-color: #212121 !important;
          border-radius: 8px !important;
          padding: 34px 16px 16px !important;
          margin: 16px 0;
          overflow-x: auto;
          box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);

          // macOS 窗口红黄绿三色圆点
          &::before {
            content: "";
            position: absolute;
            top: 12px;
            left: 14px;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #ff5f56;
            box-shadow:
              16px 0 0 #ffbd2e,
              32px 0 0 #27c93f;
          }

          code {
            font-family:
              "JetBrains Mono", "Fira Code", Consolas, Monaco, monospace;
            font-size: 13.5px;
            line-height: 1.65;
            color: #eff;
            background: transparent !important;
          }
        }

        // 表格样式强化
        table {
          width: 100%;
          border-collapse: collapse;
          margin: 16px 0;
          border-radius: 6px;
          overflow: hidden;
          font-size: 14px;

          th {
            background-color: var(--el-fill-color-light);
            font-weight: 600;
            padding: 10px 14px;
            border: 1px solid var(--el-border-color-lighter);
            text-align: left;
          }

          td {
            padding: 8px 14px;
            border: 1px solid var(--el-border-color-lighter);
          }

          tr:nth-child(even) {
            background-color: var(--el-fill-color-extra-light);
          }
        }

        // 图片圆角与阴影
        img {
          max-width: 100%;
          border-radius: 6px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
          margin: 12px 0;
        }

        // 分割线
        hr {
          border: none;
          border-top: 1px dashed var(--el-border-color);
          margin: 24px 0;
        }
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

    // 封面图管理器（4种模式：本地/外链/随机/星空）
    .cover-manager-card {
      display: flex;
      flex-direction: column;
      gap: 10px;

      .cover-preview-box {
        position: relative;
        width: 100%;
        height: 140px;
        border-radius: 8px;
        overflow: hidden;
        border: 1px solid var(--el-border-color);
        background-color: var(--el-fill-color-lighter);

        .cover-image-display {
          width: 100%;
          height: 100%;
          object-fit: cover;
          display: block;
        }

        .cover-overlay {
          position: absolute;
          inset: 0;
          background: rgba(0, 0, 0, 0.55);
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 16px;
          opacity: 0;
          transition: opacity 0.25s ease;

          .overlay-btn {
            width: 34px;
            height: 34px;
            border-radius: 50%;
            background: rgba(255, 255, 255, 0.2);
            color: #fff;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 18px;
            cursor: pointer;
            backdrop-filter: blur(4px);
            transition: all 0.2s ease;

            &:hover {
              background: rgba(255, 255, 255, 0.4);
              transform: scale(1.1);
            }

            &.danger:hover {
              background: var(--el-color-danger);
            }
          }
        }

        &:hover .cover-overlay {
          opacity: 1;
        }

        // 星空无封面占位状态
        .cover-empty-state {
          width: 100%;
          height: 100%;
          background: linear-gradient(
            135deg,
            #0f172a 0%,
            #1e1b4b 50%,
            #172554 100%
          );
          color: #e2e8f0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          gap: 6px;
          padding: 12px;
          text-align: center;

          .starry-badge {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            font-weight: 500;
            color: #93c5fd;
            background: rgba(255, 255, 255, 0.08);
            padding: 4px 12px;
            border-radius: 20px;
            backdrop-filter: blur(4px);
          }

          .starry-hint {
            font-size: 11px;
            color: #94a3b8;
          }
        }
      }

      .cover-toolbar {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 8px;

        .upload-trigger {
          display: inline-flex;
        }
      }
    }
  }

  .url-input-popover {
    display: flex;
    flex-direction: column;
    gap: 10px;

    .popover-actions {
      display: flex;
      justify-content: flex-end;
      gap: 8px;
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
