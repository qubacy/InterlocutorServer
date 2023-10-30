package overWsDto

import "encoding/json"

type Pack struct {
	Operation int                    `json:"operation"`
	RawBody   map[string]interface{} `json:"body"`
}

// -----------------------------------------------------------------------

func MakePack(operation int, body interface{}) Pack {
	bodyBytes, _ := json.Marshal(body)

	rawBody := make(map[string]interface{})
	_ = json.Unmarshal(bodyBytes, &rawBody)

	return Pack{
		Operation: operation,
		RawBody:   rawBody,
	}
}

func MakePackBytes(operation int, body interface{}) []byte {
	pack := MakePack(operation, body)

	// TODO: так легко игнорировать ошибки (?)
	packBytes, _ := json.Marshal(pack)
	return packBytes
}
