package dto

import "ilserver/domain"

/*

The default concrete Go types are:
- bool for JSON booleans;
- float64 for JSON numbers;
- string for JSON strings;
- nil for JSON null.

*/

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
	ErrorText string  `json:"errorText"`
	Idr       float64 `json:"idr"`
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
