package control

import (
	"context"
	"ilserver/service/control/dto"
)

type AuthService interface {
	SignIn(ctx context.Context, inp dto.SignInInput) (dto.SignInOutput, error)
}

type AdminService interface {
	GetAdmins(ctx context.Context) (dto.GetAdminsOutput, error)
	PostAdmin(ctx context.Context, inp dto.PostAdminInput) (dto.PostAdminOutput, error)
	DeleteAdminByIdr(ctx context.Context, idr int) (dto.DeleteAdminByIdrOutput, error)
}

type TopicService interface {
	GetTopics(ctx context.Context) (dto.GetTopicsOutput, error)
	GetRandomTopic(ctx context.Context, language int) (dto.GetRandomTopicOutput, error)

	PostTopic(ctx context.Context, inp dto.PostTopicInput) (dto.PostTopicOutput, error)
	PostTopics(ctx context.Context, inp dto.PostTopicsInput) (dto.PostTopicsOutput, error)

	DeleteTopicByIdr(ctx context.Context, idr int) (dto.DeleteTopicByIdrOutput, error)
	DeleteTopics(ctx context.Context) (dto.DeleteTopicsOutput, error)
}

type Services struct {
	AuthService
	AdminService
	TopicService
}
