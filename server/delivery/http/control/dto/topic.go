package dto

import (
	"encoding/json"
	"ilserver/domain"
	service "ilserver/service/control/dto"
)

/*

The default concrete Go types are:
- bool for JSON booleans;
- float64 for JSON numbers;
- string for JSON strings;
- nil for JSON null.

*/

// -----------------------------------------------------------------------

type GetTopicDto struct {
	Idr  float64 `json:"idr"` // ---> uint64
	Lang float64 `json:"lang"`
	Name string  `json:"name"`
}

func MakeGetTopicDto(topic domain.Topic) GetTopicDto {
	return GetTopicDto{
		Idr:  float64(topic.Idr),
		Lang: float64(topic.Lang),
		Name: topic.Name,
	}
}

type PostTopicDto struct {
	Lang float64 `json:"lang"`
	Name string  `json:"name"`
}

func (self *PostTopicDto) ToDomain() domain.Topic {
	return domain.Topic{
		Name: self.Name,
		Lang: int(self.Lang),
	}
}

// req/res
// -----------------------------------------------------------------------

type GetTopicReq struct{}

type GetTopicRes struct {
	Topics []GetTopicDto `json:"topics"`
}

func MakeGetTopicRes(topics domain.TopicList) GetTopicRes {
	result := GetTopicRes{}
	for i := range topics {
		result.Topics = append(result.Topics,
			MakeGetTopicDto(topics[i]))
	}
	return result
}

// -----------------------------------------------------------------------

type PostTopicReq struct {
	Topics []PostTopicDto `json:"topics"`
}

func MakePostTopicReqFromJson(rawJson []byte) (PostTopicReq, error) {
	result := PostTopicReq{}
	err := json.Unmarshal(rawJson, &result)
	if err != nil {
		return MakePostTopicReqEmpty(), err
	}

	return result, nil
}

func MakePostTopicReqEmpty() PostTopicReq {
	return PostTopicReq{}
}

func (self *PostTopicReq) ToServiceInp() service.PostTopicsInput {
	result := service.PostTopicsInput{}
	for i := range self.Topics {
		result.Topics = append(result.Topics,
			self.Topics[i].ToDomain())
	}

	return result
}

type PostTopicRes struct{}
