package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddProgerGo/task-kaspi/internal/handler"
	"github.com/ddProgerGo/task-kaspi/internal/middleware"
	"github.com/ddProgerGo/task-kaspi/internal/repository"
	"github.com/ddProgerGo/task-kaspi/internal/service"
	"github.com/ddProgerGo/task-kaspi/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	_ "github.com/ddProgerGo/task-kaspi/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	db, err := database.ConnectPostgres()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}

	cache := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})

	_, err = cache.Ping(context.Background()).Result()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to Redis")
	}

	log.Println("Redis подключен")

	database.RunMigrations(db)

	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Database is not reachable")
	}

	logger.Info("Connected to database successfully")

	repo := repository.NewPersonRepository(db, logger, cache)
	service := service.NewPersonService(repo, logger, cache)
	handler := handler.NewPersonHandler(service, logger)

	router := gin.Default()

	router.Use(middleware.ErrorHandlingMiddleware(logger))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/iin_check/:iin", handler.CheckIIN)
	router.POST("/people/info", handler.SavePerson)
	router.GET("/people/info/iin/:iin", handler.GetPersonByIIN)
	router.GET("/people/info/phone/:name", handler.GetPeopleByName)

	server := &http.Server{
		Addr:    os.Getenv("ADDRESS"),
		Handler: router,
	}

	go func() {
		logger.Info("Server is starting on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server startup failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Warn("Shutdown signal received, stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server shutdown failed")
	} else {
		logger.Info("Server stopped gracefully")
	}

	if err := db.Close(); err != nil {
		logger.WithError(err).Error("Error closing database connection")
	} else {
		logger.Info("Database connection closed")
	}

}
