# Task Service

Task Service - это микросервис для управления задачами в системе TeamMessenger. Обеспечивает создание, назначение, отслеживание задач и автоматическую отправку уведомлений исполнителям.

## 🎯 Основные функции

- **Создание задач** - с назначением исполнителей
- **Управление статусами** - отслеживание прогресса выполнения
- **Фильтрация задач** - по пользователям и статусам
- **Автоматические уведомления** - через Kafka при создании задач
- **Интеграция с пользователями** - получение данных из userService

## 🏗️ Архитектура

### Структура проекта
```
taskService/
├── cmd/
│   ├── server/
│   │   └── main.go                 # Основной сервер
│   └── migrate/
│       └── main.go                 # Миграции базы данных
├── internal/
│   ├── controllers/                # Бизнес-логика
│   │   └── task_controller.go
│   ├── handlers/                   # HTTP обработчики
│   │   └── task_handler.go
│   ├── models/                     # Модели данных
│   │   ├── task.go
│   │   └── task_status.go
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

#### Таблица task_statuses
```sql
CREATE TABLE task_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7) DEFAULT '#000000'
);

-- Предустановленные статусы
INSERT INTO task_statuses (name, description, color) VALUES 
('To Do', 'Задача создана, но работа не начата', '#6B7280'),
('In Progress', 'Задача выполняется', '#3B82F6'),
('Review', 'Задача на проверке', '#F59E0B'),
('Done', 'Задача завершена', '#10B981');
```

#### Таблица tasks
```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status_id INTEGER NOT NULL REFERENCES task_statuses(id),
    assigned_to UUID,  -- ID пользователя из userService
    created_by UUID NOT NULL,  -- ID создателя
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Таблица tasks_files
```sql
CREATE TABLE IF NOT EXISTS task_service.task_files (
     task_id INT REFERENCES task_service.tasks(id) ON DELETE CASCADE,
     file_id INT,
     PRIMARY KEY (task_id, file_id)
);
```

## 🔄 Интеграция с Kafka

### Уведомления о новых задачах
При создании задачи автоматически отправляется уведомление исполнителю:

```go
type NewTaskNotification struct {
    BaseNotification
    TaskTitle       string `json:"task_title"`
    TaskDescription string `json:"task_description"`
    CreatorName     string `json:"creator_name"`
}
```

## 📋 API Endpoints

### Задачи
- `POST /api/v1/tasks` - Создание новой задачи
- `GET /api/v1/tasks/:task_id` - Получение задачи по ID
- `PATCH /api/v1/tasks/:task_id/status/:status_id` - Обновление статуса задачи
- `GET /api/v1/users/:user_id/tasks` - Получение задач пользователя

### Статусы
- `GET /api/v1/task-statuses` - Получение всех доступных статусов

### Детали запросов

#### Создание задачи
```bash
curl -X POST http://localhost:8081/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Реализовать API для задач",
    "description": "Создать REST API для управления задачами с поддержкой CRUD операций",
    "assigned_to": "user-uuid-here"
  }'
```

#### Обновление статуса
```bash
curl -X PATCH http://localhost:8081/api/v1/tasks/task-uuid/status/2 \
  -H "Content-Type: application/json"
```

#### Получение задач пользователя
```bash
curl -X GET "http://localhost:8081/api/v1/users/user-uuid/tasks?status=1&limit=20&offset=0"
```

## 🔄 Интеграция

Task Service интегрирован с:
- **User Service** - получение email и имен пользователей
- **Notification Service** - через Kafka уведомления
- **API Service** - через HTTP клиенты
- **PostgreSQL** - для хранения задач и статусов