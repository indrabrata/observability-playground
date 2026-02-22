package infrastructure

import (
	"context"
	"time"

	"github.com/indrabrata/observability-playground/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
)

func NewOpenTelemetryTrace(ctx context.Context) *trace.TracerProvider {
	// Note : Ensuring trace context is passed along with requests to different services. This is handled by propagators in OpenTelemetry
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(constant.APP_NAME),
		),
	)
	if err != nil {
		zap.L().Fatal("failed to create resource", zap.Error(err))
	}

	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		zap.L().Fatal("failed to initialize trace exporter", zap.Error(err))
	}

	// Note : A Provideris a factory/registry that creates and manages telemetry instruments. It's the central configuration hub.
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)

	// Register as global so otelsql and other instrumentation libraries can find it.
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider
}
