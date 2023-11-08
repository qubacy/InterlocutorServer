package base

import (
	"ilserver/pkg/token"
	storage "ilserver/storage/control"
	"time"
)

type Services struct {
	*AuthService
	*AdminService
	*TopicService
}

type Dependencies struct {
	AccessTokenTTL time.Duration
	Storage        storage.Storage
	TokenManager   token.Manager
}

// constructor
// -----------------------------------------------------------------------

func NewServices(deps Dependencies) *Services {
	return &Services{
		AuthService:  NewAuth(deps.AccessTokenTTL, deps.Storage, deps.TokenManager),
		AdminService: NewAdminService(deps.Storage),
		TopicService: NewTopicService(deps.Storage),
	}
}
