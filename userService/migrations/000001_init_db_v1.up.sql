CREATE SCHEMA IF NOT EXISTS user_service;

-- ========================
-- SCHEMA: user_service
-- ========================

CREATE TABLE IF NOT EXISTS user_service.roles (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    name VARCHAR(100) NOT NULL UNIQUE,
                                    description TEXT
);

CREATE TABLE IF NOT EXISTS user_service.permissions (
                                          id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                          name VARCHAR(100) NOT NULL UNIQUE,
                                          description TEXT
);

CREATE TABLE IF NOT EXISTS user_service.role_permissions (
                                               role_id INT REFERENCES user_service.roles(id) ON DELETE CASCADE,
                                               permission_id INT REFERENCES user_service.permissions(id) ON DELETE CASCADE,
                                               PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS user_service.users (
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