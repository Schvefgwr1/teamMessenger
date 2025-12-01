import {
  LayoutDashboard,
  Shield,
  Key,
  Users,
  MessageSquare,
  ListChecks,
  type LucideIcon,
} from 'lucide-react';
import { ROUTES, PERMISSIONS } from '@/shared/constants';

export interface AdminNavItem {
  label: string;
  path: string;
  icon: LucideIcon;
  permission?: string;
}

/**
 * Навигация для админ-панели
 */
export const adminNavigation: AdminNavItem[] = [
  {
    label: 'Dashboard',
    path: ROUTES.ADMIN,
    icon: LayoutDashboard,
  },
  {
    label: 'Роли пользователей',
    path: ROUTES.ADMIN_ROLES,
    icon: Shield,
    permission: PERMISSIONS.PROCESS_ROLES,
  },
  {
    label: 'Разрешения',
    path: ROUTES.ADMIN_PERMISSIONS,
    icon: Key,
    permission: PERMISSIONS.GET_PERMISSIONS,
  },
  {
    label: 'Роли чатов',
    path: ROUTES.ADMIN_CHAT_ROLES,
    icon: Users,
    permission: PERMISSIONS.PROCESS_CHATS_ROLES,
  },
  {
    label: 'Разрешения чатов',
    path: ROUTES.ADMIN_CHAT_PERMISSIONS,
    icon: MessageSquare,
    permission: PERMISSIONS.PROCESS_CHATS_PERMISSIONS,
  },
  {
    label: 'Статусы задач',
    path: ROUTES.ADMIN_TASK_STATUSES,
    icon: ListChecks,
    permission: PERMISSIONS.PROCESS_TASKS_STATUSES,
  },
];

