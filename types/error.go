package types

import "fmt"

type Error struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (e Error) Error() string {
	return fmt.Sprintf("CIELO ERROR: %s", e.Message)
}
