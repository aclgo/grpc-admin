package service

import (
	"github.com/aclgo/grpc-admin/internal/admin"
	proto "github.com/aclgo/grpc-admin/proto/admin"
)

type AdminService struct {
	adminUC admin.AdminUC
	proto.UnimplementedAdminServiceServer
}

func NewAdminService(adminUC admin.AdminUC) *AdminService {
	return &AdminService{
		adminUC: adminUC,
	}
}
