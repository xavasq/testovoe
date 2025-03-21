# тестовое задание: REST API на Go

REST API для создания, получения и обновления пользователей. написано на Go с PostgreSQL, покрыто тестами и упаковано в Docker.

## стек
- **язык**: Go 1.22
- **фреймворк**: Gin
- **база**: PostgreSQL 17 (alpine)
- **драйвер**: pgx
- **миграции**: golang-migrate
- **тесты**: testify с моками
- **контейнеризация**: Docker и docker-compose

## структура
- `cmd/main.go` — точка входа
- `internal/config` — загрузка конфигурации
- `internal/database` — подключение к базе
- `internal/repository` — работа с данными
- `internal/service` — бизнес-логика
- `internal/handler` — обработка HTTP
- `internal/router` — настройка роутов
- `db/migrations` — миграции базы

## API
- `POST /users` — создать юзера  
  **тело**: `{"name": "string", "email": "string"}`  
  **ответ**: 201 Created, данные юзера
- `GET /users/:id` — получить юзера  
  **ответ**: 200 OK, данные юзера или 404 Not Found
- `PUT /users/:id` — обновить юзера  
  **тело**: `{"name": "string", "email": "string"}`  
  **ответ**: 200 OK, сообщение об успехе

## запуск

### локально
1. установи Go (1.22+), PostgreSQL и golang-migrate: