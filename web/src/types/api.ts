export interface ApiError {
  message: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  role: number;
}

export interface AdminOverviewResponse {
  total_users: number;
  active_users: number;
  disabled_users: number;
  total_files: number;
  total_folders: number;
  total_file_size: number;
  today_uploads: number;
  today_registers: number;
  ext_stats: Array<{
    ext: string;
    count: number;
  }>;
}

export interface AdminUserItem {
  identity: string;
  name: string;
  email: string;
  status: number;
  role: number;
  last_login_at: string;
  created_at: string;
}

export interface AdminUserListResponse {
  list: AdminUserItem[];
  count: number;
}

export interface AdminFileItem {
  identity: string;
  parent_id: number;
  user_identity: string;
  user_name: string;
  repository_identity: string;
  name: string;
  ext: string;
  path: string;
  size: number;
  updated_at: string;
}

export interface AdminFileListResponse {
  list: AdminFileItem[];
  count: number;
}

export interface AdminLogItem {
  identity: string;
  actor_identity: string;
  actor_name: string;
  actor_role: number;
  action: string;
  target_type: string;
  target_identity: string;
  detail: string;
  created_at: string;
}

export interface AdminLogListResponse {
  list: AdminLogItem[];
  count: number;
}

export interface UserDetailResponse {
  name: string;
  email: string;
}

export interface UserFile {
  id: number;
  identity: string;
  repository_identity: string;
  name: string;
  ext: string;
  path: string;
  size: number;
  updated_at: string;
}

export interface FileListResponse {
  list: UserFile[];
  count: number;
}

export interface FilePreUploadResponse {
  identity: string;
  upload_id: string;
  key: string;
}

export interface FileChunkUploadResponse {
  etag: string;
}

export interface FileUploadResponse {
  identity: string;
  ext: string;
  name: string;
}

export interface FileDownloadResponse {
  url: string;
}

export interface ShareCreateResponse {
  identity: string;
}

export interface ResourceInfoResponse {
  repository_identity: string;
  name: string;
  ext: string;
  size: number;
  path: string;
}
