/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	s3cli "github.com/ilia-medvedev-codefresh/aws-s3-otel-metrics/pkg/s3_client"
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("collect called")
		region, err := cmd.Flags().GetString("region")

		if err != nil {
			fmt.Println("Error getting region flag:", err)
			return
		}


		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Replace with your region
		})

		if err != nil {
			fmt.Println("Error creating session:", err)
			return
		}

		svc := s3.New(sess)

		buckets, err := cmd.Flags().GetStringArray("bucket")

		if err != nil {
			fmt.Println("Error getting bucket flag:", err)
			return
		}

		for _, bucket := range buckets {
			err := s3cli.ListObjects(svc, bucket)
			if err != nil {
				fmt.Println("Error listing objects:", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().String("region", "", "AWS region")
	collectCmd.Flags().StringArray("bucket", []string{}, "List of S3 buckets to collect metrics from")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
