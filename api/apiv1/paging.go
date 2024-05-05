package apiv1

type Pagination struct {
	HasNextPage     bool   `json:"has_next_page"`
	TotalRowsInPage int    `json:"total_rows_in_page"`
	NextObject      string `json:"next_object"`
}
