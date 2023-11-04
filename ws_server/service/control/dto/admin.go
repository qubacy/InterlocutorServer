package dto

// -----------------------------------------------------------------------

type GetAdminsOutput struct {
	AdminLogins []string
	ErrorText   string
}

// ---> 200
func MakeGetAdminsSuccess(adminLogins []string) GetAdminsOutput {
	return GetAdminsOutput{
		AdminLogins: adminLogins,
		ErrorText:   "",
	}
}

// ---> 400
func MakeGetAdminsError(errorText string) GetAdminsOutput {
	return GetAdminsOutput{
		AdminLogins: []string{},
		ErrorText:   errorText,
	}
}

// -----------------------------------------------------------------------

type PostAdminInput struct {
	AccessToken string
	Login       string
	Pass        string
}

type PostAdminOutput struct {
	Idr       int
	ErrorText string
}

// -----------------------------------------------------------------------

type DeleteAdminByIdrOutput struct {
	ErrorText string
}
