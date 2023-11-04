package control

import "fmt"

// Mixed approach: 400 - business logic errors

func ErrInputDtoIsInvalid() string {
	return fmt.Sprintf("Input data is invalid")
}

func ErrAccessTokenIsEmpty() string {
	return fmt.Sprintf("Access token is empty")
}

func ErrLoginNotFoundOrPasswordIsIncorrect(login string) string {
	return fmt.Sprintf("Login '%v' not found or password is incorrect", login)
}

func ErrLoginOrPasswordIsTooShort() string {
	return fmt.Sprintf("Login or password is too short")
}
