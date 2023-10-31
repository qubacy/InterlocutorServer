package dto

type SignInInput struct {
	Login string
	Pass  string
}

type SignInOutput struct {
	AccessToken string
	ErrorText   string
}
