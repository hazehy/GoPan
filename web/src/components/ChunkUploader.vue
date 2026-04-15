<template>
  <div class="card uploader-card">
    <input
      ref="fileInputRef"
      class="uploader-native-input"
      type="file"
      multiple
      @change="onFileSelect"
    />
    <input
      ref="bindFileInputRef"
      class="uploader-native-input"
      type="file"
      @change="onBindFileSelect"
    />
    <div class="uploader-panel x-scroll-panel uploader-scroll-panel">
      <div class="uploader-panel-inner">
      <div class="uploader-top">
        <div class="uploader-file-meta">
          <p class="uploader-caption">文件上传</p>
          <p class="uploader-file-name">
            {{
              uploadingTask
                ? `正在上传：${uploadingTask.name}${uploadingTask.ext ?? ""}`
                : "当前无进行中的上传任务"
            }}
          </p>
          <div class="uploader-summary">
            <span class="uploader-chip">总任务 {{ taskCount }}</span>
            <span class="uploader-chip">等待 {{ waitingCount }}</span>
            <span class="uploader-chip">进行中 {{ uploadingTask ? 1 : 0 }}</span>
            <span class="uploader-chip">完成 {{ completedCount }}</span>
          </div>
        </div>

        <div class="uploader-actions">
          <button class="btn btn-primary" @click="chooseFile">选择文件</button>
          <button class="btn btn-secondary" @click="startQueue" :disabled="!taskCount">
            启动队列
          </button>
          <button class="btn btn-secondary" @click="stopQueue" :disabled="!taskCount">
            停止全部
          </button>
          <button
            class="btn btn-secondary"
            @click="clearFinishedTasks"
            :disabled="!hasCompletedTasks"
          >
            清理已完成
          </button>
        </div>
      </div>

      <div v-if="!sortedTasks.length" class="uploader-empty muted">
        暂无上传任务，点击“选择文件”后将进入队列。
      </div>

      <div class="uploader-queue" v-else>
        <div class="uploader-task" v-for="task in sortedTasks" :key="task.id">
          <div class="uploader-task-head">
            <div class="uploader-task-name" :title="`${task.name}${task.ext ?? ''}`">
              {{ task.name }}{{ task.ext ?? "" }}
            </div>
            <span class="uploader-state" :class="stateClassMap[task.status]">
              {{ renderStatusText(task.status) }}
            </span>
          </div>

          <div class="uploader-progress">
            <div class="uploader-progress-bar" :style="{ width: `${task.progress}%` }"></div>
          </div>

          <div class="uploader-task-info">
            <div class="uploader-row">
              <span class="muted">进度 {{ task.progress }}%</span>
              <span class="muted"
                >大小 {{ formatSize(task.uploadedBytes) }} / {{ formatSize(task.size) }}</span
              >
              <span class="muted">速度 {{ formatSpeed(task.speedBytesPerSec) }}</span>
            </div>
            <div class="uploader-row" v-if="!task.hasLocalFile">
              <span class="uploader-tip">需绑定原文件后继续</span>
            </div>
          </div>

          <p class="error uploader-task-error" v-if="task.errorMessage">
            {{ task.errorMessage }}
          </p>

          <div class="uploader-task-actions">
            <button
              class="btn btn-secondary"
              type="button"
              @click="pauseTask(task.id)"
              :disabled="task.status !== 'uploading' && task.status !== 'waiting'"
            >
              暂停
            </button>
            <button
              class="btn btn-secondary"
              type="button"
              @click="resumeTask(task.id)"
              :disabled="!canResume(task.status) || !task.hasLocalFile"
            >
              继续
            </button>
            <button
              class="btn btn-secondary"
              type="button"
              @click="triggerBindFile(task.id)"
              :disabled="task.hasLocalFile"
            >
              绑定文件
            </button>
            <button
              class="btn btn-danger"
              type="button"
              @click="removeTask(task.id)"
            >
              移除
            </button>
          </div>
        </div>
      </div>
      </div>
    </div>

    <div
      v-if="locationDialogVisible"
      class="dialog-mask"
      @click="closeLocationDialog"
    >
      <div class="dialog-panel move-dialog-panel" @click.stop>
        <h3 class="dialog-title">选择上传保存位置</h3>
        <p class="dialog-message">请选择本次上传文件保存到哪个目录。</p>

        <div class="uploader-destination-banner">
          <span class="uploader-destination-label">当前浏览目录</span>
          <span class="uploader-destination-path">{{ currentBrowsePath }}</span>
        </div>

        <div class="path-bar" style="margin-bottom: 8px">
          <a href="javascript:void(0)" class="path-link" @click="goLocationRoot"
            >根目录</a
          >
          <template
            v-for="(node, index) in locationBreadcrumbs"
            :key="node.identity"
          >
            <span class="path-sep">/</span>
            <a
              href="javascript:void(0)"
              :class="['path-link', { 'path-link-active': index === locationBreadcrumbs.length - 1 }]"
              @click="goLocationPath(index)"
            >
              {{ node.name }}
            </a>
          </template>
        </div>

        <div class="move-folder-list" v-if="locationFolderOptions.length">
          <button
            v-for="folder in locationFolderOptions"
            :key="folder.identity"
            class="move-folder-item"
            type="button"
            @click="enterLocationFolder(folder)"
          >
            📁 {{ folder.name }}
          </button>
        </div>
        <p class="muted" v-else-if="!locationLoading">当前目录没有子文件夹</p>
        <p class="muted" v-else>目录加载中...</p>

        <p class="uploader-destination-hint">
          将上传到：{{ currentLocationPath }}
        </p>
        <p class="error" v-if="locationErrorMessage">{{ locationErrorMessage }}</p>

        <div class="dialog-actions" style="margin-top: 12px">
          <button class="btn btn-secondary" type="button" @click="closeLocationDialog">
            取消
          </button>
          <button
            class="btn btn-primary"
            type="button"
            :disabled="locationLoading"
            @click="confirmLocationAndChooseFile"
          >
            选择此位置并上传
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { fileListApi } from "@/api/modules/disk";
import type { UserFile } from "@/types/api";
import {
  type UploadTaskStatus,
  useUploadQueueStore,
} from "@/stores/uploadQueue";
import { validateFileOrFolderName } from "@/utils/validators";

const props = defineProps<{
  parentId: number;
  browsePath: Array<{ id: number; name: string }>;
}>();

const emit = defineEmits<{
  success: [];
}>();

const fileInputRef = ref<HTMLInputElement | null>(null);
const bindFileInputRef = ref<HTMLInputElement | null>(null);
const bindingTaskId = ref("");
const locationDialogVisible = ref(false);
const locationLoading = ref(false);
const locationErrorMessage = ref("");
const locationParentId = ref(0);
const selectedUploadParentId = ref(0);
const locationBreadcrumbs = ref<Array<{ id: number; identity: string; name: string }>>(
  [],
);
const locationFolderOptions = ref<UserFile[]>([]);
const uploadQueueStore = useUploadQueueStore();
uploadQueueStore.ensureInitialized();

const { sortedTasks, taskCount, waitingCount, uploadingTask, successVersion } =
  storeToRefs(uploadQueueStore);

const hasCompletedTasks = computed(() =>
  sortedTasks.value.some((item) => item.status === "success"),
);

const completedCount = computed(
  () => sortedTasks.value.filter((item) => item.status === "success").length,
);

const stateClassMap: Record<UploadTaskStatus, string> = {
  waiting: "uploader-state-waiting",
  uploading: "uploader-state-uploading",
  paused: "uploader-state-paused",
  stopped: "uploader-state-stopped",
  success: "uploader-state-success",
  failed: "uploader-state-failed",
};

function chooseFile() {
  locationParentId.value = props.parentId;
  selectedUploadParentId.value = props.parentId;
  locationBreadcrumbs.value = props.browsePath.map((node, index) => ({
    id: node.id,
    name: node.name,
    identity: `${node.id}-${index}`,
  }));
  locationErrorMessage.value = "";
  locationDialogVisible.value = true;
  void loadLocationFolderOptions();
}

function formatSize(size: number): string {
  if (size < 1024) {
    return `${size} B`;
  }
  if (size < 1024 * 1024) {
    return `${(size / 1024).toFixed(1)} KB`;
  }
  return `${(size / (1024 * 1024)).toFixed(2)} MB`;
}

function formatSpeed(speed: number) {
  if (!speed || speed <= 0) {
    return "0 B/s";
  }
  return `${formatSize(speed)}/s`;
}

function renderStatusText(status: UploadTaskStatus) {
  const textMap: Record<UploadTaskStatus, string> = {
    waiting: "排队中",
    uploading: "上传中",
    paused: "已暂停",
    stopped: "已停止",
    success: "已完成",
    failed: "失败",
  };
  return textMap[status];
}

function canResume(status: UploadTaskStatus) {
  return ["paused", "stopped", "failed"].includes(status);
}

function startQueue() {
  uploadQueueStore.startQueue();
}

function stopQueue() {
  uploadQueueStore.stopQueue();
}

function pauseTask(taskId: string) {
  uploadQueueStore.pauseTask(taskId);
}

function resumeTask(taskId: string) {
  uploadQueueStore.resumeTask(taskId);
}

function removeTask(taskId: string) {
  uploadQueueStore.removeTask(taskId);
}

function triggerBindFile(taskId: string) {
  bindingTaskId.value = taskId;
  bindFileInputRef.value?.click();
}

function clearFinishedTasks() {
  uploadQueueStore.clearFinishedTasks();
}

function isFolder(item: UserFile) {
  return !item.repository_identity;
}

async function fetchFolderListByParentId(parentId: number) {
  const chunkSize = 200;
  let currentPage = 1;
  let loaded: UserFile[] = [];
  let totalCount = 0;

  do {
    const response = await fileListApi({
      id: parentId,
      page: currentPage,
      size: chunkSize,
    });

    const list = response.list ?? [];
    loaded = loaded.concat(list);
    totalCount = response.count ?? 0;
    currentPage += 1;

    if (list.length === 0) {
      break;
    }
  } while (loaded.length < totalCount);

  return loaded.filter((item) => isFolder(item));
}

async function loadLocationFolderOptions() {
  try {
    locationLoading.value = true;
    locationErrorMessage.value = "";
    locationFolderOptions.value = await fetchFolderListByParentId(
      locationParentId.value,
    );
  } catch (error) {
    locationErrorMessage.value =
      error instanceof Error ? error.message : String(error);
  } finally {
    locationLoading.value = false;
  }
}

function closeLocationDialog() {
  locationDialogVisible.value = false;
  locationFolderOptions.value = [];
  locationErrorMessage.value = "";
}

async function goLocationRoot() {
  locationParentId.value = 0;
  locationBreadcrumbs.value = [];
  await loadLocationFolderOptions();
}

async function goLocationPath(index: number) {
  locationBreadcrumbs.value = locationBreadcrumbs.value.slice(0, index + 1);
  locationParentId.value = locationBreadcrumbs.value[index].id;
  await loadLocationFolderOptions();
}

async function enterLocationFolder(folder: UserFile) {
  locationBreadcrumbs.value.push({
    id: folder.id,
    identity: folder.identity,
    name: folder.name,
  });
  locationParentId.value = folder.id;
  await loadLocationFolderOptions();
}

function confirmLocationAndChooseFile() {
  selectedUploadParentId.value = locationParentId.value;
  closeLocationDialog();
  fileInputRef.value?.click();
}

const currentLocationPath = computed(() => {
  if (!locationBreadcrumbs.value.length) {
    return "根目录";
  }
  return `根目录 / ${locationBreadcrumbs.value
    .map((item) => item.name)
    .join(" / ")}`;
});

const currentBrowsePath = computed(() => {
  if (!props.browsePath.length) {
    return "根目录";
  }
  return `根目录 / ${props.browsePath.map((item) => item.name).join(" / ")}`;
});

const hasBrowsePath = computed(() => props.browsePath.length > 0);

watch(successVersion, (next, prev) => {
  if (next > prev) {
    emit("success");
  }
});

function onFileSelect(event: Event) {
  const target = event.target as HTMLInputElement;
  const files = Array.from(target.files ?? []);
  if (!files.length) {
    return;
  }

  const validFiles: File[] = [];

  files.forEach((file) => {
    if (file.size <= 0) {
      return;
    }
    const fileName = file.name.includes(".")
      ? file.name.slice(0, file.name.lastIndexOf("."))
      : file.name;
    const fileNameError = validateFileOrFolderName(fileName);
    if (fileNameError) {
      return;
    }
    validFiles.push(file);
  });

  if (validFiles.length > 0) {
    uploadQueueStore.enqueueFiles(validFiles, selectedUploadParentId.value);
  }

  target.value = "";
}

function onBindFileSelect(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file || !bindingTaskId.value) {
    target.value = "";
    return;
  }

  uploadQueueStore.bindTaskFile(bindingTaskId.value, file);
  bindingTaskId.value = "";
  target.value = "";
}
</script>

