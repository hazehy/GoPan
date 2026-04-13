<template>
  <div class="card admin-section">
    <div class="admin-section-head">
      <h3>用户管理</h3>
      <div class="admin-actions-inline">
        <input class="input admin-filter-input" :value="props.userKeyword" @input="onKeywordInput" placeholder="按用户名/邮箱/标识搜索" />
        <button class="btn btn-secondary" @click="props.reloadUsers">查询</button>
      </div>
    </div>
    <table class="table admin-table">
      <colgroup>
        <col class="admin-col-user-name" />
        <col class="admin-col-user-email" />
        <col class="admin-col-user-role" />
        <col class="admin-col-user-status" />
        <col class="admin-col-user-status" />
        <col class="admin-col-user-status" />
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
          <th>上传权限</th>
          <th>下载权限</th>
          <th>分享权限</th>
          <th>最近登录</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in props.users" :key="user.identity">
          <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="user.name" :title="user.name">{{ user.name }}</span></td>
          <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="user.email" :title="user.email">{{ user.email }}</span></td>
          <td>{{ user.role === 2 ? '管理员' : '普通用户' }}</td>
          <td>{{ user.status === 1 ? '正常' : '禁用' }}</td>
          <td>
            <button
              class="btn btn-secondary"
              :disabled="user.role === 2 || props.userStatusLoading"
              @click="props.toggleUserPermission(user.identity, 'upload_permission', user.upload_permission)"
            >
              {{ user.upload_permission === 1 ? '允许' : '禁止' }}
            </button>
          </td>
          <td>
            <button
              class="btn btn-secondary"
              :disabled="user.role === 2 || props.userStatusLoading"
              @click="props.toggleUserPermission(user.identity, 'download_permission', user.download_permission)"
            >
              {{ user.download_permission === 1 ? '允许' : '禁止' }}
            </button>
          </td>
          <td>
            <button
              class="btn btn-secondary"
              :disabled="user.role === 2 || props.userStatusLoading"
              @click="props.toggleUserPermission(user.identity, 'share_permission', user.share_permission)"
            >
              {{ user.share_permission === 1 ? '允许' : '禁止' }}
            </button>
          </td>
          <td><span class="admin-cell-ellipsis admin-tooltip" :data-tooltip="props.formatText(user.last_login_at)" :title="props.formatText(user.last_login_at)">{{ props.formatDateTime(user.last_login_at) }}</span></td>
          <td>
            <button
              class="btn btn-secondary"
              :disabled="user.role === 2 || props.userStatusLoading"
              @click="props.toggleUserStatus(user.identity, user.status)"
            >
              {{ user.status === 1 ? '禁用' : '启用' }}
            </button>
          </td>
        </tr>
        <tr v-if="!props.users.length">
          <td colspan="9" class="muted">暂无用户数据</td>
        </tr>
      </tbody>
    </table>
    <div class="pagination">
      <button class="btn btn-secondary" :disabled="props.userPage <= 1" @click="props.changeUserPage(props.userPage - 1)">上一页</button>
      <button class="btn btn-secondary" :disabled="props.userPage * props.pageSize >= props.userCount" @click="props.changeUserPage(props.userPage + 1)">下一页</button>
      <span class="muted">第 {{ props.userPage }} 页 / 共 {{ props.userCount }} 条</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AdminUserItem } from "@/types/api";

interface Props {
  users: AdminUserItem[];
  pageSize: number;
  userPage: number;
  userCount: number;
  userKeyword: string;
  userStatusLoading: boolean;
  reloadUsers: () => void;
  changeUserPage: (nextPage: number) => void;
  toggleUserStatus: (identity: string, currentStatus: number) => void;
  toggleUserPermission: (
    identity: string,
    field: 'upload_permission' | 'download_permission' | 'share_permission',
    currentValue: number,
  ) => void;
  formatText: (value: string) => string;
  formatDateTime: (value: string) => string;
}

const props = defineProps<Props>();
const emit = defineEmits<{ "update:userKeyword": [value: string] }>();

function onKeywordInput(event: Event) {
  emit("update:userKeyword", (event.target as HTMLInputElement).value.trim());
}
</script>
