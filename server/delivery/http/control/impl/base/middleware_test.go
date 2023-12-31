package base

import (
	"flag"
	"fmt"
	token "ilserver/pkg/token/impl"
	"net/http"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	httpMock "go.nhat.io/httpmock/mock/http"
)

func TestMain(m *testing.M) {
	fmt.Println("...start test main...")

	// no need to use server options!
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
	viper.Set(key, "test_secret")
	if len(viper.GetString(key)) == 0 {
		return fmt.Errorf("Value by key '%v' is empty", key)
	}

	// solution is not very good...
	key = "_access_token"
	viper.Set(key, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3VzZXIiLCJleHAiOjQ4NTI5Nzc2MTd9.RXYKk1gtXKZvpkn1idtGXtp4SE2Qzq9aoCM1OXjAK5M")
	if len(viper.GetString(key)) == 0 {
		return fmt.Errorf("Value by key '%v' is empty", key)
	}

	return nil
}

// tests
// -----------------------------------------------------------------------

func Test_AdminIdentity_ServeHTTP_okk(t *testing.T) {
	tokenManager := newTokenManagerWithChecks(t)
	expectedResultBytes := []byte("<protected logic>")

	middleware := NewAdminIdentity(
		tokenManager,
		func(w http.ResponseWriter, r *http.Request) {
			w.Write(expectedResultBytes) // no reaction...
			//...
		})

	// ***

	response := new(httpMock.ResponseWriter)
	response.On("Write", expectedResultBytes).
		Return(len(expectedResultBytes), nil)

	middleware.ServeHTTP(response,
		httpMock.BuildRequest().
			WithHeader("Authorization",
				fmt.Sprintf("Bearer %v", viper.GetString("_access_token"))).
			Build(),
	)

	response.AssertCalled(t, "Write", expectedResultBytes)
}

func Test_AdminIdentity_ServeHTTP_err(t *testing.T) {
	tokenManager := newTokenManagerWithChecks(t)
	expectedResultBytes := []byte("<protected logic>")

	middleware := NewAdminIdentity(
		tokenManager,
		func(w http.ResponseWriter, r *http.Request) {
			w.Write(expectedResultBytes) // no reaction...
			//...
		})

	// ***

	response := new(httpMock.ResponseWriter)
	response.On("Header", mock.Anything).Return(http.Header{})
	response.On("WriteHeader", http.StatusBadRequest)
	response.On("Write", mock.Anything).Return(0, nil)

	middleware.ServeHTTP(response,
		httpMock.BuildRequest().
			WithHeader("Authorization",
				fmt.Sprintf("Bearer %v", "123")).
			Build(),
	)

	response.AssertCalled(t, "Header", mock.Anything)
	response.AssertCalled(t, "WriteHeader", http.StatusBadRequest)
	response.AssertCalled(t, "Write", mock.Anything)
}

// experiments
// -----------------------------------------------------------------------

func Test_mock_http_Client(t *testing.T) {
	request := httpMock.BuildRequest().
		WithHeader("Authorization", "Bearer 1234567890").
		WithURI("/control/api/topics").
		Build()

	fmt.Println("Header:", request.Header)
	fmt.Println("Uri:", request.RequestURI)
}

// private
// -----------------------------------------------------------------------

func newTokenManagerWithChecks(t *testing.T) *token.Manager {
	tokenManager, err := token.NewManager(viper.GetString("control_server.token.secret"))
	if err != nil {
		t.Fatalf("Failed to get a new manager. Err: %v", err)
		return nil
	}
	return tokenManager
}
