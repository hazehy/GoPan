<template>
  <div class="container auth-container">
    <div class="card login-card">
      <div class="login-head">
        <div class="login-badge">GoPan</div>
        <h2 class="auth-title">欢迎回来</h2>
        <p class="login-subtitle">登录后继续管理你的文件与分享</p>
      </div>

      <form class="form login-form" novalidate @submit.prevent="onSubmit">
        <input
          class="input"
          v-model.trim="form.name"
          type="text"
          placeholder="用户名"
          required
        />
        <input
          class="input"
          v-model="form.password"
          type="password"
          placeholder="密码"
          required
        />
        <button class="btn btn-primary" type="submit" :disabled="loading">
          {{ loading ? "登录中..." : "登录" }}
        </button>
        <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
      </form>

      <div class="login-footer">
        <p class="muted">
          没有账号？<router-link to="/register">去注册</router-link>
        </p>
      </div>

      <div class="login-decoration" aria-hidden="true"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { loginApi } from "@/api/modules/auth";
import { useAuthStore } from "@/stores/auth";
import { validatePassword, validateUsername } from "@/utils/validators";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const form = reactive({
  name: "",
  password: "",
});
const loading = ref(false);
const errorMessage = ref("");
const lastSubmitAt = ref(0);
const SUBMIT_GUARD_MS = 1200;

async function onSubmit() {
  if (loading.value) {
    return;
  }
  const now = Date.now();
  if (now - lastSubmitAt.value < SUBMIT_GUARD_MS) {
    errorMessage.value = "操作过于频繁，请稍后再试";
    return;
  }

  const usernameError = validateUsername(form.name);
  if (usernameError) {
    errorMessage.value = usernameError;
    return;
  }
  const passwordError = validatePassword(form.password);
  if (passwordError) {
    errorMessage.value = passwordError;
    return;
  }

  try {
    lastSubmitAt.value = now;
    loading.value = true;
    errorMessage.value = "";
    const res = await loginApi(form);
    authStore.setTokens(res.token, res.refresh_token, res.role);

    const redirect =
      typeof route.query.redirect === "string"
        ? route.query.redirect
        : res.role === 2
          ? "/admin"
          : "/disk";
    await router.replace(redirect);
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    loading.value = false;
  }
}
</script>
