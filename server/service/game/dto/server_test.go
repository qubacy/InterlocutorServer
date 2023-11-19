package dto

import (
	"testing"
)

func Test_MatchedUserList_Add(t *testing.T) {
	users := MatchedUserList{}
	users.Add(MakeMatchedUser(99, "@contact"))
	if len(users) != 1 {
		t.Errorf("Something's wrong here")
	}

	// ***

	users = append(users, MakeMatchedUser(99, "@contact"))
	if len(users) != 2 {
		t.Errorf("Something's wrong here")
	}
}
