package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/file-service/internal/api"
	"github.com/tonysanin/brobar/file-service/internal/services"
	"github.com/tonysanin/brobar/pkg/helpers"
)

func main() {
	_ = godotenv.Load(".env")

	port := helpers.GetEnv("SERVER_PORT", "3001")

	fileService := services.NewFileService("./uploads")
	server := api.NewServer(fileService)

	log.Printf("Starting server on :%s", port)
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
