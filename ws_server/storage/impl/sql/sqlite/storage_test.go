package sqlite

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

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

	//...

	return nil
}

func tearDown() error {
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
	defer Free()
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
