import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '@/entities/session';
import { Spinner } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

/**
 * AdminRoute - маршрут для администраторов
 *
 * Требует:
 * 1. Авторизации
 * 2. Роли admin или наличия permission 'process_roles'
 *
 * Если не авторизован - редирект на /login
 * Если нет прав - редирект на /403
 */
export function AdminRoute() {
  const { isAuthenticated, isLoading, isAdmin, hasPermission } = useAuthStore();

  // Показываем спиннер пока проверяем авторизацию
  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-neutral-950">
        <Spinner size="lg" />
      </div>
    );
  }

  // Не авторизован - редирект на логин
  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />;
  }

  // Проверяем права администратора
  const isAdminUser = isAdmin() || hasPermission('process_roles');

  if (!isAdminUser) {
    return <Navigate to={ROUTES.FORBIDDEN} replace />;
  }

  return <Outlet />;
}

