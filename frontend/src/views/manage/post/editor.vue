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
              title="一键智能生成推荐标题、核心摘要与推荐标签"
              @click="openPublishCopilot"
            >
              <el-icon class="mr-2px"><MagicStick /></el-icon>
              AI 发文助手
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
                        title="AI 生成标题"
                        @click="generateTitle"
                        >✨</el-button
                      >
                    </template>
                  </el-popover>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="文章正文" prop="postContent">
              <!-- 续写光标流式插入的状态条：生成中/完成提示 + 停止/撤销/重试 -->
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
                    >{{ inlineReasoning }}</span
                  >
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
                placeholder="请输入文章摘要"
                maxlength="500"
                show-word-limit
              />
              <div class="form-tip">
                如果不填写，系统将自动从正文中提取
                <el-button
                  v-if="aiEnabled"
                  link
                  type="primary"
                  size="small"
                  :loading="aiSummaryLoading"
                  style="margin-left: 6px"
                  @click="generateSummary"
                >
                  ✨ AI 生成摘要
                </el-button>
              </div>
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
                    placeholder="选填，如：Spring Boot 入门系列"
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

    <!-- AI 结果抽屉：流式渲染生成内容，人工确认后才写回正文 -->
    <el-drawer
      v-model="aiDrawerVisible"
      :title="aiTaskLabel"
      :size="aiLastRequest?.task === 'proofread' ? '540px' : '480px'"
      :close-on-click-modal="!aiStreaming"
      :before-close="closeAiDrawer"
    >
      <div ref="aiResultBodyRef" class="ai-result-body" @scroll="onAiBodyScroll">
        <el-alert
          v-if="aiError"
          :title="aiError"
          type="error"
          show-icon
          :closable="false"
          class="ai-error-banner"
        />

        <!-- 任务为全文校对时的专属交互界面 -->
        <template v-if="aiLastRequest?.task === 'proofread'">
          <!-- 校对控制条 -->
          <div class="proofread-header-bar">
            <div class="proofread-stats">
              <span class="proofread-total-badge">
                发现 {{ proofreadItems.length }} 处问题
              </span>
              <span v-if="pendingProofreadCount > 0" class="text-orange-500 text-xs">
                ({{ pendingProofreadCount }} 待处理)
              </span>
              <span v-if="appliedProofreadCount > 0" class="text-green-600 text-xs">
                ({{ appliedProofreadCount }} 已采纳)
              </span>
            </div>

            <div class="proofread-header-actions">
              <el-button
                v-if="pendingProofreadCount > 0 && !aiStreaming"
                size="small"
                type="primary"
                plain
                @click="applyAllProofreadItems"
              >
                一键采纳全部 ({{ pendingProofreadCount }})
              </el-button>
              <el-radio-group v-model="proofreadViewMode" size="small">
                <el-radio-button label="cards">卡片视图</el-radio-button>
                <el-radio-button label="raw">原始输出</el-radio-button>
              </el-radio-group>
            </div>
          </div>

          <!-- 卡片视图下的筛选与列表 -->
          <div v-if="proofreadViewMode === 'cards'" class="proofread-cards-container">
            <div v-if="proofreadItems.length > 0" class="proofread-filter-row">
              <el-radio-group v-model="proofreadFilter" size="small">
                <el-radio-button label="all">全部 ({{ proofreadItems.length }})</el-radio-button>
                <el-radio-button label="pending">待处理 ({{ pendingProofreadCount }})</el-radio-button>
                <el-radio-button label="applied">已采纳 ({{ appliedProofreadCount }})</el-radio-button>
              </el-radio-group>
            </div>

            <div v-if="proofreadItems.length === 0 && !aiStreaming" class="proofread-empty">
              <el-icon class="text-green-500 text-3xl mb-2"><CircleCheck /></el-icon>
              <div>未发现明显错别字或语病，全文表述良好！</div>
            </div>

            <div class="proofread-cards-list">
              <div
                v-for="item in filteredProofreadItems"
                :key="item.id"
                class="proofread-card"
                :class="{
                  'is-applied': item.status === 'applied',
                  'is-ignored': item.status === 'ignored',
                  'is-not-found': item.status === 'not_found'
                }"
              >
                <div class="proofread-card-header">
                  <el-tag
                    size="small"
                    :type="getProofreadTagType(item.type)"
                    effect="light"
                  >
                    {{ item.type }}
                  </el-tag>
                  <span v-if="item.reason" class="proofread-reason" :title="item.reason">
                    {{ item.reason }}
                  </span>
                  <span class="proofread-status-indicator">
                    <span v-if="item.status === 'applied'" class="text-green-600 font-medium">✓ 已采纳</span>
                    <span v-else-if="item.status === 'ignored'" class="text-gray-400">已忽略</span>
                    <span v-else-if="item.status === 'not_found'" class="text-orange-400">未在正文中匹配</span>
                  </span>
                </div>

                <div class="proofread-card-diff">
                  <div class="diff-line">
                    <span class="diff-tag del">原</span>
                    <span class="diff-text del">{{ item.original }}</span>
                  </div>
                  <div class="diff-line">
                    <span class="diff-tag ins">改</span>
                    <span class="diff-text ins">{{ item.suggestion }}</span>
                  </div>
                </div>

                <div class="proofread-card-footer">
                  <el-button
                    link
                    size="small"
                    type="primary"
                    @click="locateProofreadItem(item)"
                  >
                    <el-icon class="mr-2px"><Search /></el-icon>
                    定位原文
                  </el-button>
                  <template v-if="item.status === 'pending'">
                    <el-button
                      size="small"
                      type="primary"
                      @click="applyProofreadItem(item)"
                    >
                      采纳修改
                    </el-button>
                    <el-button
                      link
                      size="small"
                      @click="ignoreProofreadItem(item)"
                    >
                      忽略
                    </el-button>
                  </template>
                </div>
              </div>
            </div>
          </div>

          <!-- 原始输出视图 -->
          <MdPreview
            v-else
            :model-value="aiResultText || (aiStreaming ? '正在校对全文…' : '')"
            preview-theme="github"
            code-theme="atom"
            class="ai-result-preview"
          />
        </template>

        <!-- 任务为润色/改写等其他任务时的常规视图 -->
        <template v-else>
          <!-- 差异对比 / 最终效果 视图切换栏（存在原始选区时展示） -->
          <div v-if="aiOriginalSelection" class="ai-view-switch-row">
            <el-radio-group v-model="aiViewMode" size="small">
              <el-radio-button label="diff">
                <el-icon class="mr-2px"><DocumentCopy /></el-icon>
                差异对比
              </el-radio-button>
              <el-radio-button label="preview">
                <el-icon class="mr-2px"><View /></el-icon>
                最终效果
              </el-radio-button>
            </el-radio-group>
            <span v-if="aiViewMode === 'diff'" class="diff-legend">
              <span class="legend-badge legend-del">红色删除</span>
              <span class="legend-badge legend-ins">绿色新增</span>
            </span>
          </div>

          <div v-if="aiReasoningText || (aiStreaming && !aiResultText)" class="ai-reasoning">
            <div class="ai-reasoning-title">
              <span class="ai-reasoning-dot" :class="{ pulse: aiStreaming && !aiResultText }"></span>
              思考过程{{ aiResultText ? "（已输出正文，可展开回看）" : "…" }}
            </div>
            <div class="ai-reasoning-text">{{ aiReasoningText }}</div>
          </div>

          <!-- 差异对比视图 -->
          <div
            v-if="aiOriginalSelection && aiViewMode === 'diff'"
            class="ai-diff-container"
          >
            <div v-if="!aiResultText && aiStreaming" class="ai-diff-placeholder">
              正在生成并计算差异…
            </div>
            <div v-else class="ai-diff-content">
              <template v-for="(chunk, idx) in diffChunks" :key="idx">
                <del v-if="chunk.type === 'removed'" class="diff-chunk diff-del">{{ chunk.value }}</del>
                <ins v-else-if="chunk.type === 'added'" class="diff-chunk diff-ins">{{ chunk.value }}</ins>
                <span v-else class="diff-chunk diff-common">{{ chunk.value }}</span>
              </template>
            </div>
          </div>

          <!-- 原 Markdown 预览视图 -->
          <MdPreview
            v-else
            :model-value="aiResultText || (aiStreaming && !aiReasoningText ? '生成中…' : '')"
            preview-theme="github"
            code-theme="atom"
            class="ai-result-preview"
          />
        </template>
      </div>
      <template #footer>
        <div class="ai-result-foot">
          <!-- 追加指令：以上一版为基础继续调整 -->
          <div
            v-if="!aiStreaming && aiResultText && !aiError"
            class="ai-refine-row"
          >
            <el-input
              v-model="aiRefineInput"
              size="small"
              placeholder="追加要求：再精简一点 / 补个示例 / 语气更正式"
              clearable
              @keydown.enter="refineAI"
            />
            <el-button
              size="small"
              type="primary"
              :disabled="!aiRefineInput.trim()"
              @click="refineAI"
            >
              继续调整
            </el-button>
          </div>
          <div class="ai-result-meta">
            <span v-if="aiHistory.length > 1" class="ai-version-nav">
              <el-button
                link
                size="small"
                :disabled="aiHistoryIdx <= 0"
                @click="aiNavHistory(-1)"
              >
                ‹
              </el-button>
              {{ aiHistoryIdx + 1 }}/{{ aiHistory.length }}
              <el-button
                link
                size="small"
                :disabled="aiHistoryIdx >= aiHistory.length - 1"
                @click="aiNavHistory(1)"
              >
                ›
              </el-button>
            </span>
            <span v-if="aiMeta && !aiError">
              {{ aiMeta.elapsed
              }}{{ aiMeta.tokens ? ` · ${aiMeta.tokens} tokens` : "" }}
            </span>
            <span v-if="aiQuotaText" class="ai-quota-text">
              {{ aiQuotaText }}
            </span>
          </div>
          <div class="ai-result-actions">
            <!-- 全文校对任务的批量采纳按钮 -->
            <el-button
              v-if="aiLastRequest?.task === 'proofread' && pendingProofreadCount > 0"
              type="primary"
              :disabled="aiStreaming"
              @click="applyAllProofreadItems"
            >
              一键采纳全部 ({{ pendingProofreadCount }})
            </el-button>
            <el-button
              v-if="aiApplyType === 'replace-selection'"
              type="primary"
              :disabled="aiStreaming || !aiResultText"
              @click="applyAIResult()"
            >
              替换选区
            </el-button>
            <el-button
              v-if="aiApplyType === 'replace-selection'"
              :disabled="aiStreaming || !aiResultText"
              title="保留原选区文本，在后方追加换行并插入本次生成内容"
              @click="insertBelowAIResult()"
            >
              在选区后插入
            </el-button>
            <el-button
              :disabled="aiStreaming || !aiResultText"
              @click="copyAIResult"
            >
              复制
            </el-button>
            <el-button
              v-if="!aiStreaming && !aiError"
              :disabled="!aiResultText"
              @click="regenAI"
            >
              重新生成
            </el-button>
            <el-button
              v-if="aiError && !aiStreaming"
              type="primary"
              plain
              @click="regenAI"
            >
              重试
            </el-button>
            <el-button
              v-if="aiStreaming"
              type="danger"
              plain
              @click="stopAI"
            >
              停止
            </el-button>
          </div>
        </div>
      </template>
    </el-drawer>

    <!-- AI 指令模板管理：覆盖任务内置指令，留空恢复默认 -->
    <el-dialog
      v-model="aiTplVisible"
      title="AI 指令模板"
      width="620px"
      append-to-body
    >
      <el-alert
        type="info"
        :closable="false"
        show-icon
        title="自定义指令会替换该任务的内置提示词；文章内容等上下文仍由系统自动附加。留空保存即恢复默认。"
        style="margin-bottom: 14px"
      />
      <el-select v-model="aiTplKey" style="width: 100%" @change="onTplKeyChange">
        <el-option
          v-for="t in aiTplItems"
          :key="t.key"
          :label="aiTplLabelMap[t.key] || t.key"
          :value="t.key"
        />
      </el-select>
      <div v-if="aiTplCurrent" class="ai-tpl-default">
        <div class="ai-tpl-default-title">内置默认指令</div>
        <div class="ai-tpl-default-text">{{ aiTplCurrent.defaultContent }}</div>
      </div>
      <el-input
        v-model="aiTplContent"
        type="textarea"
        :rows="5"
        maxlength="500"
        show-word-limit
        placeholder="输入自定义指令（留空使用默认）"
        style="margin-top: 12px"
      />
      <template #footer>
        <el-button :disabled="!aiTplContent" @click="saveTpl(true)">
          恢复默认
        </el-button>
        <el-button
          type="primary"
          :loading="aiTplSaving"
          :disabled="!aiTplContent && !aiTplCurrent?.content"
          @click="saveTpl(false)"
        >
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- AI 发文助手一站式弹窗 -->
    <el-dialog
      v-model="copilotVisible"
      title="✨ AI 发文助手"
      width="680px"
      append-to-body
      class="ai-copilot-dialog"
      :before-close="closeCopilot"
    >
      <div v-loading="copilotLoading" element-loading-text="AI 正在分析全文提炼发文元数据…" class="copilot-body">
        <el-alert
          type="info"
          :closable="false"
          show-icon
          title="一键基于全文提炼推荐标题、摘要与精准匹配标签，确认后可「一键应用到文章」。"
          style="margin-bottom: 16px"
        />

        <!-- 1. 推荐标题 -->
        <div class="copilot-section">
          <div class="copilot-section-header">
            <span class="copilot-section-title">
              <el-icon class="mr-4px text-primary"><CollectionTag /></el-icon>
              推荐标题
            </span>
            <el-button
              link
              type="primary"
              size="small"
              :loading="copilotTitleLoading"
              @click="copilotRegenTitle"
            >
              换一批
            </el-button>
          </div>
          <div class="copilot-title-list">
            <el-radio-group v-model="copilotChosenTitle" class="copilot-radio-group">
              <el-radio
                v-for="(t, idx) in copilotTitleCandidates"
                :key="idx"
                :label="t"
                class="copilot-title-radio"
              >
                {{ t }}
              </el-radio>
              <el-radio
                v-if="ruleForm.title"
                :label="ruleForm.title"
                class="copilot-title-radio text-gray-500"
              >
                保持当前原标题：{{ ruleForm.title }}
              </el-radio>
            </el-radio-group>
          </div>
        </div>

        <!-- 2. 文章摘要 -->
        <div class="copilot-section">
          <div class="copilot-section-header">
            <span class="copilot-section-title">
              <el-icon class="mr-4px text-primary"><Tickets /></el-icon>
              文章摘要
            </span>
            <el-button
              link
              type="primary"
              size="small"
              :loading="copilotSummaryLoading"
              @click="copilotRegenSummary"
            >
              重新提炼
            </el-button>
          </div>
          <el-input
            v-model="copilotSummary"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
            placeholder="AI 正在生成摘要…"
          />
        </div>

        <!-- 3. 标签匹配与建议 -->
        <div class="copilot-section">
          <div class="copilot-section-header">
            <span class="copilot-section-title">
              <el-icon class="mr-4px text-primary"><PriceTag /></el-icon>
              标签建议（智能比对标签库）
            </span>
            <el-button
              link
              type="primary"
              size="small"
              :loading="copilotTagsLoading"
              @click="copilotRegenTags"
            >
              重新建议
            </el-button>
          </div>

          <!-- 已有标签库中匹配的（推荐复用，默认高亮） -->
          <div v-if="copilotMatchedTags.length > 0" class="copilot-tag-subgroup">
            <div class="subgroup-label">命中系统已有标签库（建议勾选复用）：</div>
            <div class="copilot-tags-wrap">
              <el-check-tag
                v-for="t in copilotMatchedTags"
                :key="t"
                :checked="copilotSelectedTags.includes(t)"
                class="copilot-tag-item matched"
                @change="toggleCopilotTag(t)"
              >
                <el-icon class="mr-2px"><Check /></el-icon>
                {{ t }}
              </el-check-tag>
            </div>
          </div>

          <!-- 推荐新建的标签 -->
          <div v-if="copilotNewTags.length > 0" class="copilot-tag-subgroup">
            <div class="subgroup-label">推荐新建标签：</div>
            <div class="copilot-tags-wrap">
              <el-check-tag
                v-for="t in copilotNewTags"
                :key="t"
                :checked="copilotSelectedTags.includes(t)"
                class="copilot-tag-item new-tag"
                @change="toggleCopilotTag(t)"
              >
                + {{ t }}
              </el-check-tag>
            </div>
          </div>
          <div v-if="copilotMatchedTags.length === 0 && copilotNewTags.length === 0 && !copilotTagsLoading" class="text-xs text-gray-400">
            暂无标签建议
          </div>
        </div>
      </div>

      <template #footer>
        <div class="copilot-foot">
          <div class="copilot-foot-tip text-xs text-gray-400">
            已选：标题 {{ copilotChosenTitle ? '✓' : '-' }} · 摘要 {{ copilotSummary ? '✓' : '-' }} · 标签 ({{ copilotSelectedTags.length }})
          </div>
          <div class="copilot-foot-actions">
            <el-button @click="copilotVisible = false">取消</el-button>
            <el-button
              type="primary"
              :disabled="copilotLoading || (!copilotChosenTitle && !copilotSummary && copilotSelectedTags.length === 0)"
              @click="applyCopilotToForm"
            >
              一键应用到文章
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>

    <!-- 底部固定操作栏：长文写到结尾无需滚回顶部即可发布 -->
    <div class="editor-footer">
      <span v-if="draftSavedAtText" class="draft-indicator">
        <el-icon><DocumentChecked /></el-icon>
        草稿已保存 {{ draftSavedAtText }}
      </span>
      <span class="footer-spacer" />
      <el-button
        v-if="aiEnabled"
        type="success"
        plain
        title="一键智能生成推荐标题、核心摘要与推荐标签"
        @click="openPublishCopilot"
      >
        <el-icon class="mr-2px"><MagicStick /></el-icon>
        AI 发文助手
      </el-button>
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
  Check,
  CircleCheck,
  CollectionTag,
  Delete,
  Document,
  DocumentChecked,
  DocumentCopy,
  Download,
  Link,
  MagicStick,
  Picture,
  Plus,
  PriceTag,
  Promotion,
  Right,
  Search,
  Tickets,
  Upload,
  View,
  ZoomIn
} from "@element-plus/icons-vue";
import {
  MdEditor,
  MdPreview,
  DropdownToolbar,
  type ToolbarNames
} from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import "md-editor-v3/lib/preview.css";
import {
  getAIStatus,
  streamAIEdit,
  getAITemplates,
  saveAITemplate,
  type AIEditRequest,
  type AIStatus,
  type AIStreamHandlers,
  type AITemplateItem
} from "@/api/ai";
import { computeDiff, type DiffChunk } from "@/utils/diff";
import { parseProofreadOutput, type ProofreadItem } from "@/utils/proofread";
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

// 编辑器实例：AI 取选区/回写正文都经由它的 CodeMirror view
const editorRef = ref<any>(null);

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
  coverImage: "",
  series: ""
});

// 定时发布：开启后选择未来时间，保存时随表单提交 published=1 + pubTime
const scheduleEnabled = ref(false);
const scheduleTime = ref<string | number>("");

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

// ==================== AI 辅助写作 ====================
const aiStatus = ref<AIStatus | null>(null);
const aiEnabled = computed(() => !!aiStatus.value?.enabled);

// 抽屉状态。aiResultText/aiReasoningText 是节流后的展示值；
// 原始增量先积攒在非响应式缓冲里，定时合并刷新，避免每个 token
// 触发一次 MdPreview 全量重渲染造成卡顿
const aiDrawerVisible = ref(false);
const aiStreaming = ref(false);
const aiResultText = ref("");
const aiReasoningText = ref("");
const aiResultBodyRef = ref<any>(null);
let aiResultBuf = "";
let aiReasoningBuf = "";
let aiFlushTimer: ReturnType<typeof setTimeout> | null = null;
// 跟随滚动：用户主动上翻时暂停贴底，回到底部附近再恢复
let aiStickBottom = true;
const aiTaskLabel = ref("AI 助手");
// 本次结果的回写方式与选区快照（发起请求时记录，应用时使用）
const aiApplyType = ref<"replace-selection" | "none">("none");
const aiSelFrom = ref(0);
const aiSelTo = ref(0);
let aiAbort: (() => void) | null = null;
let aiLastRequest: AIEditRequest | null = null;
let aiStartedAt = 0;

// 选区原始快照与差异对比视图
const aiOriginalSelection = ref("");
const aiViewMode = ref<"diff" | "preview">("diff");
const diffChunks = computed<DiffChunk[]>(() => {
  if (!aiOriginalSelection.value) return [];
  return computeDiff(aiOriginalSelection.value, aiResultText.value || "");
});

// 全文校对卡片状态
const proofreadItems = ref<ProofreadItem[]>([]);
const proofreadFilter = ref<"all" | "pending" | "applied">("all");
const proofreadViewMode = ref<"cards" | "raw">("cards");

const filteredProofreadItems = computed(() => {
  if (proofreadFilter.value === "pending") {
    return proofreadItems.value.filter(i => i.status === "pending");
  }
  if (proofreadFilter.value === "applied") {
    return proofreadItems.value.filter(i => i.status === "applied");
  }
  return proofreadItems.value;
});

const pendingProofreadCount = computed(
  () => proofreadItems.value.filter(i => i.status === "pending").length
);
const appliedProofreadCount = computed(
  () => proofreadItems.value.filter(i => i.status === "applied").length
);

// 选区自定义指令（Ask AI）
const aiCustomPrompt = ref("");
const aiQuickPromptTags = [
  "转为表格",
  "提炼核心要点",
  "更口语化",
  "更严谨专业",
  "翻译为英文"
];

// AI 发文助手状态
const copilotVisible = ref(false);
const copilotLoading = ref(false);
const copilotTitleLoading = ref(false);
const copilotSummaryLoading = ref(false);
const copilotTagsLoading = ref(false);
const copilotTitleCandidates = ref<string[]>([]);
const copilotChosenTitle = ref("");
const copilotSummary = ref("");
const copilotMatchedTags = ref<string[]>([]);
const copilotNewTags = ref<string[]>([]);
const copilotSelectedTags = ref<string[]>([]);
let copilotAbortTitle: (() => void) | null = null;
let copilotAbortSummary: (() => void) | null = null;
let copilotAbortTags: (() => void) | null = null;

// 抽屉元信息：耗时与 token（done 事件带出，此前被丢弃）
const aiMeta = ref<{ elapsed: string; tokens: number } | null>(null);
// 抽屉内错误横幅（非流式信封错误与流中错误都会写入，配重试按钮）
const aiError = ref("");
// 版本栈：regen/refine 的历次结果，‹ › 导航回看
const aiHistory = ref<string[]>([]);
const aiHistoryIdx = ref(-1);
// 抽屉追加指令
const aiRefineInput = ref("");

// 配额展示：菜单底部与抽屉 footer 共用；用尽时禁用入口
const aiQuotaText = computed(() => {
  const s = aiStatus.value;
  if (!s?.enabled) return "";
  const used =
    s.dailyQuota > 0 ? `今日 ${s.todayUsed}/${s.dailyQuota} 次` : `今日 ${s.todayUsed} 次`;
  return `${used} · ${s.model}`;
});
const aiQuotaExhausted = computed(() => {
  const s = aiStatus.value;
  return !!s?.enabled && s.dailyQuota > 0 && s.todayUsed >= s.dailyQuota;
});

// ==================== 光标流式插入（续写类） ====================
// 生成内容直接写进编辑器，不再走抽屉；配浮动状态条（停止/撤销/重试）
const inlineStreaming = ref(false);
const inlineReasoning = ref("");
const inlineUndoable = ref(false);
const inlineDoneTip = ref(""); // 完成后短暂显示的提示
let inlineView: any = null;
let inlineStart = 0;
let inlinePos = 0;
let inlinePrefix = "";
let inlineWroteText = ""; // 前缀+全部已写入文本，撤销时按内容校验
let inlinePending = "";
let inlineReasoningBuf = "";
let inlineTimer: ReturnType<typeof setTimeout> | null = null;

// 光标续写的上下文窗口（与后端 prompt.go 窗口对齐）
const CURSOR_BEFORE_WINDOW = 12000;
const CURSOR_AFTER_WINDOW = 800;

// ===== 多模型切换：选择持久化，所有任务统一注入 =====
const AI_MODEL_STORAGE = "keepblog-ai-model";
const aiModelChoice = ref<string>("");
try {
  aiModelChoice.value = localStorage.getItem(AI_MODEL_STORAGE) || "";
} catch {
  /* 隐私模式等场景忽略 */
}
const aiModelOptions = computed(() => aiStatus.value?.models ?? []);
function setAIModel(name: string) {
  aiModelChoice.value = name;
  try {
    localStorage.setItem(AI_MODEL_STORAGE, name);
  } catch {
    /* 忽略 */
  }
}
// 统一流式出口：所有任务注入当前选择的模型名
function aiStream(req: AIEditRequest, handlers: AIStreamHandlers) {
  return streamAIEdit(
    { ...req, model: aiModelChoice.value || undefined },
    handlers
  );
}

// ===== 指令模板管理 =====
const aiTplVisible = ref(false);
const aiTplItems = ref<AITemplateItem[]>([]);
const aiTplKey = ref("polish:polish");
const aiTplContent = ref("");
const aiTplSaving = ref(false);
const aiTplLabelMap: Record<string, string> = {
  "polish:polish": "润色",
  "polish:expand": "扩写",
  "polish:shorten": "精简",
  "polish:translate": "翻译",
  continue: "续写",
  title: "标题生成",
  summary: "摘要生成",
  refine: "追加调整",
  tags: "标签建议",
  proofread: "全文校对"
};
const aiTplCurrent = computed(
  () => aiTplItems.value.find(t => t.key === aiTplKey.value) ?? null
);
async function openTplDialog() {
  aiMenuVisible.value = false;
  aiTplVisible.value = true;
  try {
    const res = await getAITemplates();
    if (res.code === 200) {
      aiTplItems.value = res.payload ?? [];
      const cur = aiTplItems.value.find(t => t.key === aiTplKey.value);
      aiTplContent.value = cur?.content ?? "";
    }
  } catch {
    /* 列表失败保持空 */
  }
}
function onTplKeyChange() {
  aiTplContent.value = aiTplCurrent.value?.content ?? "";
}
async function saveTpl(reset = false) {
  aiTplSaving.value = true;
  try {
    const content = reset ? "" : aiTplContent.value;
    const res = await saveAITemplate(aiTplKey.value, content);
    if (res.code === 200) {
      const cur = aiTplItems.value.find(t => t.key === aiTplKey.value);
      if (cur) cur.content = reset ? "" : aiTplContent.value.trim();
      if (reset) aiTplContent.value = "";
      message(reset ? "已恢复默认指令" : "模板已保存", { type: "success" });
    }
  } finally {
    aiTplSaving.value = false;
  }
}

// 下拉菜单数据
const aiMenuVisible = ref(false);
const aiMenuData = [
  { label: "润色选中文本", value: "polish" },
  { label: "扩写选中文本", value: "expand" },
  { label: "精简选中文本", value: "shorten" },
  { label: "翻译为英文", value: "translate" },
  { label: "续写正文（光标处）", value: "continue" },
  { label: "全文校对", value: "proofread" }
];

// 标题候选
const titlePopoverVisible = ref(false);
const aiTitleLoading = ref(false);
const aiTitleCandidates = ref<string[]>([]);
const aiSummaryLoading = ref(false);

// 标签候选
const tagsPopoverVisible = ref(false);
const aiTagsLoading = ref(false);
const aiTagCandidates = ref<string[]>([]);

// 组件卸载时中止进行中的生成
onBeforeUnmount(() => {
  aiAbort?.();
  copilotAbortTitle?.();
  copilotAbortSummary?.();
  copilotAbortTags?.();
  if (inlineTimer) {
    clearTimeout(inlineTimer);
    inlineTimer = null;
  }
});

// 探测 AI 状态（未配置则所有入口隐藏）
onMounted(async () => {
  try {
    const res = await getAIStatus();
    if (res.code === 200) aiStatus.value = res.payload;
  } catch {
    /* 探测失败按未启用处理 */
  }
});

// 每次生成完成后刷新配额（今日 x/y）
async function refreshAIStatus() {
  try {
    const res = await getAIStatus();
    if (res.code === 200) aiStatus.value = res.payload;
  } catch {
    /* 刷新失败保持旧值 */
  }
}

// 编辑器选区快照
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

function onAiMenuClick(item: any) {
  aiMenuVisible.value = false;
  if (aiQuotaExhausted.value) {
    message("今日 AI 调用配额已用完，明天再来吧", { type: "warning" });
    return;
  }
  const value = String(item?.value || "polish");
  if (value === "continue") {
    runContinue();
    return;
  }
  if (value === "proofread") {
    runProofread();
    return;
  }
  const modes = ["polish", "expand", "shorten", "translate"] as const;
  runPolish(
    (modes as readonly string[]).includes(value)
      ? (value as (typeof modes)[number])
      : "polish"
  );
}

// 选区润色/扩写/精简/翻译
function runPolish(mode: "polish" | "expand" | "shorten" | "translate") {
  const sel = captureSelection();
  if (!sel || !sel.text.trim()) {
    message("请先选中要处理的正文内容", { type: "info" });
    return;
  }
  aiSelFrom.value = sel.from;
  aiSelTo.value = sel.to;
  aiOriginalSelection.value = sel.text;
  aiViewMode.value = "diff";
  const labelMap: Record<string, string> = {
    polish: "AI 润色",
    expand: "AI 扩写",
    shorten: "AI 精简",
    translate: "AI 翻译"
  };
  startAI(
    {
      task: "polish",
      mode,
      selection: sel.text,
      // 窗口与后端 prompt.go 的 beforeWindow/afterWindow 对齐
      before: sel.view.state.sliceDoc(
        Math.max(0, sel.from - 1500),
        sel.from
      ),
      after: sel.view.state.sliceDoc(
        sel.to,
        Math.min(sel.view.state.doc.length, sel.to + 800)
      )
    },
    labelMap[mode] || "AI 润色",
    "replace-selection"
  );
}

// 选区自定义指令（Ask AI）执行
function runCustomSelectionPrompt() {
  const prompt = aiCustomPrompt.value.trim();
  if (!prompt) return;
  if (aiQuotaExhausted.value) {
    message("今日 AI 调用配额已用完", { type: "warning" });
    return;
  }
  const sel = captureSelection();
  if (!sel || !sel.text.trim()) {
    message("请先选中要处理的正文内容", { type: "info" });
    return;
  }
  aiMenuVisible.value = false;
  aiSelFrom.value = sel.from;
  aiSelTo.value = sel.to;
  aiOriginalSelection.value = sel.text;
  aiViewMode.value = "diff";
  aiCustomPrompt.value = "";
  startAI(
    {
      task: "refine",
      previous: sel.text,
      instruction: prompt,
      title: ruleForm.value.title,
      series: ruleForm.value.series
    },
    `AI：${prompt.length > 10 ? prompt.slice(0, 10) + "…" : prompt}`,
    "replace-selection"
  );
}

function applyQuickPrompt(tag: string) {
  aiCustomPrompt.value = tag;
  runCustomSelectionPrompt();
}

// 续写：光标处流式插入。上下文取光标前后窗口（后端有 Before 时不再整篇截尾），
// 生成内容直接写入编辑器，抽屉流程不介入
function runContinue() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先填写正文内容", { type: "info" });
    return;
  }
  const sel = captureSelection();
  const view = sel?.view ?? editorRef.value?.getEditorView?.();
  if (!view) return;
  const cursor = sel?.to ?? view.state.doc.length;

  inlineView = view;
  inlineStart = cursor;
  inlinePos = cursor;
  inlineWroteText = "";
  inlinePending = "";
  inlineReasoningBuf = "";
  inlineReasoning.value = "";
  inlineUndoable.value = false;
  inlineDoneTip.value = "";
  // 光标不在行首/文首时补换行，与原先抽屉确认后插入的格式一致
  const prevCh = cursor > 0 ? view.state.sliceDoc(cursor - 1, cursor) : "";
  inlinePrefix = cursor === 0 ? "" : prevCh === "\n" ? "\n" : "\n\n";

  inlineStreaming.value = true;
  aiStartedAt = Date.now();
  aiLastRequest = {
    task: "continue",
    before: view.state.sliceDoc(
      Math.max(0, cursor - CURSOR_BEFORE_WINDOW),
      cursor
    ),
    after: view.state.sliceDoc(
      cursor,
      Math.min(view.state.doc.length, cursor + CURSOR_AFTER_WINDOW)
    ),
    title: ruleForm.value.title,
    series: ruleForm.value.series
  };
  aiAbort = aiStream(aiLastRequest, {
    onDelta: text => {
      inlinePending += text;
      scheduleInlineFlush();
    },
    onReasoning: text => {
      inlineReasoningBuf += text;
      scheduleInlineFlush();
    },
    onDone: usage => {
      flushInline(true);
      inlineStreaming.value = false;
      inlineUndoable.value = inlineWroteText.length > 0;
      inlineDoneTip.value = `已续写 ${inlineWroteText.trim().length} 字${
        mkElapsedSuffix(usage)
      }，不满意可撤销`;
      refreshAIStatus();
      // 1 分钟后收起状态条（撤销入口仍在菜单里，重发起即覆盖）
      setTimeout(() => {
        if (!inlineStreaming.value) inlineUndoable.value = false;
      }, 60000);
    },
    onError: msg => {
      flushInline(true);
      inlineStreaming.value = false;
      if (inlineWroteText) inlineUndoable.value = true;
      inlineDoneTip.value = "生成失败，可撤销已写部分后重试";
      message(msg, { type: "error" });
    }
  });
}

// —— 光标流式插入的低层操作 ——

function scheduleInlineFlush() {
  if (inlineTimer) return;
  inlineTimer = setTimeout(() => {
    inlineTimer = null;
    flushInline();
  }, 120);
}

function flushInline(final = false) {
  if (!inlineView) return;
  if (inlinePending) {
    const first = inlineWroteText.length === 0;
    const chunk = first ? inlinePrefix + inlinePending : inlinePending;
    inlineView.dispatch({
      changes: { from: inlinePos, to: inlinePos, insert: chunk },
      // 光标只在首末次归位：中途抢光标会和用户手动输入打架
      ...(first || final
        ? { selection: { anchor: inlinePos + chunk.length } }
        : {}),
      scrollIntoView: true
    });
    inlinePos += chunk.length;
    inlineWroteText += chunk;
    inlinePending = "";
  }
  // 思考气泡只留末尾片段，避免浮动条被长思考撑爆
  if (inlineReasoningBuf) {
    const tail = inlineReasoningBuf.slice(-60);
    inlineReasoning.value =
      inlineReasoningBuf.length > 60 ? "…" + tail : tail;
  }
  if (final && inlineTimer) {
    clearTimeout(inlineTimer);
    inlineTimer = null;
  }
}

// 撤销：仅当插入区间内容与生成文本逐字一致时删除，用户手动改过则提示手动处理
function undoInline() {
  if (!inlineView || !inlineWroteText) return;
  const end = inlineStart + inlineWroteText.length;
  if (
    end > inlineView.state.doc.length ||
    inlineView.state.sliceDoc(inlineStart, end) !== inlineWroteText
  ) {
    message("续写内容已被手动修改，请手动删除该部分", { type: "warning" });
    resetInlineState();
    return;
  }
  inlineView.dispatch({
    changes: { from: inlineStart, to: end, insert: "" },
    selection: { anchor: inlineStart },
    scrollIntoView: true
  });
  resetInlineState();
}

function stopInline() {
  aiAbort?.();
  flushInline(true);
  inlineStreaming.value = false;
  inlineUndoable.value = inlineWroteText.length > 0;
  inlineDoneTip.value = "已停止";
}

// 失败重试：撤销已写部分后原样重发（光标回到原位，上下文重新采集）
function retryInline() {
  if (!aiLastRequest) return;
  if (inlineWroteText) {
    const end = inlineStart + inlineWroteText.length;
    if (
      end <= inlineView.state.doc.length &&
      inlineView.state.sliceDoc(inlineStart, end) === inlineWroteText
    ) {
      inlineView.dispatch({
        changes: { from: inlineStart, to: end, insert: "" },
        selection: { anchor: inlineStart }
      });
      inlinePos = inlineStart;
      inlineWroteText = "";
    }
  }
  inlinePending = "";
  inlineUndoable.value = false;
  inlineStreaming.value = true;
  aiStartedAt = Date.now();
  aiAbort = aiStream(aiLastRequest, {
    onDelta: text => {
      inlinePending += text;
      scheduleInlineFlush();
    },
    onReasoning: text => {
      inlineReasoningBuf += text;
      scheduleInlineFlush();
    },
    onDone: usage => {
      flushInline(true);
      inlineStreaming.value = false;
      inlineUndoable.value = inlineWroteText.length > 0;
      inlineDoneTip.value = `已续写 ${inlineWroteText.trim().length} 字${mkElapsedSuffix(usage)}`;
      refreshAIStatus();
    },
    onError: msg => {
      flushInline(true);
      inlineStreaming.value = false;
      if (inlineWroteText) inlineUndoable.value = true;
      message(msg, { type: "error" });
    }
  });
}

function resetInlineState() {
  inlineStreaming.value = false;
  inlineUndoable.value = false;
  inlineDoneTip.value = "";
  if (inlineTimer) {
    clearTimeout(inlineTimer);
    inlineTimer = null;
  }
}

function mkElapsedSuffix(usage?: any): string {
  const sec = ((Date.now() - aiStartedAt) / 1000).toFixed(1);
  const tokens = usage?.total_tokens;
  return `（${sec}s${tokens ? ` · ${tokens} tokens` : ""}）`;
}

function onAiBodyScroll() {
  const el = aiResultBodyRef.value;
  if (!el) return;
  aiStickBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 80;
}

function flushAIStream(final = false) {
  aiResultText.value = aiResultBuf;
  aiReasoningText.value = aiReasoningBuf;
  if (aiStickBottom && aiResultBodyRef.value) {
    aiResultBodyRef.value.scrollTop = aiResultBodyRef.value.scrollHeight;
  }
  // 全文校对任务：同步解析结构化建议
  if (aiLastRequest?.task === "proofread") {
    syncProofreadItems();
  }
  if (final && aiFlushTimer) {
    clearTimeout(aiFlushTimer);
    aiFlushTimer = null;
  }
}

function scheduleAIFlush() {
  if (aiFlushTimer) return;
  aiFlushTimer = setTimeout(() => {
    aiFlushTimer = null;
    flushAIStream();
  }, 150);
}

function startAI(
  req: AIEditRequest,
  label: string,
  apply: typeof aiApplyType.value,
  keepHistory = false
) {
  aiLastRequest = req;
  aiTaskLabel.value = label;
  aiApplyType.value = apply;
  aiResultBuf = "";
  aiReasoningBuf = "";
  aiResultText.value = "";
  aiReasoningText.value = "";
  aiStickBottom = true;
  aiError.value = "";
  aiMeta.value = null;
  if (!keepHistory) {
    aiHistory.value = [];
    aiHistoryIdx.value = -1;
  }
  aiStreaming.value = true;
  aiDrawerVisible.value = true;
  aiStartedAt = Date.now();
  aiAbort = aiStream(req, {
    onDelta: text => {
      aiResultBuf += text;
      scheduleAIFlush();
    },
    onReasoning: text => {
      aiReasoningBuf += text;
      scheduleAIFlush();
    },
    onDone: usage => {
      flushAIStream(true);
      aiStreaming.value = false;
      aiMeta.value = {
        elapsed: ((Date.now() - aiStartedAt) / 1000).toFixed(1) + "s",
        tokens: Number(usage?.total_tokens) || 0
      };
      aiHistory.value.push(aiResultBuf);
      aiHistoryIdx.value = aiHistory.value.length - 1;
      refreshAIStatus();
    },
    onError: msg => {
      flushAIStream(true);
      aiStreaming.value = false;
      aiError.value = msg;
      message(msg, { type: "error" });
    }
  });
}

function stopAI() {
  aiAbort?.();
  flushAIStream(true);
  aiStreaming.value = false;
}

function closeAiDrawer(done: () => void) {
  if (aiStreaming.value) stopAI();
  done();
}

function regenAI() {
  if (aiLastRequest) {
    const req = { ...aiLastRequest };
    const label = aiTaskLabel.value.replace(/ · 调整$/, "");
    const apply = aiApplyType.value;
    startAI(req, label, apply);
  }
}

// 抽屉内追加指令再生成：以上一版（当前展示的版本）为基础迭代
function refineAI() {
  const instruction = aiRefineInput.value.trim();
  const current = aiResultText.value;
  if (!instruction || !current || aiStreaming.value) return;
  if (aiQuotaExhausted.value) {
    message("今日 AI 调用配额已用完", { type: "warning" });
    return;
  }
  aiRefineInput.value = "";
  const base = aiLastRequest ?? ({} as AIEditRequest);
  startAI(
    {
      task: "refine",
      previous: current,
      instruction,
      title: base.title,
      series: base.series
    },
    aiTaskLabel.value.replace(/ · 调整$/, "") + " · 调整",
    aiApplyType.value,
    true
  );
}

// 版本栈导航：‹ › 回看历次结果，当前展示版本即应用/继续调整的对象
function aiNavHistory(dir: -1 | 1) {
  const next = aiHistoryIdx.value + dir;
  if (next < 0 || next >= aiHistory.value.length) return;
  aiHistoryIdx.value = next;
  aiResultText.value = aiHistory.value[next];
}

// 全文校对：解析结构化卡片，支持一键定位与采纳
function runProofread() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先填写正文内容", { type: "info" });
    return;
  }
  proofreadItems.value = [];
  proofreadFilter.value = "all";
  proofreadViewMode.value = "cards";
  startAI(
    { task: "proofread", digest: ruleForm.value.postContent },
    "AI 全文校对",
    "none"
  );
}

function syncProofreadItems() {
  const parsed = parseProofreadOutput(aiResultBuf);
  const statusMap = new Map<string, ProofreadItem["status"]>();
  proofreadItems.value.forEach(item => {
    statusMap.set(`${item.original}::${item.suggestion}`, item.status);
  });
  proofreadItems.value = parsed.map(item => {
    const key = `${item.original}::${item.suggestion}`;
    if (statusMap.has(key)) {
      item.status = statusMap.get(key)!;
    }
    return item;
  });
}

function getProofreadTagType(
  type: string
): "primary" | "danger" | "warning" | "info" | "success" {
  if (type.includes("错")) return "danger";
  if (type.includes("标点")) return "info";
  if (type.includes("语病")) return "warning";
  if (type.includes("格式")) return "primary";
  return "success";
}

function locateProofreadItem(item: ProofreadItem) {
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  const docText = view.state.doc.toString();
  const index = docText.indexOf(item.original);
  if (index === -1) {
    message(`在正文中未找到「${item.original}」，可能已被修改或删除`, {
      type: "warning"
    });
    item.status = "not_found";
    return;
  }
  const from = index;
  const to = index + item.original.length;
  view.dispatch({
    selection: { anchor: from, head: to },
    scrollIntoView: true
  });
  const lineNo = docText.slice(0, from).split("\n").length;
  message(`已在正文第 ${lineNo} 行定位并选中`, { type: "info" });
}

function applyProofreadItem(item: ProofreadItem) {
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  const docText = view.state.doc.toString();
  const index = docText.indexOf(item.original);
  if (index === -1) {
    message(`在正文中未找到「${item.original}」，无法采纳`, {
      type: "warning"
    });
    item.status = "not_found";
    return;
  }
  const from = index;
  const to = index + item.original.length;
  view.dispatch({
    changes: { from, to, insert: item.suggestion },
    selection: { anchor: from + item.suggestion.length },
    scrollIntoView: true
  });
  ruleForm.value.postContent = view.state.doc.toString();
  item.status = "applied";
  message(`已采纳修改：「${item.original}」→「${item.suggestion}」`, {
    type: "success"
  });
}

function ignoreProofreadItem(item: ProofreadItem) {
  item.status = "ignored";
}

function applyAllProofreadItems() {
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  const pending = proofreadItems.value.filter(i => i.status === "pending");
  if (pending.length === 0) return;

  const docText = view.state.doc.toString();
  const replacements: {
    item: ProofreadItem;
    from: number;
    to: number;
    insert: string;
  }[] = [];

  for (const item of pending) {
    const idx = docText.indexOf(item.original);
    if (idx !== -1) {
      replacements.push({
        item,
        from: idx,
        to: idx + item.original.length,
        insert: item.suggestion
      });
    } else {
      item.status = "not_found";
    }
  }

  if (replacements.length === 0) {
    message("未在正文中匹配到待修改内容", { type: "warning" });
    return;
  }

  // 按起始位置倒序排列，保证区间替换不引起前方位置漂移
  replacements.sort((a, b) => b.from - a.from);

  // 丢弃重叠冲突项
  const safeReplacements: typeof replacements = [];
  let lastFrom = Infinity;
  for (const rep of replacements) {
    if (rep.to <= lastFrom) {
      safeReplacements.push(rep);
      lastFrom = rep.from;
    }
  }

  const changes = safeReplacements.map(r => ({
    from: r.from,
    to: r.to,
    insert: r.insert
  }));

  view.dispatch({
    changes,
    scrollIntoView: true
  });

  ruleForm.value.postContent = view.state.doc.toString();
  safeReplacements.forEach(r => {
    r.item.status = "applied";
  });

  message(`已一键采纳 ${safeReplacements.length} 处校对建议`, {
    type: "success"
  });
}

function applyAIResult() {
  if (!aiResultText.value) return;
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  // 目前仅替换选区一种写回（续写已改为光标流式插入，不经过抽屉）
  view.dispatch({
    changes: { from: aiSelFrom.value, to: aiSelTo.value, insert: aiResultText.value }
  });
  aiDrawerVisible.value = false;
  message("已写回正文", { type: "success" });
}

// 在选区后插入生成内容（保留原选区作为参考）
function insertBelowAIResult() {
  if (!aiResultText.value) return;
  const view = editorRef.value?.getEditorView?.();
  if (!view) return;
  const insertPos = Math.max(aiSelFrom.value, aiSelTo.value);
  const chunk = `\n\n${aiResultText.value}\n`;
  view.dispatch({
    changes: { from: insertPos, to: insertPos, insert: chunk },
    selection: { anchor: insertPos + chunk.length },
    scrollIntoView: true
  });
  aiDrawerVisible.value = false;
  message("已在选区后插入内容", { type: "success" });
}

async function copyAIResult() {
  try {
    await navigator.clipboard.writeText(aiResultText.value);
    message("已复制", { type: "success" });
  } catch {
    message("复制失败", { type: "error" });
  }
}

// 生成标题：流式接收整段文本后按行拆成候选
function generateTitle() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先填写正文内容", { type: "info" });
    return;
  }
  aiTitleLoading.value = true;
  aiTitleCandidates.value = [];
  let acc = "";
  const abort = aiStream(
    { task: "title", digest: ruleForm.value.postContent },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        aiTitleLoading.value = false;
        aiTitleCandidates.value = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 3);
      },
      onError: msg => {
        aiTitleLoading.value = false;
        message(msg, { type: "error" });
      }
    }
  );
  aiAbort = abort;
}

function applyTitle(t: string) {
  ruleForm.value.title = t;
  titlePopoverVisible.value = false;
}

// 生成标签建议：按行拆候选，点击追加进已选标签
function generateTags() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先填写正文内容", { type: "info" });
    return;
  }
  if (aiQuotaExhausted.value) {
    message("今日 AI 调用配额已用完", { type: "warning" });
    return;
  }
  aiTagsLoading.value = true;
  aiTagCandidates.value = [];
  let acc = "";
  const abort = aiStream(
    {
      task: "tags",
      digest: ruleForm.value.postContent,
      title: ruleForm.value.title
    },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        aiTagsLoading.value = false;
        aiTagCandidates.value = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 8);
        refreshAIStatus();
      },
      onError: msg => {
        aiTagsLoading.value = false;
        message(msg, { type: "error" });
      }
    }
  );
  aiAbort = abort;
}

function applyTag(t: string) {
  if (!ruleForm.value.tags.includes(t)) {
    ruleForm.value.tags.push(t);
    message(`已添加标签「${t}」`, { type: "success" });
  }
}

function generateSummary() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先填写正文内容", { type: "info" });
    return;
  }
  aiSummaryLoading.value = true;
  let acc = "";
  const abort = aiStream(
    { task: "summary", digest: ruleForm.value.postContent },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        aiSummaryLoading.value = false;
        if (acc.trim()) ruleForm.value.summary = acc.trim();
      },
      onError: msg => {
        aiSummaryLoading.value = false;
        message(msg, { type: "error" });
      }
    }
  );
  aiAbort = abort;
}

// ---------- AI 发文助手 (Publish Copilot) ----------
function openPublishCopilot() {
  if (!ruleForm.value.postContent.trim()) {
    message("请先填写正文内容，以便 AI 分析提炼", { type: "info" });
    return;
  }
  if (aiQuotaExhausted.value) {
    message("今日 AI 调用配额已用完", { type: "warning" });
    return;
  }
  copilotVisible.value = true;
  copilotChosenTitle.value = ruleForm.value.title || "";
  copilotSummary.value = ruleForm.value.summary || "";
  copilotSelectedTags.value = [...ruleForm.value.tags];
  copilotMatchedTags.value = [];
  copilotNewTags.value = [];
  copilotTitleCandidates.value = [];

  // 并发触发三大生成任务
  copilotRegenTitle();
  copilotRegenSummary();
  copilotRegenTags();
}

function closeCopilot(done: () => void) {
  copilotAbortTitle?.();
  copilotAbortSummary?.();
  copilotAbortTags?.();
  copilotLoading.value = false;
  done();
}

function copilotRegenTitle() {
  copilotTitleLoading.value = true;
  copilotTitleCandidates.value = [];
  let acc = "";
  copilotAbortTitle = aiStream(
    { task: "title", digest: ruleForm.value.postContent },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        copilotTitleLoading.value = false;
        const candidates = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 4);
        copilotTitleCandidates.value = candidates;
        if (!copilotChosenTitle.value && candidates.length > 0) {
          copilotChosenTitle.value = candidates[0];
        }
      },
      onError: msg => {
        copilotTitleLoading.value = false;
        console.warn("AI 标题生成失败:", msg);
      }
    }
  );
}

function copilotRegenSummary() {
  copilotSummaryLoading.value = true;
  let acc = "";
  copilotAbortSummary = aiStream(
    { task: "summary", digest: ruleForm.value.postContent },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        copilotSummaryLoading.value = false;
        if (acc.trim()) {
          copilotSummary.value = acc.trim();
        }
      },
      onError: msg => {
        copilotSummaryLoading.value = false;
        console.warn("AI 摘要生成失败:", msg);
      }
    }
  );
}

function copilotRegenTags() {
  copilotTagsLoading.value = true;
  copilotMatchedTags.value = [];
  copilotNewTags.value = [];
  let acc = "";
  copilotAbortTags = aiStream(
    {
      task: "tags",
      digest: ruleForm.value.postContent,
      title: copilotChosenTitle.value || ruleForm.value.title
    },
    {
      onDelta: t => (acc += t),
      onDone: () => {
        copilotTagsLoading.value = false;
        const recTags = acc
          .split("\n")
          .map(x => x.replace(/^[-*\d.、\s]+/, "").trim())
          .filter(Boolean)
          .slice(0, 8);

        // 系统已有标签名列表
        const systemTagNames: string[] = tags.value.map(
          (t: any) => t.tagName || ""
        );

        const matched: string[] = [];
        const newTags: string[] = [];

        recTags.forEach(tag => {
          // 不区分大小写匹配已有标签
          const found = systemTagNames.find(
            st => st.toLowerCase() === tag.toLowerCase()
          );
          if (found) {
            if (!matched.includes(found)) matched.push(found);
            // 命中已有标签，自动勾选
            if (!copilotSelectedTags.value.includes(found)) {
              copilotSelectedTags.value.push(found);
            }
          } else {
            if (!newTags.includes(tag)) newTags.push(tag);
          }
        });

        copilotMatchedTags.value = matched;
        copilotNewTags.value = newTags;
        refreshAIStatus();
      },
      onError: msg => {
        copilotTagsLoading.value = false;
        console.warn("AI 标签生成失败:", msg);
      }
    }
  );
}

function toggleCopilotTag(t: string) {
  const idx = copilotSelectedTags.value.indexOf(t);
  if (idx > -1) {
    copilotSelectedTags.value.splice(idx, 1);
  } else {
    copilotSelectedTags.value.push(t);
  }
}

function applyCopilotToForm() {
  let appliedCount = 0;
  if (copilotChosenTitle.value && copilotChosenTitle.value !== ruleForm.value.title) {
    ruleForm.value.title = copilotChosenTitle.value;
    appliedCount++;
  }
  if (copilotSummary.value && copilotSummary.value !== ruleForm.value.summary) {
    ruleForm.value.summary = copilotSummary.value;
    appliedCount++;
  }
  if (copilotSelectedTags.value.length > 0) {
    // 合并标签（保持去重）
    const merged = Array.from(
      new Set([...ruleForm.value.tags, ...copilotSelectedTags.value])
    );
    ruleForm.value.tags = merged;
    appliedCount++;
  }
  copilotVisible.value = false;
  message(`已将发文建议应用到文章（更新了 ${appliedCount} 项）`, {
    type: "success"
  });
}

// 选中文本后浮现的快速工具栏：加粗/斜体 + AI 菜单（索引 0 = defToolbars
// 里的 DropdownToolbar，选中即弹出完整 AI 菜单，不再只是一键润色）
const floatingToolbars: ToolbarNames[] = ["bold", "italic", 0];

// Markdown 编辑器工具栏配置（支持分屏实时预览、全屏、HTML预览、大纲目录）；
// 数字 0 对应 defToolbars 插槽里的第一个组件（AI 下拉菜单）
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
  0,
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

// Esc 键统一调度：续写中优先停止续写，其次退出专注模式
function handleEscKey() {
  if (inlineStreaming.value) {
    stopInline();
    return;
  }
  if (focusMode.value) {
    exitFocusMode();
  }
}

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
          ruleForm.value = { series: "", ...res.payload };
          // 已是未来时间的定时文章：回填开关与时间
          if (
            res.payload.pubTime &&
            Number(res.payload.pubTime) > Date.now() / 1000
          ) {
            scheduleEnabled.value = true;
            scheduleTime.value = res.payload.pubTime;
          }
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
        // 定时发布：带上未来 pubTime 并直接置为已发布；未开启则不带这两项，
        // 由后端保持原有发布状态与发布时间
        const payload: Record<string, any> = { ...ruleForm.value };
        delete payload.pubTime;
        delete payload.published;
        if (scheduleEnabled.value && scheduleTime.value) {
          payload.pubTime = Number(scheduleTime.value);
          payload.published = 1;
        }
        let res;
        if (isEdit.value) {
          res = await updatePost(payload);
        } else {
          res = await savePost(payload);
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

/* ===== AI 辅助写作 ===== */
.ai-menu-overlay {
  margin: 0;
  padding: 6px;
  list-style: none;
  background: var(--el-bg-color, #fff);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.14);
  min-width: 150px;
}

.ai-menu-overlay li {
  padding: 7px 12px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
  color: var(--el-text-color-primary, #303133);
}

.ai-menu-models {
  display: block !important;
  padding: 4px 12px !important;
  cursor: default !important;
}

.ai-model-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 6px;
}

.ai-model-chip {
  padding: 3px 10px;
  border-radius: 10px;
  font-size: 12px;
  border: 1px solid var(--el-border-color, #dcdfe6);
  color: var(--el-text-color-regular, #606266);
  cursor: pointer;
  user-select: none;
  transition: all 0.15s;
}

.ai-model-chip:hover {
  border-color: var(--el-color-primary, #409eff);
}

.ai-model-chip.active {
  background: var(--el-color-primary, #409eff);
  border-color: var(--el-color-primary, #409eff);
  color: #fff;
}

.ai-menu-action {
  color: var(--el-color-primary, #409eff) !important;
  font-size: 12px !important;
}

.ai-tpl-default {
  margin-top: 12px;
  padding: 10px 12px;
  border-radius: 6px;
  background: var(--el-fill-color-light, #f5f7fa);
}

.ai-tpl-default-title {
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  margin-bottom: 4px;
}

.ai-tpl-default-text {
  font-size: 13px;
  line-height: 1.6;
  color: var(--el-text-color-regular, #606266);
}

.ai-menu-overlay li:hover {
  background: var(--el-fill-color-light, #f5f7fa);
  color: var(--el-color-primary, #409eff);
}

.ai-menu-overlay li.disabled {
  color: var(--el-text-color-placeholder, #a8abb2);
  cursor: not-allowed;
}

.ai-menu-overlay li.disabled:hover {
  background: transparent;
  color: var(--el-text-color-placeholder, #a8abb2);
}

.ai-menu-overlay li.ai-menu-quota {
  margin-top: 4px;
  padding-top: 6px;
  border-top: 1px solid var(--el-border-color-lighter, #ebeef5);
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  cursor: default;
}

.ai-menu-overlay li.ai-menu-quota:hover {
  background: transparent;
  color: var(--el-text-color-secondary, #909399);
}

/* 续写光标流式插入状态条优化：黏性吸顶保证长文滚动时不脱离视线 */
.ai-inline-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  margin-bottom: 6px;
  padding: 8px 12px;
  border: 1px solid var(--el-color-primary-light-7, #d9ecff);
  border-radius: 6px;
  background: var(--el-color-primary-light-9, #ecf5ff);
  font-size: 13px;
  color: var(--el-text-color-regular, #606266);
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.12);
}

.ai-inline-bar .ai-inline-dot {
  flex: none;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--el-color-primary, #409eff);
  animation: ai-reasoning-pulse 1.1s ease-in-out infinite;
}

.ai-inline-bar .ai-inline-label {
  flex: none;
}

.ai-inline-bar .ai-inline-reasoning {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
}

.ai-inline-bar .ai-inline-actions {
  display: inline-flex;
  gap: 6px;
  margin-left: auto;
}

/* 抽屉 footer：追加指令行 + 元信息行 */
.ai-result-foot {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.ai-refine-row {
  display: flex;
  gap: 8px;
}

.ai-result-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
}

.ai-result-meta .ai-version-nav {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.ai-result-meta .ai-quota-text {
  margin-left: auto;
}

.ai-error-banner {
  margin-bottom: 12px;
}

.ai-toolbar-trigger {
  font-size: 13px;
  line-height: 1;
  padding: 0 2px;
  user-select: none;
}

.ai-result-body {
  height: 100%;
  overflow-y: auto;
}

.ai-reasoning {
  margin-bottom: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  background: var(--el-fill-color-light, #f5f7fa);
}

.ai-reasoning-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  margin-bottom: 6px;
}

.ai-reasoning-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--el-color-primary, #409eff);
}

.ai-reasoning-dot.pulse {
  animation: ai-reasoning-pulse 1.1s ease-in-out infinite;
}

@keyframes ai-reasoning-pulse {
  0%,
  100% {
    opacity: 0.35;
  }
  50% {
    opacity: 1;
  }
}

.ai-reasoning-text {
  font-size: 12px;
  line-height: 1.7;
  color: var(--el-text-color-secondary, #909399);
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 180px;
  overflow-y: auto;
}

.ai-result-preview {
  min-height: 200px;
}

.ai-result-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.ai-title-pop .ai-title-item {
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  line-height: 1.5;
  transition: background-color 0.15s;
}

.ai-title-pop .ai-title-item:hover {
  background-color: var(--el-fill-color-light);
}

.ai-title-pop .ai-title-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  padding: 4px 2px;
}

/* 差异对比视图 */
.ai-view-switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: var(--el-fill-color-lighter, #fafafa);
  border: 1px solid var(--el-border-color-lighter, #ebeef5);
  border-radius: 6px;
}

.diff-legend {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.legend-badge {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
}

.legend-del {
  background: #ffeef0;
  color: #cf222e;
  border: 1px solid #ffd8d3;
}

.legend-ins {
  background: #dafbe1;
  color: #1a7f37;
  border: 1px solid #aceebb;
}

.ai-diff-container {
  padding: 12px;
  background: var(--el-bg-color, #fff);
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 6px;
  min-height: 200px;
  max-height: calc(100vh - 350px);
  overflow-y: auto;
  line-height: 1.8;
  font-size: 14px;
}

.ai-diff-placeholder {
  color: var(--el-text-color-secondary, #909399);
  text-align: center;
  padding: 40px 0;
  font-size: 13px;
}

.diff-chunk {
  display: inline;
  word-break: break-word;
  white-space: pre-wrap;
}

.diff-del {
  background-color: #ffeef0;
  color: #cf222e;
  text-decoration: line-through;
  padding: 1px 3px;
  border-radius: 3px;
  margin: 0 1px;
}

.diff-ins {
  background-color: #dafbe1;
  color: #1a7f37;
  text-decoration: none;
  font-weight: 500;
  padding: 1px 3px;
  border-radius: 3px;
  margin: 0 1px;
}

.diff-common {
  color: var(--el-text-color-primary, #303133);
}

/* 选区自定义指令（Ask AI）菜单项 */
.ai-menu-custom-box {
  padding: 8px 10px !important;
  cursor: default !important;
  border-bottom: 1px solid var(--el-border-color-lighter, #ebeef5);
  margin-bottom: 4px;
}

.ai-custom-prompt-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 240px;
}

.ai-send-icon {
  cursor: pointer;
  color: var(--el-text-color-placeholder, #a8abb2);
  transition: color 0.2s;
}

.ai-send-icon.active {
  color: var(--el-color-primary, #409eff);
}

.ai-send-icon:hover {
  color: var(--el-color-primary, #409eff);
}

.ai-quick-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.ai-quick-tag {
  font-size: 11px;
  padding: 1px 6px;
  background: var(--el-fill-color-light, #f5f7fa);
  color: var(--el-text-color-regular, #606266);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s;
}

.ai-quick-tag:hover {
  background: var(--el-color-primary-light-9, #ecf5ff);
  color: var(--el-color-primary, #409eff);
}

/* AI 发文助手弹窗样式 */
:deep(.ai-copilot-dialog) {
  .copilot-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
    max-height: 60vh;
    overflow-y: auto;
    padding-right: 4px;
  }

  .copilot-section {
    border: 1px solid var(--el-border-color-light, #e4e7ed);
    border-radius: 8px;
    padding: 12px 14px;
    background: var(--el-bg-color, #fff);
  }

  .copilot-section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .copilot-section-title {
    font-size: 14px;
    font-weight: 600;
    color: var(--el-text-color-primary, #303133);
    display: flex;
    align-items: center;
  }

  .copilot-radio-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 100%;
  }

  .copilot-title-radio {
    margin-right: 0;
    height: auto;
    padding: 8px 10px;
    border-radius: 6px;
    border: 1px solid var(--el-border-color-lighter, #ebeef5);
    white-space: normal;
    line-height: 1.4;
    transition: all 0.2s;

    &:hover {
      border-color: var(--el-color-primary-light-5, #a0cfff);
      background: var(--el-color-primary-light-9, #ecf5ff);
    }

    &.is-checked {
      border-color: var(--el-color-primary, #409eff);
      background: var(--el-color-primary-light-9, #ecf5ff);
    }
  }

  .copilot-tag-subgroup {
    margin-bottom: 8px;
    &:last-child {
      margin-bottom: 0;
    }
  }

  .subgroup-label {
    font-size: 12px;
    color: var(--el-text-color-secondary, #909399);
    margin-bottom: 6px;
  }

  .copilot-tags-wrap {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .copilot-tag-item {
    font-size: 13px;
    cursor: pointer;
    user-select: none;
    transition: all 0.15s;

    &.matched {
      font-weight: 500;
    }

    &.new-tag {
      border-style: dashed;
    }
  }

  .copilot-foot {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }
}

/* 全文校对交互卡片界面 */
.proofread-header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: var(--el-fill-color-lighter, #fafafa);
  border: 1px solid var(--el-border-color-lighter, #ebeef5);
  border-radius: 6px;
  gap: 8px;
}

.proofread-stats {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary, #303133);
}

.proofread-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.proofread-filter-row {
  margin-bottom: 10px;
}

.proofread-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 10px;
  color: var(--el-text-color-secondary, #909399);
  font-size: 13px;
  text-align: center;
}

.proofread-cards-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.proofread-card {
  padding: 10px 12px;
  border-radius: 6px;
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  background: var(--el-bg-color, #fff);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.2s;

  &:hover {
    border-color: var(--el-color-primary-light-5, #a0cfff);
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
  }

  &.is-applied {
    background: var(--el-fill-color-lighter, #fafafa);
    border-color: var(--el-border-color-lighter, #ebeef5);
    opacity: 0.8;
  }

  &.is-ignored {
    opacity: 0.6;
    background: var(--el-fill-color-lighter, #fafafa);
  }

  &.is-not-found {
    border-color: var(--el-color-warning-light-5, #f8e3c5);
  }
}

.proofread-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 12px;
}

.proofread-reason {
  color: var(--el-text-color-secondary, #909399);
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.proofread-status-indicator {
  margin-left: auto;
  font-size: 12px;
}

.proofread-card-diff {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px 10px;
  border-radius: 4px;
  background: var(--el-fill-color-light, #f5f7fa);
  margin-bottom: 8px;
  font-size: 13px;
  line-height: 1.6;
}

.diff-line {
  display: flex;
  align-items: flex-start;
  gap: 6px;
}

.diff-tag {
  flex: none;
  font-size: 11px;
  padding: 1px 4px;
  border-radius: 3px;
  line-height: 1.2;

  &.del {
    background: #ffeef0;
    color: #cf222e;
    border: 1px solid #ffd8d3;
  }

  &.ins {
    background: #dafbe1;
    color: #1a7f37;
    border: 1px solid #aceebb;
  }
}

.diff-text {
  word-break: break-word;
  white-space: pre-wrap;

  &.del {
    color: #cf222e;
    text-decoration: line-through;
  }

  &.ins {
    color: #1a7f37;
    font-weight: 500;
  }
}

.proofread-card-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}
</style>
