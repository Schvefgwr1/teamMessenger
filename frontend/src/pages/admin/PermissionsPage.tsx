import { useState } from 'react';
import { Search } from 'lucide-react';
import { Input } from '@/shared/ui';
import { AdminPageLayout, DataTable } from '@/widgets/AdminPanel';
import { usePermissions } from '@/entities/user';
import type { Permission } from '@/shared/types';

/**
 * Страница просмотра разрешений системы
 * Только для просмотра (permissions системные, не редактируются)
 */
export function PermissionsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const { data: permissions = [], isLoading } = usePermissions();

  const filteredPermissions = permissions.filter((perm) =>
    perm.Name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    perm.Description?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const columns = [
    {
      key: 'id',
      header: 'ID',
      render: (permission: Permission) => (
        <span className="font-mono text-xs text-neutral-400">{permission.ID}</span>
      ),
      className: 'w-20',
    },
    {
      key: 'name',
      header: 'Название',
      render: (permission: Permission) => (
        <span className="font-medium text-neutral-100">{permission.Name}</span>
      ),
    },
    {
      key: 'description',
      header: 'Описание',
      render: (permission: Permission) => (
        <span className="text-neutral-400">{permission.Description || '—'}</span>
      ),
    },
  ];

  return (
    <AdminPageLayout
      title="Разрешения системы"
      description="Системные разрешения (только просмотр)"
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
      </div>

      {/* Таблица */}
      <DataTable
        data={filteredPermissions}
        columns={columns}
        isLoading={isLoading}
        emptyMessage="Разрешения не найдены"
      />
    </AdminPageLayout>
  );
}

