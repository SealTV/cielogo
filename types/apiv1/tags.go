package apiv1

// https://developer.cielo.finance/reference/getwallettags

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

// type Tags []Tag

// func (tt *Tags) UnmarshalJSON(data []byte) error {
// 	var tags []Tag

// 	if err := json.Unmarshal(data, &tags); err != nil {
// 		return err
// 	}

// 	*tt = tags
// 	return nil
// }

type TagsResponse struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Key         TagType `json:"key"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
}

// func (t *Tag) UnmarshalJSON(data []byte) error {
// 	type Alias Tag
// 	aux := &struct {
// 		Key string `json:"key"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(t),
// 	}

// 	fmt.Sprintf("aaaaa: %s", string(data))
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}

// 	t.Key = TagType(aux.Key)

// 	return nil
// }
