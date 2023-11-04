package impl

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
		return dto.MakeGetAdminsError(
			control.ErrAccessTokenIsEmpty(),
		), nil
	}

	// *** just check that the token is valid...

	// TODO: token

	// ***

	admins, err := s.storage.AllAdmins(ctx)
	if err != nil {
		// ---> 500
		return dto.GetAdminsOutput{},
			utility.CreateCustomError(s.GetAdmins, err)
	}

	return dto.MakeGetAdminsSuccess(
		extractAdminLogins(admins),
	), nil
}

func (s *AdminService) PostAdmin(ctx context.Context, inp dto.PostAdminInput) (dto.PostAdminOutput, error) {

}

func (s *AdminService) DeleteAdminByIdr(ctx context.Context, accessToken string, idr int) (dto.DeleteAdminByIdrOutput, error) {

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
