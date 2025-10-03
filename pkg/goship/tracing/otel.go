package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer sets up a simple stdout exporter and tracer provider for development.
//
// Parameters:
//   - serviceName: Name of the service (e.g., "my-app").
//
// Example:
//
//	tp, err := InitTracer("my-app")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer tp.Shutdown(ctx)
func InitTracer(serviceName string) (*sdkTrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		// You can also configure resources, sampler, etc.
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// StartSpan starts a new span and returns the derived context and the span.
//
// Parameters:
//   - ctx: parent context (e.g., from HTTP request).
//   - name: name of the span (e.g., "db.query").
//
// Example:
//
//	ctx, span := StartSpan(ctx, "operation")
//	defer span.End()
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	tracer := otel.Tracer("app")
	return tracer.Start(ctx, name)
}
