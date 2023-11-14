package base

import (
	"context"
	"ilserver/service/control/dto"
	"ilserver/storage/game"
)

type RoomService struct {
	storage *game.Storage
}

func NewRoomService(storage *game.Storage) *RoomService {
	return &RoomService{
		storage: storage,
	}
}

// -----------------------------------------------------------------------

func (s *RoomService) GetRooms(ctx context.Context) (dto.GetRoomsOutput, error) {
	return dto.MakeGetRoomsOutputSuccess(
		s.storage.Rooms(),
	), nil
}
