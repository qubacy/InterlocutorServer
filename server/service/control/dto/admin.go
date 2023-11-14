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
func MakeGetAdminsOutputSuccess(admins domain.AdminList) GetAdminsOutput {
	result := GetAdminsOutput{}
	for i := range admins {
		result.Admins = append(result.Admins,
			MakeAdminOutput(admins[i]))
	}
	return result
}

// ---> 400
func MakeGetAdminsOutputEmpty() GetAdminsOutput {
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

func (self *PostAdminInput) ToDomain() domain.Admin {
	return domain.Admin{
		Login: self.Login,
		Pass:  self.Pass,
	}
}

// ---> 200
func MakePostAdminOutputSuccess(idr int64) PostAdminOutput {
	return PostAdminOutput{
		Idr: idr,
	}
}

// ---> 400
func MakePostAdminOutputEmpty() PostAdminOutput {
	return PostAdminOutput{}
}

// -----------------------------------------------------------------------

type DeleteAdminByIdrOutput struct {
}

// ---> 200
func MakeDeleteAdminByIdrOutputSuccess() DeleteAdminByIdrOutput {
	return DeleteAdminByIdrOutput{}
}

// ---> 400
func MakeDeleteAdminByIdrOutputEmpty() DeleteAdminByIdrOutput {
	return DeleteAdminByIdrOutput{}
}
