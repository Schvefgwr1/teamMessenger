import { Outlet, Link } from 'react-router-dom';
import { ArrowLeft, Shield } from 'lucide-react';
import { AdminSidebar } from '@/widgets/Sidebar/AdminSidebar';
import { ROUTES } from '@/shared/constants';

/**
 * Layout для админ-панели
 * Отдельный header с возвратом в приложение, AdminSidebar слева
 */
export function AdminLayout() {
  return (
    <div className="min-h-screen bg-neutral-950 text-neutral-100">
      {/* Admin Header */}
      <header className="fixed top-0 left-0 right-0 h-16 bg-neutral-900 border-b border-neutral-800 z-40 flex items-center px-6">
        <Link
          to={ROUTES.HOME}
          className="flex items-center gap-2 text-neutral-400 hover:text-neutral-100 transition-colors"
        >
          <ArrowLeft size={20} />
          <span>Вернуться в приложение</span>
        </Link>

        <div className="ml-auto flex items-center gap-3">
          <Shield size={18} className="text-primary-400" />
          <span className="text-sm text-primary-400 font-medium">
            Панель администратора
          </span>
        </div>
      </header>

      <div className="flex pt-16">
        {/* Admin Sidebar */}
        <AdminSidebar />

        {/* Main content area */}
        <main className="flex-1 ml-64 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

