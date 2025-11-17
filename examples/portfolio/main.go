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

	// Example 1: Get single wallet portfolio (V1)
	fmt.Println("=== Single Wallet Portfolio (V1) ===")
	walletAddress := "0x1234567890123456789012345678901234567890" // Replace with actual wallet

	portfolio, err := client.GetWalletPortfolioV1(ctx, walletAddress)
	if err != nil {
		log.Fatalf("Failed to get portfolio: %v", err)
	}

	fmt.Printf("Total Portfolio Value: $%.2f\n", portfolio.TotalUSDValue)
	fmt.Printf("Number of Tokens: %d\n\n", len(portfolio.Portfolio))

	// Show top 5 tokens by value
	fmt.Println("Top 5 Tokens by Value:")
	for i, token := range portfolio.Portfolio {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s (%s)\n", i+1, token.TokenName, token.TokenSymbol)
		fmt.Printf("   Balance: %.6f\n", token.Balance)
		fmt.Printf("   USD Value: $%.2f\n", token.TotalUSDValue)
		fmt.Printf("   Price: $%.6f\n\n", token.TokenPriceUSD)
	}

	// Example 2: Get aggregated multi-wallet portfolio (V2)
	fmt.Println("\n=== Aggregated Multi-Wallet Portfolio (V2) ===")
	wallets := []string{
		"0x1234567890123456789012345678901234567890",
		"0x0987654321098765432109876543210987654321",
		// Add more wallets as needed
	}

	portfolioV2, err := client.GetWalletPortfolioV2(ctx, &apiv1.WalletPortfolioV2Request{
		Wallets: wallets,
	})
	if err != nil {
		log.Fatalf("Failed to get portfolio v2: %v", err)
	}

	fmt.Printf("Aggregated Total Value: $%.2f\n", portfolioV2.TotalUSDValue)
	fmt.Printf("Total Tokens Across All Wallets: %d\n\n", len(portfolioV2.Portfolio))

	// Show chain distribution
	fmt.Println("Chain Distribution:")
	for _, chain := range portfolioV2.Chains {
		fmt.Printf("  %s: $%.2f\n", chain.Chain, chain.TotalUSD)
	}

	// Example 3: Get specific token balance for Solana wallet (V2)
	fmt.Println("\n=== Specific Token Balance (Solana Only) ===")
	solanaWallet := "SOL_WALLET_ADDRESS_HERE" // Replace with actual Solana wallet
	tokenMint := "TOKEN_MINT_ADDRESS_HERE"    // Replace with actual token mint

	tokenBalance, err := client.GetWalletPortfolioV2(ctx, &apiv1.WalletPortfolioV2Request{
		Wallets: []string{solanaWallet},
		Token:   &tokenMint,
	})
	if err != nil {
		log.Printf("Failed to get specific token balance: %v", err)
	} else if len(tokenBalance.Portfolio) > 0 {
		token := tokenBalance.Portfolio[0]
		fmt.Printf("Token: %s (%s)\n", token.TokenName, token.TokenSymbol)
		fmt.Printf("Balance: %.6f\n", token.Balance)
		fmt.Printf("USD Value: $%.2f\n", token.TotalUSDValue)
	}

	fmt.Println("\n=== Portfolio Analysis Complete ===")
	fmt.Printf("Total API Credits Used: ~%d credits\n", 20+20*len(wallets)+20)
}
