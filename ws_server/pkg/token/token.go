package token

import (
	"encoding/json"
	"fmt"
	"ilserver/pkg/utility"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/spf13/viper"
)

func NewToken(login string) (error, string) {
	key := []byte(viper.GetString("control_server.token.secret"))
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return utility.CreateCustomError(NewToken, err), ""
	}

	// ***

	claims := &jwt.RegisteredClaims{
		Subject: login,
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(
				viper.GetDuration(
					"control_server.token.duration"))),
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
