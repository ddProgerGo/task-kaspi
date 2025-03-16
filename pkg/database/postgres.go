package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectPostgres() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла:", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	log.Println("Подключение к базе данных с параметрами:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Успешное подключение к базе данных")
	return db, nil
}

func RunMigrations(db *sql.DB) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS people (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			iin VARCHAR(12) UNIQUE NOT NULL,
			phone VARCHAR(20) NOT NULL
		);`,
	}

	for _, query := range migrations {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Ошибка при выполнении миграции: %v", err)
		}
	}
	log.Println("Миграции успешно выполнены")
}
