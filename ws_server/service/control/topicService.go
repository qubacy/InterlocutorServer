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
