<template>
  <div class="container share-container">
    <div class="card">
      <h2 class="auth-title">资源分享</h2>

      <div v-if="resource">
        <div class="share-detail-grid">
          <p><strong>文件名：</strong>{{ resource.name }}{{ resource.ext }}</p>
          <p><strong>大小：</strong>{{ formatFileSize(resource.size) }}</p>
          <p><strong>扩展名：</strong>{{ resource.ext || "-" }}</p>
          <p><strong>保存目录：</strong>{{ currentSavePath }}</p>
        </div>

        <div class="share-folder-picker" v-if="authStore.isLoggedIn">
          <div class="share-folder-head">
            <span class="muted">选择保存目录</span>
            <button
              class="btn btn-secondary"
              type="button"
              @click="goParentFolder"
              :disabled="folderBreadcrumbs.length === 0 || loadingFolders"
            >
              上一级
            </button>
          </div>

          <div class="path-bar" style="margin-bottom: 8px">
            <a href="javascript:void(0)" class="path-link" @click="goRootFolder"
              >根目录</a
            >
            <template
              v-for="(node, index) in folderBreadcrumbs"
              :key="`${node.id}-${index}`"
            >
              <span class="path-sep">/</span>
              <a
                href="javascript:void(0)"
                class="path-link"
                @click="goToFolder(index)"
                >{{ node.name }}</a
              >
            </template>
          </div>

          <div class="share-folder-list" v-if="folderOptions.length">
            <button
              v-for="folder in folderOptions"
              :key="folder.identity"
              class="share-folder-item"
              type="button"
              @click="enterFolder(folder)"
            >
              📁 {{ folder.name }}
            </button>
          </div>
          <p class="muted" v-else>当前目录没有子文件夹</p>
        </div>

        <div class="stack-row mt-12">
          <button
            class="btn btn-primary"
            @click="saveToDisk"
            :disabled="saving"
          >
            {{ saving ? "保存中..." : "保存到我的网盘" }}
          </button>
          <router-link class="btn btn-secondary" to="/disk"
            >返回网盘</router-link
          >
        </div>
      </div>

      <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { ResourceInfoResponse, UserFile } from "@/types/api";
import { resourceInfoApi, resourceSaveApi } from "@/api/modules/share";
import { fileListApi } from "@/api/modules/disk";
import { formatFileSize } from "@/utils/file";
import { useAuthStore } from "@/stores/auth";
import { alertDialog } from "@/composables/useDialog";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const resource = ref<ResourceInfoResponse | null>(null);
const saving = ref(false);
const loadingFolders = ref(false);
const errorMessage = ref("");
const saveParentId = ref(0);
const folderBreadcrumbs = ref<Array<{ id: number; name: string }>>([]);
const folderOptions = ref<UserFile[]>([]);

function getShareIdentity() {
  return String(route.params.identity ?? "").trim();
}

const currentSavePath = computed(() => {
  if (!folderBreadcrumbs.value.length) {
    return "根目录";
  }
  return `根目录 / ${folderBreadcrumbs.value.map((item) => item.name).join(" / ")}`;
});

function isFolder(item: UserFile) {
  return !item.repository_identity;
}

async function loadFolders() {
  if (!authStore.isLoggedIn) {
    folderOptions.value = [];
    return;
  }

  try {
    loadingFolders.value = true;
    const response = await fileListApi({
      id: saveParentId.value,
      page: 1,
      size: 200,
    });
    folderOptions.value = (response.list ?? []).filter((item) =>
      isFolder(item),
    );
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    loadingFolders.value = false;
  }
}

async function goRootFolder() {
  folderBreadcrumbs.value = [];
  saveParentId.value = 0;
  await loadFolders();
}

async function goParentFolder() {
  if (!folderBreadcrumbs.value.length) {
    return;
  }
  folderBreadcrumbs.value.pop();
  saveParentId.value = folderBreadcrumbs.value.length
    ? folderBreadcrumbs.value[folderBreadcrumbs.value.length - 1].id
    : 0;
  await loadFolders();
}

async function goToFolder(index: number) {
  folderBreadcrumbs.value = folderBreadcrumbs.value.slice(0, index + 1);
  saveParentId.value = folderBreadcrumbs.value[index].id;
  await loadFolders();
}

async function enterFolder(folder: UserFile) {
  folderBreadcrumbs.value.push({ id: folder.id, name: folder.name });
  saveParentId.value = folder.id;
  await loadFolders();
}

async function loadResource() {
  const shareIdentity = getShareIdentity();
  if (!shareIdentity) {
    errorMessage.value = "分享链接无效：缺少 identity";
    resource.value = null;
    return;
  }

  try {
    errorMessage.value = "";
    resource.value = await resourceInfoApi(shareIdentity);
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  }
}

async function saveToDisk() {
  if (!resource.value) {
    return;
  }

  if (!authStore.isLoggedIn) {
    await router.push({ path: "/login", query: { redirect: route.fullPath } });
    return;
  }

  try {
    saving.value = true;
    const shareIdentity = getShareIdentity();
    await resourceSaveApi({
      repository_identity: resource.value.repository_identity,
      parent_id: saveParentId.value,
      share_identity: shareIdentity,
    });
    await alertDialog(`已保存到 ${currentSavePath.value}`, "保存成功");
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  await loadResource();
  await loadFolders();
});

watch(
  () => authStore.isLoggedIn,
  async (loggedIn) => {
    if (loggedIn) {
      await loadFolders();
    }
  },
);
</script>
