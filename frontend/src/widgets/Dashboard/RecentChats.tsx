import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import { MessageSquare, Users, ChevronRight, Clock } from 'lucide-react';
import { useQueries } from '@tanstack/react-query';
import { Card, Avatar, Skeleton } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import { useUserChats, chatApi } from '@/entities/chat';
import { formatMessageTime } from '@/shared/lib';
import { parseISO, subDays, isAfter } from 'date-fns';
import type { Chat, Message } from '@/shared/types';

const CHATS_LIMIT = 5;
const DAYS_AGO = 1; // Последние сутки

/**
 * Виджет последних чатов для Dashboard
 * Показывает чаты с сообщениями за последние сутки
 */
export function RecentChats() {
  const { data: chats = [], isLoading: chatsLoading, error } = useUserChats();

  // Получаем последнее сообщение для каждого чата параллельно
  const messagesQueries = useQueries({
    queries: chats.map((chat) => ({
      queryKey: ['chat-last-message', chat.id],
      queryFn: async () => {
        try {
          const response = await chatApi.getMessages(chat.id, { limit: 1 });
          const messages = response.data || [];
          return messages[0] || null;
        } catch {
          return null;
        }
      },
      enabled: !!chat.id && !chatsLoading,
      staleTime: 30 * 1000,
    })),
  });

  // Фильтруем чаты с сообщениями за последние сутки
  const recentChats = useMemo(() => {
    if (chatsLoading || messagesQueries.some((q) => q.isLoading)) {
      return [];
    }

    const oneDayAgo = subDays(new Date(), DAYS_AGO);
    const chatsWithRecentMessages: Array<{ chat: Chat; lastMessage: Message | null }> = [];

    chats.forEach((chat, index) => {
      const lastMessage = messagesQueries[index]?.data as Message | null | undefined;
      
      if (lastMessage) {
        try {
          const messageDate = parseISO(lastMessage.createdAt);
          if (isAfter(messageDate, oneDayAgo)) {
            chatsWithRecentMessages.push({ chat, lastMessage });
          }
        } catch {
          // Игнорируем ошибки парсинга даты
        }
      }
    });

    // Сортируем по времени последнего сообщения (новые первыми)
    return chatsWithRecentMessages
      .sort((a, b) => {
        if (!a.lastMessage || !b.lastMessage) return 0;
        return (
          new Date(b.lastMessage.createdAt).getTime() -
          new Date(a.lastMessage.createdAt).getTime()
        );
      })
      .slice(0, CHATS_LIMIT)
      .map((item) => ({ chat: item.chat, lastMessage: item.lastMessage }));
  }, [chats, messagesQueries, chatsLoading]);

  const isLoading = chatsLoading || messagesQueries.some((q) => q.isLoading);

  if (error) {
    return (
      <Card>
        <Card.Header>
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-neutral-200 flex items-center gap-2">
              <MessageSquare size={20} className="text-primary-400" />
              Последние чаты
            </h3>
          </div>
        </Card.Header>
        <Card.Body>
          <p className="text-sm text-neutral-500 text-center py-4">
            Ошибка загрузки чатов
          </p>
        </Card.Body>
      </Card>
    );
  }

  return (
    <Card>
      <Card.Header>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-neutral-200 flex items-center gap-2">
            <MessageSquare size={20} className="text-primary-400" />
            Последние чаты
          </h3>
          <Link
            to={ROUTES.CHATS}
            className="text-sm text-primary-400 hover:text-primary-300 transition-colors flex items-center gap-1"
          >
            Все чаты
            <ChevronRight size={16} />
          </Link>
        </div>
      </Card.Header>
      <Card.Body>
        {isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton variant="circular" width={48} height={48} />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-3 w-1/2" />
                </div>
              </div>
            ))}
          </div>
        ) : recentChats.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-2">
            {recentChats.map(({ chat, lastMessage }) => (
              <ChatItem key={chat.id} chat={chat} lastMessage={lastMessage} />
            ))}
          </div>
        )}
      </Card.Body>
    </Card>
  );
}

interface ChatItemProps {
  chat: Chat;
  lastMessage: Message | null;
}

function ChatItem({ chat, lastMessage }: ChatItemProps) {
  const membersCount = chat.users?.length || 0;
  const lastMessageTime = lastMessage ? formatMessageTime(lastMessage.createdAt) : null;

  return (
    <Link
      to={ROUTES.CHAT_DETAIL(chat.id)}
      className="flex items-center gap-3 p-3 rounded-lg hover:bg-neutral-800/50 transition-colors group"
    >
      <Avatar
        file={chat.avatarFile}
        fallback={chat.name}
        size="md"
      />
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2">
          <h4 className="font-medium text-neutral-100 truncate group-hover:text-primary-400 transition-colors">
            {chat.name}
          </h4>
          {lastMessageTime && (
            <span className="flex items-center gap-1 text-xs text-neutral-500 flex-shrink-0">
              <Clock size={12} />
              {lastMessageTime}
            </span>
          )}
        </div>
        <div className="flex items-center gap-2 mt-1">
          {chat.isGroup && (
            <span className="flex items-center gap-1 text-xs text-neutral-500">
              <Users size={12} />
              {membersCount}
            </span>
          )}
          {lastMessage && (
            <p className="text-xs text-neutral-400 truncate">
              {lastMessage.content}
            </p>
          )}
          {!lastMessage && chat.description && (
            <p className="text-xs text-neutral-400 truncate">
              {chat.description}
            </p>
          )}
        </div>
      </div>
    </Link>
  );
}

function EmptyState() {
  return (
    <div className="text-center py-8">
      <MessageSquare size={48} className="mx-auto text-neutral-600 mb-3" />
      <p className="text-sm text-neutral-500">
        У вас пока нет чатов
      </p>
      <Link
        to={ROUTES.CHATS}
        className="inline-block mt-3 text-sm text-primary-400 hover:text-primary-300 transition-colors"
      >
        Создать чат
      </Link>
    </div>
  );
}

