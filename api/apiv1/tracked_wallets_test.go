package apiv1_test

import (
	"encoding/json"
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		typ      apiv1.WalletType
		expected string
	}{
		{"unknown", apiv1.UnknownWalletType, "unknown"},
		{"evm", apiv1.EvmWalletType, "evm"},
		{"solana", apiv1.SolanaWalletType, "solana"},
		{"dydx", apiv1.DydxWalletType, "dydx"},
		{"bitcoin", apiv1.BitcoinWalletType, "bitcoin"},
		{"tron", apiv1.TronWalletType, "tron"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.typ))
		})
	}
}

func TestUpdateTrackedWalletRequest_JSON(t *testing.T) {
	listID := int64(123)
	req := apiv1.UpdateTrackedWalletRequest{
		Wallet: "0x1234567890abcdef1234567890abcdef12345678",
		Label:  "My Wallet",
		ListID: &listID,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded apiv1.UpdateTrackedWalletRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Wallet, decoded.Wallet)
	assert.Equal(t, req.Label, decoded.Label)
	require.NotNil(t, decoded.ListID)
	assert.Equal(t, *req.ListID, *decoded.ListID)
}

func TestUpdateTrackedWalletRequest_JSON_NoListID(t *testing.T) {
	req := apiv1.UpdateTrackedWalletRequest{
		Wallet: "0x1234567890abcdef1234567890abcdef12345678",
		Label:  "My Wallet",
		ListID: nil,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	// Verify that list_id is omitted when nil
	assert.NotContains(t, string(data), "list_id")

	var decoded apiv1.UpdateTrackedWalletRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Wallet, decoded.Wallet)
	assert.Equal(t, req.Label, decoded.Label)
	assert.Nil(t, decoded.ListID)
}

func TestUpdateTrackedWalletV2Request_JSON_FullFields(t *testing.T) {
	label := "Updated Wallet"
	listID := 456
	minUSD := 100.50
	newTrades := true
	telegramBot := "bot123"
	discordChannel := "channel456"

	req := apiv1.UpdateTrackedWalletV2Request{
		Label:          &label,
		ListID:         &listID,
		MinUSD:         &minUSD,
		TxTypes:        []string{"swap", "transfer"},
		Chains:         []string{"ethereum", "solana"},
		NewTrades:      &newTrades,
		TelegramBot:    &telegramBot,
		DiscordChannel: &discordChannel,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded apiv1.UpdateTrackedWalletV2Request
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.NotNil(t, decoded.Label)
	assert.Equal(t, *req.Label, *decoded.Label)
	require.NotNil(t, decoded.ListID)
	assert.Equal(t, *req.ListID, *decoded.ListID)
	require.NotNil(t, decoded.MinUSD)
	assert.InDelta(t, *req.MinUSD, *decoded.MinUSD, 0.01)
	assert.Equal(t, req.TxTypes, decoded.TxTypes)
	assert.Equal(t, req.Chains, decoded.Chains)
	require.NotNil(t, decoded.NewTrades)
	assert.Equal(t, *req.NewTrades, *decoded.NewTrades)
	require.NotNil(t, decoded.TelegramBot)
	assert.Equal(t, *req.TelegramBot, *decoded.TelegramBot)
	require.NotNil(t, decoded.DiscordChannel)
	assert.Equal(t, *req.DiscordChannel, *decoded.DiscordChannel)
}

func TestUpdateTrackedWalletV2Request_JSON_PartialFields(t *testing.T) {
	label := "Only Label Updated"
	req := apiv1.UpdateTrackedWalletV2Request{
		Label: &label,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	// Verify that only label is present and other fields are omitted
	jsonStr := string(data)
	assert.Contains(t, jsonStr, "label")
	assert.NotContains(t, jsonStr, "list_id")
	assert.NotContains(t, jsonStr, "min_usd")
	assert.NotContains(t, jsonStr, "telegram_bot")

	var decoded apiv1.UpdateTrackedWalletV2Request
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.NotNil(t, decoded.Label)
	assert.Equal(t, *req.Label, *decoded.Label)
	assert.Nil(t, decoded.ListID)
	assert.Nil(t, decoded.MinUSD)
	assert.Nil(t, decoded.NewTrades)
	assert.Nil(t, decoded.TelegramBot)
	assert.Nil(t, decoded.DiscordChannel)
}

func TestUpdateTrackedWalletV2Request_JSON_EmptyFields(t *testing.T) {
	req := apiv1.UpdateTrackedWalletV2Request{}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	// Verify that all optional fields are omitted when nil
	jsonStr := string(data)
	assert.NotContains(t, jsonStr, "label")
	assert.NotContains(t, jsonStr, "list_id")
	assert.NotContains(t, jsonStr, "min_usd")

	var decoded apiv1.UpdateTrackedWalletV2Request
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Nil(t, decoded.Label)
	assert.Nil(t, decoded.ListID)
	assert.Nil(t, decoded.MinUSD)
}

func TestTelegramBot_JSON_IntID(t *testing.T) {
	jsonData := `{
		"id": 12345,
		"label": "My Bot",
		"name": "bot_name",
		"disabled": false,
		"plan": "pro",
		"link": "https://t.me/mybot",
		"default": true,
		"is_custom": false,
		"usage": "notifications"
	}`

	var bot apiv1.TelegramBot
	err := json.Unmarshal([]byte(jsonData), &bot)
	require.NoError(t, err)

	assert.Equal(t, "My Bot", bot.Label)
	assert.Equal(t, "bot_name", bot.Name)
	assert.False(t, bot.Disabled)
	assert.Equal(t, "pro", bot.Plan)
	assert.Equal(t, "https://t.me/mybot", bot.Link)
	assert.True(t, bot.Default)
	assert.False(t, bot.IsCustom)
	assert.Equal(t, "notifications", bot.Usage)

	// Check that ID can be either int or string
	require.NotNil(t, bot.ID)
}

func TestTelegramBot_JSON_StringID(t *testing.T) {
	jsonData := `{
		"id": "bot-abc-123",
		"label": "Custom Bot",
		"name": "custom_bot",
		"disabled": true,
		"plan": "free",
		"link": "https://t.me/custombot",
		"default": false,
		"is_custom": true,
		"usage": "alerts"
	}`

	var bot apiv1.TelegramBot
	err := json.Unmarshal([]byte(jsonData), &bot)
	require.NoError(t, err)

	assert.Equal(t, "Custom Bot", bot.Label)
	assert.Equal(t, "custom_bot", bot.Name)
	assert.True(t, bot.Disabled)
	assert.Equal(t, "free", bot.Plan)
	assert.Equal(t, "https://t.me/custombot", bot.Link)
	assert.False(t, bot.Default)
	assert.True(t, bot.IsCustom)
	assert.Equal(t, "alerts", bot.Usage)

	// Check that ID can be either int or string
	require.NotNil(t, bot.ID)
}

func TestTelegramBot_JSON_Marshal(t *testing.T) {
	bot := apiv1.TelegramBot{
		ID:       "test-123",
		Label:    "Test Bot",
		Name:     "test_bot",
		Disabled: false,
		Plan:     "premium",
		Link:     "https://t.me/testbot",
		Default:  true,
		IsCustom: false,
		Usage:    "testing",
	}

	data, err := json.Marshal(bot)
	require.NoError(t, err)

	var decoded apiv1.TelegramBot
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, bot.Label, decoded.Label)
	assert.Equal(t, bot.Name, decoded.Name)
	assert.Equal(t, bot.Disabled, decoded.Disabled)
	assert.Equal(t, bot.Plan, decoded.Plan)
	assert.Equal(t, bot.Link, decoded.Link)
	assert.Equal(t, bot.Default, decoded.Default)
	assert.Equal(t, bot.IsCustom, decoded.IsCustom)
	assert.Equal(t, bot.Usage, decoded.Usage)
}

func TestGetTelegramBotsResponse_JSON(t *testing.T) {
	response := apiv1.GetTelegramBotsResponse{
		Bots: []apiv1.TelegramBot{
			{
				ID:       123,
				Label:    "Bot 1",
				Name:     "bot1",
				Disabled: false,
				Plan:     "pro",
				Link:     "https://t.me/bot1",
				Default:  true,
				IsCustom: false,
				Usage:    "main",
			},
			{
				ID:       "bot-2",
				Label:    "Bot 2",
				Name:     "bot2",
				Disabled: true,
				Plan:     "free",
				Link:     "https://t.me/bot2",
				Default:  false,
				IsCustom: true,
				Usage:    "secondary",
			},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.GetTelegramBotsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.Len(t, decoded.Bots, 2)
	assert.Equal(t, response.Bots[0].Label, decoded.Bots[0].Label)
	assert.Equal(t, response.Bots[1].Label, decoded.Bots[1].Label)
}

func TestTrackedWallet_JSON(t *testing.T) {
	listID := int64(999)
	wallet := apiv1.TrackedWallet{
		ID:     12345,
		Wallet: "0xabcdef1234567890abcdef1234567890abcdef12",
		Label:  "Test Wallet",
		Type:   apiv1.EvmWalletType,
		ListID: &listID,
	}

	data, err := json.Marshal(wallet)
	require.NoError(t, err)

	var decoded apiv1.TrackedWallet
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, wallet.ID, decoded.ID)
	assert.Equal(t, wallet.Wallet, decoded.Wallet)
	assert.Equal(t, wallet.Label, decoded.Label)
	assert.Equal(t, wallet.Type, decoded.Type)
	require.NotNil(t, decoded.ListID)
	assert.Equal(t, *wallet.ListID, *decoded.ListID)
}

func TestAddTrackedWalletRequest_JSON_WithNotifications(t *testing.T) {
	bundleID := int64(100)
	minAmountUSD := 50.75
	newTrades := true
	telegramBotID := 42
	discordChannelID := "channel-123"

	req := apiv1.AddTrackedWalletRequest{
		Wallet:           "0x1234567890abcdef1234567890abcdef12345678",
		Label:            "Whale Wallet",
		BundleID:         &bundleID,
		MinAmountUSD:     &minAmountUSD,
		Filters:          []int{1, 2, 3},
		Chains:           []int{1, 56, 137},
		NewTrades:        &newTrades,
		TelegramBotID:    &telegramBotID,
		DiscordChannelID: &discordChannelID,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded apiv1.AddTrackedWalletRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Wallet, decoded.Wallet)
	assert.Equal(t, req.Label, decoded.Label)
	require.NotNil(t, decoded.BundleID)
	assert.Equal(t, *req.BundleID, *decoded.BundleID)
	require.NotNil(t, decoded.MinAmountUSD)
	assert.InDelta(t, *req.MinAmountUSD, *decoded.MinAmountUSD, 0.01)
	assert.Equal(t, req.Filters, decoded.Filters)
	assert.Equal(t, req.Chains, decoded.Chains)
	require.NotNil(t, decoded.NewTrades)
	assert.Equal(t, *req.NewTrades, *decoded.NewTrades)
	require.NotNil(t, decoded.TelegramBotID)
	assert.Equal(t, *req.TelegramBotID, *decoded.TelegramBotID)
	require.NotNil(t, decoded.DiscordChannelID)
	assert.Equal(t, *req.DiscordChannelID, *decoded.DiscordChannelID)
}

func TestAddTrackedWalletRequest_JSON_PartialNotifications(t *testing.T) {
	minAmountUSD := 100.0
	req := apiv1.AddTrackedWalletRequest{
		Wallet:       "0x1234567890abcdef1234567890abcdef12345678",
		Label:        "Test Wallet",
		MinAmountUSD: &minAmountUSD,
		Chains:       []int{1},
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	// Verify that only specified fields are present
	jsonStr := string(data)
	assert.Contains(t, jsonStr, "min_amount_usd")
	assert.Contains(t, jsonStr, "chains")
	assert.NotContains(t, jsonStr, "bundle_id")
	assert.NotContains(t, jsonStr, "telegram_bot_id")
	assert.NotContains(t, jsonStr, "discord_channel_id")

	var decoded apiv1.AddTrackedWalletRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Wallet, decoded.Wallet)
	require.NotNil(t, decoded.MinAmountUSD)
	assert.Equal(t, req.Chains, decoded.Chains)
	assert.Nil(t, decoded.BundleID)
	assert.Nil(t, decoded.NewTrades)
	assert.Nil(t, decoded.TelegramBotID)
	assert.Nil(t, decoded.DiscordChannelID)
}

func TestAddTrackedWalletRequest_JSON_NoNotifications(t *testing.T) {
	req := apiv1.AddTrackedWalletRequest{
		Wallet: "0x1234567890abcdef1234567890abcdef12345678",
		Label:  "Simple Wallet",
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	// Verify that notification fields are omitted
	jsonStr := string(data)
	assert.NotContains(t, jsonStr, "bundle_id")
	assert.NotContains(t, jsonStr, "min_amount_usd")
	assert.NotContains(t, jsonStr, "filters")
	assert.NotContains(t, jsonStr, "chains")
	assert.NotContains(t, jsonStr, "new_trades")
	assert.NotContains(t, jsonStr, "telegram_bot_id")
	assert.NotContains(t, jsonStr, "discord_channel_id")

	var decoded apiv1.AddTrackedWalletRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Wallet, decoded.Wallet)
	assert.Equal(t, req.Label, decoded.Label)
	assert.Nil(t, decoded.BundleID)
	assert.Nil(t, decoded.MinAmountUSD)
	assert.Nil(t, decoded.Filters)
	assert.Nil(t, decoded.Chains)
	assert.Nil(t, decoded.NewTrades)
	assert.Nil(t, decoded.TelegramBotID)
	assert.Nil(t, decoded.DiscordChannelID)
}
