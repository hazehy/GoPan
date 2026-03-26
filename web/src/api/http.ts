import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/stores/auth';

interface RequestConfig extends InternalAxiosRequestConfig {
  skipAuth?: boolean;
}

const http = axios.create({
  baseURL: '/api',
  timeout: 120000,
});

let refreshing = false;
let pendingQueue: Array<(token: string) => void> = [];

function flushQueue(token: string) {
  pendingQueue.forEach((resolver) => resolver(token));
  pendingQueue = [];
}

async function refreshToken(refreshTokenValue: string) {
  const response = await axios.post('/api/token/refresh', {}, {
    headers: {
      Authorization: refreshTokenValue,
    },
    timeout: 120000,
  });
  return response.data as { token: string; refresh_token: string; role: number };
}

http.interceptors.request.use((config: RequestConfig) => {
  const authStore = useAuthStore();
  if (!config.skipAuth && authStore.token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = authStore.token;
  }
  return config;
});

http.interceptors.response.use(
  (response) => response.data,
  async (error: AxiosError) => {
    const authStore = useAuthStore();
    const originalRequest = (error.config ?? {}) as RequestConfig & { _retry?: boolean };

    if (error.response?.status !== 401 || originalRequest._retry || originalRequest.skipAuth) {
      const message =
        (error.response?.data as { message?: string; error?: string } | string | undefined) ||
        error.message ||
        '请求失败';
      return Promise.reject(typeof message === 'string' ? message : message.message || message.error || '请求失败');
    }

    if (!authStore.refreshToken) {
      authStore.clearAuth();
      return Promise.reject('登录已过期，请重新登录');
    }

    if (refreshing) {
      return new Promise((resolve) => {
        pendingQueue.push((newToken: string) => {
          originalRequest.headers = originalRequest.headers ?? {};
          originalRequest.headers.Authorization = newToken;
          resolve(http(originalRequest));
        });
      });
    }

    originalRequest._retry = true;
    refreshing = true;

    try {
      const tokenData = await refreshToken(authStore.refreshToken);
      authStore.setTokens(tokenData.token, tokenData.refresh_token, tokenData.role);
      flushQueue(tokenData.token);
      originalRequest.headers = originalRequest.headers ?? {};
      originalRequest.headers.Authorization = tokenData.token;
      return http(originalRequest);
    } catch (refreshErr) {
      authStore.clearAuth();
      return Promise.reject(refreshErr instanceof Error ? refreshErr.message : '登录已过期，请重新登录');
    } finally {
      refreshing = false;
    }
  },
);

export default http;
