## Сервис проверки валидности ссылок

## Простой веб-сервис для проверки доступности сайтов. Умеет проверять ссылки и генерировать отчеты в PDF.

### Что делает:
1. Проверяет, доступны ли сайты по ссылкам

2. Сохраняет результаты проверок

3. Позволяет скачать отчет в PDF по ID проверок

### Как использовать:

1. Проверить ссылки
curl -X POST http://localhost:8080/api/check-links \
  -H "Content-Type: application/json" \
  -d '{"links": ["google.com", "example.com"]}'
Ответ:

{
  "links": {
    "google.com": "available",
    "example.com": "not available"
  },
  "links_num": 1
}

2. Получить отчет в PDF

POST http://localhost:8080/api/generate-report \
  -H "Content-Type: application/json" \
  -d '{"links_list": [1, 2]}' \
  -o report.pdf

## Запуск сервиса:

go mod tidy

## Запустите сервер:

go run cmd/server/main.go
Сервер будет доступен на http://localhost:8080

## Структура проекта
cmd/server/main.go - точка входа

internal/app/handlers/ - обработчики HTTP запросов

internal/domain/ - модели данных и бизнес-логика

internal/infra/config/ - настройки

## Как это работает:

При проверке ссылок создается HTTP клиент с таймаутом 10 секунд

Результаты сохраняются в памяти (вместе с ID проверки)

Для генерации PDF используется список ID проверок

При остановке сервера (Ctrl+C) он завершает текущие операции

## Ограничения
Данные хранятся только в памяти и пропадают после перезапуска

## Тестирование

### Unit тесты
```bash
# Запуск всех unit тестов
go test ./...

# Запуск с подробным выводом
go test -v ./...

# Запуск тестов конкретного пакета
go test ./internal/domain/links/service
go test ./internal/domain/links/repository
go test ./internal/app/handlers/check_links_handler
go test ./internal/app/handlers/generate_report_handler
go test ./internal/infra/config

# Запуск с покрытием кода
go test -cover ./...
```

### Интеграционные тесты
```bash
# Все тесты (unit + integration)
go test ./...

# Только интеграционные тесты
go test ./test/...
go test ./internal/app/...

# С подробным выводом
go test -v ./test/...
go test -v ./internal/app/...

# С таймаутом для медленных HTTP запросов
go test -timeout 30s ./test/...
```

### Проверка сборки
```bash
# Проверка компиляции всех пакетов
go build ./...

# Сборка исполняемого файла
go build ./cmd/server

# Проверка зависимостей
go mod tidy
go mod verify
```

### Покрытие тестами
```bash
# Генерация отчета о покрытии
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Просмотр покрытия в терминале
go test -cover ./...
```

## Что покрывают тесты

### Unit тесты:
- Repository: сохранение/получение батчей, генерация ID
- Service: проверка ссылок, генерация PDF отчетов
- Handlers: HTTP обработчики с mock сервисами
- Config: загрузка конфигурации

### Интеграционные тесты:
- Полный workflow: проверка → сохранение → отчет
- HTTP API: реальные запросы через httptest
- Concurrent операции: thread safety
- Обработка ошибок и невалидных данных

## Если что-то не работает:
1. Проверьте, что ссылки в запросе указаны без http:// или https://

2. Убедитесь, что сервер запущен (curl http://localhost:8080)

3. Если нужен отчет по старым проверкам - их ID можно найти в логах при проверке