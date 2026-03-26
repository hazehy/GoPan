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

        <div v-if="activeMenu === 'overview'" class="admin-view-stack">
          <div class="admin-stats-grid">
            <div class="card admin-stat-card" v-for="item in stats" :key="item.label">
              <p class="muted">{{ item.label }}</p>
              <p class="admin-stat-value">{{ item.value }}</p>
            </div>
          </div>

          <div class="admin-chart-grid">
            <div class="card admin-chart-card">
              <h3 class="admin-subtitle">用户状态分布</h3>
              <div class="admin-bar-row">
                <span class="admin-bar-label">活跃用户</span>
                <div class="admin-bar-track">
                  <div class="admin-bar-fill admin-bar-fill-primary" :style="{ width: `${activeUserPercent}%` }"></div>
                </div>
                <span class="admin-bar-value">{{ activeUserPercent }}%</span>
              </div>
              <div class="admin-bar-row">
                <span class="admin-bar-label">禁用用户</span>
                <div class="admin-bar-track">
                  <div class="admin-bar-fill admin-bar-fill-secondary" :style="{ width: `${disabledUserPercent}%` }"></div>
                </div>
                <span class="admin-bar-value">{{ disabledUserPercent }}%</span>
              </div>
            </div>

            <div class="card admin-chart-card">
              <h3 class="admin-subtitle">文件类型分布（按扩展名）</h3>
              <div class="admin-pie-layout" v-if="extDistribution.length">
                <div class="admin-pie-chart" :style="{ background: extPieGradient }"></div>
                <div class="admin-pie-legend">
                  <div class="admin-pie-legend-item" v-for="item in extDistribution" :key="item.ext">
                    <span class="admin-pie-dot" :style="{ background: item.color }"></span>
                    <span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="`${item.ext} ${item.percent}%`" :title="`${item.ext} ${item.percent}%`">
                      {{ item.ext }} {{ item.percent }}%
                    </span>
                  </div>
                </div>
              </div>
              <p class="muted" v-if="!extDistribution.length">暂无文件类型数据</p>
            </div>
          </div>

          <div class="card admin-section">
            <h3 class="admin-subtitle">数据概览表</h3>
            <table class="table admin-overview-table">
              <thead>
                <tr>
                  <th>指标</th>
                  <th>数值</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in stats" :key="`table-${item.label}`">
                  <td>{{ item.label }}</td>
                  <td>{{ item.value }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="activeMenu === 'users'" class="card admin-section">
          <div class="admin-section-head">
            <h3>用户管理</h3>
            <div class="admin-actions-inline">
              <input class="input admin-filter-input" v-model.trim="userKeyword" placeholder="按用户名/邮箱/标识搜索" />
              <button class="btn btn-secondary" @click="reloadUsers">查询</button>
            </div>
          </div>
          <table class="table admin-table">
            <colgroup>
              <col class="admin-col-user-name" />
              <col class="admin-col-user-email" />
              <col class="admin-col-user-role" />
              <col class="admin-col-user-status" />
              <col class="admin-col-user-login" />
              <col class="admin-col-user-action" />
            </colgroup>
            <thead>
              <tr>
                <th>用户名</th>
                <th>邮箱</th>
                <th>角色</th>
                <th>状态</th>
                <th>最近登录</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.identity">
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="user.name" :title="user.name">{{ user.name }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="user.email" :title="user.email">{{ user.email }}</span></td>
                <td>{{ user.role === 2 ? '管理员' : '普通用户' }}</td>
                <td>{{ user.status === 1 ? '正常' : '禁用' }}</td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="formatText(user.last_login_at)" :title="formatText(user.last_login_at)">{{ formatDateTime(user.last_login_at) }}</span></td>
                <td>
                  <button
                    class="btn btn-secondary"
                    :disabled="user.role === 2 || userStatusLoading"
                    @click="toggleUserStatus(user.identity, user.status)"
                  >
                    {{ user.status === 1 ? '禁用' : '启用' }}
                  </button>
                </td>
              </tr>
              <tr v-if="!users.length">
                <td colspan="6" class="muted">暂无用户数据</td>
              </tr>
            </tbody>
          </table>
          <div class="pagination">
            <button class="btn btn-secondary" :disabled="userPage <= 1" @click="changeUserPage(userPage - 1)">上一页</button>
            <button class="btn btn-secondary" :disabled="userPage * pageSize >= userCount" @click="changeUserPage(userPage + 1)">下一页</button>
            <span class="muted">第 {{ userPage }} 页 / 共 {{ userCount }} 条</span>
          </div>
        </div>

        <div v-if="activeMenu === 'files'" class="card admin-section">
          <div class="admin-section-head">
            <h3>全文件管理</h3>
            <div class="admin-actions-inline">
              <input class="input admin-filter-input" v-model.trim="fileKeyword" placeholder="按名称/标识搜索" />
              <input class="input admin-filter-input" v-model.trim="fileUserName" placeholder="按用户名过滤" />
              <button class="btn btn-secondary" @click="reloadFiles">查询</button>
            </div>
          </div>
          <table class="table admin-table">
            <colgroup>
              <col class="admin-col-file-name" />
              <col class="admin-col-file-user" />
              <col class="admin-col-file-path" />
              <col class="admin-col-file-type" />
              <col class="admin-col-file-size" />
              <col class="admin-col-file-updated" />
              <col class="admin-col-file-action" />
            </colgroup>
            <thead>
              <tr>
                <th>名称</th>
                <th>用户名</th>
                <th>路径</th>
                <th>类型</th>
                <th>大小</th>
                <th>更新时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in files" :key="item.identity">
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="`${item.name}${item.ext || ''}`" :title="`${item.name}${item.ext || ''}`">{{ item.name }}{{ item.ext || '' }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.user_name || '-'" :title="item.user_name || '-'">{{ item.user_name || '-' }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="formatText(item.path)" :title="formatText(item.path)">{{ formatText(item.path) }}</span></td>
                <td>{{ item.repository_identity ? '文件' : '文件夹' }}</td>
                <td>{{ formatFileSize(item.size) }}</td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="formatText(item.updated_at)" :title="formatText(item.updated_at)">{{ formatDateTime(item.updated_at) }}</span></td>
                <td>
                  <button class="btn btn-danger" :disabled="fileDeleteLoading" @click="deleteFile(item.identity)">删除</button>
                </td>
              </tr>
              <tr v-if="!files.length">
                <td colspan="7" class="muted">暂无文件数据</td>
              </tr>
            </tbody>
          </table>
          <div class="pagination">
            <button class="btn btn-secondary" :disabled="filePage <= 1" @click="changeFilePage(filePage - 1)">上一页</button>
            <button class="btn btn-secondary" :disabled="filePage * pageSize >= fileCount" @click="changeFilePage(filePage + 1)">下一页</button>
            <span class="muted">第 {{ filePage }} 页 / 共 {{ fileCount }} 条</span>
          </div>
        </div>

        <div v-if="activeMenu === 'logs'" class="card admin-section">
          <div class="admin-section-head">
            <h3>操作日志</h3>
            <span class="muted">{{ currentLogTitle }}</span>
          </div>

          <div class="admin-log-tabs">
            <button
              v-for="item in logSubPageOptions"
              :key="item.value"
              class="admin-log-tab"
              :class="{ 'admin-log-tab-active': logSubPage === item.value }"
              @click="switchLogSubPage(item.value)"
            >
              {{ item.label }}
            </button>
          </div>

          <div class="admin-log-filters">
            <input
              v-if="logSubPage === 'login' || logSubPage === 'others'"
              class="input"
              v-model.trim="logActorName"
              placeholder="按用户名筛选"
            />

            <input
              v-if="logSubPage === 'upload' || logSubPage === 'share-create' || logSubPage === 'share-save'"
              class="input"
              v-model.trim="logFileExt"
              placeholder="按文件类型筛选，如 pdf"
            />

            <input
              v-if="logSubPage === 'share-create' || logSubPage === 'share-save'"
              class="input"
              v-model.trim="logSharerName"
              placeholder="按分享者用户名筛选"
            />

            <input
              v-if="logSubPage === 'share-save'"
              class="input"
              v-model.trim="logSaverName"
              placeholder="按保存者用户名筛选"
            />

            <input
              v-if="logSubPage === 'others'"
              class="input"
              v-model.trim="logKeyword"
              placeholder="按描述/目标标识搜索"
            />

            <select v-if="logSubPage === 'others'" class="input" v-model="logAction">
              <option value="">全部操作类型</option>
              <option v-for="item in logActionOptions" :key="item" :value="item">{{ item }}</option>
            </select>

            <input class="input" type="date" v-model="logDay" />
            <button class="btn btn-primary admin-log-filter-btn" @click="reloadLogs">查询</button>
            <button class="btn btn-secondary admin-log-filter-btn" @click="resetLogFilters">重置</button>
          </div>

          <table class="table admin-table">
            <colgroup>
              <col class="admin-col-log-time" />
              <col class="admin-col-log-user" />
              <col class="admin-col-log-action" />
              <col class="admin-col-log-target" />
              <col class="admin-col-log-detail" />
            </colgroup>
            <thead>
              <tr>
                <th>时间</th>
                <th>操作人</th>
                <th>操作类型</th>
                <th>目标标识</th>
                <th>描述</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in logs" :key="item.identity">
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.created_at" :title="item.created_at">{{ formatDateTime(item.created_at) }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.actor_name || item.actor_identity || '-'" :title="item.actor_name || item.actor_identity || '-'">{{ formatActorDisplay(item.actor_name, item.actor_identity) }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.action" :title="item.action">{{ formatActionLabel(item.action) }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.target_identity || '-'" :title="item.target_identity || '-'">{{ item.target_identity || '-' }}</span></td>
                <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.detail || '-'" :title="item.detail || '-'">{{ formatLogDetail(item.detail) }}</span></td>
              </tr>
              <tr v-if="!logs.length">
                <td colspan="5" class="muted">暂无日志数据</td>
              </tr>
            </tbody>
          </table>
          <div class="pagination">
            <button class="btn btn-secondary" :disabled="logPage <= 1" @click="changeLogPage(logPage - 1)">上一页</button>
            <button class="btn btn-secondary" :disabled="logPage * pageSize >= logCount" @click="changeLogPage(logPage + 1)">下一页</button>
            <span class="muted">第 {{ logPage }} 页 / 共 {{ logCount }} 条</span>
          </div>
        </div>

        <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
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
import type { AdminFileItem, AdminLogItem, AdminUserItem } from '@/types/api';
import { formatFileSize } from '@/utils/file';
import { confirmDialog } from '@/composables/useDialog';

const router = useRouter();
const authStore = useAuthStore();
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
const logSubPage = ref<'login' | 'upload' | 'share-create' | 'share-save' | 'others'>('login');
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

const currentLogTitle = computed(() => {
  if (logSubPage.value === 'login') {
    return '筛选项：用户名 + 日期';
  }
  if (logSubPage.value === 'upload') {
    return '筛选项：文件类型 + 日期';
  }
  if (logSubPage.value === 'share-create') {
    return '筛选项：分享者用户名 + 文件类型 + 日期';
  }
  if (logSubPage.value === 'share-save') {
    return '筛选项：分享者用户名 + 保存者用户名 + 文件类型 + 日期';
  }
  return '筛选项：操作人 + 类型 + 关键字 + 日期';
});

const logQueryParams = computed(() => {
  const base = {
    keyword: '',
    action: '',
    actor_name: '',
    file_ext: '',
    sharer_name: '',
    saver_name: '',
    day: logDay.value,
  };

  if (logSubPage.value === 'login') {
    return {
      ...base,
      action: 'USER_LOGIN',
      actor_name: logActorName.value,
    };
  }
  if (logSubPage.value === 'upload') {
    return {
      ...base,
      action: 'FILE_UPLOAD',
      file_ext: logFileExt.value,
    };
  }
  if (logSubPage.value === 'share-create') {
    return {
      ...base,
      action: 'SHARE_CREATE',
      actor_name: logSharerName.value,
      sharer_name: logSharerName.value,
      file_ext: logFileExt.value,
    };
  }
  if (logSubPage.value === 'share-save') {
    return {
      ...base,
      action: 'SHARE_RESOURCE_SAVE',
      sharer_name: logSharerName.value,
      saver_name: logSaverName.value,
      file_ext: logFileExt.value,
    };
  }
  return {
    ...base,
    keyword: logKeyword.value,
    action: logAction.value,
    actor_name: logActorName.value,
  };
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

function formatText(value: string) {
  return value || '-';
}

function formatDateTime(value: string) {
  const raw = formatText(value);
  if (raw === '-') {
    return raw;
  }

  const normalized = raw.replace('T', ' ').replace('Z', '').split('.')[0];
  const matched = normalized.match(/^(\d{4}-\d{2}-\d{2})\s(\d{2}:\d{2}:\d{2})/);
  if (matched) {
    return `${matched[1]} ${matched[2]}`;
  }

  return normalized.replace(/\s[+-]\d{4}\s[A-Za-z]+$/, '');
}

function formatActorDisplay(actorName: string, actorIdentity: string) {
  const name = (actorName || '').trim();
  if (name) {
    return name;
  }
  const identity = (actorIdentity || '').trim();
  if (!identity) {
    return '-';
  }
  if (identity.length <= 12) {
    return identity;
  }
  return `${identity.slice(0, 8)}...`;
}

function formatActionLabel(action: string) {
  const map: Record<string, string> = {
    USER_REGISTER: '用户注册',
    USER_LOGIN: '用户登录',
    FILE_UPLOAD: '文件上传',
    SHARE_CREATE: '创建分享链接',
    FILE_SAVE_REPOSITORY: '保存文件到网盘',
    SHARE_RESOURCE_SAVE: '保存分享文件',
    USER_STATUS_UPDATE: '用户状态更新',
  };
  return map[action] || action;
}

function formatLogDetail(detail: string) {
  const text = (detail || '').trim();
  if (!text) {
    return '-';
  }
  if (!text.includes('=')) {
    return text;
  }

  const map: Record<string, string> = {
    file_ext: '文件类型',
    upload_mode: '上传模式',
    file_name: '文件名',
    sharer_name: '分享者',
    saver_name: '保存者',
    expires_days: '有效期(天)',
    share_identity: '分享标识',
    saved_name: '保存名',
  };

  return text
    .split(';')
    .map((item) => item.trim())
    .filter(Boolean)
    .map((item) => {
      const index = item.indexOf('=');
      if (index <= 0) {
        return item;
      }
      const key = item.slice(0, index);
      const value = item.slice(index + 1);
      return `${map[key] || key}: ${value || '-'}`;
    })
    .join(' | ');
}

async function switchLogSubPage(next: 'login' | 'upload' | 'share-create' | 'share-save' | 'others') {
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
    errorMessage.value = error instanceof Error ? error.message : String(error);
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
    errorMessage.value = error instanceof Error ? error.message : String(error);
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
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    fileDeleteLoading.value = false;
  }
}

async function reloadUsers() {
  userPage.value = 1;
  await loadUsers();
}

async function reloadFiles() {
  filePage.value = 1;
  await loadFiles();
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
  logPage.value = 1;
  await loadLogs();
}

async function changeLogPage(nextPage: number) {
  logPage.value = nextPage;
  await loadLogs();
}

function logout() {
  authStore.clearAuth();
  router.replace('/login');
}

onMounted(() => {
  void bootstrap();
});
</script>
