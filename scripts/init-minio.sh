#!/bin/bash

# Скрипт для инициализации MinIO bucket
echo "Waiting for MinIO to be ready..."

# Ждем пока MinIO запустится
until curl -f http://localhost:9000/minio/health/live; do
  echo "Waiting for MinIO..."
  sleep 5
done

echo "MinIO is ready. Creating bucket..."

# Устанавливаем MinIO client если его нет
if ! command -v mc &> /dev/null; then
    echo "Installing MinIO client..."
    wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /tmp/mc
    chmod +x /tmp/mc
    MC_CMD="/tmp/mc"
else
    MC_CMD="mc"
fi

# Настраиваем алиас для локального MinIO
$MC_CMD alias set myminio http://localhost:9000 minioadmin minioadmin123

# Создаем bucket если его нет
$MC_CMD mb myminio/teamfiles --ignore-existing

# Устанавливаем публичную политику для чтения файлов
$MC_CMD anonymous set public myminio/teamfiles

echo "MinIO bucket 'teamfiles' created and configured successfully!"
echo "Files will be accessible via: http://localhost:8090/teamfiles/[filename]" 