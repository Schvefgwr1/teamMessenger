import { apiClient } from '@/shared/api';
import type { Chat, Message, File as ApiFile } from '@/shared/types';

// ============================================================
// Response Types
// ============================================================

/**
 * Ответ при получении чата с аватаром
 */
export interface GetChatResponse {
  chat: Chat;
  file: ApiFile | null;
}

/**
 * Ответ при получении списка чатов
 * Приходит массив ChatResponse DTO
 */
export type GetChatsResponse = Chat[];

/**
 * Ответ при получении сообщений
 * API возвращает массив сообщений напрямую
 */
export type GetMessagesResponse = Message[];

/**
 * Ответ при отправке сообщения
 * API может вернуть сообщение напрямую или в объекте
 */
export type SendMessageResponse = Message;

/**
 * Ответ при создании чата
 */
export interface CreateChatResponse {
  id: string;
  name: string;
  description?: string;
  ownerID: string;
  userIDs: string[];
  avatarFileID?: number;
}

// ============================================================
// Request Types
// ============================================================

export interface CreateChatRequest {
  name: string;
  description?: string;
  ownerID: string;
  userIDs: string[];
}

export interface UpdateChatRequest {
  name?: string;
  description?: string;
  addUserIDs?: string[];
  removeUserIDs?: string[];
}

export interface GetMessagesParams {
  offset?: number;
  limit?: number;
}

export interface SearchMessagesParams {
  query: string;
  offset?: number;
  limit?: number;
}

export interface ChangeUserRoleRequest {
  user_id: string;
  role_id: number;
}

/**
 * Ответ поиска сообщений
 */
export interface SearchMessagesResponse {
  messages: Message[] | null;
  total: number | null;
}

/**
 * Ответ с ролью текущего пользователя в чате
 */
export interface MyRoleResponse {
  roleId: number;
  roleName: string;
  permissions: Array<{ id: number; name: string }>;
}

/**
 * Участник чата
 */
export interface ChatMemberResponse {
  userId: string;
  roleId: number;
  roleName: string;
}

// ============================================================
// Chat API
// ============================================================

export const chatApi = {
  /**
   * Получить чаты пользователя
   * GET /api/v1/chats/:userId
   */
  getUserChats: (userId: string) =>
    apiClient.get<GetChatsResponse>(`/chats/${userId}`),

  /**
   * Создать чат
   * POST /api/v1/chats (multipart/form-data)
   */
  createChat: (formData: FormData) =>
    apiClient.post<CreateChatResponse>('/chats', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  /**
   * Обновить чат
   * PUT /api/v1/chats/:chatId (multipart/form-data)
   */
  updateChat: (chatId: string, formData: FormData) =>
    apiClient.put<GetChatResponse>(`/chats/${chatId}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  /**
   * Удалить чат
   * DELETE /api/v1/chats/:chatId
   */
  deleteChat: (chatId: string) =>
    apiClient.delete(`/chats/${chatId}`),

  /**
   * Забанить пользователя в чате
   * PATCH /api/v1/chats/:chatId/ban/:userId
   */
  banUser: (chatId: string, userId: string) =>
    apiClient.patch(`/chats/${chatId}/ban/${userId}`),

  /**
   * Изменить роль пользователя в чате
   * PATCH /api/v1/chats/:chatId/roles/change
   */
  changeUserRole: (chatId: string, data: ChangeUserRoleRequest) =>
    apiClient.patch(`/chats/${chatId}/roles/change`, data),

  // ============================================================
  // Messages
  // ============================================================

  /**
   * Получить сообщения чата
   * GET /api/v1/chats/messages/:chatId
   */
  getMessages: (chatId: string, params?: GetMessagesParams) =>
    apiClient.get<GetMessagesResponse>(`/chats/messages/${chatId}`, { params }),

  /**
   * Отправить сообщение
   * POST /api/v1/chats/messages/:chatId (multipart/form-data)
   */
  sendMessage: (chatId: string, formData: FormData) =>
    apiClient.post<SendMessageResponse>(`/chats/messages/${chatId}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  /**
   * Поиск сообщений в чате
   * GET /api/v1/chats/search/:chatId
   */
  searchMessages: (chatId: string, params: SearchMessagesParams) =>
    apiClient.get<SearchMessagesResponse>(`/chats/search/${chatId}`, { params }),

  /**
   * Получить свою роль в чате с permissions
   * GET /api/v1/chats/me/role/:chatId
   */
  getMyRole: (chatId: string) =>
    apiClient.get<MyRoleResponse>(`/chats/me/role/${chatId}`),

  /**
   * Получить список участников чата
   * GET /api/v1/chats/members/:chatId
   */
  getChatMembers: (chatId: string) =>
    apiClient.get<ChatMemberResponse[]>(`/chats/members/${chatId}`),
};

// ============================================================
// Chat Roles API (для админки и настроек чата)
// ============================================================

export interface ChatRoleResponse {
  id: number;
  name: string;
  permissions: Array<{ id: number; name: string }>;
}

export const chatRolesApi = {
  /**
   * Получить все роли чатов
   * GET /api/v1/chat-roles
   */
  getAllRoles: () =>
    apiClient.get<ChatRoleResponse[]>('/chat-roles'),

  /**
   * Получить роль по ID
   * GET /api/v1/chat-roles/:roleId
   */
  getRoleById: (roleId: number) =>
    apiClient.get<ChatRoleResponse>(`/chat-roles/${roleId}`),
};

// ============================================================
// Helper Functions
// ============================================================

/**
 * Создать FormData для создания чата
 * Согласно swagger - отдельные поля formData, НЕ JSON
 */
export function createChatFormData(
  data: CreateChatRequest,
  avatar?: File
): FormData {
  const formData = new FormData();
  
  // Обязательные поля
  formData.append('name', data.name);
  formData.append('ownerID', data.ownerID);
  
  // userIDs как CSV (collectionFormat: csv)
  if (data.userIDs && data.userIDs.length > 0) {
    formData.append('userIDs', data.userIDs.join(','));
  } else {
    formData.append('userIDs', ''); // required field
  }
  
  // Опциональные поля
  if (data.description) {
    formData.append('description', data.description);
  }
  
  if (avatar) {
    formData.append('avatar', avatar);
  }
  
  return formData;
}

/**
 * Создать FormData для обновления чата
 * Согласно swagger - отдельные поля formData
 */
export function updateChatFormData(
  data: UpdateChatRequest,
  avatar?: File
): FormData {
  const formData = new FormData();
  
  if (data.name) {
    formData.append('name', data.name);
  }
  
  if (data.description !== undefined) {
    formData.append('description', data.description);
  }
  
  if (data.addUserIDs && data.addUserIDs.length > 0) {
    formData.append('addUserIDs', data.addUserIDs.join(','));
  }
  
  if (data.removeUserIDs && data.removeUserIDs.length > 0) {
    formData.append('removeUserIDs', data.removeUserIDs.join(','));
  }
  
  if (avatar) {
    formData.append('avatar', avatar);
  }
  
  return formData;
}

/**
 * Создать FormData для отправки сообщения
 * Согласно swagger - content как отдельное поле
 */
export function createMessageFormData(
  content: string,
  files?: File[]
): FormData {
  const formData = new FormData();
  
  formData.append('content', content);
  
  if (files && files.length > 0) {
    files.forEach((file) => {
      formData.append('files', file);
    });
  }
  
  return formData;
}

