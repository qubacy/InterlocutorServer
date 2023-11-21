package sqlite

import (
	"context"
	"database/sql"
	"ilserver/pkg/utility"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// -----------------------------------------------------------------------

// thread safety...
type Storage struct {
	mx sync.Mutex
	db *sql.DB
}

var once = sync.Once{}
var instance *Storage = nil
var initializationError error // user friendly error

// constructor
// -----------------------------------------------------------------------

// initialize and free in one thread.
func Instance() (*Storage, error) {
	once.Do(func() {

		duration := viper.GetDuration("storage.sql.initialization_timeout")
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel() // ?

		initializationError =
			initialize(ctx)
	})
	return instance, initializationError
}
func Free() {
	// TODO: what to do with the error?
	if instance != nil {
		instance.db.Close()
	}
}

func newStorage(db *sql.DB) *Storage {
	return &Storage{
		mx: sync.Mutex{},
		db: db,
	}
}

// general
// -----------------------------------------------------------------------

func (r *Storage) RecordCountInTable(ctx context.Context, name string) (int, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.PrepareContext(ctx,
		"SELECT count(*) AS RecordCount "+
			"FROM "+name+";",
	)
	if err != nil {
		return 0, utility.CreateCustomError(r.RecordCountInTable, err)
	}
	defer stmt.Close() // ignore err.

	// ***

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return 0, utility.CreateCustomError(r.RecordCountInTable, err)
	}
	defer rows.Close() // ignore err.

	// ***

	var recordCount int
	if rows.Next() {
		err = rows.Scan(&recordCount)
		if err != nil {
			return 0, utility.CreateCustomError(r.RecordCountInTable, err)
		}
	} else {
		return 0, utility.CreateCustomError(r.RecordCountInTable, err)
	}

	return recordCount, nil
}
