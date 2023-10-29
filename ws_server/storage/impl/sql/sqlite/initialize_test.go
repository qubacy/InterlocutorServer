package sqlite

import (
	"fmt"
	"os"
	"testing"
)

// run all tests from file 'storage_test.go'

func Test_createDatabaseFile(t *testing.T) {
	// TODO:
}

func Test_openDatabaseFile(t *testing.T) {
	db, err := openDatabaseFile()
	if err != nil {
		fmt.Println("err:", err)
		t.Fail()
		return
	}
	defer db.Close()

	// ***

	stats := db.Stats()
	fmt.Println("OpenConnections:", stats.OpenConnections)
	fmt.Println("MaxOpenConnections:", stats.MaxOpenConnections)

	fmt.Println("Idle:", stats.Idle)
	fmt.Println("MaxIdleClosed:", stats.MaxIdleClosed)
	fmt.Println("MaxIdleTimeClosed:", stats.MaxIdleTimeClosed)
	fmt.Println("MaxLifetimeClosed:", stats.MaxLifetimeClosed)
	//...
}

func Test_exists(t *testing.T) {
	basicPath := databaseDirectory
	testCases := []struct {
		path string
		want bool
	}{
		{path: basicPath, want: true},
		{path: basicPath + "storage.go", want: true},
		{path: basicPath + "unknown.go", want: false},
		//...
	}

	// ***

	for i := range testCases {
		actual := exists(testCases[i].path)
		expected := testCases[i].want

		if expected != actual {
			t.Errorf("Result was incorrect, test path: %v", testCases[i].path)
		}
	}
}

// experiments
// -----------------------------------------------------------------------

func Test_Stat(t *testing.T) {
	info, err := os.Stat(databaseDirectory)
	if err != nil {
		t.Error()
		return
	}

	if info.IsDir() != true {
		t.Fail()
		return
	}

	fmt.Println("db dir mode:", info.Mode())
	fmt.Println("db dir mod time:", info.ModTime())
	fmt.Println("db dir name:", info.Name())
	fmt.Println("db dir size:", info.Size())
	fmt.Println("db dir sys:", info.Sys())
	//...
}
