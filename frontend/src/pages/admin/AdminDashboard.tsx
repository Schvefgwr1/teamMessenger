import { Shield, Key, Users, MessageSquare, ListChecks, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Card, Badge } from '@/shared/ui';
import { AdminPageLayout } from '@/widgets/AdminPanel';
import { ROUTES } from '@/shared/constants';
import { useRoles, usePermissions } from '@/entities/user';
import { useChatRoles, useChatPermissions } from '@/entities/chat';
import { useTaskStatuses } from '@/entities/task';
import { Skeleton } from '@/shared/ui';

/**
 * Главная страница админ-панели
 * Показывает статистику и быстрые ссылки
 */
export function AdminDashboard() {
  const { data: roles = [], isLoading: rolesLoading } = useRoles();
  const { data: permissions = [], isLoading: permissionsLoading } = usePermissions();
  const { data: chatRoles = [], isLoading: chatRolesLoading } = useChatRoles();
  const { data: chatPermissions = [], isLoading: chatPermissionsLoading } = useChatPermissions();
  const { data: taskStatuses = [], isLoading: taskStatusesLoading } = useTaskStatuses();

  const stats = [
    {
      label: 'Роли пользователей',
      value: roles.length,
      icon: Shield,
      color: 'text-primary-400',
      bgColor: 'bg-primary-500/20',
      link: ROUTES.ADMIN_ROLES,
      isLoading: rolesLoading,
    },
    {
      label: 'Разрешения',
      value: permissions.length,
      icon: Key,
      color: 'text-warning',
      bgColor: 'bg-warning/20',
      link: ROUTES.ADMIN_PERMISSIONS,
      isLoading: permissionsLoading,
    },
    {
      label: 'Роли чатов',
      value: chatRoles.length,
      icon: Users,
      color: 'text-success',
      bgColor: 'bg-success/20',
      link: ROUTES.ADMIN_CHAT_ROLES,
      isLoading: chatRolesLoading,
    },
    {
      label: 'Разрешения чатов',
      value: chatPermissions.length,
      icon: MessageSquare,
      color: 'text-info',
      bgColor: 'bg-info/20',
      link: ROUTES.ADMIN_CHAT_PERMISSIONS,
      isLoading: chatPermissionsLoading,
    },
    {
      label: 'Статусы задач',
      value: taskStatuses.length,
      icon: ListChecks,
      color: 'text-error',
      bgColor: 'bg-error/20',
      link: ROUTES.ADMIN_TASK_STATUSES,
      isLoading: taskStatusesLoading,
    },
  ];

  return (
    <AdminPageLayout
      title="Панель администратора"
      description="Управление системой и настройками"
      showBackButton={false}
    >
      {/* Статистика */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {stats.map((stat) => (
          <Link key={stat.link} to={stat.link}>
            <Card className="hover:bg-neutral-800/50 transition-colors group">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className={`w-12 h-12 rounded-xl ${stat.bgColor} flex items-center justify-center`}>
                    <stat.icon className={stat.color} size={24} />
                  </div>
                  <div>
                    <p className="text-sm text-neutral-400">{stat.label}</p>
                    {stat.isLoading ? (
                      <Skeleton className="h-6 w-12 mt-1" />
                    ) : (
                      <p className="text-2xl font-bold text-neutral-100 mt-1">
                        {stat.value}
                      </p>
                    )}
                  </div>
                </div>
                <ArrowRight
                  size={20}
                  className="text-neutral-600 group-hover:text-neutral-400 transition-colors"
                />
              </div>
            </Card>
          </Link>
        ))}
      </div>

      {/* Быстрые действия */}
      <Card>
        <Card.Header>
          <Card.Title>Быстрые действия</Card.Title>
        </Card.Header>
        <Card.Body>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <QuickActionLink
              to={ROUTES.ADMIN_ROLES}
              title="Управление ролями"
              description="Создание и редактирование ролей пользователей"
              icon={Shield}
            />
            <QuickActionLink
              to={ROUTES.ADMIN_CHAT_ROLES}
              title="Роли чатов"
              description="Настройка ролей для участников чатов"
              icon={Users}
            />
            <QuickActionLink
              to={ROUTES.ADMIN_TASK_STATUSES}
              title="Статусы задач"
              description="Управление статусами для Kanban-доски"
              icon={ListChecks}
            />
          </div>
        </Card.Body>
      </Card>
    </AdminPageLayout>
  );
}

interface QuickActionLinkProps {
  to: string;
  title: string;
  description: string;
  icon: React.ComponentType<{ size?: number; className?: string }>;
}

function QuickActionLink({ to, title, description, icon: Icon }: QuickActionLinkProps) {
  return (
    <Link
      to={to}
      className="flex items-start gap-4 p-4 rounded-lg hover:bg-neutral-800/50 transition-colors group"
    >
      <div className="w-10 h-10 rounded-lg bg-primary-500/20 flex items-center justify-center flex-shrink-0">
        <Icon size={20} className="text-primary-400" />
      </div>
      <div className="flex-1 min-w-0">
        <h3 className="font-semibold text-neutral-100 group-hover:text-primary-400 transition-colors">
          {title}
        </h3>
        <p className="text-sm text-neutral-400 mt-1">{description}</p>
      </div>
      <ArrowRight
        size={18}
        className="text-neutral-600 group-hover:text-neutral-400 transition-colors flex-shrink-0"
      />
    </Link>
  );
}

