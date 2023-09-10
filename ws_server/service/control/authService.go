package control

import (
	"ilserver/repository"
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

func (as *AuthService) SignIn(req controlDto.SignInReq) (error, controlDto.SignInRes) {
	var res controlDto.SignInRes

	return nil, res
}
