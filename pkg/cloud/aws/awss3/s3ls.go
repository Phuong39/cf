package awss3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
	"strconv"
)

var (
	TimestampType = util.ReturnTimestampType("aws", "s3")
	header        = []string{"序号 (SN)", "名称 (Name)", "存储桶 ACL (Bucket ACL)", "对象数量 (Object Number)", "存储桶大小 (Bucket Size)", "区域 (Region)", "存储桶地址 (Bucket URL)"}
)

func ListBuckets() []string {
	var buckets []string
	input := &s3.ListBucketsInput{}
	svc := S3Client()
	result, err := svc.ListBuckets(input)
	if err != nil {
		errutil.HandleErr(err)
	}
	for _, v := range result.Buckets {
		buckets = append(buckets, *v.Name)
	}
	return buckets
}

func PrintBucketsListRealTime(region string, s3LsObjectNumber string) {
	buckets := ListBuckets()
	log.Debugf("获取到 %d 条 S3 Bucket 信息 (Obtained %d S3 Bucket information)", len(buckets), len(buckets))

	var data = make([][]string, len(buckets))
	for i, o := range buckets {
		SN := strconv.Itoa(i + 1)
		data[i] = []string{SN, o, "", "", "", "", ""}
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
