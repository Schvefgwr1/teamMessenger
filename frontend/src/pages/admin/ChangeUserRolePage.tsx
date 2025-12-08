import { useState } from 'react';
import { Search, UserCheck } from 'lucide-react';
import { Button, Input, Avatar, Card, Badge, Skeleton } from '@/shared/ui';
import { AdminPageLayout } from '@/widgets/AdminPanel';
import { useSearchUsers, useRoles, useUpdateUserRole } from '@/entities/user';
import { useDebounce } from '@/shared/hooks';
import { cn } from '@/shared/lib';
import { toast } from '@/shared/ui';
import type { UserSearchResult } from '@/shared/types';

/**
 * Страница изменения роли пользователя
 */
export function ChangeUserRolePage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedUser, setSelectedUser] = useState<UserSearchResult | null>(null);
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null);
  const debouncedQuery = useDebounce(searchQuery, 300);
  const { data: users, isLoading: isSearching } = useSearchUsers(debouncedQuery, searchQuery.length >= 2);
  const { data: roles = [], isLoading: rolesLoading } = useRoles();
  const updateUserRoleMutation = useUpdateUserRole();

  const handleSelectUser = (user: UserSearchResult) => {
    setSelectedUser(user);
    setSearchQuery('');
    // Находим роль пользователя, если она есть
    // В UserSearchResult нет роли, поэтому оставляем null
    setSelectedRoleId(null);
  };

  const handleSave = async () => {
    if (!selectedUser || !selectedRoleId) {
      toast.error('Выберите пользователя и роль');
      return;
    }

    try {
      await updateUserRoleMutation.mutateAsync({
        userId: selectedUser.id,
        roleId: selectedRoleId,
      });
      setSelectedUser(null);
      setSelectedRoleId(null);
    } catch (error) {
      // Ошибка обработана в mutation
    }
  };

  return (
    <AdminPageLayout
      title="Изменить роль пользователя"
      description="Поиск пользователя и изменение его роли"
    >
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Левая колонка: Поиск пользователя */}
        <Card variant="elevated" className="p-4 bg-neutral-900">
          <h3 className="text-lg font-semibold text-neutral-100 mb-3">Поиск пользователя</h3>
          
          <div className="space-y-3">
            <Input
              placeholder="Введите имя или email пользователя..."
              leftIcon={<Search size={18} />}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />

            {/* Результаты поиска */}
            {searchQuery.length >= 2 && (
              <div className="border border-neutral-800 rounded-lg bg-neutral-950 max-h-96 overflow-y-auto">
                {isSearching ? (
                  <div className="p-3 space-y-2">
                    {[1, 2, 3].map((i) => (
                      <div key={i} className="flex items-center gap-3">
                        <Skeleton variant="circular" className="w-10 h-10" />
                        <div className="flex-1 space-y-2">
                          <Skeleton className="h-4 w-32" />
                          <Skeleton className="h-3 w-48" />
                        </div>
                      </div>
                    ))}
                  </div>
                ) : users && users.length > 0 ? (
                  <div className="p-2">
                    {users.map((user) => (
                      <button
                        key={user.id}
                        type="button"
                        onClick={() => handleSelectUser(user)}
                        className={cn(
                          'w-full flex items-center gap-3 p-2 rounded-lg',
                          'hover:bg-neutral-800 transition-colors',
                          'text-left',
                          selectedUser?.id === user.id && 'bg-primary-500/10 border border-primary-500/20'
                        )}
                      >
                        <Avatar
                          file={user.avatarFile}
                          fallback={user.username}
                          size="md"
                        />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-neutral-100 truncate">
                            {user.username}
                          </p>
                          <p className="text-xs text-neutral-500 truncate">
                            {user.email}
                          </p>
                        </div>
                      </button>
                    ))}
                  </div>
                ) : (
                  <div className="p-3 text-center text-sm text-neutral-500">
                    Пользователи не найдены
                  </div>
                )}
              </div>
            )}

            {/* Выбранный пользователь */}
            {selectedUser && (
              <Card variant="outlined" className="p-3 bg-neutral-900">
                <div className="flex items-center gap-3">
                  <Avatar
                    file={selectedUser.avatarFile}
                    fallback={selectedUser.username}
                    size="md"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-neutral-100 truncate">
                      {selectedUser.username}
                    </p>
                    <p className="text-xs text-neutral-500 truncate">
                      {selectedUser.email}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      setSelectedUser(null);
                      setSelectedRoleId(null);
                    }}
                  >
                    Отменить
                  </Button>
                </div>
              </Card>
            )}
          </div>
        </Card>

        {/* Правая колонка: Выбор роли */}
        <Card variant="elevated" className="p-4 bg-neutral-900">
          <h3 className="text-lg font-semibold text-neutral-100 mb-3">Выбор роли</h3>
          
          {!selectedUser ? (
            <div className="text-center py-12 text-neutral-500">
              <UserCheck size={48} className="mx-auto mb-4 text-neutral-600" />
              <p className="text-sm">Сначала выберите пользователя</p>
            </div>
          ) : rolesLoading ? (
            <div className="space-y-2">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-16 w-full" />
              ))}
            </div>
          ) : roles.length === 0 ? (
            <div className="text-center py-12 text-neutral-500">
              <p className="text-sm">Роли не найдены</p>
            </div>
          ) : (
            <div className="space-y-2">
              {roles.map((role) => {
                const isSelected = selectedRoleId === role.ID;
                return (
                  <button
                    key={role.ID}
                    type="button"
                    onClick={() => setSelectedRoleId(role.ID)}
                    className={cn(
                      'w-full p-3 rounded-lg border-2 text-left',
                      'transition-all duration-200',
                      isSelected
                        ? 'border-primary-500 bg-primary-500/10'
                        : 'border-neutral-800 bg-neutral-950 hover:border-neutral-700 hover:bg-neutral-800/50'
                    )}
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-2">
                          <h4 className={cn(
                            'text-sm font-semibold',
                            isSelected ? 'text-primary-400' : 'text-neutral-100'
                          )}>
                            {role.Name}
                          </h4>
                          {isSelected && (
                            <Badge variant="primary" size="sm">
                              Выбрано
                            </Badge>
                          )}
                        </div>
                        {role.Description && (
                          <p className="text-xs text-neutral-500 mb-2">
                            {role.Description}
                          </p>
                        )}
                        {role.Permissions && role.Permissions.length > 0 && (
                          <div className="flex flex-wrap gap-1 mt-2">
                            {role.Permissions.slice(0, 3).map((perm) => (
                              <Badge key={perm.ID} variant="default" size="sm">
                                {perm.Name}
                              </Badge>
                            ))}
                            {role.Permissions.length > 3 && (
                              <Badge variant="default" size="sm">
                                +{role.Permissions.length - 3}
                              </Badge>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                  </button>
                );
              })}

              <div className="pt-3 border-t border-neutral-800 mt-3">
                <Button
                  onClick={handleSave}
                  disabled={!selectedRoleId}
                  isLoading={updateUserRoleMutation.isPending}
                  leftIcon={<UserCheck size={18} />}
                  className="w-full"
                >
                  Сохранить изменения
                </Button>
              </div>
            </div>
          )}
        </Card>
      </div>
    </AdminPageLayout>
  );
}

