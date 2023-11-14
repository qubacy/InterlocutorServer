package dto

import (
	domain "ilserver/domain/memory"
)

// -----------------------------------------------------------------------

type GetRoomsOutput struct {
	Rooms domain.RoomList
}

// ---> 200
func MakeGetRoomsOutputSuccess(rooms domain.RoomList) GetRoomsOutput {
	return GetRoomsOutput{
		Rooms: rooms,
	}
}

// ---> 400
func MakeGetRoomsOutputEmpty() GetRoomsOutput {
	return GetRoomsOutput{}
}
