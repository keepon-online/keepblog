<script setup lang="ts">
defineOptions({
  name: "Music"
});

import { useMusic } from "./hook";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { Delete, Edit, Plus, VideoPlay, VideoPause, Close } from "@element-plus/icons-vue";

const {
  dialogFormVisible,
  loading,
  columns,
  dataList,
  form,
  rules,
  ruleFormRef,
  title,
  currentPlaying,
  isPlaying,
  togglePlay,
  stopPlay,
  submitForm,
  handleAdd,
  onSearch,
  handleUpdate,
  handleDelete
} = useMusic();
</script>

<template>
  <div class="music-manage-container">
    <PureTableBar title="博客背景音乐媒体库" :columns="columns" @refresh="onSearch">
      <template #buttons>
        <el-button
          type="primary"
          :icon="useRenderIcon(Plus)"
          @click="handleAdd"
        >
          添加音乐
        </el-button>
      </template>

      <template v-slot="{ size, dynamicColumns }">
        <el-card shadow="never" class="table-card">
          <pure-table
            border
            align-whole="center"
            showOverflowTooltip
            table-layout="auto"
            :loading="loading"
            :size="size"
            :data="dataList"
            :columns="dynamicColumns"
            :paginationSmall="size === 'small' ? true : false"
            :header-cell-style="{
              background: 'var(--el-table-row-hover-bg-color)',
              color: 'var(--el-text-color-primary)'
            }"
          >
            <template #operation="{ row }">
              <el-button
                class="reset-margin"
                link
                type="primary"
                :size="size"
                :icon="useRenderIcon(Edit)"
                @click="handleUpdate(row)"
              >
                编辑
              </el-button>
              <el-popconfirm
                :title="`确认删除歌曲【${row.name}】?`"
                @confirm="handleDelete(row)"
              >
                <template #reference>
                  <el-button
                    class="reset-margin"
                    link
                    type="danger"
                    :size="size"
                    :icon="useRenderIcon(Delete)"
                  >
                    删除
                  </el-button>
                </template>
              </el-popconfirm>
            </template>
          </pure-table>
        </el-card>
      </template>
    </PureTableBar>

    <!-- 底部悬浮试听控制器 -->
    <transition name="el-zoom-in-bottom">
      <div v-if="currentPlaying" class="floating-music-player">
        <div class="disc-cover" :class="{ 'is-spinning': isPlaying }">
          <img v-if="currentPlaying.cover" :src="currentPlaying.cover" alt="cover" />
          <div v-else class="disc-placeholder">🎵</div>
        </div>
        <div class="player-info">
          <div class="song-title">{{ currentPlaying.name }}</div>
          <div class="song-artist">{{ currentPlaying.artist }}</div>
        </div>
        <div class="player-controls">
          <el-button
            circle
            type="primary"
            :icon="useRenderIcon(isPlaying ? VideoPause : VideoPlay)"
            @click="togglePlay(currentPlaying)"
          />
          <el-button
            circle
            :icon="useRenderIcon(Close)"
            @click="stopPlay"
          />
        </div>
      </div>
    </transition>

    <!-- 添加/编辑弹窗 -->
    <el-dialog
      v-model="dialogFormVisible"
      :title="title"
      width="560px"
      destroy-on-close
    >
      <el-form
        ref="ruleFormRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-row :gutter="16">
          <el-col :span="14">
            <el-form-item label="歌曲名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入歌曲名称" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="10">
            <el-form-item label="艺术家 / 歌手" prop="artist">
              <el-input v-model="form.artist" placeholder="请输入歌手" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="音频文件直链 (URL)" prop="url">
          <el-input
            v-model="form.url"
            placeholder="例如: https://example.com/song.mp3"
            clearable
          />
        </el-form-item>

        <el-form-item label="专辑封面图片 (Cover URL)" prop="cover">
          <el-input
            v-model="form.cover"
            placeholder="例如: https://example.com/cover.jpg"
            clearable
          />
        </el-form-item>

        <el-form-item label="歌词内容 (LRC)" prop="lrc">
          <el-input
            v-model="form.lrc"
            type="textarea"
            :rows="3"
            placeholder="[00:00.00] 歌词内容..."
          />
        </el-form-item>

        <el-form-item label="排序权重" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>

        <!-- 封面实时预览 -->
        <div v-if="form.name || form.url" class="album-preview">
          <div class="album-cover">
            <img v-if="form.cover" :src="form.cover" alt="cover" />
            <div v-else class="cover-fallback">🎵</div>
          </div>
          <div class="album-details">
            <div class="album-name">{{ form.name || "未命名歌曲" }}</div>
            <div class="album-artist">{{ form.artist || "未知歌手" }}</div>
            <div class="album-url">{{ form.url || "暂未填写音频直链" }}</div>
          </div>
        </div>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogFormVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm(ruleFormRef)">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
.music-manage-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 120px);
  position: relative;

  .table-card {
    border-radius: 10px;
    border: 1px solid var(--el-border-color-lighter);

    :deep(.el-card__body) {
      padding: 0;
    }
  }

  .floating-music-player {
    position: fixed;
    bottom: 24px;
    right: 32px;
    z-index: 999;
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 10px 18px;
    background: var(--el-bg-color-overlay);
    border: 1px solid var(--el-border-color-light);
    border-radius: 36px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.16);
    backdrop-filter: blur(12px);

    .disc-cover {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      overflow: hidden;
      border: 2px solid var(--el-color-primary-light-5);
      flex-shrink: 0;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .disc-placeholder {
        width: 100%;
        height: 100%;
        background: linear-gradient(135deg, #9b59b6, #8e44ad);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 16px;
      }

      &.is-spinning {
        animation: spin 6s linear infinite;
      }
    }

    .player-info {
      max-width: 160px;

      .song-title {
        font-size: 13px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .song-artist {
        font-size: 11px;
        color: var(--el-text-color-secondary);
        margin-top: 2px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .player-controls {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
}

.album-preview {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px;
  background: var(--el-fill-color-light);
  border-radius: 8px;
  border: 1px dashed var(--el-border-color);
  margin-top: 14px;

  .album-cover {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .cover-fallback {
      width: 100%;
      height: 100%;
      background: linear-gradient(135deg, #9b59b6, #8e44ad);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 20px;
    }
  }

  .album-details {
    flex: 1;
    min-width: 0;

    .album-name {
      font-size: 14px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    .album-artist {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-top: 2px;
    }

    .album-url {
      font-size: 11px;
      color: var(--el-color-primary);
      margin-top: 2px;
      font-family: monospace;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
