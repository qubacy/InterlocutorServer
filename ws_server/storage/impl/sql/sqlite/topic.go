package sqlite

import (
	"fmt"
	"ilserver/domain"
	"strconv"
)

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
