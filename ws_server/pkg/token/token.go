package token

import (
	"time"
)

type Payload struct {
	Subject string
}

type Manager interface {
	New(Payload, time.Duration) (string, error)
	Parse(string) (Payload, error)
	Validate(string) error
}
