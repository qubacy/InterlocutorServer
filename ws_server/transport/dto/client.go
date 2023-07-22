package dto

// cli - client

type CliSearchingStartBodyClient struct {
	Profile struct {
		Username string `json:"username"`
		Contact  string `json:"contact"`
	} `json:"profile"`
}

// not used
type CliSearchingStopBody struct{}

type CliChattingNewMessageBody struct {
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

type CliChoosingUsersChosenBody struct {
	UserIdList []int `json:"userIdList"`
}
