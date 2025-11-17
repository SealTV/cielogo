package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api/apiv1"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("CIELO_API_KEY")
	if apiKey == "" {
		log.Fatal("CIELO_API_KEY environment variable is required")
	}

	// Create client
	client := cielogo.NewClient(apiKey)
	ctx := context.Background()

	// Example 1: Add a new tracked wallet with notifications
	fmt.Println("=== Adding Tracked Wallet with Notifications ===")
	walletToTrack := "0x1234567890123456789012345678901234567890" // Replace with actual wallet

	addedWallet, err := client.AddTrackedWalletsV1(ctx, &apiv1.AddTrackedWalletRequest{
		Wallet:           walletToTrack,
		Label:            "Whale Wallet #1",
		MinAmountUSD:     apiv1.ToRef(1000.0),                                   // Only notify for transactions > $1000
		TelegramBotID:    apiv1.ToRef(1),                                        // Telegram bot ID
		DiscordChannelID: apiv1.ToRef("https://discord.com/api/webhooks/..."), // Discord webhook
		NewTrades:        apiv1.ToRef(true),
		Filters:          []int{1, 2}, // Filter IDs (swap, transfer, etc.)
		Chains:           []int{1, 2}, // Chain IDs (ethereum=1, solana=2, etc.)
	})
	if err != nil {
		log.Fatalf("Failed to add tracked wallet: %v", err)
	}

	fmt.Printf("Tracked Wallet ID: %d\n", addedWallet.ID)
	fmt.Printf("Wallet: %s\n", addedWallet.Wallet)
	fmt.Printf("Label: %s\n", addedWallet.Label)
	fmt.Println()

	// Example 2: Get all tracked wallets
	fmt.Println("=== Listing All Tracked Wallets ===")
	trackedWallets, err := client.GetTrackedWalletsV1(ctx, &apiv1.GetTrackedWalletsRequest{})
	if err != nil {
		log.Fatalf("Failed to get tracked wallets: %v", err)
	}

	fmt.Printf("Total Tracked Wallets: %d\n\n", len(trackedWallets.TrackedWallets))
	for i, wallet := range trackedWallets.TrackedWallets {
		if i >= 5 { // Show only first 5
			break
		}
		fmt.Printf("%d. %s\n", i+1, wallet.Label)
		fmt.Printf("   Wallet: %s\n", wallet.Wallet)
		if wallet.ListID != nil {
			fmt.Printf("   List ID: %d\n", *wallet.ListID)
		}
		fmt.Println()
	}

	// Example 3: Update tracked wallet using V2 (partial update)
	fmt.Println("=== Updating Tracked Wallet (V2 - Partial Update) ===")
	updatedWallet, err := client.UpdateTrackedWalletV2(ctx, walletToTrack, &apiv1.UpdateTrackedWalletV2Request{
		Label:  apiv1.ToRef("Updated Whale Wallet"),
		MinUSD: apiv1.ToRef(5000.0), // Increase threshold to $5000
		// Only update label and MinUSD, all other fields remain unchanged
	})
	if err != nil {
		log.Fatalf("Failed to update tracked wallet: %v", err)
	}

	fmt.Printf("Updated Label: %s\n", updatedWallet.Label)
	if updatedWallet.ListID != nil {
		fmt.Printf("List ID: %d\n", *updatedWallet.ListID)
	}
	fmt.Println()

	// Example 4: Get wallet by address
	fmt.Println("=== Get Tracked Wallet by Address ===")
	walletByAddress, err := client.GetWalletByAddressV1(ctx, walletToTrack)
	if err != nil {
		log.Fatalf("Failed to get wallet by address: %v", err)
	}

	fmt.Printf("Wallet ID: %d\n", walletByAddress.ID)
	fmt.Printf("Label: %s\n", walletByAddress.Label)
	fmt.Printf("Wallet Address: %s\n", walletByAddress.Wallet)
	if walletByAddress.ListID != nil {
		fmt.Printf("List ID: %d\n", *walletByAddress.ListID)
	}
	fmt.Println()

	// Example 5: Get available Telegram bots
	fmt.Println("=== Available Telegram Bots for Notifications ===")
	telegramBots, err := client.GetTelegramBotsV1(ctx)
	if err != nil {
		log.Printf("Failed to get telegram bots: %v", err)
	} else {
		fmt.Printf("Available Bots: %d\n\n", len(telegramBots.Bots))
		for i, bot := range telegramBots.Bots {
			fmt.Printf("%d. %s (@%s)\n", i+1, bot.Label, bot.Name)
			fmt.Printf("   Plan: %s | Default: %v\n", bot.Plan, bot.Default)
			if bot.Link != "" {
				fmt.Printf("   Link: %s\n", bot.Link)
			}
			fmt.Println()
		}
	}

	// Example 6: Remove tracked wallet
	fmt.Println("=== Removing Tracked Wallet ===")
	err = client.RemoveTrackedWalletsV1(ctx, &apiv1.RemoveTrackedWalletsRequest{
		WalletIDs: []int64{addedWallet.ID},
	})
	if err != nil {
		log.Fatalf("Failed to remove tracked wallet: %v", err)
	}

	fmt.Printf("Successfully removed wallet ID: %d\n", addedWallet.ID)

	fmt.Println("\n=== Tracked Wallets Management Complete ===")
	fmt.Println("Total API Credits Used: ~30 credits")
}
