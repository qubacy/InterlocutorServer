package dto

type SignInInput struct {
	Login string
	Pass  string
}

type SignInOutput struct {
	AccessToken string
}

// ---> 200
func MakeSignInSuccess(accessToken string) SignInOutput {
	return SignInOutput{
		AccessToken: accessToken,
	}
}

// ---> 400
func MakeSignInEmpty() SignInOutput {
	return SignInOutput{}
}
