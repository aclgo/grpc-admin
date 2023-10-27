package usecase

import (
	"context"

	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/models"
	"github.com/aclgo/grpc-admin/pkg/logger"
)

type AdminService struct {
	adminRepo admin.AdminRepo
	logger    logger.Logger
}

func NewAdminService(adminRepo admin.AdminRepo, logger logger.Logger) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		logger:    logger,
	}
}

func (a *AdminService) Create(ctx context.Context, params *admin.ParamsCreateAdmin) (*models.ParamsUser, error) {

	created, err := a.adminRepo.Create(ctx, &models.ParamsCreateAdmin{
		Name:     params.Name,
		Lastname: params.Lastname,
		Password: params.Password,
		Email:    params.Email,
		Role:     params.Role,
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (a *AdminService) SearchUsers(ctx context.Context, params *admin.ParamsSearchUsers) (*models.DataSearchedUser, error) {
	searched, err := a.adminRepo.Search(ctx,
		&admin.ParamsSearchUsers{
			Query: params.Query,
			Role:  params.Role,
			Page:  params.Page,
			Pagination: admin.Pagination{
				OffSet: params.Pagination.OffSet,
				Limit:  params.Pagination.Limit,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return searched, nil
}
