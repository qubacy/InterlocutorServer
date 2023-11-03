package dto

type PostAdminInput struct {
	AccessToken string
	Login       string
	Pass        string
}

type PostAdminOutput struct {
	Idr       int
	ErrorText string
}

type DeleteAdminByIdrOutput struct {
	ErrorText string
}
