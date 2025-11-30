import { cn } from '@/shared/lib';

interface SkeletonProps {
  className?: string;
  variant?: 'text' | 'circular' | 'rectangular';
  width?: string | number;
  height?: string | number;
}

export function Skeleton({ className, variant = 'text', width, height }: SkeletonProps) {
  return (
    <div
      className={cn(
        'animate-pulse bg-neutral-800',
        {
          'h-4 rounded': variant === 'text',
          'rounded-full': variant === 'circular',
          'rounded-lg': variant === 'rectangular',
        },
        className
      )}
      style={{ width, height }}
    />
  );
}

// Skeleton для карточки
Skeleton.Card = function SkeletonCard({ className }: { className?: string }) {
  return (
    <div className={cn('p-4 rounded-xl bg-neutral-900 space-y-3', className)}>
      <Skeleton className="h-4 w-3/4" />
      <Skeleton className="h-4 w-1/2" />
      <Skeleton className="h-20 w-full" variant="rectangular" />
    </div>
  );
};

// Skeleton для сообщения в чате
Skeleton.ChatMessage = function SkeletonChatMessage({ className }: { className?: string }) {
  return (
    <div className={cn('flex gap-3 p-3', className)}>
      <Skeleton variant="circular" width={40} height={40} />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
      </div>
    </div>
  );
};

// Skeleton для строки таблицы
Skeleton.TableRow = function SkeletonTableRow({ columns = 4 }: { columns?: number }) {
  return (
    <div className="flex gap-4 p-3">
      {Array.from({ length: columns }).map((_, i) => (
        <Skeleton key={i} className="h-4 flex-1" />
      ))}
    </div>
  );
};

