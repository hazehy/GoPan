<template>
  <div class="card">
    <div class="toolbar">
      <button
        class="btn btn-secondary"
        @click="props.goParent"
        :disabled="props.breadcrumbs.length === 0"
      >
        返回上级
      </button>
      <button class="btn btn-primary" @click="props.createFolder">
        新建文件夹
      </button>
      <button class="btn btn-secondary" @click="props.refreshList">刷新</button>
    </div>

    <div class="path-bar">
      <span class="muted">当前位置：</span>
      <a
        href="javascript:void(0)"
        :class="[
          'path-link',
          { 'path-link-active': props.breadcrumbs.length === 0 },
        ]"
        @click="props.goToPath(-1)"
      >
        根目录
      </a>
      <template
        v-for="(node, index) in props.breadcrumbs"
        :key="`${node.id}-${index}`"
      >
        <span class="path-sep">/</span>
        <a
          href="javascript:void(0)"
          :class="[
            'path-link',
            { 'path-link-active': index === props.breadcrumbs.length - 1 },
          ]"
          @click="props.goToPath(index)"
        >
          {{ node.name }}
        </a>
      </template>
    </div>

    <div class="list-controls">
      <input
        class="input list-control-input"
        :value="props.keyword"
        @input="onKeywordInput"
        type="text"
        placeholder="按名称筛选"
      />
      <select
        class="input list-control-select"
        :value="props.typeFilter"
        @change="onTypeFilterChange"
      >
        <option value="all">全部</option>
        <option value="folder">仅文件夹</option>
        <option value="file">仅文件</option>
      </select>
      <select
        class="input list-control-select"
        :value="props.sortBy"
        @change="onSortByChange"
      >
        <option value="updated">按更新时间</option>
        <option value="name">按名称</option>
        <option value="size">按大小</option>
      </select>
      <select
        class="input list-control-select"
        :value="props.sortOrder"
        @change="onSortOrderChange"
      >
        <option value="desc">降序</option>
        <option value="asc">升序</option>
      </select>
    </div>

    <div class="x-scroll-panel">
      <div class="table-scroll-content">
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
            <tr v-for="item in props.pagedDisplayFiles" :key="item.identity">
              <td>
                <span class="file-name-cell">
                  <span class="file-icon">{{ props.getFileIcon(item) }}</span>
                  <a
                    v-if="props.isFolder(item)"
                    href="javascript:void(0)"
                    @click="props.openFolder(item)"
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
              <td>{{ props.formatFileSize(item.size) }}</td>
              <td>{{ props.formatUpdatedAt(item.updated_at) }}</td>
              <td>
                <div class="action-group">
                  <button
                    class="btn btn-secondary"
                    @click="props.downloadFile(item)"
                    :disabled="props.isFolder(item) || !item.path"
                  >
                    下载
                  </button>
                  <button class="btn btn-secondary" @click="props.renameItem(item)">
                    重命名
                  </button>
                  <button class="btn btn-secondary" @click="props.moveItem(item)">
                    移动
                  </button>
                  <button
                    class="btn btn-secondary"
                    @click="props.createShare(item)"
                    :disabled="!item.repository_identity"
                  >
                    分享
                  </button>
                  <button class="btn btn-danger" @click="props.deleteItem(item)">
                    删除
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!props.pagedDisplayFiles.length">
              <td colspan="4" class="muted">当前目录暂无文件</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="pagination">
      <button
        class="btn btn-secondary"
        :disabled="props.page <= 1"
        @click="props.changePage(props.page - 1)"
      >
        上一页
      </button>
      <button
        class="btn btn-secondary"
        :disabled="props.page * props.size >= props.total"
        @click="props.changePage(props.page + 1)"
      >
        下一页
      </button>
      <span class="muted">第 {{ props.page }} 页 / 共 {{ props.total }} 条</span>
    </div>

    <p class="error" v-if="props.errorMessage">{{ props.errorMessage }}</p>
  </div>
</template>

<script setup lang="ts">
import type { UserFile } from "@/types/api";

type TypeFilter = "all" | "folder" | "file";
type SortBy = "updated" | "name" | "size";
type SortOrder = "asc" | "desc";

interface Props {
  breadcrumbs: Array<{ id: number; name: string }>;
  keyword: string;
  typeFilter: TypeFilter;
  sortBy: SortBy;
  sortOrder: SortOrder;
  page: number;
  size: number;
  total: number;
  errorMessage: string;
  pagedDisplayFiles: UserFile[];
  goParent: () => void;
  createFolder: () => void;
  refreshList: () => void;
  goToPath: (index: number) => void;
  changePage: (nextPage: number) => void;
  openFolder: (item: UserFile) => void;
  downloadFile: (item: UserFile) => void;
  renameItem: (item: UserFile) => void;
  moveItem: (item: UserFile) => void;
  createShare: (item: UserFile) => void;
  deleteItem: (item: UserFile) => void;
  isFolder: (item: UserFile) => boolean;
  getFileIcon: (item: UserFile) => string;
  formatUpdatedAt: (value: string) => string;
  formatFileSize: (value: number) => string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  "update:keyword": [value: string];
  "update:typeFilter": [value: TypeFilter];
  "update:sortBy": [value: SortBy];
  "update:sortOrder": [value: SortOrder];
}>();

function onKeywordInput(event: Event) {
  const value = (event.target as HTMLInputElement).value.trim();
  emit("update:keyword", value);
}

function onTypeFilterChange(event: Event) {
  emit("update:typeFilter", (event.target as HTMLSelectElement).value as TypeFilter);
}

function onSortByChange(event: Event) {
  emit("update:sortBy", (event.target as HTMLSelectElement).value as SortBy);
}

function onSortOrderChange(event: Event) {
  emit("update:sortOrder", (event.target as HTMLSelectElement).value as SortOrder);
}
</script>
