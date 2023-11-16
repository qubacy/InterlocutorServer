package dto

import (
	domain "ilserver/domain/memory"
)

// parts
// -----------------------------------------------------------------------

// can remove the postfix 'dto'
type GetRoomDto struct {
	Id           string      `json:"id"` // ---> uint64
	State        interface{} `json:"state"`
	CreationTime float64     `json:"creation-time"`
	Profiles     []Profile   `json:"profile-list"`
	Language     float64     `json:"language"`
}

func MakeGetRoomDto(room domain.Room) GetRoomDto {
	return GetRoomDto{
		Id:           room.Id,
		State:        MakeRoomState(room.State),
		CreationTime: float64(room.CreationTime.Unix()),
		Profiles:     MakeProfiles(room.Profiles),
		Language:     float64(room.Language),
	}
}

// req/res
// -----------------------------------------------------------------------

type GetRoomReq struct{}

type GetRoomRes struct {
	Rooms []GetRoomDto `json:"rooms"`
}

func MakeGetRoomRes(rooms domain.RoomList) GetRoomRes {
	result := GetRoomRes{}
	for i := range rooms {
		result.Rooms = append(result.Rooms,
			MakeGetRoomDto(rooms[i]))
	}
	return result
}
