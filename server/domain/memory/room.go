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

// super state
// -----------------------------------------------------------------------

const (
	SEARCHING int = iota
	CHATTING
	CHOOSING
)

type RoomState struct {
	Name       int
	LaunchTime time.Time
}

func MakeRoomState(name int, launchTime time.Time) RoomState {
	return RoomState{
		Name:       name,
		LaunchTime: launchTime,
	}
}

// concrete
// -----------------------------------------------------------------------

type SearchingRoomState struct {
	RoomState
	LastUpdateTime time.Time
}

func MakeSearchingRoomState(time time.Time) SearchingRoomState {
	return SearchingRoomState{
		RoomState:      MakeRoomState(SEARCHING, time),
		LastUpdateTime: time,
	}
}

func NewSearchingRoomState(time time.Time) *SearchingRoomState {
	state := MakeSearchingRoomState(time)
	return &state
}

// -----------------------------------------------------------------------

type ChattingRoomState struct {
	RoomState
}

func MakeChattingRoomState(time time.Time) ChattingRoomState {
	return ChattingRoomState{
		RoomState: MakeRoomState(CHATTING, time),
	}
}

func MakeChattingRoomStateNow() ChattingRoomState {
	return ChattingRoomState{
		RoomState: MakeRoomState(CHATTING, time.Now()),
	}
}

// -----------------------------------------------------------------------

type ChoosingRoomState struct {
	RoomState
	ProfileIdAndMatchedIds map[string][]string
}

func MakeChoosingRoomState(time time.Time) ChattingRoomState {
	return ChattingRoomState{
		RoomState: MakeRoomState(CHOOSING, time),
	}
}
