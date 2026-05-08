<template>
  <div class="container auth-container register">
    <div class="card register-card">
      <div class="login-head">
        <div class="login-badge">GoPan</div>
        <h2 class="auth-title">找回密码</h2>
        <p class="login-subtitle">通过邮箱验证码重置账号密码</p>
      </div>

      <form class="form login-form" novalidate @submit.prevent="onSubmit">
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
          placeholder="新密码"
          required
        />
        <input
          class="input"
          v-model="form.confirmPassword"
          type="password"
          placeholder="确认新密码"
          required
        />
        <button class="btn btn-primary" type="submit" :disabled="loading">
          {{ loading ? "提交中..." : "重置密码" }}
        </button>
      </form>

      <div class="login-footer">
        <p class="muted">
          想起密码了？<router-link to="/login">返回登录</router-link>
        </p>
      </div>

      <div class="login-decoration" aria-hidden="true"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { resetPasswordApi, sendResetCodeApi } from "@/api/modules/auth";
import { alertDialog } from "@/composables/useDialog";
import { validateCode, validateEmail, validatePassword } from "@/utils/validators";

const router = useRouter();

const form = reactive({
  email: "",
  code: "",
  password: "",
  confirmPassword: "",
});

const loading = ref(false);
const sendingCode = ref(false);
const codeCooldown = ref(0);
const lastSubmitAt = ref(0);
const lastSendCodeAt = ref(0);
const SUBMIT_GUARD_MS = 1200;
const SEND_CODE_GUARD_MS = 1000;

let timer: number | null = null;

async function showValidationError(message: string) {
  await alertDialog(message, "提示");
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
    await showValidationError(emailError);
    return;
  }

  const now = Date.now();
  if (now - lastSendCodeAt.value < SEND_CODE_GUARD_MS) {
    await showValidationError("操作过于频繁，请稍后再试");
    return;
  }

  try {
    lastSendCodeAt.value = now;
    sendingCode.value = true;
    await sendResetCodeApi({ email: form.email });
    startCooldown();
  } catch (error) {
    await alertDialog(error instanceof Error ? error.message : String(error), "发送失败");
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
    await showValidationError("操作过于频繁，请稍后再试");
    return;
  }

  const emailError = validateEmail(form.email);
  if (emailError) {
    await showValidationError(emailError);
    return;
  }
  const codeError = validateCode(form.code);
  if (codeError) {
    await showValidationError(codeError);
    return;
  }
  const passwordError = validatePassword(form.password);
  if (passwordError) {
    await showValidationError(passwordError);
    return;
  }
  if (form.password !== form.confirmPassword) {
    await showValidationError("两次输入的密码不一致");
    return;
  }

  try {
    lastSubmitAt.value = now;
    loading.value = true;
    await resetPasswordApi({
      email: form.email,
      code: form.code,
      password: form.password,
    });
    await router.replace("/login");
  } catch (error) {
    await alertDialog(error instanceof Error ? error.message : String(error), "重置失败");
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
