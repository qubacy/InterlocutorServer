package dto

import "ilserver/domain"

/*

The default concrete Go types are:
- bool for JSON booleans;
- float64 for JSON numbers;
- string for JSON strings;
- nil for JSON null.

*/

type GetTopicDto struct {
	Idr  float64 `json:"idr"` // ---> uint64
	Lang float64 `json:"lang"`
	Name string  `json:"name"`
}

type PostTopicDto struct {
	Lang float64 `json:"lang"`
	Name string  `json:"name"`
}

// req/res
// -----------------------------------------------------------------------

type GetTopicReq struct{}

type GetTopicRes struct {
	Topics []GetTopicDto `json:"topics"`
}

type PostTopicReq struct {
	Topics []PostTopicDto `json:"topics"`
}

type PostTopicRes struct{}

// constructor
// -----------------------------------------------------------------------

func MakeGetTopicDto(topic domain.Topic) GetTopicDto {
	return GetTopicDto{
		Idr:  float64(topic.Idr),
		Lang: float64(topic.Lang),
		Name: topic.Name,
	}
}

func MakeGetTopicRes(topics domain.TopicList) GetTopicRes {
	result := GetTopicRes{}
	for i := range topics {
		result.Topics = append(result.Topics,
			MakeGetTopicDto(topics[i]))
	}
	return result
}

// converter
// -----------------------------------------------------------------------

func (self GetTopicDto) ToDomain() domain.Topic {
	return domain.Topic{
		Idr:  int(self.Idr),
		Lang: int(self.Lang),
		Name: self.Name,
	}
}

func (self PostTopicDto) ToDomain() domain.Topic {
	return domain.Topic{
		Idr:  0, // <--- not initialized!
		Lang: int(self.Lang),
		Name: self.Name,
	}
}

func (self *PostTopicReq) ToDomainList() domain.TopicList {
	result := domain.TopicList{}
	for i := range self.Topics {
		result = append(result,
			self.Topics[i].ToDomain())
	}
	return result
}
