package control

import (
	"context"
	"ilserver/service/control/dto"
)

type Auth interface {
	SignIn(ctx context.Context, dto dto.SignInInput) (dto.SignInOutput, error)
}

type Admin interface {
}

type Topic interface {
}

type Service interface {
	Auth
	Admin
	Topic
}
