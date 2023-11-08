package domain

import (
	"testing"
)

func Test_IsEqual(t *testing.T) {
	testCases := []struct {
		lhs  Topic
		rhs  Topic
		want bool
	}{
		{Topic{Idr: 1, Lang: 1, Name: "1"}, Topic{Idr: 1, Lang: 1, Name: "1"}, true},
		{Topic{Idr: 5, Lang: 1, Name: "1"}, Topic{Idr: 10, Lang: 1, Name: "1"}, true},
		{Topic{Idr: 1, Lang: 2, Name: "1"}, Topic{Idr: 1, Lang: 1, Name: "1"}, false},
		{Topic{Idr: 1, Lang: 1, Name: "123"}, Topic{Idr: 1, Lang: 1, Name: "111"}, false},
		{Topic{Idr: 1, Lang: 2, Name: "123"}, Topic{Idr: 1, Lang: 1, Name: "111"}, false},
		//...
	}

	// ***

	for i := range testCases {
		got := testCases[i].lhs.Eq(testCases[i].rhs)
		if got != testCases[i].want {
			t.Errorf("The left topic %v is not equal to the right topic %v",
				testCases[i].lhs, testCases[i].rhs)
		}
	}
}

//...
