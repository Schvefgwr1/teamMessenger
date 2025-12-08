-- Migration: Add new permissions for user profile and task statuses management
-- Created: 2024

-- Добавляем новые permissions
INSERT INTO user_service.permissions (name, description) VALUES
    ('view_full_user_profile', 'Просмотр полного профиля пользователя (только для админов)'),
    ('view_task_statuses', 'Просмотр статусов задач (доступно всем пользователям)'),
    ('manage_task_statuses', 'Управление статусами задач: создание и удаление (только для админов)')
ON CONFLICT (name) DO NOTHING;

-- Примечание: После применения миграции нужно вручную назначить permissions ролям через админку или следующим SQL:

-- Пример назначения permissions ролям (ID ролей могут отличаться, проверьте актуальные ID):
-- 
-- Для обычных пользователей (предполагается Role ID = 1, "User"):
-- INSERT INTO user_service.role_permissions (role_id, permission_id)
-- SELECT 1, id FROM user_service.permissions WHERE name = 'view_task_statuses'
-- ON CONFLICT DO NOTHING;
--
-- Для админов (предполагается Role ID = 2, "Admin"):
-- INSERT INTO user_service.role_permissions (role_id, permission_id)
-- SELECT 2, id FROM user_service.permissions 
-- WHERE name IN ('view_full_user_profile', 'view_task_statuses', 'manage_task_statuses')
-- ON CONFLICT DO NOTHING;

