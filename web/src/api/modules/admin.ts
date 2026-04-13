import http from '@/api/http';
import type {
  AdminFileListResponse,
  AdminLogListResponse,
  AdminOverviewResponse,
  AdminUserListResponse,
} from '@/types/api';

export function adminOverviewApi() {
  return http.get<never, AdminOverviewResponse>('/admin/overview');
}

export function adminUserListApi(params: { page: number; size: number; keyword: string }) {
  return http.get<never, AdminUserListResponse>('/admin/users', { params });
}

export function adminUserStatusUpdateApi(payload: {
  identity: string;
  status?: number;
  upload_permission?: number;
  download_permission?: number;
  share_permission?: number;
}) {
  return http.put('/admin/user/status', payload);
}

export function adminFileListApi(params: {
  page: number;
  size: number;
  keyword: string;
  user_name: string;
}) {
  return http.get<never, AdminFileListResponse>('/admin/files', { params });
}

export function adminFileDeleteApi(identity: string) {
  return http.delete('/admin/file', { data: { identity } });
}

export function adminLogListApi(params: {
  page: number;
  size: number;
  keyword: string;
  action: string;
  actor_name: string;
  file_ext: string;
  sharer_name: string;
  saver_name: string;
  day: string;
}) {
  return http.get<never, AdminLogListResponse>('/admin/logs', { params });
}
