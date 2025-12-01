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
  avatarFileId?: number;
  createdAt: string;
  users?: ChatUser[];
}

/**
 * Файл прикреплённый к сообщению
 */
export interface MessageFile {
  messageId: string;
  fileId: number;
  file?: File;
}

/**
 * Сообщение
 * @see chatService/internal/models/message.go
 */
export interface Message {
  id: string;
  chatId: string;
  senderId?: string;
  content: string;
  updatedAt?: string;
  createdAt: string;
  files?: MessageFile[];
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

