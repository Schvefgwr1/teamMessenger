# File Service

Сервис для управления файлами, построенный на Go с использованием Gin, MinIO и PostgreSQL.

## Описание

File Service предоставляет API для:
- Загрузки файлов
- Получения информации о файлах
- Переименования файлов
- Управления типами файлов
- Получения списка файлов с пагинацией

## Технологии

- Go 1.22
- Gin (веб-фреймворк)
- MinIO (объектное хранилище)
- PostgreSQL (база данных)
- Swagger/OpenAPI (документация API)

## Требования

- Go 1.22 или выше
- PostgreSQL
- MinIO сервер

## Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd fileService
```

2. Установите зависимости:
```bash
go mod download
```

3. Установите swag для генерации документации:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. Сгенерируйте документацию:
```bash
swag init -g cmd/server/main.go
```

## Конфигурация

Сервис использует конфигурационный файл `config/config.yaml`. Пример конфигурации:

```yaml
app:
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: file_service
  sslmode: disable

minio:
  host: localhost
  port: 9000
  access_key: your_access_key
  secret_key: your_secret_key
  bucket: files
  use_ssl: false
```

## Запуск

```bash
go run cmd/server/main.go
```

## API Endpoints

### Файлы

- `POST /files` - Загрузка файла
- `GET /files/{file_id}` - Получение информации о файле
- `PUT /files/{file_id}/rename` - Переименование файла
- `GET /files` - Получение списка файлов с пагинацией

### Типы файлов

- `POST /file-types` - Создание типа файла
- `GET /file-types/{id}` - Получение типа файла по ID
- `GET /file-types/name/{name}` - Получение типа файла по названию
- `DELETE /file-types/{id}` - Удаление типа файла

## Документация API

Swagger UI доступен по адресу:
```
http://localhost:8080/swagger/index.html
```

## Структура проекта

```
fileService/
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── docs/
│   └── swagger.json
├── internal/
│   ├── controllers/
│   │   ├── file_controller.go
│   │   └── file_type_controller.go
│   ├── handlers/
│   │   ├── file_handler.go
│   │   └── file_type_handler.go
│   ├── models/
│   │   ├── file.go
│   │   └── file_type.go
│   ├── repositories/
│   │   ├── file_repository.go
│   │   └── file_type_repository.go
│   └── routes/
│       └── routes.go
├── migrations/
├── go.mod
└── README.md
```

## Разработка

### Добавление новых эндпоинтов

1. Создайте новый метод в контроллере
2. Создайте обработчик в handlers
3. Добавьте роут в routes
4. Добавьте документацию Swagger
5. Сгенерируйте документацию: `swag init -g cmd/server/main.go`

### Миграции базы данных

Миграции находятся в директории `migrations/`. Для применения миграций используйте соответствующие инструменты PostgreSQL.

## Лицензия

MIT
