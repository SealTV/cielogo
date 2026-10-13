package api

import (
	"fmt"
	"time"
)

type Error struct {
	StatusCode         int           `json:"-"`
	RetryAfter         string        `json:"-"`
	RetryAfterDuration time.Duration `json:"-"`
	Message            string        `json:"message"`
	Status             string        `json:"status"`
}

func (e Error) Error() string {
	return fmt.Sprintf("CIELO ERROR: %s", e.Message)
}
