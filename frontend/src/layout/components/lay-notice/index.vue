<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { getNoticeList, markNoticeRead, markAllNoticeRead, type NoticeTab } from "@/api/notice";
import { useWebSocket } from "@/utils/websocket";
import NoticeList from "./components/NoticeList.vue";
import BellIcon from "~icons/ep/bell";
import { message } from "@/utils/message";

const noticesNum = ref(0);
const notices = ref<NoticeTab[]>([]);
const activeKey = ref("1");
const loading = ref(false);

// WebSocket
const { connect, on, off } = useWebSocket();

// 获取通知数据
const fetchNotices = async () => {
  loading.value = true;
  try {
    const res = await getNoticeList();
    if (res.code === 200 && res.payload) {
      notices.value = res.payload.notices || [];
      noticesNum.value = res.payload.unreadCount || 0;
      if (notices.value.length > 0) {
        activeKey.value = notices.value[0].key;
      }
    }
  } catch (error) {
    console.error("获取通知失败:", error);
  } finally {
    loading.value = false;
  }
};

// 处理新通知推送
const handleNewNotice = (data: any) => {
  // 刷新通知列表
  fetchNotices();
  message("收到新通知: " + (data.title || ""), { type: "info" });
};

// 标记已读
const handleMarkAsRead = async (id: number) => {
  try {
    await markNoticeRead(id);
    fetchNotices();
  } catch (error) {
    message("标记失败", { type: "error" });
  }
};

// 全部已读
const handleMarkAllAsRead = async () => {
  try {
    await markAllNoticeRead();
    fetchNotices();
    message("已全部标为已读", { type: "success" });
  } catch (error) {
    message("操作失败", { type: "error" });
  }
};

const getLabel = computed(
  () => (item: NoticeTab) =>
    item.name + (item.list.length > 0 ? `(${item.list.length})` : "")
);

onMounted(() => {
  fetchNotices();
  // 连接 WebSocket 并监听通知
  connect();
  on("notice", handleNewNotice);
  on("system_notice", handleNewNotice);
});

onUnmounted(() => {
  off("notice", handleNewNotice);
  off("system_notice", handleNewNotice);
});
</script>

<template>
  <el-dropdown trigger="click" placement="bottom-end">
    <span
      :class="[
        'dropdown-badge',
        'navbar-bg-hover',
        'select-none',
        Number(noticesNum) !== 0 && 'mr-[10px]'
      ]"
    >
      <el-badge :value="Number(noticesNum) === 0 ? '' : noticesNum" :max="99">
        <span class="header-notice-icon">
          <IconifyIconOffline :icon="BellIcon" />
        </span>
      </el-badge>
    </span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-tabs
          v-model="activeKey"
          :stretch="true"
          class="dropdown-tabs"
          :style="{ width: notices.length === 0 ? '200px' : '330px' }"
        >
          <el-empty
            v-if="loading"
            description="加载中..."
            :image-size="60"
          />
          <el-empty
            v-else-if="notices.length === 0"
            description="暂无消息"
            :image-size="60"
          />
          <span v-else>
            <template v-for="item in notices" :key="item.key">
              <el-tab-pane :label="getLabel(item)" :name="`${item.key}`">
                <el-scrollbar max-height="330px">
                  <div class="noticeList-container">
                    <NoticeList :list="item.list" :emptyText="item.emptyText" />
                  </div>
                </el-scrollbar>
              </el-tab-pane>
            </template>
            <div class="notice-actions">
              <el-button link type="primary" @click="handleMarkAllAsRead">
                全部已读
              </el-button>
            </div>
          </span>
        </el-tabs>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<style lang="scss" scoped>
.dropdown-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 48px;
  cursor: pointer;

  .header-notice-icon {
    font-size: 18px;
  }
}

.dropdown-tabs {
  .noticeList-container {
    padding: 15px 24px 0;
  }

  .notice-actions {
    display: flex;
    justify-content: center;
    padding: 10px 0;
    border-top: 1px solid var(--el-border-color-light);
  }

  :deep(.el-tabs__header) {
    margin: 0;
  }

  :deep(.el-tabs__nav-wrap)::after {
    height: 1px;
  }

  :deep(.el-tabs__nav-wrap) {
    padding: 0 36px;
  }
}
</style>
