package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gojo-op/todo-auth-api/internal/handler"
	"github.com/gojo-op/todo-auth-api/internal/middleware"
	"github.com/gojo-op/todo-auth-api/internal/service"
)

func Register(router *gin.Engine, authService *service.AuthService, authHandler *handler.AuthHandler, todoHandler *handler.TodoHandler) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		todos := api.Group("/todos")
		todos.Use(middleware.Auth(authService))
		{
			todos.GET("", todoHandler.List)
			todos.POST("", todoHandler.Create)
			todos.GET("/:id", todoHandler.Get)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
		}
	}
}
