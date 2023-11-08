package control

import "fmt"

// Mixed approach: 400 - business logic errors

func ErrInputDtoIsInvalid() error {
	return fmt.Errorf("Input data is invalid")
}

func ErrLoginNotFoundOrPasswordIsIncorrect(login string) error {
	return fmt.Errorf("Login '%v' not found or password is incorrect", login)
}

func ErrLoginOrPasswordIsTooShort() error {
	return fmt.Errorf("Login or password is too short")
}
