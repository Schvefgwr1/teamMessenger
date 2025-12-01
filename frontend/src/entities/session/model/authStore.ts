import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import type { User, Permission, File as ApiFile } from '@/shared/types';

/**
 * Пользователь с аватаром для хранения в сессии
 */
export interface AuthUser extends User {
  avatar?: ApiFile;
}

interface AuthState {
  // State
  token: string | null;
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setToken: (token: string) => void;
  setUser: (user: AuthUser) => void;
  login: (token: string, user: AuthUser) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;

  // Permission helpers (работают с user.Role.Permissions)
  hasPermission: (permissionName: string) => boolean;
  hasAnyPermission: (permissionNames: string[]) => boolean;
  hasAllPermissions: (permissionNames: string[]) => boolean;
  isAdmin: () => boolean;
  getPermissions: () => Permission[];
}

/**
 * Извлечь Permissions из роли пользователя
 */
function getUserPermissions(user: AuthUser | null): Permission[] {
  return user?.Role?.Permissions ?? [];
}

export const useAuthStore = create<AuthState>()(
  persist(
    immer((set, get) => ({
      // Initial state
      token: null,
      user: null,
      isAuthenticated: false,
      isLoading: true,

      // Actions
      setToken: (token) =>
        set((state) => {
          state.token = token;
          state.isAuthenticated = !!token;
        }),

      setUser: (user) =>
        set((state) => {
          state.user = user;
          // Обновляем isAuthenticated при установке пользователя
          state.isAuthenticated = !!state.token;
        }),

      login: (token, user) =>
        set((state) => {
          state.token = token;
          state.user = user;
          state.isAuthenticated = true;
          state.isLoading = false;
        }),

      logout: () =>
        set((state) => {
          state.token = null;
          state.user = null;
          state.isAuthenticated = false;
          state.isLoading = false;
        }),

      setLoading: (loading) =>
        set((state) => {
          state.isLoading = loading;
        }),

      // Permission helpers
      hasPermission: (permissionName) => {
        const Permissions = getUserPermissions(get().user);
        return Permissions.some((p) => p.Name === permissionName);
      },

      hasAnyPermission: (permissionNames) => {
        const Permissions = getUserPermissions(get().user);
        return permissionNames.some((name) =>
          Permissions.some((p) => p.Name === name)
        );
      },

      hasAllPermissions: (permissionNames) => {
        const Permissions = getUserPermissions(get().user);
        return permissionNames.every((name) =>
          Permissions.some((p) => p.Name === name)
        );
      },

      isAdmin: () => {
        const { user } = get();
        // Проверяем по имени роли или ID (admin Role обычно id=1)
        return user?.Role?.Name === 'admin' || user?.Role?.ID === 1;
      },

      getPermissions: () => getUserPermissions(get().user),
    })),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      // Сохраняем только токен, user загружаем при старте
      partialize: (state) => ({
        token: state.token,
      }),
      // При восстановлении из localStorage устанавливаем isAuthenticated на основе токена
      onRehydrateStorage: () => (state) => {
        if (state) {
          // Если есть токен, считаем что пользователь аутентифицирован (пока не загрузим данные)
          state.isAuthenticated = !!state.token;
        }
      },
    }
  )
);

// Селекторы для оптимизации ререндеров
export const useToken = () => useAuthStore((s) => s.token);
export const useUser = () => useAuthStore((s) => s.user);
export const useIsAuthenticated = () => useAuthStore((s) => s.isAuthenticated);
export const useAuthLoading = () => useAuthStore((s) => s.isLoading);

