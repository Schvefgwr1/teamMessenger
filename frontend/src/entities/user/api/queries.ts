import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userApi, UpdateUserRequest, createUpdateUserFormData } from './userApi';
import { transformUserResponse } from '../lib/transformUser';
import { useAuthStore } from '@/entities/session';
import { toast } from '@/shared/ui';

/**
 * Query keys для пользователей
 */
export const userKeys = {
  all: ['users'] as const,
  me: () => [...userKeys.all, 'me'] as const,
  detail: (userId: string) => [...userKeys.all, 'detail', userId] as const,
  brief: (userId: string, chatId: string) => [...userKeys.all, 'brief', userId, chatId] as const,
  roles: () => [...userKeys.all, 'roles'] as const,
  permissions: () => [...userKeys.all, 'permissions'] as const,
};

/**
 * Хук для получения текущего пользователя
 */
export function useCurrentUser() {
  const { isAuthenticated } = useAuthStore();

  return useQuery({
    queryKey: userKeys.me(),
    queryFn: async () => {
      const response = await userApi.getMe();
      return transformUserResponse(response.data);
    },
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 минут
  });
}

/**
 * Хук для получения пользователя по ID
 */
export function useUserById(userId: string | undefined) {
  return useQuery({
    queryKey: userKeys.detail(userId!),
    queryFn: async () => {
      const response = await userApi.getUserById(userId!);
      return transformUserResponse(response.data);
    },
    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Хук для получения краткой информации о пользователе с ролью в чате
 * Используется в UserPopover для отображения информации о собеседнике
 */
export function useUserBrief(userId: string | undefined, chatId: string | undefined) {
  return useQuery({
    queryKey: userKeys.brief(userId!, chatId!),
    queryFn: async () => {
      const response = await userApi.getUserBrief(userId!, chatId!);
      return response.data;
    },
    enabled: !!userId && !!chatId,
    staleTime: 5 * 60 * 1000, // 5 минут (кешируется на бекенде тоже)
  });
}

/**
 * Хук для поиска пользователей
 * Используется в форме создания чата для выбора участников
 */
export function useSearchUsers(query: string, enabled = true) {
  return useQuery({
    queryKey: [...userKeys.all, 'search', query],
    queryFn: async () => {
      const response = await userApi.searchUsers(query);
      return response.data.users;
    },
    enabled: enabled && query.length >= 2,
    staleTime: 30 * 1000, // 30 секунд
  });
}

/**
 * Хук для обновления профиля текущего пользователя
 */
export function useUpdateProfile() {
  const queryClient = useQueryClient();
  const { setUser } = useAuthStore();

  return useMutation({
    mutationFn: async ({
      data,
      avatar,
    }: {
      data: UpdateUserRequest;
      avatar?: File;
    }) => {
      const formData = createUpdateUserFormData(data, avatar);
      const response = await userApi.updateMe(formData);
      return transformUserResponse(response.data);
    },
    onSuccess: (updatedUser) => {
      // Обновляем данные в Auth Store СНАЧАЛА (синхронно)
      setUser(updatedUser);
      // Обновляем кеш React Query
      queryClient.setQueryData(userKeys.me(), updatedUser);
      // Инвалидируем кеш для гарантии свежести данных при следующем запросе
      queryClient.invalidateQueries({ queryKey: userKeys.me() });
      toast.success('Профиль обновлён');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка обновления профиля';
      toast.error(message);
    },
  });
}

/**
 * Хук для получения всех ролей (admin)
 */
export function useRoles() {
  return useQuery({
    queryKey: userKeys.roles(),
    queryFn: async () => {
      const response = await userApi.getRoles();
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 минут
  });
}

/**
 * Хук для получения всех permissions (admin)
 */
export function usePermissions() {
  return useQuery({
    queryKey: userKeys.permissions(),
    queryFn: async () => {
      const response = await userApi.getPermissions();
      return response.data;
    },
    staleTime: 10 * 60 * 1000,
  });
}

/**
 * Хук для создания роли (admin)
 */
export function useCreateRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: { name: string; description?: string; permissionIds?: number[] }) => {
      const response = await userApi.createRole(data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.roles() });
      toast.success('Роль создана');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка создания роли';
      toast.error(message);
    },
  });
}

/**
 * Хук для обновления permissions роли (admin)
 */
export function useUpdateRolePermissions() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ roleId, permissionIds }: { roleId: number; permissionIds: number[] }) => {
      await userApi.updateRolePermissions(roleId, permissionIds);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.roles() });
      toast.success('Разрешения роли обновлены');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка обновления разрешений роли';
      toast.error(message);
    },
  });
}

/**
 * Хук для удаления роли (admin)
 */
export function useDeleteRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (roleId: number) => {
      await userApi.deleteRole(roleId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.roles() });
      toast.success('Роль удалена');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка удаления роли';
      toast.error(message);
    },
  });
}

/**
 * Хук для изменения роли пользователя (admin)
 */
export function useUpdateUserRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ userId, roleId }: { userId: string; roleId: number }) => {
      await userApi.updateUserRole(userId, roleId);
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: userKeys.detail(variables.userId) });
      queryClient.invalidateQueries({ queryKey: userKeys.roles() });
      toast.success('Роль пользователя изменена');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка изменения роли пользователя';
      toast.error(message);
    },
  });
}

