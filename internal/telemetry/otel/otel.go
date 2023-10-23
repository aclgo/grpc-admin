package tel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials/insecure"
)

type Otel struct {
	svcName        string
	svcVersion     string
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	propagator     propagation.TextMapPropagator
	FnShutdowns    []func(context.Context) error
}

func NewOtel(svcName, svcVersion string) (*Otel, error) {
	t := Otel{
		svcName:    svcName,
		svcVersion: svcVersion,
	}

	err := t.Setup()
	if err != nil {
		return nil, fmt.Errorf("otel.Setup: %v", err)
	}

	otel.SetMeterProvider(t.MeterProvider)
	otel.SetTracerProvider(t.TracerProvider)
	otel.SetTextMapPropagator(t.Propagator())

	return &t, err
}

func (o *Otel) Setup() error {

	var (
		tr sdktrace.SpanExporter
		mt sdkmetric.Exporter
	)

	o.propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			[]attribute.KeyValue{
				semconv.ServiceName(o.svcName),
				semconv.ServiceVersion(o.svcVersion),
			}...,
		),
	)

	if err != nil {
		return fmt.Errorf("resource.Merge: %v", err)
	}

	switch exporter, ok := os.LookupEnv("OTEL_EXPORTER"); {
	case exporter == "stdout":
		tr, err = stdouttrace.New()
		if err != nil {
			return fmt.Errorf("stdouttrace: %v", err)
		}

		mt, err = stdoutmetric.New(stdoutmetric.WithEncoder(json.NewEncoder(os.Stdout)))
		if err != nil {
			return fmt.Errorf("stdoutmetric: %v", err)
		}

	case exporter == "otlp":
		tr, err = otlptracegrpc.New(context.Background(), otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("")
		}

		mt, err = otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("")
		}
	case ok:
		return fmt.Errorf("invalid param env variable OTEL_EXPORTER")
	default:
		o.TracerProvider = trace.NewNoopTracerProvider()
		o.MeterProvider = noop.NewMeterProvider()

		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.AlwaysSample(),
		),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(tr),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(mt),
		),
	)

	o.TracerProvider = tp
	o.MeterProvider = mp

	o.FnShutdowns = append(
		o.FnShutdowns,
		tp.Shutdown,
		mp.Shutdown,
	)

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
