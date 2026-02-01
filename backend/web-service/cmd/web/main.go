package main

import (
	"context"
	"log"
	"time"

	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/web-service/internal/api"
	"github.com/tonysanin/brobar/web-service/internal/repositories"
	"github.com/tonysanin/brobar/web-service/internal/services"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env")

	port := helpers.GetEnv("WEB_PORT", "3006")

	db, err := InitDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Close database failed: %v", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	settingRepo := repositories.NewSettingRepository(db)
	settingService := services.NewSettingService(settingRepo)

	reviewRepo := repositories.NewReviewRepository(db)
	reviewService := services.NewReviewService(reviewRepo)

	server := api.NewServer(settingService, reviewService)

	log.Printf("Starting server on :%s", port)
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func InitDatabase() (*sqlx.DB, error) {
	connStr := "postgres://" +
		helpers.GetEnv("DB_USER", "") + ":" +
		helpers.GetEnv("DB_PASSWORD", "") + "@" +
		helpers.GetEnv("DB_HOST", "") + ":" +
		helpers.GetEnv("DB_PORT", "") + "/" +
		helpers.GetEnv("DB_NAME", "") + "?sslmode=" +
		helpers.GetEnv("DB_SSLMODE", "disable")

	return sqlx.Connect("postgres", connStr)
}
