import http from '@/api/http';
import type { LoginResponse, UserDetailResponse } from '@/types/api';

export function loginApi(payload: { name: string; password: string }) {
  return http.post<never, LoginResponse>('/user/login', payload, { skipAuth: true });
}

export function registerApi(payload: { name: string; password: string; email: string; code: string }) {
  return http.post('/register', payload, { skipAuth: true });
}

export function sendCodeApi(payload: { email: string }) {
  return http.post('/code/send', payload, { skipAuth: true });
}

export function userDetailApi(identity: string) {
  return http.get<never, UserDetailResponse>('/user/detail', {
    params: { identity },
    skipAuth: true,
  });
}
