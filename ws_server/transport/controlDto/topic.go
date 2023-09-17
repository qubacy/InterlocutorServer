package controlDto

type PostTopicReq struct {
	AccessToken string `json:"accessToken"`
	Lang        int    `json:"lang"`
	Name        string `json:"name"`
}

type PostTopicRes struct {
	ErrorText string `json:"errorText"`
	Idr       int    `json:"idr"`
}
