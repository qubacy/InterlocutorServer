package dto

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
