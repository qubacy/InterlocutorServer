package impl

import (
	"context"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"

	"ilserver/storage"
)

type Auth struct {
	storage      storage.Storage
	tokenManager token.Manager
}

func NewAuth(storage storage.Storage, tokenManager token.Manager) *Auth {
	return &Auth{
		storage:      storage,
		tokenManager: tokenManager,
	}
}

// -----------------------------------------------------------------------

func (self *Auth) SignIn(ctx context.Context, reqDto dto.SignInInput) (dto.SignInOutput, error) {
	if !isValidSignInInput(reqDto) {
		return dto.SignInOutput{
			AccessToken: "",
			ErrorText:   control.LoginOrPasswordIsTooShort(),
		}, nil
	}

	// ***

	has, err := self.storage.HasAdminWithLoginAndPassword(ctx, reqDto.Login, reqDto.Pass)
	if err != nil {
		return dto.SignInOutput{}, utility.CreateCustomError(self.SignIn, err)
	}
	if !has {
		return dto.SignInOutput{
			AccessToken: "",
			ErrorText:   control.LoginNotFoundOrPasswordIsIncorrect(reqDto.Login),
		}, nil
	}

	// ***

	err, tokenStr := token.NewToken(reqDto.Login)
	if err != nil {
		return dto.SignInOutput{}, utility.CreateCustomError(self.SignIn, err)
	}
	return dto.SignInOutput{
		AccessToken: tokenStr,
		ErrorText:   "",
	}, nil
}

// validator
// -----------------------------------------------------------------------

func isValidSignInInput(value dto.SignInInput) bool {
	if len(value.Login) < 3 {
		return false
	}
	if len(value.Pass) < 3 {
		return false
	}
	return true
}
