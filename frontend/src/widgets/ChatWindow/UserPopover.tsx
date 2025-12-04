import { Popover, Avatar, Skeleton } from '@/shared/ui';
import { useUserBrief } from '@/entities/user';
import { User, Shield } from 'lucide-react';

interface UserPopoverProps {
  userId: string | null | undefined;
  chatId: string;
  children: React.ReactNode;
}

/**
 * Popover с краткой информацией о пользователе и его ролью в чате
 */
export function UserPopover({ userId, chatId, children }: UserPopoverProps) {
  const { data: user, isLoading } = useUserBrief(userId || undefined, chatId);

  if (!userId) {
    return <>{children}</>;
  }

  return (
    <Popover>
      <Popover.Trigger asChild>{children}</Popover.Trigger>
      <Popover.Content side="right" align="start" className="w-72">
        {isLoading ? (
          <UserPopoverSkeleton />
        ) : user ? (
          <div className="space-y-4">
            {/* Header с аватаром и именем */}
            <div className="flex items-center gap-3">
              <Avatar
                file={user.avatarFile}
                fallback={user.username}
                size="lg"
              />
              <div className="flex-1 min-w-0">
                <h3 className="font-semibold text-neutral-100 truncate">
                  {user.username}
                </h3>
                {user.email && (
                  <p className="text-sm text-neutral-400 truncate">
                    {user.email}
                  </p>
                )}
              </div>
            </div>

            {/* Разделитель */}
            <div className="h-px bg-neutral-800" />

            {/* Информация */}
            <div className="space-y-2">
              {user.age && (
                <div className="flex items-center gap-2 text-sm">
                  <User size={14} className="text-neutral-500 flex-shrink-0" />
                  <span className="text-neutral-300">{user.age} лет</span>
                </div>
              )}

              {user.description && (
                <p className="text-sm text-neutral-400 line-clamp-2">
                  {user.description}
                </p>
              )}

              {/* Роль в чате */}
              {user.chatRoleName && (
                <div className="pt-2">
                  <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium bg-primary-500/20 text-primary-400">
                    <Shield size={12} />
                    {user.chatRoleName}
                  </span>
                </div>
              )}
            </div>
          </div>
        ) : (
          <p className="text-sm text-neutral-400">Пользователь не найден</p>
        )}
      </Popover.Content>
    </Popover>
  );
}

// Skeleton для загрузки
function UserPopoverSkeleton() {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <Skeleton variant="circular" className="w-12 h-12" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-3 w-32" />
        </div>
      </div>
      <div className="h-px bg-neutral-800" />
      <div className="space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-4 w-1/2" />
      </div>
    </div>
  );
}

