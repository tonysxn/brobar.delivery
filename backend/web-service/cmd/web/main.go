package main

import (
	"context"
	"log"
	"time"

	"github.com/tonysanin/brobar/web-service/internal/api"
	"github.com/tonysanin/brobar/web-service/internal/config"
	"github.com/tonysanin/brobar/web-service/internal/repositories"
	"github.com/tonysanin/brobar/web-service/internal/services"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.NewConfig()

	db, err := InitDatabase(cfg)
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

	server := api.NewServer(cfg, settingService, reviewService)

	log.Printf("Starting server on :%s", cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", cfg.GetDatabaseURL())
}
