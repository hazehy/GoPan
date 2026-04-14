import { defineStore } from "pinia";
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

const STORAGE_KEY = "gopan_upload_tasks_v1";
const CHUNK_SIZE = 2 * 1024 * 1024;

export type UploadTaskStatus =
  | "waiting"
  | "uploading"
  | "paused"
  | "stopped"
  | "success"
  | "failed";

export interface UploadTask {
  id: string;
  name: string;
  ext: string;
  size: number;
  parentId: number;
  status: UploadTaskStatus;
  progress: number;
  uploadedBytes: number;
  speedBytesPerSec: number;
  errorMessage: string;
  createdAt: number;
  updatedAt: number;
  hasLocalFile: boolean;
  md5: string;
  uploadId: string;
  key: string;
  totalChunks: number;
  nextPartNumber: number;
  uploadParts: Array<{ part_number: number; etag: string }>;
}

interface RuntimeTaskContext {
  file: File;
  abortController: AbortController | null;
  speedStartAt: number;
  speedBaseUploadedBytes: number;
}

interface UploadQueueState {
  tasks: UploadTask[];
  initialized: boolean;
  successVersion: number;
}

const runtimeTaskMap = new Map<string, RuntimeTaskContext>();
let authClearListenerBound = false;

function persistTasks(tasks: UploadTask[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(tasks));
}

function loadTasks(): UploadTask[] {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    return [];
  }

  try {
    const parsed = JSON.parse(raw) as UploadTask[];
    if (!Array.isArray(parsed)) {
      return [];
    }

    return parsed
      .filter((item) => item && typeof item.id === "string")
      .map((item) => {
        const status = item.status === "uploading" ? "stopped" : item.status;
        const totalChunks =
          Number(item.totalChunks) > 0
            ? Number(item.totalChunks)
            : Math.max(1, Math.ceil((Number(item.size) || 0) / CHUNK_SIZE));
        const uploadParts = Array.isArray(item.uploadParts) ? item.uploadParts : [];
        const nextPartNumber =
          Number(item.nextPartNumber) > 0
            ? Number(item.nextPartNumber)
            : Math.min(totalChunks + 1, uploadParts.length + 1);
        return {
          ...item,
          status,
          hasLocalFile: false,
          speedBytesPerSec: 0,
          md5: typeof item.md5 === "string" ? item.md5 : "",
          uploadId: typeof item.uploadId === "string" ? item.uploadId : "",
          key: typeof item.key === "string" ? item.key : "",
          totalChunks,
          nextPartNumber,
          uploadParts,
          errorMessage:
            status === "stopped" && item.status === "uploading"
              ? "会话已结束，任务已停止"
              : item.errorMessage,
        };
      });
  } catch {
    return [];
  }
}

function createTaskId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

function isCanceledError(error: unknown) {
  if (!error || typeof error !== "object") {
    return false;
  }
  const maybe = error as { code?: string; name?: string; message?: string; __CANCEL__?: boolean };
  const text = `${maybe.message ?? ""}`.toLowerCase();
  return (
    maybe.code === "ERR_CANCELED" ||
    maybe.name === "CanceledError" ||
    maybe.name === "AbortError" ||
    maybe.__CANCEL__ === true ||
    text.includes("canceled") ||
    text.includes("aborted")
  );
}

function markTaskFailed(task: UploadTask, message: string) {
  task.status = "failed";
  task.errorMessage = message;
  task.speedBytesPerSec = 0;
  task.updatedAt = Date.now();
}

export const useUploadQueueStore = defineStore("uploadQueue", {
  state: (): UploadQueueState => ({
    tasks: [],
    initialized: false,
    successVersion: 0,
  }),
  getters: {
    taskCount: (state) => state.tasks.length,
    waitingCount: (state) =>
      state.tasks.filter((item) => item.status === "waiting").length,
    uploadingTask: (state) =>
      state.tasks.find((item) => item.status === "uploading") ?? null,
    sortedTasks: (state) =>
      [...state.tasks].sort((left, right) => left.createdAt - right.createdAt),
  },
  actions: {
    ensureInitialized() {
      if (!this.initialized) {
        this.tasks = loadTasks();
        this.initialized = true;
      }

      if (!authClearListenerBound && typeof window !== "undefined") {
        window.addEventListener("gopan-auth-cleared", this.haltAllByLogout);
        authClearListenerBound = true;
      }
    },
    enqueueFiles(files: File[], parentId: number) {
      this.ensureInitialized();
      const now = Date.now();

      files.forEach((file) => {
        const { name, ext } = splitNameAndExt(file.name);
        const taskId = createTaskId();
        this.tasks.push({
          id: taskId,
          name,
          ext,
          size: file.size,
          parentId,
          status: "waiting",
          progress: 0,
          uploadedBytes: 0,
          speedBytesPerSec: 0,
          errorMessage: "",
          createdAt: now,
          updatedAt: now,
          hasLocalFile: true,
          md5: "",
          uploadId: "",
          key: "",
          totalChunks: Math.ceil(file.size / CHUNK_SIZE),
          nextPartNumber: 1,
          uploadParts: [],
        });
        runtimeTaskMap.set(taskId, {
          file,
          abortController: null,
          speedStartAt: 0,
          speedBaseUploadedBytes: 0,
        });
      });

      persistTasks(this.tasks);
      void this.processNext();
    },
    bindTaskFile(taskId: string, file: File) {
      this.ensureInitialized();
      const task = this.tasks.find((item) => item.id === taskId);
      if (!task) {
        return;
      }

      const { name, ext } = splitNameAndExt(file.name);
      if (name !== task.name || ext !== task.ext || file.size !== task.size) {
        task.status = "failed";
        task.errorMessage = "请选择同名且大小一致的原文件以继续断点续传";
        task.updatedAt = Date.now();
        persistTasks(this.tasks);
        return;
      }

      runtimeTaskMap.set(taskId, {
        file,
        abortController: null,
        speedStartAt: 0,
        speedBaseUploadedBytes: 0,
      });
      task.hasLocalFile = true;
      task.errorMessage = "";
      if (task.status === "failed") {
        task.status = "stopped";
      }
      task.updatedAt = Date.now();
      persistTasks(this.tasks);
    },
    pauseTask(taskId: string) {
      this.ensureInitialized();
      const task = this.tasks.find((item) => item.id === taskId);
      if (!task) {
        return;
      }

      if (task.status === "uploading") {
        task.status = "paused";
        task.speedBytesPerSec = 0;
        task.updatedAt = Date.now();
        const runtime = runtimeTaskMap.get(taskId);
        runtime?.abortController?.abort();
      } else if (task.status === "waiting") {
        task.status = "paused";
        task.updatedAt = Date.now();
      }

      persistTasks(this.tasks);
      void this.processNext();
    },
    resumeTask(taskId: string) {
      this.ensureInitialized();
      const task = this.tasks.find((item) => item.id === taskId);
      if (!task) {
        return;
      }

      if (!task.hasLocalFile || !runtimeTaskMap.has(taskId)) {
        markTaskFailed(task, "任务文件句柄已失效，请重新选择文件上传");
        persistTasks(this.tasks);
        return;
      }

      if (["paused", "stopped", "failed"].includes(task.status)) {
        task.status = "waiting";
        task.errorMessage = "";
        task.speedBytesPerSec = 0;
        task.updatedAt = Date.now();
      }

      persistTasks(this.tasks);
      void this.processNext();
    },
    removeTask(taskId: string) {
      this.ensureInitialized();
      const task = this.tasks.find((item) => item.id === taskId);
      if (!task) {
        return;
      }

      if (task.status === "uploading") {
        const runtime = runtimeTaskMap.get(taskId);
        runtime?.abortController?.abort();
      }

      this.tasks = this.tasks.filter((item) => item.id !== taskId);
      runtimeTaskMap.delete(taskId);
      persistTasks(this.tasks);
      void this.processNext();
    },
    clearFinishedTasks() {
      this.ensureInitialized();
      this.tasks
        .filter((item) => item.status === "success")
        .forEach((item) => runtimeTaskMap.delete(item.id));
      this.tasks = this.tasks.filter((item) => item.status !== "success");
      persistTasks(this.tasks);
    },
    startQueue() {
      this.ensureInitialized();
      const now = Date.now();
      this.tasks.forEach((item) => {
        if (["paused", "stopped", "failed"].includes(item.status) && item.hasLocalFile) {
          item.status = "waiting";
          item.errorMessage = "";
          item.speedBytesPerSec = 0;
          item.updatedAt = now;
        }
      });
      persistTasks(this.tasks);
      void this.processNext();
    },
    stopQueue() {
      this.ensureInitialized();
      const now = Date.now();
      this.tasks.forEach((item) => {
        if (item.status === "waiting") {
          item.status = "stopped";
          item.updatedAt = now;
        }
        if (item.status === "uploading") {
          item.status = "stopped";
          item.speedBytesPerSec = 0;
          item.updatedAt = now;
          const runtime = runtimeTaskMap.get(item.id);
          runtime?.abortController?.abort();
        }
      });
      persistTasks(this.tasks);
    },
    haltAllByLogout() {
      this.ensureInitialized();
      const now = Date.now();
      this.tasks.forEach((item) => {
        if (item.status === "success") {
          return;
        }
        if (item.status === "uploading") {
          const runtime = runtimeTaskMap.get(item.id);
          runtime?.abortController?.abort();
        }
        item.status = "stopped";
        item.speedBytesPerSec = 0;
        item.errorMessage = "已退出登录，任务已停止";
        item.updatedAt = now;
      });
      persistTasks(this.tasks);
    },
    async processNext() {
      this.ensureInitialized();
      if (this.tasks.some((item) => item.status === "uploading")) {
        return;
      }
      const next = this.tasks.find((item) => item.status === "waiting");
      if (!next) {
        return;
      }
      await this.uploadTask(next.id);
    },
    async uploadTask(taskId: string) {
      const task = this.tasks.find((item) => item.id === taskId);
      if (!task) {
        return;
      }

      const runtime = runtimeTaskMap.get(taskId);
      if (!runtime) {
        markTaskFailed(task, "任务文件句柄已失效，请重新选择文件上传");
        persistTasks(this.tasks);
        return;
      }

      const file = runtime.file;
      if (file.size <= 0) {
        markTaskFailed(task, "文件大小不能为 0");
        persistTasks(this.tasks);
        await this.processNext();
        return;
      }

      const { name, ext } = splitNameAndExt(file.name);
      const fileNameError = validateFileOrFolderName(name);
      if (fileNameError) {
        markTaskFailed(task, fileNameError);
        persistTasks(this.tasks);
        await this.processNext();
        return;
      }

      if ((ext ?? "").length > 20) {
        markTaskFailed(task, "文件扩展名过长");
        persistTasks(this.tasks);
        await this.processNext();
        return;
      }

      const controller = new AbortController();
      runtime.abortController = controller;
      runtime.speedStartAt = Date.now();
      runtime.speedBaseUploadedBytes = task.uploadedBytes;

      task.name = name;
      task.ext = ext;
      task.size = file.size;
      task.status = "uploading";
      task.errorMessage = "";
      task.speedBytesPerSec = 0;
      task.totalChunks = Math.max(1, Math.ceil(file.size / CHUNK_SIZE));
      if (!task.nextPartNumber || task.nextPartNumber < 1) {
        task.nextPartNumber = 1;
      }
      task.updatedAt = Date.now();
      persistTasks(this.tasks);

      try {
        // Step 1: hash + pre-upload negotiation (fast pass for deduplicated file).
        let md5 = task.md5;
        if (!md5) {
          md5 = await calcFileMd5(file);
          task.md5 = md5;
          task.updatedAt = Date.now();
          persistTasks(this.tasks);
        }

        if (task.status !== "uploading") {
          return;
        }

        if (!task.uploadId || !task.key) {
          const preUploadResult = await filePreUploadApi({
            md5,
            name,
            ext,
          });

          if (task.status !== "uploading") {
            return;
          }

          if (preUploadResult.identity) {
            await userRepositoryApi({
              parent_id: task.parentId,
              repository_identity: preUploadResult.identity,
              ext,
              name,
            });
            task.uploadedBytes = task.size;
            task.progress = 100;
            task.status = "success";
            task.speedBytesPerSec = 0;
            task.nextPartNumber = task.totalChunks + 1;
            task.updatedAt = Date.now();
            this.successVersion += 1;
            persistTasks(this.tasks);
            return;
          }

          if (!preUploadResult.upload_id || !preUploadResult.key) {
            throw new Error("预上传返回数据不完整");
          }

          task.uploadId = preUploadResult.upload_id;
          task.key = preUploadResult.key;
          task.uploadParts = [];
          task.nextPartNumber = 1;
          task.uploadedBytes = 0;
          task.progress = 0;
          task.updatedAt = Date.now();
          persistTasks(this.tasks);
        }

        const uploadId = task.uploadId;
        const key = task.key;
        const totalChunks = task.totalChunks;
        const startPart = Math.max(1, task.nextPartNumber);

        // Step 2: resumable chunk upload loop; keep nextPartNumber/uploadParts durable in localStorage.
        for (let partNumber = startPart; partNumber <= totalChunks; partNumber += 1) {
          if (task.status !== "uploading") {
            return;
          }

          const start = (partNumber - 1) * CHUNK_SIZE;
          const end = Math.min(file.size, start + CHUNK_SIZE);
          const chunk = file.slice(start, end);
          const chunkSize = end - start;

          const formData = new FormData();
          formData.append("file", chunk, file.name);
          formData.append("key", key);
          formData.append("upload_id", uploadId);
          formData.append("part_number", String(partNumber));

          const chunkResp = await fileChunkUploadApi(formData, {
            signal: controller.signal,
            onUploadProgress: (loaded, total) => {
              if (task.status !== "uploading") {
                return;
              }

              const totalBytes = total && total > 0 ? total : chunkSize;
              const loadedBytes = Math.max(0, Math.min(totalBytes, loaded));
              task.uploadedBytes = Math.min(task.size, start + loadedBytes);
              task.progress = Math.min(
                100,
                Math.round((task.uploadedBytes / task.size) * 100),
              );

              const elapsedSeconds = (Date.now() - runtime.speedStartAt) / 1000;
              const sessionUploadedBytes =
                task.uploadedBytes - runtime.speedBaseUploadedBytes;
              task.speedBytesPerSec =
                elapsedSeconds > 0
                  ? Math.max(0, sessionUploadedBytes) / elapsedSeconds
                  : 0;
              task.updatedAt = Date.now();
            },
          });

          const existingPartIndex = task.uploadParts.findIndex(
            (item) => item.part_number === partNumber,
          );
          const uploadedPart = {
            part_number: partNumber,
            etag: chunkResp.etag,
          };
          if (existingPartIndex >= 0) {
            task.uploadParts[existingPartIndex] = uploadedPart;
          } else {
            task.uploadParts.push(uploadedPart);
          }

          task.uploadedBytes = end;
          task.progress = Math.min(
            100,
            Math.round((task.uploadedBytes / task.size) * 100),
          );
          task.nextPartNumber = partNumber + 1;
          const elapsedSeconds = (Date.now() - runtime.speedStartAt) / 1000;
          const sessionUploadedBytes =
            task.uploadedBytes - runtime.speedBaseUploadedBytes;
          task.speedBytesPerSec =
            elapsedSeconds > 0
              ? Math.max(0, sessionUploadedBytes) / elapsedSeconds
              : 0;
          task.updatedAt = Date.now();
          persistTasks(this.tasks);
        }

        await fileChunkUploadCompleteApi(
          {
            key,
            upload_id: uploadId,
            cos_objects: [...task.uploadParts].sort(
              (left, right) => left.part_number - right.part_number,
            ),
          },
          controller.signal,
        );

        const fileUploadResult = await fileUploadApi(
          {
            hash: md5,
            name,
            ext,
            size: file.size,
            key,
          },
          controller.signal,
        );

        if (!fileUploadResult?.identity) {
          throw new Error("文件入库失败：未获取到 repository identity");
        }

        // Step 3: link uploaded object into user's repository tree.
        await userRepositoryApi(
          {
            parent_id: task.parentId,
            repository_identity: fileUploadResult.identity,
            ext: fileUploadResult.ext,
            name,
          },
          controller.signal,
        );

        task.uploadedBytes = task.size;
        task.progress = 100;
        task.status = "success";
        task.speedBytesPerSec = 0;
        task.nextPartNumber = task.totalChunks + 1;
        task.updatedAt = Date.now();
        this.successVersion += 1;
        persistTasks(this.tasks);
      } catch (error) {
        const currentStatus = String(task.status);
        if (
          isCanceledError(error) ||
          currentStatus === "paused" ||
          currentStatus === "stopped"
        ) {
          if (task.status === "uploading") {
            task.status = "paused";
          }
          task.errorMessage = "";
          task.speedBytesPerSec = 0;
          task.updatedAt = Date.now();
          persistTasks(this.tasks);
          return;
        }

        markTaskFailed(task, error instanceof Error ? error.message : String(error));
        persistTasks(this.tasks);
      } finally {
        runtime.abortController = null;
        runtime.speedStartAt = 0;
        runtime.speedBaseUploadedBytes = task.uploadedBytes;
        await this.processNext();
      }
    },
  },
});
