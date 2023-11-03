package dto

import "ilserver/domain"

type GetTopicsOutput struct {
	Topics    domain.TopicList
	ErrorText string
}

type GetRandomTopicOutput struct {
	Topic     domain.Topic
	ErrorText string
}

// -----------------------------------------------------------------------

type PostTopicInput struct {
	AccessToken string
	Lang        int
	Name        string
}

type PostTopicOutput struct {
	Idr       int
	ErrorText string
}

// -----------------------------------------------------------------------

type PostTopicsInput struct {
	AccessToken string
	Topics      domain.TopicList
}

type PostTopicsOutput struct {
	ErrorText string
}

// -----------------------------------------------------------------------

type DeleteTopicByIdrOutput struct {
	ErrorText string
}

type DeleteTopicsOutput struct {
	ErrorText string
}
