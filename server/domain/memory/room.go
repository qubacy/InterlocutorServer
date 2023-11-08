package memory

import "time"

// -----------------------------------------------------------------------

type Room struct {
	State        interface{}
	CreationTime time.Time
	Profiles     ProfileList // index is id.
	Language     int
}

func MakeEmptyRoomWithSearchingState(lang int) Room {
	currentTime := time.Now()
	return Room{
		Profiles: ProfileList{},

		State:        NewSearchingStateRoom(currentTime),
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

type SearchingStateRoom struct {
	RoomState
	LastUpdateTime time.Time
}

func NewSearchingStateRoom(time time.Time) *SearchingStateRoom {
	return &SearchingStateRoom{
		RoomState:      MakeRoomState(SEARCHING, time),
		LastUpdateTime: time,
	}
}

type ChattingStateRoom struct {
	RoomState
}

type ChoosingStateRoom struct {
	RoomState
	ProfileIdAndMatchedIds map[string][]string // ?
}
