import { useState } from 'react';
import { Outlet, useParams } from 'react-router-dom';
import { MessageSquare } from 'lucide-react';
import { ChatList } from '@/widgets/ChatList';
import { CreateChatModal } from '@/features/chat/create-chat';
import { cn } from '@/shared/lib/cn';

/**
 * Страница чатов с боковым списком
 */
export function ChatsPage() {
  const { chatId } = useParams<{ chatId: string }>();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  return (
    <div className="flex h-full">
      {/* Sidebar с списком чатов */}
      <aside
        className={cn(
          'w-80 border-r border-neutral-800 bg-neutral-900/50 flex-shrink-0',
          'hidden md:flex flex-col',
          // На мобильных показываем только если чат не выбран
          !chatId && 'flex'
        )}
      >
        <ChatList onCreateChat={() => setIsCreateModalOpen(true)} />
      </aside>

      {/* Основная область */}
      <main className="flex-1 flex flex-col min-w-0">
        {chatId ? (
          <Outlet />
        ) : (
          <EmptyChatState />
        )}
      </main>

      {/* Модальное окно создания чата */}
      <CreateChatModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
      />
    </div>
  );
}

/**
 * Пустое состояние когда чат не выбран
 */
function EmptyChatState() {
  return (
    <div className="flex-1 flex items-center justify-center">
      <div className="text-center">
        <div className="w-20 h-20 rounded-2xl bg-neutral-800 flex items-center justify-center mx-auto mb-6">
          <MessageSquare size={40} className="text-neutral-600" />
        </div>
        <h2 className="text-xl font-semibold text-neutral-300 mb-2">
          Выберите чат
        </h2>
        <p className="text-neutral-500 max-w-sm">
          Выберите чат из списка слева или создайте новый, чтобы начать общение
        </p>
      </div>
    </div>
  );
}

