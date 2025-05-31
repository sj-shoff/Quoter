# Quoter - Мини-сервис "Цитатник"

REST API-сервис на Go для хранения и управления цитатами. Сервис позволяет добавлять, просматривать, фильтровать и удалять цитаты.

## Функциональные возможности

1. ✅ Добавление новой цитаты (POST /quotes)
2. ✅ Получение всех цитат (GET /quotes)
3. ✅ Получение случайной цитаты (GET /quotes/random)
4. ✅ Фильтрация по автору (GET /quotes?author=Confucius)
5. ✅ Удаление цитаты по ID (DELETE /quotes/{id})

## Технические особенности

- Хранение данных в памяти (данные сохраняются до перезапуска сервиса)
- Использованы только стандартные библиотеки Go (net/http, encoding/json и др.)
- Чистая архитектура с разделением на слои (handler-service-storage)
- Юнит-тесты для всех компонентов системы
- Логирование операций с использованием стандартной библиотеки slog

## Технические требования

- Go 1.22 или выше
- Только стандартные библиотеки Go

## Запуск сервиса

1. Клонируйте репозиторий:
```bash
git clone https://github.com/sj-shoff/Quoter.git
```
2. Установите зависимости:
```bash
go mod download
```
3. Запустите сервис:
```bash
make run
```

Сервис будет доступен по адресу: `http://localhost:8080`

## Проверка функционала

### 1. Добавление новой цитаты
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

Пример ответа:
```json
{"id":1,"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}
```

### 2. Получение всех цитат
```bash
curl http://localhost:8080/quotes
```

Пример ответа:
```json
[{"id":1,"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}]
```

### 3. Получение случайной цитаты
```bash
curl http://localhost:8080/quotes/random
```

Пример ответа:
```json
{"id":1,"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}
```

### 4. Фильтрация по автору
```bash
curl http://localhost:8080/quotes?author=Confucius
```

Пример ответа:
```json
[{"id":1,"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}]
```

### 5. Удаление цитаты по ID
```bash
curl -X DELETE http://localhost:8080/quotes/1
```

Статус ответа при успешном удалении: 204 No Content

## Запуск тестов

Для запуска юнит-тестов выполните:

```bash
make -B tests
```

Тесты покрывают:
- HTTP-обработчики (handlers)
- Бизнес-логику (service)
- Хранилище данных (storage)
