package infrastructure

import (
	"context"
	"os"
	"time"

	"github.com/indrabrata/observability-playground/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}

	conn, err := grpc.NewClient(otlpEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("failed to create gRPC connection to OTLP collector", zap.Error(err))
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		zap.L().Fatal("failed to initialize OTLP trace exporter", zap.Error(err))
	}

	// Note : A TracerProvider is a factory/registry that creates and manages telemetry instruments. It's the central configuration hub.
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
