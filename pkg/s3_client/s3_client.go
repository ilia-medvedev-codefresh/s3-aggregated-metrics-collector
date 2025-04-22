package s3_client

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client s3.Client
}

type S3Prefix struct {
	TotalSize   int64
	ObjectCount int64
}

func NewS3Client(region string) (*S3Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		return nil, err
	}

	cl := s3.NewFromConfig(cfg)

	if err != nil {
		return nil, err
	}

	return &S3Client{
		Client: *cl,
	}, nil
}

func (c *S3Client) AggregateObjectsByDepth(bucket string, depth int) (map[string]S3Prefix, error) {

	prefixMap := make(map[string]S3Prefix)

	paginator := s3.NewListObjectsV2Paginator(&c.Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	for paginator.HasMorePages() {

		page, err := paginator.NextPage(context.TODO())

		if err != nil {
			return nil, err
		}

		// Log the objects found
		for _, obj := range page.Contents {

			key := strings.Join(splitKeyByDepth(*aws.String(*obj.Key), depth), "/")

			if prefix, ok := prefixMap[key]; ok {
				prefix.TotalSize += *obj.Size
				prefix.ObjectCount++
				prefixMap[key] = prefix
			} else {
				prefixMap[key] = S3Prefix{
					TotalSize:   *obj.Size,
					ObjectCount: 1,
				}
			}
		}
	}

	return prefixMap, nil
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
