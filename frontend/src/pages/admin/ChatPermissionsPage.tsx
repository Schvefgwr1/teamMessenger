import { useState } from 'react';
import { Plus, Search, Trash2 } from 'lucide-react';
import { Button, Input, Modal } from '@/shared/ui';
import { AdminPageLayout, DataTable } from '@/widgets/AdminPanel';
import { useChatPermissions, chatRolesApi, chatRolesKeys } from '@/entities/chat';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from '@/shared/ui';
import type { ChatPermissionResponse } from '@/entities/chat';

/**
 * Страница управления разрешениями чатов
 */
export function ChatPermissionsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const { data: permissions = [], isLoading } = useChatPermissions();

  const filteredPermissions = permissions.filter((perm) =>
    perm.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const columns = [
    {
      key: 'id',
      header: 'ID',
      render: (permission: ChatPermissionResponse) => (
        <span className="font-mono text-xs text-neutral-400">{permission.id}</span>
      ),
      className: 'w-20',
    },
    {
      key: 'name',
      header: 'Название',
      render: (permission: ChatPermissionResponse) => (
        <span className="font-medium text-neutral-100">{permission.name}</span>
      ),
    },
    {
      key: 'actions',
      header: 'Действия',
      render: (permission: ChatPermissionResponse) => (
        <DeletePermissionButton permissionId={permission.id} permissionName={permission.name} />
      ),
      className: 'w-32',
    },
  ];

  return (
    <AdminPageLayout
      title="Разрешения чатов"
      description="Права, которые могут быть назначены ролям чатов"
    >
      {/* Toolbar */}
      <div className="flex justify-between items-center gap-4">
        <Input
          placeholder="Поиск разрешений..."
          leftIcon={<Search size={18} />}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="max-w-md"
        />
        <Button onClick={() => setIsCreateModalOpen(true)} leftIcon={<Plus size={18} />}>
          Создать разрешение
        </Button>
      </div>

      {/* Таблица */}
      <DataTable
        data={filteredPermissions}
        columns={columns}
        isLoading={isLoading}
        emptyMessage="Разрешения не найдены"
      />

      {/* Модальное окно создания */}
      <CreateChatPermissionModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
      />
    </AdminPageLayout>
  );
}

interface CreateChatPermissionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function CreateChatPermissionModal({ open, onOpenChange }: CreateChatPermissionModalProps) {
  const [name, setName] = useState('');
  const queryClient = useQueryClient();

  const createPermission = useMutation({
    mutationFn: async (data: { name: string }) => {
      const response = await chatRolesApi.createPermission(data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...chatRolesKeys.all, 'permissions'] });
      toast.success('Разрешение создано');
      onOpenChange(false);
      setName('');
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка создания разрешения';
      toast.error(message);
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createPermission.mutateAsync({ name });
  };

  return (
    <Modal open={open} onOpenChange={onOpenChange}>
      <Modal.Content title="Создать разрешение" size="sm">
        <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Название разрешения"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="Например: send_message"
        />

        <div className="flex justify-end gap-2 pt-4">
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Отмена
          </Button>
          <Button type="submit" isLoading={createPermission.isPending} disabled={!name.trim()}>
            Создать
          </Button>
        </div>
      </form>
      </Modal.Content>
    </Modal>
  );
}

interface DeletePermissionButtonProps {
  permissionId: number;
  permissionName: string;
}

function DeletePermissionButton({
  permissionId,
  permissionName,
}: DeletePermissionButtonProps) {
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const queryClient = useQueryClient();

  const deletePermission = useMutation({
    mutationFn: async () => {
      await chatRolesApi.deletePermission(permissionId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...chatRolesKeys.all, 'permissions'] });
      toast.success('Разрешение удалено');
      setIsDeleteModalOpen(false);
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка удаления разрешения';
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
        <Modal.Content title="Удалить разрешение?" size="sm">
          <div className="space-y-4">
            <p className="text-neutral-300">
              Вы уверены, что хотите удалить разрешение <strong>{permissionName}</strong>?
              Это действие нельзя отменить.
            </p>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => setIsDeleteModalOpen(false)}
                disabled={deletePermission.isPending}
              >
                Отмена
              </Button>
              <Button
                variant="danger"
                onClick={() => deletePermission.mutate()}
                isLoading={deletePermission.isPending}
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

