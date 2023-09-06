package repository

import (
	"database/sql"
	"fmt"
	"ilserver/domain"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// -----------------------------------------------------------------------

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
			returnError =
				fmt.Errorf("Init, sql.Open err: %v", err)
			return
		}

		// ***

		inst =
			newDbFacade(db)

		// ***

		err = inst.createTopics()
		if err != nil {
			returnError =
				fmt.Errorf("Init, inst.createTopics err: %v", err)
			return
		}

		// ***

		err = inst.createAdmins()
		if err != nil {
			returnError =
				fmt.Errorf("Init, inst.createAdmins err: %v", err)
			return
		}
		err = inst.inflateAdminsIfNeed()
		if err != nil {
			returnError =
				fmt.Errorf("Init, inst.inflateAdminsIfNeed err: %v", err)
			return
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

// private
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

func (r *DbFacade) inflateAdminsIfNeed() error {
	err, num :=
		inst.RecordCountInTable("Admins")
	if err != nil {
		return fmt.Errorf("inflateAdminsIfNeed, inst.RecordCountInTable"+
			"err: %v", err)
	}

	// ***

	if num == 0 {
		err = inst.inflateAdmins()
		if err != nil {
			return fmt.Errorf("inflateAdminsIfNeed, inst.inflateAdmins "+
				"err: %v", err)
		}
	}

	return nil
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

// public
// -----------------------------------------------------------------------

func (r *DbFacade) RecordCountInTable(name string) (error, int) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare("SELECT count(*) AS RecordCount " +
		"FROM " + name + ";")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err), 0
	}
	defer stmt.Close() // ignore err!

	rows, err := stmt.Query()
	if err != nil {
		return fmt.Errorf("execute query failed %v", err), 0
	}
	defer rows.Close()

	// ***

	var recordCount int
	if rows.Next() {
		err = rows.Scan(&recordCount)
		if err != nil {
			return fmt.Errorf("scan next row with err %v", err), 0
		}
	} else {
		return fmt.Errorf("rows count is zero"), 0
	}

	return nil, recordCount
}

func (r *DbFacade) HasAdminByLogin(login string) (error, bool) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare(
		"SELECT count(*) AS RecordCount " +
			"FROM [Admins] WHERE [Admins].[Login] = ?;")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err), false
	}
	defer stmt.Close() // ignore err!

	// ***

	rows, err := stmt.Query(login)
	if err != nil {
		return fmt.Errorf("execute query failed %v", err), false
	}
	defer rows.Close()

	// ***

	var recordCount int
	if rows.Next() {
		err = rows.Scan(&recordCount)
		if err != nil {
			return fmt.Errorf("scan next row with err %v", err), false
		}
	} else {
		return fmt.Errorf("rows count is zero"), false
	}

	return nil, recordCount > 0
}

func (r *DbFacade) InsertTopic(topic domain.Topic) (error, int64) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare(
		"INSERT INTO [Topics] (Lang, Name) " +
			"VALUES (?, ?);")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err), 0
	}
	defer stmt.Close() // ignore err!

	// ***

	res, err := stmt.Exec(topic.Lang, topic.Name)
	if err != nil {
		return fmt.Errorf("execute query failed %v", err), 0
	}

	// ***

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id failed %v", err), 0
	}

	return nil, lastInsertId
}

func (r *DbFacade) UpdateAdminPass(login, newPass string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare(
		"UPDATE [Admins] SET [Pass] = ? " +
			"WHERE [Login] == ?;")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err)
	}
	defer stmt.Close()

	// ***

	_, err = stmt.Exec(newPass, login)
	if err != nil {
		return fmt.Errorf("execute query failed %v", err)
	}

	// ***

	return nil
}
