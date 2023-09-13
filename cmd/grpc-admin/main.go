package main

import (
	"log"

	"github.com/aclgo/grpc-admin/internal/server"
	"github.com/aclgo/grpc-admin/pkg/postgres"
	_ "github.com/lib/pq"
)

var (
	dbDriver = "postgres"
	dbUri    = "postgres://grpc-jwt:grpc-jwt@localhost:5432/grpc-jwt?sslmode=disable"
)

func main() {
	db := postgres.Connect(dbDriver, dbUri)
	server := server.NewServer(db)

	log.Fatal(server.Run(3000))
}
