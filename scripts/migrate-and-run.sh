#!/bin/bash

set -e

echo "Starting migration and service script..."

# Параметры подключения к БД
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-team_messenger}

# Путь к миграциям (передается как аргумент)
MIGRATIONS_PATH=${1:-"./migrations"}
SERVICE_NAME=${2:-"service"}

echo "Service: $SERVICE_NAME"
echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "Migrations path: $MIGRATIONS_PATH"

# Функция ожидания готовности PostgreSQL
wait_for_postgres() {
    echo "Waiting for PostgreSQL to be ready..."
    until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do
        echo "PostgreSQL is not ready yet, waiting..."
        sleep 2
    done
    echo "PostgreSQL is ready!"
}

# Функция создания схемы для сервиса
create_schema() {
    local schema_name=$1
    echo "Creating schema '$schema_name' if not exists..."
    
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        CREATE SCHEMA IF NOT EXISTS $schema_name;
    " || echo "Schema creation failed or already exists"
}

# Функция выполнения миграций
run_migrations() {
    echo "Running migrations from $MIGRATIONS_PATH..."
    
    if [ -d "$MIGRATIONS_PATH" ]; then
        # Проверяем есть ли файлы миграций
        if ls $MIGRATIONS_PATH/*.up.sql 1> /dev/null 2>&1; then
            echo "Found migration files:"
            ls $MIGRATIONS_PATH/*.up.sql
            
            # Выполняем миграции в порядке версий
            for migration_file in $(ls $MIGRATIONS_PATH/*.up.sql | sort); do
                echo "Applying migration: $migration_file"
                
                # Извлекаем номер версии из имени файла
                version=$(basename "$migration_file" | cut -d'_' -f1)
                
                # Проверяем, не была ли уже применена эта миграция
                migration_applied=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
                    SELECT EXISTS (
                        SELECT 1 FROM information_schema.tables 
                        WHERE table_name = 'schema_migrations'
                    );
                " 2>/dev/null | tr -d ' ')

                if [ "$migration_applied" = "t" ]; then
                    # Таблица миграций существует, проверяем конкретную миграцию
                    already_applied=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
                        SELECT EXISTS (
                            SELECT 1 FROM schema_migrations 
                            WHERE version = '$version' AND service = '$SERVICE_NAME'
                        );
                    " 2>/dev/null | tr -d ' ')
                else
                    # Создаем таблицу для отслеживания миграций
                    echo "Creating schema_migrations table..."
                    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
                        CREATE TABLE IF NOT EXISTS schema_migrations (
                            version VARCHAR(50) NOT NULL,
                            service VARCHAR(100) NOT NULL,
                            applied_at TIMESTAMP DEFAULT NOW(),
                            PRIMARY KEY (version, service)
                        );
                    "
                    already_applied="f"
                fi

                if [ "$already_applied" = "f" ]; then
                    echo "Applying new migration $version for service $SERVICE_NAME..."
                    
                    # Применяем миграцию
                    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration_file"; then
                        # Записываем информацию о примененной миграции
                        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
                            INSERT INTO schema_migrations (version, service) VALUES ('$version', '$SERVICE_NAME');
                        "
                        echo "Migration $version applied successfully!"
                    else
                        echo "ERROR: Failed to apply migration $version"
                        exit 1
                    fi
                else
                    echo "Migration $version already applied for service $SERVICE_NAME, skipping..."
                fi
            done
            
            echo "All migrations completed successfully!"
        else
            echo "No migration files found in $MIGRATIONS_PATH"
        fi
    else
        echo "Migrations directory $MIGRATIONS_PATH not found, skipping migrations..."
    fi
}

# Основная логика
echo "=== Starting Migration Process ==="

# Ждем готовности PostgreSQL
wait_for_postgres

# Определяем имя схемы на основе имени сервиса
case $SERVICE_NAME in
    "user-service")
        SCHEMA_NAME="user_service"
        ;;
    "file-service")
        SCHEMA_NAME="file_service"
        ;;
    "task-service")
        SCHEMA_NAME="task_service"
        ;;
    "chat-service")
        SCHEMA_NAME="chat_service"
        ;;
    "notification-service")
        SCHEMA_NAME="notification_service"
        ;;
    *)
        SCHEMA_NAME="public"
        ;;
esac

# Создаем схему для сервиса
if [ "$SCHEMA_NAME" != "public" ]; then
    create_schema $SCHEMA_NAME
fi

# Выполняем миграции
run_migrations

echo "=== Migration Process Completed ==="

# Запускаем основное приложение
echo "Starting main application..."
exec ./main 