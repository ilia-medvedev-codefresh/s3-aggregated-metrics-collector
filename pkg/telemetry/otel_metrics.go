package telemetry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type OtelMeter struct {
    reader      metricsdk.Reader
    exporter    metricsdk.Exporter
    shutdownFunctions    []func(ctx context.Context) error
    Meter metric.Meter
    Context context.Context
}

func NewMeter(ctx context.Context, exporter metricsdk.Exporter) (*OtelMeter, error) {
    // Because we are using OTEL in a cobra-cli application, we are creating a manual reader that will read metrics on demand instead of in the background
    reader := metricsdk.NewManualReader()

    resource, err := getResource()

    if err != nil {
      return nil, fmt.Errorf("could not get resource: %w", err)
    }

      if err != nil {
          return nil, fmt.Errorf("could not get collector exporter: %w", err)
      }


      provider := metricsdk.NewMeterProvider(
          metricsdk.WithResource(resource),
          metricsdk.WithReader(reader),
      )

    if err != nil {
        return nil, fmt.Errorf("could not create meter provider: %w", err)
    }

    return &OtelMeter{
        Context: ctx,
        exporter: exporter,
        Meter: provider.Meter("s3-aggregated-otel-metrics"),
        reader: reader,
        shutdownFunctions: []func(ctx context.Context) error{
            provider.Shutdown,
        },
    }, nil
}

func (m *OtelMeter) Collect() error {
    var err error

    ctx, cancel := context.WithTimeout(m.Context, 5*time.Second)

    defer func() {
        cancel()
    }()

    collectedMetrics := &metricdata.ResourceMetrics{}

    if err := m.reader.Collect(ctx, collectedMetrics); err != nil {
        return fmt.Errorf("could not collect metrics: %w", err)
    }

    err = m.exporter.Export(context.TODO(), collectedMetrics)

    if err != nil {
        return fmt.Errorf("could not export metrics: %w", err)
    }

    return err
}

func (m *OtelMeter) Shutdown(err error) error {
    for _, fn := range m.shutdownFunctions {
		err = errors.Join(err, fn(m.Context))
	}

	return err
}

func getResource() (*resource.Resource, error) {
    resource, err := resource.Merge(
        resource.Default(),
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("s3-aggregated-otel-metrics"),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("could not merge resources: %w", err)
  }

    return resource, nil
}


func NewGRPCExporter(ctx context.Context, endpoint string) (metricsdk.Exporter, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    exporter, err := otlpmetricgrpc.New(ctx,
        otlpmetricgrpc.WithEndpoint(endpoint),
        otlpmetricgrpc.WithCompressor("gzip"),
        otlpmetricgrpc.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("could not create metric exporter: %w", err)
    }

    return exporter, nil
}
