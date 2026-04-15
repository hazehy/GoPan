<template>
  <div class="container admin-container">
    <div class="admin-layout disk-layout">
      <div v-if="sidebarVisible" class="sidebar-drawer-backdrop" @click="closeSidebar"></div>
      <aside class="card admin-sidebar disk-sidebar sidebar-drawer" :class="{ 'sidebar-drawer-open': sidebarVisible }">
        <h3 class="admin-sidebar-title">我的网盘</h3>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeTab === 'files' }"
          type="button"
          @click="selectTab('files')"
        >
          我的文件
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeTab === 'upload' }"
          type="button"
          @click="selectTab('upload')"
        >
          文件上传
        </button>
      </aside>

      <section class="admin-main disk-main">
        <div class="page-header">
          <button class="btn btn-secondary sidebar-toggle-btn" type="button" @click="toggleSidebar" aria-label="切换左侧目录">
            <span class="menu-icon" aria-hidden="true">☰</span>
            <span>目录</span>
          </button>
          <h2 class="page-title">{{ currentTitle }}</h2>
          <button class="btn btn-secondary" @click="logout">退出登录</button>
        </div>

        <ChunkUploader
          v-if="activeTab === 'upload'"
          :parent-id="currentParentId"
          :browse-path="breadcrumbs"
          @success="refreshList"
        />

        <DiskFilePanel
          v-if="activeTab === 'files'"
          :breadcrumbs="breadcrumbs"
          :keyword="keyword"
          :type-filter="typeFilter"
          :sort-by="sortBy"
          :sort-order="sortOrder"
          :page="page"
          :size="size"
          :total="total"
          :error-message="errorMessage"
          :paged-display-files="pagedDisplayFiles"
          :go-parent="goParent"
          :create-folder="createFolder"
          :refresh-list="refreshList"
          :go-to-path="goToPath"
          :change-page="changePage"
          :open-folder="openFolder"
          :download-file="downloadFile"
          :rename-item="renameItem"
          :move-item="moveItem"
          :create-share="createShare"
          :delete-item="deleteItem"
          :is-folder="isFolder"
          :get-file-icon="getFileIcon"
          :format-updated-at="formatUpdatedAt"
          :format-file-size="formatFileSize"
          @update:keyword="keyword = $event"
          @update:typeFilter="typeFilter = $event"
          @update:sortBy="sortBy = $event"
          @update:sortOrder="sortOrder = $event"
        />
      </section>
    </div>

    <DiskMoveDialog
      :visible="moveDialogVisible"
      :moving-item="movingItem"
      :move-breadcrumbs="moveBreadcrumbs"
      :move-folder-options="moveFolderOptions"
      :current-move-path="currentMovePath"
      :move-error-message="moveErrorMessage"
      :move-loading="moveLoading"
      :can-confirm-move="canConfirmMove"
      :close-move-dialog="closeMoveDialog"
      :go-move-root="goMoveRoot"
      :go-move-path="goMovePath"
      :enter-move-folder="enterMoveFolder"
      :confirm-move="confirmMove"
      :is-move-folder-self="isMoveFolderSelf"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ChunkUploader from "@/components/ChunkUploader.vue";
import DiskFilePanel from "@/components/disk/DiskFilePanel.vue";
import DiskMoveDialog from "@/components/disk/DiskMoveDialog.vue";
import { useDiskBrowser } from "@/composables/useDiskBrowser";

const sidebarVisible = ref(false);

const {
  page,
  size,
  currentParentId,
  breadcrumbs,
  errorMessage,
  keyword,
  typeFilter,
  sortBy,
  sortOrder,
  activeTab,
  moveDialogVisible,
  movingItem,
  moveBreadcrumbs,
  moveFolderOptions,
  moveLoading,
  moveErrorMessage,
  currentTitle,
  total,
  pagedDisplayFiles,
  currentMovePath,
  canConfirmMove,
  refreshList,
  changePage,
  openFolder,
  downloadFile,
  goParent,
  goToPath,
  createFolder,
  renameItem,
  deleteItem,
  moveItem,
  isMoveFolderSelf,
  closeMoveDialog,
  goMoveRoot,
  goMovePath,
  enterMoveFolder,
  confirmMove,
  createShare,
  logout,
  isFolder,
  getFileIcon,
  formatUpdatedAt,
  formatFileSize,
} = useDiskBrowser();

function closeSidebar() {
  sidebarVisible.value = false;
}

function toggleSidebar() {
  sidebarVisible.value = !sidebarVisible.value;
}

function selectTab(tab: typeof activeTab.value) {
  activeTab.value = tab;
  closeSidebar();
}
</script>
