<template>
  <el-config-provider :locale="currentLocale">
    <router-view />
    <ReDialog />
  </el-config-provider>
</template>

<script lang="ts">
import { useRouter } from "vue-router";
import { defineComponent, computed } from "vue";
import { ElConfigProvider } from "element-plus";
import { ReDialog, closeAllDialog } from "@/components/ReDialog";
import zhCn from "element-plus/es/locale/lang/zh-cn";

export default defineComponent({
  name: "app",
  components: {
    [ElConfigProvider.name]: ElConfigProvider,
    ReDialog
  },
  setup() {
    const router = useRouter();
    const currentLocale = computed(() => {
      return zhCn;
    });
    // 切换路由时关闭残留弹层（https://github.com/pure-admin/vue-pure-admin/commit/815f161b）
    router.beforeEach(() => {
      closeAllDialog();
    });
    return {
      currentLocale
    };
  }
});
</script>
