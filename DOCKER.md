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

### Frontend
- **Frontend** - React SPA приложение (статичные файлы)
- **Frontend Nginx** (8091) - Прокси для раздачи frontend приложения

### Инфраструктура
- **PostgreSQL** (5432) - Основная база данных
- **Redis** (6379) - Кеш и сессии
- **Kafka** (9092) - Очереди сообщений
- **MinIO** (9000/9001) - S3-совместимое хранилище файлов
- **Backend Nginx** (8090) - Прокси для API и MinIO (с rate limiting и WAF)
- **Zookeeper** (2181) - Координация Kafka

### Сети Docker

Система использует две изолированные сети:
- **frontend-network** - для frontend и frontend-nginx (публичный доступ)
- **backend-network** - для всех микросервисов и инфраструктуры (приватная)

## 🏗️ Архитектура конфигурации

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

#### Создание compose.env:

```bash
# Скопируйте шаблон
cp compose.env.example compose.env

# Отредактируйте compose.env под ваши нужды
nano compose.env
```

#### Основные настройки в compose.env:

```bash
# Email для уведомлений
SMTP_USERNAME=your-email@yandex.ru
SMTP_PASSWORD=your-app-password

# Пароли для безопасности (в продакшене обязательно измените!)
POSTGRES_PASSWORD=secure_password
REDIS_PASSWORD=redis_password
MINIO_ROOT_PASSWORD=minio_password

# Порты (если заняты стандартные)
USER_SERVICE_PORT=8082
API_SERVICE_PORT=8084
FRONTEND_NGINX_PORT=8091
NGINX_PORT=8090
```

#### В локальных .env файлах сервисов:

Каждый сервис имеет свой `.env` файл для специфичных настроек. Эти файлы загружаются через `godotenv.Load()` в main.go каждого сервиса.

Создайте `.env` файлы на основе `env.example` в каждом сервисе:
```bash
# Для каждого сервиса
cp userService/env.example userService/.env
cp apiService/env.example apiService/.env
# и т.д.
```

### 3. Запуск всей системы

С использованием Make:
```bash
make up
```

Или с помощью docker-compose напрямую:
```bash
docker compose -f docker-compose.yml --env-file compose.env up -d
```

**📢 Важно**: 
- При первом запуске миграции могут занять 1-2 минуты. Следите за логами сервисов.
- MinIO bucket создаётся автоматически через контейнер `minio-init`
- Frontend собирается автоматически при первом запуске

### 4. Проверка запуска

```bash
# Проверить статус контейнеров
docker compose ps

# Проверить логи
docker compose logs -f

# Проверить здоровье сервисов
curl http://localhost:8090/api/v1/health
curl http://localhost:8091/health
```

## Доступные эндпоинты

После запуска системы доступны следующие сервисы (порты настраиваются в `compose.env`):

| Сервис | URL | Описание |
|--------|-----|----------|
| **Frontend** | http://localhost:[FRONTEND_NGINX_PORT] | React SPA приложение (по умолчанию 8091) |
| **API Gateway** | http://localhost:[NGINX_PORT]/api/v1 | Основная точка входа для API (по умолчанию 8090) |
| **User Service** | http://localhost:[USER_SERVICE_PORT] | Swagger: /swagger/index.html (по умолчанию 8082) |
| **File Service** | http://localhost:[FILE_SERVICE_PORT] | Swagger: /swagger/index.html (по умолчанию 8080) |
| **Task Service** | http://localhost:[TASK_SERVICE_PORT] | Swagger: /swagger/index.html (по умолчанию 8081) |
| **Chat Service** | http://localhost:[CHAT_SERVICE_PORT] | Swagger: /swagger/index.html (по умолчанию 8083) |
| **Notification Service** | http://localhost:[NOTIFICATION_SERVICE_PORT] | Health: /health (по умолчанию 8085) |
| **MinIO Console** | http://localhost:[MINIO_CONSOLE_PORT] | admin/[MINIO_ROOT_PASSWORD] (по умолчанию 9001) |
| **PostgreSQL** | localhost:[POSTGRES_PORT] | [POSTGRES_USER]/[POSTGRES_PASSWORD] (по умолчанию 5432) |
| **Redis** | localhost:[REDIS_PORT] | Пароль из REDIS_PASSWORD (по умолчанию 6379) |
| **Kafka** | localhost:[KAFKA_PORT] | - (по умолчанию 9092) |

**Важно**: 
- API Service не экспонирует порт наружу - доступ только через Nginx на порту 8090
- Frontend доступен через отдельный Nginx на порту 8091
- Все API запросы должны идти через `/api/v1` префикс

## Доступ к файлам

Файлы доступны через Backend Nginx прокси по адресу:
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

# Запустить систему
make up

# Пересобрать и запустить
make build-up

# Остановить систему
make down

# Запуск в debug режиме
make debug

# Запуск с тестами
make test          # Unit тесты + запуск
make test-full     # Unit + Integration тесты + запуск
```

### С помощью Docker Compose напрямую

```bash
# Запустить систему
docker compose -f docker-compose.yml --env-file compose.env up -d

# Остановить систему
docker compose -f docker-compose.yml --env-file compose.env down

# Пересобрать и запустить
docker compose -f docker-compose.yml --env-file compose.env build
docker compose -f docker-compose.yml --env-file compose.env up -d

# Показать логи
docker compose logs -f

# Показать логи конкретного сервиса
docker compose logs -f api-service
docker compose logs -f frontend
docker compose logs -f nginx

# Перезапустить сервис
docker compose restart api-service

# Пересобрать конкретный сервис
docker compose build api-service
docker compose up -d api-service
```

### Проверка статуса

```bash
# Статус всех контейнеров
docker compose ps

# Проверка здоровья
docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"

# Проверка логов конкретного сервиса
docker compose logs api-service | tail -50
docker compose logs frontend | tail -50
```

## Структура конфигурации

```
teamMessenger/
├── compose.env              # 🐳 Docker инфраструктура
├── compose.env.example      # 📋 Шаблон для compose.env
├── docker-compose.yml       # 🐳 Основной compose файл
├── docker-compose.debug.yml  # 🐛 Debug режим
├── frontend/
│   ├── Dockerfile           # Frontend сборка
│   ├── nginx.conf           # Nginx для frontend контейнера
│   └── nginx-gateway.conf   # Nginx для frontend-nginx прокси
├── nginx/
│   ├── nginx.conf           # Backend Nginx конфигурация
│   └── conf.d/              # WAF правила, blacklist, whitelist
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
FRONTEND_NGINX_PORT=8092

# Перезапустите
docker compose -f docker-compose.yml --env-file compose.env restart
```

#### Для логики сервиса:
```bash
# Отредактируйте userService/.env
JWT_SECRET=new_secret_key

# Пересоберите сервис
docker compose build user-service
docker compose up -d user-service
```

#### Для Frontend:
```bash
# Отредактируйте frontend/.env.local (если есть)
VITE_API_URL=http://localhost:8090

# Пересоберите frontend
docker compose build frontend
docker compose up -d frontend
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

### Настройка Nginx для продакшена

Backend Nginx уже настроен с:
- Rate limiting для защиты от DDoS
- WAF правилами (в `nginx/conf.d/waf_rules.conf`)
- Blacklist/Whitelist IP (в `nginx/conf.d/blacklist.conf` и `whitelist.conf`)

Настройте правила под ваши нужды в `nginx/conf.d/`.

## Отладка

### Проверка здоровья сервисов

```bash
# Проверить API Gateway через Nginx
curl http://localhost:8090/api/v1/health

# Проверить Frontend
curl http://localhost:8091/health

# Проверить конкретный сервис напрямую
curl http://localhost:8082/health  # User Service
curl http://localhost:8080/health   # File Service
```

### Проверка конфигурации

```bash
# Проверить compose.env файл
cat compose.env | grep -v '^#' | grep -v '^$'

# Проверить локальные .env файлы
ls -la */.*env

# Проверить переменные в контейнере
docker compose exec user-service env | grep -E "DB_|APP_|KAFKA_"
```

### Проверка загрузки .env файлов

```bash
# Посмотреть логи загрузки .env
docker compose logs user-service | grep -i "env\|load"

# Проверить переменные внутри контейнера
docker compose exec user-service sh -c 'echo "DB_HOST=$DB_HOST, APP_PORT=$APP_PORT"'
```

### Проверка базы данных

```bash
# Подключиться к PostgreSQL
docker exec -it team-messenger-postgres psql -U postgres -d team_messenger

# Проверить таблицы
\dt

# Посмотреть схемы
\dn

# Проверить статус миграций
SELECT * FROM schema_migrations ORDER BY service, version;
```

### Проверка Frontend

```bash
# Проверить логи frontend
docker compose logs frontend

# Проверить логи frontend-nginx
docker compose logs frontend-nginx

# Проверить доступность
curl http://localhost:8091/

# Проверить переменные окружения frontend
docker compose exec frontend env | grep VITE
```

## Решение проблем

### Проблемы с конфигурацией

```bash
# Проверить compose.env файл
cat compose.env | grep -v '^#' | grep -v '^$'

# Проверить переменные в контейнере
docker compose exec api-service env | grep -E "_PORT|_PASSWORD"

# Проверить что compose.env загружается
docker compose config | grep -A 5 "environment:"
```

### Проблемы с .env файлами сервисов

```bash
# Проверить наличие .env файлов
ls -la */.*env

# Посмотреть логи загрузки .env в сервисе
docker compose logs user-service | grep -i "env file"

# Проверить содержимое .env в контейнере
docker compose exec user-service cat .env
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
docker compose exec user-service env | sort

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
docker compose logs user-service | grep -i migration

# Проверить подключение к БД
docker exec -it team-messenger-postgres pg_isready -U postgres

# Проверить схемы в БД
docker exec -it team-messenger-postgres psql -U postgres -d team_messenger -c "\dn"
```

### Проблемы с портами

Если порты заняты, измените их в `compose.env`:

```bash
# Измените нужные порты
USER_SERVICE_PORT=8182
API_SERVICE_PORT=8184
FRONTEND_NGINX_PORT=8092
NGINX_PORT=8091
POSTGRES_PORT=5433

# Перезапустите
docker compose -f docker-compose.yml --env-file compose.env down
docker compose -f docker-compose.yml --env-file compose.env up -d
```

### Проблемы с Frontend

```bash
# Проверить сборку frontend
docker compose logs frontend | grep -i "build\|error"

# Пересобрать frontend
docker compose build frontend
docker compose up -d frontend

# Проверить nginx конфигурацию frontend
docker compose exec frontend-nginx nginx -t

# Проверить доступность frontend контейнера
docker compose exec frontend-nginx wget -O- http://frontend:80/health
```

### Проблемы с Nginx

```bash
# Проверить конфигурацию backend nginx
docker compose exec nginx nginx -t

# Проверить логи nginx
docker compose logs nginx | tail -100

# Перезапустить nginx
docker compose restart nginx
```

## Остановка и очистка

```bash
# Остановить контейнеры
make down
# или
docker compose -f docker-compose.yml --env-file compose.env down

# Остановить и удалить volumes (удаляет данные!)
docker compose -f docker-compose.yml --env-file compose.env down -v

# Полная очистка (удаляет данные и образы!)
docker compose -f docker-compose.yml --env-file compose.env down -v --rmi all --remove-orphans
docker system prune -f
```

**⚠️ Внимание**: команда `down -v` удаляет все данные включая базу данных и загруженные файлы!

## Заключение

Двухуровневая архитектура конфигурации обеспечивает:

- ✅ **Совместимость** с вашими существующими .env файлами
- ✅ **Гибкость** настройки инфраструктуры отдельно от логики
- ✅ **Безопасность** изоляции секретов инфраструктуры
- ✅ **Простоту** развертывания в разных окружениях
- ✅ **Масштабируемость** через изолированные Docker сети

Ваши сервисы продолжают использовать `godotenv.Load()` и локальные .env файлы, но при работе в Docker получают правильные настройки инфраструктуры автоматически.
