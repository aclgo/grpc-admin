package server

import (
	"fmt"
	"net"

	"github.com/aclgo/grpc-admin/internal/admin/repository"
	"github.com/aclgo/grpc-admin/internal/admin/usecase"
	"github.com/aclgo/grpc-admin/internal/delivery/grpc/service"
	"github.com/aclgo/grpc-admin/pkg/logger"
	proto "github.com/aclgo/grpc-admin/proto/admin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewServer(db *sqlx.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Run(port int) error {
	adminRepo := repository.NewpostgresRepo(s.db)
	adminUC := usecase.NewAdminService(adminRepo)

	handlers := service.NewAdminService(adminUC)

	// ctx := context.Background()
	interceptor := NewInterceptor(s.logger)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptor.GrpcInterceptor),
	}

	srv := grpc.NewServer(opts...)

	proto.RegisterAdminServiceServer(srv, handlers)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "Run.Listen")
	}

	// s.logger.Infof("server grpc running port %v", port)
	if err := srv.Serve(listen); err != nil {
		return errors.Wrap(err, "Run.Serve")
	}

	return nil
}
