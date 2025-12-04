import { Outlet } from 'react-router-dom';
import { Header } from '@/widgets/Header';
import { Sidebar } from '@/widgets/Sidebar';
import { cn } from '@/shared/lib/cn';
import { useSidebarStore } from '@/widgets/Sidebar/model/sidebarStore';

/**
 * Основной layout для авторизованных пользователей
 * Header сверху, Sidebar слева, контент справа
 */
export function MainLayout() {
  const { isCollapsed } = useSidebarStore();

  return (
    <div className="h-screen bg-neutral-950 text-neutral-100 flex flex-col overflow-hidden">
      {/* Header - фиксированный сверху */}
      <Header />

      {/* Content area - занимает оставшееся пространство */}
      <div className="flex flex-1 overflow-hidden pt-16">
        {/* Sidebar - фиксированный слева */}
        <Sidebar />

        {/* Main content area */}
        <main
          className={cn(
            'flex-1 p-6 transition-all duration-300 ease-in-out',
            'overflow-y-auto overflow-x-hidden',
            isCollapsed ? 'ml-16' : 'ml-64'
          )}
        >
          <Outlet />
        </main>
      </div>
    </div>
  );
}

