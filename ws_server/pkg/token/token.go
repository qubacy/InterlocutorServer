package token

import (
	"time"
)

type Manager interface {
	New(Payload, time.Duration) (string, error)
	Parse(string) (Payload, error)
	Validate(string) error
}
