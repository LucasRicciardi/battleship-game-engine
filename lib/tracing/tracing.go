package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"os"
)

var (
	tracerProvider *trace.TracerProvider
	tracer         otel.Tracer
)

// Init initializes the OpenTelemetry tracing
func Init(serviceName string) error {
	// Create exporter based on environment
	var exporter trace.SpanExporter
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

	// Create trace provider with the exporter
	tracerProvider = trace.NewTracerProvider(
		trace.WithSyncer(exporter),
		trace.WithResource(newResource(serviceName)),
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
func Tracer() otel.Tracer {
	return tracer
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string) (context.Context, func()) {
	return tracer.Start(ctx, name)
}

// noOpExporter is a no-op span exporter for when no tracing is configured
type noOpExporter struct{}

func (e *noOpExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

func (e *noOpExporter) Shutdown(ctx context.Context) error {
	return nil
}

// newResource creates a new resource with service name
func newResource(serviceName string) *Resource {
	return &Resource{
		Attributes: []Attribute{
			{Key: "service.name", Value: StringValue(serviceName)},
			{Key: "service.version", Value: StringValue("1.0.0")},
		},
	}
}

// Resource represents an OpenTelemetry resource
type Resource struct {
	Attributes []Attribute
}

// Attribute represents an attribute key-value pair
type Attribute struct {
	Key   string
	Value Value
}

// Value represents an attribute value
type Value struct {
	StringValue string
}

// StringValue creates a string value
func StringValue(s string) Value {
	return Value{StringValue: s}
}
