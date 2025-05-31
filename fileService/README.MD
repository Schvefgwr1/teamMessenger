# File Service

File Service - это микросервис для управления файлами в системе TeamMessenger. Обеспечивает загрузку, хранение и скачивание файлов с использованием MinIO как объектного хранилища.

## 🎯 Основные функции

- **Загрузка файлов** - поддержка различных типов файлов
- **Скачивание файлов** - безопасный доступ к загруженным файлам
- **Управление метаданными** - информация о файлах в PostgreSQL
- **Интеграция с MinIO** - надежное объектное хранилище
- **Валидация файлов** - проверка размера и типа файлов
- **URL-based доступ** - прямые ссылки на файлы

## 🏗️ Архитектура

### Структура проекта
```
fileService/
├── cmd/
│   ├── server/
│   │   └── main.go                 # Основной сервер
├── internal/
│   ├── controllers/                # Бизнес-логика
│   │   └── file_controller.go
...
│   ├── handlers/                   # HTTP обработчики
│   │   └── file_handler.go
...
│   ├── models/                     # Модели данных
│   │   └── file.go
...
│   └── services/                   # Внутренние сервисы
│       └── minio_service.go
...
├── config/
│   └── config.yaml                 # Конфигурация
├── migrations/                     # SQL миграции
├── docs/                          # Swagger документация
├── go.mod
└── README.md
```

## 💾 База данных

### Таблица files
```sql
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_name VARCHAR(255) NOT NULL,
    stored_name VARCHAR(255) NOT NULL UNIQUE,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    uploaded_by UUID NOT NULL,  -- ID пользователя
    bucket_name VARCHAR(100) NOT NULL DEFAULT 'files',
    file_path VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Таблица files_types
```sql
CREATE TABLE file_service.file_types (
     id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
     name VARCHAR(20) UNIQUE NOT NULL
);
```

## 🗄️ Интеграция с MinIO

### Конфигурация MinIO
```yaml
minio:
  host: "localhost:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket: "files"
```

## 📋 API Endpoints

### Основные операции
- `POST /api/v1/files/upload` - Загрузка файла
- `GET /api/v1/files/:file_id` - Скачивание файла
- `GET /api/v1/files/:file_id/info` - Получение информации о файле
- `DELETE /api/v1/files/:file_id` - Удаление файла (если разрешено)

### Детали запросов

#### Загрузка файла
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@/path/to/your/document.pdf" \
  -H "User-ID: user-uuid-here"
```

**Response:**
```json
{
  "id": "file-uuid",
  "original_name": "document.pdf",
  "mime_type": "application/pdf",
  "size_bytes": 1024000,
  "download_url": "/api/v1/files/file-uuid"
}
```

#### Получение информации о файле
```bash
curl -X GET http://localhost:8080/api/v1/files/file-uuid/info
```

**Response:**
```json
{
  "id": "file-uuid",
  "original_name": "document.pdf",
  "mime_type": "application/pdf",
  "size_bytes": 1024000,
  "uploaded_by": "user-uuid",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Скачивание файла
```bash
curl -X GET http://localhost:8080/api/v1/files/file-uuid \
  -o downloaded_file.pdf
```

## 🔄 Интеграция

File Service интегрирован с:
- **MinIO** - объектное хранилище файлов
- **PostgreSQL** - метаданные файлов
- **API Service** - проксирование запросов
- **Chat Service** - вложения в сообщениях
- **Task Service** - файлы в задачах
