import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import {
  chatApi,
  chatRolesApi,
  createChatFormData,
  updateChatFormData,
  createMessageFormData,
  type GetChatResponse,
  type CreateChatRequest,
  type UpdateChatRequest,
  type ChangeUserRoleRequest,
  type MyRoleResponse,
} from './chatApi';
import { useAuthStore } from '@/entities/session';
import { toast } from '@/shared/ui';
import type { Chat, Message, File as ApiFile } from '@/shared/types';

// ============================================================
// Query Keys
// ============================================================

export const chatKeys = {
  all: ['chats'] as const,
  lists: () => [...chatKeys.all, 'list'] as const,
  list: (userId: string) => [...chatKeys.lists(), userId] as const,
  details: () => [...chatKeys.all, 'detail'] as const,
  detail: (chatId: string) => [...chatKeys.details(), chatId] as const,
  messages: (chatId: string) => [...chatKeys.detail(chatId), 'messages'] as const,
  messagesInfinite: (chatId: string) => [...chatKeys.messages(chatId), 'infinite'] as const,
  search: (chatId: string, query: string) => [...chatKeys.detail(chatId), 'search', query] as const,
  myRole: (chatId: string) => [...chatKeys.detail(chatId), 'myRole'] as const,
};

export const chatRolesKeys = {
  all: ['chatRoles'] as const,
  list: () => [...chatRolesKeys.all, 'list'] as const,
};

// ============================================================
// Types
// ============================================================

/**
 * Чат с файлом аватара (для UI)
 */
export interface ChatWithAvatar extends Chat {
  avatar?: ApiFile;
}

/**
 * Преобразовать GetChatResponse в ChatWithAvatar
 */
export function transformChatResponse(response: GetChatResponse): ChatWithAvatar {
  return {
    ...response.chat,
    avatar: response.file ?? undefined,
  };
}

// ============================================================
// Queries
// ============================================================

/**
 * Получить чаты текущего пользователя
 */
export function useUserChats() {
  const { user, isAuthenticated } = useAuthStore();

  return useQuery({
    queryKey: chatKeys.list(user?.ID ?? ''),
    queryFn: async () => {
      if (!user?.ID) throw new Error('User not authenticated');
      const response = await chatApi.getUserChats(user.ID);
      // Ответ приходит как массив ChatResponse[]
      return response.data;
    },
    enabled: isAuthenticated && !!user?.ID,
    staleTime: 30 * 1000, // 30 секунд
  });
}

/**
 * Получить сообщения чата
 */
export function useChatMessages(chatId: string | undefined, limit = 50) {
  return useQuery({
    queryKey: chatKeys.messages(chatId!),
    queryFn: async () => {
      const response = await chatApi.getMessages(chatId!, { limit });
      // API возвращает массив напрямую
      return response.data || [];
    },
    enabled: !!chatId,
    staleTime: 10 * 1000, // 10 секунд
    refetchInterval: 5 * 1000, // Polling каждые 5 секунд для real-time эффекта
  });
}

/**
 * Получить сообщения чата с infinite scroll
 */
export function useChatMessagesInfinite(chatId: string | undefined, pageSize = 30) {
  return useInfiniteQuery({
    queryKey: chatKeys.messagesInfinite(chatId!),
    queryFn: async ({ pageParam = 0 }) => {
      const response = await chatApi.getMessages(chatId!, {
        offset: pageParam,
        limit: pageSize,
      });
      return response.data || [];
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage.length < pageSize) return undefined;
      return allPages.flat().length;
    },
    enabled: !!chatId,
    staleTime: 10 * 1000,
  });
}

/**
 * Поиск сообщений в чате
 */
export function useSearchMessages(chatId: string | undefined, query: string) {
  return useQuery({
    queryKey: chatKeys.search(chatId!, query),
    queryFn: async () => {
      const response = await chatApi.searchMessages(chatId!, { query });
      return response.data?.messages || [];
    },
    enabled: !!chatId && query.length >= 2,
    staleTime: 60 * 1000, // 1 минута
  });
}

/**
 * Получить свою роль в чате с permissions
 */
export function useMyRoleInChat(chatId: string | undefined) {
  return useQuery({
    queryKey: chatKeys.myRole(chatId!),
    queryFn: async () => {
      const response = await chatApi.getMyRole(chatId!);
      return response.data;
    },
    enabled: !!chatId,
    staleTime: 5 * 60 * 1000, // 5 минут
  });
}

/**
 * Получить все роли чатов (для выбора при изменении роли)
 */
export function useChatRoles() {
  return useQuery({
    queryKey: chatRolesKeys.list(),
    queryFn: async () => {
      const response = await chatRolesApi.getAllRoles();
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 минут
  });
}

/**
 * Получить все permissions чатов
 */
export function useChatPermissions() {
  return useQuery({
    queryKey: [...chatRolesKeys.all, 'permissions'] as const,
    queryFn: async () => {
      const response = await chatRolesApi.getAllPermissions();
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 минут
  });
}

/**
 * Получить список участников чата
 */
export function useChatMembers(chatId: string | undefined) {
  return useQuery({
    queryKey: [...chatKeys.detail(chatId!), 'members'] as const,
    queryFn: async () => {
      const response = await chatApi.getChatMembers(chatId!);
      return response.data;
    },
    enabled: !!chatId,
    staleTime: 2 * 60 * 1000, // 2 минуты
  });
}

/**
 * Хелпер для проверки наличия permission
 */
export function hasPermission(
  myRole: MyRoleResponse | undefined,
  permissionName: string
): boolean {
  if (!myRole?.permissions) return false;
  return myRole.permissions.some(p => p.name === permissionName);
}

/**
 * Хук для обновления permissions роли чата
 */
export function useUpdateChatRolePermissions() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ roleId, permissionIds }: { roleId: number; permissionIds: number[] }) => {
      const response = await chatRolesApi.updateRolePermissions(roleId, { permissionIds });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatRolesKeys.list() });
      toast.success('Разрешения роли чата обновлены');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка обновления разрешений роли чата';
      toast.error(message);
    },
  });
}

// ============================================================
// Mutations
// ============================================================

/**
 * Создать чат
 */
export function useCreateChat() {
  const queryClient = useQueryClient();
  const { user } = useAuthStore();

  return useMutation({
    mutationFn: async ({
      data,
      avatar,
    }: {
      data: Omit<CreateChatRequest, 'ownerID'>;
      avatar?: File;
    }) => {
      const formData = createChatFormData(
        { ...data, ownerID: user!.ID },
        avatar
      );
      const response = await chatApi.createChat(formData);
      return response.data; // CreateChatResponse
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatKeys.lists() });
      toast.success('Чат создан');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка создания чата';
      toast.error(message);
    },
  });
}

/**
 * Обновить чат
 */
export function useUpdateChat(chatId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      data,
      avatar,
    }: {
      data: UpdateChatRequest;
      avatar?: File;
    }) => {
      const formData = updateChatFormData(data, avatar);
      const response = await chatApi.updateChat(chatId, formData);
      return transformChatResponse(response.data);
    },
    onSuccess: (updatedChat) => {
      queryClient.setQueryData(chatKeys.detail(chatId), updatedChat);
      queryClient.invalidateQueries({ queryKey: chatKeys.lists() });
      toast.success('Чат обновлён');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка обновления чата';
      toast.error(message);
    },
  });
}

/**
 * Удалить чат
 */
export function useDeleteChat() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (chatId: string) => chatApi.deleteChat(chatId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatKeys.lists() });
      toast.success('Чат удалён');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка удаления чата';
      toast.error(message);
    },
  });
}

/**
 * Отправить сообщение
 */
export function useSendMessage(chatId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      content,
      files,
    }: {
      content: string;
      files?: File[];
    }) => {
      const formData = createMessageFormData(content, files);
      const response = await chatApi.sendMessage(chatId, formData);
      return response.data;
    },
    onSuccess: (newMessage) => {
      if (!newMessage) return;
      // Optimistic update: добавляем сообщение в кеш
      queryClient.setQueryData<Message[]>(
        chatKeys.messages(chatId),
        (old = []) => [...old, newMessage]
      );
      // Инвалидируем для получения актуальных данных
      queryClient.invalidateQueries({ queryKey: chatKeys.messages(chatId) });
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка отправки сообщения';
      toast.error(message);
    },
  });
}

/**
 * Забанить пользователя в чате
 */
export function useBanUser(chatId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: string) => chatApi.banUser(chatId, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatKeys.detail(chatId) });
      toast.success('Пользователь забанен');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка бана пользователя';
      toast.error(message);
    },
  });
}

/**
 * Изменить роль пользователя в чате
 */
export function useChangeUserRole(chatId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ChangeUserRoleRequest) =>
      chatApi.changeUserRole(chatId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatKeys.detail(chatId) });
      toast.success('Роль изменена');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка изменения роли';
      toast.error(message);
    },
  });
}

