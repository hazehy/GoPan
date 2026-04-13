<template>
  <div v-if="props.visible" class="dialog-mask" @click="props.closeMoveDialog">
    <div class="dialog-panel move-dialog-panel" @click.stop>
      <h3 class="dialog-title">移动到文件夹</h3>
      <p class="dialog-message">
        文件：{{
          props.movingItem ? `${props.movingItem.name}${props.movingItem.ext ?? ""}` : "-"
        }}
      </p>

      <div class="path-bar" style="margin-bottom: 8px">
        <a href="javascript:void(0)" class="path-link" @click="props.goMoveRoot"
          >根目录</a
        >
        <template
          v-for="(node, index) in props.moveBreadcrumbs"
          :key="node.identity"
        >
          <span class="path-sep">/</span>
          <a
            href="javascript:void(0)"
            class="path-link"
            @click="props.goMovePath(index)"
          >
            {{ node.name }}
          </a>
        </template>
      </div>

      <div class="move-folder-list" v-if="props.moveFolderOptions.length">
        <button
          v-for="folder in props.moveFolderOptions"
          :key="folder.identity"
          class="move-folder-item"
          type="button"
          :disabled="props.isMoveFolderSelf(folder)"
          @click="props.enterMoveFolder(folder)"
        >
          📁 {{ folder.name }}
        </button>
      </div>
      <p class="muted" v-else>当前目录没有子文件夹</p>

      <p class="muted" style="margin: 8px 0 0">
        目标目录：{{ props.currentMovePath }}
      </p>
      <p class="error" v-if="props.moveErrorMessage">{{ props.moveErrorMessage }}</p>

      <div class="dialog-actions" style="margin-top: 12px">
        <button
          class="btn btn-secondary"
          type="button"
          @click="props.closeMoveDialog"
        >
          取消
        </button>
        <button
          class="btn btn-primary"
          type="button"
          :disabled="props.moveLoading || !props.canConfirmMove"
          @click="props.confirmMove"
        >
          {{ props.moveLoading ? "移动中..." : "确认移动" }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UserFile } from "@/types/api";

interface Props {
  visible: boolean;
  movingItem: UserFile | null;
  moveBreadcrumbs: Array<{ id: number; identity: string; name: string }>;
  moveFolderOptions: UserFile[];
  currentMovePath: string;
  moveErrorMessage: string;
  moveLoading: boolean;
  canConfirmMove: boolean;
  closeMoveDialog: () => void;
  goMoveRoot: () => void;
  goMovePath: (index: number) => void;
  enterMoveFolder: (folder: UserFile) => void;
  confirmMove: () => void;
  isMoveFolderSelf: (folder: UserFile) => boolean;
}

const props = defineProps<Props>();
</script>
