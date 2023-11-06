package base

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"ilserver/service/control"
	"ilserver/service/control/dto"
	"ilserver/storage"
)

type AdminService struct {
	storage      storage.Storage
	tokenManager token.Manager
}

func NewAdminService(storage storage.Storage, tokenManager token.Manager) *AdminService {
	return &AdminService{
		storage:      storage,
		tokenManager: tokenManager,
	}
}

// -----------------------------------------------------------------------

func (s *AdminService) GetAdmins(ctx context.Context, accessToken string) (dto.GetAdminsOutput, error) {
	if len(accessToken) == 0 {
		// ---> 400
		return dto.MakeGetAdminsEmpty(),
			utility.CreateCustomError(s.GetAdmins,
				control.ErrAccessTokenIsEmpty())
	}

	// *** just check that the token is valid...

	err := s.tokenManager.Validate(accessToken)
	if err != nil {
		// ---> 400
		return dto.MakeGetAdminsEmpty(),
			utility.CreateCustomError(s.GetAdmins, err)
	}

	// ***

	admins, err := s.storage.AllAdmins(ctx)
	if err != nil {
		// ---> 400
		return dto.MakeGetAdminsEmpty(),
			utility.CreateCustomError(s.GetAdmins, err)
	}

	// ---> 200
	return dto.MakeGetAdminsSuccess(
		extractAdminLogins(admins),
	), nil
}

func (s *AdminService) PostAdmin(ctx context.Context, inp dto.PostAdminInput) (dto.PostAdminOutput, error) {
	if !isValidPostAdminInput(inp) {
		// ---> 400
		return dto.MakePostAdminEmpty(),
			utility.CreateCustomError(s.PostAdmin,
				control.ErrInputDtoIsInvalid())
	}

	// ***

	err := s.tokenManager.Validate(inp.AccessToken)
	if err != nil {
		// ---> 400
		return dto.MakePostAdminEmpty(),
			utility.CreateCustomError(s.PostAdmin, err)
	}

	// ***

	idr, err := s.storage.InsertAdmin(ctx, dto.PostAdminInputToDomain(inp))
	if err != nil {
		// ---> 400
		return dto.MakePostAdminEmpty(),
			utility.CreateCustomError(s.PostAdmin, err)
	}

	// ***

	// ---> 200
	return dto.MakePostAdminSuccess(idr), nil
}

func (s *AdminService) DeleteAdminByIdr(ctx context.Context, accessToken string, idr int) (dto.DeleteAdminByIdrOutput, error) {
	if len(accessToken) == 0 {
		// ---> 400
		return dto.MakeDeleteAdminByIdrEmpty(),
			utility.CreateCustomError(s.DeleteAdminByIdr,
				control.ErrAccessTokenIsEmpty())
	}

	err := s.tokenManager.Validate(accessToken)
	if err != nil {
		// ---> 400
		return dto.MakeDeleteAdminByIdrEmpty(),
			utility.CreateCustomError(s.DeleteAdminByIdr, err)
	}

	err = s.storage.DeleteAdmin(ctx, idr)
	if err != nil {
		// ---> 400
		return dto.MakeDeleteAdminByIdrEmpty(),
			utility.CreateCustomError(s.DeleteAdminByIdr, err)
	}

	// ---> 200
	return dto.MakeDeleteAdminByIdrSuccess(), nil
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
	if len(inp.AccessToken) == 0 {
		return false
	}
	if len(inp.Login) <= 3 {
		return false
	}
	if len(inp.Pass) <= 3 {
		return false
	}
	return true
}
