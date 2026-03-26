<template>
  <teleport to="body">
    <div v-if="dialog.visible" class="dialog-mask" @click="onCancel">
      <div class="dialog-panel" @click.stop>
        <h3 class="dialog-title">{{ dialog.title }}</h3>
        <p class="dialog-message">{{ dialog.message }}</p>

        <input
          v-if="dialog.type === 'prompt'"
          v-model="dialog.inputValue"
          class="input dialog-input"
          type="text"
          :placeholder="dialog.placeholder || '请输入内容'"
          @input="clearPromptError"
          @keydown.enter="onConfirm"
        />
        <p class="error" v-if="dialog.type === 'prompt' && dialog.promptError">
          {{ dialog.promptError }}
        </p>

        <div class="dialog-actions">
          <button
            v-if="dialog.type !== 'alert'"
            class="btn btn-secondary"
            type="button"
            @click="onCancel"
          >
            {{ dialog.cancelText }}
          </button>
          <button class="btn btn-primary" type="button" @click="onConfirm">
            {{ dialog.confirmText }}
          </button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script setup lang="ts">
import {
  closeDialog,
  setDialogPromptError,
  useDialogState,
} from "@/composables/useDialog";

const dialog = useDialogState();

function onConfirm() {
  if (dialog.type === "prompt") {
    const validator = dialog.promptValidator;
    if (validator) {
      const message = validator(dialog.inputValue);
      if (message) {
        setDialogPromptError(message);
        return;
      }
    }
    closeDialog(dialog.inputValue);
    return;
  }

  if (dialog.type === "confirm") {
    closeDialog(true);
    return;
  }

  closeDialog(true);
}

function onCancel() {
  if (dialog.type === "prompt") {
    closeDialog(null);
    return;
  }

  if (dialog.type === "confirm") {
    closeDialog(false);
    return;
  }

  closeDialog(true);
}

function clearPromptError() {
  if (dialog.promptError) {
    setDialogPromptError("");
  }
}
</script>
