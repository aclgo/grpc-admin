package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aclgo/grpc-admin/config"
	"github.com/aclgo/grpc-admin/internal/admin"
	"github.com/aclgo/grpc-admin/internal/server"
	tel "github.com/aclgo/grpc-admin/internal/telemetry/otel"
	"github.com/aclgo/grpc-admin/pkg/logger"
	"github.com/aclgo/grpc-admin/pkg/postgres"
)

var (
	dbDriver = "postgres"
	dbUri    = "postgresql://grpc-admin:grpc-admin@db:5432/grpc-admin?sslmode=disable"

	collectorURL = "otel-collector:4317"
	zipkinURL    = "http://zipkin:9411/api/v2/spans"
)

func main() {

	cfg := config.NewConfig(".")
	cfg.ApiPort = "50052"
	cfg.Observability.CollectorURL = collectorURL
	cfg.Observability.ZipkinURL = zipkinURL

	logger := logger.NewapiLogger(cfg)

	otel, err := tel.NewOtel(cfg, logger, "grpc-admin", "0.01")
	if err != nil {
		logger.Errorf("otel.NewOtel: %v", err)
		return
	}

	defer func() {
		otel.Shutdowns(context.Background())
	}()

	mt := otel.Meter("grpc-jwt")
	tr := otel.Tracer("grpc-jwt")

	db := postgres.Connect(dbDriver, dbUri)

	server := server.NewServer(cfg, db, logger, &admin.Observability{Meter: mt, Trace: tr})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := server.Run(ctx); err != nil {
		logger.Errorf("server.Run: %v", err)
	}
}
