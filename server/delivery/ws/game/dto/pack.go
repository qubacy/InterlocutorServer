package dto

import (
	"encoding/json"
	"ilserver/pkg/utility"
	serviceDto "ilserver/service/game/dto"
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

// converter
// -----------------------------------------------------------------------

func (self *Pack) ToJsonBytes() ([]byte, error) {
	bytes, err := json.Marshal(*self) // <--- value
	if err != nil {
		return []byte{},
			utility.CreateCustomError(self.ToJsonBytes, err)
	}

	return bytes, nil
}

func (self *Pack) BodyToJsonBytes() ([]byte, error) {
	bytes, err := json.Marshal(self.RawBody)
	if err != nil {
		return []byte{},
			utility.CreateCustomError(self.BodyToJsonBytes, err)
	}

	return bytes, nil
}

// to concrete body
// -----------------------------------------------------------------------

func (self *Pack) AsCliSearchingStartBody() (
	serviceDto.CliSearchingStartBody, error,
) {
	data, err := self.BodyToJsonBytes()
	if err != nil {
		return serviceDto.CliSearchingStartBody{},
			utility.CreateCustomError(self.AsCliSearchingStartBody, err)
	}

	// error ignores one level!
	return serviceDto.MakeCliSearchingStartBodyFromJson(data)
}

func (self *Pack) AsCliSearchingStopBody() (
	serviceDto.CliSearchingStopBody, error,
) {
	data, err := self.BodyToJsonBytes()
	if err != nil {
		return serviceDto.CliSearchingStopBody{},
			utility.CreateCustomError(self.AsCliSearchingStopBody, err)
	}

	return serviceDto.MakeCliSearchingStopBodyFromJson(data)
}

func (self *Pack) AsCliChattingNewMessageBody() (
	serviceDto.CliChattingNewMessageBody, error,
) {
	data, err := self.BodyToJsonBytes()
	if err != nil {
		return serviceDto.MakeCliChattingNewMessageBodyEmpty(),
			utility.CreateCustomError(self.AsCliChattingNewMessageBody, err)
	}

	return serviceDto.MakeCliChattingNewMessageBodyFromJson(data)
}

func (self *Pack) AsCliChoosingUsersChosenBody() (
	serviceDto.CliChoosingUsersChosenBody, error,
) {
	data, err := self.BodyToJsonBytes()
	if err != nil {
		return serviceDto.MakeCliChoosingUsersChosenBodyEmpty(),
			utility.CreateCustomError(self.AsCliChoosingUsersChosenBody, err)
	}

	return serviceDto.MakeCliChoosingUsersChosenBodyFromJson(data)
}
