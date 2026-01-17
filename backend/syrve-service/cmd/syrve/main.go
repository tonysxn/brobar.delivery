package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/syrve-service/internal/api"
	"github.com/tonysanin/brobar/syrve-service/internal/services/syrve"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: No .env file found or failed to load")
	}

	apiLogin := helpers.GetEnv("SYRVE_TOKEN", "")
	organizationID := helpers.GetEnv("SYRVE_ORGANIZATION", "")
	port := helpers.GetEnv("SYRVE_PORT", "3011")

	client := syrve.NewClient(apiLogin, organizationID).WithTimeout(10)

	server := api.NewServer(client)

	log.Printf("Starting server on :%s", port)
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
