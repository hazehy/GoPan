export type LogSubPage = "login" | "upload" | "share-create" | "share-save" | "others";

export interface LogQueryFilterState {
  keyword: string;
  actorName: string;
  action: string;
  fileExt: string;
  sharerName: string;
  saverName: string;
  day: string;
}

export function getCurrentLogTitle(subPage: LogSubPage) {
  if (subPage === "login") {
    return "筛选项：用户名 + 日期";
  }
  if (subPage === "upload") {
    return "筛选项：文件类型 + 日期";
  }
  if (subPage === "share-create") {
    return "筛选项：分享者用户名 + 文件类型 + 日期";
  }
  if (subPage === "share-save") {
    return "筛选项：分享者用户名 + 保存者用户名 + 文件类型 + 日期";
  }
  return "筛选项：操作人 + 类型 + 关键字 + 日期";
}

// Build a mode-aware query payload: each log tab maps to a different backend filter set.
export function buildLogQueryParams(subPage: LogSubPage, state: LogQueryFilterState) {
  const base = {
    keyword: "",
    action: "",
    actor_name: "",
    file_ext: "",
    sharer_name: "",
    saver_name: "",
    day: state.day,
  };

  if (subPage === "login") {
    return {
      ...base,
      action: "USER_LOGIN",
      actor_name: state.actorName,
    };
  }
  if (subPage === "upload") {
    return {
      ...base,
      action: "FILE_UPLOAD",
      file_ext: state.fileExt,
    };
  }
  if (subPage === "share-create") {
    return {
      ...base,
      action: "SHARE_CREATE",
      actor_name: state.sharerName,
      sharer_name: state.sharerName,
      file_ext: state.fileExt,
    };
  }
  if (subPage === "share-save") {
    return {
      ...base,
      action: "SHARE_RESOURCE_SAVE",
      sharer_name: state.sharerName,
      saver_name: state.saverName,
      file_ext: state.fileExt,
    };
  }

  return {
    ...base,
    keyword: state.keyword,
    action: state.action,
    actor_name: state.actorName,
  };
}

export function formatDateTime(value: string) {
  const raw = formatText(value);
  if (raw === "-") {
    return raw;
  }

  const normalized = raw.replace("T", " ").replace("Z", "").split(".")[0];
  const matched = normalized.match(/^(\d{4}-\d{2}-\d{2})\s(\d{2}:\d{2}:\d{2})/);
  if (matched) {
    return `${matched[1]} ${matched[2]}`;
  }

  return normalized.replace(/\s[+-]\d{4}\s[A-Za-z]+$/, "");
}

export function formatText(value: string) {
  return value || "-";
}

export function formatActorDisplay(actorName: string, actorIdentity: string) {
  const name = (actorName || "").trim();
  if (name) {
    return name;
  }
  const identity = (actorIdentity || "").trim();
  if (!identity) {
    return "-";
  }
  if (identity.length <= 12) {
    return identity;
  }
  return `${identity.slice(0, 8)}...`;
}

export function formatActionLabel(action: string) {
  const map: Record<string, string> = {
    USER_REGISTER: "用户注册",
    USER_LOGIN: "用户登录",
    FILE_UPLOAD: "文件上传",
    SHARE_CREATE: "创建分享链接",
    FILE_SAVE_REPOSITORY: "保存文件到网盘",
    SHARE_RESOURCE_SAVE: "保存分享文件",
    USER_STATUS_UPDATE: "用户状态更新",
  };
  return map[action] || action;
}

export function formatLogDetail(detail: string) {
  const text = (detail || "").trim();
  if (!text) {
    return "-";
  }
  if (!text.includes("=")) {
    return text;
  }

  const map: Record<string, string> = {
    file_ext: "文件类型",
    upload_mode: "上传模式",
    file_name: "文件名",
    sharer_name: "分享者",
    saver_name: "保存者",
    expires_days: "有效期(天)",
    share_identity: "分享标识",
    saved_name: "保存名",
  };

  return text
    .split(";")
    .map((item) => item.trim())
    .filter(Boolean)
    .map((item) => {
      const index = item.indexOf("=");
      if (index <= 0) {
        return item;
      }
      const key = item.slice(0, index);
      const value = item.slice(index + 1);
      return `${map[key] || key}: ${value || "-"}`;
    })
    .join(" | ");
}
