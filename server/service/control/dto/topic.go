package dto

import "ilserver/domain"

// -----------------------------------------------------------------------

type GetTopicsOutput struct {
	Topics domain.TopicList
}

// ---> 200
func MakeGetTopicsSuccess(topics domain.TopicList) GetTopicsOutput {
	return GetTopicsOutput{
		Topics: topics,
	}
}

// ---> 400
func MakeGetTopicsEmpty() GetTopicsOutput {
	return GetTopicsOutput{}
}

// -----------------------------------------------------------------------

type GetRandomTopicOutput struct {
	Topic domain.Topic
}

// ---> 200
func MakeGetRandomTopicSuccess(topic domain.Topic) GetRandomTopicOutput {
	return GetRandomTopicOutput{
		Topic: topic,
	}
}

// ---> 400
func MakeGetRandomTopicEmpty() GetRandomTopicOutput {
	return GetRandomTopicOutput{}
}

// -----------------------------------------------------------------------

type PostTopicInput struct {
	Lang int
	Name string
}

type PostTopicOutput struct {
	Idr int64
}

func (self PostTopicInput) ToDomain() domain.Topic {
	return domain.Topic{
		Idr:  0,
		Lang: self.Lang,
		Name: self.Name,
	}
}

// ---> 200
func MakePostTopicOutputSuccess(idr int64) PostTopicOutput {
	return PostTopicOutput{
		Idr: idr,
	}
}

// ---> 400
func MakePostTopicOutputEmpty() PostTopicOutput {
	return PostTopicOutput{}
}

// -----------------------------------------------------------------------

type PostTopicsInput struct {
	Topics domain.TopicList
}

type PostTopicsOutput struct{}

// -----------------------------------------------------------------------

type DeleteTopicByIdrOutput struct{}
type DeleteTopicsOutput struct{}
