package memory

import "time"

const (
	SEARCHING int = iota
	CHATTING
	CHOOSING
)

// -----------------------------------------------------------------------

type RoomState struct {
	Name       int
	LaunchTime time.Time
}

type SearchingStateRoom struct {
	RoomState
	LastUpdateTime time.Time
}

type ChattingStateRoom struct {
	RoomState
}

type ChoosingStateRoom struct {
	RoomState
	ProfileIdAndMatchedIds map[string][]string
}

// -----------------------------------------------------------------------

type Room struct {
	State        interface{}
	CreationTime time.Time
	Profiles     []Profile // index is id
	Language     int
}
