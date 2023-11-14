package connection

import (
	"context"
	"fmt"
	"testing"
)

// experiments
// -----------------------------------------------------------------------

func Test_context_WithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	cancel()
	cancel()
	cancel()

	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
		default:
			fmt.Println("<default>")
		}
	}
}

func Test_map_getElement(t *testing.T) {
	dictionary := make(map[string]string)
	dictionary["123"] = "some name"

	key := "123"

	value, exist := dictionary[key]
	if exist {
		fmt.Printf("value `%v` found by key `%v`\n", value, key)
	} else {
		fmt.Printf("value not found by key `%v`\n", key)
	}

	// ***

	value, exist = dictionary[key]
	if exist {
		fmt.Printf("value `%v` found by key `%v`\n", value, key)
	} else {
		fmt.Printf("value not found by key `%v`\n", key)
	}
}

func Test_cutDataToSize(t *testing.T) {
	testCases := []struct {
		data string
		size int

		expected string
	}{
		{"some long string", 5, "some "},
		{"some long string", 0, ""},
		{"some long string", 500, "some long string"},
		{"some long string", 1, "s"},
		{"some long string", -1, ""},
		{"some long string", -100, ""},
	}

	for i := range testCases {
		result := cutDataToSize([]byte(testCases[i].data), testCases[i].size)
		if string(result) != testCases[i].expected {
			t.Error("Is there something wrong")
		}
	}
}

func Test_stringFromNil(t *testing.T) {
	var data []byte = nil
	fmt.Printf("'%v'\n", string(data))
}
