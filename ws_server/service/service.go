package service

import (
	"context"
	"ilserver/delivery/http/control/dto"
)

type Auth interface {
	SignIn(ctx context.Context, dto dto.SignInReq) (dto.SignInRes, error)
}

type Topic interface {
}

type Room interface {
}

type Service interface {
	Auth
	Topic
	Room
}
