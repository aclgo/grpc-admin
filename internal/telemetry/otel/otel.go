package tel

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/grpc-admin/config"
	"github.com/aclgo/grpc-admin/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Otel struct {
	svcName        string
	svcVersion     string
	logger         logger.Logger
	config         *config.Config
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	propagator     propagation.TextMapPropagator
	FnShutdowns    []func(context.Context) error
}

func NewOtel(cfg *config.Config, logger logger.Logger, svcName, svcVersion string) (*Otel, error) {
	t := Otel{
		svcName:    svcName,
		svcVersion: svcVersion,
		logger:     logger,
		config:     cfg,
	}

	t.propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			[]attribute.KeyValue{
				semconv.ServiceName(t.svcName),
				semconv.ServiceVersion(t.svcVersion),
			}...,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("resource.Merge: %v", err)
	}

	if err := t.initTracer(res); err != nil {
		return nil, err
	}

	if err := t.initMeter(res); err != nil {
		return nil, err
	}

	otel.SetMeterProvider(t.MeterProvider)
	otel.SetTracerProvider(t.TracerProvider)
	otel.SetTextMapPropagator(t.Propagator())

	return &t, err
}

func (o *Otel) initTracer(res *resource.Resource) error {
	exporter, err := zipkin.New(
		o.config.Observability.ZipkinURL,
	)

	if err != nil {
		return fmt.Errorf("zipkin.New: %v", err)
	}

	o.logger.Info("zipkin exporter init")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.AlwaysSample(),
		),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	o.TracerProvider = tp

	o.FnShutdowns = append(o.FnShutdowns, tp.Shutdown)

	return nil
}

func (o *Otel) initMeter(res *resource.Resource) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		o.config.Observability.CollectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("initMeter.DialContext: %v", err)
	}

	o.logger.Info("starting connection otel collector")

	mexp, err := otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return fmt.Errorf("otlpmetricgrpc.New: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mexp)),
		sdkmetric.WithResource(res),
	)

	o.MeterProvider = mp

	o.FnShutdowns = append(o.FnShutdowns, mp.Shutdown)

	return nil
}

func (o *Otel) Tracer(tracerName string, opts ...trace.TracerOption) trace.Tracer {
	return o.TracerProvider.Tracer(tracerName, opts...)
}

func (o *Otel) Meter(meterName string, opts ...metric.MeterOption) metric.Meter {
	return o.MeterProvider.Meter(meterName, opts...)
}

func (o *Otel) Propagator() propagation.TextMapPropagator {
	return o.propagator
}

func (o *Otel) Shutdowns(ctx context.Context) error {
	var errs []error

	for _, fn := range o.FnShutdowns {
		errs = append(errs, fn(ctx))
	}

	return errors.Join(errs...)
}
