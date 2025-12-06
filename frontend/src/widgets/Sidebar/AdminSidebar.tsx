import { NavLink } from 'react-router-dom';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { cn } from '@/shared/lib/cn';
import { useAuthStore } from '@/entities/session';
import { useSidebarStore } from './model/sidebarStore';
import { adminNavigation } from './adminNavigation';

/**
 * Sidebar для админ-панели
 * С навигацией и возможностью сворачивания
 */
export function AdminSidebar() {
  const { hasPermission } = useAuthStore();
  const { isCollapsed, toggle } = useSidebarStore();

  // Фильтруем навигацию по permissions
  const filteredNavigation = adminNavigation.filter(
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
      {/* Header */}
      {!isCollapsed && (
        <div className="p-4 border-b border-neutral-800">
          <h2 className="text-sm font-semibold text-neutral-400 uppercase tracking-wider">
            Администрирование
          </h2>
        </div>
      )}

      {/* Navigation */}
      <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
        {filteredNavigation.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === '/admin'}
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

