package dto

type SignInReq struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type SignInRes struct {
	AccessToken string `json:"accessToken"`
	ErrorText   string `json:"errorText"`
}
