<template>
  <div class="container">
    <div class="page-header">
      <h2 class="page-title">我的网盘</h2>
      <button class="btn btn-secondary" @click="logout">退出登录</button>
    </div>

    <ChunkUploader :parent-id="currentParentId" @success="refreshList" />

    <div class="card">
      <div class="toolbar">
        <button
          class="btn btn-secondary"
          @click="goParent"
          :disabled="breadcrumbs.length === 0"
        >
          返回上级
        </button>
        <button class="btn btn-primary" @click="createFolder">
          新建文件夹
        </button>
        <button class="btn btn-secondary" @click="refreshList">刷新</button>
      </div>

      <div class="path-bar">
        <span class="muted">当前位置：</span>
        <a
          href="javascript:void(0)"
          :class="[
            'path-link',
            { 'path-link-active': breadcrumbs.length === 0 },
          ]"
          @click="goToPath(-1)"
        >
          根目录
        </a>
        <template
          v-for="(node, index) in breadcrumbs"
          :key="`${node.id}-${index}`"
        >
          <span class="path-sep">/</span>
          <a
            href="javascript:void(0)"
            :class="[
              'path-link',
              { 'path-link-active': index === breadcrumbs.length - 1 },
            ]"
            @click="goToPath(index)"
          >
            {{ node.name }}
          </a>
        </template>
      </div>

      <div class="list-controls">
        <input
          class="input list-control-input"
          v-model.trim="keyword"
          type="text"
          placeholder="按名称筛选"
        />
        <select class="input list-control-select" v-model="typeFilter">
          <option value="all">全部</option>
          <option value="folder">仅文件夹</option>
          <option value="file">仅文件</option>
        </select>
        <select class="input list-control-select" v-model="sortBy">
          <option value="updated">按更新时间</option>
          <option value="name">按名称</option>
          <option value="size">按大小</option>
        </select>
        <select class="input list-control-select" v-model="sortOrder">
          <option value="desc">降序</option>
          <option value="asc">升序</option>
        </select>
      </div>

      <table class="table file-table">
        <colgroup>
          <col class="file-col-name" />
          <col class="file-col-size" />
          <col class="file-col-updated" />
          <col class="file-col-action" />
        </colgroup>
        <thead>
          <tr>
            <th>名称</th>
            <th>大小</th>
            <th>最后更新时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in pagedDisplayFiles" :key="item.identity">
            <td>
              <span class="file-name-cell">
                <span class="file-icon">{{ getFileIcon(item) }}</span>
                <a
                  v-if="isFolder(item)"
                  href="javascript:void(0)"
                  @click="openFolder(item)"
                  class="file-name-text"
                  :title="item.name"
                >
                  {{ item.name }}
                </a>
                <span
                  v-else
                  class="file-name-text"
                  :title="`${item.name}${item.ext ?? ''}`"
                >
                  {{ `${item.name}${item.ext ?? ""}` }}
                </span>
              </span>
            </td>
            <td>{{ formatFileSize(item.size) }}</td>
            <td>{{ formatUpdatedAt(item.updated_at) }}</td>
            <td>
              <div class="action-group">
                <button
                  class="btn btn-secondary"
                  @click="downloadFile(item)"
                  :disabled="isFolder(item) || !item.path"
                >
                  下载
                </button>
                <button class="btn btn-secondary" @click="renameItem(item)">
                  重命名
                </button>
                <button class="btn btn-secondary" @click="moveItem(item)">
                  移动
                </button>
                <button
                  class="btn btn-secondary"
                  @click="createShare(item)"
                  :disabled="!item.repository_identity"
                >
                  分享
                </button>
                <button class="btn btn-danger" @click="deleteItem(item)">
                  删除
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="!pagedDisplayFiles.length">
            <td colspan="4" class="muted">当前目录暂无文件</td>
          </tr>
        </tbody>
      </table>

      <div class="pagination">
        <button
          class="btn btn-secondary"
          :disabled="page <= 1"
          @click="changePage(page - 1)"
        >
          上一页
        </button>
        <button
          class="btn btn-secondary"
          :disabled="page * size >= total"
          @click="changePage(page + 1)"
        >
          下一页
        </button>
        <span class="muted">第 {{ page }} 页 / 共 {{ total }} 条</span>
      </div>

      <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
    </div>

    <div v-if="moveDialogVisible" class="dialog-mask" @click="closeMoveDialog">
      <div class="dialog-panel move-dialog-panel" @click.stop>
        <h3 class="dialog-title">移动到文件夹</h3>
        <p class="dialog-message">
          文件：{{
            movingItem ? `${movingItem.name}${movingItem.ext ?? ""}` : "-"
          }}
        </p>

        <div class="path-bar" style="margin-bottom: 8px">
          <a href="javascript:void(0)" class="path-link" @click="goMoveRoot"
            >根目录</a
          >
          <template
            v-for="(node, index) in moveBreadcrumbs"
            :key="node.identity"
          >
            <span class="path-sep">/</span>
            <a
              href="javascript:void(0)"
              class="path-link"
              @click="goMovePath(index)"
            >
              {{ node.name }}
            </a>
          </template>
        </div>

        <div class="move-folder-list" v-if="moveFolderOptions.length">
          <button
            v-for="folder in moveFolderOptions"
            :key="folder.identity"
            class="move-folder-item"
            type="button"
            :disabled="isMoveFolderSelf(folder)"
            @click="enterMoveFolder(folder)"
          >
            📁 {{ folder.name }}
          </button>
        </div>
        <p class="muted" v-else>当前目录没有子文件夹</p>

        <p class="muted" style="margin: 8px 0 0">
          目标目录：{{ currentMovePath }}
        </p>
        <p class="error" v-if="moveErrorMessage">{{ moveErrorMessage }}</p>

        <div class="dialog-actions" style="margin-top: 12px">
          <button
            class="btn btn-secondary"
            type="button"
            @click="closeMoveDialog"
          >
            取消
          </button>
          <button
            class="btn btn-primary"
            type="button"
            :disabled="moveLoading || !canConfirmMove"
            @click="confirmMove"
          >
            {{ moveLoading ? "移动中..." : "确认移动" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import ChunkUploader from "@/components/ChunkUploader.vue";
import {
  fileDeleteApi,
  fileListApi,
  fileMoveApi,
  fileRenameApi,
  folderCreateApi,
  shareCreateApi,
} from "@/api/modules/disk";
import type { UserFile } from "@/types/api";
import { formatFileSize } from "@/utils/file";
import { useAuthStore } from "@/stores/auth";
import {
  validateFileOrFolderName,
  validateShareExpiresDaysText,
} from "@/utils/validators";
import {
  alertDialog,
  confirmDialog,
  promptDialog,
} from "@/composables/useDialog";

const router = useRouter();
const authStore = useAuthStore();

const allFiles = ref<UserFile[]>([]);
const page = ref(1);
const size = ref(10);
const currentParentId = ref(0);
const breadcrumbs = ref<Array<{ id: number; name: string }>>([]);
const errorMessage = ref("");
const keyword = ref("");
const typeFilter = ref<"all" | "folder" | "file">("all");
const sortBy = ref<"updated" | "name" | "size">("updated");
const sortOrder = ref<"asc" | "desc">("desc");
const moveDialogVisible = ref(false);
const movingItem = ref<UserFile | null>(null);
const moveParentId = ref(0);
const moveBreadcrumbs = ref<
  Array<{ id: number; identity: string; name: string }>
>([]);
const moveFolderOptions = ref<UserFile[]>([]);
const moveLoading = ref(false);
const moveErrorMessage = ref("");

function isFolder(item: UserFile) {
  return !item.repository_identity;
}

const displayFiles = computed(() => {
  const filtered = allFiles.value.filter((item) => {
    const isFolderItem = isFolder(item);
    if (typeFilter.value === "folder" && !isFolderItem) {
      return false;
    }
    if (typeFilter.value === "file" && isFolderItem) {
      return false;
    }

    if (!keyword.value) {
      return true;
    }

    const fullName = isFolderItem ? item.name : `${item.name}${item.ext ?? ""}`;
    return fullName.toLowerCase().includes(keyword.value.toLowerCase());
  });

  const direction = sortOrder.value === "asc" ? 1 : -1;

  filtered.sort((left, right) => {
    const leftFolder = isFolder(left);
    const rightFolder = isFolder(right);
    if (leftFolder !== rightFolder) {
      return leftFolder ? -1 : 1;
    }

    if (sortBy.value === "name") {
      const leftName = leftFolder ? left.name : `${left.name}${left.ext ?? ""}`;
      const rightName = rightFolder
        ? right.name
        : `${right.name}${right.ext ?? ""}`;
      return leftName.localeCompare(rightName, "zh-Hans-CN") * direction;
    }

    if (sortBy.value === "size") {
      return (left.size - right.size) * direction;
    }

    const leftTime = parseUpdatedAtToMillis(left.updated_at);
    const rightTime = parseUpdatedAtToMillis(right.updated_at);
    return (leftTime - rightTime) * direction;
  });

  return filtered;
});

const total = computed(() => displayFiles.value.length);

const pagedDisplayFiles = computed(() => {
  const start = (page.value - 1) * size.value;
  const end = start + size.value;
  return displayFiles.value.slice(start, end);
});

const currentMovePath = computed(() => {
  if (!moveBreadcrumbs.value.length) {
    return "根目录";
  }
  return `根目录 / ${moveBreadcrumbs.value.map((item) => item.name).join(" / ")}`;
});

const selectedMoveTargetIdentity = computed(() => {
  if (!moveBreadcrumbs.value.length) {
    return "";
  }
  return moveBreadcrumbs.value[moveBreadcrumbs.value.length - 1].identity;
});

const canConfirmMove = computed(() => {
  if (!movingItem.value) {
    return false;
  }
  const targetIdentity = selectedMoveTargetIdentity.value;
  if (!targetIdentity) {
    return false;
  }
  if (
    isFolder(movingItem.value) &&
    targetIdentity === movingItem.value.identity
  ) {
    return false;
  }
  return true;
});

function formatUpdatedAt(value: string) {
  if (!value) {
    return "-";
  }

  const parsed = parseUpdatedAtParts(value);
  if (!parsed) {
    return value;
  }

  const pad = (num: number) => String(num).padStart(2, "0");
  return `${parsed.year}-${pad(parsed.month)}-${pad(parsed.day)} ${pad(parsed.hour)}:${pad(parsed.minute)}`;
}

function parseUpdatedAtToMillis(value: string) {
  const parsed = parseUpdatedAtParts(value);
  if (parsed) {
    return new Date(
      parsed.year,
      parsed.month - 1,
      parsed.day,
      parsed.hour,
      parsed.minute,
      parsed.second,
    ).getTime();
  }

  return Date.parse(value || "") || 0;
}

function parseUpdatedAtParts(value: string) {
  const normalized = value.replace("T", " ").replace("Z", "").split(".")[0];
  const matched = normalized.match(
    /^(\d{4})-(\d{2})-(\d{2})\s(\d{2}):(\d{2})(?::(\d{2}))?/,
  );
  if (!matched) {
    return null;
  }

  return {
    year: Number(matched[1]),
    month: Number(matched[2]),
    day: Number(matched[3]),
    hour: Number(matched[4]),
    minute: Number(matched[5]),
    second: Number(matched[6] ?? "0"),
  };
}

function getFileIcon(item: UserFile) {
  if (isFolder(item)) {
    return "📁";
  }

  const ext = (item.ext ?? "").toLowerCase();

  if (
    [".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"].includes(ext)
  ) {
    return "🖼️";
  }
  if ([".mp4", ".avi", ".mkv", ".mov", ".flv", ".wmv"].includes(ext)) {
    return "🎬";
  }
  if ([".mp3", ".wav", ".flac", ".aac", ".ogg"].includes(ext)) {
    return "🎵";
  }
  if ([".zip", ".rar", ".7z", ".tar", ".gz"].includes(ext)) {
    return "🗜️";
  }
  if ([".doc", ".docx", ".txt", ".md", ".rtf"].includes(ext)) {
    return "📄";
  }
  if ([".xls", ".xlsx", ".csv"].includes(ext)) {
    return "📊";
  }
  if ([".ppt", ".pptx"].includes(ext)) {
    return "📽️";
  }
  if (ext === ".pdf") {
    return "📕";
  }

  return "📦";
}

async function fetchList() {
  const chunkSize = 200;
  let currentPage = 1;
  let loaded: UserFile[] = [];
  let totalCount = 0;

  do {
    const response = await fileListApi({
      id: currentParentId.value,
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

  allFiles.value = loaded;
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

async function refreshList() {
  try {
    errorMessage.value = "";
    await fetchList();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  }
}

async function changePage(nextPage: number) {
  page.value = nextPage;
  await refreshList();
}

function openFolder(item: UserFile) {
  if (!isFolder(item)) {
    return;
  }
  breadcrumbs.value.push({ id: item.id, name: item.name });
  currentParentId.value = item.id;
  page.value = 1;
  refreshList();
}

function downloadFile(item: UserFile) {
  if (isFolder(item) || !item.path) {
    return;
  }

  const link = document.createElement("a");
  link.href = item.path;
  link.download = `${item.name}${item.ext ?? ""}`;
  link.target = "_blank";
  link.rel = "noopener";
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

function goParent() {
  if (breadcrumbs.value.length === 0) {
    return;
  }
  breadcrumbs.value.pop();
  currentParentId.value = breadcrumbs.value.length
    ? breadcrumbs.value[breadcrumbs.value.length - 1].id
    : 0;
  page.value = 1;
  refreshList();
}

function goToPath(index: number) {
  if (index < 0) {
    breadcrumbs.value = [];
    currentParentId.value = 0;
  } else {
    breadcrumbs.value = breadcrumbs.value.slice(0, index + 1);
    currentParentId.value = breadcrumbs.value[index].id;
  }
  page.value = 1;
  refreshList();
}

async function createFolder() {
  const name = await promptDialog("请输入文件夹名称", {
    title: "新建文件夹",
    placeholder: "请输入文件夹名称",
    promptValidator: (value: string) => validateFileOrFolderName(value.trim()),
  });
  if (!name) {
    return;
  }

  try {
    await folderCreateApi({ name, parent_id: currentParentId.value });
    await refreshList();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  }
}

async function renameItem(item: UserFile) {
  const name = await promptDialog("请输入新名称", {
    title: "重命名",
    defaultValue: item.name,
    placeholder: "请输入新名称",
    promptValidator: (value: string) => validateFileOrFolderName(value.trim()),
  });
  if (!name || name === item.name) {
    return;
  }

  try {
    await fileRenameApi({ identity: item.identity, name });
    await refreshList();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  }
}

async function deleteItem(item: UserFile) {
  const ok = await confirmDialog(
    `确认删除 ${item.name}${item.ext ?? ""} ?`,
    "删除确认",
  );
  if (!ok) {
    return;
  }
  try {
    await fileDeleteApi(item.identity);
    await refreshList();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  }
}

async function moveItem(item: UserFile) {
  moveDialogVisible.value = true;
  movingItem.value = item;
  moveParentId.value = 0;
  moveBreadcrumbs.value = [];
  moveErrorMessage.value = "";
  await loadMoveFolderOptions();
}

function isMoveFolderSelf(folder: UserFile) {
  if (!movingItem.value || !isFolder(movingItem.value)) {
    return false;
  }
  return folder.identity === movingItem.value.identity;
}

async function loadMoveFolderOptions() {
  try {
    moveLoading.value = true;
    moveErrorMessage.value = "";
    moveFolderOptions.value = await fetchFolderListByParentId(
      moveParentId.value,
    );
  } catch (error) {
    moveErrorMessage.value =
      error instanceof Error ? error.message : String(error);
  } finally {
    moveLoading.value = false;
  }
}

function closeMoveDialog() {
  moveDialogVisible.value = false;
  movingItem.value = null;
  moveParentId.value = 0;
  moveBreadcrumbs.value = [];
  moveFolderOptions.value = [];
  moveErrorMessage.value = "";
}

async function goMoveRoot() {
  moveParentId.value = 0;
  moveBreadcrumbs.value = [];
  await loadMoveFolderOptions();
}

async function goMovePath(index: number) {
  moveBreadcrumbs.value = moveBreadcrumbs.value.slice(0, index + 1);
  moveParentId.value = moveBreadcrumbs.value[index].id;
  await loadMoveFolderOptions();
}

async function enterMoveFolder(folder: UserFile) {
  if (isMoveFolderSelf(folder)) {
    moveErrorMessage.value = "不能移动到自身目录";
    return;
  }
  moveBreadcrumbs.value.push({
    id: folder.id,
    identity: folder.identity,
    name: folder.name,
  });
  moveParentId.value = folder.id;
  await loadMoveFolderOptions();
}

async function confirmMove() {
  if (!movingItem.value) {
    return;
  }
  const targetIdentity = selectedMoveTargetIdentity.value;
  if (!targetIdentity) {
    moveErrorMessage.value = "请选择目标文件夹";
    return;
  }

  try {
    moveLoading.value = true;
    moveErrorMessage.value = "";
    await fileMoveApi({
      identity: movingItem.value.identity,
      parent_identity: targetIdentity,
    });
    closeMoveDialog();
    await alertDialog("移动成功", "提示");
    await refreshList();
  } catch (error) {
    moveErrorMessage.value =
      error instanceof Error ? error.message : String(error);
  } finally {
    moveLoading.value = false;
  }
}

async function createShare(item: UserFile) {
  if (!item.repository_identity) {
    errorMessage.value = "该资源不支持分享";
    return;
  }

  const expiresText = await promptDialog("分享有效期（天），例如 7", {
    title: "创建分享",
    defaultValue: "7",
    placeholder: "请输入有效期（天）",
    promptValidator: (value: string) =>
      validateShareExpiresDaysText(value.trim()),
  });
  if (!expiresText) {
    return;
  }

  const expires = Number(expiresText);

  try {
    const response = await shareCreateApi({
      repository_identity: item.repository_identity,
      expires,
    });
    const shareIdentity = (response.identity ?? "").trim();
    if (!shareIdentity) {
      throw new Error("创建分享失败：未返回分享标识");
    }
    const link = `${window.location.origin}/share/${shareIdentity}`;
    await navigator.clipboard.writeText(link);
    await alertDialog(`分享链接已复制：${link}`, "创建成功");
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  }
}

function logout() {
  authStore.clearAuth();
  router.replace("/login");
}

onMounted(refreshList);

watch([keyword, typeFilter, sortBy, sortOrder], () => {
  page.value = 1;
});

watch(total, (nextTotal) => {
  const maxPage = Math.max(1, Math.ceil(nextTotal / size.value));
  if (page.value > maxPage) {
    page.value = maxPage;
  }
});
</script>
