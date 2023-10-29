package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"ilserver/domain"
	"ilserver/utility"
	"strconv"
)

// insert
// -----------------------------------------------------------------------

func (self *Storage) InsertTopic(ctx context.Context, topic domain.Topic) (int64, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"INSERT INTO [Topics] (Lang, Name) "+
			"VALUES (?, ?);",
	)
	if err != nil {
		return 0, utility.CreateCustomError(self.InsertTopic, err)
	}
	defer stmt.Close()

	// ***

	res, err := stmt.ExecContext(ctx, topic.Lang, topic.Name)
	if err != nil {
		return 0, utility.CreateCustomError(self.InsertTopic, err)
	}

	// ***

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return 0, utility.CreateCustomError(self.InsertTopic, err)
	}

	return lastInsertedId, nil
}

func (self *Storage) InsertTopics(ctx context.Context, topics []domain.Topic) error {
	if len(topics) == 0 {
		return nil
	}

	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	tq := "INSERT INTO [Topics] (Lang, Name) " +
		"VALUES "
	tq += topicsToSqlInsertValues(topics)
	tq += ";"

	// ***

	_, err := self.db.ExecContext(ctx, tq)
	if err != nil {
		return utility.CreateCustomError(self.AllTopics, err)
	}

	// ***

	return nil
}

// select
// -----------------------------------------------------------------------

func (self *Storage) AllTopics(ctx context.Context) ([]domain.Topic, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx, "SELECT * FROM Topics;")
	if err != nil {
		return nil, utility.CreateCustomError(self.AllTopics, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, utility.CreateCustomError(self.AllTopics, err)
	}
	defer rows.Close()

	// ***

	topics, err := sqlRowsToTopics(rows)
	if err != nil {
		return nil, utility.CreateCustomError(self.AllTopics, err)
	}
	return topics, nil
}

func (self *Storage) Topic(ctx context.Context, idr int) (domain.Topic, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"SELECT * FROM Topics "+
			"WHERE Idr = ?;",
	)
	if err != nil {
		return domain.Topic{}, utility.CreateCustomError(self.Topic, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx, idr)
	if err != nil {
		return domain.Topic{}, utility.CreateCustomError(self.Topic, err)
	}
	defer rows.Close()

	// ***

	result, err := sqlRowsToTopic(rows)
	if err != nil {
		return domain.Topic{}, utility.CreateCustomError(self.Topic, err)
	}
	return result, nil
}

func (self *Storage) RandomTopic(ctx context.Context, lang int) (domain.Topic, error) {
	self.mx.Lock()
	defer self.mx.Unlock()

	// ***

	stmt, err := self.db.PrepareContext(ctx,
		"SELECT * FROM Topics "+
			"WHERE Lang = ? "+
			"ORDER BY random() "+
			"LIMIT 1;",
	)
	if err != nil {
		return domain.Topic{}, utility.CreateCustomError(self.RandomTopic, err)
	}
	defer stmt.Close()

	// ***

	rows, err := stmt.QueryContext(ctx, lang)
	if err != nil {
		return domain.Topic{}, utility.CreateCustomError(self.RandomTopic, err)
	}
	defer rows.Close()

	// ***

	result, err := sqlRowsToTopic(rows)
	if err != nil {
		return domain.Topic{}, utility.CreateCustomError(self.RandomTopic, err)
	}
	return result, nil
}

func (self *Storage) DeleteTopic(ctx context.Context, idr int) error {
	// TODO:
	return nil
}

// private
// -----------------------------------------------------------------------

func sqlRowsToTopic(rows *sql.Rows) (domain.Topic, error) {
	result := domain.Topic{}
	if rows.Next() {
		err := rows.Scan(&result.Idr, &result.Lang, &result.Name)
		if err != nil {
			return domain.Topic{}, utility.CreateCustomError(sqlRowsToTopic, err)
		}
	} else {
		return domain.Topic{}, utility.CreateCustomError(
			sqlRowsToTopic, fmt.Errorf("rows are empty"))
	}
	return result, nil
}

func sqlRowsToTopics(rows *sql.Rows) ([]domain.Topic, error) {
	topics := []domain.Topic{}
	for rows.Next() {
		one := domain.Topic{}

		if err := rows.Scan(&one.Idr, &one.Lang, &one.Name); err != nil {
			return nil, utility.CreateCustomError(sqlRowsToTopics, err)
		}
		topics = append(topics, one)
	}
	return topics, nil
}

func topicsToSqlInsertValues(topics []domain.Topic) string {
	result := ""
	for i := range topics {
		result += "("
		result += strconv.Itoa(topics[i].Lang) + ", "
		result += ("'" + topics[i].Name + "'")
		result += ")"

		if i != len(topics)-1 {
			result += ", "
		}
	}
	return result
}
