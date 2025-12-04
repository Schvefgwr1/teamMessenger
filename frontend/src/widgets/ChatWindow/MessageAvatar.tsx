import { Avatar, Skeleton } from '@/shared/ui';
import { useUserById } from '@/entities/user';

interface MessageAvatarProps {
  senderId: string;
}

/**
 * Аватар отправителя сообщения с загрузкой данных пользователя
 */
export function MessageAvatar({ senderId }: MessageAvatarProps) {
  const { data: user, isLoading } = useUserById(senderId);

  if (isLoading) {
    return <Skeleton variant="circular" className="w-8 h-8" />;
  }

  return (
    <Avatar
      file={user?.avatar}
      fallback={user?.Username || senderId.slice(0, 2).toUpperCase()}
      size="sm"
    />
  );
}

