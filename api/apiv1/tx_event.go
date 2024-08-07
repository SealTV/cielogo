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

type ContractCreationEvent struct {
	// 	description:
	// This object provides a structure for representing a contract creation event on the blockchain.

	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// amount_usd*	number
	// example: 100.0
	// The equivalent amount in USD of the wrapped tokens.

	// contract_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the smart contract involved in the interaction.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// from_label*	string
	// example: Alice
	// A human-readable label or name associated with the originating wallet.
}

// This object provides a structure for representing a contract interaction event on the blockchain.
type ContractInteractionEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// contract_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the smart contract involved in the interaction.

	// contract_label*	string
	// example: Uniswap V3
	// A human-readable label or name associated with the smart contract.
}

// This object provides a structure for representing a flashloan event on the blockchain.
type FlashloanEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the token involved in the transaction.

	// amount*	number
	// example: 100
	// The amount of tokens involved in the transaction.

	// amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the tokens involved in the transaction.

	// dex*	string
	// example: Uniswap
	// The decentralized exchange (DEX) where the flashloan transaction took place.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// health_factor*	number
	// example: 1
	// The health factor of the wallet after the flashloan transaction.

	// name*	string
	// example: Wrapped Ether
	// The name of the token involved in the transaction.

	// platform*	string
	// example: Aave
	// The platform where the flashloan transaction took place.

	// price_usd*	number
	// example: 100
	// The price of the token in USD.

	// symbol*	string
	// example: ETH
	// The symbol of the token involved in the transaction.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.
}

// This object provides a structure for representing a NFT liquidation event on the blockchain.
type NftLiquidationEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// contract_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the NFT contract involved in the interaction.

	// currency_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the currency involved in the transaction.

	// currency_symbol*	string
	// example: ETH
	// The symbol of the currency involved in the transaction.

	// dex*	string
	// example: Uniswap
	// The decentralized exchange where the wrap transaction occurred.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// nft_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the NFT contract involved in the interaction.

	// nft_name*	string
	// example: NFT
	// The name of the NFT in the transaction.

	// nft_symbol*	string
	// example: NFT
	// The symbol of the NFT in the transaction.

	// platform*	string
	// example: Aave
	// The platform where the flashloan transaction took place.

	// price*	number
	// example: 100
	// The price of the NFT in the transaction.

	// price_usd*	number
	// example: 100
	// The price of the NFT in USD.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// token_id*	string
	// example: 1
	// The unique identifier of the NFT in the transaction.
}

// This object provides a structure for representing a NFT sweep event on the blockchain.
type NftSweepEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address involved in the NFT trading transaction.

	// wallet_label*	string
	// example: 0xa4c8...f9457d
	// A readable label for the wallet address.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique hash identifier of the NFT trading transaction.

	// tx_type*	string
	// example: nft_trade
	// Specifies the type of transaction, in this case, NFT trading.

	// chain*	string
	// example: ethereum
	// The blockchain network where the trading transaction occurred.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// from*	string
	// example: 0xebc18d25d8122da21f73a6bcb78971671f21f6ff
	// The originating wallet address of the transaction.

	// to*	string
	// example: 0xcfdbdb8723619bca23765e09d59c8f745366f8ff
	// The destination wallet address of the transaction.

	// thumbnail*	string
	// example: https://res.cloudinary.com/alchemyapi/image/upload/thumbnailv2/matic-mainnet/284ed679247466683591389baddcfc9b
	// A thumbnail image URL of the NFT involved in the transaction.

	// image*	string
	// example: https://nft-cdn.alchemy.com/matic-mainnet/284ed679247466683591389baddcfc9b
	// A full image URL of the NFT.

	// action*	string
	// example: buy
	// Describes the action taken in the NFT trade, such as 'buy' or 'sell'.

	// contract*	string
	// example: 0x00000000000000adc04c56bf30ac9d3c0aaf14dc
	// The blockchain contract address associated with the NFT.

	// marketplace*	string
	// example: OpenSea
	// The marketplace where the NFT trade occurred, such as OpenSea.

	// nft_address*	string
	// example: 0xbe9371326f91345777b04394448c23e2bfeaa826
	// The blockchain address of the NFT involved in the trade.

	// nft_name*	string
	// example: Gemesis
	// The name of the NFT traded.

	// nft_symbol*	string
	// example: OSP
	// The symbol associated with the NFT.

	// nft_token_id*	string
	// example: 32507
	// The unique token ID of the NFT involved in the trade.

	// price*	number($double)
	// example: 0.0151
	// The price at which the NFT was traded.

	// price_usd*	number($double)
	// example: 32.911356
	// The equivalent USD value of the NFT trade.

	// profit*	number($double)
	// example: 0
	// The profit earned from the trade. This may be zero in some transactions.

	// currency_symbol*	string
	// example: WETH
	// The symbol of the currency used in the trade, such as WETH or ETH.

	// buyer*	string
	// The wallet address of the buyer in the trade.

	// seller*	string
	// The wallet address of the seller in the trade.

	// token*	string
	// The token type used in the transaction.

	// first_interaction*	boolean
	// Indicates whether this was the first interaction between the buyer and seller.

	// bid_accepted*	boolean
	// Specifies if the transaction involved a bid being accepted.
}

// This object provides a structure for representing an option event on the blockchain.
type OptionEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// action*	string
	// example: buy
	// The action taken in the option event.

	// amount*	number
	// example: 100
	// The amount of tokens involved in the transaction.

	// asset*	string
	// example: ETH
	// The asset involved in the option event.

	// dex*	string
	// example: Uniswap
	// The decentralized exchange (DEX) where the option event took place.

	// direction*	string
	// example: call
	// The direction of the option event.

	// expiry*	string
	// example: 2022-12-31
	// The expiry date of the option.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// option_price_usd*	number
	// example: 100
	// The price of the option in USD.

	// position_status*	string
	// example: open
	// The status of the option position.

	// spot_price_usd*	number
	// example: 100
	// The spot price of the asset in USD.

	// status*	string
	// example: active
	// The status of the option event.

	// strike_price_usd*	number
	// example: 100
	// The strike price of the option in USD.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// type*	string
	// example: call
	// The type of option event.
}

// This object provides a structure for representing a Perpetual event on the blockchain.
type Perp struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// action*	string
	// example: buy
	// The action taken in the Perpetual event.

	// amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the tokens involved in the transaction.

	// average_price*	number
	// example: 100
	// The average price of the tokens involved in the transaction.

	// base_token_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the base token involved in the transaction.

	// base_token_amount*	number
	// example: 100
	// The amount of base tokens involved in the transaction.

	// base_token_symbol*	string
	// example: ETH
	// The symbol of the base token involved in the transaction.

	// dex*	string
	// example: Uniswap
	// The decentralized exchange where the Perpetual transaction occurred.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// liquidation*	boolean
	// example: false
	// Indicates whether the transaction was a liquidation.

	// liquidation_price*	number
	// example: 100
	// The price at which the liquidation occurred.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// trade_direction*	string
	// example: long
	// The direction of the trade in the Perpetual transaction.

	// perp_details*	string
	// example: details
	// Additional details about the Perpetual transaction.

	// token0_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the first token in the LP pair.

	// token0_amount*	number
	// example: 100
	// The amount of the first token in the LP pair.

	// token0_amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the first token in the LP pair.

	// token0_name*	string
	// example: Wrapped Ether
	// The name of the first token in the LP pair.

	// token0_price_usd*	number
	// example: 100
	// The price of the first token in the LP pair in USD.

	// token0_symbol*	string
	// example: WETH
	// The symbol of the first token in the LP pair.

	// token1_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the second token in the LP pair.

	// token1_amount*	number
	// example: 100
	// The amount of the second token in the LP pair.

	// token1_amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the second token in the LP pair.

	// token1_name*	string
	// example: Wrapped Ether
	// The name of the second token in the LP pair.

	// token1_price_usd*	number
	// example: 100
	// The price of the second token in the LP pair in USD.

	// token1_symbol*	string
	// example: WETH
	// The symbol of the second token in the LP pair.

	// realized_pnl*	number
	// example: 100
	// The realized profit and loss of the Perpetual transaction.

	// is_nft_perp*	boolean
	// example: false
	// Indicates whether the Perpetual transaction involves an NFT.

	// position_size	number
	// example: 100
	// The size of the position in the Perpetual transaction.

	// position_size_usd	number
	// example: 100
	// The equivalent amount in USD of the position size.

	// leverage	number
	// example: 100
	// The leverage used in the Perpetual transaction.

	// unrealized_pnl	number
	// example: 100
	// The unrealized profit and loss of the Perpetual transaction.
}

// This object provides a structure for representing a reward event on the blockchain.
type RewardEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the token involved in the transaction.

	// amount*	number
	// example: 100
	// The amount of tokens involved in the transaction.

	// amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the tokens involved in the transaction.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// name*	string
	// example: Wrapped Ether
	// The name of the token involved in the transaction.

	// price_usd*	number
	// example: 100
	// The price of the token in USD.

	// symbol*	string
	// example: WETH
	// The symbol of the token involved in the transaction.
}

// This object provides a structure for representing a staking event on the blockchain.
type StakingEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// from_label*	string
	// example: Alice
	// A human-readable label or name associated with the originating wallet.

	// to_label*	string
	// example: Bob
	// A human-readable label or name associated with the destination wallet.

	// amount*	number
	// example: 100
	// The amount of tokens staked in the transaction.

	// amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the staked tokens.

	// token_price_usd*	number
	// example: 100
	// The price of the token in USD.

	// contract_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the smart contract involved in the interaction.

	// symbol*	string
	// example: WETH
	// The symbol of the token staked in the transaction.

	// name*	string
	// example: Wrapped Ether
	// The name of the token staked in the transaction.

	// action*	string
	// example: stake
	// The action taken in the staking transaction.
}

// This object provides a structure for representing a Sudo Pool event on the blockchain.
type SudoPoolEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// dex	string
	// example: Uniswap
	// The decentralized exchange where the wrap transaction occurred.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// nft_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the NFT contract involved in the interaction.

	// nft_amount*	integer
	// example: 1
	// The amount of NFTs involved in the transaction.

	// nft_price*	number
	// example: 100
	// The price of the NFT in the transaction.

	// nft_symbol*	string
	// example: NFT
	// The symbol of the NFT in the transaction.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// token0_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the first token in the LP pair.

	// token0_amount*	number
	// example: 100
	// The amount of the first token in the LP pair.

	// token0_amount_usd*	number
	// example: 100
	// The equivalent amount in USD of the first token in the LP pair.

	// token0_name*	string
	// example: Wrapped Ether
	// The name of the first token in the LP pair.

	// token0_price_usd*	number
	// example: 100
	// The price of the first token in the LP pair in USD.

	// token0_symbol*	string
	// example: WETH
	// The symbol of the first token in the LP pair.
}

// This object provides a structure for representing a wrap event on the blockchain.
type WrapEvent struct {
	// wallet*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The wallet address participating in the LP transaction.

	// wallet_label*	string
	// example: vitalik.eth
	// A human-readable label or name associated with the wallet, such as a ENS name.

	// tx_hash*	string
	// example: 0xe0f84917036e7e2ebb2f8d73199547bde30d4d2918c67904a4bb200dc5bad215
	// The unique transaction hash identifying this specific LP transaction.

	// tx_type*	string
	// example: lp
	// Indicates the type of transaction, in this case, liquidity pool (LP) related.

	// chain*	string
	// example: optimism
	// The blockchain network (e.g., Ethereum, Optimism) where this transaction takes place.

	// index*	integer
	// example: 10
	// A numerical index or identifier for the transaction.

	// timestamp*	integer
	// example: 1702899395
	// The timestamp marking when the transaction was executed.

	// block*	integer
	// example: 113650309
	// The block number on the blockchain where this transaction is recorded.

	// dex*	string
	// example: Uniswap
	// The decentralized exchange where the wrap transaction occurred.

	// from*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The originating wallet address for the transaction.

	// to*	string
	// example: 0xb4fb31e7b1471a8e52dd1e962a281a732ead59c1
	// The destination wallet address for the transaction.

	// action*	string
	// amount*	number
	// example: 100
	// The amount of tokens wrapped in the transaction.

	// amount_usd*	number
	// example: 100.0
	// The equivalent amount in USD of the wrapped tokens.

	// contract_address*	string
	// example: 0x7f5c764cbc14f9669b88837ca1490cca17c31607
	// The address of the smart contract involved in the interaction.

	// name*	string
	// example: Wrapped Ether
	// The name of the token wrapped in the transaction.

	// symbol*	string
	// example: WETH
	// The symbol of the token wrapped in the transaction.

	// token_price_usd*	number
	// example: 100.0
	// The price of the token in USD at the time of the transaction.

	// token_type*	string
	// example: ERC20
	// The type of token wrapped in the transaction.
}
