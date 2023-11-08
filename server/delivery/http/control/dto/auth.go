package dto

import (
	"encoding/json"
	service "ilserver/service/control/dto"
)

type PostSignInReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PostSignInRes struct {
	AccessToken string `json:"access-token"`
}

// constructor
// -----------------------------------------------------------------------

func MakePostSignInReqEmpty() PostSignInReq {
	return PostSignInReq{}
}

func MakePostSignInReqFromJson(rawJson []byte) (PostSignInReq, error) {
	result := PostSignInReq{}
	err := json.Unmarshal(rawJson, &result)
	if err != nil {
		return MakePostSignInReqEmpty(), err
	}

	return result, nil
}

func MakePostSignInRes(serviceOut service.SignInOutput) PostSignInRes {
	return PostSignInRes{
		AccessToken: serviceOut.AccessToken,
	}
}

// converter
// -----------------------------------------------------------------------

func (self PostSignInReq) ToServiceInp() service.SignInInput {
	return service.SignInInput{
		Login: self.Login,
		Pass:  self.Password,
	}
}
