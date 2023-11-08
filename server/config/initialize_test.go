package config

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func setUp() {
	Initialize()
}

func TestInit(t *testing.T) {
	err := Initialize()
	if err != nil {
		t.Fail()
	}
}

func TestParametersAvailability(t *testing.T) {
	setUp()

	dbFileName := viper.GetString("storage.sql.sqlite.file")
	fmt.Println("database file name:", dbFileName)

	if len(dbFileName) == 0 {
		t.Fail()
	}
}
