# Notification Service

Микросервис для обработки уведомлений TeamMessenger через Kafka и отправки email уведомлений.

## Описание

NotificationService обрабатывает три типа уведомлений:

1. **Новая задача** - уведомление о назначении новой задачи
2. **Новый чат** - уведомление о создании нового чата или добавлении в группу
3. **Вход в систему** - уведомление о входе в аккаунт для безопасности

## Структура проекта

```
notificationService/
├── cmd/server/main.go          # Основной файл приложения
├── internal/
│   ├── config/config.go        # Конфигурация
│   ├── models/notification.go  # Модели уведомлений
│   └── services/
│       ├── email_service.go    # Сервис отправки email
│       └── kafka_consumer.go   # Kafka consumer
├── templates/                  # HTML шаблоны email
│   ├── new_task.html          # Шаблон для уведомления о задаче
│   ├── new_chat.html          # Шаблон для уведомления о чате
│   └── login.html             # Шаблон для уведомления о входе
├── config/config.yaml         # Конфигурационный файл
├── go.mod                     # Go модуль
└── README.md                  # Документация
```

## Конфигурация

### Переменные окружения

```bash
# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=notification-service
KAFKA_TOPIC_NOTIFICATIONS=notifications

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@teammessenger.com
FROM_NAME=TeamMessenger
TEMPLATE_PATH=./templates

# Конфигурация (опционально)
CONFIG_PATH=./config/config.yaml
```

### Настройка Gmail SMTP

1. Включите двухфакторную аутентификацию в Gmail
2. Создайте пароль приложения: https://support.google.com/accounts/answer/185833
3. Используйте пароль приложения в переменной `SMTP_PASSWORD`

## Структуры сообщений Kafka

### Базовая структура

```json
{
  "type": "new_task|new_chat|user_login",
  "payload": {
    // Конкретные данные уведомления
  }
}
```

### Новая задача

```json
{
  "type": "new_task",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "new_task",
    "email": "user@example.com",
    "created_at": "2023-12-01T12:00:00Z",
    "task_id": 123,
    "task_title": "Новая важная задача",
    "creator_name": "Иван Иванов",
    "executor_id": "550e8400-e29b-41d4-a716-446655440001",
    "due_date": "2023-12-05T18:00:00Z",
    "priority": "high"
  }
}
```

### Новый чат

```json
{
  "type": "new_chat",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "new_chat",
    "email": "user@example.com",
    "created_at": "2023-12-01T12:00:00Z",
    "chat_id": "550e8400-e29b-41d4-a716-446655440002",
    "chat_name": "Проектная группа",
    "creator_name": "Анна Петрова",
    "is_group": true,
    "description": "Обсуждение текущего проекта"
  }
}
```

### Вход в систему

```json
{
  "type": "user_login",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "user_login",
    "email": "user@example.com",
    "created_at": "2023-12-01T12:00:00Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440003",
    "username": "johndoe",
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "login_time": "2023-12-01T12:00:00Z"
  }
}
```

## Установка и запуск

### Требования

- Go 1.23+
- Apache Kafka
- SMTP сервер (Gmail, SendGrid, или другой)

### Установка зависимостей

```bash
cd notificationService
go mod download
```

### Запуск

```bash
# С конфигурационным файлом
CONFIG_PATH=./config/config.yaml go run cmd/server/main.go

# Только с переменными окружения
SMTP_USERNAME=your-email@gmail.com \
SMTP_PASSWORD=your-app-password \
go run cmd/server/main.go
```

### Сборка

```bash
go build -o notification-service cmd/server/main.go
./notification-service
```

## Интеграция с другими сервисами

Для отправки уведомлений из других сервисов, публикуйте сообщения в Kafka топик `notifications`:

### Пример отправки из Go

```go
import (
    "encoding/json"
    "github.com/IBM/sarama"
)

type KafkaMessage struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// Отправка уведомления о новой задаче
func sendTaskNotification(producer sarama.SyncProducer, taskData NewTaskNotification) error {
    message := KafkaMessage{
        Type:    "new_task",
        Payload: taskData,
    }
    
    messageBytes, _ := json.Marshal(message)
    
    _, _, err := producer.SendMessage(&sarama.ProducerMessage{
        Topic: "notifications",
        Value: sarama.StringEncoder(messageBytes),
    })
    
    return err
}
```

## Troubleshooting

### Проблемы с подключением к Kafka

- Убедитесь, что Kafka запущен на указанном адресе
- Проверьте, что топик `notifications` существует
- Проверьте настройки безопасности Kafka

### Проблемы с отправкой email

- Проверьте SMTP настройки
- Убедитесь, что используете правильный пароль приложения для Gmail
- Проверьте, что `from_email` разрешен для отправки через SMTP сервер

### Логи

Сервис выводит подробные логи в стандартный поток вывода. Для отладки можно использовать:

```bash
go run cmd/server/main.go 2>&1 | tee notification.log
```

## Мониторинг

Сервис логирует:
- Подключение к Kafka
- Получение сообщений из Kafka
- Успешную и неуспешную отправку email
- Ошибки обработки сообщений

Рекомендуется настроить мониторинг логов для отслеживания работоспособности сервиса. 