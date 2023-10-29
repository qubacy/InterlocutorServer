package sqlite

import (
	"database/sql"
	"fmt"
	"ilserver/utility"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var databaseDirectory = "../storage/"

func PathToDatabaseFile() string {
	return databaseDirectory +
		viper.GetString("storage.sql.sqlite.file")
}

func initialize() error {
	err := createDatabaseFile()
	if err != nil {
		return utility.CreateCustomError(initialize, err)
	}

	db, err := openDatabaseFile()
	if err != nil {
		return utility.CreateCustomError(initialize, err)
	}
	instance = newStorage(db)

	// TODO: добавить контекст с таймером?
	return initializeTables()
}

func initializeTables() error {
	err := instance.createTopics()
	if err != nil {
		return utility.CreateCustomError(initializeTables, err)
	}
	err = instance.createAdmins()
	if err != nil {
		return utility.CreateCustomError(initializeTables, err)
	}

	// ***

	err = instance.inflateAdminsIfNeeded()
	if err != nil {
		return utility.CreateCustomError(initializeTables, err)
	}

	return nil
}

// -----------------------------------------------------------------------

func openDatabaseFile() (*sql.DB, error) {
	if !exists(PathToDatabaseFile()) {
		return nil, utility.CreateCustomError(
			openDatabaseFile, fmt.Errorf("database file does not exist")) // base error!
	}

	return sql.Open(
		"sqlite3",
		PathToDatabaseFile(),
	)
}

func createDatabaseFile() error {
	if exists(PathToDatabaseFile()) {
		return nil
	}

	// ***

	file, err := os.Create(PathToDatabaseFile())
	if err != nil {
		return utility.CreateCustomError(
			createDatabaseFile, err)
	}
	defer file.Close()
	return nil
}

func removeDatabaseFile() error {
	if !exists(PathToDatabaseFile()) {
		return nil
	}

	// ***

	err := os.Remove(PathToDatabaseFile())
	if err != nil {
		return utility.CreateCustomError(
			removeDatabaseFile, err)
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
