import { Link } from 'react-router-dom';
import { Users } from 'lucide-react';
import { cn } from '@/shared/lib/cn';
import { Avatar } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import { formatDate } from '@/shared/lib/formatDate';
import type { Chat } from '@/shared/types';

interface ChatListItemProps {
  chat: Chat;
  isActive?: boolean;
}

/**
 * Элемент списка чатов
 */
export function ChatListItem({ chat, isActive = false }: ChatListItemProps) {
  const membersCount = chat.users?.length || 0;

  return (
    <Link
      to={ROUTES.CHAT_DETAIL(chat.id)}
      className={cn(
        'flex items-center gap-3 p-3 rounded-xl transition-colors',
        'hover:bg-neutral-800/50',
        isActive && 'bg-primary-500/10 hover:bg-primary-500/20'
      )}
    >
      {/* Avatar */}
      <Avatar
        file={chat.avatarFile}
        fallback={chat.name}
        size="lg"
      />

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2">
          <h3 className={cn(
            'font-medium truncate',
            isActive ? 'text-primary-400' : 'text-neutral-100'
          )}>
            {chat.name}
          </h3>
          <span className="text-xs text-neutral-500 flex-shrink-0">
            {formatDate(chat.createdAt, 'short')}
          </span>
        </div>

        <div className="flex items-center gap-2 mt-1">
          {chat.isGroup && (
            <span className="flex items-center gap-1 text-xs text-neutral-500">
              <Users size={12} />
              {membersCount}
            </span>
          )}
          {chat.description && (
            <p className="text-sm text-neutral-400 truncate">
              {chat.description}
            </p>
          )}
        </div>
      </div>
    </Link>
  );
}

