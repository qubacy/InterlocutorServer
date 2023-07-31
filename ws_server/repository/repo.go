package repository

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Repo struct {
	db *sql.DB
}

var once sync.Once
var inst *Repo = nil

func Inst() *Repo {
	once.Do(func() {

	})
}

func newRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func Init() error {
	db, err := sql.Open("sqlite3", "../repository/store.db")
	if err != nil {
		return err
	}

	Repo := MakeRepo(db)
	Repo.createTopics()

	return nil
}

func (r *Repo) createTopics() error {
	tq :=
		"CREATE TABLE Topics( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Lang INTEGER, " +
			"    Name TEXT, " +
			"); "

	_, err := r.db.Exec(tq)
	return err
}
