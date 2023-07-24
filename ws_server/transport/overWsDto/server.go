package overWsDto

// svr - server
// TODO: одинаковые ли дто при разных протоколах?

type SvrSearchingStartBody struct{}

type SvrSearchingGameFoundBody struct {
	FoundGameData struct {
		LocalProfileId        int    `json:"localProfileId"`
		StartSessionTime      int64  `json:"startSessionTime"`
		ChattingStageDuration int64  `json:"chattingStageDuration"`
		ChoosingStageDuration int64  `json:"choosingStageDuration"`
		ChattingTopic         string `json:"chattingTopic"`
		ProfilePublicList     []struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
		} `json:"profilePublicList"`
	}
}

type SvrChattingNewMessageBody struct {
	Message struct {
		SenderId int    `json:"senderId"`
		Text     string `json:"text"`
	} `json:"message"`
}

type SvrChattingStageIsOverBody struct{}
type SvrChoosingUsersChosenBody struct{}

type SvrChoosingStageIsOverBody struct {
	MatchedUsers []struct {
		Id      int    `json:"id"`
		Contact string `json:"contact"`
	} `json:"matchedUsers"`
}

// ***

type SvrErrBody struct {
	Err struct {
		Message string `json:"message"`
	} `json:"error"`
}
