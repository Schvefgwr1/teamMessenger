import { apiClient } from '@/shared/api';
import type { User, Role, Permission, File as ApiFile, UserBrief, UserSearchResponse } from '@/shared/types';

/**
 * Ответ API при получении пользователя
 * GET /api/v1/users/me, GET /api/v1/users/:id
 */
export interface GetUserResponse {
  file: ApiFile | null;
  user: User;
}

/**
 * Запрос на обновление профиля (JSON часть в FormData.data)
 */
export interface UpdateUserRequest {
  username?: string;
  description?: string;
  gender?: string;
  age?: number;
  roleID?: number; // Только для admin
}

/**
 * User API endpoints
 */
export const userApi = {
  /**
   * Получить текущего пользователя
   * GET /api/v1/users/me
   */
  getMe: () => apiClient.get<GetUserResponse>('/users/me'),

  /**
   * Обновить профиль текущего пользователя
   * PUT /api/v1/users/me (multipart/form-data)
   *
   * FormData:
   * - data: JSON string с UpdateUserRequest
   * - file: File (новый аватар, опционально)
   */
  updateMe: (formData: FormData) =>
    apiClient.put<GetUserResponse>('/users/me', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  /**
   * Получить пользователя по ID
   * GET /api/v1/users/:userId
   */
  getUserById: (userId: string) =>
    apiClient.get<GetUserResponse>(`/users/${userId}`),

  /**
   * Получить все системные разрешения (admin)
   * GET /api/v1/permissions
   */
  getPermissions: () => apiClient.get<Permission[]>('/permissions'),

  /**
   * Получить все роли (admin)
   * GET /api/v1/roles
   */
  getRoles: () => apiClient.get<Role[]>('/roles'),

  /**
   * Создать роль (admin)
   * POST /api/v1/roles
   */
  createRole: (data: { name: string; description?: string; permissionIds?: number[] }) =>
    apiClient.post('/roles', data),

  /**
   * Обновить permissions роли (admin)
   * PATCH /api/v1/roles/:roleId/permissions
   */
  updateRolePermissions: (roleId: number, permissionIds: number[]) =>
    apiClient.patch(`/roles/${roleId}/permissions`, { permission_ids: permissionIds }),

  /**
   * Удалить роль (admin)
   * DELETE /api/v1/roles/:roleId
   */
  deleteRole: (roleId: number) => apiClient.delete(`/roles/${roleId}`),

  /**
   * Изменить роль пользователя (admin)
   * PATCH /api/v1/users/:userId/role
   */
  updateUserRole: (userId: string, roleId: number) =>
    apiClient.patch(`/users/${userId}/role`, { role_id: roleId }),

  /**
   * Получить краткую информацию о пользователе с ролью в чате
   * GET /api/v1/users/:userId/brief?chatId=:chatId
   */
  getUserBrief: (userId: string, chatId: string) =>
    apiClient.get<UserBrief>(`/users/${userId}/brief`, {
      params: { chatId },
    }),

  /**
   * Поиск пользователей по имени или email
   * GET /api/v1/searches/users?q=:query&limit=:limit
   */
  searchUsers: (query: string, limit?: number) =>
    apiClient.get<UserSearchResponse>('/searches/users', {
      params: { q: query, limit: limit || 10 },
    }),
};

/**
 * Создать FormData для обновления профиля
 */
export function createUpdateUserFormData(
  data: UpdateUserRequest,
  avatar?: File
): FormData {
  const formData = new FormData();
  formData.append('data', JSON.stringify(data));
  if (avatar) {
    formData.append('file', avatar);
  }
  return formData;
}

