package control

import "fmt"

func LoginNotFoundOrPasswordIsIncorrect(login string) string {
	return fmt.Sprintf("Login '%v' not found or password is incorrect", login)
}

func LoginOrPasswordIsTooShort() string {
	return fmt.Sprintf("Login or password is too short")
}
