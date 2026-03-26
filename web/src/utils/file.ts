export function splitNameAndExt(filename: string): { name: string; ext: string } {
  const index = filename.lastIndexOf('.');
  if (index <= 0) {
    return { name: filename, ext: '' };
  }
  return {
    name: filename.slice(0, index),
    ext: filename.slice(index),
  };
}

export function formatFileSize(size: number): string {
  if (size === 0) {
    return '0 B';
  }
  if (!size) {
    return '-';
  }
  const units = ['B', 'KB', 'MB', 'GB'];
  let value = size;
  let index = 0;
  while (value >= 1024 && index < units.length - 1) {
    value /= 1024;
    index += 1;
  }
  return `${value.toFixed(value >= 10 || index === 0 ? 0 : 1)} ${units[index]}`;
}
