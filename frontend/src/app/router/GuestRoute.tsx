import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '@/entities/session';
import { Spinner } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';

/**
 * GuestRoute - маршрут только для гостей
 *
 * Если авторизован - редирект на главную.
 * Используется для страниц login/register.
 */
export function GuestRoute() {
  const { isAuthenticated, isLoading } = useAuthStore();

  // Показываем спиннер пока проверяем авторизацию
  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-neutral-950">
        <Spinner size="lg" />
      </div>
    );
  }

  // Уже авторизован - редирект на главную
  if (isAuthenticated) {
    return <Navigate to={ROUTES.HOME} replace />;
  }

  return <Outlet />;
}

