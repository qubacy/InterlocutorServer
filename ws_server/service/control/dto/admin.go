package dto

import "ilserver/domain"

// -----------------------------------------------------------------------

type GetAdminsOutput struct {
	AdminLogins []string
}

// ---> 200
func MakeGetAdminsSuccess(adminLogins []string) GetAdminsOutput {
	return GetAdminsOutput{
		AdminLogins: adminLogins,
	}
}

// ---> 400
func MakeGetAdminsEmpty() GetAdminsOutput {
	return GetAdminsOutput{}
}

// -----------------------------------------------------------------------

type PostAdminInput struct {
	AccessToken string
	Login       string
	Pass        string
}

type PostAdminOutput struct {
	Idr int64
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
		Idr: idr,
	}
}

// ---> 400
func MakePostAdminEmpty() PostAdminOutput {
	return PostAdminOutput{}
}

// -----------------------------------------------------------------------

type DeleteAdminByIdrOutput struct {
}

// ---> 200
func MakeDeleteAdminByIdrSuccess() DeleteAdminByIdrOutput {
	return DeleteAdminByIdrOutput{}
}

// ---> 400
func MakeDeleteAdminByIdrEmpty() DeleteAdminByIdrOutput {
	return DeleteAdminByIdrOutput{}
}
