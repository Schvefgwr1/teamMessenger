import { AxiosError } from 'axios';
import { toast } from 'sonner';

export interface ApiError {
  error: string;
  message?: string;
}

/**
 * Обработчик ошибок API
 * Показывает toast с сообщением об ошибке
 */
export function handleApiError(error: unknown): void {
  if (error instanceof AxiosError) {
    const apiError = error.response?.data as ApiError;
    const message = apiError?.error || apiError?.message || 'Произошла ошибка';
    
    switch (error.response?.status) {
      case 400:
        toast.error(`Некорректный запрос: ${message}`);
        break;
      case 401:
        toast.error('Необходима авторизация');
        break;
      case 403:
        toast.error('Доступ запрещён');
        break;
      case 404:
        toast.error('Ресурс не найден');
        break;
      case 429:
        toast.error('Слишком много запросов. Попробуйте позже');
        break;
      case 500:
        toast.error('Внутренняя ошибка сервера');
        break;
      default:
        toast.error(message);
    }
  } else {
    toast.error('Неизвестная ошибка');
  }
}

/**
 * Извлечь сообщение об ошибке из ответа API
 */
export function getErrorMessage(error: unknown): string {
  if (error instanceof AxiosError) {
    const apiError = error.response?.data as ApiError;
    return apiError?.error || apiError?.message || 'Произошла ошибка';
  }
  return 'Неизвестная ошибка';
}

