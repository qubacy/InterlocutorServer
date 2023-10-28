package sqlite

import (
	"fmt"

	"github.com/spf13/viper"
)

// returns an error without processing...

func (r *Storage) createTopics() error {
	tq :=
		"CREATE TABLE IF NOT EXISTS Topics ( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Lang INTEGER NOT NULL, " +
			"    Name TEXT NOT NULL " +
			"); "

	_, err := r.db.Exec(tq)
	return err
}

func (r *Storage) createAdmins() error {
	tq :=
		"CREATE TABLE IF NOT EXISTS Admins ( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Login TEXT UNIQUE NOT NULL, " +
			"    Pass TEXT NOT NULL " +
			"); "

	_, err := r.db.Exec(tq)
	return err
}

// -----------------------------------------------------------------------

func (r *Storage) inflateAdminsIfNeeded() error {
	err, num := instance.RecordCountInTable("Admins")
	if err != nil {
		return fmt.Errorf("")
	}

	// ***

	if num == 0 {
		err = instance.inflateAdmins()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Storage) inflateAdmins() error {
	stmt, err := r.db.Prepare(
		"INSERT INTO [Admins] (Login, Pass) " +
			"VALUES (?, ?);")

	if err != nil {
		return err
	}

	login := viper.GetString("storage.default_admin_entry.login")
	pass := viper.GetString("database.default_admin_entry.pass")

	_, err = stmt.Exec(login, pass)
	return err
}
