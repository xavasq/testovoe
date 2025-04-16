# REST API для управления пользователями на Go

Этот проект представляет собой REST API для создания, редактирования и удаления пользователей. Он написан на языке Go с использованием PostgreSQL в качестве базы данных. Проект контейнеризирован с помощью Docker и включает в себя миграции базы данных, модульные тесты и поддержку окружения через `.env`.

---

## 🛠️ Стек технологий

- **Язык программирования:** Go
- **База данных:** PostgreSQL
- **Контейнеризация:** Docker, Docker Compose
- **Фреймворк:** Gin
- **Миграции:** golang-migrate
- **Тестирование:**
  - Встроенный пакет `testing` Go
  - Библиотека `testify` для утверждений и моков
  - `testcontainers-go` для интеграционных тестов с PostgreSQL
---

## 🚀 Запуск проекта

### 1. Клонирование репозитория

```bash
git clone https://github.com/yourusername/your-repo-name.git
cd your-repo-name

2. Настройка окружения
Создайте файл .env на основе .env.example и заполните его своими данными:
cp .env.example .env

3. Запуск Docker
Для запуска проекта используйте Docker Compose:
docker-compose up --build
После запуска проект будет доступен по адресу http://localhost:8080.

📚 Методы API
Создание пользователя
Метод: POST /users

Тело запроса:
{
  "name": "Иван",
  "email": "ivan@example.com"
}

Получение информации о пользователе
Метод: GET /users/{id}

Ответ:
{
  "id": 1,
  "name": "Иван",
  "email": "ivan@example.com"
}

Обновление данных пользователя
Метод: PUT /users/{id}

Тело запроса:
{
  "name": "Иван Иванов",
  "email": "ivan.ivanov@example.com"
}
🧪 Тестирование

Для запуска модульных тестов выполните команду:
go test ./...

🐳 Docker
Проект полностью контейнеризирован. Для запуска достаточно выполнить:
docker-compose up --build

Это поднимет:
Сервер на Go (порт 8080)
PostgreSQL (порт 5432)

📁 Структура проекта (основное)
.
├── cmd
│   └── main.go              # Точка входа
├── db
│   └── migrations           # Миграции базы данных
│       ├── 000001_create_users_table.up.sql
│       └── 000001_create_users_table.down.sql
├── internal
│   ├── config               # Конфигурация приложения
│   │   └── config.go
│   ├── database             # Подключение к базе данных
│   │   └── db.go
│   ├── handler              # Обработчики HTTP-запросов
│   │   ├── user_handler.go
│   │   └── user_handler_test.go
│   ├── models               # Модели данных
│   │   └── user.go
│   ├── repository           # Логика работы с базой данных
│   │   ├── user_repository.go
│   │   └── user_repository_test.go
│   ├── router               # Маршрутизация
│   │   └── router.go
│   └── service              # Бизнес-логика
│       ├── user_service.go
│       └── user_service_test.go
├── .env.example             # Пример файла окружения
├── docker-compose.yml       # Docker Compose конфигурация
├── Dockerfile               # Dockerfile для сборки образа
├── go.mod                   # Модули Go
└── README.md                # Документация

