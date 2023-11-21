package sqlite

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/utility"
	storage "ilserver/storage/control"
	"math/rand"
	"strconv"
	"testing"
)

func Test_InsertAdmin(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Fatal()
		return
	}

	// ***

	valuesForInsert := generateFakeAdmins()
	insertAdminsWithChecks(t, storage, valuesForInsert)

	// ***

	num := recordCountInTableWithChecks(t, storage, "Admins")
	if num != len(valuesForInsert)+1 && num != len(valuesForInsert) { // plus default entity.
		t.Errorf("Record count in the table is not equal %v", len(valuesForInsert))
		return
	}

	// ***

	deleteAdminsWithChecks(t, storage)
}

func Test_HasAdminByLogin(t *testing.T) {
	storage := instanceWithChecks(t)

	newAdmin := generateFakeAdmin("new_admin_with_unknown_login")
	valuesForInsert := append(generateFakeAdmins(), newAdmin)

	insertAdminsWithChecks(t, storage, valuesForInsert)

	// ***

	has, err := storage.HasAdminByLogin(context.Background(), newAdmin.Login)
	if err != nil {
		t.Errorf("Is there something wrong. Err: %v", err)
		return
	}
	if !has {
		t.Errorf("Administrator account not found")
		return
	}

	// ***

	deleteAdminsWithChecks(t, storage)
}

func Test_AllAdmins(t *testing.T) {
	// TODO: or enough?
}

// private help functions.
// -----------------------------------------------------------------------

func insertAdminsWithChecks(t *testing.T, storage storage.Storage, admins domain.AdminList) {
	for i := range admins {
		_, err := storage.InsertAdmin(context.Background(), admins[i])
		if err != nil {
			t.Errorf("Insert admin failed. Err: %v", err)
			return
		}
	}
}

func deleteAdminsWithChecks(t *testing.T, storage storage.Storage) {
	err := storage.DeleteAdmins(context.Background())
	if err != nil {
		t.Errorf("Failed to delete admins. Err: %v", err)
		return
	}
}

// ***

func generateFakeAdmin(login string) domain.Admin {
	return domain.Admin{
		Login: login,
		Pass:  utility.RandomString(100),
	}
}

func generateFakeAdmins() domain.AdminList {
	count := rand.Intn(10) + 10
	entities := make(domain.AdminList, count)
	for i := range entities {
		entities[i] = domain.Admin{
			Login: utility.RandomString(99) + strconv.Itoa(i),
			Pass:  utility.RandomString(100),
		}
	}
	return entities
}
