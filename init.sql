-- ========================
-- Создание схем
-- ========================
CREATE SCHEMA user_service;
CREATE SCHEMA chat_service;
CREATE SCHEMA task_service;
CREATE SCHEMA file_service;
CREATE SCHEMA auth_service;

-- ========================
-- SCHEMA: file_service
-- ========================

CREATE TABLE file_service.file_types (
                                         id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                         name VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE file_service.files (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    name VARCHAR(255) NOT NULL,
                                    file_type INT REFERENCES file_service.file_types(id) ON DELETE CASCADE ,
                                    url TEXT NOT NULL,
                                    created_at TIMESTAMP DEFAULT NOW()
);

-- ========================
-- SCHEMA: user_service
-- ========================

CREATE TABLE user_service.roles (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    name VARCHAR(100) NOT NULL UNIQUE,
                                    description TEXT
);

CREATE TABLE user_service.permissions (
                                          id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                          name VARCHAR(100) NOT NULL UNIQUE,
                                          description TEXT
);

CREATE TABLE user_service.role_permissions (
                                               role_id INT REFERENCES user_service.roles(id) ON DELETE CASCADE,
                                               permission_id INT REFERENCES user_service.permissions(id) ON DELETE CASCADE,
                                               PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_service.users (
                                    id UUID PRIMARY KEY,
                                    username VARCHAR(255) UNIQUE NOT NULL,
                                    email VARCHAR(255) UNIQUE NOT NULL,
                                    password_hash VARCHAR(255) NOT NULL,
                                    description TEXT,
                                    gender VARCHAR(10),
                                    age INTEGER,
                                    avatar_file_id INT,
                                    role_id INT REFERENCES user_service.roles(id),
                                    created_at TIMESTAMP DEFAULT NOW(),
                                    updated_at TIMESTAMP DEFAULT NOW()
);

-- ========================
-- SCHEMA: chat_service
-- ========================

CREATE TABLE chat_service.chats (
                                    id UUID PRIMARY KEY,
                                    name VARCHAR(255) NOT NULL,
                                    is_group BOOLEAN DEFAULT FALSE,
                                    description TEXT,
                                    avatar_file_id INT,
                                    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chat_service.chat_permissions (
                                         id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                         chat_id UUID REFERENCES chat_service.chats(id) ON DELETE CASCADE,
                                         name VARCHAR(100) NOT NULL
);

CREATE TABLE chat_service.chat_roles (
                                         id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                         chat_id UUID REFERENCES chat_service.chats(id) ON DELETE CASCADE,
                                         name VARCHAR(100) NOT NULL
);

CREATE TABLE chat_service.chat_role_permissions (
                                                    chat_role_id INT REFERENCES chat_service.chat_roles(id) ON DELETE CASCADE,
                                                    permission_id INT REFERENCES chat_service.chat_permissions(id) ON DELETE CASCADE,
                                                    PRIMARY KEY (chat_role_id, permission_id)
);

CREATE TABLE chat_service.chat_user (
                                        chat_id UUID REFERENCES chat_service.chats(id) ON DELETE CASCADE,
                                        user_id UUID REFERENCES user_service.users(id) ON DELETE CASCADE,
                                        role_id INT REFERENCES chat_service.chat_roles(id),
                                        PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE chat_service.messages (
                                       id UUID PRIMARY KEY,
                                       chat_id UUID REFERENCES chat_service.chats(id) ON DELETE CASCADE,
                                       sender_id UUID REFERENCES user_service.users(id) ON DELETE SET NULL,
                                       content TEXT NOT NULL,
                                       updated_at TIMESTAMP,
                                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chat_service.message_files (
                                            message_id UUID REFERENCES chat_service.messages(id) ON DELETE CASCADE,
                                            file_id INT REFERENCES file_service.files(id) ON DELETE CASCADE,
                                            PRIMARY KEY (message_id, file_id)
);

-- ========================
-- SCHEMA: task_service
-- ========================

CREATE TABLE task_service.task_statuses (
                                            id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                            name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE task_service.tasks (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    title VARCHAR(255) NOT NULL,
                                    description TEXT,
                                    status INT REFERENCES task_service.task_statuses,
                                    creator_id UUID REFERENCES user_service.users(id) ON DELETE SET NULL,
                                    executor_id UUID REFERENCES user_service.users(id) ON DELETE SET NULL,
                                    chat_id UUID REFERENCES chat_service.chats(id) ON DELETE SET NULL,
                                    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE task_service.task_files (
                                         task_id INT REFERENCES task_service.tasks(id) ON DELETE CASCADE,
                                         file_id INT REFERENCES file_service.files(id) ON DELETE CASCADE,
                                         PRIMARY KEY (task_id, file_id)
);

-- ========================
-- SCHEMA: auth_service
-- ========================

CREATE TABLE auth_service.service_roles (
                                            id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                            name VARCHAR(100) UNIQUE NOT NULL,
                                            description TEXT
);

CREATE TABLE auth_service.service_permissions (
                                                  id INT PRIMARY KEY,
                                                  name VARCHAR(100) UNIQUE NOT NULL,
                                                  description TEXT
);

CREATE TABLE auth_service.service_role_permissions (
                                                       role_id INT REFERENCES auth_service.service_roles(id) ON DELETE CASCADE,
                                                       permission_id INT REFERENCES auth_service.service_permissions(id) ON DELETE CASCADE,
                                                       PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE auth_service.services (
                                       id UUID PRIMARY KEY,
                                       name VARCHAR(255) UNIQUE NOT NULL,
                                       description TEXT,
                                       role INT REFERENCES auth_service.service_roles(id),
                                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auth_service.service_tokens (
                                             id UUID PRIMARY KEY,
                                             service_id UUID REFERENCES auth_service.services(id) ON DELETE CASCADE,
                                             token VARCHAR(1000) NOT NULL,
                                             expires_at TIMESTAMP NOT NULL
);
