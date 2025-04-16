package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	s3cli "github.com/ilia-medvedev-codefresh/aws-s3-otel-metrics/pkg/s3_client"
	telemetry "github.com/ilia-medvedev-codefresh/aws-s3-otel-metrics/pkg/telemetry"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
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
				fmt.Println("Error: AWS region not specified. Use --region flag or set AWS_REGION environment variable.")
				return
			} else {
				region = reg
			}
		}

		err, s3client := s3cli.NewS3Client(region)

		if err != nil {
			fmt.Println("Error creating S3 client:", err)
			return
		}


		buckets, _ := cmd.Flags().GetStringArray("bucket")

		if len(buckets) == 0 {
			fmt.Println("Error: No S3 buckets specified")
			return
		}

		keyAggregationDepth, _ := cmd.Flags().GetInt("key-aggregation-depth")

		exp, err := stdoutmetric.New()

		if err != nil {
			log.Fatal(err)
		}

		meter, err := telemetry.NewMeter(context.TODO(), exp)

		if err != nil {
			log.Fatal(err)
		}

		gauage, _ := meter.Meter.Float64Gauge("s3.aggregated.object.size.bytes", metric.WithDescription("Aggregated object size in bytes"), metric.WithUnit("Bytes"))

		for _, bucket := range buckets {
			err, objects := s3client.ListObjectSizeBytes(bucket, keyAggregationDepth)

			if err != nil {
				fmt.Println("Error listing objects:", err)
				return
			}

			for k,v := range objects {
				gauage.Record(meter.Context, float64(v), metric.WithAttributes(
					attribute.String("bucket", bucket),
					attribute.String("aggregate.key", k),
					))
				//fmt.Printf("Bucket: %s, Object: %s, Size: %d bytes\n", bucket, k, v)
			}

			err = meter.Collect()

			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().String("region", "", "AWS region")
	collectCmd.Flags().StringArray("bucket", []string{}, "List of S3 buckets to collect metrics from")
	collectCmd.Flags().Int("key-aggregation-depth", 0, "Key depth for object size aggregation metric")
	collectCmd.MarkFlagRequired("bucket")
}
