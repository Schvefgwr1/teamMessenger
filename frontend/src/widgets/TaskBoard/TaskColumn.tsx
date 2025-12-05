import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Card, Skeleton } from '@/shared/ui';
import { cn, capitalizeFirst } from '@/shared/lib';
import { TaskCard } from './TaskCard';
import type { Task, TaskStatus } from '@/shared/types';

interface TaskColumnProps {
  status: TaskStatus;
  tasks: Task[];
  onTaskClick?: (taskId: number) => void;
}

/**
 * Колонка Kanban доски для одного статуса
 * Поддерживает drop зону для drag-and-drop
 */
export function TaskColumn({ status, tasks, onTaskClick }: TaskColumnProps) {
  const { setNodeRef, isOver } = useDroppable({
    id: status.id.toString(),
  });

  return (
    <div className="flex flex-col min-w-[280px] max-w-[280px]">
      {/* Заголовок колонки */}
      <div className="flex flex-row justify-between mb-2 px-2">
        <h3 className="text-sm font-semibold text-neutral-300 mb-1">
          {capitalizeFirst(status.name)}
        </h3>
        <div className="flex items-center gap-2">
          <span className="text-xs text-neutral-500">
            {tasks.length} {tasks.length === 1 ? 'задача' : 'задач'}
          </span>
        </div>
      </div>

      {/* Список задач */}
      <div
        ref={setNodeRef}
        className={cn(
          'rounded-lg p-2 transition-colors min-h-[422px]',
          'bg-neutral-900/50 border border-neutral-800',
          isOver && 'bg-primary-500/10 border-primary-500/50'
        )}
      >
        <SortableContext
          items={tasks.map((t) => t.id.toString())}
          strategy={verticalListSortingStrategy}
        >
          <div className="space-y-2">
            {tasks.length === 0 ? (
              <EmptyState />
            ) : (
              tasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onClick={() => onTaskClick?.(task.id)}
                />
              ))
            )}
          </div>
        </SortableContext>
      </div>
    </div>
  );
}

/**
 * Пустое состояние колонки
 */
function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-8 px-4 text-center">
      <div className="w-12 h-12 rounded-lg bg-neutral-800 flex items-center justify-center mb-3">
        <span className="text-2xl text-neutral-600">📋</span>
      </div>
      <p className="text-sm text-neutral-500">
        Нет задач в этом статусе
      </p>
    </div>
  );
}

/**
 * Skeleton для загрузки колонки
 */
export function TaskColumnSkeleton() {
  return (
    <div className="flex flex-col h-full min-w-[280px] max-w-[280px]">
      <div className="mb-4">
        <Skeleton className="h-5 w-24 mb-2" />
        <Skeleton className="h-4 w-16" />
      </div>
      <div className="flex-1 rounded-lg p-2 bg-neutral-900/50 border border-neutral-800 space-y-2">
        {[1, 2, 3].map((i) => (
          <Card key={i} variant="elevated" className="p-3 space-y-2">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-3 w-3/4" />
            <Skeleton className="h-3 w-1/2" />
          </Card>
        ))}
      </div>
    </div>
  );
}

