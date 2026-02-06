package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/pkg/syrve"
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

	apiLogin := os.Getenv("SYRVE_TOKEN")

	// We might not have OrganizationID in env if we want to discover it.
	// But Client needs it optionally.
	// For CLI, we can fetch it.
	
	if apiLogin == "" {
		log.Fatal("SYRVE_TOKEN environment variable is required")
	}

	// Initialize client without Org ID initially if we want to fetch it
	client := syrve.NewClient(apiLogin, "")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Authenticate
	tokenResp, err := client.GetAccessToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}
	token := tokenResp.Token

	// Check if organization is provided in env
	envOrgID := os.Getenv("SYRVE_ORGANIZATION")
	var orgID string

	if envOrgID != "" {
		orgID = envOrgID
		fmt.Printf("Using Organization from ENV: %s\n", orgID)
	} else {
		// Get Organizations
		orgsResp, err := client.GetOrganizations(ctx, token, syrve.OrganizationsRequest{
			IncludeDisabled: false,
		})
		if err != nil {
			log.Fatalf("Failed to get organizations: %v", err)
		}

		if len(orgsResp.Organizations) == 0 {
			log.Fatal("No organizations found")
		}

		orgID = orgsResp.Organizations[0].ID
		fmt.Printf("Using Organization (auto-discovered): %s (%s)\n", orgsResp.Organizations[0].Name, orgID)
	}

	switch os.Args[1] {
	case "install":
		installCmd.Parse(os.Args[2:])
		url := *installURL
		if url == "" {
			// Try to construct from env if not provided
			host := os.Getenv("NGINX_DOMAIN")
			if host != "" {
				url = fmt.Sprintf("https://%s/api/webhooks/syrve", host)
			}
		}

		if url == "" {
			log.Fatal("Webhook URL is required. Provide it via -url flag or ensure NGINX_DOMAIN is set in .env")
		}

		fmt.Printf("Setting webhook to: %s\n", url)
		if err := client.UpdateWebhook(ctx, token, orgID, url); err != nil {
			log.Fatalf("Failed to set webhook: %v", err)
		}
		fmt.Println("Webhook was set successfully")

		// Verify settings
		settings, err := client.GetWebhookSettings(ctx, token, orgID)
		if err != nil {
			log.Printf("Failed to verify webhook settings: %v", err)
		} else {
			fmt.Printf("Current Webhook URL: %s\n", settings.WebHooksURI)
			if settings.AuthToken != "" {
				fmt.Printf("Auth Token is SET\n")
			} else {
				fmt.Printf("Auth Token is NOT set\n")
			}
		}

	case "remove":
		removeCmd.Parse(os.Args[2:])
		fmt.Println("Removing webhook...")
		if err := client.RemoveWebhook(ctx, token, orgID); err != nil {
			log.Fatalf("Failed to delete webhook: %v", err)
		}
		fmt.Println("Webhook removed successfully")

	default:
		fmt.Println("expected 'install' or 'remove' subcommands")
		os.Exit(1)
	}
}
