package logger

import "github.com/aclgo/grpc-admin/config"

type Logger interface {
}

type apiLogger struct{}

func NewapiLogger(config *config.Config) *apiLogger {
	return &apiLogger{}
}
