package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	prefix_collector "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/prefix_collector"
	s3cli "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/s3_client"
	telemetry "github.com/ilia-medvedev-codefresh/s3-aggregated-otel-metrics/pkg/telemetry"

	"github.com/spf13/cobra"
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
			log.Print("Shutting down OTEL collector...")

			err = meter.Shutdown(err)
			if err != nil {
				log.Fatal("OTEL collection failed:", err)
			}

			log.Println("OTEL collector shutdown done!")
		}()

		if err != nil {
			log.Fatal(err)
		}

		prefixCollector, err := prefix_collector.NewPrefixCollector(s3client, meter)

		if err != nil {
			log.Fatal("Error initalizing collector:", err)
		}

		var wg sync.WaitGroup

		errCh := make(chan error, len(buckets))

		for _, bucket := range buckets {
			wg.Add(1)
			go func(b string) {
				defer wg.Done()
				err := prefixCollector.Collect(b, keyAggregationDepth)
				if err != nil {
					errCh <- fmt.Errorf("collection failed for bucket %s with the following error: %s", b, err.Error())
				}
			}(bucket)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for collectionError := range errCh {
			err = errors.Join(err, collectionError)
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
