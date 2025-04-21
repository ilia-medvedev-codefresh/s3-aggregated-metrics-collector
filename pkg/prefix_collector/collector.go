package prefix_collector

import (
	"log"
	s3cli "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/s3_client"
	telemetry "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type PrefixCollector struct {
	s3Client *s3cli.S3Client
	meter     *telemetry.OtelMeter
	sizesGauage metric.Float64Gauge
	totalObjectsGauage metric.Float64Gauge
}

func NewPrefixCollector(s3Client *s3cli.S3Client, meter *telemetry.OtelMeter) (*PrefixCollector, error) {

	sizesGauage, err := meter.Meter.Float64Gauge("s3.prefix.size", metric.WithDescription("Aggregated object size in bytes"), metric.WithUnit("bytes"))

	if err != nil {
		return nil, err
	}

	// We use gague as number of objects can decrease over time as object get deleted by lifecycle policies for example
	totalObjectsGauage, err := meter.Meter.Float64Gauge("s3.prefix.object.total", metric.WithDescription("Total objects in prefix"))

	if err != nil {
		return nil, err
	}

	return &PrefixCollector{
		s3Client: s3Client,
		meter:     meter,
		sizesGauage: sizesGauage,
		totalObjectsGauage: totalObjectsGauage,
	}, nil
}

func (pc *PrefixCollector) Collect(bucket string, keyAggregationDepth int) error {

	var err error

	log.Printf("Listing objects for bucket: %s", bucket)

	err, objects := pc.s3Client.AggregateObjectsByDepth(bucket, keyAggregationDepth)

	if err != nil {
		return err
	}

	log.Printf("Listing objects for bucket: %s Done!", bucket)

	// Record metrics
	for k,v := range objects {

		pc.sizesGauage.Record(pc.meter.Context, float64(v.TotalSize), metric.WithAttributes(
			attribute.String("bucket", bucket),
			attribute.String("prefix", k),
			))

		pc.totalObjectsGauage.Record(pc.meter.Context, float64(v.ObjectCount), metric.WithAttributes(
		attribute.String("bucket", bucket),
		attribute.String("prefix", k),
		))
	}

	log.Printf("Collecting metrics for bucket: %s", bucket)

	err = pc.meter.Collect()

	if err != nil {
		return err
	}

	log.Printf("Collecting metrics for bucket: %s Done!", bucket)

	return nil
}
