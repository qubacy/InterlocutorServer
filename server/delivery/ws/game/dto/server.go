package overWsDto

// parts
// -----------------------------------------------------------------------

type SvrMessage struct {
	SenderId int    `json:"senderId"`
	Text     string `json:"text"`
}

type MatchedUser struct {
	Id      int    `json:"id"`
	Contact string `json:"contact"`
}

type Err struct {
	Id int `json:"id"`
}

type ProfilePublic struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type FoundGameData struct {
	LocalProfileId        int    `json:"localProfileId"`
	StartSessionTime      int64  `json:"startSessionTime"`
	ChattingStageDuration int64  `json:"chattingStageDuration"`
	ChoosingStageDuration int64  `json:"choosingStageDuration"`
	ChattingTopic         string `json:"chattingTopic"`

	ProfilePublicList []ProfilePublic `json:"profilePublicList"`
}

// svr - server
// -----------------------------------------------------------------------

type SvrSearchingStartBody struct{}

type SvrSearchingGameFoundBody struct {
	FoundGameData FoundGameData `json:"foundGameData"`
}

type SvrChattingNewMessageBody struct {
	Message SvrMessage `json:"message"`
}

type SvrChattingStageIsOverBody struct{}
type SvrChoosingUsersChosenBody struct{}

type SvrChoosingStageIsOverBody struct {
	MatchedUsers []MatchedUser `json:"matchedUsers"`
}

// ***

// TODO: переделать в код ошибки (?)
type SvrErrBody struct {
	Err Err `json:"error"`
}
