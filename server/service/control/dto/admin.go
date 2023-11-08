package dto

import "ilserver/domain"

// -----------------------------------------------------------------------

type AdminWithoutPass struct {
	Idr   int64
	Login string
}

func MakeAdminOutput(admin domain.Admin) AdminWithoutPass {
	return AdminWithoutPass{
		Idr:   int64(admin.Idr),
		Login: admin.Login,
	}
}

type GetAdminsOutput struct {
	Admins []AdminWithoutPass
}

// ---> 200
func MakeGetAdminsSuccess(admins domain.AdminList) GetAdminsOutput {
	result := GetAdminsOutput{}
	for i := range admins {
		result.Admins = append(result.Admins,
			MakeAdminOutput(admins[i]))
	}
	return result
}

// ---> 400
func MakeGetAdminsEmpty() GetAdminsOutput {
	return GetAdminsOutput{}
}

// -----------------------------------------------------------------------

type PostAdminInput struct {
	Login string
	Pass  string
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
