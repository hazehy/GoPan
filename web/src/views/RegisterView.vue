<template>
  <div class="container auth-container register">
    <div class="card register-card">
      <div class="login-head">
        <div class="login-badge">GoPan</div>
        <h2 class="auth-title">创建账号</h2>
        <p class="login-subtitle">注册后即可上传、管理并分享你的文件</p>
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
          v-model.trim="form.email"
          type="email"
          placeholder="邮箱"
          required
        />
        <div class="stack-row">
          <input
            class="input"
            v-model.trim="form.code"
            type="text"
            placeholder="验证码"
            required
          />
          <button
            class="btn btn-secondary"
            type="button"
            :disabled="codeCooldown > 0 || sendingCode"
            @click="sendCode"
          >
            {{
              codeCooldown > 0
                ? `${codeCooldown}s`
                : sendingCode
                  ? "发送中"
                  : "发送验证码"
            }}
          </button>
        </div>
        <input
          class="input"
          v-model="form.password"
          type="password"
          placeholder="密码"
          required
        />
        <input
          class="input"
          v-model="form.confirmPassword"
          type="password"
          placeholder="确认密码"
          required
        />
        <button class="btn btn-primary" type="submit" :disabled="loading">
          {{ loading ? "提交中..." : "注册" }}
        </button>
        <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
      </form>

      <div class="login-footer">
        <p class="muted">
          已有账号？<router-link to="/login">去登录</router-link>
        </p>
      </div>

      <div class="login-decoration" aria-hidden="true"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { registerApi, sendCodeApi } from "@/api/modules/auth";
import {
  validateCode,
  validateEmail,
  validatePassword,
  validateUsername,
} from "@/utils/validators";

const router = useRouter();

const form = reactive({
  name: "",
  email: "",
  code: "",
  password: "",
  confirmPassword: "",
});

const loading = ref(false);
const sendingCode = ref(false);
const codeCooldown = ref(0);
const errorMessage = ref("");
const lastSubmitAt = ref(0);
const lastSendCodeAt = ref(0);
const SUBMIT_GUARD_MS = 1200;
const SEND_CODE_GUARD_MS = 1000;

let timer: number | null = null;

function showValidationError(message: string) {
  errorMessage.value = message;
}

function startCooldown() {
  codeCooldown.value = 60;
  timer = window.setInterval(() => {
    codeCooldown.value -= 1;
    if (codeCooldown.value <= 0 && timer) {
      window.clearInterval(timer);
      timer = null;
    }
  }, 1000);
}

async function sendCode() {
  if (sendingCode.value || codeCooldown.value > 0) {
    return;
  }
  const emailError = validateEmail(form.email);
  if (emailError) {
    showValidationError(emailError);
    return;
  }
  const now = Date.now();
  if (now - lastSendCodeAt.value < SEND_CODE_GUARD_MS) {
    showValidationError("操作过于频繁，请稍后再试");
    return;
  }

  try {
    lastSendCodeAt.value = now;
    sendingCode.value = true;
    errorMessage.value = "";
    await sendCodeApi({ email: form.email });
    startCooldown();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    sendingCode.value = false;
  }
}

async function onSubmit() {
  if (loading.value) {
    return;
  }
  const now = Date.now();
  if (now - lastSubmitAt.value < SUBMIT_GUARD_MS) {
    showValidationError("操作过于频繁，请稍后再试");
    return;
  }

  const usernameError = validateUsername(form.name);
  if (usernameError) {
    showValidationError(usernameError);
    return;
  }
  const emailError = validateEmail(form.email);
  if (emailError) {
    showValidationError(emailError);
    return;
  }
  const codeError = validateCode(form.code);
  if (codeError) {
    showValidationError(codeError);
    return;
  }
  const passwordError = validatePassword(form.password);
  if (passwordError) {
    showValidationError(passwordError);
    return;
  }
  if (form.password !== form.confirmPassword) {
    showValidationError("两次输入的密码不一致");
    return;
  }

  try {
    lastSubmitAt.value = now;
    loading.value = true;
    errorMessage.value = "";
    await registerApi(form);
    form.confirmPassword = "";
    await router.replace("/login");
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    loading.value = false;
  }
}

onBeforeUnmount(() => {
  if (timer) {
    window.clearInterval(timer);
  }
});
</script>
