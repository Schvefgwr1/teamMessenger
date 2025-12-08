import { useState, useCallback, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { Search, X, Loader2 } from 'lucide-react';
import { Input, Avatar } from '@/shared/ui';
import { useSearchMessages } from '@/entities/chat';
import { useDebounce } from '@/shared/hooks';
import { formatMessageTime } from '@/shared/lib/formatDate';
import { useUserBrief } from '@/entities/user';
import type { Message } from '@/shared/types';

interface SearchMessagesOverlayProps {
  chatId: string;
  isOpen: boolean;
  onClose: () => void;
  onSelectMessage: (messageId: string) => void;
}

/**
 * Оверлей для поиска сообщений в чате
 */
export function SearchMessagesOverlay({
  chatId,
  isOpen,
  onClose,
  onSelectMessage,
}: SearchMessagesOverlayProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedQuery = useDebounce(searchQuery, 300);
  const inputRef = useRef<HTMLInputElement>(null);

  // Поиск сообщений
  const { data: messages, isLoading, isFetching } = useSearchMessages(
    chatId,
    debouncedQuery
  );

  // Фокус на поле при открытии
  useEffect(() => {
    if (isOpen) {
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  }, [isOpen]);

  // Обработка клика на результат
  const handleSelectMessage = useCallback(
    (messageId: string) => {
      onSelectMessage(messageId);
      onClose();
      setSearchQuery('');
    },
    [onSelectMessage, onClose]
  );

  // Обработка Escape
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return createPortal(
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
          />

          {/* Panel */}
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.2 }}
            className="fixed top-20 left-1/2 -translate-x-1/2 z-50 w-full max-w-lg"
          >
            <div className="bg-neutral-900 border border-neutral-800 rounded-xl shadow-2xl overflow-hidden">
              {/* Search Input */}
              <div className="p-4 border-b border-neutral-800">
                <div className="relative">
                  <Input
                    ref={inputRef}
                    placeholder="Поиск сообщений..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    leftIcon={
                      isLoading || isFetching ? (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      ) : (
                        <Search className="w-4 h-4" />
                      )
                    }
                    className="pr-10"
                  />
                  <button
                    onClick={onClose}
                    className="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-neutral-800 text-neutral-400 hover:text-neutral-100 transition-colors"
                  >
                    <X size={18} />
                  </button>
                </div>
              </div>

              {/* Results */}
              <div className="max-h-80 overflow-y-auto">
                {debouncedQuery.length < 2 ? (
                  <div className="p-8 text-center text-neutral-500">
                    Введите минимум 2 символа для поиска
                  </div>
                ) : isLoading ? (
                  <div className="p-8 flex justify-center">
                    <Loader2 className="w-6 h-6 animate-spin text-primary-500" />
                  </div>
                ) : !messages || messages.length === 0 ? (
                  <div className="p-8 text-center text-neutral-500">
                    Сообщения не найдены
                  </div>
                ) : (
                  <ul className="divide-y divide-neutral-800">
                    {messages.map((message) => (
                      <SearchResultItem
                        key={message.id}
                        message={message}
                        searchQuery={debouncedQuery}
                        onClick={() => handleSelectMessage(message.id)}
                      />
                    ))}
                  </ul>
                )}
              </div>

              {/* Footer */}
              {messages && messages.length > 0 && (
                <div className="p-3 border-t border-neutral-800 bg-neutral-900/50">
                  <p className="text-xs text-neutral-500 text-center">
                    Найдено {messages.length} сообщений
                  </p>
                </div>
              )}
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>,
    document.body
  );
}

// Компонент одного результата поиска
interface SearchResultItemProps {
  message: Message;
  searchQuery: string;
  onClick: () => void;
}

function SearchResultItem({ message, searchQuery, onClick }: SearchResultItemProps) {
  // Получаем информацию об отправителе
  const { data: sender } = useUserBrief(message.senderID || undefined, undefined);

  // Подсветка найденного текста
  const highlightText = (text: string, query: string) => {
    if (!query) return text;
    
    const parts = text.split(new RegExp(`(${query})`, 'gi'));
    return parts.map((part, i) =>
      part.toLowerCase() === query.toLowerCase() ? (
        <mark key={i} className="bg-primary-500/30 text-primary-300 rounded">
          {part}
        </mark>
      ) : (
        part
      )
    );
  };

  return (
    <li>
      <button
        type="button"
        onClick={onClick}
        className="w-full px-4 py-3 flex items-start gap-3 hover:bg-neutral-800/50 transition-colors text-left"
      >
        <Avatar
          file={sender?.avatarFile}
          fallback={sender?.username || '?'}
          size="sm"
        />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-sm font-medium text-neutral-100">
              {sender?.username || 'Неизвестный'}
            </span>
            <span className="text-xs text-neutral-500">
              {formatMessageTime(message.createdAt)}
            </span>
          </div>
          <p className="text-sm text-neutral-400 line-clamp-2">
            {highlightText(message.content, searchQuery)}
          </p>
        </div>
      </button>
    </li>
  );
}

