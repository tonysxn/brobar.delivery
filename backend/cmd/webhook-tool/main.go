package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/pkg/telegram"
)

func main() {
	// Try to load .env file from varied locations
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	installURL := installCmd.String("url", "", "Webhook URL")

	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'install' or 'remove' subcommands")
		os.Exit(1)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	client := telegram.NewClient(token)

	switch os.Args[1] {
	case "install":
		installCmd.Parse(os.Args[2:])
		url := *installURL
		if url == "" {
			// Try to construct from env if not provided
			host := os.Getenv("NGINX_DOMAIN")
			if host != "" {
				url = fmt.Sprintf("https://%s/api/webhooks/telegram", host)
			}
		}

		if url == "" {
			log.Fatal("Webhook URL is required. Provide it via -url flag or ensure NGINX_DOMAIN is set in .env")
		}

		fmt.Printf("Setting webhook to: %s\n", url)
		if err := client.SetWebhook(url); err != nil {
			log.Fatalf("Failed to set webhook: %v", err)
		}
		fmt.Println("Webhook was set successfully")

	case "remove":
		removeCmd.Parse(os.Args[2:])
		fmt.Println("Removing webhook...")
		if err := client.DeleteWebhook(); err != nil {
			log.Fatalf("Failed to delete webhook: %v", err)
		}
		fmt.Println("Webhook removed successfully")

	default:
		fmt.Println("expected 'install' or 'remove' subcommands")
		os.Exit(1)
	}
}
