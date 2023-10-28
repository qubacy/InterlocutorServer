package sqlite

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestHelloEmpty(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fail()
	}

	fmt.Printf("path: %v\n", path)
	pathParts := strings.Split(path, "\\")
	pathParts = pathParts[0 : len(pathParts)-3]
	path = strings.Join(pathParts, "\\")
	fmt.Printf("updated path: %v\n", path)

	databaseDirectory = path
	fmt.Printf("database directory: %v\n",
		databaseDirectory)

	_, err = Instance()
	if err != nil {
		fmt.Println("err:", err)
		t.Fail()
	}
}
