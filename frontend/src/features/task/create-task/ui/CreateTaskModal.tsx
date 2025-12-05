import { useState, useRef, useMemo, useEffect } from 'react';
import { FileText, User, Paperclip, X } from 'lucide-react';
import { Modal, Input, Button, toast } from '@/shared/ui';
import { useCreateTask } from '@/entities/task';
import { useUserChats, useChatMembers } from '@/entities/chat';
import { useAuthStore } from '@/entities/session';
import { UserSearch } from '@/features/chat/create-chat/ui/UserSearch';
import type { UserSearchResult } from '@/shared/types';
import type { ChatMemberResponse } from '@/entities/chat';

interface CreateTaskModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

/**
 * Компонент для загрузки участников одного чата
 */
function ChatMembersLoader({ 
  chatId, 
  onMembersLoaded 
}: { 
  chatId: string; 
  onMembersLoaded: (members: ChatMemberResponse[]) => void;
}) {
  const { data } = useChatMembers(chatId);
  
  useEffect(() => {
    if (data) {
      onMembersLoaded(data);
    }
  }, [data, onMembersLoaded]);
  
  return null;
}

/**
 * Модальное окно создания новой задачи
 */
export function CreateTaskModal({ open, onOpenChange }: CreateTaskModalProps) {
  const createTask = useCreateTask();
  const { data: chats } = useUserChats();
  const { user } = useAuthStore();

  const [formData, setFormData] = useState({
    title: '',
    description: '',
  });
  const [selectedExecutor, setSelectedExecutor] = useState<UserSearchResult | null>(null);
  const [selectedChatId, setSelectedChatId] = useState<string>('');
  const [files, setFiles] = useState<File[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  // Состояние для хранения участников чатов
  const [chatMembersMap, setChatMembersMap] = useState<Record<string, ChatMemberResponse[]>>({});

  // Фильтруем чаты: показываем только те, где есть и создатель, и исполнитель
  const availableChats = useMemo(() => {
    if (!chats || !user?.ID) return [];
    if (!selectedExecutor) return []; // Если исполнитель не выбран, не показываем чаты

    return chats.filter((chat) => {
      const members = chatMembersMap[chat.id];
      if (!members || members.length === 0) return false;

      // Проверяем что создатель (текущий пользователь) в чате
      const creatorInChat = members.some((member) => member.userId === user.ID);
      if (!creatorInChat) return false;

      // Проверяем что исполнитель в чате и не забанен
      const executorInChat = members.some(
        (member) => member.userId === selectedExecutor.id && member.roleName !== 'banned'
      );
      return executorInChat;
    });
  }, [chats, user?.ID, selectedExecutor, chatMembersMap]);

  const handleSelectExecutor = (user: UserSearchResult) => {
    setSelectedExecutor(user);
  };

  const handleRemoveExecutor = () => {
    setSelectedExecutor(null);
    setSelectedChatId(''); // Сбрасываем выбранный чат при удалении исполнителя
    setChatMembersMap({}); // Очищаем кеш участников
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = Array.from(e.target.files || []);
    setFiles((prev) => [...prev, ...selectedFiles]);
  };

  const handleRemoveFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.title.trim()) {
      toast.error('Введите название задачи');
      return;
    }

    try {
      await createTask.mutateAsync({
        title: formData.title.trim(),
        description: formData.description.trim() || undefined,
        executorId: selectedExecutor?.id,
        chatId: selectedChatId || undefined,
        files: files.length > 0 ? files : undefined,
      });

      // Сброс формы и закрытие
      resetForm();
      onOpenChange(false);
    } catch {
      // Ошибка обработана в хуке
    }
  };

  const resetForm = () => {
    setFormData({ title: '', description: '' });
    setSelectedExecutor(null);
    setSelectedChatId('');
    setFiles([]);
    setChatMembersMap({}); // Очищаем кеш участников
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleClose = () => {
    resetForm();
    onOpenChange(false);
  };

  return (
    <Modal open={open} onOpenChange={handleClose}>
      <Modal.Content title="Создать задачу" description="Создайте новую задачу">
        {/* Загружаем участников для всех чатов параллельно (только если выбран исполнитель) */}
        {selectedExecutor && chats && chats.map((chat) => (
          <ChatMembersLoader
            key={chat.id}
            chatId={chat.id}
            onMembersLoaded={(members) => {
              setChatMembersMap((prev) => ({ ...prev, [chat.id]: members }));
            }}
          />
        ))}
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Название задачи */}
          <Input
            label="Название задачи"
            placeholder="Введите название..."
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            leftIcon={<FileText size={18} />}
            autoFocus
            required
          />

          {/* Описание */}
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-neutral-300">
              Описание (опционально)
            </label>
            <div className="relative">
              <FileText
                size={18}
                className="absolute left-3 top-3 text-neutral-500"
              />
              <textarea
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                placeholder="Опишите задачу..."
                rows={4}
                className="w-full pl-10 pr-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 resize-none"
              />
            </div>
          </div>

          {/* Исполнитель */}
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-neutral-300">
              Исполнитель
            </label>
            {selectedExecutor ? (
              <div className="flex items-center gap-2 px-3 py-2 rounded-lg bg-neutral-800">
                <User size={16} className="text-neutral-400" />
                <span className="flex-1 text-sm text-neutral-200">
                  {selectedExecutor.username}
                </span>
                <button
                  type="button"
                  onClick={handleRemoveExecutor}
                  className="p-1 rounded hover:bg-neutral-700 text-neutral-400 hover:text-neutral-200 transition-colors"
                >
                  <X size={14} />
                </button>
              </div>
            ) : (
              <UserSearch
                selectedUsers={[]}
                onSelect={handleSelectExecutor}
                onRemove={() => {}}
              />
            )}
          </div>

          {/* Связанный чат */}
          {selectedExecutor && (
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-neutral-300">
                Связанный чат (опционально)
              </label>
              {availableChats.length > 0 ? (
                <select
                  value={selectedChatId}
                  onChange={(e) => setSelectedChatId(e.target.value)}
                  className="w-full h-10 px-3 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                >
                  <option value="">Не выбран</option>
                  {availableChats.map((chat) => (
                    <option key={chat.id} value={chat.id}>
                      {chat.name}
                    </option>
                  ))}
                </select>
              ) : (
                <div className="px-3 py-2 rounded-lg bg-neutral-800 border border-neutral-700 text-sm text-neutral-400">
                  Нет чатов, где вы оба являетесь участниками. Выберите исполнителя или создайте общий чат.
                </div>
              )}
            </div>
          )}

          {/* Файлы */}
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-neutral-300">
              Прикрепленные файлы (опционально)
            </label>
            <div className="space-y-2">
              <input
                ref={fileInputRef}
                type="file"
                multiple
                onChange={handleFileSelect}
                className="hidden"
                accept="*/*"
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => fileInputRef.current?.click()}
                leftIcon={<Paperclip size={16} />}
              >
                Добавить файлы
              </Button>
              {files.length > 0 && (
                <div className="space-y-1">
                  {files.map((file, index) => (
                    <div
                      key={index}
                      className="flex items-center gap-2 px-3 py-2 rounded-lg bg-neutral-800"
                    >
                      <FileText size={16} className="text-neutral-400" />
                      <span className="flex-1 text-sm text-neutral-300 truncate">
                        {file.name}
                      </span>
                      <span className="text-xs text-neutral-500">
                        {(file.size / 1024).toFixed(1)} KB
                      </span>
                      <button
                        type="button"
                        onClick={() => handleRemoveFile(index)}
                        className="p-1 rounded hover:bg-neutral-700 text-neutral-400 hover:text-neutral-200 transition-colors"
                      >
                        <X size={14} />
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Кнопки */}
          <div className="flex gap-3 justify-end pt-4">
            <Button
              type="button"
              variant="secondary"
              onClick={handleClose}
            >
              Отмена
            </Button>
            <Button
              type="submit"
              isLoading={createTask.isPending}
            >
              Создать
            </Button>
          </div>
        </form>
      </Modal.Content>
    </Modal>
  );
}

