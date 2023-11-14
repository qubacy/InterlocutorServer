package base

import (
	"ilserver/pkg/token"
	storage "ilserver/storage/control"
	"ilserver/storage/game"
	"time"
)

type Services struct {
	*AuthService
	*AdminService
	*TopicService
	*RoomService
}

type Dependencies struct {
	AccessTokenTTL time.Duration
	Storage        storage.Storage // <--- controlStorage
	GameStorage    *game.Storage
	TokenManager   token.Manager
}

// constructor
// -----------------------------------------------------------------------

func NewServices(deps Dependencies) *Services {
	return &Services{
		AuthService:  NewAuth(deps.AccessTokenTTL, deps.Storage, deps.TokenManager),
		AdminService: NewAdminService(deps.Storage),
		TopicService: NewTopicService(deps.Storage),
		RoomService:  NewRoomService(deps.GameStorage),
	}
}
