package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gojo-op/todo-auth-api/internal/config"
	"github.com/gojo-op/todo-auth-api/internal/database"
	"github.com/gojo-op/todo-auth-api/internal/handler"
	"github.com/gojo-op/todo-auth-api/internal/repository"
	"github.com/gojo-op/todo-auth-api/internal/routes"
	"github.com/gojo-op/todo-auth-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	todoService := service.NewTodoService(todoRepo)

	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)

	router := gin.Default()
	router.SetTrustedProxies(nil)
	routes.Register(router, authService, authHandler, todoHandler)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("server started on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
