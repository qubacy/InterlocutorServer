package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var databaseDirectory = "../storage/"

func initialize() error {
	db, err := sql.Open("sqlite3", databaseDirectory+
		viper.GetString("storage.sql.sqlite.file"))

	if err != nil {
		return fmt.Errorf("sqlite open failed with error: %w", err)
	}
	instance = newStorage(db)

	// TODO: добавить контекст с таймером
	return initializeTables()
}

func initializeTables() error {
	err := instance.createTopics()
	if err != nil {
		return fmt.Errorf("sqlite create topics failed with error: %v", err)
	}
	err = instance.createAdmins()
	if err != nil {
		return fmt.Errorf("sqlite create admins failed with error: %v", err)
	}

	// ***

	err = instance.inflateAdminsIfNeeded()
	if err != nil {
		return fmt.Errorf("sqlite inflate admins failed with error: %v", err)
	}

	return nil
}
