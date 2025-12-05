import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  taskApi,
  createTaskFormData,
  type CreateTaskRequest,
  type GetUserTasksParams,
  type TaskToListResponse,
  type TaskResponse as ApiTaskResponse,
  type TaskServiceResponse,
} from './taskApi';
import { useAuthStore } from '@/entities/session';
import { toast } from '@/shared/ui';
import type { Task, TaskStatus } from '@/shared/types';

// ============================================================
// Query Keys
// ============================================================

export const taskKeys = {
  all: ['tasks'] as const,
  lists: () => [...taskKeys.all, 'list'] as const,
  list: (userId: string) => [...taskKeys.lists(), userId] as const,
  details: () => [...taskKeys.all, 'detail'] as const,
  detail: (taskId: number) => [...taskKeys.details(), taskId] as const,
  statuses: () => [...taskKeys.all, 'statuses'] as const,
};

// ============================================================
// Queries
// ============================================================

/**
 * Преобразовать TaskToListResponse в Task
 * Нужно сопоставить строку статуса с объектом статуса
 */
function transformTaskToList(
  taskToList: TaskToListResponse,
  statuses: TaskStatus[]
): Task {
  // Находим статус по названию (case-insensitive)
  const statusName = taskToList.status.toLowerCase();
  const status = statuses.find(
    (s) => s.name.toLowerCase() === statusName
  ) || {
    id: 0,
    name: taskToList.status,
  };

  return {
    id: taskToList.id,
    title: taskToList.title,
    status,
    creatorId: '', // Не приходит в списке
    createdAt: '', // Не приходит в списке
  };
}

/**
 * Получить задачи текущего пользователя
 */
export function useUserTasks(params?: GetUserTasksParams) {
  const { user, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  
  // Загружаем статусы параллельно
  const { data: statuses = [] } = useTaskStatuses();

  return useQuery({
    queryKey: [...taskKeys.list(user?.ID ?? ''), params, statuses],
    queryFn: async () => {
      if (!user?.ID) throw new Error('User not authenticated');
      
      const response = await taskApi.getUserTasks(user.ID, params);
      const tasksList = Array.isArray(response.data) ? response.data : [];
      
      console.log('useUserTasks - raw data:', {
        tasksList,
        statuses,
        statusesCount: statuses.length,
      });
      
      // Преобразуем TaskToListResponse[] в Task[]
      const transformedTasks = tasksList.map((taskToList) =>
        transformTaskToList(taskToList, statuses)
      );
      
      console.log('useUserTasks - transformed tasks:', transformedTasks);
      
      return transformedTasks;
    },
    enabled: isAuthenticated && !!user?.ID && statuses.length > 0,
    staleTime: 30 * 1000, // 30 секунд
    refetchInterval: 60 * 1000, // Обновление раз в минуту в фоне
  });
}

/**
 * Преобразовать TaskResponse из API в Task
 */
function transformTaskResponse(apiTask: ApiTaskResponse): Task {
  return {
    id: apiTask.id,
    title: apiTask.title,
    description: apiTask.description,
    status: apiTask.status,
    creatorId: apiTask.creatorID,
    executorId: apiTask.executorID, // Исправлено: теперь executorID тоже camelCase
    chatId: apiTask.chatID, // Исправлено: теперь chatID camelCase
    createdAt: apiTask.createdAt,
    files: undefined, // Файлы обрабатываются отдельно если нужно
  };
}

/**
 * Получить задачу по ID
 */
export function useTask(taskId: number | undefined) {
  return useQuery({
    queryKey: taskKeys.detail(taskId!),
    queryFn: async () => {
      const response = await taskApi.getTaskById(taskId!);
      const serviceResponse = response.data as TaskServiceResponse;
      
      console.log('useTask - raw response:', serviceResponse);
      
      // Извлекаем задачу из ответа
      if (!serviceResponse?.task) {
        throw new Error('Task not found in response');
      }
      
      // Преобразуем TaskResponse в Task
      const task = transformTaskResponse(serviceResponse.task);
      
      // Добавляем файлы если есть
      if (serviceResponse.files && serviceResponse.files.length > 0) {
        task.files = serviceResponse.files.map((file) => ({
          taskId: task.id,
          fileId: file.id,
          file: {
            id: file.id,
            name: file.name,
            url: file.url,
            createdAt: file.createdAt,
            fileType: file.fileType,
          },
        }));
      }
      
      console.log('useTask - transformed task:', task);
      
      return task;
    },
    enabled: !!taskId,
    staleTime: 30 * 1000,
  });
}

/**
 * Получить все статусы задач
 */
export function useTaskStatuses() {
  return useQuery({
    queryKey: taskKeys.statuses(),
    queryFn: async () => {
      const response = await taskApi.getAllStatuses();
      // Убеждаемся что возвращаем массив, а не null
      return Array.isArray(response.data) ? response.data : [];
    },
    staleTime: 10 * 60 * 1000, // 10 минут - статусы редко меняются
  });
}

// ============================================================
// Mutations
// ============================================================

/**
 * Создать задачу
 */
export function useCreateTask() {
  const queryClient = useQueryClient();
  const { user } = useAuthStore();

  return useMutation({
    mutationFn: async (data: Omit<CreateTaskRequest, 'executorId' | 'chatId'> & {
      executorId: string;
      chatId?: string;
    }) => {
      const formData = createTaskFormData(data);
      const response = await taskApi.createTask(formData);
      return response.data;
    },
    onSuccess: () => {
      if (user?.ID) {
        queryClient.invalidateQueries({ queryKey: taskKeys.list(user.ID) });
      }
      toast.success('Задача создана');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка создания задачи';
      toast.error(message);
    },
  });
}

/**
 * Обновить статус задачи
 * С оптимистичным обновлением для drag-and-drop
 * taskId передаётся в mutate, а не в хук
 */
export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();
  const { user } = useAuthStore();

  return useMutation({
    mutationFn: async ({ taskId, statusId }: { taskId: number; statusId: number }) => {
      console.log('useUpdateTaskStatus mutationFn:', { taskId, statusId });
      await taskApi.updateTaskStatus(taskId, statusId);
      return { taskId, statusId };
    },
    onMutate: async ({ taskId, statusId: newStatusId }) => {
      // Отменяем исходящие запросы
      await queryClient.cancelQueries({ queryKey: taskKeys.detail(taskId) });
      await queryClient.cancelQueries({ queryKey: taskKeys.list(user?.ID ?? '') });

      // Сохраняем предыдущее значение
      const previousTask = queryClient.getQueryData<Task>(taskKeys.detail(taskId));
      const previousTasks = queryClient.getQueryData<Task[]>(
        taskKeys.list(user?.ID ?? '')
      );

      // Оптимистично обновляем задачу
      if (previousTask) {
        queryClient.setQueryData<Task>(taskKeys.detail(taskId), (old) => {
          if (!old) return old;
          return {
            ...old,
            status: { id: newStatusId, name: old.status.name }, // Временно сохраняем старое имя
          };
        });
      }

      // Оптимистично обновляем список задач
      if (previousTasks && user?.ID) {
        queryClient.setQueryData<Task[]>(
          taskKeys.list(user.ID),
          (old) => {
            if (!old) return old;
            return old.map((task) =>
              task.id === taskId
                ? { ...task, status: { id: newStatusId, name: task.status.name } }
                : task
            );
          }
        );
      }

      return { previousTask, previousTasks, taskId };
    },
    onError: (err, { taskId }, context) => {
      // Откатываем при ошибке
      if (context?.previousTask) {
        queryClient.setQueryData(taskKeys.detail(taskId), context.previousTask);
      }
      if (context?.previousTasks && user?.ID) {
        queryClient.setQueryData(taskKeys.list(user.ID), context.previousTasks);
      }
      toast.error('Ошибка обновления статуса задачи');
    },
    onSettled: (data, error, { taskId }) => {
      // Перезапрашиваем в любом случае для получения актуальных данных
      queryClient.invalidateQueries({ queryKey: taskKeys.detail(taskId) });
      if (user?.ID) {
        queryClient.invalidateQueries({ queryKey: taskKeys.list(user.ID) });
      }
    },
  });
}

