package base

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"
	"ilserver/storage"
)

type TopicService struct {
	storage      storage.Storage
	tokenManager token.Manager
}

func NewTopicService(storage storage.Storage, tokenManager token.Manager) *TopicService {
	return &TopicService{
		storage:      storage,
		tokenManager: tokenManager,
	}
}

// -----------------------------------------------------------------------

func (s *TopicService) GetTopics(ctx context.Context, accessToken string) (dto.GetTopicsOutput, error) {
	if len(accessToken) == 0 {
		// ---> 400
		return dto.MakeGetAdminsEmpty(),
			utility.CreateCustomError(s.GetAdmins,
				control.ErrAccessTokenIsEmpty())
	}

	err := s.tokenManager.Validate(accessToken)
	if err != nil {
		// ---> 400
		return dto.MakeGetAdminsEmpty(),
			utility.CreateCustomError(s.GetAdmins, err)
	}

	// ***

	var res controlDto.GetTopicsRes

	if !tokenOk {
		res.Some = []controlDto.TopicDto{}
		res.ErrorText = "access token expired"
		return nil, res
	}

	err, topics := repository.Instance().SelectTopics()
	if err != nil {
		return err, controlDto.GetTopicsRes{}
	}

	for i := range topics {
		res.ErrorText = ""
		res.Some = append(
			res.Some, controlDto.DomainToTopicDto(
				topics[i]))
	}

	return nil, res
}

func (s *TopicService) GetRandomTopic(ctx context.Context, accessToken string) (dto.GetRandomTopicOutput, error) {

}

func (s *TopicService) PostTopic(ctx context.Context, inp dto.PostTopicInput) (dto.PostTopicOutput, error) {

}

func (s *TopicService) PostTopics(ctx context.Context, inp dto.PostTopicsInput) (dto.PostTopicsOutput, error) {

}

func (s *TopicService) DeleteTopicByIdr(ctx context.Context, accessToken string, idr int) (dto.DeleteTopicByIdrOutput, error) {

}

func (s *TopicService) DeleteTopics(ctx context.Context, accessToken string) (dto.DeleteTopicsOutput, error) {

}

// -----------------------------------------------------------------------

func (ts *TopicService) PostTopic(dtoReq controlDto.PostTopicReq) (error, controlDto.PostTopicRes) {
	err, tokenOk, _ := middleware.ParseAndVerifyToken(dtoReq.AccessToken)
	if err != nil {
		return err, controlDto.PostTopicRes{}
	}

	// ***

	var res controlDto.PostTopicRes

	if !tokenOk {
		res.Idr = 0
		res.ErrorText = "access token expired"
		return nil, res
	}

	if len(dtoReq.Name) < 5 {
		res.Idr = 0
		res.ErrorText = "topic name is too short"
		return nil, res
	}

	err, idr := repository.Instance().InsertTopic(
		domain.Topic{
			Lang: dtoReq.Lang,
			Name: dtoReq.Name,
		},
	)
	if err != nil {
		return err, controlDto.PostTopicRes{}
	}

	// ***

	res.Idr = int(idr)
	res.ErrorText = ""

	return nil, res
}
