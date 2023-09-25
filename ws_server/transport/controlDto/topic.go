package controlDto

import "ilserver/domain"

type TopicDto struct {
	Idr  int    `json:"idr"`
	Lang int    `json:"lang"`
	Name string `json:"name"`
}

type PostTopicReq struct {
	AccessToken string `json:"accessToken"`
	Lang        int    `json:"lang"`
	Name        string `json:"name"`
}

type PostTopicRes struct {
	ErrorText string `json:"errorText"`
	Idr       int    `json:"idr"`
}

// -----------------------------------------------------------------------

func DomainToTopicDto(topic domain.Topic) TopicDto {
	return TopicDto{
		Idr:  topic.Idr,
		Lang: topic.Lang,
		Name: topic.Name,
	}
}
