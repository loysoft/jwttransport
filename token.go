package jwttransport

import (
	"time"
)

// Token contains an end-user's tokens.
// This is the data you must store to persist authentication.
type Token struct {
	AccessToken string
	Expiry      time.Time // If zero the token has no (known) expiry time.
}

// Expired reports whether the token has expired or is invalid.
func (t *Token) Expired() bool {
	if t.Expiry.IsZero() || t.AccessToken == "" {
		return false
	}
	return t.Expiry.Before(time.Now())
}
