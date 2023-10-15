package logger

import "github.com/aclgo/grpc-admin/config"

type Logger interface {
	Infof(template string, args ...any)
}

type apiLogger struct{}

func NewapiLogger(config *config.Config) *apiLogger {
	return &apiLogger{}
}
