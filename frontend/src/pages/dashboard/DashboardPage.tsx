import { useAuthStore } from '@/entities/session';
import { Button, Card, Avatar, Badge } from '@/shared/ui';
import { LogOut, MessageSquare, CheckSquare, User } from 'lucide-react';

export function DashboardPage() {
  const { user, logout } = useAuthStore();

  const handleLogout = async () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-neutral-950 p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Avatar
              file={user?.avatar}
              fallback={user?.Username}
              size="lg"
            />
            <div>
              <h1 className="text-2xl font-bold text-neutral-100">
                Привет, {user?.Username}!
              </h1>
              {user?.Role && (
                <Badge variant="primary" className="mt-1">
                  {user.Role.Name}
                </Badge>
              )}
            </div>
          </div>
          <Button variant="ghost" onClick={handleLogout} leftIcon={<LogOut size={18} />}>
            Выйти
          </Button>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <Card className="flex items-center gap-4 p-6 hover:bg-neutral-800 transition-colors cursor-pointer">
            <div className="w-12 h-12 rounded-xl bg-primary-500/20 flex items-center justify-center">
              <MessageSquare className="text-primary-400" size={24} />
            </div>
            <div>
              <h3 className="font-semibold text-neutral-100">Чаты</h3>
              <p className="text-sm text-neutral-400">Перейти к сообщениям</p>
            </div>
          </Card>

          <Card className="flex items-center gap-4 p-6 hover:bg-neutral-800 transition-colors cursor-pointer">
            <div className="w-12 h-12 rounded-xl bg-success/20 flex items-center justify-center">
              <CheckSquare className="text-success" size={24} />
            </div>
            <div>
              <h3 className="font-semibold text-neutral-100">Задачи</h3>
              <p className="text-sm text-neutral-400">Управление задачами</p>
            </div>
          </Card>

          <Card className="flex items-center gap-4 p-6 hover:bg-neutral-800 transition-colors cursor-pointer">
            <div className="w-12 h-12 rounded-xl bg-warning/20 flex items-center justify-center">
              <User className="text-warning" size={24} />
            </div>
            <div>
              <h3 className="font-semibold text-neutral-100">Профиль</h3>
              <p className="text-sm text-neutral-400">Настройки аккаунта</p>
            </div>
          </Card>
        </div>

        {/* User Info */}
        <Card>
          <Card.Header>
            <Card.Title>Информация о пользователе</Card.Title>
          </Card.Header>
          <Card.Body>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-neutral-400">ID:</span>
                <span className="ml-2 text-neutral-100 font-mono">{user?.ID}</span>
              </div>
              <div>
                <span className="text-neutral-400">Username:</span>
                <span className="ml-2 text-neutral-100">{user?.Username}</span>
              </div>
              <div>
                <span className="text-neutral-400">Email:</span>
                <span className="ml-2 text-neutral-100">{user?.Email}</span>
              </div>
              <div>
                <span className="text-neutral-400">Gender:</span>
                <span className="ml-2 text-neutral-100">{user?.Gender}</span>
              </div>
              {user?.Age && (
                <div>
                  <span className="text-neutral-400">Возраст:</span>
                  <span className="ml-2 text-neutral-100">{user.Age}</span>
                </div>
              )}
              {user?.Description && (
                <div className="col-span-2">
                  <span className="text-neutral-400">Описание:</span>
                  <span className="ml-2 text-neutral-100">{user.Description}</span>
                </div>
              )}
            </div>
          </Card.Body>
        </Card>
      </div>
    </div>
  );
}

