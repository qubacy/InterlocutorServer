package dto

import "ilserver/domain"

type GetTopicsOutput struct {
	Topics domain.TopicList
}

type GetRandomTopicOutput struct {
	Topic domain.Topic
}

// -----------------------------------------------------------------------

type PostTopicInput struct {
	AccessToken string
	Lang        int
	Name        string
}

type PostTopicOutput struct {
	Idr int64
}

// -----------------------------------------------------------------------

type PostTopicsInput struct {
	AccessToken string
	Topics      domain.TopicList
}

type PostTopicsOutput struct{}

// -----------------------------------------------------------------------

type DeleteTopicByIdrOutput struct{}

type DeleteTopicsOutput struct{}
