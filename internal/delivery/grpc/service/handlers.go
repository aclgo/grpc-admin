package service

import (
	"context"

	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/models"
	proto "github.com/aclgo/grpc-admin/proto/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AdminService) Register(ctx context.Context, req *proto.ParamsCreateAdmin) (*proto.ParamsUser, error) {
	result, err := s.adminUC.Create(ctx, &admin.ParamsCreateAdmin{
		Name:     req.Name,
		Lastname: req.Lastname,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	})

	if err != nil {
		return nil, err
	}

	return parseModelProto([]*models.ParamsUser{result})[0], nil
}

func (s *AdminService) Search(ctx context.Context, req *proto.ParamsSearchRequest) (*proto.ParamsSearchResponse, error) {
	result, err := s.adminUC.SearchUsers(ctx,
		&admin.ParamsSearchUsers{
			Query: req.Query,
			Role:  req.Role,
			Page:  int(req.Page),
			Pagination: admin.Pagination{
				OffSet: int(req.Offset),
				Limit:  int(req.Limit),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &proto.ParamsSearchResponse{
		Total: int64(result.Total),
		Users: parseModelProto(result.Users),
	}, nil
}

func parseModelProto(items []*models.ParamsUser) []*proto.ParamsUser {
	var users []*proto.ParamsUser

	for _, user := range items {
		user.ClearPass()

		users = append(users, &proto.ParamsUser{
			UserId:    user.UserID,
			Name:      user.Name,
			Lastname:  user.Lastname,
			Password:  user.Password,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return users
}
