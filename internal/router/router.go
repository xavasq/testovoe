package router

import (
	"github.com/gin-gonic/gin"
	"testovoe/internal/handler"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/users")
	{
		api.POST("/", userHandler.CreateUser)
		api.GET("/:id", userHandler.GetUserByID)
		api.PUT("/:id", userHandler.UpdateUserByID)
		api.DELETE("/:id", userHandler.DeleteUserByID)
	}

	return r
}
