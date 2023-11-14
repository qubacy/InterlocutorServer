package base

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"
	storage "ilserver/storage/control"
)

type AdminService struct {
	storage storage.Storage
}

func NewAdminService(storage storage.Storage) *AdminService {
	return &AdminService{
		storage: storage,
	}
}

// -----------------------------------------------------------------------

func (s *AdminService) GetAdmins(ctx context.Context) (dto.GetAdminsOutput, error) {
	admins, err := s.storage.AllAdmins(ctx)
	if err != nil {
		// ---> 400
		return dto.MakeGetAdminsOutputEmpty(),
			utility.CreateCustomError(s.GetAdmins, err)
	}

	// ---> 200
	return dto.MakeGetAdminsOutputSuccess(admins), nil
}

func (s *AdminService) PostAdmin(ctx context.Context, inp dto.PostAdminInput) (dto.PostAdminOutput, error) {
	if !isValidPostAdminInput(inp) {
		// ---> 400
		return dto.MakePostAdminOutputEmpty(),
			utility.CreateCustomError(s.PostAdmin,
				control.ErrInputDtoIsInvalid())
	}

	// ***

	idr, err := s.storage.InsertAdmin(ctx, inp.ToDomain())
	if err != nil {
		// ---> 400
		return dto.MakePostAdminOutputEmpty(),
			utility.CreateCustomError(s.PostAdmin, err)
	}

	// ***

	// ---> 200
	return dto.MakePostAdminOutputSuccess(idr), nil
}

func (s *AdminService) DeleteAdminByIdr(ctx context.Context, idr int) (dto.DeleteAdminByIdrOutput, error) {
	err := s.storage.DeleteAdmin(ctx, idr)
	if err != nil {
		// ---> 400
		return dto.MakeDeleteAdminByIdrOutputEmpty(),
			utility.CreateCustomError(s.DeleteAdminByIdr, err)
	}

	// ---> 200
	return dto.MakeDeleteAdminByIdrOutputSuccess(), nil
}

// private
// -----------------------------------------------------------------------

func extractAdminLogins(admins domain.AdminList) []string {
	result := make([]string, 0)
	for i := range admins {
		result = append(result, admins[i].Login)
	}
	return result
}

func isValidPostAdminInput(inp dto.PostAdminInput) bool {
	if len(inp.Login) <= 3 {
		return false
	}
	if len(inp.Pass) <= 3 {
		return false
	}
	return true
}
