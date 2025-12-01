import React, { useMemo } from 'react';
import type { Permission } from '@/shared/types';

/**
 * Хелпер для проверки наличия permission у пользователя
 *
 * ВАЖНО: Это проверка только для UI (показать/скрыть элементы).
 * Реальная авторизация происходит на бекенде при каждом API запросе.
 *
 * @param userPermissions - массив permissions из user.role.permissions
 * @param permissionName - имя permission для проверки
 */
export function hasPermission(
  userPermissions: Permission[] | undefined,
  permissionName: string
): boolean {
  if (!userPermissions || userPermissions.length === 0) return false;
  return userPermissions.some((p) => p.Name === permissionName);
}

/**
 * Проверка наличия хотя бы одного из permissions
 */
export function hasAnyPermission(
  userPermissions: Permission[] | undefined,
  permissionNames: string[]
): boolean {
  if (!userPermissions || userPermissions.length === 0) return false;
  return permissionNames.some((name) => hasPermission(userPermissions, name));
}

/**
 * Проверка наличия всех permissions
 */
export function hasAllPermissions(
  userPermissions: Permission[] | undefined,
  permissionNames: string[]
): boolean {
  if (!userPermissions || userPermissions.length === 0) return false;
  return permissionNames.every((name) => hasPermission(userPermissions, name));
}

/**
 * Хук для получения списка имён permissions пользователя
 */
export function usePermissionNames(permissions: Permission[] | undefined): string[] {
  return useMemo(() => {
    if (!permissions) return [];
    return permissions.map((p) => p.Name);
  }, [permissions]);
}

/**
 * Компонент для условного рендеринга на основе permission
 *
 * Использование:
 * ```tsx
 * <Can permissions={user.role.permissions} check="process_chats">
 *   <Button>Создать чат</Button>
 * </Can>
 * ```
 */
interface CanProps {
  permissions: Permission[] | undefined;
  check: string | string[];
  mode?: 'any' | 'all';
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function Can({
  permissions,
  check,
  mode = 'any',
  children,
  fallback = null,
}: CanProps) {
  const permissionNames = Array.isArray(check) ? check : [check];

  const allowed =
    mode === 'all'
      ? hasAllPermissions(permissions, permissionNames)
      : hasAnyPermission(permissions, permissionNames);

  return allowed ? <>{children}</> : <>{fallback}</>;
}