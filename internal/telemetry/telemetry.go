package telemetry

import "context"

type Telemetry interface {
	Setup(ctx context.Context, svcName, svcVersion string) (func(context.Context), error)
}
