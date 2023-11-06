package base

import (
	"context"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"

	"ilserver/storage"

	"github.com/spf13/viper"
)

type AuthService struct {
	storage      storage.Storage
	tokenManager token.Manager
}

func NewAuth(storage storage.Storage, tokenManager token.Manager) *AuthService {
	return &AuthService{
		storage:      storage,
		tokenManager: tokenManager,
	}
}

// -----------------------------------------------------------------------

func (self *AuthService) SignIn(ctx context.Context, inp dto.SignInInput) (dto.SignInOutput, error) {
	if !isValidSignInInput(inp) {
		// ---> 400
		return dto.MakeSignInEmpty(),
			utility.CreateCustomError(self.SignIn,
				control.ErrLoginOrPasswordIsTooShort())
	}

	// ***

	has, err := self.storage.HasAdminWithLoginAndPassword(ctx, inp.Login, inp.Pass)
	if err != nil {
		// ---> 400
		return dto.MakeSignInEmpty(),
			utility.CreateCustomError(self.SignIn, err)
	}
	if !has {
		// ---> 400
		return dto.MakeSignInEmpty(),
			utility.CreateCustomError(self.SignIn, err)
	}

	// ***

	tokenValue, err := self.tokenManager.New(
		token.MakePayload(inp.Login),
		viper.GetDuration("control_server.token.duration"),
	)
	if err != nil {
		// ---> 400
		return dto.MakeSignInEmpty(),
			utility.CreateCustomError(self.SignIn, err)
	}

	// ---> 200
	return dto.MakeSignInSuccess(tokenValue), nil
}

// private
// -----------------------------------------------------------------------

func isValidSignInInput(inp dto.SignInInput) bool {
	if len(inp.Login) < 3 {
		return false
	}
	if len(inp.Pass) < 3 {
		return false
	}
	return true
}
