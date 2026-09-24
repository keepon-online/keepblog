<template>
  <div
    v-if="visible"
    ref="menuRef"
    class="slash-command-menu"
    :class="[`placement-${resolvedPlacement}`]"
    :style="menuStyle"
    @click.stop
  >
    <div class="menu-header">
      <span class="menu-title">⚡ 快速命令 (Slash Command)</span>
      <span class="menu-tip">↑↓ 切换 · 回车确认 · Esc 退出</span>
    </div>

    <div
      ref="menuListRef"
      class="menu-list"
      :style="{ maxHeight: maxListHeight }"
    >
      <div
        v-for="(cmd, index) in filteredCommands"
        :key="cmd.id"
        class="command-item"
        :class="{ active: selectedIndex === index }"
        @click="executeCommand(cmd)"
        @mouseenter="selectedIndex = index"
      >
        <span class="cmd-icon">{{ cmd.icon }}</span>
        <div class="cmd-content">
          <span class="cmd-name">{{ cmd.label }}</span>
          <span class="cmd-desc">{{ cmd.desc }}</span>
        </div>
        <span class="cmd-badge">{{ cmd.slash }}</span>
      </div>
      <div v-if="filteredCommands.length === 0" class="no-matches">
        未找到匹配命令
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from "vue";

export interface SlashCommandItem {
  id: string;
  label: string;
  desc: string;
  slash: string;
  icon: string;
  template?: string;
  action?: string;
}

export interface SlashCommandPosition {
  top: number;
  left: number;
  placement?: "top" | "bottom";
}

const props = withDefaults(
  defineProps<{
    visible: boolean;
    position: SlashCommandPosition;
    filterText?: string;
  }>(),
  {
    filterText: ""
  }
);

const emit = defineEmits<{
  (e: "select", cmd: SlashCommandItem): void;
  (e: "close"): void;
}>();

const menuRef = ref<HTMLElement | null>(null);
const menuListRef = ref<HTMLElement | null>(null);
const selectedIndex = ref(0);

const COMMANDS: SlashCommandItem[] = [
  {
    id: "ai",
    label: "AI 智能续写",
    desc: "根据前文逻辑流式生成下一自然段",
    slash: "/ai",
    icon: "✨",
    action: "ai-continue"
  },
  {
    id: "outline",
    label: "AI 文章大纲",
    desc: "依据主题与受众定位智能生成多级技术大纲",
    slash: "/outline",
    icon: "📑",
    action: "ai-outline"
  },
  {
    id: "code",
    label: "代码块 (Code Block)",
    desc: "插入带语言高亮的代码围栏",
    slash: "/code",
    icon: "💻",
    template: "```go\n// 在此编写代码\n\n```\n"
  },
  {
    id: "table",
    label: "数据表格 (Table)",
    desc: "插入 3x3 标准 Markdown 数据表格",
    slash: "/table",
    icon: "📊",
    template: "| 参数名 | 类型 | 必填 | 说明 |\n| :--- | :--- | :---: | :--- |\n| name | string | 是 | 资源名称 |\n| count | int | 否 | 数量限制 |\n"
  },
  {
    id: "mermaid",
    label: "Mermaid 流程图",
    desc: "插入 Mermaid 图表架构图模板",
    slash: "/flow",
    icon: "🗺️",
    template: "```mermaid\nflowchart TD\n    A[客户端请求] --> B[API 网关]\n    B --> C{认证校验}\n    C -->|通过| D[微服务集群]\n    C -->|拒绝| E[返回 401]\n```\n"
  },
  {
    id: "tip",
    label: "提示块 (Tip Alert)",
    desc: "插入 GitHub 风格绿色最佳实践提示块",
    slash: "/tip",
    icon: "💡",
    template: "> [!TIP]\n> 这是一个高亮提示块，用于强调关键配置与最佳实践。\n\n"
  },
  {
    id: "warning",
    label: "警告块 (Warning Alert)",
    desc: "插入 GitHub 风格黄色/橙色避坑警示",
    slash: "/warn",
    icon: "⚠️",
    template: "> [!WARNING]\n> 生产环境操作需注意并发竞争与死锁风险。\n\n"
  },
  {
    id: "note",
    label: "说明块 (Note Alert)",
    desc: "插入 GitHub 风格蓝色背景补充说明",
    slash: "/note",
    icon: "📝",
    template: "> [!NOTE]\n> 本章节仅适用于 Linux 内核 5.4 及以上版本。\n\n"
  },
  {
    id: "math",
    label: "数学公式块 (KaTeX)",
    desc: "插入独立 KaTeX 块级数学公式",
    slash: "/math",
    icon: "🔢",
    template: "$$\nE = mc^2\n$$\n"
  }
];

const filteredCommands = computed(() => {
  const query = (props.filterText || "").trim().toLowerCase();
  if (!query) return COMMANDS;
  return COMMANDS.filter(
    c =>
      c.slash.toLowerCase().includes(query) ||
      c.label.toLowerCase().includes(query) ||
      c.id.toLowerCase().includes(query)
  );
});

const resolvedPlacement = computed(() => {
  return props.position?.placement || "bottom";
});

const menuStyle = computed(() => {
  const isTop = resolvedPlacement.value === "top";
  return {
    top: `${props.position.top}px`,
    left: `${props.position.left}px`,
    transform: isTop ? "translateY(-100%)" : "translateY(0)"
  };
});

const maxListHeight = computed(() => {
  if (typeof window === "undefined") return "280px";
  const vh = window.innerHeight;
  const isTop = resolvedPlacement.value === "top";
  const topPos = props.position?.top ?? 300;

  if (isTop) {
    // 向上弹出时，可用高度是 topPos - 头部高度(约50px) - 视口安全边距(16px)
    const available = topPos - 66;
    return `${Math.max(120, Math.min(280, available))}px`;
  } else {
    // 向下弹出时，可用高度是 视口高度 - topPos - 头部高度(约50px) - 视口安全边距(16px)
    const available = vh - topPos - 66;
    return `${Math.max(120, Math.min(280, available))}px`;
  }
});

watch(
  () => props.filterText,
  () => {
    selectedIndex.value = 0;
  }
);

watch(
  () => props.visible,
  val => {
    if (val) {
      selectedIndex.value = 0;
    }
  }
);

watch(selectedIndex, async () => {
  await nextTick();
  if (!menuListRef.value) return;
  const activeEl = menuListRef.value.querySelector(".command-item.active") as HTMLElement;
  if (activeEl) {
    activeEl.scrollIntoView({ block: "nearest" });
  }
});

const executeCommand = (cmd: SlashCommandItem) => {
  emit("select", cmd);
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (!props.visible) return;

  if (e.key === "ArrowDown") {
    e.preventDefault();
    e.stopPropagation();
    if (filteredCommands.value.length > 0) {
      selectedIndex.value = (selectedIndex.value + 1) % filteredCommands.value.length;
    }
  } else if (e.key === "ArrowUp") {
    e.preventDefault();
    e.stopPropagation();
    if (filteredCommands.value.length > 0) {
      selectedIndex.value =
        (selectedIndex.value - 1 + filteredCommands.value.length) % filteredCommands.value.length;
    }
  } else if (e.key === "Enter") {
    e.preventDefault();
    e.stopPropagation();
    const cmd = filteredCommands.value[selectedIndex.value];
    if (cmd) {
      executeCommand(cmd);
    }
  } else if (e.key === "Escape") {
    e.preventDefault();
    e.stopPropagation();
    emit("close");
  }
};

const handlePointerDown = (e: PointerEvent) => {
  if (!props.visible) return;
  if (menuRef.value && !menuRef.value.contains(e.target as Node)) {
    emit("close");
  }
};

onMounted(() => {
  window.addEventListener("keydown", handleKeyDown, true);
  window.addEventListener("pointerdown", handlePointerDown, true);
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", handleKeyDown, true);
  window.removeEventListener("pointerdown", handlePointerDown, true);
});
</script>

<style scoped lang="scss">
.slash-command-menu {
  position: fixed;
  z-index: 2400;
  width: 320px;
  background: var(--el-bg-color-overlay, #ffffff);
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.16);
  padding: 8px;
  display: flex;
  flex-direction: column;

  &.placement-top {
    animation: slashFadeInTop 0.15s cubic-bezier(0.16, 1, 0.3, 1);
  }

  &.placement-bottom {
    animation: slashFadeInBottom 0.15s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 8px 8px;
    border-bottom: 1px solid var(--el-border-color-extra-light);

    .menu-title {
      font-size: 12px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    .menu-tip {
      font-size: 10px;
      color: var(--el-text-color-placeholder);
    }
  }

  .menu-list {
    max-height: 280px;
    overflow-y: auto;
    padding: 4px 0;
    overscroll-behavior: contain;
  }

  .command-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 10px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.15s ease;

    &:hover,
    &.active {
      background: var(--el-fill-color-light);
    }

    &.active {
      background: var(--el-color-primary-light-9);
      .cmd-name {
        color: var(--el-color-primary);
      }
    }

    .cmd-icon {
      font-size: 18px;
      width: 24px;
      text-align: center;
      flex-shrink: 0;
    }

    .cmd-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;

      .cmd-name {
        font-size: 13px;
        font-weight: 500;
        color: var(--el-text-color-primary);
      }

      .cmd-desc {
        font-size: 11px;
        color: var(--el-text-color-secondary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .cmd-badge {
      font-size: 11px;
      font-family: ui-monospace, SFMono-Regular, monospace;
      color: var(--el-text-color-placeholder);
      background: var(--el-fill-color);
      padding: 1px 6px;
      border-radius: 4px;
    }
  }

  .no-matches {
    padding: 16px;
    text-align: center;
    font-size: 12px;
    color: var(--el-text-color-placeholder);
  }
}

@keyframes slashFadeInBottom {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slashFadeInTop {
  from {
    opacity: 0;
    transform: translateY(calc(-100% - 6px));
  }
  to {
    opacity: 1;
    transform: translateY(-100%);
  }
}
</style>
