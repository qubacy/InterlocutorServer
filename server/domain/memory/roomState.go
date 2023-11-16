package memory

import "time"

// super state
// -----------------------------------------------------------------------

type StateName int

const (
	SEARCHING StateName = iota
	CHATTING
	CHOOSING
)

func (self StateName) ToString() string {
	switch self {
	case SEARCHING:
		return "SEARCHING"
	case CHATTING:
		return "CHATTING"
	case CHOOSING:
		return "CHOOSING"
	}
	return ""
}

type RoomState struct {
	Name       StateName
	LaunchTime time.Time
}

func MakeRoomState(name StateName, launchTime time.Time) RoomState {
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
