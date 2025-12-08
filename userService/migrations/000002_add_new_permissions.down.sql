-- Rollback: Remove new permissions
-- ВНИМАНИЕ: Удаление permissions также удалит все связи в role_permissions (ON DELETE CASCADE)

DELETE FROM user_service.permissions 
WHERE name IN (
    'view_full_user_profile',
    'view_task_statuses',
    'manage_task_statuses'
);

