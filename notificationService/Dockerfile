# Dockerfile для Notification Service
FROM golang:1.23.1-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Создаем рабочую директорию
WORKDIR /app

# Отключаем go.work для изолированной сборки
ENV GOWORK=off

# Копируем общий модуль
COPY common/ ./common/

# Копируем модуль notificationService
COPY notificationService/ ./notificationService/

# Переходим в директорию сервиса и собираем приложение
WORKDIR /app/notificationService

# Сначала очищаем и обновляем go.mod, затем загружаем зависимости
RUN go mod tidy && go mod download

# Собираем статически связанный бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

# Финальный этап - минимальный образ
FROM alpine:latest

# Устанавливаем сертификаты и timezone
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/notificationService/main .
COPY --from=builder /app/notificationService/config ./config/

# Копируем .env файл для локальных настроек (загружается через godotenv.Load())
# Создаем пустой .env если файл не существует, чтобы избежать ошибок
COPY --from=builder /app/notificationService/.env* ./
RUN touch .env

# Копируем шаблоны email уведомлений
COPY --from=builder /app/notificationService/templates ./templates/

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser
RUN chown -R appuser:appuser /root
USER appuser

# Открываем порт
EXPOSE 8085

# Запускаем приложение
CMD ["./main"] 