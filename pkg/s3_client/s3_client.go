package s3_client

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	Client *s3.S3
}

type S3Prefix struct {
	TotalSize int64
	ObjectCount int64
}

func NewS3Client(region string) (error, *S3Client) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return err, nil
	}

	client := s3.New(sess)

	return nil, &S3Client{
		Client: client,
	}
}

func(c *S3Client) AggregateObjectsByDepth(bucket string, depth int) (error, map[string]S3Prefix) {

	prefixMap := make(map[string]S3Prefix)

	// Set the parameters for listing objects in the S3 bucket
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	// Paginate through all objects in the bucket
	err := c.Client.ListObjectsV2Pages(params, func(page *s3.ListObjectsV2Output, lastPage bool) bool {

		for _, obj := range page.Contents {

				key := strings.Join(splitKeyByDepth(*aws.String(*obj.Key), depth), "/")

				if prefix, ok := prefixMap[key]; ok {
					prefix.TotalSize += *obj.Size
					prefixMap[key] = prefix
				} else {
					prefixMap[key] = S3Prefix{
						TotalSize:   *obj.Size,
						ObjectCount: 1,
				}
			}
		}

		return true
	})

	return err, prefixMap
}

func splitKeyByDepth(key string, depth int) []string {
	splitObject := strings.Split(key, "/")

	if depth > 0 && depth < len(splitObject) {
		keyParts := make([]string, depth)

		for i := 0; i < depth && i < len(splitObject); i++ {
			keyParts[i] = splitObject[i]
		}

		return keyParts
	} else {
		return splitObject
	}
}
