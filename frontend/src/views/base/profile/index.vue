<template>
  <div class="profile-container">
    <el-row :gutter="20" class="profile-layout">
      <!-- 左侧：用户身份画像卡片 -->
      <el-col :xs="24" :sm="24" :md="9" :lg="8" class="left-col">
        <el-card class="user-profile-card" :body-style="{ padding: '0px' }" shadow="hover">
          <!-- 渐变装饰顶图 -->
          <div class="card-cover-banner">
            <div class="banner-mesh-glow" />
            <div class="role-tag-float">
              <el-tag effect="dark" round class="role-pill">
                <span class="role-badge-text">👑 {{ userInfo.role || "超级管理员" }}</span>
              </el-tag>
            </div>
          </div>

          <!-- 用户头像与基本概览 -->
          <div class="user-meta-body">
            <div class="avatar-wrapper">
              <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :http-request="handleCustomUpload"
                :before-upload="beforeAvatarUpload"
              >
                <div class="avatar-box" :class="{ 'is-uploading': uploadLoading }">
                  <el-image
                    v-if="userInfo.avatar"
                    :src="userInfo.avatar"
                    fit="cover"
                    class="user-avatar"
                  >
                    <template #error>
                      <img :src="defaultAvatar" class="user-avatar" />
                    </template>
                  </el-image>
                  <img v-else :src="defaultAvatar" class="user-avatar" />

                  <!-- 上传浮层 -->
                  <div class="avatar-mask">
                    <el-icon v-if="!uploadLoading" class="mask-icon"><Camera /></el-icon>
                    <el-icon v-else class="mask-icon is-loading"><Refresh /></el-icon>
                    <span class="mask-text">{{ uploadLoading ? "上传中..." : "更换头像" }}</span>
                  </div>
                </div>
              </el-upload>
              <!-- 在线状态指示点 -->
              <span class="status-indicator-dot" title="当前在线">
                <span class="ping-circle" />
                <span class="dot-core" />
              </span>
            </div>

            <!-- 用户名与昵称 -->
            <div class="name-section">
              <h2 class="user-nickname">{{ userInfo.nickName || userInfo.username || "管理员" }}</h2>
              <div class="user-handle-wrap">
                <span class="user-handle">@{{ userInfo.username }}</span>
                <span class="user-uid">UID: {{ userInfo.userId || 1 }}</span>
              </div>
            </div>

            <!-- 励志签名 / 问候语 -->
            <div class="bio-quote-box">
              <span class="quote-symbol">“</span>
              <span class="quote-text">{{ currentTimeText }}，{{ userInfo.nickName || userInfo.username }}。生活变得再糟糕，也不妨碍我变得更好！</span>
              <span class="quote-symbol">”</span>
            </div>

            <!-- 快速统计卡片 -->
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-num">{{ userInfo.loginCount }}</div>
                <div class="stat-label">
                  <el-icon class="mr-1"><Monitor /></el-icon> 累计登录
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-num">{{ formatRegisterYear(userInfo.registerTime) }}</div>
                <div class="stat-label">
                  <el-icon class="mr-1"><Calendar /></el-icon> 加入年份
                </div>
              </div>
            </div>

            <el-divider class="card-divider" />

            <!-- 详细身份清单 -->
            <div class="detail-info-list">
              <div class="detail-item">
                <div class="item-left">
                  <span class="item-icon-box primary">
                    <el-icon><Message /></el-icon>
                  </span>
                  <span class="item-title">电子邮箱</span>
                </div>
                <div class="item-value">
                  {{ userInfo.email || "未绑定邮箱" }}
                </div>
              </div>

              <div class="detail-item">
                <div class="item-left">
                  <span class="item-icon-box success">
                    <el-icon><Phone /></el-icon>
                  </span>
                  <span class="item-title">联系手机</span>
                </div>
                <div class="item-value">
                  {{ userInfo.phonenumber || "未绑定手机" }}
                </div>
              </div>

              <div class="detail-item">
                <div class="item-left">
                  <span class="item-icon-box warning">
                    <el-icon><User /></el-icon>
                  </span>
                  <span class="item-title">性别</span>
                </div>
                <div class="item-value">
                  <el-tag v-if="userInfo.sex === 1" size="small" type="primary" effect="plain">男</el-tag>
                  <el-tag v-else-if="userInfo.sex === 2" size="small" type="danger" effect="plain">女</el-tag>
                  <el-tag v-else size="small" type="info" effect="plain">保密</el-tag>
                </div>
              </div>

              <div class="detail-item">
                <div class="item-left">
                  <span class="item-icon-box info">
                    <el-icon><Location /></el-icon>
                  </span>
                  <span class="item-title">最近登录 IP</span>
                </div>
                <div class="item-value ip-copy-val">
                  <el-tag size="small" class="ip-tag">{{ userInfo.loginIp || "未知" }}</el-tag>
                  <el-tooltip content="复制 IP 地址" placement="top">
                    <el-button
                      v-if="userInfo.loginIp && userInfo.loginIp !== '未知'"
                      link
                      type="primary"
                      :icon="DocumentCopy"
                      class="copy-btn"
                      @click="handleCopyIp"
                    />
                  </el-tooltip>
                </div>
              </div>

              <div class="detail-item">
                <div class="item-left">
                  <span class="item-icon-box purple">
                    <el-icon><Clock /></el-icon>
                  </span>
                  <span class="item-title">最近登录时间</span>
                </div>
                <div class="item-value time-val">
                  {{ userInfo.loginTime || "-" }}
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：Tab 选项卡式配置工作台 -->
      <el-col :xs="24" :sm="24" :md="15" :lg="16" class="right-col">
        <el-card class="settings-card" shadow="hover">
          <el-tabs v-model="activeTab" class="profile-tabs">
            <!-- Tab 1: 个人基本资料修改 -->
            <el-tab-pane name="basic">
              <template #label>
                <span class="tab-label-custom">
                  <el-icon class="tab-icon"><User /></el-icon>
                  基本资料
                </span>
              </template>

              <div class="tab-content-wrapper">
                <div class="section-intro-banner">
                  <div class="intro-title">基本资料设置</div>
                  <div class="intro-desc">维护您的公开资料信息，修改后将实时同步至系统全局顶部导航与操作日志。</div>
                </div>

                <el-form
                  ref="basicFormRef"
                  :model="basicForm"
                  :rules="basicRules"
                  label-position="top"
                  class="profile-edit-form"
                >
                  <el-row :gutter="20">
                    <el-col :xs="24" :sm="12">
                      <el-form-item label="登录账号 (用户名)">
                        <el-input
                          :model-value="userInfo.username"
                          disabled
                          placeholder="系统主账号"
                        >
                          <template #prefix>
                            <el-icon><Lock /></el-icon>
                          </template>
                        </el-input>
                        <span class="form-tip">账号为唯一系统标识，不可更改</span>
                      </el-form-item>
                    </el-col>

                    <el-col :xs="24" :sm="12">
                      <el-form-item label="用户昵称" prop="nickName">
                        <el-input
                          v-model="basicForm.nickName"
                          placeholder="请输入显示昵称"
                          maxlength="50"
                          show-word-limit
                          clearable
                        >
                          <template #prefix>
                            <el-icon><User /></el-icon>
                          </template>
                        </el-input>
                      </el-form-item>
                    </el-col>

                    <el-col :xs="24" :sm="12">
                      <el-form-item label="电子邮箱" prop="email">
                        <el-input
                          v-model="basicForm.email"
                          placeholder="例如: admin@example.com"
                          clearable
                        >
                          <template #prefix>
                            <el-icon><Message /></el-icon>
                          </template>
                        </el-input>
                      </el-form-item>
                    </el-col>

                    <el-col :xs="24" :sm="12">
                      <el-form-item label="手机号码" prop="phonenumber">
                        <el-input
                          v-model="basicForm.phonenumber"
                          placeholder="请输入联系电话"
                          maxlength="20"
                          clearable
                        >
                          <template #prefix>
                            <el-icon><Phone /></el-icon>
                          </template>
                        </el-input>
                      </el-form-item>
                    </el-col>

                    <el-col :xs="24" :sm="12">
                      <el-form-item label="性别设置" prop="sex">
                        <el-radio-group v-model="basicForm.sex" class="gender-radio-group">
                          <el-radio-button :label="1">男</el-radio-button>
                          <el-radio-button :label="2">女</el-radio-button>
                          <el-radio-button :label="0">保密</el-radio-button>
                        </el-radio-group>
                      </el-form-item>
                    </el-col>

                    <el-col :xs="24" :sm="12">
                      <el-form-item label="头像地址 (外链或本地上传)">
                        <el-input
                          v-model="basicForm.avatar"
                          placeholder="可直接粘贴图片 URL 或点击左侧头像上传"
                          clearable
                        >
                          <template #prefix>
                            <el-icon><Link /></el-icon>
                          </template>
                        </el-input>
                      </el-form-item>
                    </el-col>
                  </el-row>

                  <!-- 头像快捷操作行 -->
                  <div class="avatar-quick-bar">
                    <div class="preview-mini">
                      <span class="preview-label">当前头像预览：</span>
                      <el-avatar :size="40" :src="basicForm.avatar || defaultAvatar" />
                    </div>
                    <div class="quick-btns">
                      <el-button size="small" :icon="Refresh" @click="handleResetAvatar">
                        恢复默认头像
                      </el-button>
                    </div>
                  </div>

                  <div class="form-actions-bar">
                    <el-button
                      type="primary"
                      :loading="saveProfileLoading"
                      :icon="Check"
                      @click="handleSaveProfile"
                    >
                      保存修改
                    </el-button>
                    <el-button :icon="Refresh" @click="handleResetBasicForm">
                      重置内容
                    </el-button>
                  </div>
                </el-form>
              </div>
            </el-tab-pane>

            <!-- Tab 2: 账号安全与密码修改 -->
            <el-tab-pane name="security">
              <template #label>
                <span class="tab-label-custom">
                  <el-icon class="tab-icon"><Lock /></el-icon>
                  账号安全
                </span>
              </template>

              <div class="tab-content-wrapper">
                <div class="section-intro-banner">
                  <div class="intro-title">账户安全中心</div>
                  <div class="intro-desc">为保证管理控制台的最高安全等级，请定期更换复杂密码并避免在不安全终端登录。</div>
                </div>

                <!-- 安全状态指示三联卡 -->
                <el-row :gutter="16" class="security-status-row">
                  <el-col :xs="24" :sm="8">
                    <div class="sec-card">
                      <div class="sec-icon-circle green">
                        <el-icon><CircleCheckFilled /></el-icon>
                      </div>
                      <div class="sec-info">
                        <div class="sec-title">身份鉴权</div>
                        <div class="sec-desc">JWT Bearer 强认证</div>
                      </div>
                    </div>
                  </el-col>

                  <el-col :xs="24" :sm="8">
                    <div class="sec-card">
                      <div class="sec-icon-circle blue">
                        <el-icon><Key /></el-icon>
                      </div>
                      <div class="sec-info">
                        <div class="sec-title">哈希存储</div>
                        <div class="sec-desc">Bcrypt 单向加盐散列</div>
                      </div>
                    </div>
                  </el-col>

                  <el-col :xs="24" :sm="8">
                    <div class="sec-card">
                      <div class="sec-icon-circle orange">
                        <el-icon><View /></el-icon>
                      </div>
                      <div class="sec-info">
                        <div class="sec-title">审计监控</div>
                        <div class="sec-desc">IP 与访问链路记录</div>
                      </div>
                    </div>
                  </el-col>
                </el-row>

                <!-- 修改密码表单 -->
                <div class="password-form-container">
                  <div class="form-title-sub">
                    <el-icon class="mr-1 text-primary"><Key /></el-icon> 修改控制台密码
                  </div>

                  <el-form
                    ref="passwordFormRef"
                    :model="passwordForm"
                    :rules="passwordRules"
                    label-position="top"
                    class="password-edit-form"
                  >
                    <el-form-item label="当前原密码" prop="oldPassword">
                      <el-input
                        v-model="passwordForm.oldPassword"
                        type="password"
                        show-password
                        placeholder="请输入当前正在使用的原密码"
                        autocomplete="off"
                        clearable
                      >
                        <template #prefix>
                          <el-icon><Lock /></el-icon>
                        </template>
                      </el-input>
                    </el-form-item>

                    <el-form-item label="设置新密码" prop="newPassword">
                      <el-input
                        v-model="passwordForm.newPassword"
                        type="password"
                        show-password
                        placeholder="请输入 6 - 20 位新密码"
                        autocomplete="off"
                        clearable
                      >
                        <template #prefix>
                          <el-icon><Key /></el-icon>
                        </template>
                      </el-input>

                      <!-- 密码强度实时动态条 -->
                      <div v-if="passwordForm.newPassword" class="password-strength-bar">
                        <div class="strength-labels">
                          <span class="hint">密码强度评估：</span>
                          <span class="score-text" :style="{ color: passwordStrength.color }">
                            {{ passwordStrength.text }}
                          </span>
                        </div>
                        <div class="strength-segments">
                          <div
                            v-for="idx in 4"
                            :key="idx"
                            class="segment-block"
                            :class="{ active: idx <= passwordStrength.score }"
                            :style="{ backgroundColor: idx <= passwordStrength.score ? passwordStrength.color : '' }"
                          />
                        </div>
                      </div>
                    </el-form-item>

                    <el-form-item label="确认新密码" prop="confirmPassword">
                      <el-input
                        v-model="passwordForm.confirmPassword"
                        type="password"
                        show-password
                        placeholder="请再次输入新密码进行确认"
                        autocomplete="off"
                        clearable
                      >
                        <template #prefix>
                          <el-icon><Check /></el-icon>
                        </template>
                      </el-input>
                    </el-form-item>

                    <div class="password-tips-box">
                      <div class="tips-title">
                        <el-icon class="mr-1 text-warning"><WarningFilled /></el-icon>
                        安全建议提示
                      </div>
                      <ul class="tips-list">
                        <li>建议长度在 8 位以上，并混合大写字母、小写字母、数字及特殊字符。</li>
                        <li>请勿使用与其他网站相同的密码，或姓名全拼、电话号码等常见组合。</li>
                        <li>修改密码后无需重新登录当前会话，但请妥善保存新密码。</li>
                      </ul>
                    </div>

                    <div class="form-actions-bar">
                      <el-button
                        type="primary"
                        :loading="savePasswordLoading"
                        :icon="Check"
                        @click="handleSavePassword"
                      >
                        确认更新密码
                      </el-button>
                      <el-button :icon="Refresh" @click="handleResetPasswordForm">
                        清空表单
                      </el-button>
                    </div>
                  </el-form>
                </div>
              </div>
            </el-tab-pane>

            <!-- Tab 3: 会话与登录审计 -->
            <el-tab-pane name="audit">
              <template #label>
                <span class="tab-label-custom">
                  <el-icon class="tab-icon"><Monitor /></el-icon>
                  会话与审计
                </span>
              </template>

              <div class="tab-content-wrapper">
                <div class="section-intro-banner">
                  <div class="intro-title">当前会话与审计环境</div>
                  <div class="intro-desc">监控当前浏览器终端的登录环境与网络接入点，保障后台系统审计合规。</div>
                </div>

                <div class="session-detail-card">
                  <div class="session-header">
                    <div class="session-status">
                      <span class="status-pulse-dot" />
                      <span class="status-text">当前会话活跃中</span>
                    </div>
                    <el-tag type="success" effect="light" round>安全已认证</el-tag>
                  </div>

                  <el-descriptions :column="2" border class="session-descriptions">
                    <el-descriptions-item label="接入客户端操作系统">
                      <el-tag size="small" type="info" effect="plain">{{ clientEnv.os }}</el-tag>
                    </el-descriptions-item>

                    <el-descriptions-item label="当前浏览器核心">
                      <el-tag size="small" type="info" effect="plain">{{ clientEnv.browser }}</el-tag>
                    </el-descriptions-item>

                    <el-descriptions-item label="本次登录 IP">
                      <span class="font-mono">{{ userInfo.loginIp || "未知" }}</span>
                    </el-descriptions-item>

                    <el-descriptions-item label="本次登录建立时间">
                      <span>{{ userInfo.loginTime || "-" }}</span>
                    </el-descriptions-item>

                    <el-descriptions-item label="历史累计登录">
                      <span class="font-semibold text-primary">{{ userInfo.loginCount }}</span> 次
                    </el-descriptions-item>

                    <el-descriptions-item label="账号注册记录">
                      <span>{{ userInfo.registerTime || "-" }}</span>
                    </el-descriptions-item>
                  </el-descriptions>

                  <div class="audit-link-actions">
                    <div class="audit-notice">
                      <el-icon class="mr-1 text-primary"><InfoFilled /></el-icon>
                      如需查看历史每一次精确的操作记录或拦截日志，请进入系统审计中心：
                    </div>
                    <div class="audit-btn-group">
                      <el-button type="primary" plain :icon="Right" @click="goToPage('/system/log')">
                        查看完整登录日志
                      </el-button>
                      <el-button type="success" plain :icon="Right" @click="goToPage('/system/access')">
                        查看实时访问日志
                      </el-button>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, reactive, computed } from "vue";
import { useRouter } from "vue-router";
import type { FormInstance, FormRules } from "element-plus";
import {
  getUserInfo,
  updateProfile,
  changePassword,
  type UpdateProfileRequest
} from "@/api/profile";
import { upload } from "@/api/common";
import { useUserStoreHook } from "@/store/modules/user";
import { storageLocal } from "@pureadmin/utils";
import { userKey, type DataInfo } from "@/utils/auth";
import { message } from "@/utils/message";
import { formatAxis } from "@/utils/formatTime";
import {
  Camera,
  Lock,
  Phone,
  Message,
  User,
  Monitor,
  Calendar,
  Clock,
  Location,
  Key,
  Check,
  Refresh,
  DocumentCopy,
  CircleCheckFilled,
  WarningFilled,
  InfoFilled,
  Right,
  Link,
  View
} from "@element-plus/icons-vue";
import defaultAvatar from "@/assets/user.jpg";

const router = useRouter();

// 选项卡状态
const activeTab = ref("basic");

// 加载状态
const pageLoading = ref(true);
const uploadLoading = ref(false);
const saveProfileLoading = ref(false);
const savePasswordLoading = ref(false);

// 用户完整信息模型
interface UserProfileState {
  userId: number;
  username: string;
  nickName: string;
  email: string;
  phonenumber: string;
  sex: number;
  avatar: string;
  role: string;
  loginIp: string;
  loginTime: string;
  loginCount: number;
  registerTime: string;
}

const userInfo = ref<UserProfileState>({
  userId: 0,
  username: "",
  nickName: "加载中...",
  email: "",
  phonenumber: "",
  sex: 0,
  avatar: "",
  role: "超级管理员",
  loginIp: "",
  loginTime: "",
  loginCount: 0,
  registerTime: ""
});

// 基本资料表单
const basicFormRef = ref<FormInstance>();
const basicForm = reactive({
  nickName: "",
  email: "",
  phonenumber: "",
  sex: 0,
  avatar: ""
});

const basicRules = reactive<FormRules>({
  nickName: [
    { required: true, message: "请输入用户昵称", trigger: "blur" },
    { max: 50, message: "昵称长度不能超过 50 个字符", trigger: "blur" }
  ],
  email: [
    {
      type: "email",
      message: "请输入正确的邮箱格式",
      trigger: ["blur", "change"]
    }
  ],
  phonenumber: [
    {
      pattern: /^1[3-9]\d{9}$|^$/,
      message: "请输入正确的 11 位手机号码",
      trigger: "blur"
    }
  ]
});

// 修改密码表单
const passwordFormRef = ref<FormInstance>();
const passwordForm = reactive({
  oldPassword: "",
  newPassword: "",
  confirmPassword: ""
});

const equalToNewPassword = (_rule: any, value: string, callback: any) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error("两次输入的密码不一致"));
  } else {
    callback();
  }
};

const passwordRules = reactive<FormRules>({
  oldPassword: [{ required: true, message: "请输入当前原密码", trigger: "blur" }],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, max: 20, message: "密码长度必须在 6 到 20 个字符之间", trigger: "blur" }
  ],
  confirmPassword: [
    { required: true, message: "请再次输入新密码进行确认", trigger: "blur" },
    { validator: equalToNewPassword, trigger: "blur" }
  ]
});

// 问候语计算
const currentTimeText = computed(() => {
  return formatAxis(new Date());
});

// 格式化注册年份
const formatRegisterYear = (dateStr: string) => {
  if (!dateStr || dateStr === "-") return "2024";
  const year = dateStr.slice(0, 4);
  return year.length === 4 ? `${year}年` : "2024年";
};

// 客户端环境提取
const clientEnv = computed(() => {
  if (typeof navigator === "undefined") {
    return { os: "未知系统", browser: "现代浏览器" };
  }
  const ua = navigator.userAgent;
  let os = "Linux / Unix";
  if (ua.includes("Win")) os = "Windows";
  else if (ua.includes("Mac")) os = "macOS";
  else if (ua.includes("Android")) os = "Android";
  else if (ua.includes("iPhone") || ua.includes("iPad")) os = "iOS";

  let browser = "现代浏览器";
  if (ua.includes("Edg/")) browser = "Microsoft Edge";
  else if (ua.includes("Chrome/")) browser = "Google Chrome";
  else if (ua.includes("Firefox/")) browser = "Mozilla Firefox";
  else if (ua.includes("Safari/") && !ua.includes("Chrome/")) browser = "Apple Safari";

  return { os, browser };
});

// 密码强度评估
const passwordStrength = computed(() => {
  const pwd = passwordForm.newPassword;
  if (!pwd) return { score: 0, text: "未输入", color: "#909399" };
  let score = 0;
  if (pwd.length >= 6) score += 1;
  if (pwd.length >= 10) score += 1;
  if (/[0-9]/.test(pwd) && /[a-zA-Z]/.test(pwd)) score += 1;
  if (/[^0-9a-zA-Z]/.test(pwd)) score += 1;

  if (score <= 1) return { score: 1, text: "弱", color: "#f56c6c" };
  if (score === 2) return { score: 2, text: "中", color: "#e6a23c" };
  if (score === 3) return { score: 3, text: "强", color: "#67c23a" };
  return { score: 4, text: "极强", color: "#409eff" };
});

// 获取个人信息
const fetchUserInfo = async () => {
  pageLoading.value = true;
  try {
    const res = await getUserInfo();
    if (res.code === 200 && res.payload) {
      const p = res.payload;
      userInfo.value = {
        userId: p.userId || 1,
        username: p.username || "",
        nickName: p.nickName || p.username || "",
        email: p.email || "",
        phonenumber: p.phonenumber || "",
        sex: p.sex ?? 0,
        avatar: p.avatar || "",
        role: p.role || "超级管理员",
        loginIp: p.loginIp || "未知",
        loginTime: p.loginTime || "-",
        loginCount: p.loginCount || 0,
        registerTime: p.registerTime || "-"
      };

      // 同步给基本表单
      basicForm.nickName = userInfo.value.nickName;
      basicForm.email = userInfo.value.email;
      basicForm.phonenumber = userInfo.value.phonenumber;
      basicForm.sex = userInfo.value.sex;
      basicForm.avatar = userInfo.value.avatar;
    }
  } catch (error) {
    message("获取用户信息失败", { type: "error" });
  } finally {
    pageLoading.value = false;
  }
};

// 复制 IP 地址
const handleCopyIp = async () => {
  if (!userInfo.value.loginIp) return;
  try {
    await navigator.clipboard.writeText(userInfo.value.loginIp);
    message("登录 IP 地址已复制到剪贴板", { type: "success" });
  } catch {
    message("复制失败，请手动选择复制", { type: "error" });
  }
};

// 头像上传前校验
const beforeAvatarUpload = (file: File) => {
  const isImage = file.type.startsWith("image/");
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isImage) {
    message("头像图片只能是图片格式 (JPG/PNG/WEBP)!", { type: "error" });
    return false;
  }
  if (!isLt2M) {
    message("头像图片大小不能超过 2MB!", { type: "error" });
    return false;
  }
  return true;
};

// 自定义上传逻辑（带鉴权）
const handleCustomUpload = async (options: any) => {
  const file = options.file;
  if (!beforeAvatarUpload(file)) return;

  const formData = new FormData();
  formData.append("file", file);

  uploadLoading.value = true;
  try {
    const res = await upload(formData);
    if (res.code === 200 && res.payload) {
      const newAvatarUrl = res.payload;
      basicForm.avatar = newAvatarUrl;
      userInfo.value.avatar = newAvatarUrl;

      // 立即自动更新资料中的头像
      const updateRes = await updateProfile({ avatar: newAvatarUrl });
      if (updateRes.code === 200) {
        useUserStoreHook().SET_AVATAR(newAvatarUrl);
        const userStorage = storageLocal().getItem<DataInfo<number>>(userKey);
        if (userStorage) {
          storageLocal().setItem(userKey, {
            ...userStorage,
            avatar: newAvatarUrl
          });
        }
        message("头像更换并保存成功", { type: "success" });
      } else {
        message("头像上传成功，请点击“保存修改”确认资料更新", { type: "info" });
      }
    } else {
      message(res.message || "头像上传失败", { type: "error" });
    }
  } catch (error) {
    message("上传头像过程中发生错误", { type: "error" });
  } finally {
    uploadLoading.value = false;
  }
};

// 恢复默认头像
const handleResetAvatar = async () => {
  basicForm.avatar = "";
  userInfo.value.avatar = "";
  try {
    const res = await updateProfile({ avatar: "" });
    if (res.code === 200) {
      useUserStoreHook().SET_AVATAR("");
      const userStorage = storageLocal().getItem<DataInfo<number>>(userKey);
      if (userStorage) {
        storageLocal().setItem(userKey, {
          ...userStorage,
          avatar: ""
        });
      }
      message("已恢复为默认头像", { type: "success" });
    }
  } catch {
    message("重置头像失败", { type: "error" });
  }
};

// 保存资料修改
const handleSaveProfile = async () => {
  if (!basicFormRef.value) return;
  await basicFormRef.value.validate(async valid => {
    if (!valid) return;
    saveProfileLoading.value = true;
    try {
      const payload: UpdateProfileRequest = {
        nickName: basicForm.nickName,
        email: basicForm.email,
        phonenumber: basicForm.phonenumber,
        sex: Number(basicForm.sex),
        avatar: basicForm.avatar
      };
      const res = await updateProfile(payload);
      if (res.code === 200) {
        message("个人资料更新成功", { type: "success" });
        userInfo.value.nickName = basicForm.nickName;
        userInfo.value.email = basicForm.email;
        userInfo.value.phonenumber = basicForm.phonenumber;
        userInfo.value.sex = basicForm.sex;
        userInfo.value.avatar = basicForm.avatar;

        // 全局同步更新 Pinia store 与 local storage
        useUserStoreHook().SET_NICKNAME(basicForm.nickName);
        useUserStoreHook().SET_AVATAR(basicForm.avatar);
        const userStorage = storageLocal().getItem<DataInfo<number>>(userKey);
        if (userStorage) {
          storageLocal().setItem(userKey, {
            ...userStorage,
            nickname: basicForm.nickName,
            avatar: basicForm.avatar
          });
        }
      } else {
        message(res.message || "更新资料失败", { type: "error" });
      }
    } catch (error) {
      message("更新个人资料过程中发生异常", { type: "error" });
    } finally {
      saveProfileLoading.value = false;
    }
  });
};

// 重置基本表单
const handleResetBasicForm = () => {
  basicForm.nickName = userInfo.value.nickName;
  basicForm.email = userInfo.value.email;
  basicForm.phonenumber = userInfo.value.phonenumber;
  basicForm.sex = userInfo.value.sex;
  basicForm.avatar = userInfo.value.avatar;
  message("表单已重置回当前资料", { type: "info" });
};

// 提交密码修改
const handleSavePassword = async () => {
  if (!passwordFormRef.value) return;
  await passwordFormRef.value.validate(async valid => {
    if (!valid) return;
    savePasswordLoading.value = true;
    try {
      const res = await changePassword({
        oldPassword: passwordForm.oldPassword,
        newPassword: passwordForm.newPassword,
        confirmPassword: passwordForm.confirmPassword
      });
      if (res.code === 200) {
        message("密码修改成功，请妥善保管新密码", { type: "success" });
        passwordForm.oldPassword = "";
        passwordForm.newPassword = "";
        passwordForm.confirmPassword = "";
        passwordFormRef.value.resetFields();
      } else {
        message(res.message || "修改密码失败", { type: "error" });
      }
    } catch (error) {
      message("修改密码过程中发生异常", { type: "error" });
    } finally {
      savePasswordLoading.value = false;
    }
  });
};

// 清空密码表单
const handleResetPasswordForm = () => {
  if (passwordFormRef.value) {
    passwordFormRef.value.resetFields();
  }
};

// 路由跳转
const goToPage = (path: string) => {
  router.push(path);
};

onMounted(() => {
  fetchUserInfo();
});
</script>

<style scoped lang="scss">
.profile-container {
  padding: 16px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 120px);

  .profile-layout {
    margin: 0 !important;
  }

  /* 左侧用户画像卡片 */
  .user-profile-card {
    border-radius: 14px;
    border: 1px solid var(--el-border-color-lighter);
    background-color: var(--el-bg-color-overlay);
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 4px 20px -4px rgba(0, 0, 0, 0.05);

    &:hover {
      box-shadow: 0 8px 30px -4px rgba(0, 0, 0, 0.1);
    }

    .card-cover-banner {
      height: 110px;
      position: relative;
      background: linear-gradient(135deg, #3b82f6 0%, #6366f1 50%, #8b5cf6 100%);
      overflow: hidden;

      .banner-mesh-glow {
        position: absolute;
        inset: 0;
        background: radial-gradient(circle at 20% 40%, rgba(255, 255, 255, 0.25) 0%, transparent 60%),
                    radial-gradient(circle at 80% 80%, rgba(255, 255, 255, 0.2) 0%, transparent 50%);
      }

      .role-tag-float {
        position: absolute;
        top: 14px;
        right: 14px;

        .role-pill {
          background: rgba(0, 0, 0, 0.35);
          backdrop-filter: blur(8px);
          border: 1px solid rgba(255, 255, 255, 0.3);
          color: #fff;
          font-weight: 500;
          padding: 4px 10px;
        }
      }
    }

    .user-meta-body {
      padding: 0 20px 24px;
      display: flex;
      flex-direction: column;
      align-items: center;

      .avatar-wrapper {
        position: relative;
        margin-top: -50px;
        margin-bottom: 12px;

        .avatar-uploader {
          .avatar-box {
            position: relative;
            width: 96px;
            height: 96px;
            border-radius: 50%;
            overflow: hidden;
            border: 4px solid var(--el-bg-color-overlay);
            box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);
            cursor: pointer;
            transition: all 0.25s ease;

            .user-avatar {
              width: 100%;
              height: 100%;
              object-fit: cover;
              display: block;
            }

            .avatar-mask {
              position: absolute;
              inset: 0;
              background: rgba(0, 0, 0, 0.55);
              backdrop-filter: blur(2px);
              display: flex;
              flex-direction: column;
              align-items: center;
              justify-content: center;
              color: #ffffff;
              opacity: 0;
              transition: opacity 0.25s ease;

              .mask-icon {
                font-size: 22px;
                margin-bottom: 2px;
              }

              .mask-text {
                font-size: 11px;
                font-weight: 500;
                letter-spacing: 0.5px;
              }
            }

            &:hover {
              transform: scale(1.02);

              .avatar-mask {
                opacity: 1;
              }
            }

            &.is-uploading {
              .avatar-mask {
                opacity: 1;
              }
            }
          }
        }

        .status-indicator-dot {
          position: absolute;
          bottom: 4px;
          right: 4px;
          width: 14px;
          height: 14px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;

          .ping-circle {
            position: absolute;
            width: 100%;
            height: 100%;
            border-radius: 50%;
            background-color: #10b981;
            opacity: 0.75;
            animation: ping 2s cubic-bezier(0, 0, 0.2, 1) infinite;
          }

          .dot-core {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #10b981;
            border: 2px solid var(--el-bg-color-overlay);
          }
        }
      }

      .name-section {
        text-align: center;
        margin-bottom: 12px;

        .user-nickname {
          font-size: 20px;
          font-weight: 700;
          color: var(--el-text-color-primary);
          margin: 0 0 4px;
        }

        .user-handle-wrap {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
          font-size: 13px;

          .user-handle {
            color: var(--el-color-primary);
            font-weight: 500;
          }

          .user-uid {
            color: var(--el-text-color-secondary);
            font-size: 12px;
            background: var(--el-fill-color-light);
            padding: 1px 6px;
            border-radius: 4px;
          }
        }
      }

      .bio-quote-box {
        margin: 4px 0 16px;
        padding: 10px 14px;
        background-color: var(--el-fill-color-lighter);
        border-radius: 8px;
        font-size: 13px;
        line-height: 1.5;
        color: var(--el-text-color-regular);
        text-align: center;
        border-left: 3px solid var(--el-color-primary);

        .quote-symbol {
          color: var(--el-color-primary);
          font-weight: bold;
          font-size: 15px;
          padding: 0 2px;
        }
      }

      .stats-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 10px;
        width: 100%;
        margin-bottom: 4px;

        .stat-card {
          background-color: var(--el-fill-color-light);
          border-radius: 8px;
          padding: 12px 10px;
          text-align: center;
          transition: background-color 0.2s;

          &:hover {
            background-color: var(--el-fill-color);
          }

          .stat-num {
            font-size: 18px;
            font-weight: 700;
            color: var(--el-text-color-primary);
            margin-bottom: 2px;
          }

          .stat-label {
            font-size: 12px;
            color: var(--el-text-color-secondary);
            display: flex;
            align-items: center;
            justify-content: center;
          }
        }
      }

      .card-divider {
        margin: 16px 0;
        width: 100%;
      }

      .detail-info-list {
        width: 100%;
        display: flex;
        flex-direction: column;
        gap: 12px;

        .detail-item {
          display: flex;
          align-items: center;
          justify-content: space-between;
          font-size: 13px;
          padding: 6px 4px;
          border-radius: 6px;

          .item-left {
            display: flex;
            align-items: center;
            gap: 10px;
            color: var(--el-text-color-regular);

            .item-icon-box {
              width: 28px;
              height: 28px;
              border-radius: 6px;
              display: flex;
              align-items: center;
              justify-content: center;
              font-size: 14px;

              &.primary {
                background-color: var(--el-color-primary-light-9);
                color: var(--el-color-primary);
              }
              &.success {
                background-color: var(--el-color-success-light-9);
                color: var(--el-color-success);
              }
              &.warning {
                background-color: var(--el-color-warning-light-9);
                color: var(--el-color-warning);
              }
              &.info {
                background-color: var(--el-color-info-light-9);
                color: var(--el-color-info);
              }
              &.purple {
                background-color: rgba(139, 92, 246, 0.12);
                color: #8b5cf6;
              }
            }

            .item-title {
              font-weight: 500;
            }
          }

          .item-value {
            color: var(--el-text-color-primary);
            font-weight: 500;
            text-align: right;

            &.ip-copy-val {
              display: flex;
              align-items: center;
              gap: 4px;

              .ip-tag {
                font-family: monospace;
              }

              .copy-btn {
                padding: 2px 4px;
              }
            }

            &.time-val {
              font-size: 12px;
              color: var(--el-text-color-secondary);
              font-family: monospace;
            }
          }
        }
      }
    }
  }

  /* 右侧设置面板卡片 */
  .settings-card {
    border-radius: 14px;
    border: 1px solid var(--el-border-color-lighter);
    background-color: var(--el-bg-color-overlay);
    box-shadow: 0 4px 20px -4px rgba(0, 0, 0, 0.05);

    :deep(.el-card__body) {
      padding: 10px 24px 24px;
    }

    .profile-tabs {
      :deep(.el-tabs__header) {
        margin-bottom: 20px;
        border-bottom: 1px solid var(--el-border-color-lighter);
      }

      :deep(.el-tabs__item) {
        font-size: 15px;
        font-weight: 500;
        padding: 0 20px;
        height: 48px;
        line-height: 48px;
      }

      .tab-label-custom {
        display: inline-flex;
        align-items: center;
        gap: 6px;

        .tab-icon {
          font-size: 16px;
        }
      }
    }

    .tab-content-wrapper {
      .section-intro-banner {
        margin-bottom: 22px;
        padding-bottom: 14px;
        border-bottom: 1px dashed var(--el-border-color-light);

        .intro-title {
          font-size: 16px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          margin-bottom: 4px;
        }

        .intro-desc {
          font-size: 13px;
          color: var(--el-text-color-secondary);
        }
      }

      .form-tip {
        font-size: 12px;
        color: var(--el-text-color-placeholder);
        margin-top: 4px;
        display: block;
      }

      .gender-radio-group {
        width: 100%;
        display: flex;

        :deep(.el-radio-button) {
          flex: 1;

          .el-radio-button__inner {
            width: 100%;
          }
        }
      }

      .avatar-quick-bar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 10px 14px;
        background-color: var(--el-fill-color-light);
        border-radius: 8px;
        margin-bottom: 22px;

        .preview-mini {
          display: flex;
          align-items: center;
          gap: 10px;

          .preview-label {
            font-size: 13px;
            color: var(--el-text-color-secondary);
          }
        }
      }

      .form-actions-bar {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-top: 24px;
        padding-top: 16px;
        border-top: 1px solid var(--el-border-color-lighter);
      }

      /* 安全设置选项卡特定样式 */
      .security-status-row {
        margin-bottom: 24px;

        .sec-card {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 14px;
          border-radius: 10px;
          background-color: var(--el-fill-color-light);
          border: 1px solid var(--el-border-color-lighter);
          transition: all 0.25s ease;

          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
          }

          .sec-icon-circle {
            width: 38px;
            height: 38px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 18px;

            &.green {
              background-color: var(--el-color-success-light-9);
              color: var(--el-color-success);
            }
            &.blue {
              background-color: var(--el-color-primary-light-9);
              color: var(--el-color-primary);
            }
            &.orange {
              background-color: var(--el-color-warning-light-9);
              color: var(--el-color-warning);
            }
          }

          .sec-info {
            .sec-title {
              font-size: 14px;
              font-weight: 600;
              color: var(--el-text-color-primary);
              margin-bottom: 2px;
            }

            .sec-desc {
              font-size: 12px;
              color: var(--el-text-color-secondary);
            }
          }
        }
      }

      .password-form-container {
        max-width: 600px;

        .form-title-sub {
          font-size: 15px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          margin-bottom: 16px;
          display: flex;
          align-items: center;
        }

        .password-strength-bar {
          margin-top: 8px;
          width: 100%;

          .strength-labels {
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 12px;
            margin-bottom: 4px;

            .hint {
              color: var(--el-text-color-secondary);
            }

            .score-text {
              font-weight: 600;
            }
          }

          .strength-segments {
            display: flex;
            gap: 6px;

            .segment-block {
              flex: 1;
              height: 4px;
              border-radius: 2px;
              background-color: var(--el-fill-color);
              transition: all 0.3s ease;
            }
          }
        }

        .password-tips-box {
          margin-top: 18px;
          padding: 12px 16px;
          border-radius: 8px;
          background-color: var(--el-color-warning-light-9);
          border: 1px solid var(--el-color-warning-light-7);

          .tips-title {
            font-size: 13px;
            font-weight: 600;
            color: var(--el-color-warning-dark-2);
            margin-bottom: 6px;
            display: flex;
            align-items: center;
          }

          .tips-list {
            margin: 0;
            padding-left: 18px;
            font-size: 12px;
            color: var(--el-text-color-regular);
            line-height: 1.6;
          }
        }
      }

      /* 会话与审计选项卡特定样式 */
      .session-detail-card {
        .session-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 16px;

          .session-status {
            display: flex;
            align-items: center;
            gap: 8px;

            .status-pulse-dot {
              width: 10px;
              height: 10px;
              border-radius: 50%;
              background-color: var(--el-color-success);
              box-shadow: 0 0 0 3px var(--el-color-success-light-8);
              animation: pulse 2s infinite;
            }

            .status-text {
              font-size: 14px;
              font-weight: 600;
              color: var(--el-text-color-primary);
            }
          }
        }

        .session-descriptions {
          margin-bottom: 24px;
        }

        .audit-link-actions {
          padding: 16px;
          border-radius: 10px;
          background-color: var(--el-fill-color-light);
          border: 1px solid var(--el-border-color-lighter);

          .audit-notice {
            font-size: 13px;
            color: var(--el-text-color-regular);
            margin-bottom: 12px;
            display: flex;
            align-items: center;
          }

          .audit-btn-group {
            display: flex;
            gap: 12px;
            flex-wrap: wrap;
          }
        }
      }
    }
  }
}

@keyframes ping {
  75%, 100% {
    transform: scale(2);
    opacity: 0;
  }
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(103, 194, 58, 0.7);
  }
  70% {
    box-shadow: 0 0 0 8px rgba(103, 194, 58, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(103, 194, 58, 0);
  }
}

/* 响应式调整 */
@media (max-width: 992px) {
  .profile-container {
    padding: 10px;

    .left-col {
      margin-bottom: 16px;
    }
  }
}
</style>
