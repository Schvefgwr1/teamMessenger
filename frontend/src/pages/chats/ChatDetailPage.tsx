import { useState, useRef, useCallback } from 'react';
import { useParams, Navigate } from 'react-router-dom';
import { 
  useChatMessages, 
  useSendMessage, 
  useUserChats,
  useMyRoleInChat,
  hasPermission,
} from '@/entities/chat';
import { 
  MessageList, 
  MessageInput, 
  ChatHeader,
  SearchMessagesOverlay,
  ChatSettingsOverlay,
  type MessageListHandle,
} from '@/widgets/ChatWindow';
import { Skeleton } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import { CHAT_PERMISSIONS } from '@/shared/types';

/**
 * Страница детального просмотра чата с сообщениями
 */
export function ChatDetailPage() {
  const { chatId } = useParams<{ chatId: string }>();
  const messageListRef = useRef<MessageListHandle>(null);

  // Состояние оверлеев
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [highlightedMessageId, setHighlightedMessageId] = useState<string | null>(null);

  // Получаем данные чата из списка чатов
  const { data: chats, isLoading: isLoadingChats } = useUserChats();
  const chat = chats?.find((c) => c.id === chatId);

  // Получаем роль пользователя для проверки прав на настройки
  const { data: myRole } = useMyRoleInChat(chatId);

  // Проверяем есть ли хотя бы одно право для показа кнопки настроек
  const hasSettingsPermission = 
    hasPermission(myRole, CHAT_PERMISSIONS.EDIT_CHAT) ||
    hasPermission(myRole, CHAT_PERMISSIONS.DELETE_CHAT) ||
    hasPermission(myRole, CHAT_PERMISSIONS.BAN_USER) ||
    hasPermission(myRole, CHAT_PERMISSIONS.CHANGE_ROLE);

  // Получаем сообщения
  const {
    data: messages,
    isLoading: isLoadingMessages,
  } = useChatMessages(chatId);

  // Мутация для отправки сообщений
  const sendMessage = useSendMessage(chatId!);

  const handleSendMessage = (content: string, files?: File[]) => {
    sendMessage.mutate({ content, files });
  };

  // Обработка выбора сообщения из поиска
  const handleSelectMessage = useCallback((messageId: string) => {
    setHighlightedMessageId(messageId);
    messageListRef.current?.scrollToMessage(messageId);
    // Сбрасываем подсветку через некоторое время
    setTimeout(() => setHighlightedMessageId(null), 2500);
  }, []);

  // Если чат не найден после загрузки
  if (!isLoadingChats && !chat) {
    return <Navigate to={ROUTES.CHATS} replace />;
  }

  // Загрузка
  if (isLoadingChats || !chat) {
    return <ChatDetailSkeleton />;
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <ChatHeader 
        chat={chat}
        onOpenSearch={() => setIsSearchOpen(true)}
        onOpenSettings={() => setIsSettingsOpen(true)}
        showSettingsButton={hasSettingsPermission}
      />

      {/* Messages */}
      <MessageList
        ref={messageListRef}
        messages={messages || []}
        chatId={chatId!}
        isLoading={isLoadingMessages}
        highlightedMessageId={highlightedMessageId}
      />

      {/* Input */}
      <MessageInput
        onSend={handleSendMessage}
        isLoading={sendMessage.isPending}
      />

      {/* Search Overlay */}
      <SearchMessagesOverlay
        chatId={chatId!}
        isOpen={isSearchOpen}
        onClose={() => setIsSearchOpen(false)}
        onSelectMessage={handleSelectMessage}
      />

      {/* Settings Overlay */}
      <ChatSettingsOverlay
        chat={chat}
        isOpen={isSettingsOpen}
        onClose={() => setIsSettingsOpen(false)}
      />
    </div>
  );
}

// Skeleton loader
function ChatDetailSkeleton() {
  return (
    <div className="flex flex-col h-full">
      {/* Header skeleton */}
      <div className="flex items-center gap-4 px-4 py-3 border-b border-neutral-800">
        <Skeleton variant="circular" className="w-10 h-10" />
        <div className="flex-1">
          <Skeleton className="h-5 w-32 mb-1" />
          <Skeleton className="h-3 w-20" />
        </div>
      </div>

      {/* Messages skeleton */}
      <div className="flex-1 p-4 space-y-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div
            key={i}
            className={`flex gap-3 ${i % 2 === 1 ? 'justify-end' : ''}`}
          >
            {i % 2 === 0 && (
              <Skeleton variant="circular" className="w-8 h-8" />
            )}
            <Skeleton className="h-16 w-48 rounded-2xl" />
          </div>
        ))}
      </div>

      {/* Input skeleton */}
      <div className="border-t border-neutral-800 p-4">
        <Skeleton className="h-12 w-full rounded-xl" />
      </div>
    </div>
  );
}

