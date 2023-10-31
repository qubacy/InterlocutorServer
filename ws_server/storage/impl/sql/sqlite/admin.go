package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"ilserver/domain"
	"ilserver/utility"
)

func (self *Storage) InsertAdmin(ctx context.Context, admin domain.Admin) (int64, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"INSERT INTO Admins (Login, Pass) "+
			"VALUES (?, ?);",
	)
	if err != nil {
		return 0, utility.CreateCustomError(self.InsertAdmin, err)
	}
	defer stmt.Close()

	// ***

	res, err := stmt.ExecContext(ctx, admin.Login, admin.Pass)
	if err != nil {
		return 0, utility.CreateCustomError(self.InsertAdmin, err)
	}

	// ***

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return 0, utility.CreateCustomError(self.InsertAdmin, err)
	}
	return lastInsertedId, nil
}

func (self *Storage) HasAdminByLogin(ctx context.Context, login string) (bool, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"SELECT count(*) AS RecordCount "+
			"FROM Admins WHERE Admins.Login = ?;",
	)
	if err != nil {
		return false, utility.CreateCustomError(self.HasAdminByLogin, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx, login)
	if err != nil {
		return false, utility.CreateCustomError(self.HasAdminByLogin, err)
	}
	defer rows.Close()

	// ***

	recordCount, err := sqlRowsToRecordCount(rows)
	if err != nil {
		return false, utility.CreateCustomError(self.HasAdminByLogin, err)
	}
	return recordCount > 0, nil
}

func (self *Storage) HasAdminWithLoginAndPassword(ctx context.Context, login, password string) (bool, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"SELECT count(*) AS RecordCount "+
			"FROM Admins "+
			"WHERE Admins.Login = ? "+
			"AND Admins.Pass = ?;",
	)
	if err != nil {
		return false, utility.CreateCustomError(self.HasAdminWithLoginAndPassword, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx, login, password)
	if err != nil {
		return false, utility.CreateCustomError(self.HasAdminWithLoginAndPassword, err)
	}
	defer rows.Close()

	// ***

	recordCount, err := sqlRowsToRecordCount(rows)
	if err != nil {
		return false, utility.CreateCustomError(self.HasAdminWithLoginAndPassword, err)
	}
	return recordCount > 0, nil
}

func (self *Storage) AllAdmins(ctx context.Context) (domain.AdminList, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx, "SELECT * FROM Admins;")
	if err != nil {
		return nil, utility.CreateCustomError(self.AllAdmins, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, utility.CreateCustomError(self.AllAdmins, err)
	}
	defer rows.Close()

	// ***

	admins, err := sqlRowsToAdmins(rows)
	if err != nil {
		return nil, utility.CreateCustomError(self.AllAdmins, err)
	}
	return admins, nil
}

func (self *Storage) AdminByLogin(ctx context.Context, login string) (domain.Admin, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"SELECT * FROM Admins WHERE Login = ?;")
	if err != nil {
		return domain.Admin{}, utility.CreateCustomError(self.AdminByLogin, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx, login)
	if err != nil {
		return domain.Admin{}, utility.CreateCustomError(self.AdminByLogin, err)
	}
	defer rows.Close()

	// ***

	result, err := sqlRowsToAdmin(rows)
	if err != nil {
		return domain.Admin{}, utility.CreateCustomError(self.AdminByLogin, err)
	}
	return result, nil
}

func (self *Storage) UpdateAdminPasswordByLogin(ctx context.Context, login, password string) error {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"UPDATE Admins SET Pass = ? "+
			"WHERE Login = ?;",
	)
	if err != nil {
		return utility.CreateCustomError(self.UpdateAdminPasswordByLogin, err)
	}
	defer stmt.Close()

	// ***

	_, err = stmt.ExecContext(ctx, password, login)
	if err != nil {
		return utility.CreateCustomError(self.UpdateAdminPasswordByLogin, err)
	}
	return nil
}

func (self *Storage) DeleteAdminByLogin(ctx context.Context, login string) error {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	query := "DELETE FROM Admins WHERE Login = ?;"
	stmt, err := self.db.PrepareContext(ctx, query)
	if err != nil {
		return utility.CreateCustomError(self.DeleteAdminByLogin, err)
	}
	defer stmt.Close()

	// ***

	_, err = stmt.ExecContext(ctx, login)
	if err != nil {
		return utility.CreateCustomError(self.DeleteAdminByLogin, err)
	}
	return nil
}

func (self *Storage) DeleteAdmins(ctx context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	query := "DELETE FROM Admins;"
	stmt, err := self.db.PrepareContext(ctx, query)
	if err != nil {
		return utility.CreateCustomError(self.DeleteAdmins, err)
	}
	defer stmt.Close()

	// ***

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return utility.CreateCustomError(self.DeleteAdmins, err)
	}
	return nil
}

// private
// -----------------------------------------------------------------------

func sqlRowsToRecordCount(rows *sql.Rows) (int, error) {
	var recordCount int
	if rows.Next() {
		err := rows.Scan(&recordCount)
		if err != nil {
			return 0, utility.CreateCustomError(sqlRowsToRecordCount, err)
		}
	} else {
		return 0, utility.CreateCustomError(
			sqlRowsToRecordCount, fmt.Errorf("rows are empty"))
	}
	return recordCount, nil
}

func sqlRowsToAdmins(rows *sql.Rows) ([]domain.Admin, error) {
	admins := []domain.Admin{}
	for rows.Next() {
		one := domain.Admin{}

		if err := rows.Scan(&one.Idr, &one.Login, &one.Pass); err != nil {
			return nil, utility.CreateCustomError(sqlRowsToAdmins, err)
		}
		admins = append(admins, one)
	}
	return admins, nil
}

func sqlRowsToAdmin(rows *sql.Rows) (domain.Admin, error) {
	result := domain.Admin{}
	if rows.Next() {
		err := rows.Scan(&result.Idr, &result.Login, &result.Pass)
		if err != nil {
			return domain.Admin{}, utility.CreateCustomError(sqlRowsToAdmin, err)
		}
	} else {
		return domain.Admin{}, utility.CreateCustomError(
			sqlRowsToAdmin, fmt.Errorf("rows are empty"))
	}
	return result, nil
}
