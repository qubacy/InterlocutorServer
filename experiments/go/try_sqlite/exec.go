package try_sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// raw query
func insertMultipleEntries(db *sql.DB) {
	const grCount = 2
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(grCount)

	routine := func(tq string) {
		defer wg.Done()

		time.Sleep(time.Millisecond *
			time.Duration(rand.Intn(1000)))

		_, err := db.Exec(tq)
		if err != nil {
			log.Fatalf("insertMultipleEntries, sql.Open, err: %v", err)
			return
		}
	}

	// ***

	for i := 0; i < grCount; i++ {
		tq := fmt.Sprintf(
			"INSERT INTO Topics (Lang, Name) "+
				"VALUES (0, 'Topic%d');", i)

		go routine(tq)
	}

	wg.Wait()
}

// transaction
func insertMultipleEntries1(db *sql.DB) {
	const grCount = 1000
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(grCount)

	// ***

	routine := func(tq string, ctx context.Context) {
		defer wg.Done()

		time.Sleep(time.Millisecond *
			time.Duration(rand.Intn(1000)))

		// ***

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("insertMultipleEntries1, db.Begin, err: %v", err)
			return
		}

		_, err = tx.Exec(tq)
		if err != nil {
			tx.Rollback()

			log.Fatalf("insertMultipleEntries1, tx.Exec, err: %v", err)
			return
		}
		tx.Commit()
	}

	// ***

	for i := 0; i < grCount; i++ {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

		tq := fmt.Sprintf(
			"INSERT INTO Topics (Lang, Name) "+
				"VALUES (0, 'Topic%d');", i)

		go func() {
			defer cancel()
			routine(tq, ctx)
		}()
	}

	wg.Wait()
}

func Exec() {
	tq :=
		"CREATE TABLE IF NOT EXISTS Topics( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Lang INTEGER, " +
			"    Name TEXT " +
			"); "

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatalf("Exec, sql.Open, err %v:", err)
		return
	}

	// ***

	_, err = db.Exec(tq)
	if err != nil {
		log.Fatal("Exec, sql.Open, err:", err)
		return
	}

	// ***

	insertMultipleEntries1(db)
	db.Close()
}
