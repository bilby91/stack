package tracer

import (
	"context"
	"github.com/formancehq/stack/libs/go-libs/time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer = otel.Tracer("com.formance.ledger")

func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer.Start(ctx, name, opts...)
}

func TraceWithLatency[RET any](
	ctx context.Context,
	operationName string,
	fn func(ctx context.Context) (RET, error),
	finalizers ...func(ctx context.Context, ret RET),
) (RET, time.Duration, error) {
	var latency time.Duration
	ret, err := Trace(ctx, operationName, func(ctx context.Context) (RET, error) {
		now := time.Now()
		ret, err := fn(ctx)
		if err != nil {
			var zeroRet RET
			return zeroRet, err
		}

		latency = time.Since(now)

		for _, finalizer := range finalizers {
			finalizer(ctx, ret)
		}

		return ret, nil
	})
	if err != nil {
		return ret, 0, err
	}

	return ret, latency, nil
}
