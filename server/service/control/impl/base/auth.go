package base

import (
	"context"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"
	"time"

	"ilserver/storage"
)

type AuthService struct {
	accessTokenTTL time.Duration
	storage        storage.Storage
	tokenManager   token.Manager
}

func NewAuth(
	accessTokenTTL time.Duration,
	storage storage.Storage,
	tokenManager token.Manager,
) *AuthService {
	return &AuthService{
		accessTokenTTL: accessTokenTTL,
		storage:        storage,
		tokenManager:   tokenManager,
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
			utility.CreateCustomError(
				self.SignIn, control.ErrLoginNotFoundOrPasswordIsIncorrect(inp.Login))
	}

	// ***

	tokenValue, err := self.tokenManager.NewToken(
		token.MakePayload(inp.Login),
		self.accessTokenTTL,
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
