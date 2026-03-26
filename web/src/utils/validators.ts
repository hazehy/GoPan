export function isPositiveIntegerText(input: string) {
  return /^\d+$/.test(input);
}

export function validateUsername(name: string) {
  const value = name.trim();
  if (!value) {
    return '用户名不能为空';
  }
  if (value.length < 2 || value.length > 20) {
    return '用户名长度需在 2~20 位';
  }
  if (!/^[\u4e00-\u9fa5a-zA-Z0-9_]+$/.test(value)) {
    return '用户名仅支持中文、字母、数字和下划线';
  }
  return '';
}

export function validatePassword(password: string) {
  if (!password) {
    return '密码不能为空';
  }
  if (password.length < 6 || password.length > 32) {
    return '密码长度需在 6~32 位';
  }
  return '';
}

export function validateEmail(email: string) {
  const value = email.trim();
  if (!value) {
    return '邮箱不能为空';
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
    return '邮箱格式不正确';
  }
  return '';
}

export function validateCode(code: string) {
  const value = code.trim();
  if (!value) {
    return '验证码不能为空';
  }
  return '';
}

export function validateFileOrFolderName(name: string) {
  const value = name.trim();
  if (!value) {
    return '名称不能为空';
  }
  if (value.length > 100) {
    return '名称长度不能超过 100 个字符';
  }
  if (/[\\/:*?"<>|]/.test(value)) {
    return '名称不能包含 \\ / : * ? " < > |';
  }
  return '';
}

export function validateShareExpiresDaysText(input: string) {
  const value = input.trim();
  if (!value) {
    return '请输入有效期（天）';
  }
  if (!isPositiveIntegerText(value)) {
    return '有效期必须为正整数天';
  }

  const days = Number(value);
  if (!Number.isInteger(days) || days <= 0) {
    return '有效期必须大于 0 天';
  }
  if (days > 3650) {
    return '有效期不能超过 3650 天';
  }
  return '';
}
