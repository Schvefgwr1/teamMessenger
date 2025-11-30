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

  // Features
  PROCESS_CHATS: 'process_chats',
  PROCESS_TASKS: 'process_tasks',

  // Admin - system
  GET_PERMISSIONS: 'get_permissions',
  PROCESS_ROLES: 'process_roles',

  // Admin - chat system
  PROCESS_CHATS_ROLES: 'process_chats_roles',
  PROCESS_CHATS_PERMISSIONS: 'process_chats_permissions',

  // Admin - task system
  PROCESS_TASKS_STATUSES: 'process_tasks_statuses',
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
