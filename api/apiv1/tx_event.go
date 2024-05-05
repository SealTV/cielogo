package apiv1

import (
	"encoding/json"
	"fmt"

	"github.com/sealtv/cielogo/api/chains"
)

type TxType string

const (
	TxTypeBridge      TxType = "bridge"
	TxTypeLending     TxType = "lending"
	TxTypeLP          TxType = "lp"
	TxTypeNftLending  TxType = "nft_lending"
	TxTypeNftMitnt    TxType = "nft_mint"
	TxTypeNftTrade    TxType = "nft_trade"
	TxTypeNftTransfer TxType = "nft_transfer"
	TxTypeSwap        TxType = "swap"
	TxTypeTransfer    TxType = "transfer"
)

type TxEvent struct {
	Wallet      string           `json:"wallet"`
	WalletLabel string           `json:"wallet_label"`
	TxHash      string           `json:"tx_hash"`
	TxType      TxType           `json:"tx_type"`
	Chain       chains.ChainType `json:"chain"`
	Index       int              `json:"index"`
	Timestamp   int64            `json:"timestamp"`
	Block       int              `json:"block"`
	Data        TransactionEvent `json:"-"`
}

func (t *TxEvent) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Wallet      string           `json:"wallet"`
		WalletLabel string           `json:"wallet_label"`
		TxHash      string           `json:"tx_hash"`
		TxType      TxType           `json:"tx_type"`
		Chain       chains.ChainType `json:"chain"`
		Index       int              `json:"index"`
		Timestamp   int64            `json:"timestamp"`
		Block       int              `json:"block"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal tx event: %w", err)
	}

	t.Wallet = tmp.Wallet
	t.WalletLabel = tmp.WalletLabel
	t.TxHash = tmp.TxHash
	t.TxType = tmp.TxType
	t.Chain = tmp.Chain
	t.Index = tmp.Index
	t.Timestamp = tmp.Timestamp
	t.Block = tmp.Block

	switch tmp.TxType {
	case TxTypeBridge:
		t.Data = &BridgeEvent{}
	case TxTypeLending:
		t.Data = &LendingEvent{}
	case TxTypeLP:
		t.Data = &LpEvent{}
	case TxTypeNftLending:
		t.Data = &NftLendingEvent{}
	case TxTypeNftMitnt:
		t.Data = &NftMintEvent{}
	case TxTypeNftTrade:
		t.Data = &NftTradeEvent{}
	case TxTypeNftTransfer:
		t.Data = &NftTransferEvent{}
	case TxTypeSwap:
		t.Data = &SwapEvent{}
	case TxTypeTransfer:
		t.Data = &TransferEvent{}
	}

	if err := json.Unmarshal(data, t.Data); err != nil {
		return fmt.Errorf("failed to unmarshal tx event %q data: %w", string(t.TxType), err)
	}

	return nil
}

type TransactionEvent interface {
	GetType() TxType
}

type BridgeEvent struct {
	Wallet       string           `json:"wallet"`
	WalletLabel  string           `json:"wallet_label"`
	TxHash       string           `json:"tx_hash"`
	TxType       TxType           `json:"tx_type"`
	Chain        chains.ChainType `json:"chain"`
	Index        int              `json:"index"`
	Timestamp    int64            `json:"timestamp"`
	Block        int              `json:"block"`
	From         string           `json:"from"`
	To           string           `json:"to"`
	FromLabel    string           `json:"from_label"`
	ToLabel      string           `json:"to_label"`
	TokenAddress string           `json:"token_address"`
	TokenName    string           `json:"token_name"`
	TokenSymbol  string           `json:"token_symbol"`
	Amoun        float64          `json:"amount"`
	AmountUSD    float64          `json:"amount_usd"`
	FromChain    chains.ChainType `json:"from_chain"`
	ToChain      chains.ChainType `json:"to_chain"`
	Platfrom     string           `json:"platform"`
	Price        float64          `json:"price"`
	Type         string           `json:"type"`
}

func (b *BridgeEvent) GetType() TxType {
	return TxTypeBridge
}

type LendingEvent struct {
	Wallet       string           `json:"wallet"`
	WalletLabel  string           `json:"wallet_label"`
	TxHash       string           `json:"tx_hash"`
	TxType       TxType           `json:"tx_type"`
	Chain        chains.ChainType `json:"chain"`
	Index        int              `json:"index"`
	Timestamp    int64            `json:"timestamp"`
	Block        int              `json:"block"`
	From         string           `json:"from"`
	FromLabel    string           `json:"from_label"`
	Action       string           `json:"action"`
	Address      string           `json:"address"`
	Amount       float64          `json:"amount"`
	AmountUSD    float64          `json:"amount_usd"`
	Dex          string           `json:"dex"`
	HealthFactor float64          `json:"health_factor"`
	Name         string           `json:"name"`
	Platform     string           `json:"platform"`
	PriceUSD     float64          `json:"price_usd"`
	Symbol       string           `json:"symbol"`
}

func (l *LendingEvent) GetType() TxType {
	return TxTypeLending
}

type LpType string

const (
	LpTypeAddLpType = "add"
	LpTypeRemove    = "remove"
)

type LpEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           int              `json:"block"`
	Dex             string           `json:"dex"`
	From            string           `json:"from"`
	Type            LpType           `json:"type"`
	Token0Address   string           `json:"token0_address"`
	Token0Amount    float64          `json:"token0_amount"`
	Token0AmountUSD float64          `json:"token0_amount_usd"`
	Token0Name      string           `json:"token0_name"`
	Token0PriceUSD  float64          `json:"token0_price_usd"`
	Token0Symbol    string           `json:"token0_symbol"`
	Token1Address   string           `json:"token1_address"`
	Token1Amount    float64          `json:"token1_amount"`
	Token1AmountUSD float64          `json:"token1_amount_usd"`
	Token1Name      string           `json:"token1_name"`
	Token1PriceUSD  float64          `json:"token1_price_usd"`
	Token1Symbol    string           `json:"token1_symbol"`
	LowerBound      float64          `json:"lower_bound"`
	UpperBound      float64          `json:"upper_bound"`
}

func (l *LpEvent) GetType() TxType {
	return TxTypeLP
}

type NftLendingEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           int              `json:"block"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	FromLabel       string           `json:"from_label"`
	ToLabel         string           `json:"to_label"`
	Thumbnail       string           `json:"thumbnail"`
	Image           string           `json:"image"`
	Action          string           `json:"action"`
	CurrencyAddress string           `json:"currency_address"`
	CurrenctSymbol  string           `json:"currency_symbol"`
	Inserte         float64          `json:"interest"`
	NftAddress      string           `json:"nft_address"`
	NftName         string           `json:"nft_name"`
	NftSymbol       string           `json:"nft_symbol"`
	Platform        string           `json:"platform"`
	NftTokenId      string           `json:"nft_token_id"`
	Price           float64          `json:"price"`
	PriceUSD        float64          `json:"price_usd"`
	Terms           float64          `json:"terms"`
	Refinance       bool             `json:"refinance"`
}

func (n *NftLendingEvent) GetType() TxType {
	return TxTypeNftLending
}

type NftMintEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           int              `json:"block"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	FromLabel       string           `json:"from_label"`
	ToLabel         string           `json:"to_label"`
	Thumbnail       string           `json:"thumbnail"`
	Image           string           `json:"image"`
	Amount          float64          `json:"amount"`
	ContractAddress string           `json:"contract_address"`
	ContractType    string           `json:"contract_type"`
	Fee             float64          `json:"fee"`
	NftName         string           `json:"nft_name"`
	NftSymbol       string           `json:"nft_symbol"`
	NftToekenId     string           `json:"nft_token_id"`
	CurrencySymbol  string           `json:"currency_symbol"`
	Type            string           `json:"type"`
	Value           float64          `json:"value"`
	ValueUsd        float64          `json:"value_usd"`
}

func (n *NftMintEvent) GetType() TxType {
	return TxTypeNftMitnt
}

type NftTradeEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           int              `json:"block"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	Thtumbnail      string           `json:"thumbnail"`
	Image           string           `json:"image"`
	Action          string           `json:"action"`
	Contract        string           `json:"contract"`
	Marketplace     string           `json:"marketplace"`
	NftAddress      string           `json:"nft_address"`
	NftName         string           `json:"nft_name"`
	NftSymbol       string           `json:"nft_symbol"`
	NftTokenId      string           `json:"nft_token_id"`
	Price           float64          `json:"price"`
	PriceUsd        float64          `json:"price_usd"`
	Profit          float64          `json:"profit"`
	CurrencySymbol  string           `json:"currency_symbol"`
	Buyer           string           `json:"buyer"`
	Seller          string           `json:"seller"`
	Token           string           `json:"token"`
	FirsInteraction bool             `json:"first_interaction"`
	BidAccepted     bool             `json:"bid_accepted"`
}

func (n *NftTradeEvent) GetType() TxType {
	return TxTypeNftTrade
}

type NftTransferEvent struct {
	Wallet           string           `json:"wallet"`
	WalletLabel      string           `json:"wallet_label"`
	TxHash           string           `json:"tx_hash"`
	TxType           TxType           `json:"tx_type"`
	Chain            chains.ChainType `json:"chain"`
	Index            int              `json:"index"`
	Timestamp        int64            `json:"timestamp"`
	Block            int              `json:"block"`
	From             string           `json:"from"`
	To               string           `json:"to"`
	FromLabel        string           `json:"from_label"`
	ToLabel          string           `json:"to_label"`
	Thtumbnail       string           `json:"thumbnail"`
	Image            string           `json:"image"`
	ConstractAddress string           `json:"contract_address"`
	ConstractType    string           `json:"contract_type"`
	Fee              float64          `json:"fee"`
	NftName          string           `json:"nft_name"`
	NftSymbol        string           `json:"nft_symbol"`
	NftTokenId       string           `json:"nft_token_id"`
	Type             string           `json:"type"`
	Value            float64          `json:"value"`
}

func (n *NftTransferEvent) GetType() TxType {
	return TxTypeNftTransfer
}

type SwapEvent struct {
	Wallet       string           `json:"wallet"`
	WalletLabel  string           `json:"wallet_label"`
	TxHash       string           `json:"tx_hash"`
	TxType       TxType           `json:"tx_type"`
	Chain        chains.ChainType `json:"chain"`
	Index        int              `json:"index"`
	Timestamp    int64            `json:"timestamp"`
	Block        int              `json:"block"`
	From         string           `json:"from"`
	To           string           `json:"to"`
	FromLabel    string           `json:"from_label"`
	ToLabel      string           `json:"to_label"`
	TokenAddress string           `json:"token_address"`
	TokenName    string           `json:"token_name"`
	TokenSymbol  string           `json:"token_symbol"`
	Amount       float64          `json:"amount"`
	AmountUsd    float64          `json:"amount_usd"`
	FromChain    string           `json:"from_chain"`
	ToChain      string           `json:"to_chain"`
	Platform     string           `json:"platform"`
	Price        float64          `json:"price"`
	Type         string           `json:"type"`
}

func (s *SwapEvent) GetType() TxType {
	return TxTypeSwap
}

type TransferEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           int              `json:"block"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	FromLabel       string           `json:"from_label"`
	ToLabel         string           `json:"to_label"`
	AmountUsd       float64          `json:"amount_usd"`
	ContractAddress string           `json:"contract_address"`
	Name            string           `json:"name"`
	Symbol          string           `json:"symbol"`
	TokenPriceUsd   float64          `json:"token_price_usd"`
	Type            string           `json:"type"`
}

func (t *TransferEvent) GetType() TxType {
	return TxTypeTransfer
}
