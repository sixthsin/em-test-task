# Enrichment API

## Возможности

- RESTful API, построенный с использованием Gin
- База данных PostgreSQL для хранения данных с использованием gorm
- Логирование с использованием Logrus
- Конфигурация окружения с помощью Godotenv
- Документация Swagger
- Автоматизированный запуск c помощью docker-compose

## Предварительные требования

- Docker
- Docker Compose

## Начало работы

### Переменные окружения

Создайте файл `.env` в корневом каталоге и настройте переменные окружения
Например:
```env
PORT=":8080"
DSN="host=localhost port=5432 dbname=mydatabase user=myuser password=mypassword"
```

### Сборка и запуск с Docker Compose

Для сборки и запуска приложения с использованием Docker Compose выполните следующую команду:

```bash
docker-compose up --build
```

API будет доступен по адресу `http://localhost:8080/api/v1/`.

### Документация API

Документация Swagger доступна по адресу `http://localhost:8080/api/v1/swagger/index.html`.

## Структура проекта

- `cmd/`: Содержит точку входа в основное приложение.
- `pkg/`: Содержит основную логику приложения.
- `internal/`: Содержит внутренние пакеты, не предназначенные для публичного использования.
- `cfg/`: Содержит файлы конфигурации.
- `migrations/`: Содержит миграции базы данных.
- `docs/`: Содержит файлы документации.

## Зависимости

Проект использует следующие модули Go:

- `github.com/gin-gonic/gin`
- `github.com/go-playground/validator/v10`
- `github.com/joho/godotenv`
- `github.com/sirupsen/logrus`
- `gorm.io/driver/postgres`
- `gorm.io/gorm`
