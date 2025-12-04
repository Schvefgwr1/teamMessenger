import { Link } from 'react-router-dom';
import { useAuthStore } from '@/entities/session';
import { Card, Avatar, Badge } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import { MessageSquare, CheckSquare, User, ChevronRight } from 'lucide-react';

export function DashboardPage() {
  const { user } = useAuthStore();

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      {/* Welcome section */}
      <div className="flex items-center gap-4">
        <Avatar
          file={user?.avatar}
          fallback={user?.Username}
          size="xl"
        />
        <div>
          <h1 className="text-2xl font-bold text-neutral-100">
            Привет, {user?.Username}!
          </h1>
          <p className="text-neutral-400 mt-1">
            Добро пожаловать в Team Messenger
          </p>
          {user?.Role && (
            <Badge variant="primary" className="mt-2">
              {user.Role.Name}
            </Badge>
          )}
        </div>
      </div>

      {/* Quick Actions */}
      <section>
        <h2 className="text-lg font-semibold text-neutral-200 mb-4">
          Быстрые действия
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <QuickActionCard
            to={ROUTES.CHATS}
            icon={MessageSquare}
            iconBg="bg-primary-500/20"
            iconColor="text-primary-400"
            title="Чаты"
            description="Перейти к сообщениям"
          />
          <QuickActionCard
            to={ROUTES.TASKS}
            icon={CheckSquare}
            iconBg="bg-success/20"
            iconColor="text-success"
            title="Задачи"
            description="Управление задачами"
          />
          <QuickActionCard
            to={ROUTES.PROFILE}
            icon={User}
            iconBg="bg-warning/20"
            iconColor="text-warning"
            title="Профиль"
            description="Настройки аккаунта"
          />
        </div>
      </section>

      {/* User Info Card */}
      <section>
        <h2 className="text-lg font-semibold text-neutral-200 mb-4">
          Информация о пользователе
        </h2>
        <Card>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            <InfoItem label="ID" value={user?.ID} mono />
            <InfoItem label="Username" value={user?.Username} />
            <InfoItem label="Email" value={user?.Email} />
            <InfoItem label="Роль" value={user?.Role?.Name} />
            {user?.Gender && (
              <InfoItem label="Пол" value={user.Gender} />
            )}
            {user?.Age && (
              <InfoItem label="Возраст" value={String(user.Age)} />
            )}
            {user?.Description && (
              <div className="md:col-span-2">
                <InfoItem label="Описание" value={user.Description} />
              </div>
            )}
          </div>
        </Card>
      </section>

      {/* Recent Activity (placeholder) */}
      <section>
        <h2 className="text-lg font-semibold text-neutral-200 mb-4">
          Последняя активность
        </h2>
        <Card className="text-center py-12">
          <p className="text-neutral-500">
            Здесь будет отображаться ваша последняя активность
          </p>
        </Card>
      </section>
    </div>
  );
}

// Quick Action Card component
interface QuickActionCardProps {
  to: string;
  icon: React.ComponentType<{ size?: number | string; className?: string }>;
  iconBg: string;
  iconColor: string;
  title: string;
  description: string;
}

function QuickActionCard({
  to,
  icon: Icon,
  iconBg,
  iconColor,
  title,
  description,
}: QuickActionCardProps) {
  return (
    <Link to={to}>
      <Card className="flex items-center justify-between p-5 hover:bg-neutral-800/50 transition-colors group cursor-pointer">
        <div className="flex items-center gap-4">
          <div className={`w-12 h-12 rounded-xl ${iconBg} flex items-center justify-center`}>
            <Icon className={iconColor} size={24} />
          </div>
          <div>
            <h3 className="font-semibold text-neutral-100">{title}</h3>
            <p className="text-sm text-neutral-400">{description}</p>
          </div>
        </div>
        <ChevronRight
          size={20}
          className="text-neutral-600 group-hover:text-neutral-400 transition-colors"
        />
      </Card>
    </Link>
  );
}

// Info Item component
interface InfoItemProps {
  label: string;
  value?: string;
  mono?: boolean;
}

function InfoItem({ label, value, mono }: InfoItemProps) {
  return (
    <div>
      <span className="text-neutral-400">{label}:</span>
      <span className={`ml-2 text-neutral-100 ${mono ? 'font-mono text-xs' : ''}`}>
        {value || '—'}
      </span>
    </div>
  );
}
