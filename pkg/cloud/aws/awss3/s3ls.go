package awss3

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
)

var (
	TimestampType = util.ReturnTimestampType("aws", "s3")
	header        = []string{"序号 (SN)", "名称 (Name)", "存储桶 ACL (Bucket ACL)", "对象数量 (Object Number)", "存储桶大小 (Bucket Size)", "区域 (Region)", "存储桶地址 (Bucket URL)"}
)

func ListBuckets() []string {
	var buckets []string
	input := &s3.ListBucketsInput{}
	svc := S3Client("")
	result, err := svc.ListBuckets(input)
	if err != nil {
		errutil.HandleErr(err)
	}

	for _, v := range result.Buckets {
		buckets = append(buckets, *v.Name)
	}
	return buckets
}

func GetBucketRegion(bucket string) string {
	region, err := s3manager.GetBucketRegionWithClient(context.Background(), S3Client(""), bucket)
	if err != nil {
		errutil.HandleErr(err)
	}
	return region
}

func FindBucketAcl(bucket string, region string) string {
	var (
		bucketACL string
		read      int
		write     int
		readACP   int
		writeACP  int
		sum       int
	)
	input := &s3.GetBucketAclInput{
		Bucket: aws.String(bucket),
	}
	svc := S3Client(region)
	result, err := svc.GetBucketAcl(input)
	if err != nil {
		errutil.HandleErr(err)
	}
	for _, v := range result.Grants {
		if *v.Grantee.Type == "Group" {
			if *v.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
				switch *v.Permission {
				case "READ":
					read = 1
				case "WRITE":
					write = 2
				case "READ_ACP":
					readACP = 4
				case "WRITE_ACP":
					writeACP = 8
				}
			}
		}
	}
	sum = read + write + readACP + writeACP
	switch sum {
	case 0:
		bucketACL = "私有 (Private)"
	case 1:
		bucketACL = "公共读 (Public Read)"
	case 2:
		bucketACL = "公共写 (Public Write)"
	case 3:
		bucketACL = "公共读写 (Public Read Write)"
	case 4:
		bucketACL = "ACL 公共读 (Bucket ACL Public Read)"
	case 5:
		bucketACL = "存储桶和 ACL 公共读 (Bucket and Bucket ACL are Public Read)"
	case 6:
		bucketACL = "存储桶公共写、ACL 公共读 (Bucket Public Write, Bucket ACL Public Read)"
	case 7:
		bucketACL = "存储桶公共写读写、ACL 公共读 (Bucket Public Read Write, Bucket ACL Public Read)"
	case 8:
		bucketACL = "ACL 公共写 (Bucket ACL Public Write)"
	case 9:
		bucketACL = "存储桶公共读、ACL 公共写 (Bucket Public Read, Bucket ACL Public Write)"
	case 10:
		bucketACL = "存储桶和 ACL 公共写 (Bucket and Bucket ACL are Public Write)"
	case 11:
		bucketACL = "存储桶公共读写、ACL 公共写 (Bucket Public Read Write, Bucket ACL Public Write)"
	case 12:
		bucketACL = "ACL 公共读写 (Bucket ACL Public Read Write)"
	case 13:
		bucketACL = "存储桶公共读、ACL 公共读写 (Bucket Public Read, Bucket ACL Public Read Write)"
	case 14:
		bucketACL = "存储桶公共写、ACL 公共读写 (Bucket Public Write, Bucket ACL Public Read Write)"
	case 15:
		bucketACL = "存储桶和 ACL 公共读写 (Bucket and Bucket ACL are Public Read Write)"
	}
	return bucketACL
}

func PrintBucketsListRealTime(region string, s3LsObjectNumber string) {
	buckets := ListBuckets()
	log.Debugf("获取到 %d 条 S3 Bucket 信息 (Obtained %d S3 Bucket information)", len(buckets), len(buckets))

	var data = make([][]string, len(buckets))
	for i, o := range buckets {
		SN := strconv.Itoa(i + 1)
		bucketRegion := GetBucketRegion(o)
		bucketACL := FindBucketAcl(o, bucketRegion)
		data[i] = []string{SN, o, bucketACL, "", "", bucketRegion, ""}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("没发现存储桶 (No Buckets Found)")
	} else {
		Caption := "AWS 资源 (AWS resources)"
		cloud.PrintTable(td, Caption)
	}
	cmdutil.WriteCacheFile(td, "aws", "s3", "all", "all")
	util.WriteTimestamp(TimestampType)
}

func PrintBucketsListHistory(region string) {
	cmdutil.PrintOSSCacheFile(header, region, "aws", "s3")
}

func PrintBucketsList(region string, lsFlushCache bool, s3LsObjectNumber string) {
	if lsFlushCache {
		PrintBucketsListRealTime(region, s3LsObjectNumber)
	} else {
		oldTimestamp := util.ReadTimestamp(TimestampType)
		if oldTimestamp == 0 {
			PrintBucketsListRealTime(region, s3LsObjectNumber)
		} else if util.IsFlushCache(oldTimestamp) {
			PrintBucketsListRealTime(region, s3LsObjectNumber)
		} else {
			util.TimeDifference(oldTimestamp)
			PrintBucketsListHistory(region)
		}
	}
}
