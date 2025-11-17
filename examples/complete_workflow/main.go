package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/api/chains"
)

// This example demonstrates a complete workflow for tracking and analyzing crypto wallets:
// 1. Create a wallet list
// 2. Add wallets to tracking with notifications
// 3. Analyze wallet performance
// 4. Monitor portfolio
// 5. Get transaction feed
// 6. Find related wallets
// 7. Clean up resources

func main() {
	apiKey := os.Getenv("CIELO_API_KEY")
	if apiKey == "" {
		log.Fatal("CIELO_API_KEY environment variable is required")
	}

	client := cielogo.NewClient(apiKey)
	ctx := context.Background()

	// Target wallet for analysis (replace with actual wallet)
	targetWallet := "0x1234567890123456789012345678901234567890"

	fmt.Println("=== Complete Wallet Analysis Workflow ===")
	fmt.Println()

	// Step 1: Create a wallet list
	fmt.Println("Step 1: Creating Wallet List...")
	walletList, err := client.AddWalletsListV1(ctx, &apiv1.AddWalletsListRequest{
		Name:        "High-Value Traders",
		Description: "List of wallets with significant trading activity",
		IsPublic:    false,
	})
	if err != nil {
		log.Fatalf("Failed to create wallet list: %v", err)
	}
	fmt.Printf("✓ Created wallet list: %s (ID: %d)\n\n", walletList.Name, walletList.ID)

	// Step 2: Add wallet to tracking with notifications
	fmt.Println("Step 2: Adding Wallet to Tracking...")
	trackedWallet, err := client.AddTrackedWalletsV1(ctx, &apiv1.AddTrackedWalletRequest{
		Wallet:           targetWallet,
		Label:            "Main Analysis Target",
		ListID:           &walletList.ID,
		MinAmountUSD:     apiv1.ToRef(1000.0),
		Filters:          []int{1, 2}, // Filter IDs for swap, transfer
		Chains:           []int{1, 2}, // Chain IDs (ethereum=1, solana=2)
		NewTrades:        apiv1.ToRef(true),
		TelegramBotID:    apiv1.ToRef(1),
		DiscordChannelID: apiv1.ToRef("https://discord.com/api/webhooks/..."),
	})
	if err != nil {
		log.Fatalf("Failed to add tracked wallet: %v", err)
	}
	fmt.Printf("✓ Added wallet to tracking (ID: %d)\n", trackedWallet.ID)
	fmt.Println("  Notifications enabled for transactions > $1000")
	fmt.Println()

	// Step 3: Analyze wallet trading performance
	fmt.Println("Step 3: Analyzing Trading Performance...")
	sevenDays := apiv1.Timeframe7Days
	tradingStats, err := client.GetTradingStatsV1(ctx, &apiv1.TradingStatsRequest{
		Wallet: targetWallet,
		Days:   &sevenDays,
	})
	if err != nil {
		log.Printf("⚠ Failed to get trading stats: %v\n", err)
	} else {
		fmt.Printf("✓ Trading Stats Retrieved:\n")
		fmt.Printf("  Total Profit: $%.2f\n", tradingStats.TotalProfit)
		fmt.Printf("  Win Rate: %.1f%%\n", tradingStats.WinRate*100)
		fmt.Printf("  Total Trades: %d\n", tradingStats.TotalTrades)
		fmt.Printf("  Average Profit: $%.2f\n\n", tradingStats.AverageProfit)
	}

	// Step 4: Get current portfolio
	fmt.Println("Step 4: Fetching Current Portfolio...")
	portfolio, err := client.GetWalletPortfolioV1(ctx, targetWallet)
	if err != nil {
		log.Fatalf("Failed to get portfolio: %v", err)
	}
	fmt.Printf("✓ Portfolio Retrieved:\n")
	fmt.Printf("  Total Value: $%.2f\n", portfolio.TotalUSDValue)
	fmt.Printf("  Number of Tokens: %d\n", len(portfolio.Portfolio))
	if len(portfolio.Portfolio) > 0 {
		fmt.Printf("  Top Token: %s (%.2f %s, $%.2f)\n\n",
			portfolio.Portfolio[0].TokenSymbol,
			portfolio.Portfolio[0].Balance,
			portfolio.Portfolio[0].TokenSymbol,
			portfolio.Portfolio[0].TotalUSDValue)
	}

	// Step 5: Get token PnL with only active positions
	fmt.Println("Step 5: Analyzing Token PnL (Active Positions Only)...")
	tokenPnL, err := client.GetTokensPnlV1(ctx, &apiv1.TokensPnLRequest{
		Wallet:              targetWallet,
		Chains:              []chains.ChainType{chains.Ethereum},
		ActivePositionsOnly: apiv1.ToRef(true),
	})
	if err != nil {
		log.Printf("⚠ Failed to get token PnL: %v\n", err)
	} else {
		fmt.Printf("✓ Active Positions: %d\n", len(tokenPnL.Items))
		if len(tokenPnL.Items) > 0 {
			fmt.Printf("  Best Performer: %s (PnL: $%.2f, ROI: %.2f%%)\n\n",
				tokenPnL.Items[0].Symbol,
				tokenPnL.Items[0].TotalPnlUSD,
				tokenPnL.Items[0].RoiPercentage)
		}
	}

	// Step 6: Get recent transaction feed
	fmt.Println("Step 6: Fetching Recent Transaction Feed...")
	maxUSD := 100000.0
	listID := int(walletList.ID)
	feed, err := client.GetFeedV1(ctx, &apiv1.FeedRequest{
		List:             &listID,
		MaxUSD:           &maxUSD,
		IncludeMarketCap: apiv1.ToRef(false), // Keep costs down
	})
	if err != nil {
		log.Printf("⚠ Failed to get feed: %v\n", err)
	} else {
		fmt.Printf("✓ Feed Retrieved: %d recent transactions\n", len(feed.Items))
		if len(feed.Items) > 0 {
			fmt.Printf("  Latest: %s on %s (%s)\n\n",
				feed.Items[0].TxType,
				feed.Items[0].Chain,
				time.Unix(feed.Items[0].Timestamp, 0).Format("2006-01-02 15:04"))
		}
	}

	// Step 7: Find related wallets
	fmt.Println("Step 7: Discovering Related Wallets...")
	sortCriteria := apiv1.RelatedWalletsSortingInflowDesc
	relatedWallets, err := client.GetRelatedWalletsV1(ctx, &apiv1.RelatedWalletsRequest{
		Wallet:       targetWallet,
		SortCriteria: &sortCriteria,
	})
	if err != nil {
		log.Printf("⚠ Failed to get related wallets: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d related wallets\n", len(relatedWallets.RelatedWallets))
		for i, wallet := range relatedWallets.RelatedWallets {
			if i >= 3 { // Show top 3
				break
			}
			fmt.Printf("  %d. %s - Inflow: $%.2f | Outflow: $%.2f\n",
				i+1, wallet.Wallet, wallet.InflowUSD, wallet.OutflowUSD)
		}
		fmt.Println()
	}

	// Step 8: Get aggregated stats
	fmt.Println("Step 8: Getting Aggregated Statistics...")
	aggStats, err := client.GetAggregatedTokenPnLV1(ctx, &apiv1.AggregatedTokenPnLRequest{
		Wallet: targetWallet,
		Chains: []chains.ChainType{chains.Ethereum, chains.Solana},
	})
	if err != nil {
		log.Printf("⚠ Failed to get aggregated stats: %v\n", err)
	} else {
		fmt.Printf("✓ Aggregated Stats:\n")
		fmt.Printf("  Tokens Traded: %d\n", aggStats.TokensTraded)
		fmt.Printf("  Realized PnL: $%.2f\n", aggStats.RealizedPnlUsd)
		fmt.Printf("  Combined ROI: %.2f%%\n\n", aggStats.CombinedRoiPercentage)
	}

	// Step 9: Summary and cleanup prompt
	fmt.Println("=== Workflow Summary ===")
	fmt.Printf("Wallet List: %s (ID: %d)\n", walletList.Name, walletList.ID)
	fmt.Printf("Tracked Wallet: %s (ID: %d)\n", trackedWallet.Label, trackedWallet.ID)
	fmt.Printf("Portfolio Value: $%.2f\n", portfolio.TotalUSDValue)
	if tradingStats != nil {
		fmt.Printf("Total Profit: $%.2f | Win Rate: %.1f%%\n", tradingStats.TotalProfit, tradingStats.WinRate*100)
	}

	fmt.Println("\n=== Estimated API Credits Used ===")
	fmt.Println("Wallet List Creation: 5 credits")
	fmt.Println("Add Tracked Wallet: 5 credits")
	fmt.Println("Trading Stats: 30 credits")
	fmt.Println("Portfolio: 20 credits")
	fmt.Println("Token PnL: 5 credits")
	fmt.Println("Feed: 5 credits")
	fmt.Println("Related Wallets: 10 credits")
	fmt.Println("Aggregated Stats: 20 credits")
	fmt.Println("TOTAL: ~100 credits")

	fmt.Println("\n=== Cleanup (optional) ===")
	fmt.Println("To remove the tracked wallet and list, uncomment the cleanup code below:")
	fmt.Println("// err = client.RemoveTrackedWalletsV1(ctx, &apiv1.RemoveTrackedWalletsRequest{IDs: []int{int(trackedWallet.ID)}})")
	fmt.Println("// err = client.DeleteWalletsListV1(ctx, walletList.ID, true)")

	// Uncomment to clean up resources
	/*
		fmt.Println("\nCleaning up resources...")
		err = client.RemoveTrackedWalletsV1(ctx, &apiv1.RemoveTrackedWalletsRequest{
			IDs: []int{int(trackedWallet.ID)},
		})
		if err != nil {
			log.Printf("Failed to remove tracked wallet: %v", err)
		} else {
			fmt.Println("✓ Removed tracked wallet")
		}

		err = client.DeleteWalletsListV1(ctx, walletList.ID, true)
		if err != nil {
			log.Printf("Failed to delete wallet list: %v", err)
		} else {
			fmt.Println("✓ Deleted wallet list")
		}
	*/

	fmt.Println("\n=== Complete Workflow Finished ===")
}
