package game

import (
	"ilserver/service/game/dto"
	"ilserver/storage/control"
	"ilserver/storage/game"
)

type Service struct {
	gameStorage    *game.Storage // <--- impl
	controlStorage control.Storage
}

func NewService(
	gameStorage *game.Storage,
	controlStorage control.Storage,
) *Service {
	return &Service{
		gameStorage:    gameStorage,
		controlStorage: controlStorage,
	}
}

// public
// -----------------------------------------------------------------------

func (s *Service) SearchingStart(profileId string, body dto.CliSearchingStartBody) (
	dto.SvrSearchingStartBody, error,
) {
	if !body.IsValid() {
		return dto.MakeSvrSearchingStartBodyEmpty(), ErrInvalidClientPackBody
	}

	s.gameStorage.InsertRoomWithSearchingState(body.Profile.Language)

	return dto.SvrSearchingStartBody{}, nil
}
