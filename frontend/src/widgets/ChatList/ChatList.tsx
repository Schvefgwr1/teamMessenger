import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Search, MessageSquarePlus } from 'lucide-react';
import { Input, Button, Skeleton } from '@/shared/ui';
import { useUserChats } from '@/entities/chat';
import { ChatListItem } from './ChatListItem';
import { useDebounce } from '@/shared/hooks';

interface ChatListProps {
  onCreateChat?: () => void;
}

/**
 * Список чатов пользователя
 */
export function ChatList({ onCreateChat }: ChatListProps) {
  const { chatId } = useParams<{ chatId: string }>();
  const { data: chats, isLoading, error } = useUserChats();

  const [searchQuery, setSearchQuery] = useState('');
  const debouncedSearch = useDebounce(searchQuery, 300);

  // Фильтрация чатов по поиску
  const filteredChats = chats?.filter((chat) =>
    chat.name.toLowerCase().includes(debouncedSearch.toLowerCase())
  ) || [];

  if (isLoading) {
    return <ChatListSkeleton />;
  }

  if (error) {
    return (
      <div className="p-4 text-center text-error">
        Ошибка загрузки чатов
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-4 border-b border-neutral-800">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-neutral-100">Чаты</h2>
          {onCreateChat && (
            <Button
              size="sm"
              variant="ghost"
              onClick={onCreateChat}
              leftIcon={<MessageSquarePlus size={18} />}
            >
              Новый
            </Button>
          )}
        </div>
        <Input
          placeholder="Поиск чатов..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          leftIcon={<Search size={18} />}
        />
      </div>

      {/* Chat list */}
      <div className="flex-1 overflow-y-auto p-2">
        {filteredChats.length === 0 ? (
          <EmptyState
            hasSearch={debouncedSearch.length > 0}
            onCreateChat={onCreateChat}
          />
        ) : (
          <div className="space-y-1">
            {filteredChats.map((chat) => (
              <ChatListItem
                key={chat.id}
                chat={chat}
                isActive={chat.id === chatId}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

// Empty state component
interface EmptyStateProps {
  hasSearch: boolean;
  onCreateChat?: () => void;
}

function EmptyState({ hasSearch, onCreateChat }: EmptyStateProps) {
  if (hasSearch) {
    return (
      <div className="flex flex-col items-center justify-center h-full py-12 text-center">
        <Search size={48} className="text-neutral-700 mb-4" />
        <p className="text-neutral-400">Чаты не найдены</p>
        <p className="text-sm text-neutral-500 mt-1">
          Попробуйте изменить запрос
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center h-full py-12 text-center">
      <MessageSquarePlus size={48} className="text-neutral-700 mb-4" />
      <p className="text-neutral-400">У вас пока нет чатов</p>
      {onCreateChat && (
        <Button
          variant="primary"
          size="sm"
          className="mt-4"
          onClick={onCreateChat}
        >
          Создать первый чат
        </Button>
      )}
    </div>
  );
}

// Skeleton loader
function ChatListSkeleton() {
  return (
    <div className="p-4 space-y-4">
      <div className="flex items-center justify-between">
        <Skeleton className="h-6 w-20" />
        <Skeleton className="h-8 w-20" />
      </div>
      <Skeleton className="h-10 w-full" />
      <div className="space-y-2 mt-4">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="flex items-center gap-3 p-3">
            <Skeleton variant="circular" className="w-12 h-12" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-3 w-1/2" />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

