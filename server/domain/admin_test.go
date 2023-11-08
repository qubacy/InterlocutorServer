package domain

import "testing"

func Test_AdminList_SomeEq(t *testing.T) {
	admins := AdminList{
		{Idr: 1, Login: "test", Pass: "test"},
	}
	if admins.Eq(nil) {
		t.Error("Something wrong")
		return
	}

	// ***

	otherAdmins := admins
	admins = nil

	if admins.Eq(otherAdmins) {
		t.Error("Something wrong")
		return
	}

	// ***

	otherAdmins = nil
	admins = nil

	if admins.Eq(otherAdmins) {
		t.Error("Something wrong")
		return
	}
}
