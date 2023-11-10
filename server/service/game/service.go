package game

import (
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
