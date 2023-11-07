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

// help
// -----------------------------------------------------------------------

func ErrorToJson(errorObj Error) (string, error) {
	bytes, err := json.Marshal(errorObj)
	if err != nil {
		return "", utility.CreateCustomError(ErrorToJson, err)
	}
	return string(bytes), nil
}
