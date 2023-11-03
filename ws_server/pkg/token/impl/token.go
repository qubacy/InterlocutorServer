package impl

import (
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/spf13/viper"
)

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if len(signingKey) == 0 {

	}

	return &Manager{
		signingKey: signingKey,
	}, nil
}

// -----------------------------------------------------------------------

func (s *Manager) NewToken(payload token.Payload, duration time.Duration) (error, string) {
	key := []byte(viper.GetString("control_server.token.secret"))
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return utility.CreateCustomError(s.NewToken, err), ""
	}

	claims := &jwt.RegisteredClaims{
		Subject: payload.Subject,
		ExpiresAt: jwt.NewNumericDate(
			time.Now().UTC().Add(
				viper.GetDuration(
					"control_server.token.duration"))),
	}

	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return utility.CreateCustomError(s.NewToken, err), ""
	}
	return nil, token.String()
}
