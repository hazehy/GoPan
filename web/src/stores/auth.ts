import { defineStore } from 'pinia';

const TOKEN_KEY = 'gopan_token';
const REFRESH_TOKEN_KEY = 'gopan_refresh_token';
const ROLE_KEY = 'gopan_role';

interface AuthState {
  token: string;
  refreshToken: string;
  role: number;
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: localStorage.getItem(TOKEN_KEY) ?? '',
    refreshToken: localStorage.getItem(REFRESH_TOKEN_KEY) ?? '',
    role: Number(localStorage.getItem(ROLE_KEY) ?? 0),
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    isAdmin: (state) => state.role === 2,
  },
  actions: {
    setTokens(token: string, refreshToken: string, role: number) {
      this.token = token;
      this.refreshToken = refreshToken;
      this.role = role;
      localStorage.setItem(TOKEN_KEY, token);
      localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
      localStorage.setItem(ROLE_KEY, String(role));
    },
    clearAuth() {
      this.token = '';
      this.refreshToken = '';
      this.role = 0;
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem(ROLE_KEY);
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('gopan-auth-cleared'));
      }
    },
  },
});
