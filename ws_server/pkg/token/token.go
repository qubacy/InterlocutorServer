package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cristalhq/jwt/v5"
)

type Payload struct {
	Subject string
}

type Manager interface {
	New(Payload, time.Duration) (string, error)
	Parse(string) (Payload, error)
	Verify(string) error
}

// -----------------------------------------------------------------------

func VerifyToken(token string) {

}

func ParseAndVerifyToken(token string) (error, bool, string) {
	key := []byte(`secret`)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return err, false, ""
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return err, false, ""
	}

	// or just verify it's signature
	err = verifier.Verify(newToken)
	if err != nil {
		return err, false, ""
	}

	// get Registered claims
	var newClaims jwt.RegisteredClaims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	if errClaims != nil {
		return errClaims, false, ""
	}

	// or parse only claims
	errParseClaims := jwt.ParseClaims(tokenBytes, verifier, &newClaims)
	if errParseClaims != nil {
		return errParseClaims, false, ""
	}

	// verify claims as you wish
	if len(newClaims.Audience) != 1 {
		return fmt.Errorf("audience has not login"), false, ""
	}
	var isValid bool = newClaims.IsValidAt(time.Now())
	if !isValid {
		return nil, false, ""
	}

	aud := newClaims.Audience
	return nil, true, aud[0]
}
