import { ReactNode } from 'react';
import { LucideIcon } from 'lucide-react';
import { Button, type ButtonProps } from '../Button';
import { cn } from '@/shared/lib/cn';

export interface EmptyStateProps {
  /** Иконка для отображения */
  icon?: LucideIcon;
  /** Заголовок */
  title?: string;
  /** Описание */
  description?: string;
  /** Дополнительный контент */
  children?: ReactNode;
  /** Действие (кнопка) */
  action?: {
    label: string;
    onClick: () => void;
    variant?: ButtonProps['variant'];
    icon?: ReactNode;
  };
  /** Размер */
  size?: 'sm' | 'md' | 'lg';
  /** Кастомный класс */
  className?: string;
}

/**
 * Универсальный компонент для отображения пустого состояния
 * 
 * @example
 * ```tsx
 * <EmptyState
 *   icon={MessageSquare}
 *   title="Нет чатов"
 *   description="Создайте первый чат, чтобы начать общение"
 *   action={{
 *     label: "Создать чат",
 *     onClick: () => setCreateModalOpen(true),
 *     icon: <Plus size={18} />
 *   }}
 * />
 * ```
 */
export function EmptyState({
  icon: Icon,
  title,
  description,
  children,
  action,
  size = 'md',
  className,
}: EmptyStateProps) {
  const sizeClasses = {
    sm: {
      icon: 'w-12 h-12',
      iconSize: 24,
      title: 'text-base',
      description: 'text-sm',
    },
    md: {
      icon: 'w-16 h-16',
      iconSize: 32,
      title: 'text-lg',
      description: 'text-sm',
    },
    lg: {
      icon: 'w-20 h-20',
      iconSize: 40,
      title: 'text-xl',
      description: 'text-base',
    },
  };

  const sizes = sizeClasses[size];

  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center text-center py-12 px-4',
        className
      )}
    >
      {/* Иконка */}
      {Icon && (
        <div
          className={cn(
            'rounded-full bg-neutral-800/50 flex items-center justify-center mb-4',
            sizes.icon
          )}
        >
          <Icon className="text-neutral-400" size={sizes.iconSize} />
        </div>
      )}

      {/* Заголовок */}
      {title && (
        <h3 className={cn('font-semibold text-neutral-200 mb-2', sizes.title)}>
          {title}
        </h3>
      )}

      {/* Описание */}
      {description && (
        <p className={cn('text-neutral-400 max-w-md mb-6', sizes.description)}>
          {description}
        </p>
      )}

      {/* Дополнительный контент */}
      {children}

      {/* Действие */}
      {action && (
        <Button
          variant={action.variant || 'primary'}
          onClick={action.onClick}
          leftIcon={action.icon}
          size={size === 'sm' ? 'sm' : 'md'}
        >
          {action.label}
        </Button>
      )}
    </div>
  );
}

