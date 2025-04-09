package s3_client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ListObjects(svc *s3.S3, bucket string) error {
	// Set the parameters for listing objects in the S3 bucket
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	// Paginate through all objects in the bucket
	err := svc.ListObjectsV2Pages(params, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			// Print object key and size
			fmt.Printf("Object Key: %s, Size: %d bytes\n", *obj.Key, *obj.Size)
		}
		return true
	})

	return err
}
