package dto

import (
	"encoding/json"
	service "ilserver/service/control/dto"
)

// -----------------------------------------------------------------------

type GetAdminDto struct {
	Idr   float64 `json:"idr"` // ---> uint64
	Login string  `json:"login"`
}

func MakeGetAdminDto(adminWithoutPass service.AdminWithoutPass) GetAdminDto {
	return GetAdminDto{
		Idr:   float64(adminWithoutPass.Idr),
		Login: adminWithoutPass.Login,
	}
}

type PostAdminDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// -----------------------------------------------------------------------

type GetAdminReq struct{}

type GetAdminRes struct {
	Admins []GetAdminDto `json:"admins"`
}

func MakeGetAdminRes(serviceOut service.GetAdminsOutput) GetAdminRes {
	result := GetAdminRes{}
	for i := range serviceOut.Admins {
		result.Admins = append(result.Admins,
			MakeGetAdminDto(serviceOut.Admins[i]))
	}
	return result
}

// -----------------------------------------------------------------------

type PostAdminReq struct {
	Admin PostAdminDto `json:"admin"`
}

func MakePostAdminReqFromJson(rawJson []byte) (PostAdminReq, error) {
	result := PostAdminReq{}
	err := json.Unmarshal(rawJson, &result)
	if err != nil {
		return MakePostAdminReqEmpty(), err
	}

	return result, nil
}

func MakePostAdminReqEmpty() PostAdminReq {
	return PostAdminReq{}
}

func (self *PostAdminReq) ToServiceInp() service.PostAdminInput {
	return service.PostAdminInput{
		Login: self.Admin.Login,
		Pass:  self.Admin.Password,
	}
}

type PostAdminRes struct {
	Idr float64 `json:"idr"`
}

func MakePostAdminRes(serviceOut service.PostAdminOutput) PostAdminRes {
	return PostAdminRes{
		Idr: float64(serviceOut.Idr),
	}
}
