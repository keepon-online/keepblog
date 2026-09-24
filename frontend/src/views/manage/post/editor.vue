<template>
  <div
    class="editor-container"
    :class="{ 'focus-mode': focusMode }"
    @keydown.esc.exact="handleEscKey"
  >
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ isEdit ? "编辑文章" : "新增文章" }}</span>
          <div class="header-actions">
            <el-button
              v-if="aiEnabled"
              type="success"
              plain
              title="智能发文助手与 SEO 质量体检"
              @click="openPublishModal"
            >
              <el-icon class="mr-2px"><MagicStick /></el-icon>
              发文体检与 AI 助手
            </el-button>
            <el-button
              v-if="aiEnabled"
              type="primary"
              plain
              title="根据主题与读者深度智能生成多级文章大纲与写作要点"
              @click="openOutlineDialog"
            >
              <el-icon class="mr-2px"><DocumentCopy /></el-icon>
              ✨ 生成大纲
            </el-button>
            <el-button
              title="查看本地修订快照与历史差异对比"
              @click="openHistoryModal"
            >
              <el-icon class="mr-2px"><Tickets /></el-icon>
              ⏳ 草稿时光机
            </el-button>
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
              >
                <template #append>
                  <el-popover
                    v-if="aiEnabled"
                    v-model:visible="titlePopoverVisible"
                    placement="bottom"
                    width="320"
                    trigger="click"
                  >
                    <div class="ai-title-pop">
                      <div v-if="aiTitleLoading" class="ai-title-tip">
                        正在生成候选标题…
                      </div>
                      <template v-else>
                        <div
                          v-for="t in aiTitleCandidates"
                          :key="t"
                          class="ai-title-item"
                          @click="applyTitle(t)"
                        >
                          {{ t }}
                        </div>
                        <div v-if="!aiTitleCandidates.length" class="ai-title-tip">
                          点击右侧 ✨ 按钮生成
                        </div>
                      </template>
                    </div>
                    <template #reference>
                      <el-button
                        :loading="aiTitleLoading"
                        title="AI 生成推荐标题"
                        @click="generateTitle"
                      >
                        ✨
                      </el-button>
                    </template>
                  </el-popover>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="文章正文" prop="postContent">
              <!-- 空白页灵感启动台 -->
              <div
                v-if="aiEnabled && !ruleForm.postContent.trim() && showOutlineInspiration"
                class="ai-outline-banner"
              >
                <div class="banner-left">
                  <span class="banner-badge">✨ 灵感启动台</span>
                  <span class="banner-text">
                    文章正文尚为空白？输入一个主题，让 AI 协助梳理严谨的多级技术大纲与写作要点
                  </span>
                </div>
                <div class="banner-right">
                  <el-button
                    size="small"
                    type="primary"
                    @click="openOutlineDialog"
                  >
                    生成文章大纲
                  </el-button>
                  <el-button
                    size="small"
                    text
                    title="收起引导"
                    @click="showOutlineInspiration = false"
                  >
                    <el-icon><Close /></el-icon>
                  </el-button>
                </div>
              </div>

              <!-- 续写光标流式插入状态条 -->
              <div
                v-if="aiEnabled && (inlineStreaming || inlineDoneTip)"
                class="ai-inline-bar"
                :class="{ 'is-streaming': inlineStreaming }"
              >
                <template v-if="inlineStreaming">
                  <span class="ai-inline-dot"></span>
                  <span class="ai-inline-label">AI 续写中…</span>
                  <span
                    v-if="inlineReasoning"
                    class="ai-inline-reasoning"
                    :title="inlineReasoning"
                  >
                    {{ inlineReasoning }}
                  </span>
                  <span class="ai-inline-actions">
                    <el-button size="small" type="danger" plain @click="stopInline">
                      停止
                    </el-button>
                  </span>
                </template>
                <template v-else>
                  <span class="ai-inline-label">{{ inlineDoneTip }}</span>
                  <span class="ai-inline-actions">
                    <el-button
                      v-if="inlineUndoable"
                      size="small"
                      @click="undoInline"
                    >
                      撤销
                    </el-button>
                    <el-button
                      v-if="inlineUndoable"
                      size="small"
                      type="primary"
                      plain
                      @click="retryInline"
                    >
                      重试
                    </el-button>
                  </span>
                </template>
              </div>

              <!-- 编辑器主体 -->
              <div class="editor-wrap-relative">
                <MdEditor
                  ref="editorRef"
                  v-model="ruleForm.postContent"
                  :toolbars="toolbars"
                  :floating-toolbars="aiEnabled ? floatingToolbars : []"
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
                  @input="handleEditorInput"
                >
                  <template #defToolbars>
                    <DropdownToolbar
                      v-if="aiEnabled"
                      title="AI 助手"
                      :visible="aiMenuVisible"
                      :on-change="(v: boolean) => (aiMenuVisible = v)"
                    >
                      <template #trigger>
                        <span class="ai-toolbar-trigger">✨ AI</span>
                      </template>
                      <template #overlay>
                        <ul class="ai-menu-overlay">
                          <!-- 自由指令输入框（Ask AI） -->
                          <li class="ai-menu-custom-box" @click.stop>
                            <div class="ai-custom-prompt-wrap">
                              <el-input
                                v-model="aiCustomPrompt"
                                size="small"
                                placeholder="对选中内容提要求（Enter发送）"
                                clearable
                                :disabled="aiQuotaExhausted"
                                @keydown.enter.stop="runCustomSelectionPrompt"
                              >
                                <template #suffix>
                                  <el-icon
                                    class="ai-send-icon"
                                    :class="{ active: !!aiCustomPrompt.trim() }"
                                    @click.stop="runCustomSelectionPrompt"
                                  >
                                    <Promotion />
                                  </el-icon>
                                </template>
                              </el-input>
                              <div class="ai-quick-tags">
                                <span
                                  v-for="tag in aiQuickPromptTags"
                                  :key="tag"
                                  class="ai-quick-tag"
                                  @click.stop="applyQuickPrompt(tag)"
                                >
                                  {{ tag }}
                                </span>
                              </div>
                            </div>
                          </li>
                          <li
                            v-for="m in aiMenuData"
                            :key="m.value"
                            :class="{ disabled: aiQuotaExhausted }"
                            @click="onAiMenuClick(m)"
                          >
                            {{ m.label }}
                          </li>
                          <li class="ai-menu-divider" />
                          <li class="ai-menu-section-header">
                            <span class="ai-section-title">💻 代码块助手</span>
                          </li>
                          <li
                            v-for="c in aiCodeMenuData"
                            :key="c.value"
                            :class="{ disabled: aiQuotaExhausted }"
                            @click="onAiCodeMenuClick(c)"
                          >
                            {{ c.label }}
                          </li>
                          <li class="ai-menu-divider" />
                          <li
                            v-if="aiModelOptions.length > 1"
                            class="ai-menu-models"
                            @click.stop
                          >
                            <span class="ai-menu-quota">模型</span>
                            <div class="ai-model-chips">
                              <span
                                v-for="mo in aiModelOptions"
                                :key="mo.name"
                                class="ai-model-chip"
                                :class="{
                                  active:
                                    (aiModelChoice || 'default') === mo.name
                                }"
                                :title="mo.model"
                                @click="setAIModel(mo.name)"
                              >
                                {{ mo.name === "default" ? mo.model : mo.name }}
                              </span>
                            </div>
                          </li>
                          <li
                            class="ai-menu-action"
                            @click="openTplDialog"
                          >
                            自定义指令模板
                          </li>
                          <li v-if="aiQuotaText" class="ai-menu-quota">
                            {{ aiQuotaText }}
                          </li>
                        </ul>
                      </template>
                    </DropdownToolbar>
                  </template>
                </MdEditor>

                <!-- 浮动 Tab 智能续写幽灵胶囊 -->
                <div
                  v-if="ghostTipVisible"
                  class="ghost-tip-capsule"
                  :style="{ top: `${ghostTipPos.top}px`, left: `${ghostTipPos.left}px` }"
                  @click="runContinue"
                >
                  <span class="ghost-sparkle">✨</span>
                  <span class="ghost-text">按 <kbd>Tab</kbd> 让 AI 智能续写</span>
                </div>
              </div>

              <div class="form-tip">
                支持 Markdown 语法格式 · 键入 <code>/</code> 快速唤起命令菜单 · 按 <code>⌘K</code> 快速发起 AI 创作
                <span v-if="draftSavedAtText" class="draft-indicator">
                  · 草稿自动保存于 {{ draftSavedAtText }}
                </span>
                <el-button
                  link
                  type="primary"
                  size="small"
                  class="ml-2"
                  @click="openHistoryModal"
                >
                  查看历史快照
                </el-button>
              </div>
            </el-form-item>
          </el-col>

          <!-- 右栏元数据面板 -->
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
              <div v-if="aiEnabled" class="form-tip">
                <el-popover
                  v-model:visible="tagsPopoverVisible"
                  placement="bottom"
                  width="300"
                  trigger="click"
                >
                  <div class="ai-title-pop">
                    <div v-if="aiTagsLoading" class="ai-title-tip">
                      正在分析正文生成标签建议…
                    </div>
                    <template v-else>
                      <div
                        v-for="t in aiTagCandidates"
                        :key="t"
                        class="ai-title-item"
                        @click="applyTag(t)"
                      >
                        {{ t }}
                      </div>
                      <div v-if="!aiTagCandidates.length" class="ai-title-tip">
                        点击右侧 ✨ 按钮生成
                      </div>
                    </template>
                  </div>
                  <template #reference>
                    <el-button
                      link
                      type="primary"
                      size="small"
                      :loading="aiTagsLoading"
                      @click="generateTags"
                    >
                      ✨ AI 建议标签
                    </el-button>
                  </template>
                </el-popover>
              </div>
            </el-form-item>

            <el-form-item label="文章摘要" prop="summary">
              <el-input
                v-model="ruleForm.summary"
                type="textarea"
                :rows="4"
                placeholder="请输入文章摘要，或使用 AI 发文助手一键提炼"
                maxlength="500"
                show-word-limit
              />
              <div v-if="aiEnabled" class="form-tip">
                <el-button
                  link
                  type="primary"
                  size="small"
                  :loading="aiSummaryLoading"
                  @click="generateSummary"
                >
                  ✨ AI 提炼摘要
                </el-button>
              </div>
            </el-form-item>

            <!-- 文章封面管理器模块 -->
            <el-form-item label="文章封面" class="cover-form-item">
              <CoverManager
                v-model="ruleForm.coverImage"
                @open-tech-cover="techCoverVisible = true"
              />
              <div class="form-tip">
                支持上传、外链、壁纸库或 AI 极客技术封面；留空启用动态星空流星效果
              </div>
            </el-form-item>

            <el-row :gutter="20">
              <el-col :xs="24" :span="24">
                <el-form-item label="发布时间">
                  <el-switch
                    v-model="scheduleEnabled"
                    active-text="定时发布"
                    inactive-text="立即"
                    style="margin-bottom: 8px"
                  />
                  <el-date-picker
                    v-if="scheduleEnabled"
                    v-model="scheduleTime"
                    type="datetime"
                    placeholder="选择上线时间"
                    format="YYYY-MM-DD HH:mm"
                    value-format="X"
                    :clearable="false"
                    style="width: 100%"
                  />
                  <div v-if="scheduleEnabled" class="form-tip">
                    保存后文章进入“已发布”状态，到所选时间才对访客可见
                  </div>
                  <div v-else class="form-tip">
                    保存后仍为草稿状态，发布请在文章列表操作
                  </div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :span="24">
                <el-form-item label="文章系列">
                  <el-input
                    v-model="ruleForm.series"
                    placeholder="选填，如：Go 并发编程核心原理"
                    clearable
                    maxlength="50"
                  />
                  <div class="form-tip">
                    同名系列的文章会在正文下方生成连载导航
                  </div>
                </el-form-item>
              </el-col>

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

    <!-- 隐藏的本地 Markdown 文件上传 input -->
    <input
      ref="mdFileInputRef"
      type="file"
      accept=".md,.markdown"
      style="display: none"
      @change="handleImportMdFile"
    />

    <!-- AI 结果流式抽屉组件 -->
    <AiResultDrawer
      v-model:visible="aiDrawerVisible"
      v-model:result-text="aiResultText"
      :last-request="aiLastRequest"
      :task-label="aiTaskLabel"
      :streaming="aiStreaming"
      :error="aiError"
      :reasoning-text="aiReasoningText"
      :meta="aiMeta"
      :quota-text="aiQuotaText"
      :history="aiHistory"
      :history-idx="aiHistoryIdx"
      :original-selection="aiOriginalSelection"
      :apply-type="aiApplyType"
      :proofread-items="proofreadItems"
      :pending-proofread-count="pendingProofreadCount"
      :applied-proofread-count="appliedProofreadCount"
      @close="aiDrawerVisible = false"
      @apply-result="applyAIResult"
      @insert-below="insertBelowAIResult"
      @copy-result="copyAIResult"
      @regen="regenAI"
      @stop="stopAI"
      @refine="refineAI"
      @nav-history="aiNavHistory"
      @select-history="selectHistoryIndex"
      @apply-proofread="applyProofreadItem"
      @ignore-proofread="ignoreProofreadItem"
      @locate-proofread="locateProofreadItem"
      @apply-all-proofread="applyAllProofreadItems"
    />

    <!-- 就地校对悬浮卡片 -->
    <ProofreadHoverCard
      :visible="hoverCardVisible"
      :item="hoverCardItem"
      :position="hoverCardPos"
      @apply="applyProofreadItem"
      @ignore="ignoreProofreadItem"
      @close="hideHoverCard"
    />

    <!-- 快速命令 Slash 菜单 -->
    <SlashCommand
      :visible="slashCommandVisible"
      :position="slashCommandPos"
      :filter-text="slashFilterText"
      @select="handleSelectSlashCommand"
      @close="slashCommandVisible = false"
    />

    <!-- Cmd+K 浮动 AI 指令中心 -->
    <AiFloatingBar
      :visible="floatingAiBarVisible"
      :has-selection="!!floatingAiSelection"
      :selection-length="floatingAiSelection.length"
      :models="aiModelOptions"
      :current-model="aiModelChoice"
      @update:current-model="setAIModel"
      @submit="handleFloatingAiSubmit"
      @close="floatingAiBarVisible = false"
    />

    <!-- 草稿时光机快照对比弹窗 -->
    <DraftHistoryModal
      v-model:visible="historyModalVisible"
      :revisions="revisionHistory"
      :current-content="ruleForm.postContent"
      @restore="restoreRevision"
      @delete="deleteRevision"
      @clear-all="clearAllRevisions"
    />

    <!-- AI 指令模板管理弹窗 -->
    <AiTemplateModal v-model:visible="aiTplVisible" />

    <!-- AI 多级大纲生成器弹窗 -->
    <AiOutlineModal
      v-model:visible="outlineDialogVisible"
      :initial-topic="ruleForm.title"
      :ai-model-choice="aiModelChoice"
      :ai-quota-exhausted="aiQuotaExhausted"
      @apply="handleApplyOutline"
    />

    <!-- AI 极客技术封面工作台 -->
    <TechCoverModal
      v-model:visible="techCoverVisible"
      :post-title="ruleForm.title"
      :post-tags="ruleForm.tags"
      :post-category-name="currentCategoryName"
      @applied="onTechCoverApplied"
    />

    <!-- AI 发文助手与 SEO 质量体检一站式弹窗 -->
    <AiPublishModal
      v-model:visible="publishModalVisible"
      :current-title="ruleForm.title"
      :current-summary="ruleForm.summary"
      :current-tags="ruleForm.tags"
      :current-category-id="ruleForm.categoryId"
      :current-cover="ruleForm.coverImage"
      :post-content="ruleForm.postContent"
      :system-tags="tags"
      :submit-loading="saveLoading"
      @apply="handleApplyPublishMeta"
    />

    <!-- 发文成功喜报弹窗 -->
    <PostPublishModal
      :visible="postPublishCelebrationVisible"
      :is-edit="isEdit"
      :post-id="publishedPostId"
      :post-title="publishedPostTitle"
      @continue-edit="handleContinueEdit"
      @write-another="handleWriteAnother"
      @back-to-list="router.push({ name: '内容管理' })"
    />

    <!-- 语言转换目标选择弹窗 -->
    <el-dialog
      v-model="codeConvertVisible"
      title="🔄 跨语言代码转换"
      width="420px"
      append-to-body
    >
      <div style="padding: 10px 0;">
        <div style="margin-bottom: 8px; font-size: 13px; color: #606266;">
          选择目标编程语言：
        </div>
        <el-radio-group v-model="targetConvertLang" style="display: flex; flex-wrap: wrap; gap: 8px;">
          <el-radio-button
            v-for="lang in ['Go', 'TypeScript', 'Python', 'Rust', 'Java', 'C++', 'PHP']"
            :key="lang"
            :label="lang"
          />
        </el-radio-group>
      </div>
      <template #footer>
        <el-button @click="codeConvertVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCodeConvert">开始转换</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
defineOptions({
  name: "Editor"
});

import {
  ref,
  reactive,
  computed,
  watch,
  onMounted,
  onBeforeUnmount,
  nextTick
} from "vue";
import { useRouter, useRoute } from "vue-router";
import type { FormInstance, FormRules } from "element-plus";
import {
  Aim,
  Close,
  DocumentCopy,
  Download,
  MagicStick,
  Promotion,
  Tickets,
  Upload
} from "@element-plus/icons-vue";
import {
  MdEditor,
  DropdownToolbar,
  config,
  type ToolbarNames
} from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import "md-editor-v3/lib/preview.css";
import { EditorView } from "@codemirror/view";
import { useDebounceFn } from "@vueuse/core";

import { getCategoryList } from "@/api/category";
import { getTagList } from "@/api/tag";
import { getPost, savePost, updatePost } from "@/api/post";
import { upload } from "@/api/common";
import { message } from "@/utils/message";
import { streamAIEdit } from "@/api/ai";
import { parseProofreadOutput, type ProofreadItem } from "@/utils/proofread";
import { detectCodeBlock, guessLanguage, wrapCodeFence } from "@/utils/codeBlock";

// Hooks
import { useDraftHistory } from "./hooks/useDraftHistory";
import { useAiStreaming } from "./hooks/useAiStreaming";
import { useProofreadHighlight, proofreadStateField } from "./hooks/useProofreadHighlight";

// Components
import CoverManager from "./components/CoverManager.vue";
import DraftHistoryModal from "./components/DraftHistoryModal.vue";
import AiTemplateModal from "./components/AiTemplateModal.vue";
import AiOutlineModal from "./components/AiOutlineModal.vue";
import TechCoverModal from "./components/TechCoverModal.vue";
import AiPublishModal from "./components/AiPublishModal.vue";
import PostPublishModal from "./components/PostPublishModal.vue";
import ProofreadHoverCard from "./components/ProofreadHoverCard.vue";
import SlashCommand, { type SlashCommandItem } from "./components/SlashCommand.vue";
import AiFloatingBar from "./components/AiFloatingBar.vue";
import AiResultDrawer from "./components/AiResultDrawer.vue";

const router = useRouter();
const route = useRoute();
const ruleFormRef = ref<FormInstance>();
const saveLoading = ref(false);
const editorRef = ref<any>(null);

const isEdit = computed(() => !!route.params.id);
const currentPostId = computed(() => (route.params.id as string) || undefined);

// 专注模式
const focusMode = ref(false);
const mainColProps = computed(() =>
  focusMode.value
    ? { xs: 24, sm: 24, md: 24, lg: 24 }
    : { xs: 24, sm: 24, md: 16, lg: 18 }
);

const editorStyle = computed(() => {
  const offset = focusMode.value ? 260 : 320;
  return {
    height: `calc(100vh - ${offset}px)`,
    minHeight: focusMode.value ? "600px" : "500px"
  };
});

const tags = ref<any[]>([]);
const categories = ref<any[]>([]);

const ruleForm = ref({
  postId: undefined as number | undefined,
  postSlug: "",
  title: "",
  categoryId: undefined,
  tags: [] as string[],
  summary: "",
  status: 1,
  type: 1,
  postContent: "",
  postContentHtml: "",
  coverImage: "",
  series: ""
});

const currentCategoryName = computed(() => {
  const cat = categories.value.find(c => c.categoryId === ruleForm.value.categoryId);
  return cat?.categoryName || "";
});

// 定时发布
const scheduleEnabled = ref(false);
const scheduleTime = ref<string | number>("");

// 弹窗显隐控制
const aiTplVisible = ref(false);
const outlineDialogVisible = ref(false);
const techCoverVisible = ref(false);
const publishModalVisible = ref(false);
const postPublishCelebrationVisible = ref(false);
const publishedPostId = ref<string | number>("");
const publishedPostTitle = ref("");
const codeConvertVisible = ref(false);
const targetConvertLang = ref("TypeScript");
const showOutlineInspiration = ref(true);

// 导入导出 Markdown
const mdFileInputRef = ref<HTMLInputElement | null>(null);

// ==================== Hooks 集成 ====================
// 1. 草稿时光机
const {
  draftSavedAtText,
  historyModalVisible,
  revisionHistory,
  persistDraftDebounced,
  handleSaveDraft,
  tryRestoreDraft,
  clearDraft,
  openHistoryModal,
  restoreRevision,
  deleteRevision,
  clearAllRevisions
} = useDraftHistory(currentPostId, ruleForm, isEdit);

// 2. AI 流式服务
const {
  aiEnabled,
  aiQuotaText,
  aiQuotaExhausted,
  aiModelOptions,
  aiModelChoice,
  setAIModel,
  refreshAIStatus,
  aiDrawerVisible,
  aiStreaming,
  aiResultText,
  aiReasoningText,
  aiError,
  aiMeta,
  aiTaskLabel,
  aiApplyType,
  aiOriginalSelection,
  aiSelRange,
  aiLastRequest,
  aiHistory,
  aiHistoryIdx,
  aiNavHistory,
  selectHistoryIndex,
  startAiStream,
  stopAI,
  regenAI,
  refineAI
} = useAiStreaming();

// 3. 全文校对与波浪线高亮
const proofreadItems = ref<ProofreadItem[]>([]);
const pendingProofreadCount = computed(
  () => proofreadItems.value.filter(i => !i.status || i.status === "pending").length
);
const appliedProofreadCount = computed(
  () => proofreadItems.value.filter(i => i.status === "applied").length
);

const {
  hoverCardVisible,
  hoverCardItem,
  hoverCardPos,
  refreshDecorations,
  clearDecorations,
  locateProofreadItem,
  applyProofreadItem,
  ignoreProofreadItem,
  applyAllProofreadItems,
  showHoverCard,
  hideHoverCard
} = useProofreadHighlight(editorRef, proofreadItems, newContent => {
  ruleForm.value.postContent = newContent;
});

// ==================== CodeMirror 扩展与键盘交互 ====================
// 浮动 AI 指令中心 (Cmd+K)
const floatingAiBarVisible = ref(false);
const floatingAiSelection = ref("");

// Slash 命令菜单 (/)
const slashCommandVisible = ref(false);
const slashCommandPos = ref<{
  top: number;
  left: number;
  placement?: "top" | "bottom";
}>({ top: 0, left: 0, placement: "bottom" });
const slashFilterText = ref("");

// Ghost 智能续写胶囊 (Tab)
const ghostTipVisible = ref(false);
const ghostTipPos = ref({ top: 0, left: 0 });

const openFloatingAiBar = (view: EditorView) => {
  if (!aiEnabled.value) return;
  const sel = view.state.selection.main;
  floatingAiSelection.value = sel.empty ? "" : view.state.sliceDoc(sel.from, sel.to);
  floatingAiBarVisible.value = true;
};

const handleEditorMouseOver = (event: MouseEvent) => {
  const target = event.target as HTMLElement | null;
  if (!target) return;
  const typoEl = target.closest(
    ".cm-proofread-typo, .cm-proofread-grammar, .cm-proofread-style"
  ) as HTMLElement | null;

  if (typoEl) {
    const id = typoEl.getAttribute("data-proofread-id");
    const item = proofreadItems.value.find(i => i.id === id);
    if (item && item.status !== "applied" && item.status !== "ignored") {
      const rect = typoEl.getBoundingClientRect();
      showHoverCard(item, {
        top: rect.bottom + 6,
        left: Math.max(10, rect.left - 20)
      });
    }
  }
};

// 注册 CodeMirror 扩展
config({
  codeMirrorExtensions(extensions) {
    const keyHandler = EditorView.domEventHandlers({
      keydown(event, view) {
        if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
          event.preventDefault();
          openFloatingAiBar(view);
          return true;
        }
        if (event.key === "Tab" && ghostTipVisible.value) {
          event.preventDefault();
          runContinue();
          ghostTipVisible.value = false;
          return true;
        }
        return false;
      },
      mouseover(event) {
        handleEditorMouseOver(event);
      },
      scroll(event, view) {
        if (slashCommandVisible.value) {
          checkSlashCommand();
        }
      }
    });

    return [
      ...extensions,
      {
        type: "proofread",
        extension: proofreadStateField
      },
      {
        type: "keyHandler",
        extension: keyHandler
      }
    ];
  }
});

// 编辑器输入监听 (Slash 命令与 Ghost 提示)
const handleEditorInput = (_event?: Event) => {
  checkSlashCommand();
  checkGhostTip();
};

const checkSlashCommand = () => {
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  const pos = view.state.selection.main.head;
  const line = view.state.doc.lineAt(pos);
  const textBeforeCursor = line.text.slice(0, pos - line.from);

  if (textBeforeCursor.startsWith("/") && textBeforeCursor.length <= 12) {
    const coords = view.coordsAtPos(pos);
    if (coords && coords.bottom >= 0 && coords.top <= window.innerHeight) {
      const menuHeight = 320;
      const menuWidth = 320;
      const vh = window.innerHeight;
      const vw = window.innerWidth;
      const spaceBelow = vh - coords.bottom;
      const spaceAbove = coords.top;

      // 当底部可用空间不足预估高度(320px)且上方空间比下方更宽敞时，自动向上翻转弹出
      const placeUp = spaceBelow < menuHeight && spaceAbove > spaceBelow;

      slashCommandPos.value = {
        top: placeUp ? coords.top - 6 : coords.bottom + 6,
        left: Math.min(Math.max(12, coords.left - 10), Math.max(12, vw - menuWidth - 20)),
        placement: placeUp ? "top" : "bottom"
      };
      slashFilterText.value = textBeforeCursor.slice(1);
      slashCommandVisible.value = true;
      return;
    }
  }
  slashCommandVisible.value = false;
};

const checkGhostTip = useDebounceFn(() => {
  const view = editorRef.value?.getEditorView?.();
  if (!view || !aiEnabled.value || aiStreaming.value) {
    ghostTipVisible.value = false;
    return;
  }
  const pos = view.state.selection.main.head;
  const doc = view.state.doc;
  if (pos < 10) {
    ghostTipVisible.value = false;
    return;
  }
  const line = doc.lineAt(pos);
  if (pos === line.to && line.text.trim().length > 6) {
    const coords = view.coordsAtPos(pos);
    if (coords) {
      ghostTipPos.value = {
        top: coords.top - 28,
        left: Math.min(window.innerWidth - 200, coords.left + 10)
      };
      ghostTipVisible.value = true;
      return;
    }
  }
  ghostTipVisible.value = false;
}, 1200);

const handleSelectSlashCommand = (cmd: SlashCommandItem) => {
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  const pos = view.state.selection.main.head;
  const line = view.state.doc.lineAt(pos);

  view.dispatch({
    changes: { from: line.from, to: pos, insert: "" }
  });

  slashCommandVisible.value = false;

  if (cmd.template) {
    view.dispatch({
      changes: { from: line.from, insert: cmd.template },
      selection: { anchor: line.from + cmd.template.length },
      scrollIntoView: true
    });
    ruleForm.value.postContent = view.state.doc.toString();
  } else if (cmd.action === "ai-continue") {
    runContinue();
  } else if (cmd.action === "ai-outline") {
    openOutlineDialog();
  }
};

const handleFloatingAiSubmit = (payload: { prompt: string; task?: string; model?: string }) => {
  const sel = captureSelection();
  const text = sel?.text || ruleForm.value.postContent;
  startAiStream(
    {
      task: (payload.task as any) || "polish",
      selection: sel?.text,
      digest: ruleForm.value.postContent,
      instruction: payload.prompt,
      model: payload.model
    },
    {
      label: payload.prompt.slice(0, 10),
      applyType: sel?.text ? "replace-selection" : "insert-cursor",
      selectionText: sel?.text,
      selectionRange: sel ? { from: sel.from, to: sel.to } : undefined
    }
  );
};

// ==================== 选区捕获与常用 AI 动作 ====================
function captureSelection() {
  const view = editorRef.value?.getEditorView?.();
  const sel = view?.state?.selection?.main;
  if (!view || !sel) return null;
  return {
    view,
    from: sel.from,
    to: sel.to,
    text: view.state.sliceDoc(sel.from, sel.to)
  };
}

// 续写状态
const inlineStreaming = ref(false);
const inlineDoneTip = ref("");
const inlineReasoning = ref("");
const inlineUndoable = ref(false);
let inlineAbort: (() => void) | null = null;
let inlineInsertedRange: { from: number; to: number } | null = null;

function runContinue() {
  if (aiQuotaExhausted.value) {
    message("今日 AI 调用配额已用完", { type: "warning" });
    return;
  }
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;

  const sel = view.state.selection.main;
  const insertPos = sel.head;
  const beforeText = view.state.sliceDoc(0, insertPos);
  const afterText = view.state.sliceDoc(insertPos);

  inlineStreaming.value = true;
  inlineDoneTip.value = "";
  inlineReasoning.value = "";
  inlineUndoable.value = false;
  inlineInsertedRange = { from: insertPos, to: insertPos };

  let currentEnd = insertPos;

  inlineAbort = streamAIEdit(
    {
      task: "continue",
      before: beforeText,
      after: afterText,
      title: ruleForm.value.title,
      series: ruleForm.value.series,
      model: aiModelChoice.value === "default" ? undefined : aiModelChoice.value
    },
    {
      onDelta: delta => {
        view.dispatch({
          changes: { from: currentEnd, insert: delta },
          selection: { anchor: currentEnd + delta.length },
          scrollIntoView: true
        });
        currentEnd += delta.length;
        if (inlineInsertedRange) inlineInsertedRange.to = currentEnd;
      },
      onReasoning: r => {
        inlineReasoning.value += r;
      },
      onDone: () => {
        inlineStreaming.value = false;
        inlineDoneTip.value = "AI 续写完成";
        inlineUndoable.value = true;
        ruleForm.value.postContent = view.state.doc.toString();
        refreshAIStatus();
      },
      onError: err => {
        inlineStreaming.value = false;
        inlineDoneTip.value = "续写失败";
        message(err || "续写中断", { type: "error" });
      }
    }
  );
}

function stopInline() {
  inlineAbort?.();
  inlineAbort = null;
  inlineStreaming.value = false;
}

function undoInline() {
  const view = editorRef.value?.getEditorView?.();
  if (!view || !inlineInsertedRange) return;
  view.dispatch({
    changes: {
      from: inlineInsertedRange.from,
      to: inlineInsertedRange.to,
      insert: ""
    },
    selection: { anchor: inlineInsertedRange.from }
  });
  ruleForm.value.postContent = view.state.doc.toString();
  inlineUndoable.value = false;
  inlineDoneTip.value = "已撤销续写内容";
}

function retryInline() {
  undoInline();
  runContinue();
}

// 应用抽屉结果回写选区
function applyAIResult() {
  const view = editorRef.value?.getEditorView?.();
  if (!view || !aiResultText.value) return;

  const sel = captureSelection();
  const targetRange = aiSelRange.value || (sel ? { from: sel.from, to: sel.to } : null);
  const from = targetRange?.from ?? 0;
  const to = targetRange?.to ?? 0;

  view.dispatch({
    changes: { from, to, insert: aiResultText.value },
    selection: { anchor: from + aiResultText.value.length },
    scrollIntoView: true
  });

  ruleForm.value.postContent = view.state.doc.toString();
  aiDrawerVisible.value = false;
  message("已应用 AI 内容至选区", { type: "success" });
}

function insertBelowAIResult() {
  const view = editorRef.value?.getEditorView?.();
  if (!view || !aiResultText.value) return;

  const sel = captureSelection();
  const targetRange = aiSelRange.value || (sel ? { from: sel.from, to: sel.to } : null);
  const pos = targetRange?.to ?? view.state.doc.length;
  const insertText = "\n\n" + aiResultText.value;

  view.dispatch({
    changes: { from: pos, insert: insertText },
    selection: { anchor: pos + insertText.length },
    scrollIntoView: true
  });

  ruleForm.value.postContent = view.state.doc.toString();
  aiDrawerVisible.value = false;
  message("已在下方插入生成内容", { type: "success" });
}

function copyAIResult() {
  if (!aiResultText.value) return;
  navigator.clipboard.writeText(aiResultText.value);
  message("已复制到剪贴板", { type: "success" });
}

// 标题生成
const titlePopoverVisible = ref(false);
const aiTitleLoading = ref(false);
const aiTitleCandidates = ref<string[]>([]);

function generateTitle() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先输入正文内容，以便 AI 准确提炼标题", { type: "warning" });
    return;
  }
  aiTitleLoading.value = true;
  aiTitleCandidates.value = [];
  let acc = "";

  streamAIEdit(
    {
      task: "title",
      digest: ruleForm.value.postContent,
      model: aiModelChoice.value === "default" ? undefined : aiModelChoice.value
    },
    {
      onDelta: d => (acc += d),
      onDone: () => {
        aiTitleLoading.value = false;
        aiTitleCandidates.value = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 5);
        refreshAIStatus();
      },
      onError: () => {
        aiTitleLoading.value = false;
      }
    }
  );
}

function applyTitle(t: string) {
  ruleForm.value.title = t;
  titlePopoverVisible.value = false;
  message("已应用标题", { type: "success" });
}

// 标签生成
const tagsPopoverVisible = ref(false);
const aiTagsLoading = ref(false);
const aiTagCandidates = ref<string[]>([]);

function generateTags() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先输入正文内容", { type: "warning" });
    return;
  }
  aiTagsLoading.value = true;
  aiTagCandidates.value = [];
  let acc = "";

  streamAIEdit(
    {
      task: "tags",
      digest: ruleForm.value.postContent,
      title: ruleForm.value.title,
      model: aiModelChoice.value === "default" ? undefined : aiModelChoice.value
    },
    {
      onDelta: d => (acc += d),
      onDone: () => {
        aiTagsLoading.value = false;
        aiTagCandidates.value = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 8);
        refreshAIStatus();
      },
      onError: () => {
        aiTagsLoading.value = false;
      }
    }
  );
}

function applyTag(t: string) {
  if (!ruleForm.value.tags.includes(t)) {
    ruleForm.value.tags.push(t);
  }
}

// 摘要提炼
const aiSummaryLoading = ref(false);

function generateSummary() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先输入正文内容", { type: "warning" });
    return;
  }
  aiSummaryLoading.value = true;
  let acc = "";

  streamAIEdit(
    {
      task: "summary",
      digest: ruleForm.value.postContent,
      model: aiModelChoice.value === "default" ? undefined : aiModelChoice.value
    },
    {
      onDelta: d => (acc += d),
      onDone: () => {
        aiSummaryLoading.value = false;
        if (acc.trim()) ruleForm.value.summary = acc.trim();
        refreshAIStatus();
      },
      onError: () => {
        aiSummaryLoading.value = false;
      }
    }
  );
}

// Toolbar 菜单项
const aiMenuVisible = ref(false);
const aiCustomPrompt = ref("");
const aiQuickPromptTags = [
  "转为表格",
  "提炼核心要点",
  "生成 Mermaid 流程图",
  "为代码添加注释",
  "排查代码 Bug",
  "更口语化"
];

const aiMenuData = [
  { label: "✨ 智能润色", value: "polish" },
  { label: "✍️ 深度续写", value: "continue" },
  { label: "🔍 全文校对 (错别字/语病)", value: "proofread" },
  { label: "💡 提炼摘要", value: "summary" },
  { label: "📑 生成文章大纲", value: "outline" },
  { label: "⚡ 精简内容", value: "shorten" },
  { label: "📖 丰富扩写", value: "expand" },
  { label: "🎓 学术专业化", value: "academic" },
  { label: "☕ 口语平易化", value: "casual" },
  { label: "🌐 翻译为地道英文", value: "translate" }
];

const aiCodeMenuData = [
  { label: "📝 为代码添加精细注释", value: "code_comment" },
  { label: "🐞 排查潜在 Bug 与并发安全", value: "code_bug" },
  { label: "⚡ 性能与可读性优化", value: "code_optimize" },
  { label: "🧪 自动生成单元测试", value: "code_test" },
  { label: "🔄 跨语言代码转换", value: "code_convert" },
  { label: "💡 深入解释核心逻辑", value: "code_explain" },
  { label: "🗺️ 生成 Mermaid 架构流程图", value: "mermaid" }
];

function applyQuickPrompt(tag: string) {
  aiCustomPrompt.value = tag;
  runCustomSelectionPrompt();
}

function runCustomSelectionPrompt() {
  const prompt = aiCustomPrompt.value.trim();
  if (!prompt) return;
  aiMenuVisible.value = false;
  aiCustomPrompt.value = "";

  const sel = captureSelection();
  const text = sel?.text || ruleForm.value.postContent;

  startAiStream(
    {
      task: "polish",
      selection: sel?.text,
      digest: ruleForm.value.postContent,
      instruction: prompt
    },
    {
      label: prompt.slice(0, 10),
      applyType: sel?.text ? "replace-selection" : "insert-cursor",
      selectionText: sel?.text,
      selectionRange: sel ? { from: sel.from, to: sel.to } : undefined
    }
  );
}

function onAiMenuClick(item: any) {
  aiMenuVisible.value = false;
  const value = item.value;

  if (value === "continue") {
    runContinue();
    return;
  }
  if (value === "outline") {
    openOutlineDialog();
    return;
  }
  if (value === "proofread") {
    startAiStream(
      {
        task: "proofread",
        digest: ruleForm.value.postContent
      },
      {
        label: "全文校对",
        onDone: fullText => {
          proofreadItems.value = parseProofreadOutput(fullText);
          refreshDecorations();
          if (proofreadItems.value.length > 0) {
            message(
              `全文校对完成，发现 ${proofreadItems.value.length} 处待优化建议，已在正文中以波浪线高亮标出`,
              { type: "info" }
            );
          } else {
            message("全文校对完成，未发现明显语病或错别字！", { type: "success" });
          }
        }
      }
    );
    return;
  }

  const sel = captureSelection();
  const text = sel?.text || ruleForm.value.postContent;

  startAiStream(
    {
      task: "polish",
      selection: sel?.text,
      digest: ruleForm.value.postContent,
      mode: value
    },
    {
      label: item.label,
      applyType: sel?.text ? "replace-selection" : "insert-cursor",
      selectionText: sel?.text,
      selectionRange: sel ? { from: sel.from, to: sel.to } : undefined
    }
  );
}

let pendingCodeBlockInfo: ReturnType<typeof detectCodeBlock> = null;
let pendingCodeSel: ReturnType<typeof captureSelection> = null;

function onAiCodeMenuClick(item: any) {
  aiMenuVisible.value = false;
  const sel = captureSelection();
  const view = editorRef.value?.getEditorView?.();
  const doc = view?.state?.doc?.toString() || ruleForm.value.postContent;
  const from = sel ? sel.from : 0;
  const to = sel ? sel.to : doc.length;
  const codeInfo = detectCodeBlock(doc, from, to);

  if (item.value === "code_convert") {
    pendingCodeBlockInfo = codeInfo;
    pendingCodeSel = sel;
    codeConvertVisible.value = true;
    return;
  }

  const isBlock = !!codeInfo;
  const code = codeInfo?.code || sel?.text || ruleForm.value.postContent;

  startAiStream(
    {
      task: "polish",
      mode: item.value,
      selection: code,
      digest: ruleForm.value.postContent
    },
    {
      label: item.label,
      applyType: "replace-selection",
      selectionText: isBlock ? codeInfo.fullBlock : sel?.text,
      selectionRange: isBlock
        ? { from: codeInfo.from, to: codeInfo.to }
        : sel
        ? { from: sel.from, to: sel.to }
        : undefined
    }
  );
}

function confirmCodeConvert() {
  codeConvertVisible.value = false;
  const sel = pendingCodeSel || captureSelection();
  const view = editorRef.value?.getEditorView?.();
  const doc = view?.state?.doc?.toString() || ruleForm.value.postContent;
  const from = sel ? sel.from : 0;
  const to = sel ? sel.to : doc.length;
  const codeInfo = pendingCodeBlockInfo || detectCodeBlock(doc, from, to);
  pendingCodeBlockInfo = null;
  pendingCodeSel = null;

  const isBlock = !!codeInfo;
  const code = codeInfo?.code || sel?.text || ruleForm.value.postContent;

  startAiStream(
    {
      task: "polish",
      mode: "code_convert",
      selection: code,
      digest: ruleForm.value.postContent,
      instruction: `将上述代码转换为 ${targetConvertLang.value} 语言实现`
    },
    {
      label: `转为 ${targetConvertLang.value}`,
      applyType: "replace-selection",
      selectionText: isBlock ? codeInfo.fullBlock : sel?.text,
      selectionRange: isBlock
        ? { from: codeInfo.from, to: codeInfo.to }
        : sel
        ? { from: sel.from, to: sel.to }
        : undefined
    }
  );
}

function openTplDialog() {
  aiMenuVisible.value = false;
  aiTplVisible.value = true;
}

function openOutlineDialog() {
  outlineDialogVisible.value = true;
}

function handleApplyOutline(payload: {
  content: string;
  mode: "replace" | "append" | "sync-title";
  title?: string;
  summary?: string;
}) {
  if (payload.mode === "replace") {
    ruleForm.value.postContent = payload.content;
  } else if (payload.mode === "append") {
    ruleForm.value.postContent =
      ruleForm.value.postContent.trimEnd() + "\n\n" + payload.content;
  } else if (payload.mode === "sync-title") {
    ruleForm.value.postContent = payload.content;
    if (payload.title) ruleForm.value.title = payload.title;
    if (payload.summary) ruleForm.value.summary = payload.summary;
  }
  showOutlineInspiration.value = false;
  message("大纲已成功应用到编辑器", { type: "success" });
}

function onTechCoverApplied(url: string) {
  ruleForm.value.coverImage = url;
}

function openPublishModal() {
  publishModalVisible.value = true;
}

function handleApplyPublishMeta(payload: {
  title: string;
  summary: string;
  tags: string[];
  publishImmediately: boolean;
}) {
  if (payload.title) ruleForm.value.title = payload.title;
  if (payload.summary) ruleForm.value.summary = payload.summary;
  if (payload.tags && payload.tags.length > 0) ruleForm.value.tags = payload.tags;

  if (payload.publishImmediately) {
    submitForm(ruleFormRef.value);
  } else {
    message("已更新文章标题、摘要与标签", { type: "success" });
  }
}

function handleContinueEdit() {
  postPublishCelebrationVisible.value = false;
  if (!isEdit.value && publishedPostId.value) {
    router.replace({
      name: "内容编辑",
      params: { id: String(publishedPostId.value) }
    });
  }
}

function handleWriteAnother() {
  postPublishCelebrationVisible.value = false;
  ruleForm.value = {
    postId: undefined,
    postSlug: "",
    title: "",
    categoryId: undefined,
    tags: [],
    summary: "",
    status: 1,
    type: 1,
    postContent: "",
    postContentHtml: "",
    coverImage: "",
    series: ""
  };
  clearDecorations();
  router.push({ path: "/manage/editor" });
}

// 专注模式切换
const toggleFocusMode = () => {
  focusMode.value = !focusMode.value;
  message(focusMode.value ? "已进入专注写作模式 (按 Esc 退出)" : "已退出专注模式", {
    type: "info"
  });
};

const handleEscKey = () => {
  if (focusMode.value) {
    focusMode.value = false;
  }
};

// 导入导出 Markdown
const triggerImportMd = () => {
  mdFileInputRef.value?.click();
};

const handleImportMdFile = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = e => {
    const raw = (e.target?.result as string) || "";
    let body = raw;

    // 解析 FrontMatter
    const fmMatch = raw.match(/^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/);
    if (fmMatch) {
      body = fmMatch[2].trimStart();
      const yaml = fmMatch[1];
      const titleMatch = yaml.match(/(?:^|\n)title:\s*["']?([^\n"']+)["']?/);
      if (titleMatch && !ruleForm.value.title) {
        ruleForm.value.title = titleMatch[1].trim();
      }
      const summaryMatch = yaml.match(/(?:^|\n)description:\s*["']?([^\n"']+)["']?/);
      if (summaryMatch && !ruleForm.value.summary) {
        ruleForm.value.summary = summaryMatch[1].trim();
      }
    }

    ruleForm.value.postContent = body;
    message("Markdown 内容已成功导入", { type: "success" });
  };
  reader.readAsText(file, "UTF-8");
  target.value = "";
};

const handleExportMd = () => {
  const f = ruleForm.value;
  const tagListStr = (f.tags || []).map((t: string) => `  - "${t}"`).join("\n");
  const frontMatter = [
    "---",
    `title: "${(f.title || "未命名博文").replace(/"/g, '\\"')}"`,
    f.summary ? `description: "${f.summary.replace(/"/g, '\\"')}"` : "",
    f.tags && f.tags.length ? `tags:\n${tagListStr}` : "",
    `date: "${new Date().toISOString()}"`,
    "---",
    "",
    ""
  ]
    .filter(Boolean)
    .join("\n");

  const blob = new Blob([frontMatter + (f.postContent || "")], {
    type: "text/markdown;charset=utf-8"
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `${(f.title || "post").slice(0, 30)}.md`;
  a.click();
  URL.revokeObjectURL(url);
  message("Markdown 文件已开始下载", { type: "success" });
};

// 表单校验规则
const rules = reactive<FormRules>({
  title: [
    { required: true, message: "请输入文章标题", trigger: "blur" },
    { min: 2, max: 200, message: "长度在 2 到 200 个字符", trigger: "blur" }
  ],
  categoryId: [
    { required: true, message: "请选择文章分类", trigger: "change" }
  ],
  postContent: [
    { required: true, message: "请输入文章内容", trigger: "blur" }
  ]
});

// 编辑器工具栏配置
const toolbars: ToolbarNames[] = [
  "bold",
  "underline",
  "italic",
  "strikeThrough",
  "-",
  "title",
  "sub",
  "sup",
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

const floatingToolbars: ToolbarNames[] = [
  "bold",
  "italic",
  "underline",
  "strikeThrough",
  "-",
  "title",
  "quote",
  "code",
  "link"
];

const onHtmlChanged = (html: string) => {
  ruleForm.value.postContentHtml = html;
};

const onGetCatalog = () => {};

const onUploadImg = async (
  files: File[],
  callback: (urls: string[]) => void
) => {
  try {
    const res = await Promise.all(
      files.map(file => {
        return new Promise<string>((resolve, reject) => {
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
    callback(res);
  } catch (e: any) {
    message(e?.message || "上传图片失败", { type: "error" });
  }
};

// 提交表单
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;

  const handleInvalid = () => {
    if (focusMode.value) focusMode.value = false;
    message("表单验证失败，请检查必填项", { type: "warning" });
  };

  await formEl.validate(async (valid) => {
    if (valid) {
      saveLoading.value = true;
      try {
        const payload: Record<string, any> = { ...ruleForm.value };
        delete payload.pubTime;
        delete payload.published;
        if (scheduleEnabled.value && scheduleTime.value) {
          payload.pubTime = Number(scheduleTime.value);
          payload.published = 1;
        } else if (!isEdit.value) {
          // 新建文章：若状态为公开(1)，默认发布 published=1；若为草稿(0)，则 published=0
          payload.published = ruleForm.value.status === 1 ? 1 : 0;
        }

        let res;
        if (isEdit.value) {
          res = await updatePost(payload);
        } else {
          res = await savePost(payload);
        }

        if (res.code === 200) {
          await clearDraft();
          message(`${isEdit.value ? "更新" : "发布"}文章成功`, {
            type: "success"
          });
          const targetSlug =
            res.payload?.postSlug ||
            ruleForm.value.postSlug ||
            (route.params.id as string) ||
            res.payload?.postId ||
            res.payload?.id ||
            ruleForm.value.postId ||
            "";
          publishedPostId.value = targetSlug;
          publishedPostTitle.value = ruleForm.value.title;
          if (res.payload?.postId) ruleForm.value.postId = res.payload.postId;
          if (res.payload?.postSlug) ruleForm.value.postSlug = res.payload.postSlug;
          postPublishCelebrationVisible.value = true;
        } else {
          message(
            `${isEdit.value ? "更新" : "发布"}文章失败: ${res.message || "未知错误"}`,
            { type: "error" }
          );
        }
      } catch (error: any) {
        message(`操作失败: ${error?.message || error}`, { type: "error" });
      } finally {
        saveLoading.value = false;
      }
    } else {
      handleInvalid();
    }
  });
};

// 监听表单改动自动保存草稿
watch(
  ruleForm,
  () => {
    persistDraftDebounced();
  },
  { deep: true }
);

// 初始化数据
const initData = async () => {
  try {
    const [catRes, tagRes] = await Promise.all([
      getCategoryList(),
      getTagList()
    ]);
    if (catRes.code === 200) categories.value = catRes.payload ?? [];
    if (tagRes.code === 200) tags.value = tagRes.payload ?? [];
  } catch (e) {
    console.warn("加载分类或标签失败", e);
  }

  await refreshAIStatus();

  if (isEdit.value) {
    try {
      const postRes = await getPost(route.params.id);
      if (postRes.code === 200 && postRes.payload) {
        const p = postRes.payload;
        ruleForm.value = {
          postId: p.postId,
          postSlug: p.postSlug || (route.params.id as string) || "",
          title: p.title || "",
          categoryId: p.categoryId,
          tags: Array.isArray(p.tags) ? p.tags.map((t: any) => t.tagName || t) : [],
          summary: p.summary || "",
          status: p.status ?? 1,
          type: p.type ?? 1,
          postContent: p.postContent || "",
          postContentHtml: p.postContentHtml || "",
          coverImage: p.coverImage || "",
          series: p.series || ""
        };

        if (p.published === 1 && p.pubTime) {
          scheduleEnabled.value = true;
          scheduleTime.value = p.pubTime;
        }

        await tryRestoreDraft(p.updateTime ? new Date(p.updateTime).getTime() : undefined);
      }
    } catch (e) {
      console.warn("获取文章详情失败", e);
    }
  } else {
    await tryRestoreDraft();
  }
};

onMounted(() => {
  initData();
});
</script>

<style scoped lang="scss">
.editor-container {
  padding: 16px;
  min-height: 100vh;
  box-sizing: border-box;
  background-color: var(--el-bg-color-page);
  transition: all 0.3s ease;

  &.focus-mode {
    padding: 8px;
    background-color: var(--el-bg-color);

    .editor-card {
      box-shadow: none;
      border: none;
    }
  }
}

.editor-card {
  border-radius: 8px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;

    .card-title {
      font-size: 16px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    .header-actions {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;
    }
  }
}

.editor-form {
  margin-top: 8px;
}

.editor-wrap-relative {
  position: relative;
  width: 100%;
}

.markdown-editor:not(.md-editor-fullscreen) {
  isolation: isolate;
  position: relative;
  z-index: 1;
}

:deep(.md-editor-preview .md-editor-code .md-editor-code-head) {
  z-index: 2 !important;
}

:deep(.md-editor-catalog-fixed) {
  z-index: 10 !important;
}

:deep(.md-editor.md-editor-fullscreen) {
  z-index: 1500 !important;
}

.ghost-tip-capsule {
  position: fixed;
  z-index: 2100;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--el-bg-color-overlay, #ffffff);
  border: 1px solid var(--el-color-primary-light-5);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.12);
  padding: 3px 10px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 11px;
  color: var(--el-color-primary);
  animation: ghostPulse 1.8s infinite;

  kbd {
    background: var(--el-fill-color-darker);
    border-radius: 3px;
    padding: 1px 4px;
    font-size: 10px;
  }

  &:hover {
    background: var(--el-color-primary-light-9);
    border-color: var(--el-color-primary);
  }
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin-top: 6px;
  line-height: 1.4;

  code {
    background: var(--el-fill-color);
    padding: 1px 4px;
    border-radius: 3px;
    font-family: ui-monospace, SFMono-Regular, monospace;
    font-size: 11px;
  }

  .draft-indicator {
    color: var(--el-color-success);
    margin-left: 6px;
  }
}

.radio-group-block {
  display: flex;
  gap: 12px;

  :deep(.el-radio) {
    margin-right: 0;
    flex: 1;
    display: flex;
    justify-content: center;
  }
}

.ai-outline-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.08), rgba(103, 194, 58, 0.06));
  border: 1px dashed var(--el-color-primary-light-5);
  border-radius: 6px;
  padding: 8px 14px;
  margin-bottom: 8px;

  .banner-left {
    display: flex;
    align-items: center;
    gap: 8px;

    .banner-badge {
      font-size: 12px;
      font-weight: 600;
      color: var(--el-color-primary);
    }

    .banner-text {
      font-size: 12px;
      color: var(--el-text-color-regular);
    }
  }

  .banner-right {
    display: flex;
    align-items: center;
    gap: 6px;
  }
}

.ai-inline-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  border-radius: 6px;
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 8px;
  font-size: 12px;

  &.is-streaming {
    background: rgba(64, 158, 255, 0.08);
    border-color: rgba(64, 158, 255, 0.3);
  }

  .ai-inline-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--el-color-primary);
    animation: inlineDotPulse 1s infinite;
  }

  .ai-inline-label {
    font-weight: 500;
    color: var(--el-text-color-primary);
  }

  .ai-inline-reasoning {
    color: var(--el-text-color-placeholder);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 260px;
  }

  .ai-inline-actions {
    margin-left: auto;
    display: flex;
    gap: 6px;
  }
}

.ai-toolbar-trigger {
  font-weight: 600;
  color: var(--el-color-primary);
  cursor: pointer;
}

.ai-menu-overlay {
  list-style: none;
  padding: 6px 0;
  margin: 0;
  min-width: 240px;

  li {
    padding: 7px 14px;
    font-size: 13px;
    color: var(--el-text-color-primary);
    cursor: pointer;
    transition: background 0.15s ease;

    &:hover:not(.disabled):not(.ai-menu-custom-box):not(.ai-menu-section-header):not(.ai-menu-models) {
      background: var(--el-fill-color-light);
      color: var(--el-color-primary);
    }

    &.disabled {
      color: var(--el-text-color-disabled);
      cursor: not-allowed;
    }
  }

  .ai-menu-custom-box {
    padding: 8px 12px 10px;
    cursor: default;

    .ai-custom-prompt-wrap {
      display: flex;
      flex-direction: column;
      gap: 6px;

      .ai-send-icon {
        cursor: pointer;
        color: var(--el-text-color-placeholder);
        &.active {
          color: var(--el-color-primary);
        }
      }

      .ai-quick-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;

        .ai-quick-tag {
          font-size: 10px;
          padding: 1px 6px;
          border-radius: 3px;
          background: var(--el-fill-color);
          color: var(--el-text-color-secondary);
          cursor: pointer;

          &:hover {
            background: var(--el-color-primary-light-9);
            color: var(--el-color-primary);
          }
        }
      }
    }
  }

  .ai-menu-divider {
    height: 1px;
    padding: 0;
    margin: 4px 0;
    background: var(--el-border-color-extra-light);
  }

  .ai-menu-section-header {
    padding: 4px 14px;
    font-size: 11px;
    font-weight: 600;
    color: var(--el-text-color-placeholder);
    cursor: default;
  }

  .ai-menu-models {
    padding: 6px 14px;
    cursor: default;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .ai-model-chips {
      display: flex;
      gap: 4px;

      .ai-model-chip {
        font-size: 10px;
        padding: 1px 6px;
        border-radius: 3px;
        background: var(--el-fill-color);
        color: var(--el-text-color-secondary);
        cursor: pointer;

        &.active {
          background: var(--el-color-primary-light-8);
          color: var(--el-color-primary);
          font-weight: 600;
        }
      }
    }
  }

  .ai-menu-action {
    font-size: 12px;
    color: var(--el-color-primary);
  }

  .ai-menu-quota {
    font-size: 11px;
    color: var(--el-text-color-placeholder);
    padding: 4px 14px;
    cursor: default;
  }
}

.ai-title-pop {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 220px;
  overflow-y: auto;

  .ai-title-tip {
    font-size: 12px;
    color: var(--el-text-color-placeholder);
    text-align: center;
    padding: 10px 0;
  }

  .ai-title-item {
    font-size: 13px;
    padding: 6px 8px;
    border-radius: 4px;
    cursor: pointer;
    color: var(--el-text-color-primary);

    &:hover {
      background: var(--el-color-primary-light-9);
      color: var(--el-color-primary);
    }
  }
}

/* CodeMirror 6 波浪线高亮样式 */
:deep(.cm-proofread-typo) {
  text-decoration: underline wavy #ef4444 2px !important;
  text-underline-offset: 3px !important;
  background: rgba(239, 68, 68, 0.08) !important;
  cursor: pointer !important;
}

:deep(.cm-proofread-grammar) {
  text-decoration: underline wavy #f59e0b 2px !important;
  text-underline-offset: 3px !important;
  background: rgba(245, 158, 11, 0.08) !important;
  cursor: pointer !important;
}

:deep(.cm-proofread-style) {
  text-decoration: underline wavy #3b82f6 2px !important;
  text-underline-offset: 3px !important;
  background: rgba(59, 130, 246, 0.08) !important;
  cursor: pointer !important;
}

@keyframes ghostPulse {
  0% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-2px);
  }
  100% {
    transform: translateY(0);
  }
}

@keyframes inlineDotPulse {
  0% {
    opacity: 0.4;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.4;
  }
}
</style>
