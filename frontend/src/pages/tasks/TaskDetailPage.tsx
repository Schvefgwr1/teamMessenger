import {Link, useNavigate, useParams} from 'react-router-dom';
import {ArrowLeft, Calendar, Clock, ExternalLink, FileText, MessageSquare} from 'lucide-react';
import {Avatar, Badge, Button, buttonVariants, Card, Skeleton} from '@/shared/ui';
import {useTask, useTaskStatuses, useUpdateTaskStatus} from '@/entities/task';
import {useUserById} from '@/entities/user';
import {cn, formatDate, getFileUrl, isValidUuid} from '@/shared/lib';
import {ROUTES} from '@/shared/constants';

/**
 * Детальная страница задачи
 */
export function TaskDetailPage() {
    const {taskId} = useParams<{ taskId: string }>();
    const navigate = useNavigate();
    const taskIdNum = taskId ? parseInt(taskId, 10) : undefined;

    const {data: task, isLoading, error: taskError} = useTask(taskIdNum);
    const {data: statuses = []} = useTaskStatuses();
    const updateTaskStatus = useUpdateTaskStatus();

    // Загружаем данные создателя
    const {data: creator} = useUserById(task?.creatorId);

    const handleStatusChange = (newStatusId: number) => {
        if (!taskIdNum) return;
        console.log('TaskDetailPage handleStatusChange:', {
            taskId: taskIdNum,
            currentStatusId: task?.status?.id,
            newStatusId,
            newStatusName: statuses.find(s => s.id === newStatusId)?.name,
            allStatuses: statuses.map(s => ({id: s.id, name: s.name}))
        });
        updateTaskStatus.mutate({taskId: taskIdNum, statusId: newStatusId});
    };

    // Отладочная информация
    console.log('TaskDetailPage render:', {
        taskIdNum,
        task,
        taskError,
        statuses,
        hasStatus: !!task?.status,
        statusId: task?.status?.id,
        chatId: task?.chatId,
        hasChatId: !!task?.chatId,
        chatIdType: typeof task?.chatId,
    });

    if (isLoading) {
        return <TaskDetailSkeleton/>;
    }

    if (taskError) {
        return (
            <div className="flex flex-col items-center justify-center py-12">
                <p className="text-error mb-4">Ошибка загрузки задачи</p>
                <p className="text-sm text-neutral-500 mb-4">
                    {taskError.message || 'Неизвестная ошибка'}
                </p>
                <Button onClick={() => navigate('/tasks')}>Вернуться к задачам</Button>
            </div>
        );
    }

    if (!task) {
        return (
            <div className="flex flex-col items-center justify-center py-12">
                <p className="text-neutral-400 mb-4">Задача не найдена</p>
                <Button onClick={() => navigate('/tasks')}>Вернуться к задачам</Button>
            </div>
        );
    }

    // Проверяем что статус загружен
    if (!task.status || !task.status.id) {
        return (
            <div className="flex flex-col items-center justify-center py-12">
                <p className="text-neutral-400 mb-4">Ошибка: статус задачи не загружен</p>
                <p className="text-sm text-neutral-500 mb-4">
                    Данные задачи: {JSON.stringify(task, null, 2)}
                </p>
                <Button onClick={() => navigate('/tasks')}>Вернуться к задачам</Button>
            </div>
        );
    }

    // Определяем цвет статуса для Badge
    const getStatusVariant = (statusName: string): 'default' | 'primary' | 'success' | 'warning' | 'error' => {
        const name = statusName.toLowerCase();
        if (name.includes('created') || name.includes('создан')) return 'primary';
        if (name.includes('in_progress') || name.includes('в работе')) return 'warning';
        if (name.includes('completed') || name.includes('завершен')) return 'success';
        if (name.includes('canceled') || name.includes('отменен')) return 'error';
        return 'default';
    };

    return (
        <div className="max-w-4xl mx-auto space-y-6 pb-8">
            {/* Header */}
            <Card variant="elevated">
                <div className="flex items-center gap-4">
                    <div className="flex-1">
                        <h1 className="text-2xl font-bold text-neutral-100 mb-2">{task.title}</h1>
                        <div className="flex items-center gap-3">
                            <Badge variant={getStatusVariant(task.status.name)} size="md">
                                {task.status.name}
                            </Badge>
                            {updateTaskStatus.isPending && (
                                <span className="text-xs text-neutral-500 flex items-center gap-1">
                                    <Clock size={12}/>
                                    Обновление...
                                </span>
                            )}
                        </div>
                    </div>
                    <Button
                        variant="ghost"
                        size="md"
                        onClick={() => navigate(ROUTES.TASKS)}
                        leftIcon={<ArrowLeft size={20}/>}
                    >
                        Назад
                    </Button>
                </div>
            </Card>

            {/* Основная информация */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Основной контент */}
                <div className="lg:col-span-2 space-y-6">
                    {/* Описание */}
                    {task.description ? (
                        <Card variant="elevated">
                            <Card.Header>
                                <Card.Title>Описание</Card.Title>
                            </Card.Header>
                            <Card.Body>
                                <p className="text-neutral-300 whitespace-pre-wrap leading-relaxed">
                                    {task.description}
                                </p>
                            </Card.Body>
                        </Card>
                    ) : (
                        <Card variant="outlined">
                            <Card.Body>
                                <p className="text-neutral-500 text-sm italic">Описание отсутствует</p>
                            </Card.Body>
                        </Card>
                    )}

                    {/* Прикрепленные файлы */}
                    {task.files && task.files.length > 0 && (
                        <Card variant="elevated">
                            <Card.Title className="flex items-center gap-2">
                                Прикрепленные файлы
                                <Badge variant="primary" size="sm">{task.files.length}</Badge>
                            </Card.Title>
                            <Card.Body>
                                <div className="space-y-2">
                                    {task.files.map((taskFile, index) => {
                                        const fileUrl = taskFile.file?.url
                                            ? getFileUrl(taskFile.file)
                                            : null;

                                        return (
                                            <div
                                                key={index}
                                                className="flex items-center gap-3 rounded-lg bg-neutral-800/50 hover:bg-neutral-800 transition-colors border border-neutral-800"
                                            >
                                                <div className="p-2 rounded-lg bg-primary-500/10">
                                                    <FileText size={18} className="text-primary-400"/>
                                                </div>
                                                <div className="flex-1 min-w-0">
                                                    <p className="text-sm font-medium text-neutral-200 truncate">
                                                        {taskFile.file?.name || `Файл #${taskFile.fileId}`}
                                                    </p>
                                                    {taskFile.file?.fileType && (
                                                        <p className="text-xs text-neutral-500 mt-0.5">
                                                            {taskFile.file.fileType.name}
                                                        </p>
                                                    )}
                                                </div>
                                                {fileUrl && (
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                    >
                                                        <a
                                                            href={fileUrl}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                            className="flex items-center gap-1"
                                                        >
                                                            <ExternalLink size={14}/>
                                                            Открыть
                                                        </a>
                                                    </Button>
                                                )}
                                            </div>
                                        );
                                    })}
                                </div>
                            </Card.Body>
                        </Card>
                    )}
                </div>

                {/* Боковая панель */}
                <div className="space-y-6">
                    {/* Статус */}
                    <Card variant="elevated">
                        <Card.Title>Статус задачи</Card.Title>
                        <Card.Footer>
                            <select
                                value={task.status?.id ?? ''}
                                onChange={(e) => handleStatusChange(parseInt(e.target.value, 10))}
                                className="w-full px-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all"
                                disabled={updateTaskStatus.isPending}
                            >
                                {statuses.map((status) => (
                                    <option key={status.id} value={status.id}>
                                        {status.name}
                                    </option>
                                ))}
                            </select>
                        </Card.Footer>
                    </Card>

                    {/* Метаданные */}
                    <Card variant="elevated">
                        <Card.Title>Информация</Card.Title>
                        <Card.Body className="space-y-4">
                            {/* Создатель */}
                            {creator && (
                                <div
                                    className="flex items-center gap-3 rounded-lg bg-neutral-800/50 border border-neutral-800">
                                    <Avatar
                                        file={creator.avatar}
                                        fallback={creator.Username}
                                        size="sm"
                                    />
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center gap-2">
                                            <span className="text-sm font-medium text-neutral-200 truncate">
                                                {creator.Username}
                                            </span>
                                        </div>
                                        <p className="text-xs text-neutral-500 mb-1">Создатель</p>
                                    </div>
                                </div>
                            )}

                            {/* Связанный чат */}
                            <div className="rounded-lg bg-neutral-800/50 border border-neutral-800">
                                <div className="flex items-center gap-2 mb-2">
                                    <MessageSquare size={14} className="text-neutral-500"/>
                                    <p className="text-xs text-neutral-500">Связанный чат</p>
                                </div>
                                {task.chatId && isValidUuid(task.chatId) ? (
                                    <Link
                                        to={ROUTES.CHAT_DETAIL(task.chatId)}
                                        className={cn(
                                            buttonVariants({variant: 'secondary', size: 'sm'}),
                                            'w-full flex items-center justify-center gap-2 hover:bg-neutral-700 transition-colors'
                                        )}
                                    >
                                        <MessageSquare size={16}/>
                                        Перейти к чату
                                        <ExternalLink size={14}/>
                                    </Link>
                                ) : (
                                    <div className="text-xs text-neutral-500 italic py-2">
                                        Чат не прикреплен
                                    </div>
                                )}
                            </div>

                            {/* Дата создания */}
                            <div
                                className="flex items-center gap-3 rounded-lg bg-neutral-800/50 border border-neutral-800">
                                <div className="p-2 rounded-lg bg-primary-500/10">
                                    <Calendar size={18} className="text-primary-400"/>
                                </div>
                                 <div className="flex-1 min-w-0 flex flex-col">
                                    <p className="text-xs text-neutral-500 mb-1">Создана</p>
                                    <div className="flex flex-row gap-1">
                                        <p className="text-sm font-medium text-neutral-200">
                                            {formatDate(task.createdAt, 'dd.MM.yyyy')}
                                        </p>
                                        <p className="text-sm text-neutral-500">
                                            {formatDate(task.createdAt, 'HH:mm')}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </Card.Body>
                    </Card>
                </div>
            </div>
        </div>
    );
}

/**
 * Skeleton для загрузки детальной страницы
 */
function TaskDetailSkeleton() {
    return (
        <div className="max-w-4xl mx-auto space-y-6 pb-8">
            <Card variant="elevated">
                <div className="flex items-center gap-4">
                    <Skeleton className="h-10 w-10 rounded-lg"/>
                    <div className="flex-1">
                        <Skeleton className="h-8 w-64 mb-2"/>
                        <Skeleton className="h-6 w-24"/>
                    </div>
                </div>
            </Card>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2 space-y-6">
                    <Card variant="elevated">
                        <Card.Header>
                            <Skeleton className="h-6 w-32"/>
                        </Card.Header>
                        <Card.Body>
                            <Skeleton className="h-24 w-full"/>
                        </Card.Body>
                    </Card>
                </div>
                <div className="space-y-6">
                    <Card variant="elevated">
                        <Card.Header>
                            <Skeleton className="h-6 w-32"/>
                        </Card.Header>
                        <Card.Body>
                            <Skeleton className="h-10 w-full"/>
                        </Card.Body>
                    </Card>
                    <Card variant="elevated">
                        <Card.Header>
                            <Skeleton className="h-6 w-32"/>
                        </Card.Header>
                        <Card.Body className="space-y-4">
                            <Skeleton className="h-16 w-full"/>
                            <Skeleton className="h-16 w-full"/>
                        </Card.Body>
                    </Card>
                </div>
            </div>
        </div>
    );
}

