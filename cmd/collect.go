package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	s3cli "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/s3_client"
	telemetry "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/telemetry"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Run collection once",

	Run: func(cmd *cobra.Command, args []string) {

		region, _ := cmd.Flags().GetString("region")

		if region == "" {
			if reg, exists := os.LookupEnv("AWS_REGION"); !exists {
				log.Fatal("Error: AWS region not specified. Use --region flag or set AWS_REGION environment variable.")
			} else {
				region = reg
			}
		}

		err, s3client := s3cli.NewS3Client(region)

		if err != nil {
			log.Fatal("Error creating S3 client:", err)
		}


		buckets, _ := cmd.Flags().GetStringArray("bucket")

		if len(buckets) == 0 {
			log.Fatal("Error: No S3 buckets specified")
		}

		keyAggregationDepth, _ := cmd.Flags().GetInt("key-aggregation-depth")

		otelContext := context.TODO()

		otelGrpcEndpoint, _ := cmd.Flags().GetString("otel-grpc-endpoint")

		exp,err := telemetry.NewGRPCExporter(otelContext, otelGrpcEndpoint)

		if err != nil {
			log.Fatal("Error creating OTEL GRPC exporter:", err)
		}

		meter, err := telemetry.NewMeter(otelContext, exp)

		defer func() {
			err = meter.Shutdown(err)
			if err != nil {
				log.Fatal("OTEL collection failed:", err)
			}
		}()

		if err != nil {
			log.Fatal(err)
		}

		sizesGauage, _ := meter.Meter.Float64Gauge("s3.prefix.size", metric.WithDescription("Aggregated object size in bytes"), metric.WithUnit("bytes"))
		// We use gague as number of objects can decrease over time as object get deleted by lifecycle policies for example
		totalObjectsGauage, _ := meter.Meter.Float64Gauge("s3.prefix.object.total", metric.WithDescription("Total objects in prefix"))

		for _, bucket := range buckets {
			err, objects := s3client.AggregateObjectsByDepth(bucket, keyAggregationDepth)

			if err != nil {
				fmt.Println("Error listing objects:", err)
				return
			}

			// Record metrics
			for k,v := range objects {
				sizesGauage.Record(meter.Context, float64(v.TotalSize), metric.WithAttributes(
					attribute.String("bucket", bucket),
					attribute.String("prefix", k),
					))

					totalObjectsGauage.Record(meter.Context, float64(v.ObjectCount), metric.WithAttributes(
					attribute.String("bucket", bucket),
					attribute.String("prefix", k),
					))
			}

			err = meter.Collect()
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().String("region", "", "AWS region")
	collectCmd.Flags().StringArray("bucket", []string{}, "List of S3 buckets to collect metrics from")
	collectCmd.Flags().String("otel-grpc-endpoint", "localhost:4317", "Open Telemetry receiver gRPC endpoint")
	collectCmd.Flags().Int("key-aggregation-depth", 0, "Key depth for object size aggregation metric")
	_ = collectCmd.MarkFlagRequired("bucket")
	_ = collectCmd.MarkFlagRequired("otel-grpc-endpoint")
}
