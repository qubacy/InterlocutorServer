package dto

import (
	domain "ilserver/domain/memory"
	"time"
)

// parts
// -----------------------------------------------------------------------

type SvrMessage struct {
	SenderId int    `json:"senderId"`
	Text     string `json:"text"`
}

func MakeSvrMessage(senderId int, text string) SvrMessage {
	return SvrMessage{
		SenderId: senderId,
		Text:     text,
	}
}

type MatchedUser struct {
	Id      int    `json:"id"`
	Contact string `json:"contact"`
}

func MakeMatchedUser(id int, contact string) MatchedUser {
	return MatchedUser{
		Id:      id,
		Contact: contact,
	}
}

type MatchedUserList []MatchedUser

func (self MatchedUserList) AddMatchedUser(matchedUser MatchedUser) {
	self = append(self, matchedUser)
}

type ProfilePublic struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func MakeProfilePublic(id int, username string) ProfilePublic {
	return ProfilePublic{
		Id:       id, // <--- profile local id (index in array)
		Username: username,
	}
}

func MakeProfilePublicList() []ProfilePublic {
	return []ProfilePublic{}
}

type FoundGameData struct {
	LocalProfileId        int    `json:"localProfileId"`
	StartSessionTime      int64  `json:"startSessionTime"`
	ChattingStageDuration int64  `json:"chattingStageDuration"`
	ChoosingStageDuration int64  `json:"choosingStageDuration"`
	ChattingTopic         string `json:"chattingTopic"`

	ProfilePublicList []ProfilePublic `json:"profilePublicList"`
}

func MakeFoundGameData(localProfileId int,
	chattingStageDuration time.Duration,
	choosingStageDuration time.Duration,
	topicName string,
) FoundGameData {
	return FoundGameData{
		LocalProfileId: localProfileId,

		StartSessionTime:      time.Now().Unix(),                    // <--- seconds
		ChattingStageDuration: chattingStageDuration.Milliseconds(), // <--- ms, strange...
		ChoosingStageDuration: choosingStageDuration.Milliseconds(),

		ChattingTopic:     topicName,
		ProfilePublicList: MakeProfilePublicList(),
	}
}

func (self *FoundGameData) AddProfiles(profiles domain.ProfileList) {
	for i := range profiles {
		self.AddProfilePublic(
			MakeProfilePublic(
				i, profiles[i].Username,
			),
		)
	}
}

func (self *FoundGameData) AddProfilePublic(value ProfilePublic) {
	self.ProfilePublicList = append(self.ProfilePublicList, value)
}

// svr - server
// -----------------------------------------------------------------------

type SvrSearchingStartBody struct{}

func MakeSvrSearchingStartBodyEmpty() SvrSearchingStartBody {
	return SvrSearchingStartBody{}
}

// -----------------------------------------------------------------------

type SvrSearchingGameFoundBody struct {
	FoundGameData FoundGameData `json:"foundGameData"`
}

func MakeSvrSearchingGameFoundBody(foundGameData FoundGameData) SvrSearchingGameFoundBody {
	return SvrSearchingGameFoundBody{
		FoundGameData: foundGameData,
	}
}

// -----------------------------------------------------------------------

type SvrChattingNewMessageBody struct {
	Message SvrMessage `json:"message"`
}

func MakeSvrChattingNewMessageBodyFromParts(senderId int, text string) SvrChattingNewMessageBody {
	return SvrChattingNewMessageBody{
		Message: MakeSvrMessage(senderId, text),
	}
}

func MakeSvrChattingNewMessageBodyEmpty() SvrChattingNewMessageBody {
	return SvrChattingNewMessageBody{}
}

// -----------------------------------------------------------------------

type SvrChattingStageIsOverBody struct{}

func MakeSvrChattingStageIsOverBodyEmpty() SvrChattingStageIsOverBody {
	return SvrChattingStageIsOverBody{}
}

// -----------------------------------------------------------------------

type SvrChoosingUsersChosenBody struct{}

type SvrChoosingStageIsOverBody struct {
	MatchedUsers []MatchedUser `json:"matchedUsers"`
}

func MakeSvrChoosingStageIsOverBody(matchedUsers []MatchedUser) SvrChoosingStageIsOverBody {
	return SvrChoosingStageIsOverBody{
		MatchedUsers: matchedUsers,
	}
}
