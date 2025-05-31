# Chat Service

Chat Service - это микросервис для управления чатами и обмена сообщениями в системе TeamMessenger. Поддерживает групповые и личные чаты с уведомлениями о новых участниках.

## 🎯 Основные функции

- **Создание чатов** - групповые и личные чаты
- **Обмен сообщениями** - отправка и получение сообщений
- **Поиск по сообщениям** - полнотекстовый поиск
- **Управление участниками** - добавление пользователей в чаты
- **Уведомления** - автоматическая отправка уведомлений через Kafka
- **Интеграция с файлами** - поддержка вложений

## 🏗️ Архитектура

### Структура проекта
```
chatService/
├── cmd/
│   ├── server/
│   │   └── main.go                 # Основной сервер
│   └── migrate/
│       └── main.go                 # Миграции базы данных
├── internal/
│   ├── controllers/                # Бизнес-логика
│   │   └── chat_controller.go
│   ├── handlers/                   # HTTP обработчики
│   │   └── chat_handler.go
│   ├── models/                     # Модели данных
│   │   ├── chat.go
│   │   ├── message.go
│   │   └── chat_user.go
│   └── services/                   # Внутренние сервисы
│       └── notification_service.go
├── config/
│   └── config.yaml                 # Конфигурация
├── migrations/                     # SQL миграции
├── docs/                          # Swagger документация
├── go.mod
└── README.md
```

## 💾 База данных

### Основные таблицы

#### Таблица chats
```sql
CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_group BOOLEAN NOT NULL DEFAULT false,
    created_by UUID NOT NULL,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Таблица messages
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    content TEXT NOT NULL,
    file_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Таблица chat_users
```sql
CREATE TABLE chat_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔄 Интеграция с Kafka

### Уведомления о новых чатах
При создании чата или добавлении пользователей автоматически отправляются уведомления:

```go
type NewChatNotification struct {
    BaseNotification
    ChatName    string `json:"chat_name"`
    CreatorName string `json:"creator_name"`
    IsGroup     bool   `json:"is_group"`
}
```

## 📋 API Endpoints

### Чаты
- `GET /api/v1/chats/:user_id` - Получение всех чатов пользователя
- `POST /api/v1/chats` - Создание нового чата
- `POST /api/v1/chats/:chat_id/users` - Добавление пользователей в чат

### Сообщения
- `POST /api/v1/messages/:chat_id` - Отправка сообщения в чат
- `GET /api/v1/messages/:chat_id` - Получение сообщений чата (с пагинацией)
- `GET /api/v1/messages/:chat_id/search` - Поиск сообщений в чате

### Детали запросов

#### Создание чата
```bash
curl -X POST http://localhost:8083/api/v1/chats \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Команда разработки",
    "description": "Обсуждение проекта",
    "is_group": true,
    "user_ids": ["user1-uuid", "user2-uuid"]
  }'
```

#### Отправка сообщения
```bash
curl -X POST http://localhost:8083/api/v1/messages/chat-uuid \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Привет всем!",
    "file_url": "https://example.com/file.jpg"
  }'
```

#### Поиск сообщений
```bash
curl -X GET "http://localhost:8083/api/v1/messages/chat-uuid/search?query=привет&limit=20&offset=0"
```

## 🔧 Функциональность

### Типы чатов
1. **Личные чаты** (`is_group: false`) - между двумя пользователями
2. **Групповые чаты** (`is_group: true`) - множество участников

### Поиск сообщений
- **Полнотекстовый поиск** по содержимому сообщений
- **Пагинация** результатов поиска
- **Фильтрация** по дате создания

### Управление файлами
- Поддержка **вложений** в сообщениях
- Интеграция с **File Service** для загрузки файлов
- Сохранение **URL файлов** в базе данных

## 📊 Мониторинг

### Логирование
- Создание новых чатов
- Отправка сообщений
- Добавление участников
- Отправка уведомлений
- Ошибки базы данных

### Производительность
- **Индексы** на chat_id и sender_id для быстрого поиска
- **Пагинация** для больших списков сообщений
- **Кэширование** часто запрашиваемых данных

## 🔄 Интеграция

Chat Service интегрирован с:
- **User Service** - получение информации о пользователях
- **File Service** - загрузка и получение файлов
- **Notification Service** - через Kafka уведомления
- **API Service** - через HTTP клиенты