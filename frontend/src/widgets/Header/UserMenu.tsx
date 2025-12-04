import { Link, useNavigate } from 'react-router-dom';
import { Shield, User, LogOut, Settings } from 'lucide-react';
import { useAuthStore } from '@/entities/session';
import { useCurrentUser } from '@/entities/user';
import { Avatar, Dropdown } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

/**
 * Dropdown меню пользователя в Header
 * Аватар, информация, навигация, выход
 */
export function UserMenu() {
  const navigate = useNavigate();
  const { logout } = useAuthStore();
  // Используем React Query для получения актуальных данных пользователя
  const { data: user, isLoading } = useCurrentUser();
  const isAdminUser = user?.Role?.Name === 'admin' || user?.Role?.ID === 1;

  const handleLogout = () => {
    logout();
    navigate(ROUTES.LOGIN);
  };

  // Если данные ещё загружаются, показываем fallback
  if (isLoading || !user) {
    return (
      <div className="flex items-center gap-2 p-1.5 rounded-lg">
        <div className="w-8 h-8 rounded-full bg-neutral-800 animate-pulse" />
        <div className="w-20 h-4 bg-neutral-800 rounded animate-pulse hidden md:block" />
      </div>
    );
  }

  return (
    <Dropdown>
      <Dropdown.Trigger>
        <button className="flex items-center gap-2 p-1.5 rounded-lg hover:bg-neutral-800 transition-colors">
          <Avatar
            file={user.avatar}
            fallback={user.Username}
            size="sm"
          />
          <span className="text-sm font-medium text-neutral-300 hidden md:block max-w-[120px] truncate">
            {user.Username}
          </span>
        </button>
      </Dropdown.Trigger>

      <Dropdown.Content align="end" className="w-56">
        {/* User info header */}
        <div className="px-3 py-2.5 border-b border-neutral-800">
          <p className="font-medium text-neutral-100 truncate">
            {user?.Username}
          </p>
          <p className="text-sm text-neutral-400 truncate">{user?.Email}</p>
          {user?.Role && (
            <p className="text-xs text-neutral-500 mt-1">
              Роль: {user.Role.Name}
            </p>
          )}
        </div>

        {/* Menu items */}
        <div className="p-1.5">
          <Dropdown.Item asChild>
            <Link
              to={ROUTES.PROFILE}
              className="flex items-center gap-2 w-full"
            >
              <User size={16} />
              <span>Профиль</span>
            </Link>
          </Dropdown.Item>

          <Dropdown.Item asChild>
            <Link
              to={ROUTES.PROFILE}
              className="flex items-center gap-2 w-full"
            >
              <Settings size={16} />
              <span>Настройки</span>
            </Link>
          </Dropdown.Item>

          {isAdminUser && (
            <Dropdown.Item asChild>
              <Link
                to={ROUTES.ADMIN}
                className="flex items-center gap-2 w-full text-primary-400"
              >
                <Shield size={16} />
                <span>Админ-панель</span>
              </Link>
            </Dropdown.Item>
          )}
        </div>

        {/* Logout */}
        <div className="p-1.5 border-t border-neutral-800">
          <Dropdown.Item
            onClick={handleLogout}
            className="flex items-center gap-2 w-full text-error hover:bg-error/10"
          >
            <LogOut size={16} />
            <span>Выйти</span>
          </Dropdown.Item>
        </div>
      </Dropdown.Content>
    </Dropdown>
  );
}

