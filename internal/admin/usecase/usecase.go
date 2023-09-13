package usecase

import (
	"context"

	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/models"
)

type AdminService struct {
	adminRepo admin.AdminRepo
}

func NewAdminService(adminRepo admin.AdminRepo) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

type ParamsCreate struct {
	Name     string
	Lastname string
	Password string
	Email    string
	Role     string
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

func (a *AdminService) SearchUsers(ctx context.Context, params *admin.ParamsSearchUsers) (*models.DataSearchUser, error) {
	searched, err := a.adminRepo.Search(ctx,
		&admin.SearchParams{
			QueryString: params.Query,
			Role:        params.Role,
			Page:        params.Page,
			Pagination:  admin.Pagination{OffSet: params.OffSet, Limit: params.Limit},
		},
	)

	if err != nil {
		return nil, err
	}

	return searched, nil
}

// type ParamsFindByID struct {
// 	UserID string
// }

// func (a *AdminService) FindByID(ctx context.Context, params *ParamsFindByID) (*models.ParamsUser, error) {
// 	return nil, nil
// }

// type ParamsDeleteUser struct {
// 	UserID string
// }

// func (a *AdminService) DeleteUser(ctx context.Context, params *ParamsDeleteUser) (*models.ParamsUser, error) {
// 	return nil, nil
// }
