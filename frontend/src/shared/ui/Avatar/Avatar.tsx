import { useState } from 'react';
import { cn, getAvatarUrl } from '@/shared/lib';
import type { File as ApiFile } from '@/shared/types';

interface AvatarProps {
  /** Объект File из API */
  file?: ApiFile | null;
  /** Прямой URL (альтернатива file) */
  src?: string | null;
  /** Alt текст */
  alt?: string;
  /** Текст для инициалов (fallback) */
  fallback?: string;
  /** Размер аватара */
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  /** Индикатор статуса */
  status?: 'online' | 'offline' | 'away' | 'busy';
  /** Дополнительные классы */
  className?: string;
}

const sizeClasses = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-base',
  lg: 'w-12 h-12 text-lg',
  xl: 'w-16 h-16 text-xl',
};

const statusClasses = {
  online: 'bg-success',
  offline: 'bg-neutral-500',
  away: 'bg-warning',
  busy: 'bg-error',
};

const statusSizeClasses = {
  xs: 'w-2 h-2',
  sm: 'w-2.5 h-2.5',
  md: 'w-3 h-3',
  lg: 'w-3.5 h-3.5',
  xl: 'w-4 h-4',
};

export function Avatar({
  file,
  src,
  alt,
  fallback,
  size = 'md',
  status,
  className,
}: AvatarProps) {
  const [imageError, setImageError] = useState(false);

  // Приоритет: file.url > src > null
  const avatarUrl = file ? getAvatarUrl(file) : src;
  const initials = fallback?.slice(0, 2).toUpperCase() || '?';

  return (
    <div className={cn('relative inline-block', className)}>
      {avatarUrl && !imageError ? (
        <img
          src={avatarUrl}
          alt={alt || fallback || 'Avatar'}
          onError={() => setImageError(true)}
          className={cn('rounded-full object-cover bg-neutral-800', sizeClasses[size])}
        />
      ) : (
        <div
          className={cn(
            'rounded-full flex items-center justify-center',
            'bg-gradient-to-br from-primary-500 to-primary-700',
            'text-white font-medium',
            sizeClasses[size]
          )}
        >
          {initials}
        </div>
      )}
      {status && (
        <span
          className={cn(
            'absolute bottom-0 right-0 block rounded-full ring-2 ring-neutral-900',
            statusSizeClasses[size],
            statusClasses[status]
          )}
        />
      )}
    </div>
  );
}

