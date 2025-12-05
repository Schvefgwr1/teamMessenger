import { useState, useMemo } from 'react';
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
  closestCorners,
} from '@dnd-kit/core';
import { Plus } from 'lucide-react';
import { Button } from '@/shared/ui';
import { useUserTasks, useTaskStatuses, useUpdateTaskStatus } from '@/entities/task';
import { TaskColumn, TaskColumnSkeleton } from './TaskColumn';
import { TaskCard } from './TaskCard';
import type { Task } from '@/shared/types';

interface TaskBoardProps {
  onTaskClick?: (taskId: number) => void;
  onCreateTask?: () => void;
}

/**
 * Kanban доска задач
 * Группирует задачи по статусам и поддерживает drag-and-drop
 */
export function TaskBoard({ onTaskClick, onCreateTask }: TaskBoardProps) {
  const [activeTaskId, setActiveTaskId] = useState<number | null>(null);

  // Загрузка данных
  const { 
    data: tasks = [], 
    isLoading: tasksLoading,
    error: tasksError,
  } = useUserTasks({
    limit: 100, // Загружаем достаточно задач для Kanban
    offset: 0,
  });
  const { 
    data: statuses = [], 
    isLoading: statusesLoading,
    error: statusesError,
  } = useTaskStatuses();

  // Группировка задач по статусам
  const tasksByStatus = useMemo(() => {
    const grouped: Record<number, Task[]> = {};
    
    // Проверяем что статусы загружены
    if (!statuses || statuses.length === 0) {
      return grouped;
    }
    
    // Инициализируем все статусы пустыми массивами
    statuses.forEach((status) => {
      grouped[status.id] = [];
    });

    // Распределяем задачи по статусам
    if (tasks && tasks.length > 0) {
      tasks.forEach((task) => {
        const statusId = task.status?.id;
        if (statusId && grouped[statusId] !== undefined) {
          grouped[statusId].push(task);
        }
      });
    }

    return grouped;
  }, [tasks, statuses]);

  // Настройка сенсоров для drag-and-drop
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8, // Минимальное расстояние для начала перетаскивания
      },
    })
  );

  // Обработка начала перетаскивания
  const handleDragStart = (event: DragStartEvent) => {
    const taskId = parseInt(event.active.id as string, 10);
    setActiveTaskId(taskId);
  };

  // Хук для обновления статуса
  const updateTaskStatus = useUpdateTaskStatus();

  // Обработка окончания перетаскивания
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTaskId(null);

    if (!over) return;

    const taskId = parseInt(active.id as string, 10);
    const overId = over.id as string;
    
    // Находим задачу, которую перетаскиваем
    const task = tasks?.find((t) => t.id === taskId);
    if (!task) return;

    // Определяем новый статус
    let newStatusId: number | null = null;
    
    // Проверяем, является ли over.id ID статуса (колонки)
    const statusId = parseInt(overId, 10);
    if (!isNaN(statusId) && statuses.some((s) => s.id === statusId)) {
      // Это ID статуса (колонки)
      newStatusId = statusId;
    } else {
      // Это может быть ID задачи - находим статус этой задачи
      const overTask = tasks?.find((t) => t.id === parseInt(overId, 10));
      if (overTask) {
        newStatusId = overTask.status.id;
      } else {
        // Не удалось определить статус
        console.error('Cannot determine new status for task', { 
          overId, 
          taskId, 
          availableStatusIds: statuses.map(s => s.id),
          availableTaskIds: tasks?.map(t => t.id)
        });
        return;
      }
    }

    // Логируем для отладки
    console.log('TaskBoard handleDragEnd:', {
      taskId,
      currentStatusId: task.status.id,
      currentStatusName: task.status.name,
      newStatusId,
      newStatusName: statuses.find(s => s.id === newStatusId)?.name,
      overId,
      allStatuses: statuses.map(s => ({ id: s.id, name: s.name }))
    });

    // Если статус не изменился, ничего не делаем
    if (task.status.id === newStatusId) return;

    // Обновляем статус через mutation
    updateTaskStatus.mutate({ taskId, statusId: newStatusId });
  };

  const isLoading = tasksLoading || statusesLoading;

  // Отладочная информация
  console.log('TaskBoard render:', {
    isLoading,
    tasksLoading,
    statusesLoading,
    tasksCount: tasks?.length ?? 0,
    statusesCount: statuses?.length ?? 0,
    tasks: tasks,
    statuses: statuses,
    tasksByStatus,
    tasksError,
    statusesError,
  });

  // Обработка ошибок
  if (tasksError || statusesError) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <p className="text-error mb-2">Ошибка загрузки данных</p>
          <p className="text-sm text-neutral-500 mb-4">
            {tasksError?.message || statusesError?.message || 'Неизвестная ошибка'}
          </p>
          <Button onClick={() => window.location.reload()}>
            Обновить страницу
          </Button>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="flex gap-4 overflow-x-auto pb-4">
        {[1, 2, 3, 4].map((i) => (
          <TaskColumnSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (!statuses || statuses.length === 0) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <p className="text-neutral-400 mb-4">Статусы задач не загружены</p>
          <p className="text-xs text-neutral-500">
            Проверьте консоль для отладки
          </p>
        </div>
      </div>
    );
  }

  if ((!tasks || tasks.length === 0) && !isLoading) {
    return (
      <div className="flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-neutral-100 mb-1">Задачи</h1>
            <p className="text-sm text-neutral-500">
              У вас пока нет задач
            </p>
          </div>
          {onCreateTask && (
            <Button onClick={onCreateTask} leftIcon={<Plus size={18} />}>
              Создать задачу
            </Button>
          )}
        </div>

        {/* Empty state */}
        <div className="flex items-center justify-center py-12 border border-neutral-800 rounded-lg bg-neutral-900/50">
          <div className="text-center">
            <p className="text-neutral-400 mb-4">Нет задач</p>
            <p className="text-sm text-neutral-500 mb-4">
              Создайте первую задачу, чтобы начать работу
            </p>
            {onCreateTask && (
              <Button onClick={onCreateTask} leftIcon={<Plus size={18} />}>
                Создать задачу
              </Button>
            )}
          </div>
        </div>
      </div>
    );
  }

  // Находим активную задачу для DragOverlay
  const activeTask = activeTaskId && tasks
    ? tasks.find((t) => t.id === activeTaskId)
    : null;

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCorners}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className="flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-neutral-100 mb-1">Задачи</h1>
            <p className="text-sm text-neutral-500">
              Всего задач: {tasks?.length ?? 0}
            </p>
          </div>
          {onCreateTask && (
            <Button onClick={onCreateTask} leftIcon={<Plus size={18} />}>
              Создать задачу
            </Button>
          )}
        </div>

        {/* Kanban доска */}
        <div className="overflow-x-auto pb-4">
          <div className="flex gap-4 min-h-[422px]">
            {statuses.map((status) => (
              <TaskColumn
                key={status.id}
                status={status}
                tasks={tasksByStatus[status.id] || []}
                onTaskClick={onTaskClick}
              />
            ))}
          </div>
        </div>

        {/* Drag Overlay - показывает карточку при перетаскивании */}
        <DragOverlay>
          {activeTask ? (
            <div className="rotate-3 opacity-90">
              <TaskCard task={activeTask} />
            </div>
          ) : null}
        </DragOverlay>
      </div>
    </DndContext>
  );
}

