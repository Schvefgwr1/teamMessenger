import { useEffect } from 'react';
import { useAuthStore } from '@/entities/session';
import { userApi, transformUserResponse } from '@/entities/user';

interface AuthProviderProps {
  children: React.ReactNode;
}

/**
 * AuthProvider - инициализирует сессию при старте приложения
 *
 * 1. Проверяет наличие токена в localStorage
 * 2. Если токен есть - загружает данные пользователя
 * 3. Если токен невалиден - очищает сессию
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const { token, setUser, setLoading, logout } = useAuthStore();

  useEffect(() => {
    const initAuth = async () => {
      // Если нет токена - сразу завершаем загрузку
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        // Загружаем данные пользователя
        const response = await userApi.getMe();
        const user = transformUserResponse(response.data);
        setUser(user);
      } catch (error: any) {
        // Проверяем тип ошибки
        // Если это 401 (Unauthorized) - токен невалидный, выходим
        // Если это CORS или другая ошибка сети - не выходим, просто логируем
        if (error?.response?.status === 401) {
          console.error('Token invalid, logging out:', error);
          logout();
        } else {
          // Другие ошибки (CORS, network, etc.) - не выходим, но логируем
          console.error('Failed to fetch user (non-auth error):', error);
          // Не вызываем logout() при CORS/network ошибках
          // Пользователь останется "аутентифицированным" по токену
        }
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, [token, setUser, setLoading, logout]);

  return <>{children}</>;
}

