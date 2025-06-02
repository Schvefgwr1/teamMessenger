# Team Messenger - Развертывание с Docker

Это руководство описывает как развернуть Team Messenger с использованием Docker и docker-compose.

## Архитектура

Система состоит из следующих компонентов:

### Микросервисы
- **API Service** (8084) - Gateway для всех запросов
- **User Service** (8082) - Управление пользователями и аутентификация
- **File Service** (8080) - Загрузка и управление файлами
- **Chat Service** (8083) - Обмен сообщениями
- **Task Service** (8081) - Управление задачами
- **Notification Service** (8085) - Email уведомления

### Инфраструктура
- **PostgreSQL** (5432) - Основная база данных
- **Redis** (6379) - Кеш и сессии
- **Kafka** (9092) - Очереди сообщений
- **MinIO** (9000/9001) - S3-совместимое хранилище файлов
- **Nginx** (8090) - Прокси для доступа к файлам
- **Zookeeper** (2181) - Координация Kafka

## 🏗️ Новая архитектура конфигурации

### Двухуровневая система переменных окружения

**1. compose.env** - Настройки Docker инфраструктуры:
- Порты сервисов и инфраструктуры
- Пароли для PostgreSQL, Redis, MinIO
- Настройки Kafka и других внешних сервисов

**2. Локальные .env файлы** - В каждом сервисе:
- Специфичные настройки каждого сервиса
- Загружаются через `godotenv.Load()` в main.go
- Переопределяются переменными из docker-compose при работе в контейнерах

### Приоритет переменных (от высшего к низшему):
1. **Environment из docker-compose.yml** - переменные инфраструктуры (DB_HOST, KAFKA_BROKERS)
2. **Локальный .env файл сервиса** - загружается через godotenv.Load()
3. **Системные переменные окружения** - если установлены в системе

### Преимущества новой архитектуры:
- 🔄 **Совместимость**: Ваши существующие .env файлы продолжают работать
- 🐳 **Docker-ready**: Переменные инфраструктуры переопределяются для контейнеров
- 🔧 **Гибкость**: Можно настраивать отдельно инфраструктуру и логику сервисов
- 🛡️ **Безопасность**: Пароли инфраструктуры изолированы от кода сервисов

## 🔄 Система миграций

**Автоматические миграции**: Каждый микросервис автоматически выполняет свои миграции при запуске. Миграции находятся в папке `migrations/` каждого сервиса и имеют формат `000001_init_db.up.sql`.

**Особенности**:
- Миграции выполняются только один раз
- Каждый сервис имеет свою схему в БД
- Отслеживание применённых миграций в таблице `schema_migrations`
- Автоматическое создание схем для сервисов

## Быстрый старт

### Предварительные требования

- Docker
- Docker Compose 
- Make (опционально, для удобства)

### 1. Клонирование и подготовка

```bash
git clone <repository-url>
cd teamMessenger
```

### 2. Настройка переменных окружения

#### Автоматическая настройка (рекомендуется):
```bash
make setup-env
```

Это создаст:
- `compose.env` - для Docker инфраструктуры
- `.env` файлы во всех сервисах - для локальной логики

#### Раздельная настройка:
```bash
# Только Docker инфраструктура
make setup-compose

# Только .env файлы сервисов
make setup-services

# Для разработки
make env-dev
```

### 3. Основные настройки

#### В compose.env (инфраструктура):
```bash
# Email для уведомлений
SMTP_USERNAME=your-email@yandex.ru
SMTP_PASSWORD=your-app-password

# Пароли для безопасности (в продакшене)
POSTGRES_PASSWORD=secure_password
REDIS_PASSWORD=redis_password
MINIO_ROOT_PASSWORD=minio_password

# Порты (если заняты стандартные)
USER_SERVICE_PORT=8082
API_SERVICE_PORT=8084
```

#### В локальных .env файлах сервисов:
Каждый сервис имеет свой .env файл для специфичных настроек. Эти файлы загружаются через `godotenv.Load()` в main.go каждого сервиса.

### 4. Запуск всей системы

С использованием Make:
```bash
make up
```

Или с помощью docker-compose:
```bash
docker-compose up -d
chmod +x scripts/init-minio.sh
./scripts/init-minio.sh
```

**📢 Важно**: При первом запуске миграции могут занять 1-2 минуты. Следите за логами сервисов.

### 5. Проверка запуска

```bash
make status
# или
make show-config  # показать текущие настройки
```

## Управление конфигурацией

### Команды для работы с файлами конфигурации

```bash
# Создать все файлы конфигурации
make setup-env

# Только Docker инфраструктура
make setup-compose

# Только .env файлы сервисов
make setup-services

# Использовать настройки для разработки
make env-dev

# Показать рекомендации для продакшена  
make env-prod

# Показать текущие настройки
make show-config
```

### Структура конфигурации

```
teamMessenger/
├── compose.env              # 🐳 Docker инфраструктура
├── compose.env.example      # 📋 Шаблон для compose.env
├── compose.development      # 🛠️ Настройки для разработки
├── userService/
│   ├── .env                 # ⚙️ Локальные настройки userService
│   └── env.example          # 📋 Шаблон для userService
├── apiService/
│   ├── .env                 # ⚙️ Локальные настройки apiService
│   └── env.example          # 📋 Шаблон для apiService
└── ... (остальные сервисы)
```

### Как работает система переменных

1. **В main.go сервиса**:
```go
// Загружаем локальный .env файл
if err := godotenv.Load(); err != nil {
    log.Printf("No .env file found: %v", err)
}
```

2. **Docker Compose переопределяет**:
```yaml
environment:
  - DB_HOST=postgres          # Переопределяет localhost
  - DB_USER=${POSTGRES_USER}  # Из compose.env
```

3. **Результат**: Сервис получает правильные настройки для Docker окружения

### Изменение настроек

#### Для Docker инфраструктуры (порты, пароли):
```bash
# Отредактируйте compose.env
POSTGRES_PORT=5433
API_SERVICE_PORT=8184

# Перезапустите
make restart
```

#### Для логики сервиса:
```bash
# Отредактируйте userService/.env
JWT_SECRET=new_secret_key

# Пересоберите сервис
docker-compose build user-service
docker-compose restart user-service
```

## Доступные эндпоинты

После запуска системы доступны следующие сервисы (порты читаются из `compose.env`):

| Сервис | URL | Описание |
|--------|-----|----------|
| API Gateway | http://localhost:[API_SERVICE_PORT] | Основная точка входа |
| User Service | http://localhost:[USER_SERVICE_PORT] | Swagger: /swagger/index.html |
| File Service | http://localhost:[FILE_SERVICE_PORT] | Swagger: /swagger/index.html |
| Task Service | http://localhost:[TASK_SERVICE_PORT] | Swagger: /swagger/index.html |
| Chat Service | http://localhost:[CHAT_SERVICE_PORT] | Swagger: /swagger/index.html |
| Notification Service | http://localhost:[NOTIFICATION_SERVICE_PORT] | Health: /health |
| MinIO Console | http://localhost:[MINIO_CONSOLE_PORT] | admin/[MINIO_ROOT_PASSWORD] |
| PostgreSQL | localhost:[POSTGRES_PORT] | [POSTGRES_USER]/[POSTGRES_PASSWORD] |
| Redis | localhost:[REDIS_PORT] | Пароль из REDIS_PASSWORD |
| Kafka | localhost:[KAFKA_PORT] | - |

Чтобы увидеть актуальные порты:
```bash
make show-config
```

## Доступ к файлам

Файлы доступны через Nginx прокси по адресу:
```
http://localhost:[NGINX_PORT]/teamfiles/[filename]
```

Например, если файл загружен как `avatar.jpg`:
```
http://localhost:8090/teamfiles/avatar.jpg
```

## Управление контейнерами

### С помощью Make

```bash
# Показать все команды
make help

# Собрать образы
make build

# Запустить систему
make up

# Остановить систему
make down

# Перезапустить
make restart

# Показать логи
make logs

# Показать логи в реальном времени
make logs-follow

# Показать статус
make status

# Показать текущие настройки
make show-config

# Пересобрать и перезапустить
make rebuild

# Полная очистка
make clean
```

### Управление конфигурацией

```bash
# Настройка всех файлов конфигурации
make setup-env

# Только Docker инфраструктура
make setup-compose

# Только локальные .env файлы сервисов
make setup-services

# Настройки для разработки
make env-dev

# Информация для продакшена
make env-prod
```

### Управление миграциями

```bash
# Проверить статус миграций
make check-migrations

# Принудительно запустить миграции для сервиса
make migrate-user
make migrate-file
make migrate-chat
make migrate-task
```

### Частичное управление

```bash
# Запустить только инфраструктуру
make up-infra

# Запустить только микросервисы
make up-services

# Остановить только микросервисы
make down-services
```

### Логи отдельных сервисов

```bash
make logs-api          # API Gateway
make logs-user         # User Service
make logs-file         # File Service
make logs-chat         # Chat Service
make logs-task         # Task Service
make logs-notification # Notification Service
```

## Конфигурация для продакшена

### Безопасность

Обязательно измените в `compose.env`:

```bash
# Сильные пароли
POSTGRES_PASSWORD=very_secure_db_password_123
REDIS_PASSWORD=secure_redis_password_456  
MINIO_ROOT_PASSWORD=secure_minio_password_789

# Реальные email настройки
SMTP_USERNAME=noreply@yourdomain.com
SMTP_PASSWORD=real_app_password
FROM_EMAIL=noreply@yourdomain.com
```

### Настройка локальных .env файлов

Проверьте и настройте локальные .env файлы в каждом сервисе:

```bash
# userService/.env
JWT_SECRET=production_jwt_secret
AUTH_TIMEOUT=30m

# apiService/.env
RATE_LIMIT=1000
CORS_ORIGINS=https://yourdomain.com

# fileService/.env
MAX_FILE_SIZE=50MB
ALLOWED_TYPES=jpg,png,pdf,doc
```

## Отладка

### Проверка здоровья сервисов

```bash
# Проверить все сервисы
curl http://localhost:$(grep '^API_SERVICE_PORT=' compose.env | cut -d'=' -f2)/health

# Проверить конкретный сервис
curl http://localhost:$(grep '^USER_SERVICE_PORT=' compose.env | cut -d'=' -f2)/swagger/index.html
```

### Проверка конфигурации

```bash
# Показать все настройки
make show-config

# Проверить compose.env файл
cat compose.env | grep -v '^#' | grep -v '^$'

# Проверить локальные .env файлы
ls -la */.*env

# Проверить переменные в контейнере
docker-compose exec user-service env | grep -E "DB_|APP_|KAFKA_"
```

### Проверка загрузки .env файлов

```bash
# Посмотреть логи загрузки .env
docker-compose logs user-service | grep -i "env\|load"

# Проверить переменные внутри контейнера
docker-compose exec user-service sh -c 'echo "DB_HOST=$DB_HOST, APP_PORT=$APP_PORT"'
```

### Проверка базы данных

```bash
# Подключиться к PostgreSQL
docker exec -it team-messenger-postgres psql -U $(grep '^POSTGRES_USER=' compose.env | cut -d'=' -f2) -d $(grep '^POSTGRES_DB=' compose.env | cut -d'=' -f2)

# Проверить таблицы
\dt

# Посмотреть схемы
\dn

# Проверить статус миграций
SELECT * FROM schema_migrations ORDER BY service, version;
```

## Решение проблем

### Проблемы с конфигурацией

```bash
# Проверить все файлы конфигурации
make show-config

# Пересоздать все файлы конфигурации
rm compose.env */.*env
make setup-env

# Применить настройки разработки
make env-dev

# Проверить переменные в контейнере
docker-compose exec api-service env | grep -E "_PORT|_PASSWORD"
```

### Проблемы с .env файлами сервисов

```bash
# Проверить наличие .env файлов
make show-config

# Создать недостающие .env файлы
make setup-services

# Посмотреть логи загрузки .env в сервисе
docker-compose logs user-service | grep -i "env file"

# Проверить содержимое .env в контейнере
docker-compose exec user-service cat .env
```

### Конфликты переменных

Если есть конфликты между переменными:

1. **Проверьте приоритет**:
   - Docker environment (высший)
   - Локальный .env файл
   - Системные переменные

2. **Отладьте переменные**:
```bash
# В контейнере
docker-compose exec user-service env | sort

# Сравните с локальным .env
cat userService/.env
```

3. **Измените docker-compose.yml** если нужно:
```yaml
environment:
  - DB_HOST=postgres  # Это переопределит .env
  # Уберите строку если хотите использовать из .env
```

### Проблемы с миграциями

```bash
# Посмотреть логи миграций
docker-compose logs [service-name] | grep -i migration

# Проверить подключение к БД
docker exec -it team-messenger-postgres pg_isready -U postgres

# Принудительно запустить миграции
make migrate-user  # или другой сервис

# Проверить схемы в БД
docker exec -it team-messenger-postgres psql -U postgres -d team_messenger -c "\dn"
```

### Проблемы с портами

Если порты заняты, измените их в `compose.env`:

```bash
# Измените нужные порты
USER_SERVICE_PORT=8182
API_SERVICE_PORT=8184
POSTGRES_PORT=5433

# Перезапустите
make restart
```

## Остановка и очистка

```bash
# Остановить контейнеры
make down

# Полная очистка (удаляет данные!)
make clean

# Или вручную
docker-compose down -v --rmi all --remove-orphans
docker system prune -f
```

**⚠️ Внимание**: команда `make clean` удаляет все данные включая базу данных и загруженные файлы!

## Заключение

Новая двухуровневая архитектура конфигурации обеспечивает:

- ✅ **Совместимость** с вашими существующими .env файлами
- ✅ **Гибкость** настройки инфраструктуры отдельно от логики
- ✅ **Безопасность** изоляции секретов инфраструктуры
- ✅ **Простоту** развертывания в разных окружениях

Ваши сервисы продолжают использовать `godotenv.Load()` и локальные .env файлы, но при работе в Docker получают правильные настройки инфраструктуры автоматически. 