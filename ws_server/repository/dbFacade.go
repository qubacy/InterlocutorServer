package repository

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

type DbFacade struct {
	mx sync.Mutex
	db *sql.DB
}

var inst *DbFacade = nil

func Instance() *DbFacade {
	return inst
}

// -----------------------------------------------------------------------

var once = sync.Once{}

func Init() error {
	var returnError error = nil
	once.Do(func() {
		db, err := sql.Open("sqlite3", "../repository/"+
			viper.GetString("database.file"))
		if err != nil {
			returnError = fmt.Errorf("Init, sql.Open err: %v", err)
		}

		// ***

		inst =
			newDbFacade(db)

		// ***

		err = inst.createTopics()
		if err != nil {
			returnError = fmt.Errorf("Init, inst.createTopics err: %v", err)
		}

		err = inst.createAdmins()
		if err != nil {
			returnError = fmt.Errorf("Init, inst.createAdmins err: %v", err)
		}

		err = inst.inflateAdmins()
		if err != nil {
			returnError = fmt.Errorf("Init, inst.inflateAdmins err: %v", err)
		}
	})

	return returnError
}

func newDbFacade(db *sql.DB) *DbFacade {
	return &DbFacade{
		mx: sync.Mutex{},
		db: db,
	}
}

// -----------------------------------------------------------------------

func (r *DbFacade) createTopics() error {
	tq :=
		"CREATE TABLE IF NOT EXISTS Topics( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Lang INTEGER NOT NULL, " +
			"    Name TEXT NOT NULL " +
			"); "

	_, err := r.db.Exec(tq)
	return err
}

func (r *DbFacade) createAdmins() error {
	tq :=
		"CREATE TABLE IF NOT EXISTS Admins( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Login TEXT UNIQUE NOT NULL, " +
			"    Pass TEXT NOT NULL, " +
			"    RefreshTokenHash TEXT NULL" +
			"); "

	_, err := r.db.Exec(tq)
	return err
}

func (r *DbFacade) inflateAdmins() error {
	stmt, err := r.db.Prepare(
		"INSERT INTO [Admins] (Login, Pass) " +
			"VALUES (?, ?);")

	if err != nil {
		return err
	}

	login := viper.GetString("database.default_admin_entry.login")
	pass := viper.GetString("database.default_admin_entry.pass")

	_, err = stmt.Exec(login, pass)
	return err
}

// -----------------------------------------------------------------------

func (r *DbFacade) hasAdminByLogin(login string) (error, bool) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	return nil, true
}
