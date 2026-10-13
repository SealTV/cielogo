package cielogo

import (
	"testing"
	"time"
)

func FuzzRetryAfter(f *testing.F) {
	for _, seed := range []string{"10", "0", "-1", "9223372036854775807", "Tue, 13 Oct 2026 23:00:00 GMT", "bad"} {
		f.Add(seed)
	}
	now := time.Date(2026, time.October, 13, 20, 0, 0, 0, time.UTC)
	f.Fuzz(func(t *testing.T, value string) {
		if delay := parseRetryAfter(value, now); delay < 0 {
			t.Fatalf("negative retry delay: %v", delay)
		}
	})
}
