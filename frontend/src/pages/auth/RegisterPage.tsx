import { useState, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { authApi, createRegisterFormData, DEFAULT_ROLE_ID } from '@/entities/session';
import { Button, Input, Card, toast } from '@/shared/ui';
import { ROUTES } from '@/shared/constants';
import { formatFileSize, MAX_FILE_SIZE } from '@/shared/lib/formatFileSize';
import { User, Mail, Lock, Eye, EyeOff, Calendar, Upload, X } from 'lucide-react';

export function RegisterPage() {
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    description: '',
    gender: '',
    age: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [avatar, setAvatar] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const registerMutation = useMutation({
    mutationFn: async (data: typeof formData) => {
      // Парсим age в число, проверяя на валидность
      const ageValue = data.age ? parseInt(data.age, 10) : undefined;
      if (ageValue === undefined || isNaN(ageValue) || ageValue <= 0) {
        throw new Error('Укажите корректный возраст');
      }

      const registerData = createRegisterFormData({
        username: data.username,
        email: data.email,
        password: data.password,
        description: data.description || undefined,
        gender: data.gender || undefined,
        age: ageValue,
        roleID: DEFAULT_ROLE_ID,
      }, avatar || undefined);
      return authApi.register(registerData);
    },
    onSuccess: () => {
      toast.success('Регистрация успешна! Войдите в систему.');
      navigate(ROUTES.LOGIN);
    },
    onError: (error: Error & { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Ошибка регистрации';
      toast.error(message);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Валидация
    if (!formData.username || !formData.email || !formData.password) {
      toast.error('Заполните обязательные поля');
      return;
    }
    if (formData.password !== formData.confirmPassword) {
      toast.error('Пароли не совпадают');
      return;
    }
    if (formData.password.length < 6) {
      toast.error('Пароль должен быть не менее 6 символов');
      return;
    }
    if (!formData.age || parseInt(formData.age, 10) <= 0) {
      toast.error('Укажите корректный возраст');
      return;
    }

    registerMutation.mutate(formData);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Проверка размера файла
    if (file.size > MAX_FILE_SIZE) {
      toast.error(`Файл слишком большой. Максимальный размер: ${formatFileSize(MAX_FILE_SIZE)}`);
      return;
    }

    // Проверка типа файла (только изображения)
    if (!file.type.startsWith('image/')) {
      toast.error('Пожалуйста, выберите изображение');
      return;
    }

    setAvatar(file);

    // Создаем превью
    const reader = new FileReader();
    reader.onloadend = () => {
      setAvatarPreview(reader.result as string);
    };
    reader.readAsDataURL(file);
  };

  const handleRemoveAvatar = () => {
    setAvatar(null);
    setAvatarPreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <div className="min-h-screen bg-neutral-950 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-gradient-to-br from-primary-400 to-primary-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
            <span className="text-2xl font-bold text-white">TM</span>
          </div>
          <h1 className="text-2xl font-bold text-neutral-100">Регистрация</h1>
          <p className="text-neutral-400 mt-2">Создайте новый аккаунт</p>
        </div>

        <Card>
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Имя пользователя"
              placeholder="username"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              leftIcon={<User size={18} />}
              autoComplete="username"
            />

            <Input
              label="Email"
              type="email"
              placeholder="user@example.com"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              leftIcon={<Mail size={18} />}
              autoComplete="email"
            />

            <Input
              label="Пароль"
              type={showPassword ? 'text' : 'password'}
              placeholder="••••••••"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              leftIcon={<Lock size={18} />}
              rightIcon={
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="hover:text-neutral-300 transition-colors"
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              }
              autoComplete="new-password"
              hint="Минимум 6 символов"
            />

            <Input
              label="Подтвердите пароль"
              type={showPassword ? 'text' : 'password'}
              placeholder="••••••••"
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
              leftIcon={<Lock size={18} />}
              autoComplete="new-password"
            />

            <Input
              label="Возраст"
              type="number"
              placeholder="25"
              value={formData.age}
              onChange={(e) => setFormData({ ...formData, age: e.target.value })}
              leftIcon={<Calendar size={18} />}
              min="1"
              max="150"
            />

            {/* Description */}
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-neutral-300">
                Описание (опционально)
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Расскажите о себе..."
                rows={3}
                className="w-full px-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 resize-none"
              />
            </div>

            {/* Gender */}
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-neutral-300">
                Пол (опционально)
              </label>
              <select
                value={formData.gender}
                onChange={(e) => setFormData({ ...formData, gender: e.target.value })}
                className="w-full h-10 px-3 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
              >
                <option value="">Не указан</option>
                <option value="male">Мужской</option>
                <option value="female">Женский</option>
                <option value="other">Другой</option>
              </select>
            </div>

            {/* Avatar Upload */}
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-neutral-300">
                Аватар (опционально)
              </label>
              {avatarPreview ? (
                <div className="space-y-2">
                  <div className="relative inline-block">
                    <div className="w-24 h-24 rounded-full overflow-hidden border-2 border-primary-500">
                      <img
                        src={avatarPreview}
                        alt="Avatar preview"
                        className="w-full h-full object-cover"
                      />
                    </div>
                    <button
                      type="button"
                      onClick={handleRemoveAvatar}
                      className="absolute -top-1 -right-1 p-1 rounded-full bg-error text-white hover:bg-red-600 transition-colors shadow-md"
                      title="Удалить аватар"
                    >
                      <X size={14} />
                    </button>
                  </div>
                  {avatar && (
                    <p className="text-xs text-neutral-400">
                      {avatar.name} ({formatFileSize(avatar.size)})
                    </p>
                  )}
                </div>
              ) : (
                <div>
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    onChange={handleFileChange}
                    className="hidden"
                    id="avatar-upload"
                  />
                  <label
                    htmlFor="avatar-upload"
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-300 hover:bg-neutral-800 hover:text-neutral-100 cursor-pointer transition-colors w-fit"
                  >
                    <Upload size={18} />
                    <span>Загрузить аватар</span>
                  </label>
                  <p className="text-xs text-neutral-500 mt-1">
                    Максимальный размер: {formatFileSize(MAX_FILE_SIZE)}
                  </p>
                </div>
              )}
            </div>

            <Button
              type="submit"
              className="w-full"
              isLoading={registerMutation.isPending}
            >
              Зарегистрироваться
            </Button>
          </form>

          <div className="mt-6 text-center text-sm">
            <span className="text-neutral-400">Уже есть аккаунт? </span>
            <Link
              to={ROUTES.LOGIN}
              className="text-primary-400 hover:text-primary-300 font-medium"
            >
              Войти
            </Link>
          </div>
        </Card>
      </div>
    </div>
  );
}

