package sqlite

import (
	"database/sql"
	"fmt"
	"ilserver/domain"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// -----------------------------------------------------------------------

type Storage struct {
	mx sync.Mutex
	db *sql.DB
}

var once = sync.Once{}
var instance *Storage = nil
var initializationError error // user friendly error

// constructor
// -----------------------------------------------------------------------

func Instance() (*Storage, error) {
	once.Do(func() {
		initializationError =
			initialize()
	})
	return instance, initializationError
}
func Free() {
	// TODO: что делать с ошибкой?
	instance.db.Close()
}

func newStorage(db *sql.DB) *Storage {
	return &Storage{
		mx: sync.Mutex{},
		db: db,
	}
}

// admins
// -----------------------------------------------------------------------

func (r *Storage) RecordCountInTable(name string) (error, int) {
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

// topics
// -----------------------------------------------------------------------

func (r *Storage) SelectRandomOneTopic(lang int) (error, domain.Topic) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare(
		"SELECT * FROM Topics " +
			"WHERE Lang = ? " +
			"ORDER BY random() " +
			"LIMIT 1;")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err),
			domain.Topic{}
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.Query(lang)
	if err != nil {
		return fmt.Errorf("execute query failed %v", err),
			domain.Topic{}
	}
	defer rows.Close()

	// ***

	tc := domain.Topic{}
	if rows.Next() {
		err = rows.Scan(&tc.Idr, &tc.Lang, &tc.Name)
		if err != nil {
			return fmt.Errorf("scan next row with err %v", err),
				domain.Topic{}
		}
	} else {
		return fmt.Errorf("rows count is zero"),
			domain.Topic{}
	}
	return nil, tc
}

func (r *Storage) SelectTopics() (error, []domain.Topic) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	stmt, err := r.db.Prepare(
		"SELECT * FROM Topics;")
	if err != nil {
		return fmt.Errorf("prepare query failed %v", err),
			nil
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.Query()
	if err != nil {
		return fmt.Errorf("execute query failed %v", err),
			nil
	}
	defer rows.Close()

	// ***

	topics := []domain.Topic{}
	for rows.Next() {
		one := domain.Topic{}
		if err := rows.Scan(&one.Idr, &one.Lang, &one.Name); err != nil {
			return fmt.Errorf("scan next row with err %v", err),
				[]domain.Topic{}
		}

		topics = append(topics, one)
	}
	return nil, topics
}

func (r *Storage) InsertTopic(topic domain.Topic) (error, int64) {
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

func (r *Storage) InsertTopics(topics []domain.Topic) error {
	if len(topics) == 0 {
		return nil
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	// ***

	tq := "INSERT INTO [Topics] (Lang, Name) " +
		"VALUES "
	for i := range topics {
		tq += "("
		tq += strconv.Itoa(topics[i].Lang) + ", "
		tq += ("'" + topics[i].Name + "'")
		tq += ")"

		if i != len(topics)-1 {
			tq += ", "
		}
	}
	tq += ";"

	// ***

	_, err := r.db.Exec(tq)
	if err != nil {
		return fmt.Errorf("execute query failed %v", err)
	}

	// ***

	return nil
}
