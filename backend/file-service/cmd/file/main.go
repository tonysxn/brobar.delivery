package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/file-service/internal/api"
	"github.com/tonysanin/brobar/file-service/internal/config"
	"github.com/tonysanin/brobar/file-service/internal/services"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.NewConfig()

	fileService := services.NewFileService(cfg.UploadDir)
	server := api.NewServer(fileService)

	log.Printf("Starting server on :%s", cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
