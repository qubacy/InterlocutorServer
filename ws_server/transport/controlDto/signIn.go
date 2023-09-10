package controlDto

type SignInReq struct {
	Login string
	Pass  string
}

type SignInRes struct {
	AccessToken string
}
