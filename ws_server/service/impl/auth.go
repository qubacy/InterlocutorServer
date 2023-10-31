package control

import (
	"ilserver/repository"
	"ilserver/storage"
	"ilserver/transport/control/middleware"
	"ilserver/transport/controlDto"
)

type AuthService struct {
	storage storage.Storage
}

func NewAuthService() *AuthService {
	return &AuthService{
		storage: repository.Instance(),
	}
}

// -----------------------------------------------------------------------

func (self *AuthService) SignIn(req controlDto.SignInReq) (error, controlDto.SignInRes) {
	err, has := self.storage.HasAdminWithLoginAndPass(req.Login, req.Pass)
	if err != nil {
		return err, controlDto.SignInRes{}
	}

	// ***

	var res controlDto.SignInRes

	if !has {
		res.AccessToken = ""
		res.ErrorText = "user not found"
		return nil, res
	}

	// ***

	err, tokenStr := middleware.BuildNewToken(req.Login)
	if err != nil {
		return err, controlDto.SignInRes{}
	}

	res.AccessToken = tokenStr
	res.ErrorText = ""

	return nil, res
}
