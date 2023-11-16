package dto

import (
	"encoding/json"
	domain "ilserver/domain/memory"
	"ilserver/pkg/utility"
)

// parts
// -----------------------------------------------------------------------

type Profile struct {
	Username string `json:"username"`
	Contact  string `json:"contact"`
	Language int    `json:"language"`
}

type CliMessage struct {
	Text string `json:"text"`
}

// cli - client
// -----------------------------------------------------------------------

type CliSearchingStartBody struct {
	Profile Profile `json:"profile"`
}

func MakeCliSearchingStartBodyFromJson(data []byte) (
	CliSearchingStartBody, error,
) {
	dto := CliSearchingStartBody{}
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return MakeCliSearchingStartBodyEmpty(),
			utility.CreateCustomError(
				MakeCliSearchingStartBodyFromJson, err)
	}
	return dto, nil
}

func MakeCliSearchingStartBodyEmpty() CliSearchingStartBody {
	return CliSearchingStartBody{}
}

// -----------------------------------------------------------------------

// not used.
type CliSearchingStopBody struct{}

func MakeCliSearchingStopBodyFromJson(data []byte) (
	CliSearchingStopBody, error,
) {
	return CliSearchingStopBody{}, nil
}

func MakeCliSearchingStopBodyEmpty() CliSearchingStopBody {
	return CliSearchingStopBody{}
}

// -----------------------------------------------------------------------

type CliChattingNewMessageBody struct {
	Message CliMessage `json:"message"`
}

func MakeCliChattingNewMessageBodyFromJson(data []byte) (
	CliChattingNewMessageBody, error,
) {
	dto := CliChattingNewMessageBody{}
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return MakeCliChattingNewMessageBodyEmpty(),
			utility.CreateCustomError(
				MakeCliChattingNewMessageBodyFromJson, err)
	}
	return dto, nil
}

func MakeCliChattingNewMessageBodyEmpty() CliChattingNewMessageBody {
	return CliChattingNewMessageBody{}
}

// -----------------------------------------------------------------------

type CliChoosingUsersChosenBody struct {
	UserIdList []int `json:"userIdList"`
}

func MakeCliChoosingUsersChosenBodyFromJson(data []byte) (
	CliChoosingUsersChosenBody, error,
) {
	dto := CliChoosingUsersChosenBody{}
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return MakeCliChoosingUsersChosenBodyEmpty(),
			utility.CreateCustomError(
				MakeCliChoosingUsersChosenBodyFromJson, err)
	}
	return dto, nil
}

func MakeCliChoosingUsersChosenBodyEmpty() CliChoosingUsersChosenBody {
	return CliChoosingUsersChosenBody{}
}

// validator (can be moved to service)
// -----------------------------------------------------------------------

func (dto *CliSearchingStartBody) IsValid() bool {
	return dto.Profile.Contact != "" && dto.Profile.Username != ""
}

func (dto *CliChattingNewMessageBody) IsValid() bool {
	return dto.Message.Text != ""
}

// adapter
// -----------------------------------------------------------------------

func MakeProfileFromReqDto(id string, dto CliSearchingStartBody) domain.Profile {
	return domain.Profile{
		Id: id,

		Username: dto.Profile.Username,
		Contact:  dto.Profile.Contact,
	}
}
