# Common Module

Common Module - это общий модуль, содержащий переиспользуемые компоненты, модели данных, конфигурации и утилиты для всех микросервисов в системе TeamMessenger.

## 🎯 Назначение

Модуль предоставляет:
- **Общие модели данных** - структуры для уведомлений, конфигурации
- **Kafka интеграцию** - producers, consumers и конфигурацию
- **Redis клиенты** - подключение и операции с Redis
- **HTTP клиенты** - для взаимодействия между сервисами
- **Конфигурацию** - загрузка и парсинг YAML конфигураций
- **Database утилиты** - подключение к PostgreSQL
- **MinIO интеграцию** - работа с объектным хранилищем

## 🏗️ Структура модуля

```
common/
├── config/
│   └── config.go               # Загрузка и структуры конфигурации
├── contracts/
│   ├── user_client.go         # Интерфейсы для HTTP клиентов
│   ├── chat_client.go
│   ├── task_client.go
│   └── file_client.go
├── db/
│   └── postgres.go            # Подключение к PostgreSQL
├── http_clients/
│   ├── user_client.go         # HTTP клиенты для сервисов
│   ├── chat_client.go
│   ├── task_client.go
│   └── file_client.go
├── kafka/
│   ├── config.go              # Конфигурация Kafka
│   ├── producer.go            # Общий producer
│   └── key_producer.go        # Producer для обновлений ключей
├── minio/
│   └── client.go              # MinIO клиент
├── models/
│   ├── notification.go        # Модели уведомлений
│   └── key_update.go          # Модели обновления ключей
├── redis/
│   └── client.go              # Redis клиент
├── go.mod
└── README.md
```

## 📋 Основные компоненты

### 1. Конфигурация (config/)

#### Структуры конфигурации
```go
type Config struct {
    Database struct {
        Host     string `yaml:"host"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
        Port     int    `yaml:"port"`
    } `yaml:"db"`
    MinIO MinIO       `yaml:"minio"`
    Redis Redis       `yaml:"redis"`
    App   AppConfig   `yaml:"app"`
    Keys  KeysConfig  `yaml:"keys"`
    Kafka KafkaConfig `yaml:"kafka"`
}
```

### 2. Модели уведомлений (models/)

#### Базовая структура
```go
type BaseNotification struct {
    Type      string `json:"type"`
    Email     string `json:"email"`
    Timestamp int64  `json:"timestamp"`
}
```

#### Типы уведомлений
```go
// Уведомления о новых задачах
type NewTaskNotification struct {
    BaseNotification
    TaskTitle       string `json:"task_title"`
    TaskDescription string `json:"task_description"`
    CreatorName     string `json:"creator_name"`
}

// Уведомления о новых чатах
type NewChatNotification struct {
    BaseNotification
    ChatName    string `json:"chat_name"`
    CreatorName string `json:"creator_name"`
    IsGroup     bool   `json:"is_group"`
}

// Уведомления о входе в систему
type LoginNotification struct {
    BaseNotification
    UserAgent string `json:"user_agent"`
    IPAddress string `json:"ip_address"`
}
```

### 3. Kafka интеграция (kafka/)

#### Конфигурация
```go
func GetKafkaBrokers() []string {
    brokers := os.Getenv("KAFKA_BROKERS")
    if brokers == "" {
        return []string{"localhost:9092"}
    }
    return strings.Split(brokers, ",")
}

func GetNotificationsTopic() string {
    topic := os.Getenv("NOTIFICATIONS_TOPIC")
    if topic == "" {
        return "notifications"
    }
    return topic
}

func GetKeyUpdatesTopic() string {
    topic := os.Getenv("KEY_UPDATES_TOPIC")
    if topic == "" {
        return "key_updates"
    }
    return topic
}
```

#### Producer для уведомлений
```go
type Producer struct {
    producer sarama.SyncProducer
    topic    string
}

func (p *Producer) SendNotification(notification interface{}) error {
    data, err := json.Marshal(notification)
    if err != nil {
        return err
    }

    message := &sarama.ProducerMessage{
        Topic: p.topic,
        Value: sarama.StringEncoder(data),
    }

    _, _, err = p.producer.SendMessage(message)
    return err
}
```

#### Producer для обновлений ключей
```go
type KeyProducer struct {
    producer sarama.SyncProducer
    topic    string
}

func (kp *KeyProducer) SendKeyUpdate(keyUpdate models.PublicKeyUpdate) error {
    data, err := json.Marshal(keyUpdate)
    if err != nil {
        return err
    }

    message := &sarama.ProducerMessage{
        Topic: kp.topic,
        Value: sarama.StringEncoder(data),
    }

    _, _, err = kp.producer.SendMessage(message)
    return err
}
```

### 4. Redis клиент (redis/)

```go
func NewRedisClient(config *config.Redis) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
        Password: config.Password,
        DB:       config.DB,
    })
}
```

### 5. PostgreSQL подключение (db/)

```go
func ConnectPostgres(config *config.Config) (*sql.DB, error) {
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        config.Database.Host,
        config.Database.Port,
        config.Database.User,
        config.Database.Password,
        config.Database.Name,
    )
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### 6. HTTP клиенты (http_clients/)

#### User Client
```go
type UserClient struct {
    baseURL string
    client  *http.Client
}

func (uc *UserClient) GetUserEmail(userID uuid.UUID) (string, error) {
    resp, err := uc.client.Get(fmt.Sprintf("%s/api/v1/users/%s/email", uc.baseURL, userID))
    // ... обработка ответа
}

func (uc *UserClient) GetUserName(userID uuid.UUID) (string, error) {
    resp, err := uc.client.Get(fmt.Sprintf("%s/api/v1/users/%s/name", uc.baseURL, userID))
    // ... обработка ответа
}
```

### 7. MinIO клиент (minio/)

```go
func NewMinIOClient(config *config.MinIO) (*minio.Client, error) {
    return minio.New(config.Host, &minio.Options{
        Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
        Secure: false, // Настраивается через конфигурацию
    })
}
```

## 🔄 Интеграция

Common модуль используется во всех микросервисах:
- **apiService** - HTTP клиенты, Redis, конфигурация
- **userService** - Kafka producers, модели, конфигурация
- **chatService** - Kafka producers, HTTP клиенты, база данных
- **taskService** - HTTP клиенты, Kafka producers, модели
- **fileService** - MinIO клиент, база данных, конфигурация
- **notificationService** - Модели уведомлений, Kafka consumers 