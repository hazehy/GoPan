import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import {
  fileDeleteApi,
  fileDownloadApi,
  fileMoveApi,
  fileRenameApi,
  folderCreateApi,
  shareCreateApi,
} from '@/api/modules/disk';
import type { UserFile } from '@/types/api';
import { formatFileSize } from '@/utils/file';
import { toErrorMessage } from '@/utils/error';
import {
  formatUpdatedAt,
  getFileIcon,
  isFolder,
  parseUpdatedAtToMillis,
} from '@/utils/diskView';
import { useAuthStore } from '@/stores/auth';
import { useUploadQueueStore } from '@/stores/uploadQueue';
import {
  fetchAllFilesByParentId,
  fetchFolderListByParentId,
} from '@/composables/useUserFileList';
import {
  validateFileOrFolderName,
  validateShareExpiresDaysText,
} from '@/utils/validators';
import {
  alertDialog,
  confirmDialog,
  promptDialog,
} from '@/composables/useDialog';

export function useDiskBrowser() {
  const router = useRouter();
  const authStore = useAuthStore();
  const uploadQueueStore = useUploadQueueStore();

  const allFiles = ref<UserFile[]>([]);
  const page = ref(1);
  const size = ref(10);
  const currentParentId = ref(0);
  const breadcrumbs = ref<Array<{ id: number; name: string }>>([]);
  const errorMessage = ref('');
  const keyword = ref('');
  const typeFilter = ref<'all' | 'folder' | 'file'>('all');
  const sortBy = ref<'updated' | 'name' | 'size'>('updated');
  const sortOrder = ref<'asc' | 'desc'>('desc');
  const activeTab = ref<'files' | 'upload'>('files');
  const moveDialogVisible = ref(false);
  const movingItem = ref<UserFile | null>(null);
  const moveParentId = ref(0);
  const moveBreadcrumbs = ref<
    Array<{ id: number; identity: string; name: string }>
  >([]);
  const moveFolderOptions = ref<UserFile[]>([]);
  const moveLoading = ref(false);
  const moveErrorMessage = ref('');

  const currentTitle = computed(() =>
    activeTab.value === 'files' ? '我的文件' : '文件上传',
  );

  const displayFiles = computed(() => {
    const filtered = allFiles.value.filter((item) => {
      const isFolderItem = isFolder(item);
      if (typeFilter.value === 'folder' && !isFolderItem) {
        return false;
      }
      if (typeFilter.value === 'file' && isFolderItem) {
        return false;
      }

      if (!keyword.value) {
        return true;
      }

      const fullName = isFolderItem ? item.name : `${item.name}${item.ext ?? ''}`;
      return fullName.toLowerCase().includes(keyword.value.toLowerCase());
    });

    const direction = sortOrder.value === 'asc' ? 1 : -1;

    filtered.sort((left, right) => {
      const leftFolder = isFolder(left);
      const rightFolder = isFolder(right);
      if (leftFolder !== rightFolder) {
        return leftFolder ? -1 : 1;
      }

      if (sortBy.value === 'name') {
        const leftName = leftFolder ? left.name : `${left.name}${left.ext ?? ''}`;
        const rightName = rightFolder
          ? right.name
          : `${right.name}${right.ext ?? ''}`;
        return leftName.localeCompare(rightName, 'zh-Hans-CN') * direction;
      }

      if (sortBy.value === 'size') {
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
      return '根目录';
    }
    return `根目录 / ${moveBreadcrumbs.value.map((item) => item.name).join(' / ')}`;
  });

  const selectedMoveTargetIdentity = computed(() => {
    if (!moveBreadcrumbs.value.length) {
      return '';
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
    if (isFolder(movingItem.value) && targetIdentity === movingItem.value.identity) {
      return false;
    }
    return true;
  });

  async function fetchList() {
    allFiles.value = await fetchAllFilesByParentId(currentParentId.value);
  }

  async function refreshList() {
    try {
      errorMessage.value = '';
      await fetchList();
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
  }

  function changePage(nextPage: number) {
    page.value = nextPage;
  }

  function openFolder(item: UserFile) {
    if (!isFolder(item)) {
      return;
    }
    breadcrumbs.value.push({ id: item.id, name: item.name });
    currentParentId.value = item.id;
    page.value = 1;
    void refreshList();
  }

  async function downloadFile(item: UserFile) {
    if (isFolder(item) || !item.path) {
      return;
    }

    try {
      const filename = `${item.name}${item.ext ?? ''}`;
      const { url } = await fileDownloadApi({
        repository_identity: item.repository_identity,
        filename,
      });
      if (!url) {
        throw new Error('下载链接生成失败');
      }

      const link = document.createElement('a');
      link.href = url;
      link.target = '_self';
      link.rel = 'noopener';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
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
    void refreshList();
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
    void refreshList();
  }

  async function createFolder() {
    const name = await promptDialog('请输入文件夹名称', {
      title: '新建文件夹',
      placeholder: '请输入文件夹名称',
      promptValidator: (value: string) => validateFileOrFolderName(value.trim()),
    });
    if (!name) {
      return;
    }

    try {
      await folderCreateApi({ name, parent_id: currentParentId.value });
      await refreshList();
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
  }

  async function renameItem(item: UserFile) {
    const name = await promptDialog('请输入新名称', {
      title: '重命名',
      defaultValue: item.name,
      placeholder: '请输入新名称',
      promptValidator: (value: string) => validateFileOrFolderName(value.trim()),
    });
    if (!name || name === item.name) {
      return;
    }

    try {
      await fileRenameApi({ identity: item.identity, name });
      await refreshList();
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
  }

  async function deleteItem(item: UserFile) {
    const ok = await confirmDialog(
      `确认删除 ${item.name}${item.ext ?? ''} ?`,
      '删除确认',
    );
    if (!ok) {
      return;
    }
    try {
      await fileDeleteApi(item.identity);
      await refreshList();
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
  }

  async function moveItem(item: UserFile) {
    moveDialogVisible.value = true;
    movingItem.value = item;
    moveParentId.value = 0;
    moveBreadcrumbs.value = [];
    moveErrorMessage.value = '';
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
      moveErrorMessage.value = '';
      moveFolderOptions.value = await fetchFolderListByParentId(moveParentId.value);
    } catch (error) {
      moveErrorMessage.value = toErrorMessage(error);
    } finally {
      moveLoading.value = false;
    }
  }

  function resetMoveDialogState() {
    movingItem.value = null;
    moveDialogVisible.value = false;
    moveParentId.value = 0;
    moveBreadcrumbs.value = [];
    moveFolderOptions.value = [];
    moveErrorMessage.value = '';
  }

  function closeMoveDialog() {
    resetMoveDialogState();
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
      moveErrorMessage.value = '不能移动到自身目录';
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
      moveErrorMessage.value = '请选择目标文件夹';
      return;
    }

    try {
      moveLoading.value = true;
      moveErrorMessage.value = '';
      await fileMoveApi({
        identity: movingItem.value.identity,
        parent_identity: targetIdentity,
      });
      resetMoveDialogState();
      await alertDialog('移动成功', '提示');
      await refreshList();
    } catch (error) {
      moveErrorMessage.value = toErrorMessage(error);
    } finally {
      moveLoading.value = false;
    }
  }

  async function createShare(item: UserFile) {
    if (!item.repository_identity) {
      errorMessage.value = '该资源不支持分享';
      return;
    }

    const expiresText = await promptDialog('分享有效期（天），例如 7', {
      title: '创建分享',
      defaultValue: '7',
      placeholder: '请输入有效期（天）',
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
      const shareIdentity = (response.identity ?? '').trim();
      if (!shareIdentity) {
        throw new Error('创建分享失败：未返回分享标识');
      }
      const link = `${window.location.origin}/share/${shareIdentity}`;
      await navigator.clipboard.writeText(link);
      await alertDialog(`分享链接已复制：${link}`, '创建成功');
    } catch (error) {
      errorMessage.value = toErrorMessage(error);
    }
  }

  function logout() {
    uploadQueueStore.haltAllByLogout();
    authStore.clearAuth();
    router.replace('/login');
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

  return {
    allFiles,
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
  };
}
