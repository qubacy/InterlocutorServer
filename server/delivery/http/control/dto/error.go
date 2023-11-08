package dto

import (
	"encoding/json"
	"ilserver/pkg/utility"
)

type Error struct {
	Text    string `json:"text"`
	Details string `json:"details"`
}

func MakeError(text, details string) Error {
	return Error{
		Text:    text,
		Details: details,
	}
}

// methods
// -----------------------------------------------------------------------

func (self Error) ToJson() (string, error) {
	bytes, err := json.Marshal(self)
	if err != nil {
		return "", utility.CreateCustomError(self.ToJson, err)
	}
	return string(bytes), nil
}
