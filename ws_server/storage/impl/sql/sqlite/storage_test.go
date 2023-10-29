package sqlite

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// executed even if exactly one test is started!
func TestMain(m *testing.M) {
	fmt.Println("...start test main...")

	if err := setUp(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return
	}

	// ***

	flag.Parse() // ?
	code := m.Run()

	// ***

	if err := tearDown(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return
	}

	fmt.Println("...finish test main...")
	os.Exit(code)
}

// set up / tear down
// -----------------------------------------------------------------------

func setUp() error {
	err := setUpDatabaseDirectory()
	if err != nil {
		return err
	}

	err = setUpViper()
	if err != nil {
		return err
	}

	return nil
}

func setUpDatabaseDirectory() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	pathParts := strings.Split(path, "\\")
	pathParts = pathParts[0 : len(pathParts)-3]
	path = strings.Join(pathParts, "\\")

	databaseDirectory = path + "\\" // !
	return nil
}

func setUpViper() error {
	key := "storage.sql.sqlite.file"
	viper.Set(key, "storage.db")
	if len(viper.GetString(key)) == 0 {
		return fmt.Errorf("Value by key '%v' is empty", key)
	}

	key = "storage.default_admin_entry.login"
	viper.Set(key, "admin")
	if len(viper.GetString(key)) == 0 {
		return fmt.Errorf("Value by key '%v' is empty", key)
	}

	key = "storage.default_admin_entry.pass"
	viper.Set(key, "admin")
	if len(viper.GetString(key)) == 0 {
		return fmt.Errorf("Value by key '%v' is empty", key)
	}

	key = "storage.sql.initialization_timeout"
	viper.Set(key, 1*time.Second)
	if viper.GetDuration(key) == 0 {
		return fmt.Errorf("Value by key '%v' is zero", key)
	}

	//...

	return nil
}

func tearDown() error {
	Free() // close all!
	err := removeDatabaseFile()
	if err != nil {
		return err
	}
	return nil
}

// tests
// -----------------------------------------------------------------------

func Test_Instance(t *testing.T) {
	_, err := Instance()
	if err != nil {
		fmt.Println("err:", err)
		t.Fail()
		return
	}
}

func Test_RecordCountInTable_okk(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Errorf("failed to get storage. Err: %v", err)
		return
	}

	// ***

	// attention empty context!
	num, err := storage.RecordCountInTable(context.Background(), "Admins")
	if err != nil {
		t.Errorf("request failed. Err: %v", err)
		return
	}
	if num != 1 {
		t.Errorf("row count not equal 1")
		return
	}

	// ***

	num, err = storage.RecordCountInTable(context.Background(), "Topics")
	if err != nil {
		t.Errorf("request failed. Err: %v", err)
		return
	}
	if num != 0 {
		t.Errorf("row count not equal 0")
		return
	}
}

func Test_RecordCountInTable_err(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Errorf("failed to get storage. Err: %v", err)
		return
	}

	// ***

	_, err = storage.RecordCountInTable(context.Background(), "UnknownTable")
	if err == nil {
		t.Errorf("table does not exist, but query completed")
		return
	}
}

// experiments
// -----------------------------------------------------------------------

func Test_Getwd(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println("path:", path)
}
