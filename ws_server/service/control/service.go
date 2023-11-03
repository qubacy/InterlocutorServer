package control

import (
	"context"
	"ilserver/domain"
	"ilserver/service/control/dto"
)

type AuthService interface {
	SignIn(ctx context.Context, inp dto.SignInInput) (dto.SignInOutput, error)
}

type AdminService interface {
	GetAdmins(ctx context.Context, accessToken string) domain.AdminList
	PostAdmin(ctx context.Context, inp dto.PostAdminInput) (dto.PostAdminOutput, error)
	DeleteAdminByIdr(ctx context.Context, accessToken string, idr int) (dto.DeleteAdminByIdrOutput, error)
}

type TopicService interface {
	GetTopics(ctx context.Context, accessToken string) (dto.GetTopicsOutput, error)
	GetRandomTopic(ctx context.Context, accessToken string) (dto.GetRandomTopicOutput, error)

	PostTopic(ctx context.Context, inp dto.PostTopicInput) (dto.PostTopicOutput, error)
	PostTopics(ctx context.Context, inp dto.PostTopicsInput) (dto.PostTopicsOutput, error)

	DeleteTopicByIdr(ctx context.Context, accessToken string, idr int) (dto.DeleteTopicByIdrOutput, error)
	DeleteTopics(ctx context.Context, accessToken string) (dto.DeleteTopicsOutput, error)
}

type Service interface {
	AuthService
	AdminService
	TopicService
}
