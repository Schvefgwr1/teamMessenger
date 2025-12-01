import { useState, useEffect } from 'react';
import { User, FileText, Users, Calendar } from 'lucide-react';
import { Button, Input, Card, toast } from '@/shared/ui';
import { useUpdateProfile, UpdateUserRequest } from '@/entities/user';
import { AvatarUpload } from '@/features/profile/upload-avatar';
import type { AuthUser } from '@/entities/session';

interface EditProfileFormProps {
  /** Текущие данные пользователя */
  user: AuthUser;
  /** Callback после успешного сохранения */
  onSuccess?: () => void;
  /** Callback при отмене */
  onCancel?: () => void;
}

/**
 * Форма редактирования профиля пользователя
 */
export function EditProfileForm({ user, onSuccess, onCancel }: EditProfileFormProps) {
  const updateProfile = useUpdateProfile();

  // Состояние формы
  const [formData, setFormData] = useState<UpdateUserRequest>({
    username: user.Username || '',
    description: user.Description || '',
    gender: user.Gender || '',
    age: user.Age || undefined,
  });

  // Состояние аватара
  const [newAvatar, setNewAvatar] = useState<File | null>(null);

  // Сбрасываем форму при изменении пользователя
  useEffect(() => {
    setFormData({
      username: user.Username || '',
      description: user.Description || '',
      gender: user.Gender || '',
      age: user.Age || undefined,
    });
    setNewAvatar(null);
  }, [user]);

  /**
   * Создаёт объект только с изменёнными полями
   * Сравнивает текущие значения формы с исходными значениями пользователя
   */
  const getChangedFields = (): UpdateUserRequest => {
    const changed: UpdateUserRequest = {};

    // Username - только если изменился
    if (formData.username !== user.Username) {
      changed.username = formData.username || undefined;
    }

    // Description - только если изменился (учитываем пустую строку vs undefined)
    const currentDescription = formData.description || '';
    const originalDescription = user.Description || '';
    if (currentDescription !== originalDescription) {
      changed.description = currentDescription || undefined;
    }

    // Gender - только если изменился
    const currentGender = formData.gender || '';
    const originalGender = user.Gender || '';
    if (currentGender !== originalGender) {
      changed.gender = currentGender || undefined;
    }

    // Age - только если изменился
    if (formData.age !== user.Age) {
      changed.age = formData.age;
    }

    // RoleID НЕ включаем - это только для админов

    return changed;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Валидация
    if (!formData.username?.trim()) {
      toast.error('Имя пользователя обязательно');
      return;
    }

    if (formData.age !== undefined && (formData.age < 1 || formData.age > 150)) {
      toast.error('Укажите корректный возраст');
      return;
    }

    // Получаем только изменённые поля
    const changedFields = getChangedFields();

    // Если нет изменений в данных и нет нового аватара - ничего не отправляем
    if (Object.keys(changedFields).length === 0 && !newAvatar) {
      toast.info('Нет изменений для сохранения');
      return;
    }

    try {
      await updateProfile.mutateAsync({
        data: changedFields,
        avatar: newAvatar || undefined,
      });
      onSuccess?.();
    } catch {
      // Ошибка обработана в хуке
    }
  };

  const handleChange = (
    field: keyof UpdateUserRequest,
    value: string | number | undefined
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  // Проверяем, были ли изменения
  const hasChanges =
    formData.username !== user.Username ||
    formData.description !== (user.Description || '') ||
    formData.gender !== (user.Gender || '') ||
    formData.age !== user.Age ||
    newAvatar !== null;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Avatar section */}
      <div className="flex flex-col items-center gap-4">
        <AvatarUpload
          currentFile={user.avatar}
          fallback={user.Username}
          onFileSelect={setNewAvatar}
          size="xl"
        />
        <p className="text-sm text-neutral-400">
          Нажмите для изменения аватара
        </p>
      </div>

      {/* Form fields */}
      <div className="space-y-4">
        <Input
          label="Имя пользователя"
          placeholder="username"
          value={formData.username || ''}
          onChange={(e) => handleChange('username', e.target.value)}
          leftIcon={<User size={18} />}
        />

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-neutral-300">
            О себе
          </label>
          <div className="relative">
            <FileText
              size={18}
              className="absolute left-3 top-3 text-neutral-500"
            />
            <textarea
              value={formData.description || ''}
              onChange={(e) => handleChange('description', e.target.value)}
              placeholder="Расскажите о себе..."
              rows={3}
              className="w-full pl-10 pr-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 resize-none"
            />
          </div>
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-neutral-300">
            Пол
          </label>
          <div className="relative">
            <Users
              size={18}
              className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500"
            />
            <select
              value={formData.gender || ''}
              onChange={(e) => handleChange('gender', e.target.value)}
              className="w-full h-10 pl-10 pr-3 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 appearance-none"
            >
              <option value="">Не указан</option>
              <option value="male">Мужской</option>
              <option value="female">Женский</option>
              <option value="other">Другой</option>
            </select>
          </div>
        </div>

        <Input
          label="Возраст"
          type="number"
          placeholder="25"
          value={formData.age?.toString() || ''}
          onChange={(e) =>
            handleChange(
              'age',
              e.target.value ? parseInt(e.target.value, 10) : undefined
            )
          }
          leftIcon={<Calendar size={18} />}
          min={1}
          max={150}
        />
      </div>

      {/* Read-only fields */}
      <Card variant="outlined" className="p-4 space-y-2">
        <p className="text-xs text-neutral-500 uppercase tracking-wider mb-3">
          Информация (только чтение)
        </p>
        <div className="grid grid-cols-2 gap-2 text-sm">
          <div>
            <span className="text-neutral-500">Email:</span>
            <span className="ml-2 text-neutral-300">{user.Email}</span>
          </div>
          <div>
            <span className="text-neutral-500">Роль:</span>
            <span className="ml-2 text-neutral-300">{user.Role?.Name}</span>
          </div>
        </div>
      </Card>

      {/* Actions */}
      <div className="flex gap-3 justify-end pt-4 border-t border-neutral-800">
        {onCancel && (
          <Button type="button" variant="secondary" onClick={onCancel}>
            Отмена
          </Button>
        )}
        <Button
          type="submit"
          isLoading={updateProfile.isPending}
          disabled={!hasChanges}
        >
          Сохранить изменения
        </Button>
      </div>
    </form>
  );
}

