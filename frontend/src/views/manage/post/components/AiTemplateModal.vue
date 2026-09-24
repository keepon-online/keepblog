<template>
  <el-dialog
    :model-value="visible"
    title="AI 指令模板管理"
    width="620px"
    append-to-body
    @close="$emit('update:visible', false)"
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
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  getAITemplates,
  saveAITemplate,
  type AITemplateItem
} from "@/api/ai";
import { message } from "@/utils/message";

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  (e: "update:visible", val: boolean): void;
}>();

const aiTplItems = ref<AITemplateItem[]>([]);
const aiTplKey = ref("polish:polish");
const aiTplContent = ref("");
const aiTplSaving = ref(false);

const aiTplLabelMap: Record<string, string> = {
  "polish:polish": "润色",
  "polish:simplify": "精简",
  "polish:expand": "扩写",
  "polish:academic": "学术化",
  "polish:casual": "口语化",
  "title:title": "生成标题",
  "summary:summary": "提炼摘要",
  "tags:tags": "推荐标签",
  "outline:outline": "生成大纲",
  "continue:continue": "智能续写",
  "proofread:proofread": "全文校对"
};

const aiTplCurrent = computed(
  () => aiTplItems.value.find(t => t.key === aiTplKey.value) ?? null
);

const fetchTemplates = async () => {
  try {
    const res = await getAITemplates();
    if (res.code === 200) {
      aiTplItems.value = res.payload ?? [];
      const cur = aiTplItems.value.find(t => t.key === aiTplKey.value);
      aiTplContent.value = cur?.content ?? "";
    }
  } catch {
    // ignore
  }
};

watch(
  () => props.visible,
  val => {
    if (val) {
      fetchTemplates();
    }
  }
);

const onTplKeyChange = () => {
  aiTplContent.value = aiTplCurrent.value?.content ?? "";
};

const saveTpl = async (reset = false) => {
  aiTplSaving.value = true;
  try {
    const content = reset ? "" : aiTplContent.value;
    const res = await saveAITemplate(aiTplKey.value, content);
    if (res.code === 200) {
      const cur = aiTplItems.value.find(t => t.key === aiTplKey.value);
      if (cur) cur.content = reset ? "" : aiTplContent.value.trim();
      if (reset) aiTplContent.value = "";
      message(reset ? "已恢复默认指令" : "指令模板保存成功", {
        type: "success"
      });
    } else {
      message(res.message || "保存失败", { type: "error" });
    }
  } catch (e: any) {
    message(e?.message || "网络异常", { type: "error" });
  } finally {
    aiTplSaving.value = false;
  }
};
</script>

<style scoped lang="scss">
.ai-tpl-default {
  margin-top: 12px;
  padding: 10px 12px;
  background: var(--el-fill-color-light);
  border-radius: 4px;
  font-size: 12px;

  .ai-tpl-default-title {
    color: var(--el-text-color-secondary);
    margin-bottom: 4px;
    font-weight: 500;
  }

  .ai-tpl-default-text {
    color: var(--el-text-color-regular);
    white-space: pre-wrap;
    line-height: 1.5;
  }
}
</style>
