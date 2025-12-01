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

