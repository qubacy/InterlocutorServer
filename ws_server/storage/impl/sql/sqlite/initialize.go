package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"ilserver/utility"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var databaseDirectory = "../storage/" // ?

func PathToDatabaseFile() string {
	return databaseDirectory +
		viper.GetString("storage.sql.sqlite.file")
}

func initialize(ctx context.Context) error {
	err := createDatabaseFile()
	if err != nil {
		return utility.CreateCustomError(initialize, err)
	}

	// ***

	db, err := openDatabaseFile()
	if err != nil {
		return utility.CreateCustomError(initialize, err)
	}
	instance = newStorage(db)

	// ***

	// TODO: добавить контекст с таймером?
	return initializeTables(ctx)
}

func initializeTables(ctx context.Context) error {
	err := instance.createTopics(ctx)
	if err != nil {
		return utility.CreateCustomError(initializeTables, err)
	}
	err = instance.createAdmins(ctx)
	if err != nil {
		return utility.CreateCustomError(initializeTables, err)
	}

	// ***

	err = instance.inflateAdminsIfNeeded(ctx)
	if err != nil {
		return utility.CreateCustomError(initializeTables, err)
	}

	return nil
}

// work with file system
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
