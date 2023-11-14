package dto

// direct mapping
// -----------------------------------------------------------------------

type RoomState struct {
	Name       int     `json:"name"`
	LaunchTime float64 `json:"launch-time"`
}

type SearchingRoomState struct {
	RoomState      `json:"room-state"`
	LastUpdateTime float64 `json:"last-update-time"`
}

type ChattingRoomState struct {
	RoomState `json:"room-state"`
}

type ChoosingRoomState struct {
	RoomState
	ProfileIdAndMatchedIds map[string][]string // ?
}

// parts
// -----------------------------------------------------------------------

// can remove the postfix 'dto'
type GetRoomDto struct {
	Id           float64     `json:"id"` // ---> uint64
	State        interface{} `json:"state"`
	CreationTime float64     `json:"creation-time"`
	Profiles     []Profile   `json:"profile-list"`
	Language     float64     `json:"language"`
}

// req/res
// -----------------------------------------------------------------------

type GetRoomReq struct{}

type GetRoomRes struct {
	Rooms []GetRoomDto `json:"rooms"`
}
