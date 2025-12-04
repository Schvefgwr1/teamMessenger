import { Link } from 'react-router-dom';
import { ArrowLeft, Settings, Users, Search } from 'lucide-react';
import { Avatar, Badge, Button } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import type { Chat } from '@/shared/types';

interface ChatHeaderProps {
  chat: Chat;
  onOpenSettings?: () => void;
  onOpenSearch?: () => void;
  showSettingsButton?: boolean;
}

/**
 * Заголовок чата
 */
export function ChatHeader({ 
  chat, 
  onOpenSettings, 
  onOpenSearch,
  showSettingsButton = true,
}: ChatHeaderProps) {
  const membersCount = chat.users?.length || 0;

  return (
    <div className="flex items-center gap-4 px-4 py-3 border-b border-neutral-800 bg-neutral-900/50">
      {/* Кнопка назад (мобильная) */}
      <Link
        to={ROUTES.CHATS}
        className="md:hidden p-2 -ml-2 rounded-lg text-neutral-400 hover:text-neutral-100 hover:bg-neutral-800 transition-colors"
      >
        <ArrowLeft size={20} />
      </Link>

      {/* Avatar */}
      <Avatar
        file={chat.avatarFile}
        fallback={chat.name}
        size="md"
      />

      {/* Info */}
      <div className="flex-1 min-w-0">
        <h2 className="font-semibold text-neutral-100 truncate">
          {chat.name}
        </h2>
        <div className="flex items-center gap-2">
          {chat.isGroup ? (
            <span className="flex items-center gap-1 text-xs text-neutral-400">
              <Users size={12} />
              {membersCount} участник{membersCount === 1 ? '' : membersCount < 5 ? 'а' : 'ов'}
            </span>
          ) : (
            <Badge variant="success" size="sm">Онлайн</Badge>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center gap-1">
        {/* Поиск (всегда доступен) */}
        {onOpenSearch && (
          <Button
            variant="ghost"
            size="icon"
            onClick={onOpenSearch}
            title="Поиск сообщений"
          >
            <Search size={20} />
          </Button>
        )}

        {/* Настройки (только если есть права) */}
        {showSettingsButton && onOpenSettings && (
          <Button
            variant="ghost"
            size="icon"
            onClick={onOpenSettings}
            title="Настройки чата"
          >
            <Settings size={20} />
          </Button>
        )}
      </div>
    </div>
  );
}

