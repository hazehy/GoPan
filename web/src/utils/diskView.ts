import type { UserFile } from "@/types/api";

const imageExtSet = new Set([".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"]);
const videoExtSet = new Set([".mp4", ".avi", ".mkv", ".mov", ".flv", ".wmv"]);
const audioExtSet = new Set([".mp3", ".wav", ".flac", ".aac", ".ogg"]);
const archiveExtSet = new Set([".zip", ".rar", ".7z", ".tar", ".gz"]);
const documentExtSet = new Set([".doc", ".docx", ".txt", ".md", ".rtf"]);
const sheetExtSet = new Set([".xls", ".xlsx", ".csv"]);
const slideExtSet = new Set([".ppt", ".pptx"]);

interface UpdatedAtParts {
  year: number;
  month: number;
  day: number;
  hour: number;
  minute: number;
  second: number;
}

export function isFolder(item: UserFile) {
  return !item.repository_identity;
}

export function getFileIcon(item: UserFile) {
  if (isFolder(item)) {
    return "📁";
  }

  const ext = (item.ext ?? "").toLowerCase();
  if (imageExtSet.has(ext)) {
    return "🖼️";
  }
  if (videoExtSet.has(ext)) {
    return "🎬";
  }
  if (audioExtSet.has(ext)) {
    return "🎵";
  }
  if (archiveExtSet.has(ext)) {
    return "🗜️";
  }
  if (documentExtSet.has(ext)) {
    return "📄";
  }
  if (sheetExtSet.has(ext)) {
    return "📊";
  }
  if (slideExtSet.has(ext)) {
    return "📽️";
  }
  if (ext === ".pdf") {
    return "📕";
  }

  return "📦";
}

export function parseUpdatedAtParts(value: string): UpdatedAtParts | null {
  const normalized = value.replace("T", " ").replace("Z", "").split(".")[0];
  const matched = normalized.match(
    /^(\d{4})-(\d{2})-(\d{2})\s(\d{2}):(\d{2})(?::(\d{2}))?/,
  );
  if (!matched) {
    return null;
  }

  return {
    year: Number(matched[1]),
    month: Number(matched[2]),
    day: Number(matched[3]),
    hour: Number(matched[4]),
    minute: Number(matched[5]),
    second: Number(matched[6] ?? "0"),
  };
}

export function parseUpdatedAtToMillis(value: string) {
  const parsed = parseUpdatedAtParts(value);
  if (parsed) {
    return new Date(
      parsed.year,
      parsed.month - 1,
      parsed.day,
      parsed.hour,
      parsed.minute,
      parsed.second,
    ).getTime();
  }

  return Date.parse(value || "") || 0;
}

export function formatUpdatedAt(value: string) {
  if (!value) {
    return "-";
  }

  const parsed = parseUpdatedAtParts(value);
  if (!parsed) {
    return value;
  }

  const pad = (num: number) => String(num).padStart(2, "0");
  return `${parsed.year}-${pad(parsed.month)}-${pad(parsed.day)} ${pad(parsed.hour)}:${pad(parsed.minute)}`;
}
