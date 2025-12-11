CREATE SCHEMA IF NOT EXISTS task_service;

-- ========================
-- SCHEMA: task_service
-- ========================

CREATE TABLE IF NOT EXISTS task_service.task_statuses (
                                            id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                            name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS task_service.tasks (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    title VARCHAR(255) NOT NULL,
                                    description TEXT,
                                    status INT REFERENCES task_service.task_statuses,
                                    creator_id UUID,
                                    executor_id UUID,
                                    chat_id UUID,
                                    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS task_service.task_files (
                                         task_id INT REFERENCES task_service.tasks(id) ON DELETE CASCADE,
                                         file_id INT,
                                         PRIMARY KEY (task_id, file_id)
);

INSERT INTO task_service.task_statuses (name) VALUES
    ('canseled'),
    ('created')
;
