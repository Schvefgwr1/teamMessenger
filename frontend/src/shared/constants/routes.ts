/**
 * Маршруты приложения
 */
export const ROUTES = {
  // Auth
  LOGIN: '/login',
  REGISTER: '/register',

  // Main
  HOME: '/',
  DASHBOARD: '/',

  // Profile
  PROFILE: '/profile',
  USER_PROFILE: (userId: string) => `/users/${userId}`,

  // Chats
  CHATS: '/chats',
  CHAT_DETAIL: (chatId: string) => `/chats/${chatId}`,

  // Tasks
  TASKS: '/tasks',
  TASK_DETAIL: (taskId: number | string) => `/tasks/${taskId}`,

  // Admin
  ADMIN: '/admin',
  ADMIN_ROLES: '/admin/roles',
  ADMIN_PERMISSIONS: '/admin/permissions',
  ADMIN_CHAT_ROLES: '/admin/chat-roles',
  ADMIN_CHAT_PERMISSIONS: '/admin/chat-permissions',
  ADMIN_TASK_STATUSES: '/admin/task-statuses',

  // Errors
  FORBIDDEN: '/403',
  NOT_FOUND: '/404',
} as const;

/**
 * Публичные маршруты (доступны без авторизации)
 */
export const PUBLIC_ROUTES = [ROUTES.LOGIN, ROUTES.REGISTER];

/**
 * Маршруты только для гостей (редирект если авторизован)
 */
export const GUEST_ONLY_ROUTES = [ROUTES.LOGIN, ROUTES.REGISTER];

/**
 * Административные маршруты
 */
export const ADMIN_ROUTES = [
  ROUTES.ADMIN,
  ROUTES.ADMIN_ROLES,
  ROUTES.ADMIN_PERMISSIONS,
  ROUTES.ADMIN_CHAT_ROLES,
  ROUTES.ADMIN_CHAT_PERMISSIONS,
  ROUTES.ADMIN_TASK_STATUSES,
];

