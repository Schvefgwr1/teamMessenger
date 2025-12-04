import {useEffect, useState} from 'react';
import {createPortal} from 'react-dom';
import {AnimatePresence, motion} from 'framer-motion';
import {useNavigate} from 'react-router-dom';
import {AlertTriangle, Ban, Camera, Edit3, Loader2, Settings, Trash2, UserCog, Users, X,} from 'lucide-react';
import {useForm} from 'react-hook-form';
import {zodResolver} from '@hookform/resolvers/zod';
import {z} from 'zod';
import {Avatar, Badge, Button, Input, Modal} from '@/shared/ui';
import {cn} from '@/shared/lib/cn';
import {
    chatKeys,
    type ChatMemberResponse,
    hasPermission,
    useBanUser,
    useChangeUserRole,
    useChatMembers,
    useChatRoles,
    useDeleteChat,
    useMyRoleInChat,
    useUpdateChat,
} from '@/entities/chat';
import {useUserBrief} from '@/entities/user';
import {useAuthStore} from '@/entities/session';
import {useQueryClient} from '@tanstack/react-query';
import {ROUTES} from '@/shared/constants';
import {type Chat, CHAT_PERMISSIONS} from '@/shared/types';

interface ChatSettingsOverlayProps {
    chat: Chat;
    isOpen: boolean;
    onClose: () => void;
}

/**
 * Оверлей для настроек чата
 * Отображается если у пользователя есть хотя бы один из permissions:
 * - edit_chat: редактирование чата
 * - delete_chat: удаление чата
 * - ban_user: бан пользователей
 * - change_role: изменение ролей
 */
export function ChatSettingsOverlay({
                                        chat,
                                        isOpen,
                                        onClose,
                                    }: ChatSettingsOverlayProps) {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const {user} = useAuthStore();

    // Получаем роль текущего пользователя в чате
    const {data: myRole, isLoading: isLoadingRole} = useMyRoleInChat(chat.id);

    // Проверяем permissions
    const canEditChat = hasPermission(myRole, CHAT_PERMISSIONS.EDIT_CHAT);
    const canDeleteChat = hasPermission(myRole, CHAT_PERMISSIONS.DELETE_CHAT);
    const canBanUser = hasPermission(myRole, CHAT_PERMISSIONS.BAN_USER);
    const canChangeRole = hasPermission(myRole, CHAT_PERMISSIONS.CHANGE_ROLE);

    const hasAnyPermission = canEditChat || canDeleteChat || canBanUser || canChangeRole;

    // Загружаем участников чата только если есть права на бан или изменение ролей
    const {data: members, isLoading: isLoadingMembers} = useChatMembers(
        (canBanUser || canChangeRole) ? chat.id : undefined
    );

    // Модалки подтверждения
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [userToBan, setUserToBan] = useState<string | null>(null);
    const [userToChangeRole, setUserToChangeRole] = useState<string | null>(null);

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

    // Мутации
    const deleteChat = useDeleteChat();
    const banUser = useBanUser(chat.id);
    const changeUserRole = useChangeUserRole(chat.id);

    // Обработка удаления чата
    const handleDeleteChat = () => {
        deleteChat.mutate(chat.id, {
            onSuccess: () => {
                setShowDeleteConfirm(false);
                onClose();
                // Перенаправляем на страницу чатов без выбранного чата
                navigate(ROUTES.CHATS);
            },
        });
    };

    // Обработка бана
    const handleBanUser = (userId: string) => {
        banUser.mutate(userId, {
            onSuccess: () => {
                setUserToBan(null);
                // Инвалидируем кеш чатов и участников
                queryClient.invalidateQueries({queryKey: chatKeys.lists()});
                queryClient.invalidateQueries({queryKey: [...chatKeys.detail(chat.id), 'members']});
            },
        });
    };

    if (!isOpen) return null;

    return createPortal(
        <AnimatePresence>
            {isOpen && (
                <>
                    {/* Backdrop */}
                    <motion.div
                        initial={{opacity: 0}}
                        animate={{opacity: 1}}
                        exit={{opacity: 0}}
                        onClick={onClose}
                        className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
                    />

                    {/* Panel */}
                    <motion.div
                        initial={{opacity: 0, x: 300}}
                        animate={{opacity: 1, x: 0}}
                        exit={{opacity: 0, x: 300}}
                        transition={{duration: 0.2, ease: 'easeOut'}}
                        className="fixed top-0 right-0 h-full w-full max-w-md z-50 bg-neutral-900 border-l border-neutral-800 shadow-2xl overflow-hidden flex flex-col"
                    >
                        {/* Header */}
                        <div className="flex items-center justify-between px-6 py-4 border-b border-neutral-800">
                            <div className="flex items-center gap-3">
                                <Settings size={20} className="text-primary-500"/>
                                <h2 className="text-lg font-semibold text-neutral-100">
                                    Настройки чата
                                </h2>
                            </div>
                            <button
                                onClick={onClose}
                                className="p-2 rounded-lg text-neutral-400 hover:text-neutral-100 hover:bg-neutral-800 transition-colors"
                            >
                                <X size={20}/>
                            </button>
                        </div>

                        {/* Content */}
                        <div className="flex-1 overflow-y-auto">
                            {isLoadingRole ? (
                                <div className="flex items-center justify-center py-12">
                                    <Loader2 className="w-8 h-8 animate-spin text-primary-500"/>
                                </div>
                            ) : !hasAnyPermission ? (
                                <div className="p-6 text-center text-neutral-500">
                                    У вас нет прав для управления этим чатом
                                </div>
                            ) : (
                                <div className="p-6 space-y-6">
                                    {/* Редактирование чата */}
                                    {canEditChat && (
                                        <EditChatSection chat={chat} onClose={onClose}/>
                                    )}

                                    {/* Управление участниками */}
                                    {(canBanUser || canChangeRole) && (
                                        isLoadingMembers ? (
                                            <div className="flex items-center justify-center py-8">
                                                <Loader2 className="w-6 h-6 animate-spin text-primary-500"/>
                                            </div>
                                        ) : members && members.length > 0 ? (
                                            <MembersSection
                                                chatId={chat.id}
                                                members={members}
                                                canBanUser={canBanUser}
                                                canChangeRole={canChangeRole}
                                                onBanUser={setUserToBan}
                                                onChangeRole={setUserToChangeRole}
                                                myRoleName={myRole?.roleName}
                                                currentUserId={user?.ID}
                                            />
                                        ) : null
                                    )}

                                    {/* Удаление чата */}
                                    {canDeleteChat && (
                                        <div className="pt-4 border-t border-neutral-800">
                                            <h3 className="text-sm font-medium text-error mb-3 flex items-center gap-2">
                                                <AlertTriangle size={16}/>
                                                Опасная зона
                                            </h3>
                                            <Button
                                                variant="danger"
                                                className="w-full"
                                                onClick={() => setShowDeleteConfirm(true)}
                                                leftIcon={<Trash2 size={16}/>}
                                            >
                                                Удалить чат
                                            </Button>
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>

                        {/* Role info */}
                        {myRole && (
                            <div className="px-6 py-3 border-t border-neutral-800 bg-neutral-900/50">
                                <p className="text-xs text-neutral-500">
                                    Ваша роль: <span className="text-primary-400">{myRole.roleName}</span>
                                </p>
                            </div>
                        )}
                    </motion.div>

                    {/* Delete Confirmation Modal */}
                    {showDeleteConfirm && (
                        <Modal
                            open={showDeleteConfirm}
                            onOpenChange={(open) => {
                                if (!open) {
                                    setShowDeleteConfirm(false);
                                }
                            }}
                        >
                            <Modal.Content
                                title="Удалить чат?"
                                description="Это действие нельзя отменить. Все сообщения будут удалены."
                            >
                                <div className="flex gap-3 mt-6">
                                    <Button
                                        variant="secondary"
                                        className="flex-1"
                                        onClick={() => setShowDeleteConfirm(false)}
                                    >
                                        Отмена
                                    </Button>
                                    <Button
                                        variant="danger"
                                        className="flex-1"
                                        onClick={handleDeleteChat}
                                        isLoading={deleteChat.isPending}
                                    >
                                        Удалить
                                    </Button>
                                </div>
                            </Modal.Content>
                        </Modal>
                    )}

                    {/* Ban Confirmation Modal */}
                    {userToBan && (
                        <Modal
                            open={!!userToBan}
                            onOpenChange={(open) => {
                                if (!open) {
                                    setUserToBan(null);
                                }
                            }}
                        >
                            <Modal.Content
                                title="Забанить пользователя?"
                                description="Пользователь не сможет отправлять сообщения в этот чат."
                            >
                                <div className="flex gap-3 mt-6">
                                    <Button
                                        variant="secondary"
                                        className="flex-1"
                                        onClick={() => setUserToBan(null)}
                                    >
                                        Отмена
                                    </Button>
                                    <Button
                                        variant="danger"
                                        className="flex-1"
                                        onClick={() => {
                                            if (userToBan) {
                                                handleBanUser(userToBan);
                                            }
                                        }}
                                        isLoading={banUser.isPending}
                                    >
                                        Забанить
                                    </Button>
                                </div>
                            </Modal.Content>
                        </Modal>
                    )}

                    {/* Change Role Modal */}
                    {userToChangeRole && (
                        <ChangeRoleModal
                            chatId={chat.id}
                            userId={userToChangeRole}
                            isOpen={!!userToChangeRole}
                            onClose={() => setUserToChangeRole(null)}
                            changeUserRole={changeUserRole}
                            myRoleName={myRole?.roleName}
                        />
                    )}
                </>
            )}
        </AnimatePresence>,
        document.body
    );
}

// Секция редактирования чата
const editChatSchema = z.object({
    name: z.string().min(1, 'Название обязательно').max(100),
    description: z.string().max(500).optional(),
});

type EditChatFormData = z.infer<typeof editChatSchema>;

function EditChatSection({chat, onClose}: { chat: Chat; onClose: () => void }) {
    const updateChat = useUpdateChat(chat.id);
    const [avatarFile, setAvatarFile] = useState<File | null>(null);

    const form = useForm<EditChatFormData>({
        resolver: zodResolver(editChatSchema),
        defaultValues: {
            name: chat.name,
            description: chat.description || '',
        },
    });

    const onSubmit = (data: EditChatFormData) => {
        updateChat.mutate(
            {
                data: {
                    name: data.name,
                    description: data.description,
                },
                avatar: avatarFile || undefined,
            },
            {
                onSuccess: () => {
                    onClose();
                },
            }
        );
    };

    return (
        <div>
            <h3 className="text-sm font-medium text-neutral-300 mb-4 flex items-center gap-2">
                <Edit3 size={16}/>
                Редактирование чата
            </h3>

            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                {/* Avatar */}
                <div className="flex items-center gap-4">
                    <div className="relative">
                        <Avatar
                            file={chat.avatarFile}
                            fallback={chat.name}
                            size="xl"
                        />
                        <label
                            className="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full opacity-0 hover:opacity-100 cursor-pointer transition-opacity">
                            <Camera size={20} className="text-white"/>
                            <input
                                type="file"
                                accept="image/*"
                                className="sr-only"
                                onChange={(e) => setAvatarFile(e.target.files?.[0] || null)}
                            />
                        </label>
                    </div>
                    {avatarFile && (
                        <span className="text-sm text-primary-400">
              Новый аватар выбран
            </span>
                    )}
                </div>

                <Input
                    label="Название"
                    {...form.register('name')}
                    error={form.formState.errors.name?.message}
                />

                <div className="space-y-1.5">
                    <label className="text-sm font-medium text-neutral-300">
                        Описание
                    </label>
                    <textarea
                        {...form.register('description')}
                        className="w-full h-20 px-3 py-2 rounded-lg bg-neutral-900 border border-neutral-800 text-neutral-100 placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 resize-none"
                        placeholder="Описание чата..."
                    />
                </div>

                <Button
                    type="submit"
                    className="w-full"
                    isLoading={updateChat.isPending}
                >
                    Сохранить изменения
                </Button>
            </form>
        </div>
    );
}

// Секция участников
function MembersSection({
                            chatId,
                            members,
                            canBanUser,
                            canChangeRole,
                            onBanUser,
                            onChangeRole,
                            myRoleName,
                            currentUserId,
                        }: {
    chatId: string;
    members: ChatMemberResponse[];
    canBanUser: boolean;
    canChangeRole: boolean;
    onBanUser: (userId: string) => void;
    onChangeRole: (userId: string) => void;
    myRoleName?: string;
    currentUserId?: string;
}) {
    return (
        <div>
            <h3 className="text-sm font-medium text-neutral-300 mb-4 flex items-center gap-2">
                <Users size={16}/>
                Участники ({members.length})
            </h3>

            <ul className="space-y-2">
                {members.map((member) => (
                    <MemberItem
                        key={member.userId}
                        member={member}
                        chatId={chatId}
                        canBanUser={canBanUser}
                        canChangeRole={canChangeRole}
                        onBanUser={onBanUser}
                        onChangeRole={onChangeRole}
                        myRoleName={myRoleName}
                        currentUserId={currentUserId}
                    />
                ))}
            </ul>
        </div>
    );
}

// Один участник
function MemberItem({
                        member,
                        chatId,
                        canBanUser,
                        canChangeRole,
                        onBanUser,
                        onChangeRole,
                        myRoleName,
                        currentUserId,
                    }: {
    member: ChatMemberResponse;
    chatId: string;
    canBanUser: boolean;
    canChangeRole: boolean;
    onBanUser: (userId: string) => void;
    onChangeRole: (userId: string) => void;
    myRoleName?: string;
    currentUserId?: string;
}) {
    const {data: userInfo} = useUserBrief(member.userId, chatId);
    const isBanned = member.roleName === 'banned';
    const isOwner = member.roleName === 'owner';
    const isCurrentUser = member.userId === currentUserId;

    // Логика отображения кнопок:
    // - Нельзя менять роль самого себя
    // - Owner'а можно менять/банить только если текущий пользователь тоже owner
    // - Можно менять роль забаненных (чтобы разбанить через изменение роли)
    // - Нельзя банить уже забаненных (они уже забанены)
    const cannotChangeRole = isCurrentUser || (isOwner && myRoleName !== 'owner');
    const cannotBan = isCurrentUser || isBanned || (isOwner && myRoleName !== 'owner');

    // Показываем кнопки если есть permissions и можно выполнить действие
    const canChangeRoleForUser = canChangeRole && !cannotChangeRole;
    const canBanThisUser = canBanUser && !cannotBan;
    const showActions = canChangeRoleForUser || canBanThisUser;

    // Отладка (можно убрать позже)
    // console.log('MemberItem:', {
    //   userId: member.userId,
    //   roleName: member.roleName,
    //   isBanned,
    //   isOwner,
    //   isCurrentUser,
    //   currentUserId,
    //   myRoleName,
    //   canChangeRole,
    //   canBanUser,
    //   canChangeRoleForUser,
    //   canBanThisUser,
    //   showActions,
    // });

    return (
        <li className="flex items-center gap-3 p-3 rounded-lg bg-neutral-800/30 hover:bg-neutral-800/50 transition-colors">
            <Avatar
                file={userInfo?.avatarFile}
                fallback={userInfo?.username || '?'}
                size="sm"
            />

            <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-neutral-100 truncate">
                    {userInfo?.username || 'Загрузка...'}
                </p>
                <div className="flex items-center gap-2">
                    <Badge
                        variant={isBanned ? 'error' : isOwner ? 'primary' : 'default'}
                        size="sm"
                    >
                        {member.roleName || 'member'}
                    </Badge>
                </div>
            </div>

            {showActions && (
                <div className="flex items-center gap-1">
                    {canChangeRoleForUser && (
                        <button
                            onClick={() => onChangeRole(member.userId)}
                            className="p-2 rounded-lg text-neutral-400 hover:text-primary-400 hover:bg-neutral-700 transition-colors"
                            title="Изменить роль"
                        >
                            <UserCog size={16}/>
                        </button>
                    )}
                    {canBanThisUser && (
                        <button
                            onClick={() => onBanUser(member.userId)}
                            className="p-2 rounded-lg text-neutral-400 hover:text-error hover:bg-neutral-700 transition-colors"
                            title="Забанить"
                        >
                            <Ban size={16}/>
                        </button>
                    )}
                </div>
            )}
        </li>
    );
}

// Модалка изменения роли
function ChangeRoleModal({
                             chatId,
                             userId,
                             isOpen,
                             onClose,
                             changeUserRole,
                             myRoleName,
                         }: {
    chatId: string;
    userId: string;
    isOpen: boolean;
    onClose: () => void;
    changeUserRole: ReturnType<typeof useChangeUserRole>;
    myRoleName?: string;
}) {
    const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null);
    const {data: roles, isLoading: isLoadingRoles} = useChatRoles();
    const {data: userInfo} = useUserBrief(userId, chatId);

    // Фильтруем роли - owner может назначать только если он сам owner
    const availableRoles = roles?.filter((role) => {
        // owner роль может назначать только owner
        if (role.name === 'owner' && myRoleName !== 'owner') return false;
        // banned роль не показываем - для этого есть кнопка бана
        if (role.name === 'banned') return false;
        return true;
    });

    const handleChangeRole = () => {
        if (!selectedRoleId) return;

        changeUserRole.mutate(
            {user_id: userId, role_id: selectedRoleId},
            {
                onSuccess: () => {
                    setSelectedRoleId(null);
                    onClose();
                },
            }
        );
    };

    const handleClose = () => {
        setSelectedRoleId(null);
        onClose();
    };

    return (
        <Modal
            open={isOpen}
            onOpenChange={(open) => {
                if (!open) {
                    handleClose();
                }
            }}
        >
            <Modal.Content
                title="Изменить роль"
                description={`Выберите новую роль для ${userInfo?.username || 'пользователя'}`}
            >
                <div className="space-y-4">
                    {isLoadingRoles ? (
                        <div className="flex justify-center py-4">
                            <Loader2 className="w-6 h-6 animate-spin text-primary-500"/>
                        </div>
                    ) : (
                        <div className="space-y-2">
                            {availableRoles?.map((role) => (
                                <button
                                    key={role.id}
                                    type="button"
                                    onClick={() => setSelectedRoleId(role.id)}
                                    className={cn(
                                        'w-full p-3 rounded-lg border text-left transition-colors',
                                        selectedRoleId === role.id
                                            ? 'border-primary-500 bg-primary-500/10'
                                            : 'border-neutral-800 hover:border-neutral-700 bg-neutral-900'
                                    )}
                                >
                                    <p className="font-medium text-neutral-100">{role.name}</p>
                                    {role.permissions.length > 0 && (
                                        <div className="mt-2 flex flex-wrap gap-1">
                                            {role.permissions.map((p) => (
                                                <Badge key={p.id} variant="default" size="sm">
                                                    {p.name}
                                                </Badge>
                                            ))}
                                        </div>
                                    )}
                                </button>
                            ))}
                        </div>
                    )}

                    <div className="flex gap-3 mt-6">
                        <Button variant="secondary" className="flex-1" onClick={handleClose}>
                            Отмена
                        </Button>
                        <Button
                            className="flex-1"
                            onClick={handleChangeRole}
                            disabled={!selectedRoleId}
                            isLoading={changeUserRole.isPending}
                        >
                            Изменить
                        </Button>
                    </div>
                </div>
            </Modal.Content>
        </Modal>
    );
}

