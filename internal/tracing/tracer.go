package tracing

import (
	"fmt"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func New(cfg Config, logger *zap.Logger) trace.Tracer { //nolint: ireturn
	if !cfg.Enabled {
		return trace.NewNoopTracerProvider().Tracer("snapp.dispatching")
	}

	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Agent.Host), jaeger.WithAgentPort(cfg.Agent.Port)),
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to initialize export pipeline: %v", err).Error())
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNamespaceKey.String("snapp.dispatching"),
			semconv.ServiceNameKey.String("soteria"),
		),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Ratio))),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	// register the TraceContext propagator globally.
	var tc propagation.TraceContext

	otel.SetTextMapPropagator(tc)

	tracer := otel.Tracer("dispatching/soteria")

	return tracer
}
