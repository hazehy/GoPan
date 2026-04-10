import http from '@/api/http';
import type {
  FileChunkUploadResponse,
  FileDownloadResponse,
  FileListResponse,
  FilePreUploadResponse,
  FileUploadResponse,
  ShareCreateResponse,
} from '@/types/api';

export function fileListApi(params: { id: number; page: number; size: number }) {
  return http.get<never, FileListResponse>('/file/list', { params });
}

export function folderCreateApi(payload: { name: string; parent_id: number }) {
  return http.post('/folder/create', payload);
}

export function fileRenameApi(payload: { identity: string; name: string }) {
  return http.post('/file/rename', payload);
}

export function fileMoveApi(payload: { identity: string; parent_identity: string }) {
  return http.put('/file/move', payload);
}

export function fileDeleteApi(identity: string) {
  return http.delete('/file/delete', {
    data: { identity },
  });
}

export function fileDownloadApi(payload: { repository_identity: string; filename: string }) {
  return http.get<never, FileDownloadResponse>('/file/download', {
    params: payload,
  });
}

export function shareCreateApi(payload: { repository_identity: string; expires: number }) {
  return http.post<never, ShareCreateResponse>('/share/create', payload);
}

export function filePreUploadApi(payload: { md5: string; name: string; ext: string }) {
  return http.post<never, FilePreUploadResponse>('/file/preupload', payload);
}

export function fileChunkUploadApi(formData: FormData, signal?: AbortSignal) {
  return http.post<never, FileChunkUploadResponse>('/file/chunkupload', formData, {
    signal,
  });
}

export function fileChunkUploadCompleteApi(payload: {
  key: string;
  upload_id: string;
  cos_objects: Array<{ part_number: number; etag: string }>;
}, signal?: AbortSignal) {
  return http.post('/file/chunkupload/complete', payload, { signal });
}

export function fileUploadApi(payload: {
  hash: string;
  name: string;
  ext: string;
  size: number;
  path?: string;
  key?: string;
}, signal?: AbortSignal) {
  return http.post<never, FileUploadResponse>('/file/upload', payload, { signal });
}

export function userRepositoryApi(payload: {
  parent_id: number;
  repository_identity: string;
  ext: string;
  name: string;
}, signal?: AbortSignal) {
  return http.post('/user/repository', payload, { signal });
}
