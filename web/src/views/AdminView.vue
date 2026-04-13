<template>
  <div class="container admin-container">
    <div class="admin-layout">
      <aside class="card admin-sidebar">
        <h3 class="admin-sidebar-title">管理目录</h3>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'overview' }"
          @click="activeMenu = 'overview'"
        >
          数据总览
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'users' }"
          @click="activeMenu = 'users'"
        >
          用户管理
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'files' }"
          @click="activeMenu = 'files'"
        >
          全文件管理
        </button>
        <button
          class="admin-menu-item"
          :class="{ 'admin-menu-item-active': activeMenu === 'logs' }"
          @click="activeMenu = 'logs'"
        >
          操作日志
        </button>
      </aside>

      <section class="admin-main">
        <div class="page-header">
          <h2 class="page-title">{{ currentTitle }}</h2>
          <button class="btn btn-secondary" @click="logout">退出登录</button>
        </div>

        <AdminOverviewPanel
          v-if="activeMenu === 'overview'"
          :stats="stats"
          :active-user-percent="activeUserPercent"
          :disabled-user-percent="disabledUserPercent"
          :ext-distribution="extDistribution"
          :ext-pie-gradient="extPieGradient"
        />

        <AdminUsersPanel
          v-if="activeMenu === 'users'"
          :users="users"
          :page-size="pageSize"
          :user-page="userPage"
          :user-count="userCount"
          :user-keyword="userKeyword"
          :user-status-loading="userStatusLoading"
          :reload-users="reloadUsers"
          :change-user-page="changeUserPage"
          :update-user-status="updateUserStatus"
          :update-user-permission="updateUserPermission"
          :format-text="formatText"
          :format-date-time="formatDateTime"
          @update:userKeyword="userKeyword = $event"
        />

        <AdminFilesPanel
          v-if="activeMenu === 'files'"
          :files="files"
          :page-size="pageSize"
          :file-page="filePage"
          :file-count="fileCount"
          :file-keyword="fileKeyword"
          :file-user-name="fileUserName"
          :file-delete-loading="fileDeleteLoading"
          :reload-files="reloadFiles"
          :change-file-page="changeFilePage"
          :delete-file="deleteFile"
          :format-text="formatText"
          :format-date-time="formatDateTime"
          :format-file-size="formatFileSize"
          @update:fileKeyword="fileKeyword = $event"
          @update:fileUserName="fileUserName = $event"
        />

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
import { onMounted } from 'vue';
import AdminOverviewPanel from '@/components/admin/AdminOverviewPanel.vue';
import AdminUsersPanel from '@/components/admin/AdminUsersPanel.vue';
import AdminFilesPanel from '@/components/admin/AdminFilesPanel.vue';
import AdminLogsPanel from '@/components/admin/AdminLogsPanel.vue';
import { useAdminDashboard } from '@/composables/useAdminDashboard';

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

onMounted(() => {
  void bootstrap();
});
</script>
