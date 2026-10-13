package apiv1

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sealtv/cielogo/api/chains"
)

type TxType string

const (
	TxTypeBridge      TxType = "bridge"
	TxTypeLending     TxType = "lending"
	TxTypeLP          TxType = "lp"
	TxTypeNftLending  TxType = "nft_lending"
	TxTypeNftMint     TxType = "nft_mint"
	TxTypeNftTrade    TxType = "nft_trade"
	TxTypeNftTransfer TxType = "nft_transfer"
	TxTypeSwap        TxType = "swap"
	TxTypeTransfer    TxType = "transfer"

	TxTypeContractCreation    TxType = "contract_creation"
	TxTypeContractInteraction TxType = "contract_interaction"
	TxTypeFlashloan           TxType = "flashloan"
	TxTypeNftLiquidation      TxType = "nft_liquidation"
	TxTypeNftSweep            TxType = "nft_sweep"
	TxTypeOption              TxType = "option"
	TxTypePerp                TxType = "perp"
	TxTypeReward              TxType = "reward"
	TxTypeStaking             TxType = "staking"
	TxTypeSudoPool            TxType = "sudo_pool"
	TxTypeWrap                TxType = "wrap"
)

// Additional event and virtual filter types supported by the API.
const (
	TxTypeNftMitnt          TxType = TxTypeNftMint
	TxTypePerpOrder         TxType = "perp_order"
	TxTypePerpTwap          TxType = "perp_twap"
	TxTypePerpImminentLiq   TxType = "perp_imminent_liq"
	TxTypeHyperliquidLedger TxType = "hyperliquid_ledger"
	TxTypeSwapBuy           TxType = "swap_buy"
	TxTypeSwapSell          TxType = "swap_sell"
	TxTypeTransferIn        TxType = "transfer_in"
	TxTypeTransferOut       TxType = "transfer_out"
)

// TokenMarketCap contains enrichment requested with IncludeMarketCap.
type TokenMarketCap struct {
	TokenAddress string  `json:"token_address"`
	MarketCap    float64 `json:"market_cap"`
	Liquidity    float64 `json:"liquidity"`
}

type TxEvent struct {
	Wallet         string           `json:"wallet"`
	WalletLabel    string           `json:"wallet_label"`
	TxHash         string           `json:"tx_hash"`
	TxType         TxType           `json:"tx_type"`
	Chain          chains.ChainType `json:"chain"`
	Index          int              `json:"index"`
	Timestamp      int64            `json:"timestamp"`
	Block          uint             `json:"block"`
	TokenMarketCap *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Verification   json.RawMessage  `json:"verification,omitempty"`
	Migration      json.RawMessage  `json:"migration,omitempty"`
	Data           TransactionEvent `json:"-"`
}

func (t *TxEvent) UnmarshalJSON(data []byte) error {
	type header TxEvent
	var decoded header
	if err := json.Unmarshal(data, &decoded); err != nil {
		return fmt.Errorf("failed to unmarshal tx event: %w", err)
	}
	if decoded.TxType == "" {
		return errors.New("tx event requires tx_type")
	}
	decoded.Data = createTransactionEventByType(decoded.TxType)
	if decoded.Data == nil {
		decoded.Data = &UnknownTransactionEvent{Type: decoded.TxType, Raw: append(json.RawMessage(nil), data...)}
	} else if err := json.Unmarshal(data, decoded.Data); err != nil {
		return fmt.Errorf("failed to unmarshal tx event %q data: %w", decoded.TxType, err)
	}
	*t = TxEvent(decoded)
	return nil
}

//nolint:funlen // One explicit case per supported protocol event.
func createTransactionEventByType(txType TxType) TransactionEvent {
	switch txType {
	case TxTypeBridge:
		return &BridgeEvent{}
	case TxTypeLending:
		return &LendingEvent{}
	case TxTypeLP:
		return &LpEvent{}
	case TxTypeNftLending:
		return &NftLendingEvent{}
	case TxTypeNftMint:
		return &NftMintEvent{}
	case TxTypeNftTrade:
		return &NftTradeEvent{}
	case TxTypeNftTransfer:
		return &NftTransferEvent{}
	case TxTypeSwap:
		return &SwapEvent{}
	case TxTypeTransfer:
		return &TransferEvent{}
	case TxTypeContractCreation:
		return &ContractCreationEvent{}
	case TxTypeContractInteraction:
		return &ContractInteractionEvent{}
	case TxTypeFlashloan:
		return &FlashloanEvent{}
	case TxTypeNftLiquidation:
		return &NftLiquidationEvent{}
	case TxTypeNftSweep:
		return &NftSweepEvent{}
	case TxTypeOption:
		return &OptionEvent{}
	case TxTypePerp:
		return &PerpEvent{}
	case TxTypeReward:
		return &RewardEvent{}
	case TxTypeStaking:
		return &StakingEvent{}
	case TxTypeSudoPool:
		return &SudoPoolEvent{}
	case TxTypeWrap:
		return &WrapEvent{}
	case TxTypePerpOrder:
		return &PerpOrderEvent{}
	case TxTypePerpTwap:
		return &PerpTwapEvent{}
	case TxTypePerpImminentLiq:
		return &PerpImminentLiqEvent{}
	case TxTypeHyperliquidLedger:
		return &HyperliquidLedgerEvent{}
	case TxTypeSwapBuy, TxTypeSwapSell, TxTypeTransferIn, TxTypeTransferOut:
		return nil // Virtual query filters are not response event types.
	default:
		return nil
	}
}

type TransactionEvent interface{ GetType() TxType }

// UnknownTransactionEvent preserves a future event without failing the whole page.
type UnknownTransactionEvent struct {
	Type TxType
	Raw  json.RawMessage
}

func (e *UnknownTransactionEvent) GetType() TxType { return e.Type }

// LP action values.
type LpType string

const (
	LpTypeAdd       LpType = "add"
	LpTypeAddLpType        = LpTypeAdd
	LpTypeRemove           = "remove"
)

// BridgeEvent is the API payload for TxTypeBridge.
type BridgeEvent struct {
	Wallet         string           `json:"wallet"`
	WalletLabel    string           `json:"wallet_label,omitempty"`
	TxHash         string           `json:"tx_hash"`
	TxType         TxType           `json:"tx_type"`
	Chain          chains.ChainType `json:"chain"`
	Index          int              `json:"index"`
	Timestamp      int64            `json:"timestamp"`
	Block          uint             `json:"block"`
	Verification   json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration      json.RawMessage  `json:"migration,omitempty"`
	From           string           `json:"from"`
	To             string           `json:"to"`
	FromLabel      string           `json:"from_label"`
	ToLabel        string           `json:"to_label"`
	TokenAddress   string           `json:"token_address"`
	TokenName      string           `json:"token_name"`
	TokenSymbol    string           `json:"token_symbol"`
	TokenIconLink  string           `json:"token_icon_link"`
	Amount         float64          `json:"amount"`
	AmountUSD      float64          `json:"amount_usd"`
	FromChain      string           `json:"from_chain"`
	ToChain        string           `json:"to_chain"`
	Platform       string           `json:"platform"`
	Price          float64          `json:"price"`
	Type           string           `json:"type"`
}

func (e *BridgeEvent) GetType() TxType { return TxTypeBridge }

// LendingEvent is the API payload for TxTypeLending.
type LendingEvent struct {
	Wallet         string           `json:"wallet"`
	WalletLabel    string           `json:"wallet_label,omitempty"`
	TxHash         string           `json:"tx_hash"`
	TxType         TxType           `json:"tx_type"`
	Chain          chains.ChainType `json:"chain"`
	Index          int              `json:"index"`
	Timestamp      int64            `json:"timestamp"`
	Block          uint             `json:"block"`
	Verification   json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration      json.RawMessage  `json:"migration,omitempty"`
	From           string           `json:"from"`
	FromLabel      string           `json:"from_label"`
	Action         string           `json:"action"`
	Address        string           `json:"address"`
	Amount         float64          `json:"amount"`
	AmountUSD      float64          `json:"amount_usd"`
	Dex            string           `json:"dex"`
	HealthFactor   float64          `json:"health_factor"`
	Name           string           `json:"name"`
	Platform       string           `json:"platform"`
	PriceUSD       float64          `json:"price_usd"`
	Symbol         string           `json:"symbol"`
	TokenIconLink  string           `json:"token_icon_link"`
}

func (e *LendingEvent) GetType() TxType { return TxTypeLending }

// LpEvent is the API payload for TxTypeLP.
type LpEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	Dex             string           `json:"dex"`
	From            string           `json:"from"`
	Type            LpType           `json:"type"`
	Token0Address   string           `json:"token0_address"`
	Token0Amount    float64          `json:"token0_amount"`
	Token0AmountUSD float64          `json:"token0_amount_usd"`
	Token0Name      string           `json:"token0_name"`
	Token0PriceUSD  float64          `json:"token0_price_usd"`
	Token0Symbol    string           `json:"token0_symbol"`
	Token0IconLink  string           `json:"token0_icon_link"`
	Token1Address   string           `json:"token1_address"`
	Token1Amount    float64          `json:"token1_amount"`
	Token1AmountUSD float64          `json:"token1_amount_usd"`
	Token1Name      string           `json:"token1_name"`
	Token1PriceUSD  float64          `json:"token1_price_usd"`
	Token1Symbol    string           `json:"token1_symbol"`
	Token1IconLink  string           `json:"token1_icon_link"`
	LowerBound      float64          `json:"lower_bound"`
	UpperBound      float64          `json:"upper_bound"`
	FromLabel       string           `json:"from_label,omitempty"`
	ToLabel         string           `json:"to_label,omitempty"`
}

func (e *LpEvent) GetType() TxType { return TxTypeLP }

// NftLendingEvent is the API payload for TxTypeNftLending.
type NftLendingEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
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

func (e *NftLendingEvent) GetType() TxType { return TxTypeNftLending }

// NftMintEvent is the API payload for TxTypeNftMint.
type NftMintEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
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

func (e *NftMintEvent) GetType() TxType { return TxTypeNftMint }

// NftTradeEvent is the API payload for TxTypeNftTrade.
type NftTradeEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
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
	IsSell          bool             `json:"is_sell"`
	BidAccepted     bool             `json:"bid_accepted"`
	FromLabel       string           `json:"from_label,omitempty"`
	ToLabel         string           `json:"to_label,omitempty"`
}

func (e *NftTradeEvent) GetType() TxType { return TxTypeNftTrade }

// NftTransferEvent is the API payload for TxTypeNftTransfer.
type NftTransferEvent struct {
	Wallet           string           `json:"wallet"`
	WalletLabel      string           `json:"wallet_label,omitempty"`
	TxHash           string           `json:"tx_hash"`
	TxType           TxType           `json:"tx_type"`
	Chain            chains.ChainType `json:"chain"`
	Index            int              `json:"index"`
	Timestamp        int64            `json:"timestamp"`
	Block            uint             `json:"block"`
	Verification     json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap   *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration        json.RawMessage  `json:"migration,omitempty"`
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

func (e *NftTransferEvent) GetType() TxType { return TxTypeNftTransfer }

// SwapEvent is the API payload for TxTypeSwap.
type SwapEvent struct {
	Wallet           string           `json:"wallet"`
	WalletLabel      string           `json:"wallet_label,omitempty"`
	TxHash           string           `json:"tx_hash"`
	TxType           TxType           `json:"tx_type"`
	Chain            chains.ChainType `json:"chain"`
	Index            int              `json:"index"`
	Timestamp        int64            `json:"timestamp"`
	Block            uint             `json:"block"`
	Verification     json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap   *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration        json.RawMessage  `json:"migration,omitempty"`
	Dex              string           `json:"dex"`
	From             string           `json:"from"`
	To               string           `json:"to"`
	Token0Address    string           `json:"token0_address"`
	Token0Amount     float64          `json:"token0_amount"`
	Token0AmountUsd  float64          `json:"token0_amount_usd"`
	Token0PriceUsd   float64          `json:"token0_price_usd"`
	Token0Name       string           `json:"token0_name"`
	Token0Symbol     string           `json:"token0_symbol"`
	Token0IconLink   string           `json:"token0_icon_link"`
	Token1Address    string           `json:"token1_address"`
	Token1Amount     float64          `json:"token1_amount"`
	Token1AmountUsd  float64          `json:"token1_amount_usd"`
	Token1PriceUsd   float64          `json:"token1_price_usd"`
	Token1Name       string           `json:"token1_name"`
	Token1Symbol     string           `json:"token1_symbol"`
	Token1IconLink   string           `json:"token1_icon_link"`
	FirstInteraction bool             `json:"first_interaction"`
	IsSell           bool             `json:"is_sell"`
	Platform         string           `json:"platform"`
	PlatformURI      string           `json:"platform_uri"`
	PlatformUI       string           `json:"platform_ui"`
	FromLabel        string           `json:"from_label,omitempty"`
	ToLabel          string           `json:"to_label,omitempty"`
	IsOutcome        bool             `json:"is_outcome,omitempty"`
	OutcomeQuestion  string           `json:"outcome_question,omitempty"`
	OutcomeSide      string           `json:"outcome_side,omitempty"`
	BuilderTag       string           `json:"builder_tag,omitempty"`
	IsSettlement     bool             `json:"is_settlement,omitempty"`
	ResolutionSide   string           `json:"resolution_side,omitempty"`
	OutcomeClosedPnl float64          `json:"outcome_closed_pnl,omitempty"`
	QuoteCoin        string           `json:"quote_coin,omitempty"`
	AssetID          string           `json:"asset_id,omitempty"`
}

func (e *SwapEvent) GetType() TxType { return TxTypeSwap }

// TransferEvent is the API payload for TxTypeTransfer.
type TransferEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	FromLabel       string           `json:"from_label"`
	ToLabel         string           `json:"to_label"`
	Amount          float64          `json:"amount"`
	AmountUsd       float64          `json:"amount_usd"`
	ContractAddress string           `json:"contract_address"`
	Name            string           `json:"name"`
	Symbol          string           `json:"symbol"`
	TokenPriceUsd   float64          `json:"token_price_usd"`
	Type            string           `json:"type"`
	TokenIconLink   string           `json:"token_icon_link"`
	IsMempool       *bool            `json:"mempool,omitempty"`
}

func (e *TransferEvent) GetType() TxType { return TxTypeTransfer }

// ContractCreationEvent is the API payload for TxTypeContractCreation.
type ContractCreationEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	AmountUsd       int              `json:"amount_usd"`
	ContractAddress string           `json:"contract_address"`
	From            string           `json:"from"`
	FromLabel       string           `json:"from_label"`
}

func (e *ContractCreationEvent) GetType() TxType { return TxTypeContractCreation }

// ContractInteractionEvent is the API payload for TxTypeContractInteraction.
type ContractInteractionEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	ContractAddress string           `json:"contract_address"`
	ContractLabel   string           `json:"contract_label"`
	FromLabel       string           `json:"from_label,omitempty"`
	ToLabel         string           `json:"to_label,omitempty"`
}

func (e *ContractInteractionEvent) GetType() TxType { return TxTypeContractInteraction }

// FlashloanEvent is the API payload for TxTypeFlashloan.
type FlashloanEvent struct {
	Wallet         string           `json:"wallet"`
	WalletLabel    string           `json:"wallet_label,omitempty"`
	TxHash         string           `json:"tx_hash"`
	TxType         TxType           `json:"tx_type"`
	Chain          chains.ChainType `json:"chain"`
	Index          int              `json:"index"`
	Timestamp      int64            `json:"timestamp"`
	Block          uint             `json:"block"`
	Verification   json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration      json.RawMessage  `json:"migration,omitempty"`
	Address        string           `json:"address"`
	Amount         float64          `json:"amount"`
	AmountUsd      float64          `json:"amount_usd"`
	Dex            string           `json:"dex"`
	From           string           `json:"from"`
	HealthFactor   float64          `json:"health_factor"`
	Name           string           `json:"name"`
	Platform       string           `json:"platform"`
	PriceUsd       float64          `json:"price_usd"`
	Symbol         string           `json:"symbol"`
	TokenIconLink  string           `json:"token_icon_link"`
	FromLabel      string           `json:"from_label,omitempty"`
}

func (e *FlashloanEvent) GetType() TxType { return TxTypeFlashloan }

// NftLiquidationEvent is the API payload for TxTypeNftLiquidation.
type NftLiquidationEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	ContractAddress string           `json:"contract_address"`
	CurrencyAddress string           `json:"currency_address"`
	CurrencySymbol  string           `json:"currency_symbol"`
	Dex             string           `json:"dex"`
	From            string           `json:"from"`
	NftAddress      string           `json:"nft_address"`
	NftName         string           `json:"nft_name"`
	NftSymbol       string           `json:"nft_symbol"`
	Platform        string           `json:"platform"`
	Price           float64          `json:"price"`
	PriceUsd        float64          `json:"price_usd"`
	To              string           `json:"to"`
	TokenId         string           `json:"token_id"`
}

func (e *NftLiquidationEvent) GetType() TxType { return TxTypeNftLiquidation }

// NftSweepEvent is the API payload for TxTypeNftSweep.
type NftSweepEvent struct {
	Wallet           string           `json:"wallet"`
	WalletLabel      string           `json:"wallet_label,omitempty"`
	TxHash           string           `json:"tx_hash"`
	TxType           TxType           `json:"tx_type"`
	Chain            chains.ChainType `json:"chain"`
	Index            int              `json:"index"`
	Timestamp        int64            `json:"timestamp"`
	Block            uint             `json:"block"`
	Verification     json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap   *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration        json.RawMessage  `json:"migration,omitempty"`
	From             string           `json:"from"`
	To               string           `json:"to"`
	Thumbnail        string           `json:"thumbnail"`
	Image            string           `json:"image"`
	Action           string           `json:"action"`
	Contract         string           `json:"contract"`
	Marketplace      string           `json:"marketplace"`
	NftAddress       string           `json:"nft_address"`
	NftName          string           `json:"nft_name"`
	NftSymbol        string           `json:"nft_symbol"`
	NftTokenId       string           `json:"nft_token_id"`
	Price            float64          `json:"price"`
	PriceUsd         float64          `json:"price_usd"`
	Profit           float64          `json:"profit"`
	CurrencySymbol   string           `json:"currency_symbol"`
	Buyer            string           `json:"buyer"`
	Seller           string           `json:"seller"`
	Token            string           `json:"token"`
	FirstInteraction bool             `json:"first_interaction"`
	IsSell           bool             `json:"is_sell"`
	BidAccepted      bool             `json:"bid_accepted"`
	FromLabel        string           `json:"from_label,omitempty"`
	ToLabel          string           `json:"to_label,omitempty"`
}

func (e *NftSweepEvent) GetType() TxType { return TxTypeNftSweep }

// OptionEvent is the API payload for TxTypeOption.
type OptionEvent struct {
	Wallet         string           `json:"wallet"`
	WalletLabel    string           `json:"wallet_label,omitempty"`
	TxHash         string           `json:"tx_hash"`
	TxType         TxType           `json:"tx_type"`
	Chain          chains.ChainType `json:"chain"`
	Index          int              `json:"index"`
	Timestamp      int64            `json:"timestamp"`
	Block          uint             `json:"block"`
	Verification   json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration      json.RawMessage  `json:"migration,omitempty"`
	Action         string           `json:"action"`
	Amount         float64          `json:"amount"`
	Asset          string           `json:"asset"`
	Dex            string           `json:"dex"`
	Direction      string           `json:"direction"`
	Expiry         string           `json:"expiry"`
	From           string           `json:"from"`
	OptionPriceUsd float64          `json:"option_price_usd"`
	PositionStatus string           `json:"position_status"`
	SpotPriceUsd   float64          `json:"spot_price_usd"`
	Status         string           `json:"status"`
	StrikePriceUsd float64          `json:"strike_price_usd"`
	To             string           `json:"to"`
	Type           string           `json:"type"`
}

func (e *OptionEvent) GetType() TxType { return TxTypeOption }

// PerpEvent is the API payload for TxTypePerp.
type PerpEvent struct {
	Wallet           string           `json:"wallet"`
	WalletLabel      string           `json:"wallet_label,omitempty"`
	TxHash           string           `json:"tx_hash"`
	TxType           TxType           `json:"tx_type"`
	Chain            chains.ChainType `json:"chain"`
	Index            int              `json:"index"`
	Timestamp        int64            `json:"timestamp"`
	Block            uint             `json:"block"`
	Verification     json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap   *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration        json.RawMessage  `json:"migration,omitempty"`
	Action           string           `json:"action"`
	AmountUsd        float64          `json:"amount_usd"`
	AveragePrice     float64          `json:"average_price"`
	BaseTokenAddress string           `json:"base_token_address"`
	BaseTokenAmount  float64          `json:"base_token_amount"`
	BaseTokenSymbol  string           `json:"base_token_symbol"`
	Dex              string           `json:"dex"`
	From             string           `json:"from"`
	Liquidation      bool             `json:"liquidation"`
	LiquidationPrice float64          `json:"liquidation_price,omitempty"`
	To               string           `json:"to"`
	TradeDirection   string           `json:"trade_direction"`
	PerpDetails      string           `json:"perp_details"`
	Token0Address    string           `json:"token0_address"`
	Token0Amount     float64          `json:"token0_amount"`
	Token0AmountUsd  float64          `json:"token0_amount_usd"`
	Token0Name       string           `json:"token0_name"`
	Token0PriceUsd   float64          `json:"token0_price_usd"`
	Token0Symbol     string           `json:"token0_symbol"`
	Token0IconLink   string           `json:"token0_icon_link"`
	Token1Address    string           `json:"token1_address"`
	Token1Amount     float64          `json:"token1_amount"`
	Token1AmountUsd  float64          `json:"token1_amount_usd"`
	Token1Name       string           `json:"token1_name"`
	Token1PriceUsd   float64          `json:"token1_price_usd"`
	Token1Symbol     string           `json:"token1_symbol"`
	RealizedPnl      float64          `json:"realized_pnl"`
	IsNftPerp        bool             `json:"is_nft_perp"`
	Token1IconLink   string           `json:"token1_icon_link"`
	PositionSize     *float64         `json:"position_size"`
	PositionSizeUsd  *float64         `json:"position_size_usd"`
	Leverage         *float64         `json:"leverage"`
	UnrealizedPnl    *float64         `json:"unrealized_pnl"`
	FromLabel        string           `json:"from_label,omitempty"`
	ToLabel          string           `json:"to_label,omitempty"`
}

func (e *PerpEvent) GetType() TxType { return TxTypePerp }

// RewardEvent is the API payload for TxTypeReward.
type RewardEvent struct {
	Wallet         string           `json:"wallet"`
	WalletLabel    string           `json:"wallet_label,omitempty"`
	TxHash         string           `json:"tx_hash"`
	TxType         TxType           `json:"tx_type"`
	Chain          chains.ChainType `json:"chain"`
	Index          int              `json:"index"`
	Timestamp      int64            `json:"timestamp"`
	Block          uint             `json:"block"`
	Verification   json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration      json.RawMessage  `json:"migration,omitempty"`
	Address        string           `json:"address"`
	Amount         float64          `json:"amount"`
	AmountUsd      float64          `json:"amount_usd"`
	From           string           `json:"from"`
	Name           string           `json:"name"`
	PriceUsd       float64          `json:"price_usd"`
	Symbol         string           `json:"symbol"`
	FromLabel      string           `json:"from_label,omitempty"`
}

func (e *RewardEvent) GetType() TxType { return TxTypeReward }

// StakingEvent is the API payload for TxTypeStaking.
type StakingEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	FromLabel       string           `json:"from_label"`
	ToLabel         string           `json:"to_label"`
	Amount          float64          `json:"amount"`
	AmountUsd       float64          `json:"amount_usd"`
	TokenPriceUsd   float64          `json:"token_price_usd"`
	ContractAddress string           `json:"contract_address"`
	Symbol          string           `json:"symbol"`
	Name            string           `json:"name"`
	Action          string           `json:"action"`
}

func (e *StakingEvent) GetType() TxType { return TxTypeStaking }

// SudoPoolEvent is the API payload for TxTypeSudoPool.
type SudoPoolEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	Dex             string           `json:"dex"`
	From            string           `json:"from"`
	NftAddress      string           `json:"nft_address"`
	NftAmount       int              `json:"nft_amount"`
	NftPrice        float64          `json:"nft_price"`
	NftSymbol       string           `json:"nft_symbol"`
	To              string           `json:"to"`
	Token0Address   string           `json:"token0_address"`
	Token0Amount    float64          `json:"token0_amount"`
	Token0AmountUsd float64          `json:"token0_amount_usd"`
	Token0Name      string           `json:"token0_name"`
	Token0PriceUsd  float64          `json:"token0_price_usd"`
	Token0Symbol    string           `json:"token0_symbol"`
	Token0IconLink  string           `json:"token0_icon_link"`
	FromLabel       string           `json:"from_label,omitempty"`
	ToLabel         string           `json:"to_label,omitempty"`
}

func (e *SudoPoolEvent) GetType() TxType { return TxTypeSudoPool }

// WrapEvent is the API payload for TxTypeWrap.
type WrapEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	Dex             string           `json:"dex"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	Action          string           `json:"action"`
	Amount          float64          `json:"amount"`
	AmountUsd       float64          `json:"amount_usd"`
	ContractAddress string           `json:"contract_address"`
	Name            string           `json:"name"`
	Symbol          string           `json:"symbol"`
	TokenPriceUsd   float64          `json:"token_price_usd"`
	TokenType       string           `json:"token_type"`
	TokenIconLink   string           `json:"token_icon_link"`
	FromLabel       string           `json:"from_label,omitempty"`
	ToLabel         string           `json:"to_label,omitempty"`
}

func (e *WrapEvent) GetType() TxType { return TxTypeWrap }

// PerpOrderEvent is the API payload for TxTypePerpOrder.
type PerpOrderEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	Action          string           `json:"action"`
	Status          string           `json:"status"`
	TradeDirection  string           `json:"trade_direction"`
	Side            string           `json:"side"`
	Dex             string           `json:"dex"`
	BaseTokenSymbol string           `json:"base_token_symbol"`
	BaseTokenAmount float64          `json:"base_token_amount"`
	AmountUsd       float64          `json:"amount_usd"`
	AveragePrice    float64          `json:"average_price"`
	LimitPx         float64          `json:"limit_px"`
	Sz              float64          `json:"sz"`
	OrigSz          float64          `json:"orig_sz"`
	Oid             int64            `json:"oid"`
	OrderType       string           `json:"order_type"`
	Tif             string           `json:"tif"`
	ReduceOnly      bool             `json:"reduce_only"`
	Token0Address   string           `json:"token0_address"`
	Token0Amount    float64          `json:"token0_amount"`
	Token0AmountUsd float64          `json:"token0_amount_usd"`
	Token0Name      string           `json:"token0_name"`
	Token0PriceUsd  float64          `json:"token0_price_usd"`
	Token0Symbol    string           `json:"token0_symbol"`
	Token1Address   string           `json:"token1_address"`
	Token1Amount    float64          `json:"token1_amount"`
	Token1AmountUsd float64          `json:"token1_amount_usd"`
	Token1Name      string           `json:"token1_name"`
	Token1PriceUsd  float64          `json:"token1_price_usd"`
	Token1Symbol    string           `json:"token1_symbol"`
}

func (e *PerpOrderEvent) GetType() TxType { return TxTypePerpOrder }

// PerpTwapEvent is the API payload for TxTypePerpTwap.
type PerpTwapEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	Action          string           `json:"action"`
	Status          string           `json:"status"`
	TradeDirection  string           `json:"trade_direction"`
	Side            string           `json:"side"`
	Dex             string           `json:"dex"`
	BaseTokenSymbol string           `json:"base_token_symbol"`
	BaseTokenAmount float64          `json:"base_token_amount"`
	AmountUsd       float64          `json:"amount_usd"`
	AveragePrice    float64          `json:"average_price"`
	TwapId          int64            `json:"twap_id"`
	TargetSize      float64          `json:"target_size"`
	ExecutedSz      float64          `json:"executed_sz"`
	ExecutedNtl     float64          `json:"executed_ntl"`
	ProgressPct     float64          `json:"progress_pct"`
	Minutes         int              `json:"minutes"`
	ReduceOnly      bool             `json:"reduce_only"`
	Randomize       bool             `json:"randomize"`
	ErrorMsg        string           `json:"error_msg,omitempty"`
	Token0Address   string           `json:"token0_address"`
	Token0Amount    float64          `json:"token0_amount"`
	Token0AmountUsd float64          `json:"token0_amount_usd"`
	Token0Name      string           `json:"token0_name"`
	Token0PriceUsd  float64          `json:"token0_price_usd"`
	Token0Symbol    string           `json:"token0_symbol"`
	Token1Address   string           `json:"token1_address"`
	Token1Amount    float64          `json:"token1_amount"`
	Token1AmountUsd float64          `json:"token1_amount_usd"`
	Token1Name      string           `json:"token1_name"`
	Token1PriceUsd  float64          `json:"token1_price_usd"`
	Token1Symbol    string           `json:"token1_symbol"`
}

func (e *PerpTwapEvent) GetType() TxType { return TxTypePerpTwap }

// PerpImminentLiqEvent is the API payload for TxTypePerpImminentLiq.
type PerpImminentLiqEvent struct {
	Wallet           string           `json:"wallet"`
	WalletLabel      string           `json:"wallet_label,omitempty"`
	TxHash           string           `json:"tx_hash"`
	TxType           TxType           `json:"tx_type"`
	Chain            chains.ChainType `json:"chain"`
	Index            int              `json:"index"`
	Timestamp        int64            `json:"timestamp"`
	Block            uint             `json:"block"`
	Verification     json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap   *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration        json.RawMessage  `json:"migration,omitempty"`
	Action           string           `json:"action"`
	TradeDirection   string           `json:"trade_direction"`
	Dex              string           `json:"dex"`
	BaseTokenSymbol  string           `json:"base_token_symbol"`
	BaseTokenAmount  float64          `json:"base_token_amount"`
	AmountUsd        float64          `json:"amount_usd"`
	AveragePrice     float64          `json:"average_price"`
	PositionSize     float64          `json:"position_size"`
	PositionSizeUsd  float64          `json:"position_size_usd"`
	DistancePct      float64          `json:"distance_pct"`
	LiquidationPrice float64          `json:"liquidation_price"`
	MarkPrice        float64          `json:"mark_price"`
	Leverage         int              `json:"leverage"`
	Token0Address    string           `json:"token0_address"`
	Token0Amount     float64          `json:"token0_amount"`
	Token0AmountUsd  float64          `json:"token0_amount_usd"`
	Token0Name       string           `json:"token0_name"`
	Token0PriceUsd   float64          `json:"token0_price_usd"`
	Token0Symbol     string           `json:"token0_symbol"`
	Token1Address    string           `json:"token1_address"`
	Token1Amount     float64          `json:"token1_amount"`
	Token1AmountUsd  float64          `json:"token1_amount_usd"`
	Token1Name       string           `json:"token1_name"`
	Token1PriceUsd   float64          `json:"token1_price_usd"`
	Token1Symbol     string           `json:"token1_symbol"`
}

func (e *PerpImminentLiqEvent) GetType() TxType { return TxTypePerpImminentLiq }

// HyperliquidLedgerEvent is the API payload for TxTypeHyperliquidLedger.
type HyperliquidLedgerEvent struct {
	Wallet          string           `json:"wallet"`
	WalletLabel     string           `json:"wallet_label,omitempty"`
	TxHash          string           `json:"tx_hash"`
	TxType          TxType           `json:"tx_type"`
	Chain           chains.ChainType `json:"chain"`
	Index           int              `json:"index"`
	Timestamp       int64            `json:"timestamp"`
	Block           uint             `json:"block"`
	Verification    json.RawMessage  `json:"verification,omitempty"`
	TokenMarketCap  *TokenMarketCap  `json:"token_market_cap,omitempty"`
	Migration       json.RawMessage  `json:"migration,omitempty"`
	Action          string           `json:"action"`
	LedgerType      string           `json:"ledger_type"`
	Token           string           `json:"token"`
	Amount          string           `json:"amount"`
	UsdcValue       string           `json:"usdc_value"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	Hash            string           `json:"hash"`
	Dex             string           `json:"dex"`
	BaseTokenSymbol string           `json:"base_token_symbol"`
	BaseTokenAmount float64          `json:"base_token_amount"`
	AmountUsd       float64          `json:"amount_usd"`
	AveragePrice    float64          `json:"average_price"`
	Token0Address   string           `json:"token0_address"`
	Token0Amount    float64          `json:"token0_amount"`
	Token0AmountUsd float64          `json:"token0_amount_usd"`
	Token0Name      string           `json:"token0_name"`
	Token0PriceUsd  float64          `json:"token0_price_usd"`
	Token0Symbol    string           `json:"token0_symbol"`
}

func (e *HyperliquidLedgerEvent) GetType() TxType { return TxTypeHyperliquidLedger }

// MarshalJSON flattens the typed transaction payload into the same shape as the
// API feed. The outer event header takes precedence over duplicated payload fields.
func (t TxEvent) MarshalJSON() ([]byte, error) {
	type header TxEvent
	if t.TxType == "" && t.Data != nil {
		t.TxType = t.Data.GetType()
	}
	if t.TxType == "" {
		return nil, errors.New("tx event requires tx_type")
	}
	var payload map[string]json.RawMessage
	if t.Data != nil {
		encoded, err := json.Marshal(t.Data)
		if err != nil {
			return nil, fmt.Errorf("encode transaction payload: %w", err)
		}
		if err := json.Unmarshal(encoded, &payload); err != nil {
			return nil, fmt.Errorf("decode transaction payload: %w", err)
		}
	}
	if payload == nil {
		payload = make(map[string]json.RawMessage)
	}
	encoded, err := json.Marshal(header(t))
	if err != nil {
		return nil, fmt.Errorf("encode event header: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		return nil, fmt.Errorf("decode event header: %w", err)
	}
	for key, value := range fields {
		payload[key] = value
	}
	return json.Marshal(payload)
}

func (e UnknownTransactionEvent) MarshalJSON() ([]byte, error) {
	if !json.Valid(e.Raw) {
		return nil, errors.New("unknown event contains invalid JSON")
	}
	return e.Raw, nil
}
