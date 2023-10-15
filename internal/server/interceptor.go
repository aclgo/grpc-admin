package server

import (
	"context"
	"time"

	"github.com/aclgo/grpc-admin/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Interceptor struct {
	logger logger.Logger
}

func NewInterceptor(logger logger.Logger) *Interceptor {
	return &Interceptor{
		logger: logger,
	}
}

func (i *Interceptor) GrpcInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	i.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}
