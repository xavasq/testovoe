package main

import (
	"log"
	"testovoe/internal/database"
	"testovoe/internal/handler"
	"testovoe/internal/repository"
	"testovoe/internal/router"
	"testovoe/internal/service"
)

func main() {
	database.ConnectDB()
	defer database.CloseDB()

	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := router.SetupRouter(userHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("ошибка при запуске сервера: %v", err)
	}
}
