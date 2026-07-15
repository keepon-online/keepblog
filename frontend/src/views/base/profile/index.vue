<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">个人信息</span>
        </div>
      </template>

      <div class="profile-content">
        <div class="profile-header">
          <div class="avatar-section">
            <el-upload
              class="avatar-uploader"
              :action="uploadAction"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
            >
              <img
                v-if="userInfo.avatar"
                :src="userInfo.avatar"
                class="avatar"
              />
              <el-avatar v-else :size="100" :src="defaultAvatar" />
              <div class="avatar-overlay">
                <el-icon class="avatar-icon"><Camera /></el-icon>
                <span class="avatar-text">更换头像</span>
              </div>
            </el-upload>
            <div class="avatar-tip">点击更换头像</div>
          </div>

          <div class="user-info">
            <div class="welcome-text">
              {{ currentTime }}，{{
                userInfo.nickname
              }}，生活变的再糟糕，也不妨碍我变得更好！
            </div>

            <el-row :gutter="20" class="info-grid">
              <el-col :xs="24" :sm="12" :md="8">
                <div class="info-item">
                  <div class="info-label">昵称：</div>
                  <div class="info-value">{{ userInfo.nickname }}</div>
                </div>
              </el-col>
              <el-col :xs="24" :sm="12" :md="8">
                <div class="info-item">
                  <div class="info-label">身份：</div>
                  <div class="info-value">{{ userInfo.role }}</div>
                </div>
              </el-col>
              <el-col :xs="24" :sm="12" :md="8">
                <div class="info-item">
                  <div class="info-label">登录IP：</div>
                  <div class="info-value">{{ userInfo.loginIp }}</div>
                </div>
              </el-col>
              <el-col :xs="24" :sm="12" :md="8">
                <div class="info-item">
                  <div class="info-label">登录时间：</div>
                  <div class="info-value">{{ userInfo.loginTime }}</div>
                </div>
              </el-col>
              <el-col :xs="24" :sm="12" :md="8">
                <div class="info-item">
                  <div class="info-label">累计登录：</div>
                  <div class="info-value">{{ userInfo.loginCount }}次</div>
                </div>
              </el-col>
              <el-col :xs="24" :sm="12" :md="8">
                <div class="info-item">
                  <div class="info-label">注册时间：</div>
                  <div class="info-value">{{ userInfo.registerTime }}</div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="security-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">账号安全</span>
        </div>
      </template>

      <div class="security-content">
        <div
          class="security-item"
          v-for="item in securityItems"
          :key="item.key"
        >
          <div class="security-item-main">
            <div class="security-item-icon">
              <el-icon :size="20"><component :is="item.icon" /></el-icon>
            </div>
            <div class="security-item-info">
              <div class="security-item-title">{{ item.title }}</div>
              <div class="security-item-desc">{{ item.desc }}</div>
            </div>
          </div>
          <el-button
            type="primary"
            link
            @click="handleSecurityAction(item.key)"
          >
            {{ item.action }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="dialogFormVisible"
      title="修改密码"
      width="500px"
      :before-close="handleClose"
    >
      <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        label-width="90px"
        status-icon
      >
        <el-form-item
          label="原密码"
          prop="oldPassword"
        >
          <el-input
            v-model="ruleForm.oldPassword"
            type="password"
            show-password
            placeholder="请输入原密码"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item
          label="新密码"
          prop="newPassword"
        >
          <el-input
            v-model="ruleForm.newPassword"
            type="password"
            show-password
            placeholder="请输入新密码"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item
          label="确认密码"
          prop="confirmPassword"
        >
          <el-input
            v-model="ruleForm.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
            autocomplete="off"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelForm">取消</el-button>
          <el-button
            type="primary"
            :loading="submitLoading"
            @click="submitForm(ruleFormRef)"
          >
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, reactive, computed } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { changePassword, getUserInfo, updateProfile, type UserInfoResponse } from "@/api/profile";
import { upload } from "@/api/common";
import { message } from "@/utils/message";
import { formatAxis } from "@/utils/formatTime";
import {
  Camera,
  Lock,
  Iphone,
  ChatLineSquare,
  Link
} from "@element-plus/icons-vue";
import defaultAvatar from "@/assets/user.jpg";

const dialogFormVisible = ref(false);
const submitLoading = ref(false);
const loading = ref(true);
const uploadAction = import.meta.env.VITE_BASE_URL + "/api/upload/images";

interface UserInfo {
  avatar: string;
  nickname: string;
  role: string;
  loginIp: string;
  loginTime: string;
  loginCount: number;
  registerTime: string;
}

const userInfo = ref<UserInfo>({
  avatar: "",
  nickname: "加载中...",
  role: "",
  loginIp: "",
  loginTime: "",
  loginCount: 0,
  registerTime: ""
});

interface SecurityItem {
  key: string;
  title: string;
  desc: string;
  action: string;
  icon: any;
}

const securityItems: SecurityItem[] = [
  {
    key: "password",
    title: "账户密码",
    desc: "当前密码强度：强",
    action: "修改",
    icon: Lock
  },
  {
    key: "phone",
    title: "密保手机",
    desc: "已绑定手机：132****4108",
    action: "修改",
    icon: Iphone
  },
  {
    key: "question",
    title: "密保问题",
    desc: "已设置密保问题，账号安全大幅度提升",
    action: "设置",
    icon: ChatLineSquare
  },
  {
    key: "qq",
    title: "绑定QQ",
    desc: "已绑定QQ：110****566",
    action: "设置",
    icon: Link
  }
];

interface RuleForm {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
}

const ruleFormRef = ref<FormInstance>();
const ruleForm = reactive<RuleForm>({
  oldPassword: "",
  newPassword: "",
  confirmPassword: ""
});

const equalToPassword = (rule: any, value: string, callback: any) => {
  if (ruleForm.newPassword !== value) {
    callback(new Error("两次输入的密码不一致"));
  } else {
    callback();
  }
};

const rules = reactive<FormRules>({
  oldPassword: [{ required: true, message: "请输入原密码", trigger: "blur" }],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, max: 20, message: "长度在 6 到 20 个字符", trigger: "blur" }
  ],
  confirmPassword: [
    { required: true, message: "请再次输入新密码", trigger: "blur" },
    { validator: equalToPassword, trigger: "blur" }
  ]
});

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;

  await formEl.validate(async (valid, fields) => {
    if (valid) {
      submitLoading.value = true;
      try {
        const res = await changePassword(ruleForm);
        if (res.code === 200) {
          message("密码修改成功", { type: "success" });
          dialogFormVisible.value = false;
          resetForm(formEl);
        } else {
          message(res.message || "修改失败", { type: "error" });
        }
      } catch (error) {
        message("修改过程中发生错误", { type: "error" });
      } finally {
        submitLoading.value = false;
      }
    }
  });
};

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields();
};

const cancelForm = () => {
  dialogFormVisible.value = false;
  resetForm(ruleFormRef.value);
};

const handleClose = (done: () => void) => {
  if (submitLoading.value) {
    message("请等待提交完成", { type: "warning" });
    return;
  }
  done();
};

const handleSecurityAction = (key: string) => {
  switch (key) {
    case "password":
      handleUpdatePassword();
      break;
    default:
      message(`暂不支持${key}操作`, { type: "info" });
  }
};

const handleUpdatePassword = () => {
  ruleForm.newPassword = "";
  ruleForm.oldPassword = "";
  ruleForm.confirmPassword = "";
  dialogFormVisible.value = true;
};

// 头像上传相关
const handleAvatarSuccess = async (response: any, _file: any) => {
  if (response.code === 200) {
    // 上传成功后更新用户头像
    try {
      const res = await updateProfile({ avatar: response.payload });
      if (res.code === 200) {
        userInfo.value.avatar = response.payload;
        message("头像更换成功", { type: "success" });
      } else {
        message(res.message || "保存头像失败", { type: "error" });
      }
    } catch (error) {
      message("保存头像失败", { type: "error" });
    }
  } else {
    message(response.message || "上传失败", { type: "error" });
  }
};

const beforeAvatarUpload = (file: any) => {
  const isJPG = file.type === "image/jpeg" || file.type === "image/png";
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isJPG) {
    message("头像图片只能是 JPG 或 PNG 格式!", { type: "error" });
  }
  if (!isLt2M) {
    message("头像图片大小不能超过 2MB!", { type: "error" });
  }
  return isJPG && isLt2M;
};

// 当前时间提示语
const currentTime = computed(() => {
  return formatAxis(new Date());
});

// 获取用户信息
const fetchUserInfo = async () => {
  loading.value = true;
  try {
    const res = await getUserInfo();
    if (res.code === 200 && res.payload) {
      userInfo.value = {
        avatar: res.payload.avatar || "",
        nickname: res.payload.nickName || res.payload.username || "",
        role: res.payload.role || "用户",
        loginIp: res.payload.loginIp || "未知",
        loginTime: res.payload.loginTime || "-",
        loginCount: res.payload.loginCount || 0,
        registerTime: res.payload.registerTime || "-"
      };
    }
  } catch (error) {
    message("获取用户信息失败", { type: "error" });
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchUserInfo();
});
</script>

<style scoped lang="scss">
.profile-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 150px);

  .profile-card,
  .security-card {
    margin-bottom: 20px;
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

  .profile-content {
    .profile-header {
      display: flex;
      gap: 30px;
      padding: 20px;

      .avatar-section {
        display: flex;
        flex-direction: column;
        align-items: center;

        .avatar-uploader {
          position: relative;
          width: 100px;
          height: 100px;
          border-radius: 50%;
          overflow: hidden;
          cursor: pointer;

          .avatar {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }

          .avatar-overlay {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.6);
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            color: white;
            opacity: 0;
            transition: opacity 0.3s;

            .avatar-icon {
              font-size: 24px;
              margin-bottom: 5px;
            }

            .avatar-text {
              font-size: 12px;
            }
          }

          &:hover {
            .avatar-overlay {
              opacity: 1;
            }
          }
        }

        .avatar-tip {
          margin-top: 10px;
          font-size: 12px;
          color: var(--el-color-info);
        }
      }

      .user-info {
        flex: 1;

        .welcome-text {
          font-size: 18px;
          font-weight: 500;
          margin-bottom: 20px;
          color: var(--el-text-color-primary);
        }

        .info-grid {
          .info-item {
            display: flex;
            margin-bottom: 15px;

            .info-label {
              width: 80px;
              color: var(--el-text-color-secondary);
              font-weight: 500;
            }

            .info-value {
              flex: 1;
              color: var(--el-text-color-primary);
            }
          }
        }
      }
    }
  }

  .security-content {
    padding: 10px 0;

    .security-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 15px 20px;
      border-bottom: 1px solid var(--el-border-color-light);

      &:last-child {
        border-bottom: none;
      }

      .security-item-main {
        display: flex;
        align-items: center;
        gap: 15px;

        .security-item-icon {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          background-color: var(--el-color-primary-light-9);
          display: flex;
          justify-content: center;
          align-items: center;
          color: var(--el-color-primary);
        }

        .security-item-info {
          .security-item-title {
            font-size: 15px;
            font-weight: 500;
            color: var(--el-text-color-primary);
            margin-bottom: 5px;
          }

          .security-item-desc {
            font-size: 13px;
            color: var(--el-text-color-secondary);
          }
        }
      }
    }
  }
}

// 响应式优化
@media (max-width: 768px) {
  .profile-container {
    padding: 12px;

    .profile-content {
      .profile-header {
        flex-direction: column;
        align-items: center;
        gap: 20px;
        padding: 15px;

        .user-info {
          .welcome-text {
            font-size: 16px;
            text-align: center;
          }

          .info-grid {
            .el-col {
              margin-bottom: 10px;
            }
          }
        }
      }
    }

    .security-content {
      .security-item {
        padding: 12px 15px;
      }
    }
  }
}
</style>
