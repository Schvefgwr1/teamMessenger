import { useState, useEffect } from 'react';
import { Plus, Search, Edit, Trash2, Check } from 'lucide-react';
import { Button, Input, Badge, Modal } from '@/shared/ui';
import { cn } from '@/shared/lib';
import { AdminPageLayout, DataTable } from '@/widgets/AdminPanel';
import {
  useChatRoles,
  useChatPermissions,
  useUpdateChatRolePermissions,
  chatRolesApi,
  chatRolesKeys,
} from '@/entities/chat';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from '@/shared/ui';
import type { ChatRoleResponse, ChatPermissionResponse } from '@/entities/chat';

/**
 * Страница управления ролями чатов
 */
export function ChatRolesPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<ChatRoleResponse | null>(null);
  const { data: roles = [], isLoading } = useChatRoles();
  const { data: permissions = [] } = useChatPermissions();

  const filteredRoles = roles.filter((role) =>
    role.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const columns = [
    {
      key: 'id',
      header: 'ID',
      render: (role: ChatRoleResponse) => (
        <span className="font-mono text-xs text-neutral-400">{role.id}</span>
      ),
      className: 'w-20',
    },
    {
      key: 'name',
      header: 'Название',
      render: (role: ChatRoleResponse) => (
        <span className="font-medium text-neutral-100">{role.name}</span>
      ),
    },
    {
      key: 'permissions',
      header: 'Разрешения',
      render: (role: ChatRoleResponse) => (
        <div className="flex flex-wrap gap-1">
          {role.permissions?.length > 0 ? (
            role.permissions.map((perm) => (
              <Badge key={perm.id} variant="primary" size="sm">
                {perm.name}
              </Badge>
            ))
          ) : (
            <span className="text-xs text-neutral-500">Нет разрешений</span>
          )}
        </div>
      ),
    },
    {
      key: 'actions',
      header: 'Действия',
      render: (role: ChatRoleResponse) => (
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => setEditingRole(role)}
          >
            <Edit size={16} />
          </Button>
          <DeleteRoleButton roleId={role.id} roleName={role.name} />
        </div>
      ),
      className: 'w-32',
    },
  ];

  return (
    <AdminPageLayout
      title="Роли чатов"
      description="Управление ролями для участников чатов"
    >
      {/* Toolbar */}
      <div className="flex justify-between items-center gap-4">
        <Input
          placeholder="Поиск ролей..."
          leftIcon={<Search size={18} />}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="max-w-md"
        />
        <Button onClick={() => setIsCreateModalOpen(true)} leftIcon={<Plus size={18} />}>
          Создать роль
        </Button>
      </div>

      {/* Таблица */}
      <DataTable
        data={filteredRoles}
        columns={columns}
        isLoading={isLoading}
        emptyMessage="Роли не найдены"
      />

      {/* Модальное окно создания */}
      <CreateChatRoleModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        permissions={permissions}
      />

      {/* Модальное окно редактирования */}
      {editingRole && (
        <EditChatRoleModal
          open={!!editingRole}
          onOpenChange={(open) => !open && setEditingRole(null)}
          role={editingRole}
          permissions={permissions}
        />
      )}
    </AdminPageLayout>
  );
}

interface CreateChatRoleModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  permissions: ChatPermissionResponse[];
}

function CreateChatRoleModal({
  open,
  onOpenChange,
  permissions,
}: CreateChatRoleModalProps) {
  const [name, setName] = useState('');
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<number[]>([]);
  const queryClient = useQueryClient();

  const createRole = useMutation({
    mutationFn: async (data: { name: string; permissionIds?: number[] }) => {
      const response = await chatRolesApi.createRole(data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatRolesKeys.all });
      toast.success('Роль создана');
      onOpenChange(false);
      setName('');
      setSelectedPermissionIds([]);
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка создания роли';
      toast.error(message);
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createRole.mutateAsync({
      name,
      permissionIds: selectedPermissionIds.length > 0 ? selectedPermissionIds : undefined,
    });
  };

  const togglePermission = (permissionId: number) => {
    setSelectedPermissionIds((prev) =>
      prev.includes(permissionId)
        ? prev.filter((id) => id !== permissionId)
        : [...prev, permissionId]
    );
  };

  return (
    <Modal open={open} onOpenChange={onOpenChange}>
      <Modal.Content title="Создать роль чата" size="lg">
        <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Название роли"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="Например: moderator"
        />

        <div>
          <label className="text-sm font-medium text-neutral-300 mb-2 block">
            Разрешения
          </label>
          <div className="max-h-60 overflow-y-auto bg-neutral-900 rounded-lg border border-neutral-800">
            {permissions.length === 0 ? (
              <div className="p-4">
                <p className="text-sm text-neutral-500">Нет доступных разрешений</p>
              </div>
            ) : (
              <div className="p-2 space-y-1">
                {permissions.map((permission) => {
                  const isChecked = selectedPermissionIds.includes(permission.id);
                  return (
                    <label
                      key={permission.id}
                      className={cn(
                        'flex items-center gap-3 p-3 rounded-lg',
                        'hover:bg-neutral-800/50 cursor-pointer',
                        'transition-colors duration-150',
                        isChecked && 'bg-primary-500/10 hover:bg-primary-500/20'
                      )}
                    >
                      {/* Стилизованный checkbox */}
                      <div className="relative flex-shrink-0">
                        <input
                          type="checkbox"
                          checked={isChecked}
                          onChange={() => togglePermission(permission.id)}
                          className="sr-only"
                        />
                        <div
                          className={cn(
                            'w-5 h-5 rounded border-2 flex items-center justify-center',
                            'transition-all duration-200',
                            isChecked
                              ? 'bg-primary-500 border-primary-500'
                              : 'bg-neutral-800 border-neutral-700 hover:border-primary-500/50'
                          )}
                        >
                          {isChecked && (
                            <Check size={14} className="text-white" strokeWidth={3} />
                          )}
                        </div>
                      </div>
                      <span className={cn(
                        'text-sm font-medium',
                        isChecked ? 'text-primary-400' : 'text-neutral-100'
                      )}>
                        {permission.name}
                      </span>
                    </label>
                  );
                })}
              </div>
            )}
          </div>
        </div>

        <div className="flex justify-end gap-2 pt-4">
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Отмена
          </Button>
          <Button type="submit" isLoading={createRole.isPending} disabled={!name.trim()}>
            Создать
          </Button>
        </div>
      </form>
      </Modal.Content>
    </Modal>
  );
}

interface DeleteRoleButtonProps {
  roleId: number;
  roleName: string;
}

function DeleteRoleButton({ roleId, roleName }: DeleteRoleButtonProps) {
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const queryClient = useQueryClient();

  const deleteRole = useMutation({
    mutationFn: async () => {
      await chatRolesApi.deleteRole(roleId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: chatRolesKeys.all });
      toast.success('Роль удалена');
      setIsDeleteModalOpen(false);
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка удаления роли';
      toast.error(message);
    },
  });

  return (
    <>
      <Button
        variant="ghost"
        size="icon-sm"
        onClick={() => setIsDeleteModalOpen(true)}
      >
        <Trash2 size={16} className="text-error" />
      </Button>

      <Modal
        open={isDeleteModalOpen}
        onOpenChange={setIsDeleteModalOpen}
      >
        <Modal.Content title="Удалить роль?" size="sm">
          <div className="space-y-4">
            <p className="text-neutral-300">
              Вы уверены, что хотите удалить роль <strong>{roleName}</strong>?
              Это действие нельзя отменить.
            </p>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => setIsDeleteModalOpen(false)}
                disabled={deleteRole.isPending}
              >
                Отмена
              </Button>
              <Button
                variant="danger"
                onClick={() => deleteRole.mutate()}
                isLoading={deleteRole.isPending}
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

interface EditChatRoleModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  role: ChatRoleResponse;
  permissions: ChatPermissionResponse[];
}

function EditChatRoleModal({ open, onOpenChange, role, permissions }: EditChatRoleModalProps) {
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<number[]>(
    role.permissions?.map((p) => p.id) || []
  );
  const updateRolePermissions = useUpdateChatRolePermissions();

  // Обновляем выбранные permissions при изменении роли
  useEffect(() => {
    if (open && role) {
      setSelectedPermissionIds(role.permissions?.map((p) => p.id) || []);
    }
  }, [open, role]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      await updateRolePermissions.mutateAsync({
        roleId: role.id,
        permissionIds: selectedPermissionIds,
      });
      onOpenChange(false);
    } catch (error) {
      // Ошибка обработана в mutation
    }
  };

  const togglePermission = (permissionId: number) => {
    setSelectedPermissionIds((prev) =>
      prev.includes(permissionId)
        ? prev.filter((id) => id !== permissionId)
        : [...prev, permissionId]
    );
  };

  return (
    <Modal open={open} onOpenChange={onOpenChange}>
      <Modal.Content title={`Редактировать роль: ${role.name}`} size="lg">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-sm font-medium text-neutral-300 mb-2 block">
              Разрешения
            </label>
            <div className="max-h-60 overflow-y-auto bg-neutral-900 rounded-lg border border-neutral-800">
              {permissions.length === 0 ? (
                <div className="p-4">
                  <p className="text-sm text-neutral-500">Нет доступных разрешений</p>
                </div>
              ) : (
                <div className="p-2 space-y-1">
                  {permissions.map((permission) => {
                    const isChecked = selectedPermissionIds.includes(permission.id);
                    return (
                      <label
                        key={permission.id}
                        className={cn(
                          'flex items-start gap-3 p-3 rounded-lg',
                          'hover:bg-neutral-800/50 cursor-pointer',
                          'transition-colors duration-150',
                          isChecked && 'bg-primary-500/10 hover:bg-primary-500/20'
                        )}
                      >
                        {/* Стилизованный checkbox */}
                        <div className="relative flex-shrink-0 mt-0.5">
                          <input
                            type="checkbox"
                            checked={isChecked}
                            onChange={() => togglePermission(permission.id)}
                            className="sr-only"
                          />
                          <div
                            className={cn(
                              'w-5 h-5 rounded border-2 flex items-center justify-center',
                              'transition-all duration-200',
                              isChecked
                                ? 'bg-primary-500 border-primary-500'
                                : 'bg-neutral-800 border-neutral-700 hover:border-primary-500/50'
                            )}
                          >
                            {isChecked && (
                              <Check size={14} className="text-white" strokeWidth={3} />
                            )}
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <span className={cn(
                            'text-sm font-medium block',
                            isChecked ? 'text-primary-400' : 'text-neutral-100'
                          )}>
                            {permission.name}
                          </span>
                        </div>
                      </label>
                    );
                  })}
                </div>
              )}
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Отмена
            </Button>
            <Button
              type="submit"
              isLoading={updateRolePermissions.isPending}
            >
              Сохранить
            </Button>
          </div>
        </form>
      </Modal.Content>
    </Modal>
  );
}

