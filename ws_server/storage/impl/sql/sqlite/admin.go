package sqlite

import "fmt"

func (r *Storage) HasAdminByLogin(login string) (error, bool) {
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

func (r *Storage) UpdateAdminPass(login, newPass string) error {
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

func (r *Storage) HasAdminWithLoginAndPass(login, pass string) (error, bool) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare(
		"SELECT count(*) AS RecordCount " +
			"FROM [Admins] " +
			"WHERE [Admins].[Login] = ? " +
			"AND [Admins].[Pass] = ?;")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err), false
	}
	defer stmt.Close() // ignore err!

	// ***

	rows, err := stmt.Query(login, pass)
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
