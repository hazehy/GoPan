<template>
  <div class="card admin-section">
    <div class="admin-section-head">
      <h3>操作日志</h3>
      <span class="muted">{{ props.currentLogTitle }}</span>
    </div>

    <div class="admin-log-tabs">
      <button
        v-for="item in props.logSubPageOptions"
        :key="item.value"
        class="admin-log-tab"
        :class="{ 'admin-log-tab-active': props.logSubPage === item.value }"
        @click="props.switchLogSubPage(item.value)"
      >
        {{ item.label }}
      </button>
    </div>

    <div class="admin-log-filters">
      <input
        v-if="props.logSubPage === 'login' || props.logSubPage === 'others'"
        class="input"
        :value="props.logActorName"
        @input="onLogActorNameInput"
        placeholder="按用户名筛选"
      />

      <input
        v-if="props.logSubPage === 'upload' || props.logSubPage === 'share-create' || props.logSubPage === 'share-save'"
        class="input"
        :value="props.logFileExt"
        @input="onLogFileExtInput"
        placeholder="按文件类型筛选，如 pdf"
      />

      <input
        v-if="props.logSubPage === 'share-create' || props.logSubPage === 'share-save'"
        class="input"
        :value="props.logSharerName"
        @input="onLogSharerNameInput"
        placeholder="按分享者用户名筛选"
      />

      <input
        v-if="props.logSubPage === 'share-save'"
        class="input"
        :value="props.logSaverName"
        @input="onLogSaverNameInput"
        placeholder="按保存者用户名筛选"
      />

      <input
        v-if="props.logSubPage === 'others'"
        class="input"
        :value="props.logKeyword"
        @input="onLogKeywordInput"
        placeholder="按描述/目标标识搜索"
      />

      <select
        v-if="props.logSubPage === 'others'"
        class="input"
        :value="props.logAction"
        @change="onLogActionChange"
      >
        <option value="">全部操作类型</option>
        <option v-for="item in props.logActionOptions" :key="item" :value="item">{{ item }}</option>
      </select>

      <input class="input" type="date" :value="props.logDay" @input="onLogDayInput" />
      <button class="btn btn-primary admin-log-filter-btn" @click="props.reloadLogs">查询</button>
      <button class="btn btn-secondary admin-log-filter-btn" @click="props.resetLogFilters">重置</button>
    </div>

    <div class="x-scroll-panel">
      <div class="table-scroll-content">
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
            <tr v-for="item in props.logs" :key="item.identity">
              <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.created_at" :title="item.created_at">{{ props.formatDateTime(item.created_at) }}</span></td>
              <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.actor_name || item.actor_identity || '-'" :title="item.actor_name || item.actor_identity || '-'">{{ props.formatActorDisplay(item.actor_name, item.actor_identity) }}</span></td>
              <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.action" :title="item.action">{{ props.formatActionLabel(item.action) }}</span></td>
              <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.target_identity || '-'" :title="item.target_identity || '-'">{{ item.target_identity || '-' }}</span></td>
              <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.detail || '-'" :title="item.detail || '-'">{{ props.formatLogDetail(item.detail) }}</span></td>
            </tr>
            <tr v-if="!props.logs.length">
              <td colspan="5" class="muted">暂无日志数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="pagination">
      <button class="btn btn-secondary" :disabled="props.logPage <= 1" @click="props.changeLogPage(props.logPage - 1)">上一页</button>
      <button class="btn btn-secondary" :disabled="props.logPage * props.pageSize >= props.logCount" @click="props.changeLogPage(props.logPage + 1)">下一页</button>
      <span class="muted">第 {{ props.logPage }} 页 / 共 {{ props.logCount }} 条</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AdminLogItem } from "@/types/api";
import type { LogSubPage } from "@/utils/adminLog";

interface Props {
  logs: AdminLogItem[];
  pageSize: number;
  logPage: number;
  logCount: number;
  currentLogTitle: string;
  logSubPage: LogSubPage;
  logSubPageOptions: ReadonlyArray<{ value: LogSubPage; label: string }>;
  logActionOptions: string[];
  logKeyword: string;
  logActorName: string;
  logAction: string;
  logFileExt: string;
  logSharerName: string;
  logSaverName: string;
  logDay: string;
  switchLogSubPage: (next: LogSubPage) => void;
  resetLogFilters: () => void;
  reloadLogs: () => void;
  changeLogPage: (nextPage: number) => void;
  formatDateTime: (value: string) => string;
  formatActorDisplay: (actorName: string, actorIdentity: string) => string;
  formatActionLabel: (action: string) => string;
  formatLogDetail: (detail: string) => string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  "update:logKeyword": [value: string];
  "update:logActorName": [value: string];
  "update:logAction": [value: string];
  "update:logFileExt": [value: string];
  "update:logSharerName": [value: string];
  "update:logSaverName": [value: string];
  "update:logDay": [value: string];
}>();

function onLogKeywordInput(event: Event) {
  emit("update:logKeyword", (event.target as HTMLInputElement).value.trim());
}

function onLogActorNameInput(event: Event) {
  emit("update:logActorName", (event.target as HTMLInputElement).value.trim());
}

function onLogActionChange(event: Event) {
  emit("update:logAction", (event.target as HTMLSelectElement).value);
}

function onLogFileExtInput(event: Event) {
  emit("update:logFileExt", (event.target as HTMLInputElement).value.trim());
}

function onLogSharerNameInput(event: Event) {
  emit("update:logSharerName", (event.target as HTMLInputElement).value.trim());
}

function onLogSaverNameInput(event: Event) {
  emit("update:logSaverName", (event.target as HTMLInputElement).value.trim());
}

function onLogDayInput(event: Event) {
  emit("update:logDay", (event.target as HTMLInputElement).value);
}
</script>
