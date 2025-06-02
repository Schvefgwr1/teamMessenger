CREATE SCHEMA IF NOT EXISTS chat_service;

-- ========================
-- SCHEMA: chat_service
-- ========================

CREATE TABLE chat_service.chats (
    id             uuid         not null primary key,
    name           varchar(255) not null,
    is_group       boolean   default false,
    description    text,
    avatar_file_id integer,
    created_at     timestamp default now()
);

CREATE TABLE chat_service.chat_permissions (
    id   integer generated always as identity primary key,
    name varchar(100) not null
);

CREATE TABLE chat_service.chat_roles (
    id   integer generated always as identity primary key,
    name varchar(100) not null
);

CREATE TABLE chat_service.chat_role_permissions (
    chat_role_id       integer not null
        references chat_service.chat_roles(id) on delete cascade,
    chat_permission_id integer not null
        constraint chat_role_permissions_permission_id_fkey
            references chat_service.chat_permissions(id) on delete cascade,
    primary key (chat_role_id, chat_permission_id)
);

CREATE TABLE chat_service.chat_user (
    chat_id uuid not null
        references chat_service.chats(id) on delete cascade,
    user_id uuid not null,
    role_id integer
        references chat_service.chat_roles(id),
    primary key (chat_id, user_id)
);

CREATE TABLE chat_service.messages (
    id         uuid not null primary key,
    chat_id    uuid,
    sender_id  uuid,
    content    text not null,
    updated_at timestamp,
    created_at timestamp default now(),
    constraint messages_sender_id_user_id_fk
        foreign key (sender_id, chat_id) references chat_service.chat_user (user_id, chat_id)
);

CREATE TABLE chat_service.message_files (
    message_id uuid    not null
        references chat_service.messages(id) on delete cascade,
    file_id    integer not null,
    primary key (message_id, file_id)
);

-- Insert default data
INSERT INTO chat_service.chat_permissions (name) VALUES ('send_message'), ('view_messages'), ('change_role'), ('ban_user'), ('edit_chat'), ('delete_chat');
INSERT INTO chat_service.chat_roles (name) VALUES ('owner'), ('banned'), ('main');
INSERT INTO chat_service.chat_role_permissions VALUES
   (1, 1),
   (1, 2),
   (1, 3),
   (1, 4),
   (1, 5),
   (1, 6),
   (2, 2),
   (3, 1),
   (3, 2)
;