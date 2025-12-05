import {useSortable} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';
import {Badge, Card} from '@/shared/ui';
import {cn, capitalizeFirst} from '@/shared/lib';
import type {Task} from '@/shared/types';

interface TaskCardProps {
    task: Task;
    onClick?: () => void;
}

/**
 * Карточка задачи для Kanban доски
 * Поддерживает drag-and-drop
 */
export function TaskCard({task, onClick}: TaskCardProps) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({
        id: task.id.toString(),
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
    };

    return (
        <div
            ref={setNodeRef}
            style={style}
            {...attributes}
            {...listeners}
            className={cn(
                'cursor-grab active:cursor-grabbing',
                isDragging && 'opacity-50 rotate-2'
            )}
        >
            <Card
                variant="elevated"
                className={cn(
                    'hover:shadow-lg transition-all duration-200',
                    'p-3 space-y-2'
                )}
                onClick={onClick}
            >
                {/* Заголовок */}
                <div className="flex items-start justify-between gap-2 min-w-full">
                    <Badge variant="primary" size="sm">
                        {task.id}
                    </Badge>
                    <h3 className="font-medium text-neutral-100 text-sm line-clamp-2 flex-1">
                        {capitalizeFirst(task.title)}
                    </h3>
                </div>
            </Card>
        </div>
    );
}

