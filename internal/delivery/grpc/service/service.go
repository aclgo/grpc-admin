package service

import (
	"github.com/aclgo/grpc-admin/internal/admin"
)

type AdminService struct {
	adminUC admin.AdminUC
}

func NewAdminService(adminUC admin.AdminUC) *AdminService {
	return &AdminService{
		adminUC: adminUC,
	}
}
