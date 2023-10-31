package sqlite

import (
	"context"
	"ilserver/pkg/utility"

	"github.com/spf13/viper"
)

// create
// -----------------------------------------------------------------------

func (r *Storage) createTopics(ctx context.Context) error {
	tq :=
		"CREATE TABLE IF NOT EXISTS Topics ( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Lang INTEGER NOT NULL, " +
			"    Name TEXT NOT NULL " +
			"); "

	_, err := r.db.ExecContext(ctx, tq)
	if err != nil {
		return utility.CreateCustomError(
			r.createTopics, err)
	}
	return nil
}

func (r *Storage) createAdmins(ctx context.Context) error {
	tq :=
		"CREATE TABLE IF NOT EXISTS Admins ( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Login TEXT UNIQUE NOT NULL, " +
			"    Pass TEXT NOT NULL " +
			"); "

	_, err := r.db.ExecContext(ctx, tq)
	if err != nil {
		return utility.CreateCustomError(
			r.createAdmins, err)
	}
	return nil
}

// inflate
// -----------------------------------------------------------------------

func (r *Storage) inflateAdminsIfNeeded(ctx context.Context) error {
	num, err := instance.RecordCountInTable(ctx, "Admins")
	if err != nil {
		return utility.CreateCustomError(
			r.inflateAdminsIfNeeded, err)
	}

	// ***

	if num == 0 {
		err = instance.inflateAdmins(ctx)
		if err != nil {
			return utility.CreateCustomError(
				r.inflateAdminsIfNeeded, err)
		}
	}

	return nil
}

func (r *Storage) inflateAdmins(ctx context.Context) error {
	stmt, err := r.db.PrepareContext(ctx,
		"INSERT INTO [Admins] (Login, Pass) "+
			"VALUES (?, ?);")

	if err != nil {
		return utility.CreateCustomError(
			r.inflateAdmins, err)
	}

	login := viper.GetString("storage.default_admin_entry.login")
	pass := viper.GetString("storage.default_admin_entry.pass")

	_, err = stmt.ExecContext(ctx, login, pass)
	if err != nil {
		return utility.CreateCustomError(
			r.inflateAdmins, err)
	}

	return nil
}
