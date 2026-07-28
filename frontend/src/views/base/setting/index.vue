<template>
  <div class="setting-container">
    <el-card class="setting-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">网站设置</span>
          <el-button
            type="primary"
            :loading="saveLoading"
            @click="submitWebsiteForm()"
          >
            保存设置
          </el-button>
        </div>
      </template>

      <el-tabs
        v-model="activeName"
        type="card"
        class="setting-tabs"
        @tab-click="handleClick"
      >
        <el-tab-pane name="SEO">
          <template #label>
            <div class="tab-label">
              <el-icon><Notification /></el-icon>
              <span>SEO设置</span>
            </div>
          </template>

          <el-form
            ref="seoFormRef"
            label-position="top"
            :model="WebsiteForm"
            class="setting-form"
            status-icon
          >
            <el-row :gutter="20">
              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="网站标题" prop="title">
                  <el-input
                    v-model="WebsiteForm.title"
                    placeholder="请输入网站标题"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="填写网站的标题，显示在浏览器标签页上"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">
                    显示在浏览器标签页和搜索引擎结果中的标题
                  </div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="网站域名" prop="url">
                  <el-input
                    v-model="WebsiteForm.url"
                    placeholder="例如: https://www.example.com"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="填写网站的完整域名，包含协议部分"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">完整的网站访问地址</div>
                </el-form-item>
              </el-col>

              <el-col :span="24">
                <el-form-item label="关键词" prop="keywords">
                  <el-input
                    v-model="WebsiteForm.keywords"
                    placeholder="请输入网站关键词，多个关键词用逗号分隔"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="填写网站的关键词，便于搜索引擎收录"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">
                    多个关键词请用英文逗号分隔，如：Java开发,后端实战,SpringBoot
                  </div>
                </el-form-item>
              </el-col>

              <el-col :span="24">
                <el-form-item label="页面描述" prop="description">
                  <el-input
                    v-model="WebsiteForm.description"
                    type="textarea"
                    :rows="3"
                    placeholder="请输入网站描述信息"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="网站的描述信息，显示在搜索引擎结果中"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">简洁明了地描述网站的主要内容和服务</div>
                </el-form-item>
              </el-col>

              <el-col :span="24">
                <el-form-item label="页脚版权" prop="copyright">
                  <el-input
                    v-model="WebsiteForm.copyright"
                    placeholder="例如: © 2023 Your Company. All rights reserved."
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="网站底部的版权声明信息"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">显示在网站底部的版权信息</div>
                </el-form-item>
              </el-col>

              <el-col :span="24">
                <el-form-item label="网站公告" prop="notice">
                  <el-input
                    v-model="WebsiteForm.notice"
                    type="textarea"
                    :rows="3"
                    placeholder="请输入网站公告内容"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="网站首页显示的重要通知信息"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">首页顶部滚动显示的重要通知</div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="record">
          <template #label>
            <div class="tab-label">
              <el-icon><Document /></el-icon>
              <span>备案统计</span>
            </div>
          </template>

          <el-form
            label-position="top"
            :model="WebsiteForm"
            class="setting-form"
            status-icon
          >
            <el-row :gutter="20">
              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="ICP备案号" prop="icp">
                  <el-input
                    v-model="WebsiteForm.icp"
                    placeholder="例如: 黔ICP备19012300号"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="工信部备案管理系统中的备案号"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">在工信部备案管理系统中获得的备案号</div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="百度收录校验码" prop="site">
                  <el-input
                    v-model="WebsiteForm.site"
                    placeholder="例如: IIO8ueBTje"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="百度站长平台的站点验证标识"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">用于百度搜索资源平台验证网站所有权</div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="百度统计Key" prop="stat">
                  <el-input
                    v-model="WebsiteForm.stat"
                    placeholder="例如: 0fef4e6e7db66c9ddcf3cb59160b2ba7"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="百度统计的站点识别码"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">百度统计用于识别网站的唯一标识</div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="social">
          <template #label>
            <div class="tab-label">
              <el-icon><Link /></el-icon>
              <span>社交链接</span>
            </div>
          </template>

          <el-form
            label-position="top"
            :model="WebsiteForm"
            class="setting-form"
            status-icon
          >
            <el-row :gutter="20">
              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="GitHub地址" prop="github">
                  <el-input
                    v-model="WebsiteForm.github"
                    placeholder="例如: https://github.com/username"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="项目的GitHub仓库地址"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">开源项目托管地址</div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="Gitee地址" prop="gitee">
                  <el-input
                    v-model="WebsiteForm.gitee"
                    placeholder="例如: https://gitee.com/username"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip content="项目的码云仓库地址" placement="top">
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">国内代码托管平台地址</div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="网站创建日期" prop="siteStartDate">
                  <el-date-picker
                    v-model="WebsiteForm.siteStartDate"
                    type="date"
                    placeholder="选择网站创建日期"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                  <div class="form-tip">
                    用于计算网站运行天数，显示在侧边栏网站资讯中
                  </div>
                </el-form-item>
              </el-col>

              <el-col :xs="24" :sm="24" :md="12" :lg="12">
                <el-form-item label="联系邮箱" prop="email">
                  <el-input
                    v-model="WebsiteForm.email"
                    placeholder="例如: example@outlook.com"
                    clearable
                  >
                    <template #suffix>
                      <el-tooltip
                        content="网站联系邮箱，显示在首页和侧边栏"
                        placement="top"
                      >
                        <el-icon class="help-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </template>
                  </el-input>
                  <div class="form-tip">用于访客联系的邮箱地址</div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {
  Notification,
  Document,
  Link,
  InfoFilled
} from "@element-plus/icons-vue";
import { onMounted, ref } from "vue";
import { getWebsite, updateWebsite } from "@/api/setting";
import type { TabsPaneContext, FormInstance } from "element-plus";

import { message } from "@/utils/message";

const activeName = ref("SEO");
const saveLoading = ref(false);
const seoFormRef = ref<FormInstance>();

const WebsiteForm = ref({
  id: undefined,
  title: "",
  url: "",
  notice: "",
  description: "",
  stat: "",
  site: "",
  icp: "",
  copyright: "",
  keywords: "",
  github: "",
  gitee: "",
  siteStartDate: "",
  email: ""
});

const handleClick = (tab: TabsPaneContext) => {
  if (tab.paneName === "SEO") {
    getWebsite().then(res => {
      if (res.code === 200) {
        WebsiteForm.value = res.payload;
      }
    });
  }
};

const submitWebsiteForm = async () => {
  if (!WebsiteForm.value.id) {
    message("正在加载数据，请稍后再试", { type: "warning" });
    return;
  }

  saveLoading.value = true;
  try {
    // 创建一个新对象，只包含需要提交的字段
    const formData = {
      id: WebsiteForm.value.id,
      title: WebsiteForm.value.title,
      url: WebsiteForm.value.url,
      notice: WebsiteForm.value.notice,
      description: WebsiteForm.value.description,
      stat: WebsiteForm.value.stat,
      site: WebsiteForm.value.site,
      icp: WebsiteForm.value.icp,
      copyright: WebsiteForm.value.copyright,
      keywords: WebsiteForm.value.keywords,
      github: WebsiteForm.value.github,
      gitee: WebsiteForm.value.gitee,
      siteStartDate: WebsiteForm.value.siteStartDate,
      email: WebsiteForm.value.email
    };

    const res = await updateWebsite(formData);
    if (res.code === 200) {
      message("设置保存成功", { type: "success" });
    } else {
      message(`保存失败: ${res.message || "未知错误"}`, { type: "error" });
    }
  } catch (error: any) {
    message(`保存过程中发生错误: ${error.message || "未知错误"}`, {
      type: "error"
    });
  } finally {
    saveLoading.value = false;
  }
};

onMounted(() => {
  getWebsite().then(res => {
    if (res.code === 200) {
      WebsiteForm.value = res.payload;
    }
  });
});
</script>

<style scoped lang="scss">
.setting-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);

  .setting-card {
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

  .setting-tabs {
    :deep(.el-tabs__header) {
      margin-bottom: 25px;

      .el-tabs__nav {
        border: none;
      }

      .el-tabs__item {
        height: 50px;
        line-height: 50px;
        font-size: 15px;
        font-weight: 500;
        color: var(--el-text-color-secondary);
        border: none !important;
        border-radius: 6px 6px 0 0;
        transition: all 0.3s;

        &.is-active {
          color: var(--el-color-primary);
          background-color: var(--el-color-primary-light-9);
          font-weight: 600;
        }

        &:hover {
          color: var(--el-color-primary);
        }

        .tab-label {
          display: flex;
          align-items: center;
          gap: 8px;
        }
      }
    }
  }

  .setting-form {
    :deep(.el-form-item) {
      margin-bottom: 24px;

      .el-form-item__label {
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin-bottom: 8px;
      }

      .el-input,
      .el-textarea {
        .help-icon {
          cursor: pointer;
          color: var(--el-color-info);
          font-size: 16px;

          &:hover {
            color: var(--el-color-primary);
          }
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
  .setting-container {
    padding: 12px;

    :deep(.el-card__header) {
      padding: 15px;
    }
  }
}
</style>
