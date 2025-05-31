# Notification Service

Notification Service - это микросервис для обработки и отправки уведомлений в системе TeamMessenger. Получает сообщения через Kafka и отправляет красивые HTML email уведомления пользователям.

## 🎯 Основные функции

- **Обработка уведомлений через Kafka** - асинхронное получение сообщений
- **Отправка HTML email** - красивые, responsive шаблоны уведомлений
- **Поддержка различных типов** - задачи, чаты, входы в систему
- **SMTP интеграция** - поддержка Yandex провайдера
- **Graceful shutdown** - корректное завершение обработки сообщений

## 🏗️ Архитектура

### Структура проекта
```
notificationService/
├── cmd/
│   └── server/
│       └── main.go             # Основной сервер
├── internal/
│   ├── config/
│   │   └── config.go           # Конфигурация
│   └── services/
│       ├── email_service.go    # Сервис отправки email
│       └── kafka_consumer.go   # Kafka consumer
├── templates/                  # HTML шаблоны
│   ├── new_task.html          # Шаблон для уведомлений о задачах
│   ├── new_chat.html          # Шаблон для уведомлений о чатах
│   └── login.html             # Шаблон для уведомлений о входе
├── config/
│   └── config.yaml            # Конфигурация
├── go.mod
└── README.md
```

## 📧 Типы уведомлений

### 1. Новая задача (new_task)
**Когда отправляется:** При создании задачи с назначением исполнителя
**Содержимое:**
- Название и описание задачи
- Имя создателя задачи
- Ссылка на задачу в системе

### 2. Новый чат (new_chat)
**Когда отправляется:** При создании чата или добавлении в существующий
**Содержимое:**
- Название чата
- Имя создателя/добавившего
- Тип чата (личный/групповой)

### 3. Вход в систему (user_login)
**Когда отправляется:** При каждом входе пользователя для безопасности
**Содержимое:**
- IP адрес
- Браузер и устройство
- Время входа

## 🎨 HTML шаблоны

### Дизайн шаблонов
- **Responsive дизайн** - корректное отображение на всех устройствах
- **Современный стиль** - Material Design элементы
- **Брендинг TeamMessenger** - фирменные цвета и логотип
- **Dark/Light темы** - поддержка различных тем почтовых клиентов

### Структура шаблона
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TeamMessenger - {{.Title}}</title>
    <style>
        /* Responsive CSS стили */
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>TeamMessenger</h1>
        </div>
        <div class="content">
            <!-- Контент уведомления -->
        </div>
        <div class="footer">
            <!-- Подвал с контактами -->
        </div>
    </div>
</body>
</html>
```

## 🔄 Интеграция с Kafka

### Consumer Configuration
```go
type KafkaConsumer struct {
    consumer       sarama.ConsumerGroup
    emailService   *EmailService
    brokers        []string
    topic          string
    groupID        string
}
```

### Структуры сообщений

#### Базовая структура из Common модуля
```go
type BaseNotification struct {
    Type      string `json:"type"`
    Email     string `json:"email"`
    Timestamp int64  `json:"timestamp"`
}
```

#### Обработка сообщений
```go
func (nc *NotificationConsumer) processMessage(message *sarama.ConsumerMessage) error {
    var notification models.BaseNotification
    if err := json.Unmarshal(message.Value, &notification); err != nil {
        return err
    }

    switch notification.Type {
    case "new_task":
        return nc.handleNewTaskNotification(message.Value)
    case "new_chat":
        return nc.handleNewChatNotification(message.Value)
    case "user_login":
        return nc.handleLoginNotification(message.Value)
    default:
        return fmt.Errorf("unknown notification type: %s", notification.Type)
    }
}
```

## 🔄 Интеграция

Notification Service интегрирован с:
- **Kafka** - получение уведомлений от всех сервисов
- **SMTP сервер** - отправка email
- **Common модуль** - модели уведомлений
- **Task Service** - уведомления о новых задачах
- **Chat Service** - уведомления о чатах
- **User Service** - уведомления о входах

### Отправка уведомлений из других сервисов

```go
import (
    "common/kafka"
    "common/models"
)

// Пример из taskService
notification := models.NewTaskNotification{
    BaseNotification: models.BaseNotification{
        Type:      "new_task",
        Email:     executorEmail,
        Timestamp: time.Now().Unix(),
    },
    TaskTitle:       task.Title,
    TaskDescription: task.Description,
    CreatorName:     creatorName,
}

err := producer.SendNotification(notification)
```