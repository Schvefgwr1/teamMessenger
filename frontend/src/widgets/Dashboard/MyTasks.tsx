import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import { CheckSquare, ChevronRight, Clock } from 'lucide-react';
import { Card, Badge, Skeleton } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import { useUserTasks } from '@/entities/task';
import { formatTime, capitalizeFirst } from '@/shared/lib';
import { parseISO, subDays, isAfter } from 'date-fns';
import type { Task } from '@/shared/types';

const TASKS_LIMIT = 5;
const DAYS_AGO = 1; // Последние сутки

/**
 * Виджет новых задач для Dashboard
 * Показывает задачи, созданные за последние сутки
 */
export function MyTasks() {
  const { data: tasks = [], isLoading, error } = useUserTasks({
    limit: 100, // Загружаем достаточно для фильтрации
    offset: 0,
  });

  // Фильтруем задачи, созданные за последние сутки
  const recentTasks = useMemo(() => {
    if (!tasks || tasks.length === 0) return [];
    
    const oneDayAgo = subDays(new Date(), DAYS_AGO);
    
    const filtered = tasks
      .filter((task) => {
        // Проверяем наличие createdAt
        if (!task.createdAt || task.createdAt.trim() === '') {
          console.warn('Task without createdAt:', task);
          return false;
        }
        
        try {
          const taskDate = parseISO(task.createdAt);
          const isValid = isAfter(taskDate, oneDayAgo);
          if (!isValid) {
            console.log('Task filtered out (too old):', {
              taskId: task.id,
              createdAt: task.createdAt,
              oneDayAgo: oneDayAgo.toISOString(),
            });
          }
          return isValid;
        } catch (error) {
          console.error('Error parsing task date:', {
            taskId: task.id,
            createdAt: task.createdAt,
            error,
          });
          return false;
        }
      })
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, TASKS_LIMIT);
    
    console.log('MyTasks - filtered tasks:', {
      totalTasks: tasks.length,
      filteredCount: filtered.length,
      oneDayAgo: oneDayAgo.toISOString(),
    });
    
    return filtered;
  }, [tasks]);

  if (error) {
    return (
      <Card>
        <Card.Header>
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-neutral-200 flex items-center gap-2">
              <CheckSquare size={20} className="text-success" />
              Новые задачи
            </h3>
          </div>
        </Card.Header>
        <Card.Body>
          <p className="text-sm text-neutral-500 text-center py-4">
            Ошибка загрузки задач
          </p>
        </Card.Body>
      </Card>
    );
  }

  return (
    <Card>
      <Card.Header>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-neutral-200 flex items-center gap-2">
            <CheckSquare size={20} className="text-success" />
            Новые задачи
          </h3>
          {tasks.length > TASKS_LIMIT && (
            <Link
              to={ROUTES.TASKS}
              className="text-sm text-primary-400 hover:text-primary-300 transition-colors flex items-center gap-1"
            >
              Показать все
              <ChevronRight size={16} />
            </Link>
          )}
        </div>
      </Card.Header>
      <Card.Body>
        {isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="space-y-2">
                <div className="flex items-center justify-between">
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-5 w-16 rounded-full" />
                </div>
                <Skeleton className="h-3 w-1/2" />
              </div>
            ))}
          </div>
        ) : recentTasks.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-3">
            {recentTasks.map((task) => (
              <TaskItem key={task.id} task={task} />
            ))}
          </div>
        )}
      </Card.Body>
    </Card>
  );
}

interface TaskItemProps {
  task: Task;
}

function TaskItem({ task }: TaskItemProps) {
  const getStatusVariant = (statusName: string): 'primary' | 'success' | 'warning' | 'error' => {
    const name = statusName.toLowerCase();
    if (name.includes('создан') || name.includes('created')) return 'primary';
    if (name.includes('выполнен') || name.includes('completed') || name.includes('done')) return 'success';
    if (name.includes('в работе') || name.includes('in progress') || name.includes('progress')) return 'warning';
    if (name.includes('отменен') || name.includes('canceled') || name.includes('cancelled')) return 'error';
    return 'primary';
  };

  return (
    <Link
      to={ROUTES.TASK_DETAIL(task.id)}
      className="block p-3 rounded-lg hover:bg-neutral-800/50 transition-colors group"
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h4 className="font-medium text-neutral-100 truncate group-hover:text-primary-400 transition-colors">
              {capitalizeFirst(task.title)}
            </h4>
            <Badge variant={getStatusVariant(task.status.name)} size="sm">
              {task.status.name}
            </Badge>
          </div>
          {task.description && (
            <p className="text-xs text-neutral-400 line-clamp-2 mt-1">
              {task.description}
            </p>
          )}
          <div className="flex items-center gap-3 mt-2">
            <span className="flex items-center gap-1 text-xs text-neutral-500">
              <Clock size={12} />
              {formatTime(task.createdAt)}
            </span>
          </div>
        </div>
      </div>
    </Link>
  );
}

function EmptyState() {
  return (
    <div className="text-center py-8">
      <CheckSquare size={48} className="mx-auto text-neutral-600 mb-3" />
      <p className="text-sm text-neutral-500">
        У вас пока нет задач
      </p>
      <Link
        to={ROUTES.TASKS}
        className="inline-block mt-3 text-sm text-primary-400 hover:text-primary-300 transition-colors"
      >
        Создать задачу
      </Link>
    </div>
  );
}

