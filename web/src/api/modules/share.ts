import http from '@/api/http';
import type { ResourceInfoResponse } from '@/types/api';

export function resourceInfoApi(identity: string) {
  return http.get<never, ResourceInfoResponse>('/resource/info', {
    params: { identity },
    skipAuth: true,
  });
}

export function resourceSaveApi(payload: { repository_identity: string; parent_id: number; share_identity?: string }) {
  return http.post('/resource/save', payload);
}
