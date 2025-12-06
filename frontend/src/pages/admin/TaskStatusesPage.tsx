import { useState } from 'react';
import { Plus, Search, Trash2 } from 'lucide-react';
import { Button, Input, Badge, Modal } from '@/shared/ui';
import { AdminPageLayout, DataTable } from '@/widgets/AdminPanel';
import { useTaskStatuses, taskApi, taskKeys } from '@/entities/task';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from '@/shared/ui';
import type { TaskStatus } from '@/shared/types';

/**
 * Страница управления статусами задач
 */
export function TaskStatusesPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const { data: statuses = [], isLoading } = useTaskStatuses();

  const filteredStatuses = statuses.filter((status) =>
    status.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const columns = [
    {
      key: 'id',
      header: 'ID',
      render: (status: TaskStatus) => (
        <span className="font-mono text-xs text-neutral-400">{status.id}</span>
      ),
      className: 'w-20',
    },
    {
      key: 'name',
      header: 'Название',
      render: (status: TaskStatus) => (
        <div className="flex items-center gap-3">
          <Badge variant="primary" size="sm">
            {status.name}
          </Badge>
        </div>
      ),
    },
    {
      key: 'actions',
      header: 'Действия',
      render: (status: TaskStatus) => (
        <DeleteStatusButton statusId={status.id} statusName={status.name} />
      ),
      className: 'w-32',
    },
  ];

  return (
    <AdminPageLayout
      title="Статусы задач"
      description="Статусы для Kanban-доски задач"
    >
      {/* Toolbar */}
      <div className="flex justify-between items-center gap-4">
        <Input
          placeholder="Поиск статусов..."
          leftIcon={<Search size={18} />}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="max-w-md"
        />
        <Button onClick={() => setIsCreateModalOpen(true)} leftIcon={<Plus size={18} />}>
          Добавить статус
        </Button>
      </div>

      {/* Таблица */}
      <DataTable
        data={filteredStatuses}
        columns={columns}
        isLoading={isLoading}
        emptyMessage="Статусы не найдены"
      />

      {/* Модальное окно создания */}
      <CreateStatusModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
      />
    </AdminPageLayout>
  );
}

interface CreateStatusModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function CreateStatusModal({ open, onOpenChange }: CreateStatusModalProps) {
  const [name, setName] = useState('');
  const queryClient = useQueryClient();

  const createStatus = useMutation({
    mutationFn: async (data: { name: string }) => {
      const response = await taskApi.createStatus(data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: taskKeys.statuses() });
      toast.success('Статус создан');
      onOpenChange(false);
      setName('');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка создания статуса';
      toast.error(message);
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createStatus.mutateAsync({ name });
  };

  return (
    <Modal open={open} onOpenChange={onOpenChange}>
      <Modal.Content title="Добавить статус" size="sm">
        <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Название статуса"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="Например: in_progress"
        />

        <div className="flex justify-end gap-2 pt-4">
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Отмена
          </Button>
          <Button type="submit" isLoading={createStatus.isPending} disabled={!name.trim()}>
            Создать
          </Button>
        </div>
      </form>
      </Modal.Content>
    </Modal>
  );
}

interface DeleteStatusButtonProps {
  statusId: number;
  statusName: string;
}

function DeleteStatusButton({ statusId, statusName }: DeleteStatusButtonProps) {
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const queryClient = useQueryClient();

  const deleteStatus = useMutation({
    mutationFn: async () => {
      await taskApi.deleteStatus(statusId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: taskKeys.statuses() });
      toast.success('Статус удалён');
      setIsDeleteModalOpen(false);
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка удаления статуса';
      toast.error(message);
    },
  });

  return (
    <>
      <Button variant="ghost" size="icon-sm" onClick={() => setIsDeleteModalOpen(true)}>
        <Trash2 size={16} className="text-error" />
      </Button>

      <Modal
        open={isDeleteModalOpen}
        onOpenChange={setIsDeleteModalOpen}
      >
        <Modal.Content title="Удалить статус?" size="sm">
          <div className="space-y-4">
            <p className="text-neutral-300">
              Вы уверены, что хотите удалить статус <strong>{statusName}</strong>?
              Это действие нельзя отменить.
            </p>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => setIsDeleteModalOpen(false)}
                disabled={deleteStatus.isPending}
              >
                Отмена
              </Button>
              <Button
                variant="danger"
                onClick={() => deleteStatus.mutate()}
                isLoading={deleteStatus.isPending}
              >
                Удалить
              </Button>
            </div>
          </div>
        </Modal.Content>
      </Modal>
    </>
  );
}

