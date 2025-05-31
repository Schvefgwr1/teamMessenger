# User Service

User Service - это микросервис, отвечающий за управление пользователями, аутентификацию и систему безопасности с автоматической ротацией ключей шифрования.

## 🎯 Основные функции

- **Управление пользователями** - регистрация, обновление профилей
- **Аутентификация** - вход в систему с генерацией JWT токенов
- **Генерация и ротация RSA ключей** для подписи JWT токенов
- **Автоматическое обновление ключей** по расписанию
- **Отправка уведомлений** о входах в систему через Kafka
- **API для получения публичных ключей** другими сервисами

## 🏗️ Архитектура

### Структура проекта
```
userService/
├── cmd/
│   ├── server/
│   │   └── main.go                 # Основной сервер
│   └── migrate/
│       └── main.go                 # Миграции базы данных
├── internal/
│   ├── controllers/                # Бизнес-логика
│   │   ├── auth_controller.go
│   │   └── user_controller.go
│   ├── handlers/                   # HTTP обработчики
│   │   ├── auth_handler.go
│   │   └── user_handler.go
│   ├── models/                     # Модели данных
│   │   └── user.go
│   ├── services/                   # Внутренние сервисы
│   │   ├── key_management_service.go
│   │   ├── key_scheduler_service.go
│   │   └── notification_service.go
│   └── utils/                      # Утилиты
│       └── utils.go
├── config/
│   └── config.yaml                 # Конфигурация
├── migrations/                     # SQL миграции
├── docs/                          # Swagger документация
├── go.mod
└── README.md
```

## 🔐 Система безопасности

### Автоматическая ротация ключей

User Service реализует уникальную систему автоматической ротации RSA ключей:

#### KeyManagementService
```go
type KeyManagementService struct {
    keyProducer   *kafka.KeyProducer
    serviceName   string
    currentVersion int
}
```

**Функции:**
- Генерация новых RSA ключей (2048 бит)
- Сохранение приватного ключа локально
- Отправка публичного ключа в Kafka для других сервисов
- Версионирование ключей

#### KeySchedulerService
```go
type KeySchedulerService struct {
    keyManagement *KeyManagementService
    interval      time.Duration
    ticker        *time.Ticker
    stopChannel   chan bool
}
```

**Функции:**
- Автоматическое обновление ключей по расписанию
- Настраиваемый интервал ротации (по умолчанию 24 часа)
- Graceful остановка планировщика

### Процесс ротации ключей

1. **Планировщик** срабатывает по расписанию
2. **Генерируются** новые RSA ключи
3. **Сохраняется** приватный ключ в `private_key.pem`
4. **Отправляется** публичный ключ в Kafka топик `key_updates`
5. **API Service** получает обновление и инвалидирует все сессии
6. **Все новые токены** подписываются новым ключом

## 🔄 Интеграция с Kafka

### Отправка обновлений ключей
```go
type PublicKeyUpdate struct {
    ServiceName  string `json:"service_name"`
    PublicKeyPEM string `json:"public_key_pem"`
    KeyVersion   int    `json:"key_version"`
    Timestamp    int64  `json:"timestamp"`
}
```

### Отправка уведомлений о входе
```go
type LoginNotification struct {
    BaseNotification
    UserAgent string `json:"user_agent"`
    IPAddress string `json:"ip_address"`
}
```

## 💾 База данных

### Таблица users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Дополнительные таблицы
```sql
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
```

## 📋 API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация нового пользователя
- `POST /api/v1/auth/login` - Вход в систему
- `GET /api/v1/auth/public-key` - Получение текущего публичного ключа

### Пользователи
- `GET /api/v1/users/:id` - Получение пользователя по ID
- `PUT /api/v1/users/:id` - Обновление пользователя
- `GET /api/v1/users/:id/email` - Получение email пользователя

### Управление ключами
- `POST /api/v1/keys/regenerate` - Принудительная генерация новых ключей

## 🔄 Интеграция

User Service интегрирован с:
- **API Service** - предоставляет публичные ключи
- **Kafka** - отправляет обновления ключей и уведомления
- **Notification Service** - через Kafka уведомления
- **PostgreSQL** - для хранения данных пользователей 