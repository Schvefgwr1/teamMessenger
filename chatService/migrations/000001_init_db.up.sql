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
                                               name VARCHAR(100) NOT NULL
);

CREATE TABLE chat_service.chat_roles (
                                         id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                         name VARCHAR(100) NOT NULL
);

CREATE TABLE chat_service.chat_role_permissions (
                                                    chat_role_id INT REFERENCES chat_service.chat_roles(id) ON DELETE CASCADE,
                                                    chat_permission_id INT REFERENCES chat_service.chat_permissions(id) ON DELETE CASCADE,
                                                    PRIMARY KEY (chat_role_id, chat_permission_id)
);

CREATE TABLE chat_service.chat_user (
                                        chat_id UUID REFERENCES chat_service.chats(id) ON DELETE CASCADE,
                                        user_id UUID,
                                        role_id INT REFERENCES chat_service.chat_roles(id),
                                        PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE chat_service.messages (
                                       id UUID PRIMARY KEY,
                                       chat_id UUID REFERENCES chat_user(chat_id),
                                       sender_id UUID REFERENCES chat_user(user_id),
                                       content TEXT NOT NULL,
                                       updated_at TIMESTAMP,
                                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chat_service.message_files (
                                            message_id UUID REFERENCES chat_service.messages(id) ON DELETE CASCADE,
                                            file_id INT,
                                            PRIMARY KEY (message_id, file_id)
);