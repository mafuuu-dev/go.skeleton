package types

import "time"

type RateLimitType struct {
	Max        int
	Expiration time.Duration
}
