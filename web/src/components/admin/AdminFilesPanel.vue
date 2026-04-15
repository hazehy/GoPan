<template>
  <div class="card admin-section">
    <div class="admin-section-head">
      <h3>全文件管理</h3>
      <div class="admin-actions-inline">
        <input class="input admin-filter-input" :value="props.fileKeyword" @input="onFileKeywordInput" placeholder="按名称/标识搜索" />
        <input class="input admin-filter-input" :value="props.fileUserName" @input="onFileUserNameInput" placeholder="按用户名过滤" />
        <button class="btn btn-secondary" @click="props.reloadFiles">查询</button>
      </div>
    </div>
    <div class="x-scroll-panel">
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
          <tr v-for="item in props.files" :key="item.identity">
            <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="`${item.name}${item.ext || ''}`" :title="`${item.name}${item.ext || ''}`">{{ item.name }}{{ item.ext || '' }}</span></td>
            <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="item.user_name || '-'" :title="item.user_name || '-'">{{ item.user_name || '-' }}</span></td>
            <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="props.formatText(item.path)" :title="props.formatText(item.path)">{{ props.formatText(item.path) }}</span></td>
            <td>{{ item.repository_identity ? '文件' : '文件夹' }}</td>
            <td>{{ props.formatFileSize(item.size) }}</td>
            <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="props.formatText(item.updated_at)" :title="props.formatText(item.updated_at)">{{ props.formatDateTime(item.updated_at) }}</span></td>
            <td>
              <button class="btn btn-danger" :disabled="props.fileDeleteLoading" @click="props.deleteFile(item.identity)">删除</button>
            </td>
          </tr>
          <tr v-if="!props.files.length">
            <td colspan="7" class="muted">暂无文件数据</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="pagination">
      <button class="btn btn-secondary" :disabled="props.filePage <= 1" @click="props.changeFilePage(props.filePage - 1)">上一页</button>
      <button class="btn btn-secondary" :disabled="props.filePage * props.pageSize >= props.fileCount" @click="props.changeFilePage(props.filePage + 1)">下一页</button>
      <span class="muted">第 {{ props.filePage }} 页 / 共 {{ props.fileCount }} 条</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AdminFileItem } from "@/types/api";

interface Props {
  files: AdminFileItem[];
  pageSize: number;
  filePage: number;
  fileCount: number;
  fileKeyword: string;
  fileUserName: string;
  fileDeleteLoading: boolean;
  reloadFiles: () => void;
  changeFilePage: (nextPage: number) => void;
  deleteFile: (identity: string) => void;
  formatText: (value: string) => string;
  formatDateTime: (value: string) => string;
  formatFileSize: (value: number) => string;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  "update:fileKeyword": [value: string];
  "update:fileUserName": [value: string];
}>();

function onFileKeywordInput(event: Event) {
  emit("update:fileKeyword", (event.target as HTMLInputElement).value.trim());
}

function onFileUserNameInput(event: Event) {
  emit("update:fileUserName", (event.target as HTMLInputElement).value.trim());
}
</script>
