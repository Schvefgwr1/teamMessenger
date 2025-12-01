import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/entities/session';
import { Spinner } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

/**
 * ProtectedRoute - защищённый маршрут
 *
 * Требует авторизации.
 * Если не авторизован - редирект на /login с сохранением целевого URL.
 */
export function ProtectedRoute() {
  const { isAuthenticated, isLoading } = useAuthStore();
  const location = useLocation();

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
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />;
  }

  return <Outlet />;
}

