import { useState, useRef, useEffect } from 'react';
import { Search, X, Mail, User } from 'lucide-react';
import { Input, Avatar, Skeleton } from '@/shared/ui';
import { useSearchUsers } from '@/entities/user';
import { useDebounce } from '@/shared/hooks';
import { cn } from '@/shared/lib/cn';
import type { UserSearchResult } from '@/shared/types';

interface UserSearchProps {
  selectedUsers: UserSearchResult[];
  onSelect: (user: UserSearchResult) => void;
  onRemove: (userId: string) => void;
  excludeUserIds?: string[];
}

/**
 * Компонент поиска и выбора пользователей
 * Определяет тип поиска (email или username) автоматически
 */
export function UserSearch({
  selectedUsers,
  onSelect,
  onRemove,
  excludeUserIds = [],
}: UserSearchProps) {
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const debouncedQuery = useDebounce(query, 300);
  const containerRef = useRef<HTMLDivElement>(null);

  // Определяем тип поиска
  const isEmailSearch = query.includes('@');

  // Поиск пользователей
  const { data: users, isLoading } = useSearchUsers(debouncedQuery, isOpen);

  // Фильтруем уже выбранных и исключённых
  const filteredUsers = users?.filter(
    (user) =>
      !selectedUsers.some((s) => s.id === user.id) &&
      !excludeUserIds.includes(user.id)
  );

  // Закрытие dropdown при клике вне
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleSelect = (user: UserSearchResult) => {
    onSelect(user);
    setQuery('');
    setIsOpen(false);
  };

  return (
    <div className="space-y-3">
      {/* Выбранные пользователи */}
      {selectedUsers.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selectedUsers.map((user) => (
            <div
              key={user.id}
              className="flex items-center gap-2 px-2 py-1 rounded-lg bg-neutral-800 text-sm"
            >
              <Avatar
                file={user.avatarFile}
                fallback={user.username}
                size="xs"
              />
              <span className="text-neutral-200">{user.username}</span>
              <button
                type="button"
                onClick={() => onRemove(user.id)}
                className="p-0.5 rounded hover:bg-neutral-700 text-neutral-400 hover:text-neutral-200 transition-colors"
              >
                <X size={14} />
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Поле поиска */}
      <div ref={containerRef} className="relative">
        <Input
          placeholder="Поиск по имени или email..."
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setIsOpen(true);
          }}
          onFocus={() => setIsOpen(true)}
          leftIcon={<Search size={18} />}
          rightIcon={
            query.length > 0 ? (
              <span className="text-xs text-neutral-500 flex items-center gap-1">
                {isEmailSearch ? (
                  <>
                    <Mail size={12} /> email
                  </>
                ) : (
                  <>
                    <User size={12} /> имя
                  </>
                )}
              </span>
            ) : undefined
          }
        />

        {/* Dropdown с результатами */}
        {isOpen && query.length >= 2 && (
          <div className="absolute z-50 w-full mt-1 py-1 rounded-lg bg-neutral-900 border border-neutral-800 shadow-xl max-h-60 overflow-y-auto">
            {isLoading ? (
              <div className="p-2 space-y-2">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="flex items-center gap-3 p-2">
                    <Skeleton variant="circular" className="w-8 h-8" />
                    <div className="flex-1 space-y-1">
                      <Skeleton className="h-4 w-24" />
                      <Skeleton className="h-3 w-32" />
                    </div>
                  </div>
                ))}
              </div>
            ) : filteredUsers && filteredUsers.length > 0 ? (
              filteredUsers.map((user) => (
                <button
                  key={user.id}
                  type="button"
                  onClick={() => handleSelect(user)}
                  className={cn(
                    'w-full flex items-center gap-3 px-3 py-2',
                    'hover:bg-neutral-800 transition-colors',
                    'text-left'
                  )}
                >
                  <Avatar
                    file={user.avatarFile}
                    fallback={user.username}
                    size="sm"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-neutral-100 truncate">
                      {user.username}
                    </p>
                    <p className="text-xs text-neutral-500 truncate">
                      {user.email}
                    </p>
                  </div>
                </button>
              ))
            ) : (
              <p className="px-3 py-4 text-sm text-neutral-500 text-center">
                {query.length < 2
                  ? 'Введите минимум 2 символа'
                  : 'Пользователи не найдены'}
              </p>
            )}
          </div>
        )}
      </div>

      {/* Подсказка */}
      <p className="text-xs text-neutral-500">
        Введите имя пользователя или email для поиска
      </p>
    </div>
  );
}

