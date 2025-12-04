import { useState } from 'react';
import { Edit3, Shield, Calendar, Mail, User as UserIcon } from 'lucide-react';
import { Card, Avatar, Badge, Button, Skeleton } from '@/shared/ui';
import { EditProfileForm } from '@/features/profile';
import { useCurrentUser } from '@/entities/user';

/**
 * Страница профиля текущего пользователя
 */
export function ProfilePage() {
  const { data: user, isLoading, refetch } = useCurrentUser();
  const [isEditing, setIsEditing] = useState(false);

  if (isLoading || !user) {
    return <ProfileSkeleton />;
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-neutral-100">Мой профиль</h1>
          <p className="text-neutral-400 mt-1">
            Управление вашими личными данными
          </p>
        </div>
        {!isEditing && (
          <Button
            variant="secondary"
            leftIcon={<Edit3 size={18} />}
            onClick={() => setIsEditing(true)}
          >
            Редактировать
          </Button>
        )}
      </div>

      {isEditing ? (
        /* Edit mode */
        <Card className="p-6">
          <h2 className="text-lg font-semibold text-neutral-100 mb-6">
            Редактирование профиля
          </h2>
          <EditProfileForm
            user={user}
            onSuccess={() => {
              setIsEditing(false);
              // Принудительно перезагружаем данные с сервера
              refetch();
            }}
            onCancel={() => setIsEditing(false)}
          />
        </Card>
      ) : (
        /* View mode */
        <>
          {/* Profile card */}
          <Card className="p-6">
            <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6">
              {/* Avatar */}
              <Avatar
                file={user.avatar}
                fallback={user.Username}
                size="xl"
                className="w-32 h-32 text-3xl"
              />

              {/* Main info */}
              <div className="flex-1 text-center sm:text-left">
                <h2 className="text-2xl font-bold text-neutral-100">
                  {user.Username}
                </h2>
                <p className="text-neutral-400 mt-1">{user.Email}</p>

                <div className="flex flex-wrap gap-2 mt-3 justify-center sm:justify-start">
                  {user.Role && (
                    <Badge variant="primary">
                      <Shield size={12} className="mr-1" />
                      {user.Role.Name}
                    </Badge>
                  )}
                  {user.Age && (
                    <Badge variant="default">
                      <Calendar size={12} className="mr-1" />
                      {user.Age} лет
                    </Badge>
                  )}
                  {user.Gender && (
                    <Badge variant="default">
                      {user.Gender === 'male'
                        ? '♂ Мужской'
                        : user.Gender === 'female'
                        ? '♀ Женский'
                        : '⚧ Другой'}
                    </Badge>
                  )}
                </div>

                {user.Description && (
                  <p className="text-neutral-300 mt-4 leading-relaxed">
                    {user.Description}
                  </p>
                )}
              </div>
            </div>
          </Card>

          {/* Details card */}
          <Card className="p-6">
            <h3 className="text-lg font-semibold text-neutral-100 mb-4">
              Детальная информация
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <InfoRow
                icon={<UserIcon size={18} />}
                label="ID пользователя"
                value={user.ID}
                mono
              />
              <InfoRow
                icon={<Mail size={18} />}
                label="Email"
                value={user.Email}
              />
              <InfoRow
                icon={<Shield size={18} />}
                label="Роль"
                value={user.Role?.Name || 'Не назначена'}
              />
              <InfoRow
                icon={<Calendar size={18} />}
                label="Возраст"
                value={user.Age ? `${user.Age} лет` : 'Не указан'}
              />
            </div>
          </Card>

          {/* Permissions card */}
          {user.Role?.Permissions && user.Role.Permissions.length > 0 && (
            <Card className="p-6">
              <h3 className="text-lg font-semibold text-neutral-100 mb-4">
                Ваши разрешения
              </h3>
              <div className="flex flex-wrap gap-2">
                {user.Role.Permissions.map((permission) => (
                  <Badge key={permission.ID} variant="default">
                    {permission.Name}
                  </Badge>
                ))}
              </div>
              <p className="text-xs text-neutral-500 mt-4">
                Разрешения определяются вашей ролью и управляются администратором
              </p>
            </Card>
          )}
        </>
      )}
    </div>
  );
}

// Info row component
interface InfoRowProps {
  icon: React.ReactNode;
  label: string;
  value: string;
  mono?: boolean;
}

function InfoRow({ icon, label, value, mono }: InfoRowProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg bg-neutral-800/50">
      <div className="text-neutral-500">{icon}</div>
      <div className="flex-1 min-w-0">
        <p className="text-xs text-neutral-500">{label}</p>
        <p
          className={`text-neutral-200 truncate ${
            mono ? 'font-mono text-sm' : ''
          }`}
        >
          {value}
        </p>
      </div>
    </div>
  );
}

// Skeleton loader
function ProfileSkeleton() {
  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-4 w-64 mt-2" />
        </div>
        <Skeleton className="h-10 w-32" />
      </div>
      <Card className="p-6">
        <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6">
          <Skeleton variant="circular" className="w-32 h-32" />
          <div className="flex-1 space-y-3">
            <Skeleton className="h-8 w-48" />
            <Skeleton className="h-4 w-32" />
            <div className="flex gap-2">
              <Skeleton className="h-6 w-20" />
              <Skeleton className="h-6 w-20" />
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}

