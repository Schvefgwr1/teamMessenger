/**
 * ВАЖНО: Эти константы используются только для удобства работы с именами permissions.
 * Реальная проверка прав происходит:
 * 1. На бекенде при каждом API запросе (middleware)
 * 2. На фронтенде через user.role.permissions[] для UI (показать/скрыть элементы)
 *
 * Фронтенд НЕ является источником истины для permissions.
 * Если permission не найден у пользователя — элемент скрывается,
 * но даже если хакер покажет элемент, бекенд отклонит запрос.
 */

/**
 * Известные пользовательские permissions (для автодополнения и типизации)
 * Актуальный список получаем из user.role.permissions при авторизации
 */
export const KNOWN_PERMISSIONS = {
  // User management
  PROCESS_YOUR_ACC: 'process_your_acc',
  WATCH_USERS: 'watch_users',
  VIEW_FULL_USER_PROFILE: 'view_full_user_profile', // Только админы

  // Features
  PROCESS_CHATS: 'process_chats',
  PROCESS_TASKS: 'process_tasks',

  // Task statuses
  VIEW_TASK_STATUSES: 'view_task_statuses', // Пользовательские
  MANAGE_TASK_STATUSES: 'manage_task_statuses', // Только админы

  // Admin - system
  GET_PERMISSIONS: 'get_permissions',
  PROCESS_ROLES: 'process_roles',
  PROCESS_USERS_ROLES: 'process_users_roles',

  // Admin - chat system
  PROCESS_CHATS_ROLES: 'process_chats_roles',
  PROCESS_CHATS_PERMISSIONS: 'process_chats_permissions',

  // Legacy (deprecated, но может использоваться для обратной совместимости)
  PROCESS_TASKS_STATUSES: 'process_tasks_statuses', // @deprecated - используйте VIEW_TASK_STATUSES или MANAGE_TASK_STATUSES
} as const;

/**
 * Тип для известных permission names
 */
export type KnownPermission = (typeof KNOWN_PERMISSIONS)[keyof typeof KNOWN_PERMISSIONS];

/**
 * Проверка является ли строка известным permission
 * (для типобезопасности, не для авторизации)
 */
export function isKnownPermission(name: string): name is KnownPermission {
  return Object.values(KNOWN_PERMISSIONS).includes(name as KnownPermission);
}
