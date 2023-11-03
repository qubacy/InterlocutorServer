package impl

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/token"
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

func (s *AdminService) GetAdmins(ctx context.Context, accessToken string) domain.AdminList {

}

func (s *AdminService) PostAdmin(ctx context.Context, inp dto.PostAdminInput) (dto.PostAdminOutput, error) {

}

func (s *AdminService) DeleteAdminByIdr(ctx context.Context, accessToken string, idr int) (dto.DeleteAdminByIdrOutput, error) {

}

// private
// -----------------------------------------------------------------------
