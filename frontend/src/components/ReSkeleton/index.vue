<script setup lang="ts">
import { computed } from "vue";

interface Props {
  loading?: boolean;
  rows?: number;
  animated?: boolean;
  throttle?: number;
  variant?: "card" | "table" | "list" | "chart";
}

const props = withDefaults(defineProps<Props>(), {
  loading: true,
  rows: 4,
  animated: true,
  throttle: 0,
  variant: "card"
});

// 根据变体类型计算行高配置
const rowConfig = computed(() => {
  switch (props.variant) {
    case "chart":
      return { rows: 1, height: "300px" };
    case "table":
      return { rows: props.rows, height: "40px" };
    case "list":
      return { rows: props.rows, height: "60px" };
    case "card":
    default:
      return { rows: props.rows, height: "20px" };
  }
});
</script>

<template>
  <el-skeleton
    :loading="loading"
    :animated="animated"
    :throttle="throttle"
    :rows="rowConfig.rows"
    class="re-skeleton"
    :class="[`re-skeleton--${variant}`]"
  >
    <template #template>
      <!-- 卡片变体 -->
      <template v-if="variant === 'card'">
        <el-skeleton-item variant="h3" style="width: 30%; margin-bottom: 16px" />
        <el-skeleton-item
          v-for="i in rows"
          :key="i"
          variant="text"
          :style="{ width: i === rows ? '60%' : '100%', marginBottom: '12px' }"
        />
      </template>

      <!-- 图表变体 -->
      <template v-else-if="variant === 'chart'">
        <div class="chart-skeleton">
          <el-skeleton-item variant="rect" style="width: 100%; height: 300px" />
        </div>
      </template>

      <!-- 表格变体 -->
      <template v-else-if="variant === 'table'">
        <div class="table-skeleton">
          <div class="table-skeleton__header">
            <el-skeleton-item
              v-for="i in 5"
              :key="i"
              variant="text"
              style="width: 18%; height: 40px"
            />
          </div>
          <div
            v-for="row in rows"
            :key="row"
            class="table-skeleton__row"
          >
            <el-skeleton-item
              v-for="col in 5"
              :key="col"
              variant="text"
              style="width: 18%; height: 32px"
            />
          </div>
        </div>
      </template>

      <!-- 列表变体 -->
      <template v-else-if="variant === 'list'">
        <div
          v-for="i in rows"
          :key="i"
          class="list-skeleton__item"
        >
          <el-skeleton-item variant="circle" style="width: 48px; height: 48px" />
          <div class="list-skeleton__content">
            <el-skeleton-item variant="h3" style="width: 30%" />
            <el-skeleton-item variant="text" style="width: 80%; margin-top: 8px" />
          </div>
        </div>
      </template>
    </template>

    <!-- 默认插槽：实际内容 -->
    <template #default>
      <slot />
    </template>
  </el-skeleton>
</template>

<style scoped lang="scss">
.re-skeleton {
  width: 100%;

  &--chart {
    .chart-skeleton {
      padding: 16px;
    }
  }

  &--table {
    .table-skeleton {
      &__header {
        display: flex;
        gap: 12px;
        padding: 12px 16px;
        background: var(--el-fill-color-light);
        border-radius: 4px 4px 0 0;
      }

      &__row {
        display: flex;
        gap: 12px;
        padding: 12px 16px;
        border-bottom: 1px solid var(--el-border-color-lighter);
      }
    }
  }

  &--list {
    .list-skeleton__item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 12px 0;
      border-bottom: 1px solid var(--el-border-color-lighter);

      &:last-child {
        border-bottom: none;
      }
    }

    .list-skeleton__content {
      flex: 1;
    }
  }
}
</style>
