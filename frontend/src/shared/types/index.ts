// ============================================================
// FILE SERVICE MODELS
// ============================================================

/**
 * Тип файла (image, document, video и т.д.)
 * @see fileService/internal/models/file_type.go
 */
export interface FileType {
  id: number;
  name: string;
}

/**
 * Файл (аватары, прикреплённые файлы)
 * @see fileService/internal/models/file.go
 */
export interface File {
  id: number;
  name: string;
  url: string;
  createdAt: string;
  fileType: FileType;
}

// ============================================================
// USER SERVICE MODELS
// ============================================================

/**
 * Разрешение пользователя
 * @see userService/internal/models/permission.go
 */
export interface Permission {
  ID: number;
  Name: string;
  Description?: string;
}

/**
 * Роль пользователя
 * @see userService/internal/models/role.go
 */
export interface Role {
  ID: number;
  Name: string;
  Description?: string;
  Permissions: Permission[];
}

/**
 * Пользователь
 * @see userService/internal/models/user.go
 */
export interface User {
  ID: string;
  Username: string;
  Email: string;
  Description?: string;
  Gender?: string;
  Age?: number;
  Role: Role
}

/**
 * Пользователь с файлом аватара
 */
export interface UserWithAvatar extends User {
  avatar?: File;
}

/**
 * Краткая информация о пользователе (для popover в чате)
 * @see apiService/internal/dto/user_brief_dto.go
 */
export interface UserBrief {
  username: string;
  email: string;
  age?: number;
  description?: string;
  avatarFile?: File;
  chatRoleName?: string;
}

/**
 * Результат поиска пользователя
 * @see apiService/internal/dto/user_search_dto.go
 */
export interface UserSearchResult {
  id: string;
  username: string;
  email: string;
  avatarFile?: File;
}

/**
 * Ответ на поиск пользователей
 */
export interface UserSearchResponse {
  users: UserSearchResult[];
}

/**
 * Ответ API при получении пользователя
 */
export interface GetUserResponse {
  file: File | null;
  user: User;
}

/**
 * Запрос на обновление пользователя
 */
export interface UpdateUserRequest {
  username?: string;
  description?: string;
  gender?: string;
  age?: number;
  roleID?: number;
}

// ============================================================
// CHAT SERVICE MODELS
// ============================================================

/**
 * Разрешение в чате
 * @see chatService/internal/models/chat_permission.go
 */
export interface ChatPermission {
  id: number;
  name: string;
}

/**
 * Роль в чате
 * @see chatService/internal/models/chat_role.go
 */
export interface ChatRole {
  id: number;
  name: string;
  permissions: ChatPermission[];
}

/**
 * Ответ с ролью текущего пользователя в чате
 * @see common/contracts/api-chat/chat.go MyRoleResponse
 */
export interface MyRoleInChat {
  roleId: number;
  roleName: string;
  permissions: ChatPermission[];
}

/**
 * Известные chat permissions для проверки на фронтенде
 */
export const CHAT_PERMISSIONS = {
  EDIT_CHAT: 'edit_chat',
  DELETE_CHAT: 'delete_chat',
  BAN_USER: 'ban_user',
  CHANGE_ROLE: 'change_role',
  SEND_MESSAGE: 'send_message',
} as const;

export type ChatPermissionName = typeof CHAT_PERMISSIONS[keyof typeof CHAT_PERMISSIONS];

/**
 * Участник чата
 * @see chatService/internal/models/chat_user.go
 */
export interface ChatUser {
  chatId: string;
  userId: string;
  role: ChatRole;
  isBanned?: boolean;
}

/**
 * Чат
 * @see chatService/internal/models/chat.go
 */
export interface Chat {
  id: string;
  name: string;
  isGroup: boolean;
  description?: string;
  avatarFileID?: number;
  avatarFile?: File;
  createdAt: string;
  users?: ChatUser[];
}

/**
 * Сообщение
 * @see chatService/internal/handlers/dto/get_chat_messages_dto.go
 */
export interface Message {
  id: string;
  chatID: string;
  senderID?: string | null;
  content: string;
  updatedAt?: string | null;
  createdAt: string;
  files?: File[] | null;
}

// ============================================================
// TASK SERVICE MODELS
// ============================================================

/**
 * Статус задачи
 * @see taskService/internal/models/task_status.go
 */
export interface TaskStatus {
  id: number;
  name: string;
}

/**
 * Файл прикреплённый к задаче
 */
export interface TaskFile {
  taskId: number;
  fileId: number;
  file?: File;
}

/**
 * Задача
 * @see taskService/internal/models/task.go
 */
export interface Task {
  id: number;
  title: string;
  description?: string;
  status: TaskStatus;
  creatorId: string;
  executorId?: string;
  chatId?: string;
  createdAt: string;
  files?: TaskFile[];
}

// ============================================================
// UTILITY TYPES
// ============================================================

/**
 * Ответ со списком с пагинацией
 */
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  offset: number;
  limit: number;
}

