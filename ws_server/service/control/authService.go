package control

import (
	"ilserver/repository"
	"ilserver/transport/controlDto"
	"time"

	"github.com/cristalhq/jwt/v5"
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

func buildNewToken(login string) (error, string) {
	key := []byte(`secret`)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return err, ""
	}

	// create claims (you can create your own, see: ExampleBuilder_withUserClaims)
	claims := &jwt.RegisteredClaims{
		Audience: []string{login},
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(5 * time.Minute)),
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

// -----------------------------------------------------------------------

func (as *AuthService) SignIn(req controlDto.SignInReq) (error, controlDto.SignInRes) {
	var res controlDto.SignInRes
	err, has := as.dbFacade.HasAdminWithLoginAndPass(req.Login, req.Pass)
	if err != nil {
		return err, controlDto.SignInRes{}
	}

	if !has {
		res.AccessToken = ""
		res.ErrorText = "user not found"
		return nil, res
	}

	// ***

	err, tokenStr := buildNewToken(req.Login)
	if err != nil {
		res.AccessToken = ""
		res.ErrorText = "build token failed"
		return nil, res
	}

	res.AccessToken = tokenStr
	res.ErrorText = ""

	return nil, res
}
