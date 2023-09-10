package controlDto

type PostTopicReq struct {
	AccessToken string
	Lang        int
	Name        int
}

type PostTopicRes struct {
	Idr int
}
