package main

import (
	"log"

	"github.com/ddProgerGo/task-kaspi/internal/handler"
	"github.com/ddProgerGo/task-kaspi/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных", err)
	}
	
	database.RunMigrations(db)

	r := gin.Default()

	handler.RegisterRoutes(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера", err)
	}
}
