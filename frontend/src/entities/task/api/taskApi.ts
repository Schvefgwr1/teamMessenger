import { apiClient } from '@/shared/api';
import type { Task, TaskStatus } from '@/shared/types';

// ============================================================
// Response Types
// ============================================================

/**
 * Ответ API при получении задачи по ID
 * GET /api/v1/tasks/:task_id
 * API возвращает TaskServiceResponse с полем task
 */
export interface TaskServiceResponse {
  task: TaskResponse;
  files?: Array<{
    id: number;
    name: string;
    url: string;
    createdAt: string;
    fileType: {
      id: number;
      name: string;
    };
  }>;
}

/**
 * Ответ задачи от API (внутри TaskServiceResponse)
 */
export interface TaskResponse {
  id: number;
  title: string;
  description?: string;
  creatorID: string; // camelCase в API
  executorID?: string; // camelCase в API
  chatID?: string; // camelCase в API
  status: TaskStatus;
  createdAt: string;
}

/**
 * Ответ API при получении списка задач пользователя
 * GET /api/v1/users/:user_id/tasks
 * API возвращает TaskToList[] где status - это строка (название статуса)
 */
export interface TaskToListResponse {
  id: number;
  title: string;
  status: string; // Название статуса (например, "created")
  createdAt?: string; // Дата создания (может приходить с бекенда, но не указана в swagger)
}

/**
 * Ответ при получении списка задач пользователя
 */
export type GetUserTasksResponse = TaskToListResponse[];

/**
 * Ответ при получении всех статусов
 * GET /api/v1/tasks/statuses
 */
export type GetTaskStatusesResponse = TaskStatus[];

/**
 * Ответ при создании задачи
 * POST /api/v1/tasks
 */
export type CreateTaskResponse = Task;

// ============================================================
// Request Types
// ============================================================

/**
 * Параметры запроса для получения задач пользователя
 */
export interface GetUserTasksParams {
  limit?: number;
  offset?: number;
}

/**
 * Данные для создания задачи
 * Используется для формирования FormData
 */
export interface CreateTaskRequest {
  title: string;
  description?: string;
  executorId: string; // UUID
  chatId?: string; // UUID
  files?: File[];
}

// ============================================================
// Task API
// ============================================================

export const taskApi = {
  /**
   * Создать задачу
   * POST /api/v1/tasks (multipart/form-data)
   */
  createTask: (formData: FormData) =>
    apiClient.post<CreateTaskResponse>('/tasks', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  /**
   * Получить задачу по ID
   * GET /api/v1/tasks/:task_id
   * Возвращает TaskServiceResponse с полем task
   */
  getTaskById: (taskId: number) =>
    apiClient.get<TaskServiceResponse>(`/tasks/${taskId}`),

  /**
   * Обновить статус задачи
   * PATCH /api/v1/tasks/:task_id/status/:status_id
   */
  updateTaskStatus: (taskId: number, statusId: number) =>
    apiClient.patch(`/tasks/${taskId}/status/${statusId}`),

  /**
   * Получить задачи пользователя
   * GET /api/v1/users/:user_id/tasks
   */
  getUserTasks: (userId: string, params?: GetUserTasksParams) =>
    apiClient.get<GetUserTasksResponse>(`/users/${userId}/tasks`, { params }),

  /**
   * Получить все статусы задач
   * GET /api/v1/tasks/statuses
   */
  getAllStatuses: () =>
    apiClient.get<GetTaskStatusesResponse>('/tasks/statuses'),
};

// ============================================================
// Helper Functions
// ============================================================

/**
 * Создать FormData для создания задачи
 * Согласно swagger - отдельные поля formData
 */
export function createTaskFormData(data: CreateTaskRequest): FormData {
  const formData = new FormData();

  // Обязательные поля
  formData.append('title', data.title);

  // Опциональные поля
  if (data.description) {
    formData.append('description', data.description);
  }

  if (data.executorId) {
    formData.append('executor_id', data.executorId);
  }

  if (data.chatId) {
    formData.append('chat_id', data.chatId);
  }

  // Файлы (множественные)
  if (data.files && data.files.length > 0) {
    data.files.forEach((file) => {
      formData.append('files', file);
    });
  }

  return formData;
}

