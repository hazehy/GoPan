import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import {
  adminFileDeleteApi,
  adminFileListApi,
  adminLogListApi,
  adminOverviewApi,
  adminUserListApi,
  adminUserStatusUpdateApi,
} from '@/api/modules/admin';
import { useAuthStore } from '@/stores/auth';
import { useUploadQueueStore } from '@/stores/uploadQueue';
import type { AdminFileItem, AdminLogItem, AdminUserItem } from '@/types/api';
import { formatFileSize } from '@/utils/file';
import { toErrorMessage } from '@/utils/error';
import {
  buildLogQueryParams,
  formatActionLabel,
  formatActorDisplay,
  formatDateTime,
  formatLogDetail,
  formatText,
  getCurrentLogTitle,
  type LogSubPage,
} from '@/utils/adminLog';
import { resetPageAndLoad } from '@/composables/usePaging';
import { confirmDialog } from '@/composables/useDialog';

export function useAdminDashboard() {
  const router = useRouter();
  const authStore = useAuthStore();
  const uploadQueueStore = useUploadQueueStore();

  const activeMenu = ref<'overview' | 'users' | 'files' | 'logs'>('overview');

  const overview = ref({
    total_users: 0,
    active_users: 0,
    disabled_users: 0,
    total_files: 0,
    total_folders: 0,
    total_file_size: 0,
    today_uploads: 0,
    today_registers: 0,
    ext_stats: [] as Array<{ ext: string; count: number }>,
  });
  const users = ref<AdminUserItem[]>([]);
  const files = ref<AdminFileItem[]>([]);
  const logs = ref<AdminLogItem[]>([]);

  const pageSize = 10;
  const userPage = ref(1);
  const userCount = ref(0);
  const userKeyword = ref('');
  const filePage = ref(1);
  const fileCount = ref(0);
  const fileKeyword = ref('');
  const fileUserName = ref('');
  const userStatusLoading = ref(false);
  const fileDeleteLoading = ref(false);
  const logPage = ref(1);
  const logCount = ref(0);
  const logKeyword = ref('');
  const logActorName = ref('');
  const logAction = ref('');
  const logFileExt = ref('');
  const logSharerName = ref('');
  const logSaverName = ref('');
  const logDay = ref('');
  const logSubPage = ref<LogSubPage>('login');
  const errorMessage = ref('');

  const logActionOptions = [
    'USER_REGISTER',
    'USER_LOGIN',
    'FILE_UPLOAD',
    'SHARE_CREATE',
    'FILE_SAVE_REPOSITORY',
    'SHARE_RESOURCE_SAVE',
    'USER_STATUS_UPDATE',
  ];

  const logSubPageOptions = [
    { value: 'login', label: '用户登录日志' },
    { value: 'upload', label: '文件上传日志' },
    { value: 'share-create', label: '分享链接生成日志' },
    { value: 'share-save', label: '分享链接文件保存日志' },
    { value: 'others', label: '其他日志' },
  ] as const;

  const stats = computed(() => [
    { label: '总用户数', value: overview.value.total_users },
    { label: '活跃用户', value: overview.value.active_users },
    { label: '禁用用户', value: overview.value.disabled_users },
    { label: '总文件数', value: overview.value.total_files },
    { label: '总文件夹数', value: overview.value.total_folders },
    { label: '总文件容量', value: formatFileSize(overview.value.total_file_size) },
    { label: '今日上传', value: overview.value.today_uploads },
    { label: '今日注册', value: overview.value.today_registers },
  ]);

  const currentTitle = computed(() => {
    if (activeMenu.value === 'users') {
      return '用户管理';
    }
    if (activeMenu.value === 'files') {
      return '全文件管理';
    }
    if (activeMenu.value === 'logs') {
      return '操作日志';
    }
    return '数据总览';
  });

  const currentLogTitle = computed(() => getCurrentLogTitle(logSubPage.value));

  const logQueryParams = computed(() => {
    return buildLogQueryParams(logSubPage.value, {
      keyword: logKeyword.value,
      actorName: logActorName.value,
      action: logAction.value,
      fileExt: logFileExt.value,
      sharerName: logSharerName.value,
      saverName: logSaverName.value,
      day: logDay.value,
    });
  });

  const activeUserPercent = computed(() => {
    const total = overview.value.total_users || 0;
    if (total <= 0) {
      return 0;
    }
    return Math.round((overview.value.active_users / total) * 100);
  });

  const disabledUserPercent = computed(() => {
    const total = overview.value.total_users || 0;
    if (total <= 0) {
      return 0;
    }
    return Math.round((overview.value.disabled_users / total) * 100);
  });

  const extDistribution = computed(() => {
    const colors = ['#2563eb', '#3b82f6', '#60a5fa', '#93c5fd', '#1d4ed8', '#38bdf8', '#0ea5e9', '#0284c7'];
    const list = overview.value.ext_stats ?? [];
    const total = list.reduce((sum, item) => sum + item.count, 0);
    if (total <= 0) {
      return [];
    }
    return list.map((item, index) => ({
      ext: item.ext,
      count: item.count,
      percent: Math.round((item.count / total) * 100),
      color: colors[index % colors.length],
    }));
  });

  const extPieGradient = computed(() => {
    if (!extDistribution.value.length) {
      return 'conic-gradient(#e2e8f0 0deg 360deg)';
    }
    let start = 0;
    const segments: string[] = [];
    for (const item of extDistribution.value) {
      const deg = Math.round((item.percent / 100) * 360);
      const end = start + deg;
      segments.push(`${item.color} ${start}deg ${end}deg`);
      start = end;
    }
    if (start < 360) {
      segments.push(`${extDistribution.value[extDistribution.value.length - 1].color} ${start}deg 360deg`);
    }
    return `conic-gradient(${segments.join(', ')})`;
  });

  async function switchLogSubPage(next: LogSubPage) {
    if (logSubPage.value === next) {
      return;
    }
    logSubPage.value = next;
    await resetLogFilters();
  }

  async function resetLogFilters() {
    logKeyword.value = '';
    logActorName.value = '';
    logAction.value = '';
    logFileExt.value = '';
    logSharerName.value = '';
    logSaverName.value = '';
    logDay.value = '';
    await reloadLogs();
  }

  async function loadOverview() {
    overview.value = await adminOverviewApi();
  }

  async function loadUsers() {
    const res = await adminUserListApi({
      page: userPage.value,
      size: pageSize,
      keyword: userKeyword.value,
    });
    users.value = res.list;
    userCount.value = res.count;
  }

  async function loadFiles() {
    const res = await adminFileListApi({
      page: filePage.value,
      size: pageSize,
      keyword: fileKeyword.value,
      user_name: fileUserName.value,
    });
    files.value = res.list;
    fileCount.value = res.count;
  }

  async function loadLogs() {
    const params = logQueryParams.value;
    const res = await adminLogListApi({
      page: logPage.value,
      size: pageSize,
      keyword: params.keyword,
      action: params.action,
      actor_name: params.actor_name,
      file_ext: params.file_ext,
      sharer_name: params.sharer_name,
      saver_name: params.saver_name,
      day: params.day,
    });
    logs.value = res.list;
    logCount.value = res.count;
  }

  async function bootstrap() {
    try {
      errorMessage.value = '';
      await Promise.all([loadOverview(), loadUsers(), loadFiles(), loadLogs()]);
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
  }

  async function toggleUserStatus(identity: string, currentStatus: number) {
    if (userStatusLoading.value) {
      return;
    }
    userStatusLoading.value = true;
    try {
      await adminUserStatusUpdateApi({
        identity,
        status: currentStatus === 1 ? 2 : 1,
      });
      await Promise.all([loadUsers(), loadOverview()]);
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    } finally {
      userStatusLoading.value = false;
    }
  }

  async function toggleUserPermission(
    identity: string,
    field: 'upload_permission' | 'download_permission' | 'share_permission',
    currentValue: number,
  ) {
    if (userStatusLoading.value) {
      return;
    }
    userStatusLoading.value = true;
    try {
      await adminUserStatusUpdateApi({
        identity,
        [field]: currentValue === 1 ? 2 : 1,
      });
      await loadUsers();
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    } finally {
      userStatusLoading.value = false;
    }
  }

  async function deleteFile(identity: string) {
    if (fileDeleteLoading.value) {
      return;
    }
    const confirmed = await confirmDialog('确认删除该文件/文件夹及其子节点吗？', '删除确认');
    if (!confirmed) {
      return;
    }
    fileDeleteLoading.value = true;
    try {
      await adminFileDeleteApi(identity);
      await Promise.all([loadFiles(), loadOverview()]);
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    } finally {
      fileDeleteLoading.value = false;
    }
  }

  async function reloadUsers() {
    await resetPageAndLoad(userPage, loadUsers);
  }

  async function reloadFiles() {
    await resetPageAndLoad(filePage, loadFiles);
  }

  async function changeUserPage(nextPage: number) {
    userPage.value = nextPage;
    await loadUsers();
  }

  async function changeFilePage(nextPage: number) {
    filePage.value = nextPage;
    await loadFiles();
  }

  async function reloadLogs() {
    await resetPageAndLoad(logPage, loadLogs);
  }

  async function changeLogPage(nextPage: number) {
    logPage.value = nextPage;
    await loadLogs();
  }

  function logout() {
    uploadQueueStore.haltAllByLogout();
    authStore.clearAuth();
    router.replace('/login');
  }

  return {
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
    toggleUserStatus,
    toggleUserPermission,
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
  };
}
