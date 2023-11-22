package memory

import (
	"time"
)

// -----------------------------------------------------------------------

type Room struct {
	Id           string
	State        interface{}
	CreationTime time.Time
	Profiles     ProfileList // index is id.
	Language     int
}

type RoomList []Room

func MakeEmptyRoomWithSearchingState(lang int) Room {
	currentTime := time.Now()
	return Room{
		Id:       "",
		Profiles: ProfileList{},

		State:        MakeSearchingRoomState(currentTime), // value!
		CreationTime: currentTime,
		Language:     lang,
	}
}

func (self *Room) RemoveProfile(index int) {
	self.Profiles = append(
		self.Profiles[:index],
		self.Profiles[index+1:]...,
	)
}

func (self *Room) StateName() (StateName, bool) {
	switch self.State.(type) {
	case SearchingRoomState:
		return SEARCHING, true
	case ChattingRoomState:
		return CHATTING, true
	case ChoosingRoomState:
		return CHOOSING, true
	}
	return 0, false // <--- name cannot be used!
}
