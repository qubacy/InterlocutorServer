package token

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

	fmt.Println("...finish test main...")
	os.Exit(code)
}

func setUp() error {
	err := setUpViper()
	if err != nil {
		return err
	}

	return nil
}

func setUpViper() error {
	key := "control_server.token.secret"
	viper.Set(key, "secret")
	if len(viper.GetString(key)) == 0 {
		return fmt.Errorf("Value by key '%v' is empty", key)
	}
	return nil
}

// -----------------------------------------------------------------------

func Test_NewToken(t *testing.T) {
	err, tokenString := NewToken("test")
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}

	// ***

	fmt.Println(tokenString)
	tokenParts := strings.Split(tokenString, ".")
	if len(tokenParts) != 3 {
		t.Errorf("Is there something wrong. Err: %v", err)
	}

	// ***

}

// experiments
// -----------------------------------------------------------------------

func Test_base64(t *testing.T) {
	// TODO:

}
