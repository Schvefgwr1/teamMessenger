import { useEffect, useRef, useImperativeHandle, forwardRef, useCallback } from 'react';
import { cn } from '@/shared/lib/cn';
import { Skeleton } from '@/shared/ui';
import { formatMessageTime, formatChatDate } from '@/shared/lib/formatDate';
import { useAuthStore } from '@/entities/session';
import { UserPopover } from './UserPopover';
import { MessageAvatar } from './MessageAvatar';
import type { Message } from '@/shared/types';

export interface MessageListHandle {
  scrollToMessage: (messageId: string) => void;
}

interface MessageListProps {
  messages: Message[];
  chatId: string;
  isLoading?: boolean;
  highlightedMessageId?: string | null;
}

/**
 * Список сообщений чата
 */
export const MessageList = forwardRef<MessageListHandle, MessageListProps>(
  function MessageList({ messages, chatId, isLoading, highlightedMessageId }, ref) {
  const { user } = useAuthStore();
  const bottomRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const messageRefs = useRef<Map<string, HTMLDivElement>>(new Map());

  // Функция для скролла к сообщению
  const scrollToMessage = useCallback((messageId: string) => {
    const messageElement = messageRefs.current.get(messageId);
    if (messageElement) {
      messageElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
      // Подсветка сообщения временно
      messageElement.classList.add('ring-2', 'ring-primary-500', 'ring-offset-2', 'ring-offset-neutral-950');
      setTimeout(() => {
        messageElement.classList.remove('ring-2', 'ring-primary-500', 'ring-offset-2', 'ring-offset-neutral-950');
      }, 2000);
    }
  }, []);

  // Expose scrollToMessage через ref
  useImperativeHandle(ref, () => ({
    scrollToMessage,
  }), [scrollToMessage]);

  // Прокрутка к последнему сообщению
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages.length]);

  // Подсветка сообщения при изменении highlightedMessageId
  useEffect(() => {
    if (highlightedMessageId) {
      scrollToMessage(highlightedMessageId);
    }
  }, [highlightedMessageId, scrollToMessage]);

  if (isLoading) {
    return <MessageListSkeleton />;
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <p className="text-neutral-500">Сообщений пока нет. Начните диалог!</p>
      </div>
    );
  }

  // Переворачиваем массив - старые сверху, новые снизу
  const sortedMessages = [...messages].reverse();

  return (
    <div ref={containerRef} className="flex-1 overflow-y-auto p-4 space-y-3">
      {sortedMessages.map((message, index) => {
        const prevMessage = index > 0 ? sortedMessages[index - 1] : null;
        
        // Сравниваем даты для разделителя (только дата, без времени)
        const getDateKey = (dateStr: string) => {
          try {
            const date = new Date(dateStr);
            return `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
          } catch {
            return dateStr;
          }
        };
        
        const currentDateKey = getDateKey(message.createdAt);
        const prevDateKey = prevMessage ? getDateKey(prevMessage.createdAt) : null;
        const showDateSeparator = prevDateKey !== currentDateKey;
        const dateLabel = formatChatDate(message.createdAt);

        return (
          <div 
            key={message.id}
            ref={(el) => {
              if (el) messageRefs.current.set(message.id, el);
            }}
          >
            {/* Разделитель даты только при смене даты */}
            {showDateSeparator && (
              <div className="flex items-center gap-4 my-4">
                <div className="flex-1 h-px bg-neutral-800" />
                <span className="text-xs text-neutral-500 px-2 font-medium">
                  {dateLabel}
                </span>
                <div className="flex-1 h-px bg-neutral-800" />
              </div>
            )}

            {/* Сообщение */}
            <MessageItem
              message={message}
              chatId={chatId}
              isOwn={message.senderID === user?.ID}
            />
          </div>
        );
      })}
      <div ref={bottomRef} />
    </div>
  );
});

// Компонент одного сообщения
interface MessageItemProps {
  message: Message;
  chatId: string;
  isOwn: boolean;
}

function MessageItem({ message, chatId, isOwn }: MessageItemProps) {
  return (
    <div
      className={cn(
        'flex gap-3 max-w-[80%]',
        isOwn && 'ml-auto flex-row-reverse'
      )}
    >
      {/* Avatar (только для чужих сообщений) с Popover */}
      {!isOwn && message.senderID && (
        <UserPopover userId={message.senderID} chatId={chatId}>
          <button
            type="button"
            className="flex-shrink-0 cursor-pointer hover:opacity-80 transition-opacity"
          >
            <MessageAvatar senderId={message.senderID} />
          </button>
        </UserPopover>
      )}

      {/* Контент сообщения */}
      <div
        className={cn(
          'rounded-2xl px-4 py-2',
          isOwn
            ? 'bg-primary-500 text-white rounded-br-md'
            : 'bg-neutral-800 text-neutral-100 rounded-bl-md'
        )}
      >
        <p className="text-sm whitespace-pre-wrap break-words">
          {message.content}
        </p>

        {/* Прикрепленные файлы */}
        {message.files && message.files.length > 0 && (
          <div className="mt-2 space-y-1">
            {message.files.map((file) => (
              <a
                key={file.id}
                href={file.url}
                target="_blank"
                rel="noopener noreferrer"
                className={cn(
                  'text-xs px-2 py-1 rounded block hover:opacity-80',
                  isOwn ? 'bg-primary-600' : 'bg-neutral-700'
                )}
              >
                📎 {file.name}
              </a>
            ))}
          </div>
        )}

        {/* Время */}
        <p
          className={cn(
            'text-xs mt-1',
            isOwn ? 'text-primary-200' : 'text-neutral-500'
          )}
        >
          {formatMessageTime(message.createdAt)}
        </p>
      </div>
    </div>
  );
}

// Skeleton loader
function MessageListSkeleton() {
  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {Array.from({ length: 6 }).map((_, i) => (
        <div
          key={i}
          className={cn('flex gap-3', i % 2 === 0 ? '' : 'flex-row-reverse')}
        >
          {i % 2 === 0 && (
            <Skeleton variant="circular" className="w-8 h-8" />
          )}
          <div className="space-y-1">
            <Skeleton
              className={cn(
                'h-16 rounded-2xl',
                i % 2 === 0 ? 'w-64' : 'w-48'
              )}
            />
          </div>
        </div>
      ))}
    </div>
  );
}

