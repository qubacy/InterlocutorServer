package dto

import "ilserver/domain"

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
	Idr       int64
	ErrorText string
}

func PostAdminInputToDomain(inp PostAdminInput) domain.Admin {
	return domain.Admin{
		Login: inp.Login,
		Pass:  inp.Pass,
	}
}

// ---> 200
func MakePostAdminSuccess(idr int64) PostAdminOutput {
	return PostAdminOutput{
		Idr:       idr,
		ErrorText: "",
	}
}

// ---> 400
func MakePostAdminError(errorText string) PostAdminOutput {
	return PostAdminOutput{
		Idr:       0,
		ErrorText: errorText,
	}
}

// -----------------------------------------------------------------------

type DeleteAdminByIdrOutput struct {
	ErrorText string
}

// ---> 200
func MakeDeleteAdminByIdrSuccess() DeleteAdminByIdrOutput {
	return DeleteAdminByIdrOutput{}
}

// ---> 400
func MakeDeleteAdminByIdrError(errorText string) DeleteAdminByIdrOutput {
	return DeleteAdminByIdrOutput{
		ErrorText: errorText,
	}
}
