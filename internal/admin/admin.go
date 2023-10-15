package admin

import (
	"context"
	"strconv"

	"github.com/aclgo/grpc-admin/internal/models"
	"github.com/pkg/errors"
)

type AdminUC interface {
	Create(ctx context.Context, params *ParamsCreateAdmin) (*models.ParamsUser, error)
	SearchUsers(ctx context.Context, params *ParamsSearchUsers) (*models.DataSearchUser, error)
}

type AdminRepo interface {
	Create(context.Context, *models.ParamsCreateAdmin) (*models.ParamsUser, error)
	// Find(context.Context, *models.ParamsFind) (*models.ParamsUser, error)
	Search(context.Context, *ParamsSearchUsers) (*models.DataSearchUser, error)
}

type ParamsCreateAdmin struct {
	Name     string `json:"name"`
	Lastname string `json:"last_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type ParamsSearchUsers struct {
	Query      string
	Role       string
	Page       int
	Pagination Pagination
}

type Pagination struct {
	OffSet int
	Limit  int
}

func NewParamsSearchUsers(query, role, page, offset, limit string) (*ParamsSearchUsers, error) {

	parsedPage := 1
	parsedOffSet := 0
	parsedLimit := 0

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return nil, errors.Wrap(err, "NewParamsSearchUsers: invalid page")
		}
		if pageInt > 0 {
			parsedPage = pageInt
		}
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errors.Wrap(err, "NewParamsSearchUsers: invalid limit")
		}

		if limitInt >= 0 {
			parsedLimit = limitInt
		}
	}

	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, errors.Wrap(err, "NewParamsSearchUsers: invalid offset")
		}
		if offsetInt >= 0 {
			parsedOffSet = parsedLimit * (parsedPage - 1)
		}
	}

	return &ParamsSearchUsers{
		Query: query,
		Role:  role,
		Page:  parsedPage,
		Pagination: Pagination{
			OffSet: parsedOffSet,
			Limit:  parsedLimit,
		},
	}, nil
}
