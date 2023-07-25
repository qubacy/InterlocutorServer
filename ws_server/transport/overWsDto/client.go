package overWsDto

// TODO: общие DTOs независимо от протокола?

import (
	"ilserver/domain"
)

// parts
// -----------------------------------------------------------------------
type Profile struct {
	Username string `json:"username"`
	Contact  string `json:"contact"`
}

type CliMessage struct {
	Text string `json:"text"`
}

// cli - client
// -----------------------------------------------------------------------

type CliSearchingStartBodyClient struct {
	Profile Profile `json:"profile"`
}

// not used
type CliSearchingStopBody struct{}

type CliChattingNewMessageBody struct {
	Message CliMessage `json:"message"`
}

type CliChoosingUsersChosenBody struct {
	UserIdList []int `json:"userIdList"`
}

// validator
// -----------------------------------------------------------------------

func (dto *CliSearchingStartBodyClient) IsValid() bool {
	return dto.Profile.Contact != "" && dto.Profile.Username != ""
}

func (dto *CliChattingNewMessageBody) IsValid() bool {
	return dto.Message.Text != ""
}

func (dto *CliChoosingUsersChosenBody) IsValid() bool {
	return len(dto.UserIdList) > 0
}

// adapter
// -----------------------------------------------------------------------

func MakeProfileFromReqDto(id string, dto CliSearchingStartBodyClient) domain.Profile {
	return domain.Profile{
		Id: id,

		Username: dto.Profile.Username,
		Contact:  dto.Profile.Contact,
	}
}
