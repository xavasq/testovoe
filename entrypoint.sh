#!/bin/sh
# Ожидаем, пока база данных станет доступной
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"; do
  echo "Ожидаем базу данных..."
  sleep 1
done

echo "База готова, применяем миграции..."
# Проверяем, что migrate работает и выводим результат
/usr/local/bin/migrate -path /root/db/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up || {
  echo "Ошибка при применении миграций!"
  exit 1
}

echo "Миграции успешно применены, запускаем приложение..."
# Запускаем сервер
/root/main