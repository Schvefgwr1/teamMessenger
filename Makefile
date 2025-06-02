.PHONY: help build up down restart logs clean init-minio

# Показать справку
help:
	@echo "Team Messenger Docker Commands:"
	@echo "  make build      - Собрать все образы"
	@echo "  make up         - Запустить все контейнеры"
	@echo "  make down       - Остановить все контейнеры"
	@echo "  make restart    - Перезапустить все контейнеры"
	@echo "  make logs       - Показать логи всех контейнеров"
	@echo "  make logs-follow - Показать логи в реальном времени"
	@echo "  make clean      - Удалить все контейнеры и образы"
	@echo "  make init-minio - Инициализировать MinIO bucket"
	@echo "  make status     - Показать статус контейнеров"
	@echo "  make rebuild    - Пересобрать и перезапустить все"
	@echo ""
	@echo "📁 Environment files:"
	@echo "  make setup-compose  - Настроить compose.env для Docker инфраструктуры"
	@echo "  make env-dev        - Использовать настройки для разработки"
	@echo "  make env-prod       - Показать пример для продакшена"
	@echo "  make setup-services - Создать .env файлы для всех сервисов"
	@echo ""
	@echo "ℹ️  Новая архитектура:"
	@echo "  📄 compose.env      - Настройки Docker инфраструктуры"
	@echo "  📄 сервисы/.env     - Локальные настройки каждого сервиса"
	@echo "  🔄 Миграции БД выполняются автоматически при запуске"

# Настройка compose.env файла для Docker инфраструктуры
setup-compose:
	@if [ ! -f compose.env ]; then \
		echo "Создаем compose.env файл из шаблона..."; \
		cp compose.env.example compose.env; \
		echo "✅ Файл compose.env создан! Отредактируйте его под ваши нужды."; \
		echo "📧 Особенно важно настроить SMTP для email уведомлений."; \
	else \
		echo "⚠️  Файл compose.env уже существует. Если хотите пересоздать его:"; \
		echo "   rm compose.env && make setup-compose"; \
	fi

# Создать .env файлы для всех сервисов из шаблонов
setup-services:
	@echo "🔧 Создаем .env файлы для каждого сервиса..."
	@for service in userService apiService fileService chatService taskService notificationService; do \
		if [ -f $$service/env.example ]; then \
			if [ ! -f $$service/.env ]; then \
				cp $$service/env.example $$service/.env; \
				echo "✅ Создан $$service/.env"; \
			else \
				echo "⚠️  $$service/.env уже существует"; \
			fi; \
		else \
			echo "❌ Не найден $$service/env.example"; \
		fi; \
	done
	@echo "📋 Все .env файлы сервисов готовы к настройке!"

# Полная настройка окружения (compose + сервисы)
setup-env: setup-compose setup-services
	@echo ""
	@echo "🎉 Настройка окружения завершена!"
	@echo ""
	@echo "📝 Теперь отредактируйте файлы:"
	@echo "  📄 compose.env - порты, пароли инфраструктуры"
	@echo "  📄 каждый_сервис/.env - локальные настройки сервисов"

# Использовать настройки для разработки
env-dev:
	@echo "Копируем настройки для разработки..."
	@cp compose.development compose.env
	@echo "✅ Настройки Docker инфраструктуры для разработки применены!"
	@echo "📝 Локальные .env файлы сервисов остаются без изменений"

# Показать пример для продакшена
env-prod:
	@echo "📋 Для продакшена рекомендуется изменить следующие переменные в compose.env:"
	@echo ""
	@echo "# Безопасные пароли:"
	@echo "POSTGRES_PASSWORD=secure_db_password_123"
	@echo "REDIS_PASSWORD=secure_redis_password"
	@echo "MINIO_ROOT_PASSWORD=secure_minio_password_123"
	@echo ""
	@echo "# Email настройки:"
	@echo "SMTP_USERNAME=your-real-email@domain.com"
	@echo "SMTP_PASSWORD=your-real-app-password"
	@echo ""
	@echo "# При необходимости, другие порты:"
	@echo "USER_SERVICE_PORT=8082"
	@echo "API_SERVICE_PORT=8084"
	@echo "# и т.д."
	@echo ""
	@echo "📝 Также проверьте локальные .env файлы в каждом сервисе!"

# Собрать все образы
build:
	@make setup-env
	docker-compose --env-file compose.env build

# Запустить все контейнеры
up:
	@make setup-env
	docker-compose --env-file compose.env up -d
	@echo "Ожидание готовности сервисов..."
	@echo "🔄 Миграции БД выполняются автоматически..."
	@sleep 45
	@echo "Инициализация MinIO bucket..."
	@make init-minio
	@echo ""
	@echo "🚀 Team Messenger запущен!"
	@echo ""
	@echo "📋 Доступные сервисы:"
	@echo "  API Gateway:      http://localhost:$$(grep '^API_SERVICE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8084')"
	@echo "  User Service:     http://localhost:$$(grep '^USER_SERVICE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8082')"
	@echo "  File Service:     http://localhost:$$(grep '^FILE_SERVICE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8080')" 
	@echo "  Task Service:     http://localhost:$$(grep '^TASK_SERVICE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8081')"
	@echo "  Chat Service:     http://localhost:$$(grep '^CHAT_SERVICE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8083')"
	@echo "  Notification:     http://localhost:$$(grep '^NOTIFICATION_SERVICE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8085')"
	@echo ""
	@echo "🗄️  Инфраструктура:"
	@echo "  PostgreSQL:       localhost:$$(grep '^POSTGRES_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '5432')"
	@echo "  Redis:            localhost:$$(grep '^REDIS_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '6379')"
	@echo "  Kafka:            localhost:$$(grep '^KAFKA_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '9092')"
	@echo "  MinIO Console:    http://localhost:$$(grep '^MINIO_CONSOLE_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '9001')"
	@echo "  MinIO API:        http://localhost:$$(grep '^MINIO_API_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '9000')"
	@echo "  MinIO Proxy:      http://localhost:$$(grep '^NGINX_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8090')"
	@echo ""
	@echo "📁 Файлы доступны по: http://localhost:$$(grep '^NGINX_PORT=' compose.env 2>/dev/null | cut -d'=' -f2 || echo '8090')/teamfiles/[filename]"
	@echo "💾 Миграции автоматически применены для всех сервисов"

# Остановить все контейнеры
down:
	docker-compose --env-file compose.env down

# Перезапустить все контейнеры
restart:
	docker-compose --env-file compose.env restart

# Показать логи
logs:
	docker-compose --env-file compose.env logs

# Показать логи в реальном времени
logs-follow:
	docker-compose --env-file compose.env logs -f

# Полная очистка
clean:
	docker-compose --env-file compose.env down -v --rmi all --remove-orphans
	docker system prune -f

# Инициализировать MinIO bucket
init-minio:
	@echo "Настройка MinIO bucket..."
	@chmod +x scripts/init-minio.sh
	@./scripts/init-minio.sh

# Показать статус контейнеров
status:
	docker-compose --env-file compose.env ps

# Показать текущие настройки из compose.env
show-config:
	@echo "📋 Текущие настройки из compose.env файла:"
	@echo ""
	@if [ -f compose.env ]; then \
		echo "Порты сервисов:"; \
		grep -E "SERVICE_PORT|POSTGRES_PORT|REDIS_PORT|KAFKA_PORT|NGINX_PORT|MINIO.*PORT" compose.env | sed 's/^/  /' || echo "  Не найдено"; \
		echo ""; \
		echo "Настройки БД:"; \
		grep -E "POSTGRES_" compose.env | sed 's/^/  /' || echo "  Не найдено"; \
		echo ""; \
		echo "Настройки MinIO:"; \
		grep -E "MINIO_" compose.env | sed 's/^/  /' || echo "  Не найдено"; \
	else \
		echo "❌ Файл compose.env не найден. Выполните: make setup-compose"; \
	fi
	@echo ""
	@echo "📂 Локальные .env файлы сервисов:"
	@for service in userService apiService fileService chatService taskService notificationService; do \
		if [ -f $$service/.env ]; then \
			echo "  ✅ $$service/.env"; \
		else \
			echo "  ❌ $$service/.env"; \
		fi; \
	done

# Пересобрать и перезапустить
rebuild:
	docker-compose --env-file compose.env down
	docker-compose --env-file compose.env build --no-cache
	docker-compose --env-file compose.env up -d

# Запуск только инфраструктуры (БД, Redis, Kafka, MinIO)
up-infra:
	@make setup-env
	docker-compose --env-file compose.env up -d postgres redis zookeeper kafka minio nginx
	@echo "Инфраструктура запущена"

# Запуск только микросервисов
up-services:
	docker-compose --env-file compose.env up -d user-service file-service task-service chat-service notification-service api-service
	@echo "Микросервисы запущены (миграции выполняются автоматически)"

# Остановить только микросервисы  
down-services:
	docker-compose --env-file compose.env stop user-service file-service task-service chat-service notification-service api-service

# Просмотр логов конкретного сервиса
logs-user:
	docker-compose --env-file compose.env logs -f user-service

logs-api:
	docker-compose --env-file compose.env logs -f api-service

logs-file:
	docker-compose --env-file compose.env logs -f file-service

logs-chat:
	docker-compose --env-file compose.env logs -f chat-service

logs-task:
	docker-compose --env-file compose.env logs -f task-service

logs-notification:
	docker-compose --env-file compose.env logs -f notification-service

# Проверка статуса миграций
check-migrations:
	@echo "Проверка статуса миграций в БД..."
	@docker exec -it team-messenger-postgres psql -U postgres -d team_messenger -c "SELECT * FROM schema_migrations ORDER BY service, version;"

# Принудительный запуск миграций для конкретного сервиса
migrate-user:
	docker-compose --env-file compose.env exec user-service ./migrate-and-run.sh ./migrations user-service

migrate-file:
	docker-compose --env-file compose.env exec file-service ./migrate-and-run.sh ./migrations file-service

migrate-chat:
	docker-compose --env-file compose.env exec chat-service ./migrate-and-run.sh ./migrations chat-service

migrate-task:
	docker-compose --env-file compose.env exec task-service ./migrate-and-run.sh ./migrations task-service 