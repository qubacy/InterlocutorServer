package utility

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func IsFunction(i interface{}) bool {
	funType := reflect.TypeOf(i)
	funTypeName := funType.String()
	return strings.Contains(funTypeName, "func")
}

func GetFunctionName(i interface{}) string {
	if !IsFunction(i) {
		return ""
	}

	fullFunctionName := runtime.FuncForPC(
		reflect.ValueOf(
			i).Pointer()).Name()

	parts := strings.Split(fullFunctionName, "/")
	if len(parts) == 0 { // impossible?
		return ""
	}

	shortFunctionName := parts[len(parts)-1]
	return shortFunctionName
}

func CreateCustomError(i interface{}, err error) error {
	return fmt.Errorf(GetFunctionName(i)+"\n with an error/in: %w", err)
}
