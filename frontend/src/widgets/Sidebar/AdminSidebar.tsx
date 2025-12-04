import { NavLink } from 'react-router-dom';
import { cn } from '@/shared/lib/cn';
import { useAuthStore } from '@/entities/session';
import { adminNavigation } from './adminNavigation';

/**
 * Sidebar для админ-панели
 * Всегда развёрнут, без сворачивания
 */
export function AdminSidebar() {
  const { hasPermission } = useAuthStore();

  // Фильтруем навигацию по permissions
  const filteredNavigation = adminNavigation.filter(
    (item) => !item.permission || hasPermission(item.permission)
  );

  return (
    <aside className="fixed left-0 top-16 bottom-0 w-64 z-30 bg-neutral-900 border-r border-neutral-800 flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-neutral-800">
        <h2 className="text-sm font-semibold text-neutral-400 uppercase tracking-wider">
          Администрирование
        </h2>
      </div>

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
                isActive && 'bg-primary-500/10 text-primary-400 hover:bg-primary-500/20'
              )
            }
          >
            <item.icon size={20} className="flex-shrink-0" />
            <span className="font-medium">{item.label}</span>
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="p-4 border-t border-neutral-800">
        <p className="text-xs text-neutral-600 text-center">
          Team Messenger Admin
        </p>
      </div>
    </aside>
  );
}

