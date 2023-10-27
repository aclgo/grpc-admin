package logger

import (
	"log"
	"os"

	"github.com/aclgo/grpc-admin/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Init() error
	Debug(args ...any)
	Debugf(template string, args ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Warn(args ...any)
	Warnf(template string, args ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	Panic(args ...any)
	Panicf(template string, args ...any)
	Fatal(args ...any)
	Fatalf(template string, args ...any)
}

type apiLogger struct {
	config        *config.Config
	sugaredLogger *zap.SugaredLogger
}

func NewapiLogger(config *config.Config) *apiLogger {

	logger := &apiLogger{
		config: config,
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("NewapiLogger.Init: %v", err)
	}

	return logger
}

func (a *apiLogger) getLogLevel() zapcore.Level {
	levels := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
	}

	level, ok := levels[a.config.LogLevel]
	if ok {
		return level
	}

	return zapcore.DebugLevel
}

func (a *apiLogger) Init() error {
	logLevel := a.getLogLevel()

	logWriter := zapcore.AddSync(os.Stdout)

	encConfig := zapcore.EncoderConfig{}

	if a.config.DevMode {
		encConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encConfig = zap.NewProductionEncoderConfig()
	}

	encConfig.NameKey = "[SERVICE]"
	encConfig.TimeKey = "[TIME]"
	encConfig.LevelKey = "[LEVEL]"
	encConfig.CallerKey = "[LINE]"
	encConfig.MessageKey = "[MESSAGE]"
	encConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encConfig.EncodeDuration = zapcore.StringDurationEncoder

	var enc zapcore.Encoder
	if a.config.Encoding == "console" {
		encConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encConfig.EncodeCaller = zapcore.FullCallerEncoder
		encConfig.ConsoleSeparator = " | "
		enc = zapcore.NewConsoleEncoder(encConfig)
	} else {
		encConfig.FunctionKey = "[CALLER]"
		encConfig.EncodeName = zapcore.FullNameEncoder
		enc = zapcore.NewJSONEncoder(encConfig)
	}

	core := zapcore.NewCore(enc, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	a.sugaredLogger = logger.Sugar()
	if err := a.sugaredLogger.Sync(); err != nil {
		log.Printf("sugaredLogger.Sync: %v", err)
	}

	return nil
}

func (a *apiLogger) Debug(args ...any) {
	a.sugaredLogger.Debug(args...)
}
func (a *apiLogger) Debugf(template string, args ...any) {
	a.sugaredLogger.Debugf(template, args...)
}
func (a *apiLogger) Info(args ...any) {
	a.sugaredLogger.Info(args...)
}
func (a *apiLogger) Infof(template string, args ...any) {
	a.sugaredLogger.Infof(template, args...)
}
func (a *apiLogger) Warn(args ...any) {
	a.sugaredLogger.Warn(args...)
}
func (a *apiLogger) Warnf(template string, args ...any) {
	a.sugaredLogger.Warnf(template, args...)
}
func (a *apiLogger) Error(args ...any) {
	a.sugaredLogger.Error(args...)
}
func (a *apiLogger) Errorf(template string, args ...any) {
	a.sugaredLogger.Errorf(template, args...)
}
func (a *apiLogger) Panic(args ...any) {
	a.sugaredLogger.Panic(args...)
}
func (a *apiLogger) Panicf(template string, args ...any) {
	a.sugaredLogger.Panicf(template, args...)
}
func (a *apiLogger) Fatal(args ...any) {
	a.sugaredLogger.Fatal(args...)
}
func (a *apiLogger) Fatalf(template string, args ...any) {
	a.sugaredLogger.Fatalf(template, args...)
}
