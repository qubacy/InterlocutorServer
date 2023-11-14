package base

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"
	storage "ilserver/storage/control"
)

type TopicService struct {
	storage storage.Storage
}

func NewTopicService(storage storage.Storage) *TopicService {
	return &TopicService{
		storage: storage,
	}
}

// read
// -----------------------------------------------------------------------

func (s *TopicService) GetTopics(ctx context.Context) (dto.GetTopicsOutput, error) {
	topics, err := s.storage.AllTopics(ctx)
	if err != nil {
		// ---> 400
		return dto.MakeGetTopicsOutputEmpty(),
			utility.CreateCustomError(s.GetTopics, err)
	}

	// ---> 200
	return dto.MakeGetTopicsOutputSuccess(topics), nil
}

func (s *TopicService) GetRandomTopic(ctx context.Context, language int) (dto.GetRandomTopicOutput, error) {
	topic, err := s.storage.RandomTopic(ctx, language)
	if err != nil {
		// ---> 400
		return dto.MakeGetRandomTopicOutputEmpty(),
			utility.CreateCustomError(s.GetRandomTopic, err)
	}

	return dto.MakeGetRandomTopicOutputSuccess(topic), nil
}

// insert
// -----------------------------------------------------------------------

func (s *TopicService) PostTopic(ctx context.Context, inp dto.PostTopicInput) (dto.PostTopicOutput, error) {
	if !isValidPostTopicInput(inp) {
		// ---> 400
		return dto.MakePostTopicOutputEmpty(),
			utility.CreateCustomError(s.PostTopic,
				control.ErrInputDtoIsInvalid())
	}

	// ***

	idr, err := s.storage.InsertTopic(ctx, inp.ToDomain())
	if err != nil {
		return dto.MakePostTopicOutputEmpty(),
			utility.CreateCustomError(s.PostTopic, err)
	}

	return dto.PostTopicOutput{Idr: idr}, nil
}

func (s *TopicService) PostTopics(ctx context.Context, inp dto.PostTopicsInput) (dto.PostTopicsOutput, error) {
	if !isValidPostTopicsInput(inp) {
		// ---> 400
		return dto.PostTopicsOutput{},
			utility.CreateCustomError(s.PostTopics,
				control.ErrInputDtoIsInvalid())
	}

	// ***

	err := s.storage.InsertTopics(ctx, inp.Topics)
	if err != nil {
		// ---> 400
		return dto.PostTopicsOutput{},
			utility.CreateCustomError(s.PostTopics,
				control.ErrInputDtoIsInvalid())
	}

	// ---> 200
	return dto.PostTopicsOutput{}, nil
}

// delete
// -----------------------------------------------------------------------

func (s *TopicService) DeleteTopicByIdr(ctx context.Context, idr int) (dto.DeleteTopicByIdrOutput, error) {
	err := s.storage.DeleteTopic(ctx, idr)
	if err != nil {
		// ---> 400
		return dto.DeleteTopicByIdrOutput{},
			utility.CreateCustomError(s.DeleteTopicByIdr,
				control.ErrInputDtoIsInvalid())
	}

	// ---> 200
	return dto.DeleteTopicByIdrOutput{}, nil
}

func (s *TopicService) DeleteTopics(ctx context.Context) (dto.DeleteTopicsOutput, error) {
	err := s.storage.DeleteTopics(ctx)
	if err != nil {
		// ---> 400
		return dto.DeleteTopicsOutput{},
			utility.CreateCustomError(s.DeleteTopicByIdr,
				control.ErrInputDtoIsInvalid())
	}

	// ---> 200
	return dto.DeleteTopicsOutput{}, nil
}

// private
// -----------------------------------------------------------------------

func isValidPostTopicInput(inp dto.PostTopicInput) bool {
	return isValidTopic(
		inp.ToDomain())
}

func isValidPostTopicsInput(inp dto.PostTopicsInput) bool {
	for i := range inp.Topics {
		topic := inp.Topics[i]
		if !isValidTopic(topic) {
			return false
		}
	}
	return true
}

func isValidTopic(topic domain.Topic) bool {
	if topic.Lang < 0 {
		return false
	}
	if len(topic.Name) <= 3 {
		return false
	}
	return true
}
