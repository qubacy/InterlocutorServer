package n1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// to err: database is locked
func insertMultipleEntries(db *sql.DB) {
	const grCount = 100
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

// to err: database is locked
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

// ok!
func insertMultipleEntries3(db *sql.DB) {
	var wg sync.WaitGroup = sync.WaitGroup{}
	const grCount = 100
	wg.Add(grCount)

	var mx = sync.Mutex{}
	routine := func(tq string) {
		defer wg.Done()

		time.Sleep(time.Millisecond *
			time.Duration(rand.Intn(1000)))

		mx.Lock()
		_, err := db.Exec(tq)
		mx.Unlock()

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

func Exec() {
	pathToDatabaseCatalog := "./n1/storage"

	// ***

	if _, err := os.Stat(pathToDatabaseCatalog); !errors.Is(err, os.ErrNotExist) {
		err = os.RemoveAll(pathToDatabaseCatalog)
		if err != nil {
			log.Fatalf("Exec, os.RemoveAll, err %v:", err)
			return
		}
		fmt.Println("rm ok")
	}
	err := os.Mkdir(pathToDatabaseCatalog, os.ModePerm)
	if err != nil {
		log.Fatalf("Exec, os.Mkdir, err %v:", err)
		return
	}

	// ***

	queryText :=
		"CREATE TABLE IF NOT EXISTS Topics( " +
			"    Idr INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"    Lang INTEGER, " +
			"    Name TEXT " +
			"); "

	db, err := sql.Open("sqlite3", pathToDatabaseCatalog+"/topics.db")
	if err != nil {
		log.Fatalf("Exec, sql.Open, err %v:", err)
		return
	}
	defer db.Close()

	// ***

	_, err = db.Exec(queryText)
	if err != nil {
		log.Fatal("Exec, sql.Open, err:", err)
		return
	}

	// ***

	insertMultipleEntries3(db)
}
