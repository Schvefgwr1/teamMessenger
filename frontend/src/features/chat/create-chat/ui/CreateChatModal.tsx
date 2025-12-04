import { useState } from 'react';
import { Users, FileText } from 'lucide-react';
import { Modal, Input, Button, toast } from '@/shared/ui';
import { useCreateChat } from '@/entities/chat';
import { useAuthStore } from '@/entities/session';
import { AvatarUpload } from '@/features/profile/upload-avatar';
import { UserSearch } from './UserSearch';
import type { UserSearchResult } from '@/shared/types';

interface CreateChatModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

/**
 * Модальное окно создания нового чата
 */
export function CreateChatModal({ open, onOpenChange }: CreateChatModalProps) {
  const createChat = useCreateChat();
  const { user } = useAuthStore();

  const [formData, setFormData] = useState({
    name: '',
    description: '',
  });
  const [selectedUsers, setSelectedUsers] = useState<UserSearchResult[]>([]);
  const [avatar, setAvatar] = useState<File | null>(null);

  const handleSelectUser = (user: UserSearchResult) => {
    setSelectedUsers((prev) => [...prev, user]);
  };

  const handleRemoveUser = (userId: string) => {
    setSelectedUsers((prev) => prev.filter((u) => u.id !== userId));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim()) {
      toast.error('Введите название чата');
      return;
    }

    // Собираем ID выбранных пользователей
    const userIDs = selectedUsers.map((u) => u.id);

    try {
      await createChat.mutateAsync({
        data: {
          name: formData.name.trim(),
          description: formData.description.trim() || undefined,
          userIDs,
        },
        avatar: avatar || undefined,
      });

      // Сброс формы и закрытие
      resetForm();
      onOpenChange(false);
    } catch {
      // Ошибка обработана в хуке
    }
  };

  const resetForm = () => {
    setFormData({ name: '', description: '' });
    setSelectedUsers([]);
    setAvatar(null);
  };

  const handleClose = () => {
    resetForm();
    onOpenChange(false);
  };

  return (
    <Modal open={open} onOpenChange={handleClose}>
      <Modal.Content title="Создать чат" description="Создайте новый групповой чат">
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Название чата"
            placeholder="Введите название..."
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            leftIcon={<Users size={18} />}
            autoFocus
          />

          {/* Avatar Upload */}
          <div className="flex flex-col items-center gap-2">
            <label className="text-sm font-medium text-neutral-300 self-start">
              Аватар чата (опционально)
            </label>
            <AvatarUpload
              currentFile={null}
              fallback={formData.name || 'Чат'}
              onFileSelect={setAvatar}
              size="lg"
            />
          </div>

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
                placeholder="О чём этот чат..."
                rows={3}
                className="w-full pl-10 pr-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 resize-none"
              />
            </div>
          </div>

          {/* Поиск и выбор участников */}
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-neutral-300">
              Участники
            </label>
            <UserSearch
              selectedUsers={selectedUsers}
              onSelect={handleSelectUser}
              onRemove={handleRemoveUser}
              excludeUserIds={user?.ID ? [user.ID] : []}
            />
          </div>

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
              isLoading={createChat.isPending}
            >
              Создать
            </Button>
          </div>
        </form>
      </Modal.Content>
    </Modal>
  );
}

