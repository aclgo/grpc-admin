package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aclgo/grpc-admin/internal/admin/repository"
	"github.com/aclgo/grpc-admin/internal/admin/usecase"
	"github.com/aclgo/grpc-admin/internal/delivery/grpc/service"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Server struct {
	db *sqlx.DB
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

	ctx := context.Background()

	http.HandleFunc("/create", handlers.Register(ctx))
	http.HandleFunc("/search", handlers.Search(ctx))

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		return errors.Wrap(err, "Run.ListenAndServe")
	}

	return nil
}
