package dto

import (
	"encoding/json"
	"ilserver/pkg/utility"
)

type KeyValueStorage map[string]interface{}

type Pack struct {
	Operation int             `json:"operation"`
	RawBody   KeyValueStorage `json:"body"`
}

// constructor
// -----------------------------------------------------------------------

func MakePack(op int, rawBody KeyValueStorage) Pack {
	return Pack{
		Operation: op,
		RawBody:   rawBody,
	}
}

func MakeEmptyPack() Pack {
	return Pack{}
}

func MakePackFromJson(data []byte) (Pack, error) {
	pack := Pack{}
	err := json.Unmarshal(data, &pack)
	if err != nil {
		return MakeEmptyPack(), utility.CreateCustomError(MakePackFromJson, err)
	}
	return pack, nil
}

func MakePackFromAnyBody(operation int, body interface{}) (Pack, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return MakeEmptyPack(),
			utility.CreateCustomError(MakePackFromAnyBody, err)
	}

	rawBody := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &rawBody)
	if err != nil {
		return MakeEmptyPack(),
			utility.CreateCustomError(MakePackFromAnyBody, err)
	}

	return MakePack(operation, rawBody), nil
}

func MakePackAsJsonBytes(operation int, body interface{}) ([]byte, error) {
	pack, err := MakePackFromAnyBody(operation, body)
	if err != nil {
		return []byte{},
			utility.CreateCustomError(MakePackAsJsonBytes, err)
	}

	return pack.ToJsonBytes()
}

// methods
// -----------------------------------------------------------------------

func (self *Pack) ToJsonBytes() ([]byte, error) {
	bytes, err := json.Marshal(*self) // <--- value
	if err != nil {
		return []byte{},
			utility.CreateCustomError(self.ToJsonBytes, err)
	}

	return bytes, nil
}
