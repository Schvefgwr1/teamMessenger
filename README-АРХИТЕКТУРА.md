# 🏗️ Архитектура конфигурации Team Messenger

## Проблема, которую решаем

У вас уже были локальные `.env` файлы в каждом сервисе, которые загружаются через `godotenv.Load()` в `main.go`. При контейнеризации с Docker возникли конфликты:

1. **Docker Compose** нужны переменные для настройки инфраструктуры (порты, пароли)
2. **Каждый сервис** имеет свои настройки в локальном `.env` файле
3. **Конфликт**: В контейнере нужны другие значения (например, `DB_HOST=postgres` вместо `localhost`)

## ✅ Новая двухуровневая архитектура

### 1️⃣ Уровень инфраструктуры: `compose.env`

**Назначение**: Настройки Docker Compose и инфраструктуры
**Управляет**:
- Портами сервисов и баз данных
- Паролями для PostgreSQL, Redis, MinIO
- URL подключений Kafka, Redis
- Настройками email для уведомлений

**Пример `compose.env`**:
```bash
# Порты
POSTGRES_PORT=5432
USER_SERVICE_PORT=8082
API_SERVICE_PORT=8084

# Пароли инфраструктуры  
POSTGRES_PASSWORD=secure_password
REDIS_PASSWORD=redis_pass
MINIO_ROOT_PASSWORD=minio_pass

# Email уведомления
SMTP_USERNAME=your-email@yandex.ru
SMTP_PASSWORD=your-app-password
```

### 2️⃣ Уровень сервисов: локальные `.env`

**Назначение**: Специфичные настройки каждого сервиса
**Управляет**:
- JWT секретами
- Таймаутами и лимитами
- Бизнес-логикой сервиса
- Локальными настройками разработки

**Пример `userService/.env`**:
```bash
# Локальные настройки userService
JWT_SECRET=your_jwt_secret
AUTH_TIMEOUT=30m
RATE_LIMIT=100

# Для локальной разработки
DB_HOST=localhost
DB_PORT=5432
```

## 🔄 Как это работает

### 1. В main.go каждого сервиса (без изменений!)

```go
func main() {
    // Загружаем локальный .env файл (как и раньше)
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found: %v", err)
    }
    
    // Ваш код остается прежним!
    cfg, err := config.LoadConfig("config/config.yaml")
    // ...
}
```

### 2. Docker Compose переопределяет нужные переменные

```yaml
# docker-compose.yml
user-service:
  environment:
    # Переопределяем для контейнера
    - DB_HOST=postgres              # вместо localhost из .env
    - DB_USER=${POSTGRES_USER}      # из compose.env  
    - DB_PASSWORD=${POSTGRES_PASSWORD}  # из compose.env
    - KAFKA_BROKERS=kafka:9092      # вместо localhost:9092
```

### 3. Приоритет переменных (от высшего к низшему)

1. **Docker environment** (docker-compose.yml) - переменные инфраструктуры
2. **Локальный .env файл** - через `godotenv.Load()`
3. **Системные переменные** - если установлены в системе

### 4. Результат 🎯

- **В локальной разработке**: Сервис использует `localhost` из `.env`
- **В Docker**: Сервис использует `postgres` из docker-compose environment
- **Ваш код**: Остается без изменений!

## 📁 Структура файлов

```
teamMessenger/
├── compose.env                   # 🐳 Настройки Docker инфраструктуры
├── compose.env.example           # 📋 Шаблон
├── compose.development           # 🛠️ Для разработки
├── docker-compose.yml            # 🐳 Использует compose.env
├── userService/
│   ├── .env                      # ⚙️ Локальные настройки userService
│   ├── env.example               # 📋 Шаблон для userService
│   ├── main.go                   # 🔄 godotenv.Load() (без изменений!)
│   └── Dockerfile                # 📦 Копирует .env в контейнер
├── apiService/
│   ├── .env                      # ⚙️ Локальные настройки apiService  
│   ├── env.example               # 📋 Шаблон для apiService
│   └── main.go                   # 🔄 godotenv.Load() (без изменений!)
└── ... (остальные сервисы)
```

## 🚀 Быстрый старт

### 1. Настройка всех файлов конфигурации
```bash
make setup-env
# Создает compose.env и все локальные .env файлы
```

### 2. Для разработки
```bash
make env-dev
# Применяет настройки разработки для Docker инфраструктуры
```

### 3. Запуск системы
```bash
make up
# Запускает всю систему с правильными настройками
```

## 🔧 Управление настройками

### Docker инфраструктура (порты, пароли)
```bash
# Отредактируйте compose.env
nano compose.env

# Перезапустите инфраструктуру
make restart
```

### Настройки конкретного сервиса
```bash
# Отредактируйте локальный .env
nano userService/.env

# Пересоберите только этот сервис
docker-compose build user-service
docker-compose restart user-service
```

### Проверка конфигурации
```bash
# Показать все настройки
make show-config

# Проверить переменные в контейнере
docker-compose exec user-service env | grep -E "DB_|APP_"
```

## 🎯 Преимущества новой архитектуры

### ✅ Совместимость
- Ваши существующие `.env` файлы работают без изменений
- `godotenv.Load()` в `main.go` остается как есть
- Нет необходимости менять код сервисов

### ✅ Гибкость
- Инфраструктура настраивается отдельно от логики сервисов
- Можно легко менять порты, пароли, URL без изменения кода
- Разные настройки для development/staging/production

### ✅ Безопасность
- Пароли инфраструктуры изолированы в `compose.env`
- Секреты сервисов остаются в их локальных `.env`
- Можно использовать разные уровни доступа

### ✅ Простота развертывания
- Один команда `make up` запускает всю систему
- Автоматическое переопределение переменных для Docker
- Поддержка разных окружений через разные `compose.env`

## 🔍 Отладка

### Проверить загрузку .env файлов
```bash
# Логи загрузки локального .env
docker-compose logs user-service | grep -i "env file"

# Переменные в контейнере
docker-compose exec user-service env | sort
```

### Проверить приоритет переменных
```bash
# Сравнить локальный .env и контейнер
echo "=== Локальный .env ==="
cat userService/.env

echo "=== В контейнере ==="  
docker-compose exec user-service env | grep -E "DB_|APP_"
```

### Конфликты переменных
Если есть конфликты, проверьте:
1. **docker-compose.yml environment** - имеет высший приоритет
2. **Локальный .env** - через `godotenv.Load()`
3. **Системные переменные** - самый низкий приоритет

## 📚 Команды Make

```bash
# Настройка
make setup-env          # Создать все файлы конфигурации
make setup-compose      # Только Docker инфраструктура
make setup-services     # Только локальные .env файлы

# Управление
make up                 # Запустить всю систему
make down               # Остановить систему
make restart            # Перезапустить систему
make show-config        # Показать все настройки

# Разработка
make env-dev            # Настройки для разработки
make env-prod           # Рекомендации для продакшена
```

## 🎉 Результат

Теперь у вас есть гибкая система конфигурации, которая:
- Сохраняет вашу существующую архитектуру с `.env` файлами
- Добавляет мощные возможности Docker Compose
- Обеспечивает правильную работу в контейнерах
- Позволяет легко управлять разными окружениями

Ваш код остается без изменений, но получает все преимущества контейнеризации! 🚀 