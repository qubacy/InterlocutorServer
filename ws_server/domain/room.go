package domain

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
}

// -----------------------------------------------------------------------

// TODO: можно ли заменить на структуру?
type Room struct {
	State        interface{}
	CreationTime time.Time
	Profiles     []Profile // index is id
}

func MakeRoom() Room {
	return Room{
		State:        SearchingStateRoom{},
		CreationTime: time.Now(),
		Profiles:     []Profile{},
	}
}
