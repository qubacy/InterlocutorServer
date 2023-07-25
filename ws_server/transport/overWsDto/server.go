package overWsDto

// parts
// -----------------------------------------------------------------------

type SvrMessage struct {
	SenderId int    `json:"senderId"`
	Text     string `json:"text"`
}

type MatchedUsers struct {
	Id      int    `json:"id"`
	Contact string `json:"contact"`
}

type Err struct {
	Message string `json:"message"`
}

type ProfilePublicList struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type FoundGameData struct {
	LocalProfileId        int    `json:"localProfileId"`
	StartSessionTime      int64  `json:"startSessionTime"`
	ChattingStageDuration int64  `json:"chattingStageDuration"`
	ChoosingStageDuration int64  `json:"choosingStageDuration"`
	ChattingTopic         string `json:"chattingTopic"`

	ProfilePublicList []ProfilePublicList `json:"profilePublicList"`
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
	MatchedUsers []MatchedUsers `json:"matchedUsers"`
}

// ***

type SvrErrBody struct {
	Err Err `json:"error"`
}
