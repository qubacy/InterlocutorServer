package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/spf13/viper"
)

func BuildNewToken(login string) (error, string) {
	key := []byte(`secret`)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return err, ""
	}

	// create claims
	claims := &jwt.RegisteredClaims{
		Audience: []string{login},
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(
				viper.GetDuration(
					"control_server.token_duration"))),
	}

	// create a Builder
	builder := jwt.NewBuilder(signer)

	// and build a Token
	token, err := builder.Build(claims)
	if err != nil {
		return err, ""
	}

	// here is token as a string
	return nil, token.String()
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
