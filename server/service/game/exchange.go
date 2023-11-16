package game

type AsyncResponse struct {
	ProfileId  string
	ServerBody interface{}
}

func MakeAsyncResponse(profileId string, serverBody interface{}) AsyncResponse {
	return AsyncResponse{
		ProfileId:  profileId,
		ServerBody: serverBody,
	}
}

// -----------------------------------------------------------------------

type AsyncResponseAboutError struct {
	ProfileId string
	Err       error
}

func MakeAsyncResponseAboutError(profileId string, err error) AsyncResponseAboutError {
	return AsyncResponseAboutError{
		ProfileId: profileId,
		Err:       err,
	}
}
