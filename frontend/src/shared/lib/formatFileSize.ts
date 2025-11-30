/**
 * Форматировать размер файла в человекочитаемый формат
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Б';

  const k = 1024;
  const sizes = ['Б', 'КБ', 'МБ', 'ГБ', 'ТБ'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Проверить допустимый размер файла
 */
export function isFileSizeValid(bytes: number, maxSizeMB: number): boolean {
  return bytes <= maxSizeMB * 1024 * 1024;
}

/**
 * Максимальный размер файла по умолчанию (10 МБ)
 */
export const MAX_FILE_SIZE = 10 * 1024 * 1024;

