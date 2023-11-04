package impl

import (
	"encoding/json"
	"errors"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"time"

	"github.com/cristalhq/jwt/v5"
)

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if len(signingKey) == 0 {
		return nil, errors.New(token.ErrSigningKeyIsEmpty)
	}

	return &Manager{
		signingKey: signingKey,
	}, nil
}

// ---> 500
// -----------------------------------------------------------------------

func (s *Manager) NewToken(payload token.Payload, duration time.Duration) (error, string) {
	key := []byte(s.signingKey)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return utility.CreateCustomError(s.NewToken, err), ""
	}

	claims := &jwt.RegisteredClaims{
		Subject: payload.Subject,
		ExpiresAt: jwt.NewNumericDate(
			time.Now().UTC().Add(duration)), // !
	}

	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return utility.CreateCustomError(s.NewToken, err), ""
	}
	return nil, token.String()
}

// with user-friendly errors?
// -----------------------------------------------------------------------

func (s *Manager) Parse(tokenValue string) (token.Payload, error) {
	tokenObject, _, err := prepareAndCheck(s.signingKey, tokenValue)
	if err != nil {
		return token.Payload{},
			utility.CreateCustomError(s.Parse, err)
	}

	// ***

	registeredClaims, err := extractAndCheck(tokenObject)
	if err != nil {
		return token.Payload{},
			utility.CreateCustomError(s.Parse, err)
	}

	// ***

	return token.Payload{
		Subject: registeredClaims.Subject,
	}, nil
}

func (s *Manager) Validate(tokenValue string) error {
	_, verifier, err := prepareAndCheck(s.signingKey, tokenValue)
	if err != nil {
		return utility.CreateCustomError(prepareAndCheck, err)
	}

	// parse only claims!
	var registeredClaims jwt.RegisteredClaims
	err = jwt.ParseClaims([]byte(tokenValue), verifier, &registeredClaims)
	if err != nil {
		return utility.CreateCustomError(prepareAndCheck, err)
	}

	if err = checkRegisteredClaims(registeredClaims); err != nil {
		return utility.CreateCustomError(prepareAndCheck, err)
	}
	return nil
}

// private
// -----------------------------------------------------------------------

func prepareAndCheck(signingKey, tokenValue string) (*jwt.Token, *jwt.HSAlg, error) {
	verifierHs, err := jwt.NewVerifierHS(jwt.HS256, []byte(signingKey))
	if err != nil {
		return nil, nil, utility.CreateCustomError(prepareAndCheck, err)
	}

	tokenObject, err := jwt.Parse([]byte(tokenValue), verifierHs)
	if err != nil {
		return nil, nil, utility.CreateCustomError(prepareAndCheck, err)
	}

	// ***

	err = verifierHs.Verify(tokenObject)
	if err != nil {
		return nil, nil, utility.CreateCustomError(prepareAndCheck, err)
	}
	return tokenObject, verifierHs, nil
}

func extractAndCheck(tokenObject *jwt.Token) (jwt.RegisteredClaims, error) {
	var registeredClaims jwt.RegisteredClaims
	err := json.Unmarshal(tokenObject.Claims(), &registeredClaims)
	if err != nil {
		return jwt.RegisteredClaims{},
			utility.CreateCustomError(extractAndCheck, err)
	}

	// ***

	if err = checkRegisteredClaims(registeredClaims); err != nil {
		return jwt.RegisteredClaims{},
			utility.CreateCustomError(extractAndCheck, err)
	}

	return registeredClaims, nil
}

// ***

func checkRegisteredClaims(registeredClaims jwt.RegisteredClaims) error {
	if len(registeredClaims.Subject) == 0 {
		return errors.New(token.ErrSubjectIsEmpty)
	}
	if !registeredClaims.IsValidAt(time.Now()) {
		return errors.New(token.ErrTokenExpired)
	}
	return nil
}
