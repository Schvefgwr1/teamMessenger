import { NavLink } from 'react-router-dom';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { cn } from '@/shared/lib/cn';
import { useCurrentUser } from '@/entities/user';
import { useSidebarStore } from './model/sidebarStore';
import { userNavigation } from './navigation';

/**
 * Sidebar для основного приложения
 * С навигацией и возможностью сворачивания
 */
export function Sidebar() {
  const { isCollapsed, toggle } = useSidebarStore();
  // Используем React Query для получения актуальных данных пользователя
  const { data: user, isLoading } = useCurrentUser();

  // Функция проверки permissions из актуальных данных пользователя
  const hasPermission = (permissionName: string): boolean => {
    if (!user?.Role?.Permissions) return false;
    return user.Role.Permissions.some((p) => p.Name === permissionName);
  };

  // Фильтруем навигацию по permissions
  // Если данные загружаются, показываем все элементы (чтобы не было пустого sidebar)
  const filteredNavigation = isLoading
    ? userNavigation
    : userNavigation.filter(
        (item) => !item.permission || hasPermission(item.permission)
      );

  return (
    <aside
      className={cn(
        'fixed left-0 top-16 bottom-0 z-30',
        'bg-neutral-900 border-r border-neutral-800',
        'flex flex-col transition-all duration-300 ease-in-out',
        isCollapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Navigation */}
      <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
        {filteredNavigation.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 px-3 py-2.5 rounded-lg',
                'text-neutral-400 hover:text-neutral-100 hover:bg-neutral-800',
                'transition-colors duration-200',
                isActive && 'bg-primary-500/10 text-primary-400 hover:bg-primary-500/20',
                isCollapsed && 'justify-center'
              )
            }
            title={isCollapsed ? item.label : undefined}
          >
            <item.icon size={20} className="flex-shrink-0" />
            {!isCollapsed && (
              <span className="font-medium truncate">{item.label}</span>
            )}
          </NavLink>
        ))}
      </nav>

      {/* Collapse toggle button */}
      <div className="p-3 border-t border-neutral-800">
        <button
          onClick={toggle}
          className={cn(
            'flex items-center gap-3 w-full px-3 py-2.5 rounded-lg',
            'text-neutral-500 hover:text-neutral-300 hover:bg-neutral-800',
            'transition-colors duration-200',
            isCollapsed && 'justify-center'
          )}
          title={isCollapsed ? 'Развернуть' : 'Свернуть'}
        >
          {isCollapsed ? (
            <ChevronRight size={20} />
          ) : (
            <>
              <ChevronLeft size={20} />
              <span className="text-sm">Свернуть</span>
            </>
          )}
        </button>
      </div>
    </aside>
  );
}

