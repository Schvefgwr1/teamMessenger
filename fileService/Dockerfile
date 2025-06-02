# Dockerfile для File Service
FROM golang:1.23.1-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Создаем рабочую директорию
WORKDIR /app

# Отключаем go.work для изолированной сборки
ENV GOWORK=off

# Копируем общий модуль
COPY common/ ./common/

# Копируем модуль fileService
COPY fileService/ ./fileService/

# Копируем скрипт миграций
COPY scripts/ ./scripts/

# Переходим в директорию сервиса и собираем приложение
WORKDIR /app/fileService

# Сначала очищаем и обновляем go.mod, затем загружаем зависимости
RUN go mod tidy && go mod download

# Собираем статически связанный бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

# Финальный этап - минимальный образ
FROM alpine:latest

# Устанавливаем сертификаты, timezone и PostgreSQL client
RUN apk --no-cache add ca-certificates tzdata postgresql-client bash

WORKDIR /root/

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/fileService/main .
COPY --from=builder /app/fileService/config ./config/

# Копируем .env файл для локальных настроек (загружается через godotenv.Load())
# Создаем пустой .env если файл не существует, чтобы избежать ошибок
COPY --from=builder /app/fileService/.env* ./
RUN touch .env

# Копируем миграции
COPY --from=builder /app/fileService/migrations ./migrations/

# Копируем скрипт миграций из builder этапа
COPY --from=builder /app/scripts/migrate-and-run.sh ./migrate-and-run.sh
RUN chmod +x ./migrate-and-run.sh

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser
RUN chown -R appuser:appuser /root
USER appuser

# Открываем порт
EXPOSE 8080

# Используем скрипт миграций как entrypoint
CMD ["./migrate-and-run.sh", "./migrations", "file-service"] 