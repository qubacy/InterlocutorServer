package utility

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGetFunctionName(t *testing.T) {
	testCases := []struct {
		param interface{}
		want  string
	}{
		{111, ""},
		{111.111, ""},
		{"string", ""},
		{strings.Compare, "strings.Compare"},
		{strings.Contains, "strings.Contains"},
		{strings.Count, "strings.Count"},
		{GetFunctionName, "utility.GetFunctionName"},
		{IsFunction, "utility.IsFunction"},
		//...
	}

	for i := range testCases {
		actual := GetFunctionName(testCases[i].param)
		if actual != testCases[i].want {
			t.Fail()
			return
		}
	}
}

func TestIsFunction(t *testing.T) {
	testCases := []struct {
		param interface{}
		want  bool
	}{
		{111, false},
		{111.111, false},
		{"string", false},
		{strings.Compare, true},
		{strings.Contains, true},
		{strings.Count, true},
		//...
	}

	// ***

	for i := range testCases {
		actual := IsFunction(testCases[i].param)
		if actual != testCases[i].want {
			t.Fail()
			return
		}
	}
}

func Test_CreateCustomError(t *testing.T) {
	err := CreateCustomError(GetFunctionName, fmt.Errorf("<some library error>"))
	err = CreateCustomError(GetFunctionName, err)
	err = CreateCustomError(strconv.Atoi, err)
	err = CreateCustomError(IsFunction, err)
	//...
	fmt.Println(err.Error())
}

// benchmarks
// -----------------------------------------------------------------------

func BenchmarkGetFunctionName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if x := GetFunctionName(strings.Contains); x != "strings.Contains" {
			b.Fatalf("Unexpected result: %v", x)
		}
	}
}

func BenchmarkIsFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if x := IsFunction(strings.Contains); x != true {
			b.Fatalf("Unexpected result: %v", x)
		}
	}
}

// experiments
// -----------------------------------------------------------------------

func TestSplit(t *testing.T) {
	testCases := []struct {
		param string
		sep   string
		want  []string
	}{
		{"", ",", []string{""}},
		{"ilserver/utility.GetFunctionName", "/", []string{"ilserver", "utility.GetFunctionName"}},
		{"123 123 123", " ", []string{"123", "123", "123"}},
		{"123   123   123", " ", []string{"123", "", "", "123", "", "", "123"}},
	}

	// ***

	for i := range testCases {
		actual := strings.Split(
			testCases[i].param,
			testCases[i].sep,
		)

		if len(actual) != len(testCases[i].want) {
			t.Fail()
			return
		}
		if reflect.TypeOf(actual) != reflect.TypeOf(testCases[i].want) {
			t.Fail()
			return
		}
		if !reflect.DeepEqual(actual, testCases[i].want) {
			t.Fail()
			return
		}
	}
}

func TestFunctionType(t *testing.T) {
	var fun interface{} = GetFunctionName

	funType := reflect.TypeOf(fun)
	fmt.Println("fun type:", funType)

	funTypeName := funType.String()
	fmt.Println("fun type name:", funTypeName)

	if !strings.Contains(funType.String(), "func") {
		t.Fail()
		return
	}
}

func TestFuncName(t *testing.T) {
	pc, file, line, ok := runtime.Caller(0)
	if !ok {
		t.Fail()
		return
	}

	fmt.Printf("pc: %v\n", pc)
	fmt.Printf("file: %v\n", file)
	fmt.Printf("line: %v\n", line)

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		t.Fail()
		return
	}

	fmt.Println()
	fmt.Printf("full function name: %v\n", fn.Name())

	shortFuncName := strings.Split(fn.Name(), ".")[1]
	fmt.Printf("short function name: %v\n", shortFuncName)

	if shortFuncName != "TestFuncName" {
		t.Fail()
		return
	}
}

func TestFuncNameV1(t *testing.T) {
	functionName := runtime.FuncForPC(
		reflect.ValueOf(
			TestFuncNameV1).Pointer()).Name()

	fmt.Println("function name:", functionName)

	shortFuncName := strings.Split(functionName, ".")[1]
	fmt.Printf("short function name: %v\n", shortFuncName)

	if shortFuncName != "TestFuncNameV1" {
		t.Fail()
		return
	}
}

func TestWrapError(t *testing.T) {
	err := errors.New("0")
	err = fmt.Errorf("1 %w", err)
	err = fmt.Errorf("2 %w", err)

	// ***

	if err.Error() != "2 1 0" {
		t.Fail()
		return
	}
	err = errors.Unwrap(err)
	if err.Error() != "1 0" {
		t.Fail()
		return
	}
	err = errors.Unwrap(err)
	if err.Error() != "0" {
		t.Fail()
		return
	}
	err = errors.Unwrap(err)
	if err != nil {
		t.Fail()
		return
	}
}

func TestExecutable(t *testing.T) {
	path, err := os.Executable()
	if err != nil {
		fmt.Println("err:", err.Error())
		t.Fail()
		return
	}

	fmt.Printf("path: %v\n", path)
}

func TestGetwd(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("err:", err.Error())
		t.Fail()
		return
	}

	fmt.Printf("path: %v\n", path)
}
