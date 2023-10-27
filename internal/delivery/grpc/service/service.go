package service

import (
	"github.com/aclgo/grpc-admin/internal/admin"
	proto "github.com/aclgo/grpc-admin/proto/admin"
)

type AdminService struct {
	adminUC  admin.AdminUC
	observer *admin.Observability
	proto.UnimplementedAdminServiceServer
}

func NewAdminService(adminUC admin.AdminUC, observer *admin.Observability) *AdminService {
	return &AdminService{
		adminUC:  adminUC,
		observer: observer,
	}
}
