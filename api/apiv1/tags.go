package apiv1

type TagType string

const (
	TagTypeHighVolumeDexTrader TagType = "high_volume_dex_trader"
	TagTypeEarlyDefi           TagType = "early_defi"
	TagTypeMultichain          TagType = "multichain"
	TagTypeNewWallet           TagType = "new_wallet"
	TagTypeHighLeverageTrader  TagType = "high_leverage_trader"
	TagTypeNftTrader           TagType = "nft_trader"
	TagTypeNftfi               TagType = "nftfi"
	TagTypeNftHighPnl          TagType = "nft_high_pnl"
	TagTypePopularWallet       TagType = "popular_wallet"
	TagTypeAirdropHunter       TagType = "airdrop_hunter"
	TagTypeGemFinder           TagType = "gem_finder"
	TagTypeHighWinRate         TagType = "high_win_rate"
	TagTypeNewWhale            TagType = "new_whale"
	TagTypeFlipper             TagType = "flipper"
	TagTypeHoneypot            TagType = "honeypot"
	TagTypeMev                 TagType = "mev"
)

type GetWalletTagsRequest struct {
	Wallet string `json:"wallet"`
}

type GetWalletTagsResponse struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Key         TagType `json:"key"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
}
