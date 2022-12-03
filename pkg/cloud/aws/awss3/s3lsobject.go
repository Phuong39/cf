package awss3

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/teamssix/cf/pkg/util/errutil"
)

var (
	objectsKey  []string
	objectsSize []int64
)

func ListObjectsV2(bucket string, region string, s3LsObjectNumber string, NextContinuationToken string) {
	var (
		s3LsObjectNumberInt64 int64
		MaxKeys               int64
		err                   error
		IsTruncated           bool
		input                 *s3.ListObjectsV2Input
	)
	if s3LsObjectNumber != "all" {
		s3LsObjectNumberInt64, err = strconv.ParseInt(s3LsObjectNumber, 10, 64)
		errutil.HandleErr(err)
	}
	if s3LsObjectNumberInt64 == 0 || s3LsObjectNumberInt64 > 1000 {
		MaxKeys = 1000
	} else {
		MaxKeys = s3LsObjectNumberInt64
	}
	if NextContinuationToken == "" {
		input = &s3.ListObjectsV2Input{
			Bucket:  aws.String(bucket),
			MaxKeys: aws.Int64(MaxKeys),
		}
	} else {
		input = &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			MaxKeys:           aws.Int64(MaxKeys),
			ContinuationToken: aws.String(NextContinuationToken),
		}
	}
	svc := S3Client(region)
	result, err := svc.ListObjectsV2(input)
	IsTruncated = *result.IsTruncated
	errutil.HandleErr(err)
	for _, v := range result.Contents {
		objectsKey = append(objectsKey, *v.Key)
		objectsSize = append(objectsSize, *v.Size)
	}
	objectNum := len(objectsKey)
	if objectNum%10000 == 0 && objectNum != 0 {
		log.Infof("当前已获取到 %d 条数据 (%d pieces of data have been obtained)", objectNum, objectNum)
	}
	if objectNum == 100000 {
		var name bool
		prompt := &survey.Confirm{
			Message: "已查询到 10w 条对象，是否继续？如果继续可能会耗费较长时间。(Found up to 100,000 objects, want to continue? If you continue, it may take a long time)",
			Default: true,
		}
		_ = survey.AskOne(prompt, &name)
		if !name {
			IsTruncated = false
			log.Infoln("已停止继续查询对象，您还可以通过 -n 参数指定您想要查询对象的数量。(Has stopped continuing to query objects. You can specify the number of objects to query with the -n parameter.)")
		}
	}
	if s3LsObjectNumber != "all" {
		s3LsObjectNumberInt64, err = strconv.ParseInt(s3LsObjectNumber, 10, 64)
		errutil.HandleErr(err)
		s3LsObjectNumberInt := int(s3LsObjectNumberInt64)
		if s3LsObjectNumberInt >= objectNum {
			IsTruncated = false
			objectNum = s3LsObjectNumberInt
			objectsKey = objectsKey[0:objectNum]
			objectsSize = objectsSize[0:objectNum]
		}
	}
	if IsTruncated {
		ListObjectsV2(bucket, region, s3LsObjectNumber, *result.NextContinuationToken)
	}
}
