<template>
  <div class="ai-config">
    <el-card class="ai-config-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">AI 配置</span>
          <el-button type="primary" :loading="saving" @click="save">
            保存
          </el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        title="此处配置为覆盖层：字段留空即回落服务器 config.yaml 基线；保存后立即生效（无需重启）。apiKey 只存服务端数据库，界面永不明文回显。"
        style="margin-bottom: 18px"
      />

      <el-form label-width="110px" class="ai-config-form">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="接入点">
              <el-input
                v-model="form.baseURL"
                placeholder="https://open.bigmodel.cn/api/coding/paas/v4"
                clearable
              >
                <template #suffix>
                  <SourceTag :from-db="fromDB.baseURL" />
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="API Key">
              <el-input
                v-model="form.apiKey"
                type="password"
                show-password
                :placeholder="keyPlaceholder"
                clearable
              >
                <template #suffix>
                  <SourceTag :from-db="fromDB.apiKey" />
                </template>
              </el-input>
              <div class="form-tip">
                <el-checkbox v-model="form.clearApiKey" size="small">
                  清除已配置的 Key（回落 yaml）
                </el-checkbox>
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="默认模型">
              <el-input
                v-model="form.model"
                placeholder="glm-5.3"
                clearable
              >
                <template #suffix>
                  <SourceTag :from-db="fromDB.model" />
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item label="单次上限">
              <el-input-number
                v-model="form.maxTokens"
                :min="0"
                :step="1024"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item label="每日配额">
              <el-input-number
                v-model="form.dailyQuota"
                :min="0"
                :step="50"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item label="超时(秒)">
              <el-input-number
                v-model="form.timeout"
                :min="0"
                :step="30"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="文风设定">
          <el-input
            v-model="form.styleHint"
            placeholder="技术博客写作助手，文风简洁准确，避免空洞修饰"
            clearable
            maxlength="120"
          >
            <template #suffix>
              <SourceTag :from-db="fromDB.styleHint" />
            </template>
          </el-input>
          <div class="form-tip">注入 system prompt，影响所有 AI 任务的输出风格</div>
        </el-form-item>

        <el-divider content-position="left">
          多模型切换（可选）
          <SourceTag :from-db="fromDB.models" style="margin-left: 6px" />
        </el-divider>
        <div class="form-tip" style="margin-bottom: 10px">
          配置后编辑器 AI 菜单出现模型选择；name 为显示名，baseURL / apiKey /
          maxTokens 留空则继承上方默认配置
        </div>

        <div v-for="(m, i) in form.models" :key="i" class="model-row">
          <el-input
            v-model="m.name"
            placeholder="名称（如 DeepSeek）"
            class="model-cell"
          />
          <el-input
            v-model="m.model"
            placeholder="模型名（如 deepseek-chat）"
            class="model-cell"
          />
          <el-input
            v-model="m.baseURL"
            placeholder="接入点（留空继承）"
            class="model-cell wide"
          />
          <el-input
            v-model="m.apiKey"
            type="password"
            show-password
            placeholder="Key（留空继承）"
            class="model-cell"
          />
          <el-button
            type="danger"
            :icon="Delete"
            circle
            plain
            @click="form.models.splice(i, 1)"
          />
        </div>
        <el-button
          :icon="Plus"
          plain
          style="margin-top: 6px"
          @click="form.models.push({ name: '', model: '' })"
        >
          添加模型
        </el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { Delete, Plus } from "@element-plus/icons-vue";
import { message } from "@/utils/message";
import {
  getAIConfig,
  updateAIConfig,
  type AIConfigInfo
} from "@/api/ai";

defineOptions({ name: "AIConfig" });

const saving = ref(false);
const fromDB = ref<Record<string, boolean>>({});
const keyPlaceholder = ref("未配置");

const form = reactive({
  baseURL: "",
  apiKey: "",
  clearApiKey: false,
  model: "",
  maxTokens: 0,
  dailyQuota: 0,
  styleHint: "",
  timeout: 0,
  models: [] as Array<{
    name: string;
    model: string;
    baseURL?: string;
    apiKey?: string;
    maxTokens?: number;
  }>
});

// 来源标记小组件（覆盖/基线），内联定义避免单文件多组件注册
const SourceTag = defineComponent({
  props: { fromDb: Boolean },
  setup(props) {
    return () =>
      h(
        "span",
        {
          class: "source-tag",
          title: props.fromDb ? "来自后台覆盖" : "来自 config.yaml 基线"
        },
        props.fromDb ? "覆盖" : "基线"
      );
  }
});

async function load() {
  try {
    const res = await getAIConfig();
    if (res.code === 200) {
      const info = res.payload as AIConfigInfo;
      form.baseURL = info.baseURL || "";
      form.model = info.model || "";
      form.maxTokens = info.maxTokens || 0;
      form.dailyQuota = info.dailyQuota || 0;
      form.styleHint = info.styleHint || "";
      form.timeout = info.timeout || 0;
      form.models = (info.models || [])
        .filter(m => m.name !== "default")
        .map(m => ({ ...m }));
      form.apiKey = "";
      form.clearApiKey = false;
      fromDB.value = info.fromDB || {};
      keyPlaceholder.value = info.apiKeySet
        ? `已配置：${info.apiKeyMasked}（留空保持不变）`
        : "未配置";
    }
  } catch {
    /* 加载失败静默，保存时会再暴露 */
  }
}

async function save() {
  // 清除勾选与填写新值互斥提示
  if (form.clearApiKey && form.apiKey) {
    message("已勾选清除，请清空 Key 输入框", { type: "warning" });
    return;
  }
  saving.value = true;
  try {
    const res = await updateAIConfig({
      baseURL: form.baseURL.trim(),
      apiKey: form.apiKey.trim(),
      clearApiKey: form.clearApiKey,
      model: form.model.trim(),
      maxTokens: form.maxTokens,
      dailyQuota: form.dailyQuota,
      styleHint: form.styleHint.trim(),
      timeout: form.timeout,
      models: form.models.filter(m => m.name.trim() && m.model.trim())
    });
    if (res.code === 200) {
      message("已保存并即时生效", { type: "success" });
      await load();
    }
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<script lang="ts">
import { defineComponent, h } from "vue";
export default { name: "AIConfigPage" };
</script>

<style scoped>
.ai-config-card {
  max-width: 1100px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-weight: 600;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.6;
  margin-top: 4px;
}

.model-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.model-cell {
  width: 180px;
}

.model-cell.wide {
  flex: 1;
}

.source-tag {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  padding: 0 4px;
  white-space: nowrap;
}
</style>
