package apiv1

type WalletsListOrdering string

const (
	PopularWalletsListOrdering WalletsListOrdering = "popular"
	NewWalletsListOrdering     WalletsListOrdering = "new"
)

type GetAllWalletsListsRequest struct {
	FollowOnly bool
	Order      *WalletsListOrdering
	NextObject *string
}

type GetAllWalletsListsResponse struct {
	List   []WalletList `json:"list"`
	Paging Pagination   `json:"paging"`
}

type AddWalletsListRequest struct {
	Name           string   `json:"name"`
	IsPublic       bool     `json:"is_public,omitempty"`
	Wallets        []string `json:"wallets,omitempty"`
	FollowedListID int64    `json:"followed_list_id,omitempty"`
	Description    string   `json:"description,omitempty"`
}

type UpdateWalletsListRequest struct {
	ListID int64 `json:"-"`

	Name           string   `json:"name"`
	IsPublic       bool     `json:"is_public,omitempty"`
	Wallets        []string `json:"wallets,omitempty"`
	FollowedListID int64    `json:"followed_list_id,omitempty"`
	Description    string   `json:"description,omitempty"`
}

type ToggleFollowWalletsListResponce struct {
	Followed bool `json:"followed"`
}

type WalletList struct {
	ID            int64  `json:"id,omitempty"`
	Name          string `json:"name"`
	CreatedAt     int64  `json:"created_at"`
	BotID         int64  `json:"bot_id"`
	Description   string `json:"description"`
	IsPublic      bool   `json:"is_public"`
	FollowedCount int64  `json:"followed_count"`
	WalletsCount  int64  `json:"wallets_count"`
	ShareURL      string `json:"share_url"`
	ImageURL      string `json:"image_url"`
	Followed      bool   `json:"followed"`
	IsCreator     bool   `json:"is_creator"`
}
