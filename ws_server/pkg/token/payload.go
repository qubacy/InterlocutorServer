package token

type Payload struct {
	Subject string
}

func MakePayload(sub string) Payload {
	return Payload{Subject: sub}
}
