import {
  Home,
  MessageSquare,
  CheckSquare,
  User,
  type LucideIcon,
} from 'lucide-react';
import { ROUTES, PERMISSIONS } from '@/shared/constants';

export interface NavItem {
  label: string;
  path: string;
  icon: LucideIcon;
  permission?: string;
}

/**
 * Навигация для основного приложения
 */
export const userNavigation: NavItem[] = [
  {
    label: 'Главная',
    path: ROUTES.HOME,
    icon: Home,
  },
  {
    label: 'Чаты',
    path: ROUTES.CHATS,
    icon: MessageSquare,
    permission: PERMISSIONS.PROCESS_CHATS,
  },
  {
    label: 'Задачи',
    path: ROUTES.TASKS,
    icon: CheckSquare,
    permission: PERMISSIONS.PROCESS_TASKS,
  },
  {
    label: 'Профиль',
    path: ROUTES.PROFILE,
    icon: User,
    permission: PERMISSIONS.PROCESS_YOUR_ACC,
  },
];

