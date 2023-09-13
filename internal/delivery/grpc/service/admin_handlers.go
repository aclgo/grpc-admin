package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aclgo/grpc-admin/internal/admin"
)

func (s *AdminService) Register(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params admin.ParamsCreateAdmin

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		created, err := s.adminUC.Create(ctx, &params)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(created)
	}
}

func (s *AdminService) Search(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		role := r.URL.Query().Get("role")
		page := r.URL.Query().Get("page")
		offset := r.URL.Query().Get("offset")
		limit := r.URL.Query().Get("limit")

		params, err := admin.NewParamsSearchUsers(query, role, page, offset, limit)
		if err != nil {
			log.Println(err)
			return
		}

		users, err := s.adminUC.SearchUsers(ctx, params)
		if err != nil {
			log.Println(err)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}
