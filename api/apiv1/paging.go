package apiv1

type Pagination struct {
	HasNextPage     bool   `json:"has_next_page"`
	TotalRowsInPage int    `json:"total_rows_in_page"`
	NextObject      string `json:"next_object"`
}

// FeedPagination contains the opaque base64 cursor used by startFrom.
type FeedPagination struct {
	HasNextPage     bool   `json:"has_next_page"`
	TotalRowsInPage int    `json:"total_rows_in_page"`
	NextCursor      string `json:"next_cursor,omitempty"`
}
