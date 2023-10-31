package control

import (
	"ilserver/domain"
	"ilserver/repository"
	"ilserver/transport/control/middleware"
	"ilserver/transport/controlDto"
)

type TopicService struct {
	dbFacade *repository.DbFacade
}

func NewTopicService() *TopicService {
	return &TopicService{
		dbFacade: repository.Instance(),
	}
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

func (ts *TopicService) GetTopics(dtoReq controlDto.GetTopicsReq) (error, controlDto.GetTopicsRes) {
	err, tokenOk, _ := middleware.ParseAndVerifyToken(dtoReq.AccessToken)
	if err != nil {
		return err, controlDto.GetTopicsRes{}
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
