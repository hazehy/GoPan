<template>
  <div class="container admin-container">
    <div
      v-if="sidebarVisible"
      class="sidebar-drawer-backdrop"
      @click="closeSidebar"
    ></div>
    <div class="admin-layout">
      <aside class="card admin-sidebar sidebar-drawer" :class="{ 'sidebar-drawer-open': sidebarVisible }">
        <h3 class="admin-sidebar-title">管理目录</h3>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'overview' }"
          @click="selectMenu('overview')"
        >
          数据总览
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'users' }"
          @click="selectMenu('users')"
        >
          用户管理
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'files' }"
          @click="selectMenu('files')"
        >
          全文件管理
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'logs' }"
          @click="selectMenu('logs')"
        >
          操作日志
        </button>
      </aside>

      <section class="admin-main">
        <div class="page-header">
          <button class="btn btn-secondary sidebar-toggle-btn" type="button" @click="toggleSidebar" aria-label="打开左侧目录">
            <span class="menu-icon" aria-hidden="true">☰</span>
            <span>目录</span>
          </button>
                <div v-if="sidebarVisible" class="sidebar-drawer-backdrop" @click="closeSidebar"></div>
          <h2 class="page-title">{{ currentTitle }}</h2>
                  <aside class="card admin-sidebar sidebar-drawer" :class="{ 'sidebar-drawer-open': sidebarVisible }">
        </div>

        <AdminOverviewPanel
          v-if="activeMenu === 'overview'"
                      @click="selectMenu('overview')"
          :active-user-percent="activeUserPercent"
          :disabled-user-percent="disabledUserPercent"
          :ext-distribution="extDistribution"
          :ext-pie-gradient="extPieGradient"
        />

                      @click="selectMenu('users')"
          v-if="activeMenu === 'users'"
          :users="users"
          :page-size="pageSize"
          :user-page="userPage"
          :user-count="userCount"
          :user-keyword="userKeyword"
                      @click="selectMenu('files')"
          :reload-users="reloadUsers"
          :change-user-page="changeUserPage"
          :update-user-status="updateUserStatus"
          :update-user-permission="updateUserPermission"
          :format-text="formatText"
          :format-date-time="formatDateTime"
                      @click="selectMenu('logs')"
        />

        <AdminFilesPanel
          v-if="activeMenu === 'files'"
          :files="files"
          :page-size="pageSize"
          :file-page="filePage"
                      <button class="btn btn-secondary sidebar-toggle-btn" type="button" @click="toggleSidebar" aria-label="切换左侧目录">
                        <span class="menu-icon" aria-hidden="true">☰</span>
                        <span>目录</span>
                      </button>
          :file-count="fileCount"
          :file-keyword="fileKeyword"
          :file-user-name="fileUserName"
          :file-delete-loading="fileDeleteLoading"
          :reload-files="reloadFiles"
          :change-file-page="changeFilePage"
          :delete-file="deleteFile"
          :format-text="formatText"
          import { ref, onMounted } from 'vue';
          :format-file-size="formatFileSize"
          @update:fileKeyword="fileKeyword = $event"
          @update:fileUserName="fileUserName = $event"
        />


          const sidebarVisible = ref(false);
        <AdminLogsPanel
          v-if="activeMenu === 'logs'"
          :logs="logs"
          :page-size="pageSize"
          :log-page="logPage"
          :log-count="logCount"
          :current-log-title="currentLogTitle"
          :log-sub-page="logSubPage"
          :log-sub-page-options="logSubPageOptions"
          :log-action-options="logActionOptions"
          :log-keyword="logKeyword"
          :log-actor-name="logActorName"
          :log-action="logAction"
          :log-file-ext="logFileExt"
          :log-sharer-name="logSharerName"
          :log-saver-name="logSaverName"
          :log-day="logDay"
          :switch-log-sub-page="switchLogSubPage"
          :reset-log-filters="resetLogFilters"
          :reload-logs="reloadLogs"
          :change-log-page="changeLogPage"
          :format-date-time="formatDateTime"
          :format-actor-display="formatActorDisplay"
          :format-action-label="formatActionLabel"
          :format-log-detail="formatLogDetail"
          @update:logKeyword="logKeyword = $event"
          @update:logActorName="logActorName = $event"
          @update:logAction="logAction = $event"
          @update:logFileExt="logFileExt = $event"
          @update:logSharerName="logSharerName = $event"
          @update:logSaverName="logSaverName = $event"
          @update:logDay="logDay = $event"
        />

        <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { onMounted } from 'vue';
import AdminOverviewPanel from '@/components/admin/AdminOverviewPanel.vue';
import AdminUsersPanel from '@/components/admin/AdminUsersPanel.vue';
import AdminFilesPanel from '@/components/admin/AdminFilesPanel.vue';
import AdminLogsPanel from '@/components/admin/AdminLogsPanel.vue';
import { useAdminDashboard } from '@/composables/useAdminDashboard';

const sidebarVisible = ref(false);

const {
  activeMenu,
  users,
  files,
  logs,
  pageSize,
  userPage,
  userCount,
  userKeyword,
  filePage,
  fileCount,
  fileKeyword,
  fileUserName,
  userStatusLoading,
  fileDeleteLoading,
  logPage,
  logCount,
  logKeyword,
  logActorName,
  logAction,
  logFileExt,
  logSharerName,
  logSaverName,
  logDay,
  logSubPage,
  errorMessage,
  logActionOptions,
  logSubPageOptions,
  stats,
  currentTitle,
  currentLogTitle,
  activeUserPercent,
  disabledUserPercent,
  extDistribution,
  extPieGradient,
  switchLogSubPage,
  resetLogFilters,
  updateUserStatus,
  updateUserPermission,
  deleteFile,
  reloadUsers,
  reloadFiles,
  changeUserPage,
  changeFilePage,
  reloadLogs,
  changeLogPage,
  bootstrap,
  logout,
  formatText,
  formatDateTime,
  formatFileSize,
  formatActorDisplay,
  formatActionLabel,
  formatLogDetail,
} = useAdminDashboard();

function openSidebar() {
  sidebarVisible.value = true;
}

function closeSidebar() {
  sidebarVisible.value = false;
}

function toggleSidebar() {
  sidebarVisible.value = !sidebarVisible.value;
}

function selectMenu(menu: typeof activeMenu.value) {
  activeMenu.value = menu;
  closeSidebar();
}

onMounted(() => {
  void bootstrap();
});
</script>
