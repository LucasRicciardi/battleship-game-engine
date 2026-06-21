package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"os"
)

var (
	tracerProvider *tracesdk.TracerProvider
	tracer         trace.Tracer
)

// Init initializes the OpenTelemetry tracing
func Init(serviceName string) error {
	// Create exporter based on environment
	var exporter tracesdk.SpanExporter
	if os.Getenv("TRACING_EXPORT") == "stdout" {
		var err error
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return fmt.Errorf("failed to create stdout exporter: %w", err)
		}
	} else {
		// For production, you would use Jaeger, Datadog, or other exporters
		// For now, use a no-op exporter
		exporter = &noOpExporter{}
	}

	// Create resource with service name
	res, err := resource.New(context.Background(), resource.WithAttributes(
		attribute.String("service.name", serviceName),
		attribute.String("service.version", "1.0.0"),
	))
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider with the exporter
	tracerProvider = tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)
	tracer = otel.Tracer(serviceName)

	return nil
}

// Shutdown shuts down the tracer provider
func Shutdown() error {
	if tracerProvider != nil {
		return tracerProvider.Shutdown(context.Background())
	}
	return nil
}

// Tracer returns the global tracer
func Tracer() trace.Tracer {
	return tracer
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string) (context.Context, func()) {
	_, span := tracer.Start(ctx, name)
	return ctx, func() { span.End() }
}

// noOpExporter is a no-op span exporter for when no tracing is configured
type noOpExporter struct{}

func (e *noOpExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	return nil
}

func (e *noOpExporter) Shutdown(ctx context.Context) error {
	return nil
}
