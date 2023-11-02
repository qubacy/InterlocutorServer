package token

import (
	"bytes"
	"encoding/base64"
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

func Test_base64_StdEncoding_EncodeToString(t *testing.T) {
	result := base64.StdEncoding.EncodeToString([]byte("one 🐘 and three 🐋"))
	_, err := fmt.Println(result)
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
}

func Test_base64_StdEncoding_EncodeToString_v1(t *testing.T) {
	result := base64.StdEncoding.EncodeToString([]byte("deja vu"))
	_, err := fmt.Println(result)
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
}

// ***

func Test_base64_StdEncoding_DecodeString(t *testing.T) {
	str := "YmFzZTY0LlN0ZEVuY29kaW5nLkRlY29kZVN0cmluZw=="
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
		return
	}

	_, err = fmt.Println(string(data))
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
}

func Test_base64_StdEncoding_DecodeString_v1(t *testing.T) {
	str := "ZGVqYSB2dQ=="
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	_, err = fmt.Println(string(data))
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
}

// ***

func Test_base64_NewEncoder(t *testing.T) {
	buffer := bytes.NewBufferString("")
	encoding := base64.StdEncoding.WithPadding(base64.StdPadding)
	encoder := base64.NewEncoder(encoding, buffer)
	defer encoder.Close()

	// ***

	n, err := encoder.Write([]byte("deja vu")) // TODO: does not work!?
	fmt.Printf("%v bytes written to buffer\n", n)
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
	if n == 0 {
		t.Errorf("Is there something wrong")
	}

	// ***

	fmt.Println("Buffer string:", buffer.String())
	fmt.Println("Buffer bytes:", string(buffer.Bytes()))

	fmt.Println("Cap buffer:", buffer.Cap())
	fmt.Println("Len buffer:", buffer.Len())

	if string(buffer.Bytes()) == "ZGVqYSB2dQ==" { // !=
		t.Errorf("Is there something wrong")
	}

	// ***

	rawBuffer := make([]byte, 512)
	n, err = buffer.Read(rawBuffer)
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
	fmt.Printf("%v bytes written to buffer\n", n)
	if n == 0 {
		t.Errorf("Is there something wrong")
	}
	fmt.Println("Raw buffer as string:", string(rawBuffer))

	fmt.Println("Cap buffer:", buffer.Cap())
	fmt.Println("Len buffer:", buffer.Len())
}

func Test_bytes_NewBufferString(t *testing.T) {
	buffer := bytes.NewBufferString("")
	_, err := fmt.Println("Cap buffer:", buffer.Cap())
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}

	// ***

	n, err := buffer.Write([]byte("12345678901234567890123456789012345"))
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
	if n != 35 {
		t.Errorf("Is there something wrong")
	}

	// ***

	_, err = fmt.Println("Cap buffer:", buffer.Cap())
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
	}
}
