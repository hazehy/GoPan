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
      </form>

      <div class="login-footer">
        <p class="muted">
          没有账号？<router-link to="/register">去注册</router-link>
        </p>
        <p class="muted mt-12">
          忘记密码？<router-link to="/forgot-password">找回密码</router-link>
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
import { alertDialog } from "@/composables/useDialog";
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
const lastSubmitAt = ref(0);
const SUBMIT_GUARD_MS = 1200;

async function onSubmit() {
  if (loading.value) {
    return;
  }
  const now = Date.now();
  if (now - lastSubmitAt.value < SUBMIT_GUARD_MS) {
    await alertDialog("操作过于频繁，请稍后再试", "提示");
    return;
  }

  const usernameError = validateUsername(form.name);
  if (usernameError) {
    await alertDialog(usernameError, "提示");
    return;
  }
  const passwordError = validatePassword(form.password);
  if (passwordError) {
    await alertDialog(passwordError, "提示");
    return;
  }

  try {
    lastSubmitAt.value = now;
    loading.value = true;
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
    await alertDialog(error instanceof Error ? error.message : String(error), "登录失败");
  } finally {
    loading.value = false;
  }
}
</script>
