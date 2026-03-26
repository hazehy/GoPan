<template>
  <div class="card uploader-card">
    <input
      ref="fileInputRef"
      class="uploader-native-input"
      type="file"
      @change="onFileSelect"
    />

    <div class="uploader-top">
      <div class="uploader-file-meta">
        <p class="uploader-file-name">
          {{ selectedFile ? selectedFile.name : "请选择要上传的文件" }}
        </p>
        <p class="muted" v-if="selectedFile || totalBytes > 0">
          总大小：{{ formatSize(totalBytes || selectedFile?.size || 0) }}
        </p>
      </div>

      <div class="uploader-actions">
        <button
          class="btn btn-secondary"
          :disabled="uploading"
          @click="chooseFile"
        >
          选择文件
        </button>
        <button
          class="btn btn-primary"
          :disabled="!selectedFile || uploading"
          @click="startUpload"
        >
          {{ uploading ? "上传中..." : "开始上传" }}
        </button>
      </div>
    </div>

    <div class="uploader-progress" v-if="showUploadPanel">
      <div
        class="uploader-progress-bar"
        :style="{ width: `${progress}%` }"
      ></div>
    </div>

    <div class="uploader-row" v-if="showUploadPanel">
      <span class="uploader-state">{{
        progress === 100 ? "上传完成" : "上传处理中"
      }}</span>
      <span class="muted">{{ progress }}%</span>
      <button
        v-if="canCloseUploadPanel"
        class="uploader-close-btn"
        type="button"
        aria-label="关闭上传状态"
        @click="closeUploadPanel"
      >
        ×
      </button>
    </div>

    <div class="uploader-stats" v-if="showUploadPanel">
      <span class="muted">已上传：{{ formatSize(uploadedBytes) }}</span>
      <span class="muted">总大小：{{ formatSize(totalBytes) }}</span>
      <span class="muted">上传速度：{{ formatSpeed(speedBytesPerSec) }}</span>
    </div>

    <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import {
  fileChunkUploadApi,
  fileChunkUploadCompleteApi,
  filePreUploadApi,
  fileUploadApi,
  userRepositoryApi,
} from "@/api/modules/disk";
import { calcFileMd5 } from "@/utils/hash";
import { splitNameAndExt } from "@/utils/file";
import { validateFileOrFolderName } from "@/utils/validators";

const props = defineProps<{
  parentId: number;
}>();

const emit = defineEmits<{
  success: [];
}>();

const selectedFile = ref<File | null>(null);
const uploading = ref(false);
const errorMessage = ref("");
const fileInputRef = ref<HTMLInputElement | null>(null);
const uploadedBytes = ref(0);
const totalBytes = ref(0);
const uploadStartAt = ref(0);
const speedBytesPerSec = ref(0);
const showUploadPanel = ref(false);

const progress = computed(() => {
  if (totalBytes.value <= 0) {
    return 0;
  }
  return Math.min(
    100,
    Math.round((uploadedBytes.value / totalBytes.value) * 100),
  );
});

const canCloseUploadPanel = computed(
  () => !uploading.value && progress.value === 100,
);

const CHUNK_SIZE = 2 * 1024 * 1024;

function chooseFile() {
  fileInputRef.value?.click();
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

function closeUploadPanel() {
  if (!canCloseUploadPanel.value) {
    return;
  }
  showUploadPanel.value = false;
}

function onFileSelect(event: Event) {
  const target = event.target as HTMLInputElement;
  selectedFile.value = target.files?.[0] ?? null;
  uploadedBytes.value = 0;
  totalBytes.value = 0;
  speedBytesPerSec.value = 0;
  errorMessage.value = "";
  showUploadPanel.value = false;
}

function buildPathByKey(key: string) {
  const host = import.meta.env.VITE_COS_BUCKET_URL?.trim();
  if (!host) {
    return key;
  }
  return `${host.replace(/\/$/, "")}/${key}`;
}

async function startUpload() {
  if (!selectedFile.value || uploading.value) {
    return;
  }

  const file = selectedFile.value;
  if (file.size <= 0) {
    errorMessage.value = "文件大小不能为 0";
    return;
  }

  const { name, ext } = splitNameAndExt(file.name);
  const fileNameError = validateFileOrFolderName(name);
  if (fileNameError) {
    errorMessage.value = fileNameError;
    return;
  }

  if ((ext ?? "").length > 20) {
    errorMessage.value = "文件扩展名过长";
    return;
  }

  try {
    uploading.value = true;
    errorMessage.value = "";
    uploadedBytes.value = 0;
    totalBytes.value = file.size;
    speedBytesPerSec.value = 0;
    uploadStartAt.value = Date.now();
    showUploadPanel.value = true;

    const md5 = await calcFileMd5(file);

    const preUploadResult = await filePreUploadApi({ md5, name, ext });

    if (preUploadResult.identity) {
      await userRepositoryApi({
        parent_id: props.parentId,
        repository_identity: preUploadResult.identity,
        ext,
        name,
      });
      uploadedBytes.value = totalBytes.value;
      emit("success");
      return;
    }

    const uploadId = preUploadResult.upload_id;
    const key = preUploadResult.key;

    if (!uploadId || !key) {
      throw new Error("预上传返回数据不完整");
    }

    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    const parts: Array<{ part_number: number; etag: string }> = [];

    for (let index = 0; index < totalChunks; index += 1) {
      const start = index * CHUNK_SIZE;
      const end = Math.min(file.size, start + CHUNK_SIZE);
      const chunk = file.slice(start, end);
      const partNumber = index + 1;

      const formData = new FormData();
      formData.append("file", chunk, file.name);
      formData.append("key", key);
      formData.append("upload_id", uploadId);
      formData.append("part_number", String(partNumber));

      const chunkResp = await fileChunkUploadApi(formData);
      parts.push({
        part_number: partNumber,
        etag: chunkResp.etag,
      });

      uploadedBytes.value = end;
      const elapsedSeconds = (Date.now() - uploadStartAt.value) / 1000;
      speedBytesPerSec.value =
        elapsedSeconds > 0 ? uploadedBytes.value / elapsedSeconds : 0;
    }

    await fileChunkUploadCompleteApi({
      key,
      upload_id: uploadId,
      cos_objects: parts,
    });

    const fileUploadResult = await fileUploadApi({
      hash: md5,
      name,
      ext,
      size: file.size,
      path: buildPathByKey(key),
    });

    if (!fileUploadResult?.identity) {
      throw new Error("文件入库失败：未获取到 repository identity");
    }

    await userRepositoryApi({
      parent_id: props.parentId,
      repository_identity: fileUploadResult.identity,
      ext: fileUploadResult.ext,
      name,
    });

    uploadedBytes.value = totalBytes.value;
    emit("success");
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    uploading.value = false;
    if (fileInputRef.value) {
      fileInputRef.value.value = "";
    }
    selectedFile.value = null;
  }
}
</script>
