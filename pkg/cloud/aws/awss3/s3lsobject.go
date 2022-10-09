package awss3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/teamssix/cf/pkg/util/errutil"
)

var objectsnum int

func ListObjectsV2(bucket string, region string, s3LsObjectNumberInt64 int64, MaxKeys int64) ([]string, []int64) {
	var (
		objectsKey  []string
		objectsSize []int64
	)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(MaxKeys),
	}
	svc := S3Client(region)
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		errutil.HandleErr(err)
	}
	for _, v := range result.Contents {
		objectsKey = append(objectsKey, *v.Key)
		objectsSize = append(objectsSize, *v.Size)
	}
	fmt.Println(result)
	return objectsKey, objectsSize
}
