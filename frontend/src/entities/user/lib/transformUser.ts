import type { GetUserResponse } from '../api/userApi';
import type { AuthUser } from '@/entities/session';

/**
 * Преобразует ответ API (GetUserResponse) в формат AuthUser
 * для хранения в Auth Store
 */
export function transformUserResponse(response: GetUserResponse): AuthUser {
  return {
    ...response.user,
    avatar: response.file ?? undefined,
  };
}

