package token

import (
	"time"
)

type Manager interface {
	NewToken(Payload, time.Duration) (string, error)
	Parse(string) (Payload, error)
	Validate(string) error
}
