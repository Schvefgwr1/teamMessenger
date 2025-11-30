import { API_BASE_URL } from '@/shared/api';
import type { File } from '@/shared/types';

/**
 * Получить полный URL файла
 * File.url может хранить относительный путь, нужно добавить базовый URL
 */
export function getFileUrl(file: File | null | undefined): string | undefined {
  if (!file?.url) return undefined;
  
  // Если URL уже абсолютный
  if (file.url.startsWith('http://') || file.url.startsWith('https://')) {
    return file.url;
  }
  
  // Добавляем базовый URL
  return `${API_BASE_URL}${file.url}`;
}

/**
 * Получить URL аватара
 */
export function getAvatarUrl(file: File | null | undefined): string | undefined {
  return getFileUrl(file);
}

