package base

import (
	"ilserver/pkg/token"
	"ilserver/storage"
)

type Services struct {
	*AuthService
	*AdminService
	*TopicService
}

func NewServices(storage storage.Storage, tokenManager token.Manager) *Services {
	return &Services{
		AuthService:  NewAuth(storage, tokenManager),
		AdminService: NewAdminService(storage, tokenManager),
		TopicService: NewTopicService(storage, tokenManager),
	}
}
