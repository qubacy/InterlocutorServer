package control

import (
	"ilserver/repository"
	"ilserver/transport/control/middleware"
	"ilserver/transport/controlDto"
)

type AuthService struct {
	dbFacade *repository.DbFacade
}

func NewAuthService() *AuthService {
	return &AuthService{
		dbFacade: repository.Instance(),
	}
}

// -----------------------------------------------------------------------

func (as *AuthService) SignIn(req controlDto.SignInReq) (error, controlDto.SignInRes) {
	err, has := as.dbFacade.HasAdminWithLoginAndPass(req.Login, req.Pass)
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
