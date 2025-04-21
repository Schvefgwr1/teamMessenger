CREATE SCHEMA IF NOT EXISTS file_service;

CREATE TABLE IF NOT EXISTS file_service.file_types (
                                         id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                         name VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS file_service.files (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    name VARCHAR(255) NOT NULL,
                                    file_type INT REFERENCES file_service.file_types(id) ON DELETE CASCADE ,
                                    url TEXT NOT NULL,
                                    created_at TIMESTAMP DEFAULT NOW()
);