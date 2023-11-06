package dto

import "ilserver/domain"

type TopicDto struct {
	Idr  float64 `json:"idr"` // ---> uint64
	Lang float64 `json:"lang"`
	Name string  `json:"name"`
}

type PostTopicReq struct {
	AccessToken string  `json:"accessToken"`
	Lang        float64 `json:"lang"`
	Name        string  `json:"name"`
}

type PostTopicRes struct {
	Idr float64 `json:"idr"`
}

// converter
// -----------------------------------------------------------------------

func TopicDomainToDto(topic domain.Topic) TopicDto {
	return TopicDto{
		Idr:  float64(topic.Idr),
		Lang: float64(topic.Lang),
		Name: topic.Name,
	}
}

func TopicDtoToDomain(topic TopicDto) domain.Topic {
	return domain.Topic{
		Idr: int(topic.Idr),
	}
}

type UntypedTopicParts []interface{}

type PostTopicsReq struct {
	Some []UntypedTopicParts `json:"some"`
}

type GetTopicsReq struct {
	AccessToken string `json:"accessToken"`
}

type GetTopicsRes struct {
	ErrorText string     `json:"errorText"`
	Some      []TopicDto `json:"some"`
}
