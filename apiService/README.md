# API Service

API Service - это точка входа (API Gateway) в систему TeamMessenger, которая обеспечивает аутентификацию, авторизацию и маршрутизацию запросов к другим микросервисам.

## 🎯 Основные функции

- **API Gateway** - единая точка входа для всех клиентских запросов
- **JWT аутентификация** с динамической ротацией ключей
- **Управление сессиями** через Redis
- **Маршрутизация запросов** к соответствующим микросервисам
- **Автоматическая инвалидация сессий** при обновлении ключей шифрования
- **Кэширование** часто используемых данных

## 🏗️ Архитектура

### Структура проекта
```
apiService/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа приложения
├── internal/
│   ├── controllers/                # Бизнес-логика контроллеров
│   │   ├── auth_controller.go
│   │   ├── chat_controller.go
│   │   ├── task_controller.go
│   │   └── user_controller.go
│   ├── handlers/                   # HTTP обработчики
│   │   ├── auth_handler.go
│   │   ├── chat_handler.go
│   │   ├── task_handler.go
│   │   └── user_handler.go
│   ├── http_clients/              # HTTP клиенты для других сервисов
│   │   ├── chat_client.go
│   │   ├── file_client.go
│   │   ├── task_client.go
│   │   └── user_client.go
│   ├── middlewares/               # Middleware для аутентификации
│   │   └── jwt_check_middleware.go
│   ├── routes/                    # Определение маршрутов
│   │   ├── auth_routes.go
│   │   ├── chat_routes.go
│   │   ├── task_routes.go
│   │   └── user_routes.go
│   └── services/                  # Внутренние сервисы
│       ├── cache_service.go
│       ├── key_service.go
│       ├── key_update_consumer.go
│       ├── public_key_manager.go
│       └── session_service.go
├── config/
│   └── config.yaml               # Конфигурация сервиса
├── docs/                         # Swagger документация
├── go.mod
└── README.md
```

## 🔐 Система безопасности

### JWT с динамической ротацией ключей

API Service реализует уникальную систему безопасности с автоматической ротацией ключей:

1. **PublicKeyManager** - thread-safe менеджер публичных ключей
2. **KeyUpdateConsumer** - Kafka consumer для получения обновлений ключей
3. **Автоматическая инвалидация сессий** при смене ключей

### Процесс аутентификации

1. Клиент отправляет JWT токен в заголовке `Authorization: Bearer <token>`
2. Middleware извлекает текущий публичный ключ из PublicKeyManager
3. Проверяется подпись токена
4. Проверяется валидность сессии в Redis
5. При успехе запрос проксируется к соответствующему сервису

## 🔄 Интеграция с Kafka

### Прием обновлений ключей

```go
type KeyUpdateConsumer struct {
    consumer         sarama.ConsumerGroup
    publicKeyManager *PublicKeyManager
    sessionService   *SessionService
    redisClient      *redis.Client
}
```

Потребляет сообщения из топика `key_updates` и автоматически:
- Обновляет публичный ключ в PublicKeyManager
- Инвалидирует все активные сессии в Redis
- Логирует процесс обновления

## 📋 API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Вход в систему
- `POST /api/v1/auth/logout` - Выход из системы (требует аутентификации)

### Пользователи
- `GET /api/v1/users/me` - Получение профиля текущего пользователя
- `PUT /api/v1/users/me` - Обновление профиля

### Чаты
- `GET /api/v1/chats/:user_id` - Получение чатов пользователя
- `POST /api/v1/chats` - Создание нового чата
- `POST /api/v1/chats/messages/:chat_id` - Отправка сообщения
- `GET /api/v1/chats/messages/:chat_id` - Получение сообщений чата
- `GET /api/v1/chats/search/:chat_id` - Поиск по сообщениям

### Задачи
- `POST /api/v1/tasks` - Создание задачи
- `PATCH /api/v1/tasks/:task_id/status/:status_id` - Обновление статуса задачи
- `GET /api/v1/tasks/:task_id` - Получение задачи по ID
- `GET /api/v1/users/:user_id/tasks` - Получение задач пользователя

## 🔧 Конфигурация

### Redis
Используется для:
- Хранения активных сессий
- Кэширования данных пользователей и чатов
- Быстрого доступа к часто запрашиваемой информации

### HTTP клиенты
API Service взаимодействует с другими сервисами через HTTP:
- **UserClient** - работа с пользователями
- **ChatClient** - работа с чатами
- **TaskClient** - работа с задачами
- **FileClient** - работа с файлами

## 📊 Мониторинг

### Логирование
- Все HTTP запросы логируются с деталями
- Ошибки аутентификации отслеживаются
- Обновления ключей логируются
- Процесс инвалидации сессий отслеживается

### Graceful Shutdown
Сервис корректно завершает работу:
- Останавливает Kafka consumer
- Закрывает соединения с Redis
- Завершает обработку текущих запросов

## 🔄 Интеграция

API Service тесно интегрирован с другими компонентами системы:
- **UserService** - для получения начального публичного ключа
- **Kafka** - для получения обновлений ключей
- **Redis** - для управления сессиями
- **Все микросервисы** - для проксирования запросов