package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"testovoe/internal/config"
)

var DB *pgxpool.Pool

func ConnectDB() {
	cfg := config.LoadEnv()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatalf("ошибка при подключении к базе данных: %v", err)
	}
	DB = pool
	fmt.Println("база данных успешно подключена")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("подключение к базе данных закрыто")
	}
}
