package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"

	otlpmetricgrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Setup(
    ctx context.Context,
    serviceName string,
    endpoint string,
) (func(context.Context) error, error) {

    traceExporter, err := otlptracegrpc.New(
        ctx,
        otlptracegrpc.WithEndpoint(endpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    metricExporter, err := otlpmetricgrpc.New(
        ctx,
        otlpmetricgrpc.WithEndpoint(endpoint),
        otlpmetricgrpc.WithInsecure(),
    )
    if err != nil {
        _ = traceExporter.Shutdown(ctx)
        return nil, err
    }

    res := resource.NewWithAttributes(
        "",
        attribute.String("service.name", serviceName),
    )

    tracerProvider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(traceExporter),
        sdktrace.WithResource(res),
    )

    meterProvider := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(
            sdkmetric.NewPeriodicReader(
                metricExporter,
                sdkmetric.WithInterval(15*time.Second),
            ),
        ),
        sdkmetric.WithResource(res),
    )

    otel.SetTracerProvider(tracerProvider)
    otel.SetMeterProvider(meterProvider)

    otel.SetTextMapPropagator(
        propagation.NewCompositeTextMapPropagator(
            propagation.TraceContext{},
            propagation.Baggage{},
        ),
    )

    shutdown := func(ctx context.Context) error {
        return errors.Join(
            meterProvider.Shutdown(ctx),
            tracerProvider.Shutdown(ctx),
        )
    }

    return shutdown, nil
}

var (
    RecordsReceived  metric.Int64Counter
    RecordsProcessed metric.Int64Counter
)

func InitMetrics() {
    meter := otel.Meter("blockberg-terminal/handlers")

    var err error

    RecordsReceived, err = meter.Int64Counter("records.received")
    if err != nil {
        panic(err)
    }

    RecordsProcessed, err = meter.Int64Counter("records.processed")
    if err != nil {
        panic(err)
    }
}