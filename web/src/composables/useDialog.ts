import { reactive } from 'vue';

type DialogType = 'alert' | 'confirm' | 'prompt';

interface DialogState {
  visible: boolean;
  type: DialogType;
  title: string;
  message: string;
  confirmText: string;
  cancelText: string;
  placeholder: string;
  inputValue: string;
  promptError: string;
  promptValidator?: (value: string) => string;
  resolve?: (value: string | boolean | null) => void;
}

const state = reactive<DialogState>({
  visible: false,
  type: 'alert',
  title: '提示',
  message: '',
  confirmText: '确定',
  cancelText: '取消',
  placeholder: '',
  inputValue: '',
  promptError: '',
  promptValidator: undefined,
  resolve: undefined,
});

function openDialog(options: {
  type: DialogType;
  title?: string;
  message: string;
  defaultValue?: string;
  placeholder?: string;
  confirmText?: string;
  cancelText?: string;
  promptValidator?: (value: string) => string;
}) {
  return new Promise<string | boolean | null>((resolve) => {
    state.type = options.type;
    state.title = options.title ?? '提示';
    state.message = options.message;
    state.confirmText = options.confirmText ?? '确定';
    state.cancelText = options.cancelText ?? '取消';
    state.placeholder = options.placeholder ?? '';
    state.inputValue = options.defaultValue ?? '';
    state.promptError = '';
    state.promptValidator = options.promptValidator;
    state.visible = true;
    state.resolve = resolve;
  });
}

export function useDialogState() {
  return state;
}

export function closeDialog(result: string | boolean | null) {
  state.visible = false;
  state.promptError = '';
  state.promptValidator = undefined;
  state.resolve?.(result);
  state.resolve = undefined;
}

export function setDialogPromptError(message: string) {
  state.promptError = message;
}

export async function alertDialog(message: string, title?: string) {
  await openDialog({
    type: 'alert',
    title,
    message,
  });
}

export async function confirmDialog(message: string, title?: string) {
  const result = await openDialog({
    type: 'confirm',
    title,
    message,
  });
  return result === true;
}

export async function promptDialog(
  message: string,
  options?: {
    title?: string;
    defaultValue?: string;
    placeholder?: string;
    confirmText?: string;
    cancelText?: string;
    promptValidator?: (value: string) => string;
  },
) {
  const result = await openDialog({
    type: 'prompt',
    message,
    title: options?.title,
    defaultValue: options?.defaultValue,
    placeholder: options?.placeholder,
    confirmText: options?.confirmText,
    cancelText: options?.cancelText,
    promptValidator: options?.promptValidator,
  });

  if (typeof result !== 'string') {
    return null;
  }

  return result.trim();
}
