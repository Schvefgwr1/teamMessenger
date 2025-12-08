import { useState, useEffect } from 'react';
import { Plus, Search, Edit, Trash2, Check } from 'lucide-react';
import { Button, Input, Badge, Modal } from '@/shared/ui';
import { cn } from '@/shared/lib';
import { AdminPageLayout, DataTable } from '@/widgets/AdminPanel';
import {
  useRoles,
  usePermissions,
  useCreateRole,
  useUpdateRolePermissions,
  useDeleteRole,
} from '@/entities/user';
import { toast } from '@/shared/ui';
import type { Role, Permission } from '@/shared/types';

/**
 * Страница управления ролями пользователей
 */
export function RolesPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const { data: roles = [], isLoading } = useRoles();
  const { data: permissions = [] } = usePermissions();

  const filteredRoles = roles.filter((role) =>
    role.Name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    role.Description?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const columns = [
    {
      key: 'id',
      header: 'ID',
      render: (role: Role) => (
        <span className="font-mono text-xs text-neutral-400">{role.ID}</span>
      ),
      className: 'w-20',
    },
    {
      key: 'name',
      header: 'Название',
      render: (role: Role) => (
        <span className="font-medium text-neutral-100">{role.Name}</span>
      ),
    },
    {
      key: 'description',
      header: 'Описание',
      render: (role: Role) => (
        <span className="text-neutral-400">{role.Description || '—'}</span>
      ),
    },
    {
      key: 'permissions',
      header: 'Разрешения',
      render: (role: Role) => (
        <div className="flex flex-wrap gap-1">
          {role.Permissions?.length > 0 ? (
            role.Permissions.map((perm) => (
              <Badge key={perm.ID} variant="primary" size="sm">
                {perm.Name}
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
      render: (role: Role) => (
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => setEditingRole(role)}
          >
            <Edit size={16} />
          </Button>
          <DeleteRoleButton roleId={role.ID} roleName={role.Name} />
        </div>
      ),
      className: 'w-32',
    },
  ];

  return (
    <AdminPageLayout
      title="Роли пользователей"
      description="Управление ролями и их разрешениями"
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
        data={filteredRoles.map(role => ({ ...role, id: role.ID }))}
        columns={columns}
        isLoading={isLoading}
        emptyMessage="Роли не найдены"
      />

      {/* Модальное окно создания */}
      <CreateRoleModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        permissions={permissions}
      />

      {/* Модальное окно редактирования */}
      {editingRole && (
        <EditRoleModal
          open={!!editingRole}
          onOpenChange={(open) => !open && setEditingRole(null)}
          role={editingRole}
          permissions={permissions}
        />
      )}
    </AdminPageLayout>
  );
}

interface CreateRoleModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  permissions: Permission[];
}

function CreateRoleModal({ open, onOpenChange, permissions }: CreateRoleModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<number[]>([]);
  const createRole = useCreateRole();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      await createRole.mutateAsync({
        name,
        description: description || undefined,
        permissionIds: selectedPermissionIds.length > 0 ? selectedPermissionIds : undefined,
      });
      toast.success('Роль создана');
      onOpenChange(false);
      setName('');
      setDescription('');
      setSelectedPermissionIds([]);
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
      <Modal.Content title="Создать роль" size="lg">
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
            Описание
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full h-24 px-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
            placeholder="Описание роли..."
          />
        </div>

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
                  const isChecked = selectedPermissionIds.includes(permission.ID);
                  return (
                    <label
                      key={permission.ID}
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
                          onChange={() => togglePermission(permission.ID)}
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
                          {permission.Name}
                        </span>
                        {permission.Description && (
                          <p className="text-xs text-neutral-500 mt-1">
                            {permission.Description}
                          </p>
                        )}
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
            isLoading={createRole.isPending}
            disabled={!name.trim()}
          >
            Создать
          </Button>
        </div>
      </form>
      </Modal.Content>
    </Modal>
  );
}

interface EditRoleModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  role: Role;
  permissions: Permission[];
}

function EditRoleModal({ open, onOpenChange, role, permissions }: EditRoleModalProps) {
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<number[]>(
    role.Permissions?.map((p) => p.ID) || []
  );
  const updateRolePermissions = useUpdateRolePermissions();

  // Обновляем выбранные permissions при изменении роли
  useEffect(() => {
    if (open && role) {
      setSelectedPermissionIds(role.Permissions?.map((p) => p.ID) || []);
    }
  }, [open, role]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      await updateRolePermissions.mutateAsync({
        roleId: role.ID,
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
      <Modal.Content title={`Редактировать роль: ${role.Name}`} size="lg">
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
                    const isChecked = selectedPermissionIds.includes(permission.ID);
                    return (
                      <label
                        key={permission.ID}
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
                            onChange={() => togglePermission(permission.ID)}
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
                            {permission.Name}
                          </span>
                          {permission.Description && (
                            <p className="text-xs text-neutral-500 mt-1">
                              {permission.Description}
                            </p>
                          )}
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

interface DeleteRoleButtonProps {
  roleId: number;
  roleName: string;
}

function DeleteRoleButton({ roleId, roleName }: DeleteRoleButtonProps) {
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const deleteRole = useDeleteRole();

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
                onClick={() => deleteRole.mutate(roleId)}
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

