package huaweiobs

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"strconv"
)

var (
	TimestampType = util.ReturnTimestampType("huawei", "obs")
	header        = []string{"序号 (SN)", "名称 (Name)", "存储桶 ACL (Bucket ACL)", "对象数量 (Object Number)", "存储桶大小 (Bucket Size)", "区域 (Region)", "存储桶地址 (Bucket URL)"}
)

func ListBuckets() []string {
	var buckets []string
	listBucketsResult, err := obsClient("all").ListBuckets(nil)
	errutil.HandleErr(err)
	for _, v := range listBucketsResult.Buckets {
		buckets = append(buckets, v.Name)
	}
	return buckets
}

func GetBucketRegion(bucketName string) string {
	bucketRegion, err := obsClient("all").GetBucketLocation(bucketName)
	errutil.HandleErr(err)
	return bucketRegion.Location
}

func GetBucketAcl(bucketName string) string {
	var (
		bucketAclStr string
		read         int
		write        int
		readACP      int
		writeACP     int
		sum          int
	)
	bucketACL, err := obsClient("all").GetBucketAcl(bucketName)
	errutil.HandleErr(err)
	for _, v := range bucketACL.Grants {
		if v.Grantee.Type == "Group" && v.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
			switch v.Permission {
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
	sum = read + write + readACP + writeACP
	switch sum {
	case 0:
		bucketAclStr = "私有 (Private)"
	case 1:
		bucketAclStr = "公共读 (Public Read)"
	case 2:
		bucketAclStr = "公共写 (Public Write)"
	case 3:
		bucketAclStr = "公共读写 (Public Read Write)"
	case 4:
		bucketAclStr = "ACL 公共读 (Bucket ACL Public Read)"
	case 5:
		bucketAclStr = "存储桶和 ACL 公共读 (Bucket and Bucket ACL are Public Read)"
	case 6:
		bucketAclStr = "存储桶公共写、ACL 公共读 (Bucket Public Write, Bucket ACL Public Read)"
	case 7:
		bucketAclStr = "存储桶公共写读写、ACL 公共读 (Bucket Public Read Write, Bucket ACL Public Read)"
	case 8:
		bucketAclStr = "ACL 公共写 (Bucket ACL Public Write)"
	case 9:
		bucketAclStr = "存储桶公共读、ACL 公共写 (Bucket Public Read, Bucket ACL Public Write)"
	case 10:
		bucketAclStr = "存储桶和 ACL 公共写 (Bucket and Bucket ACL are Public Write)"
	case 11:
		bucketAclStr = "存储桶公共读写、ACL 公共写 (Bucket Public Read Write, Bucket ACL Public Write)"
	case 12:
		bucketAclStr = "ACL 公共读写 (Bucket ACL Public Read Write)"
	case 13:
		bucketAclStr = "存储桶公共读、ACL 公共读写 (Bucket Public Read, Bucket ACL Public Read Write)"
	case 14:
		bucketAclStr = "存储桶公共写、ACL 公共读写 (Bucket Public Write, Bucket ACL Public Read Write)"
	case 15:
		bucketAclStr = "存储桶和 ACL 公共读写 (Bucket and Bucket ACL are Public Read Write)"
	}
	return bucketAclStr
}

func getBucketObjectSum(bucket string, region string, s3LsObjectNumber string) (string, string) {
	log.Infof("正在获取 %s 存储桶的数据 (Fetching data for %s bucket)", bucket, bucket)
	var (
		objectsKeyNum  string
		objectsSizeSum string
		n              int64
	)
	ListObjects(bucket, region, s3LsObjectNumber, "")
	objectsKeyNum = strconv.Itoa(len(objectsKey))
	for _, v := range objectsSize {
		n += v
	}
	objectsSizeSum = pubutil.FormatFileSize(n)
	return objectsKeyNum, objectsSizeSum
}

func PrintBucketsListRealTime(region string, obsLsObjectNumber string) {
	var (
		num     int
		dataLen int
	)
	buckets := ListBuckets()
	log.Infof("在全部区域下获取到 %d 条 obs Bucket 信息 (Find %d obs Bucket under all areas)", len(buckets), len(buckets))
	var data = make([][]string, len(buckets))
	for i, o := range buckets {
		SN := strconv.Itoa(i + 1)
		bucketRegion := GetBucketRegion(o)
		if region == bucketRegion {
			bucketACL := GetBucketAcl(o)
			objectsKeyNum, objectsSizeSum := getBucketObjectSum(o, region, obsLsObjectNumber)
			data[num] = []string{SN, o, bucketACL, objectsKeyNum, objectsSizeSum, region, "https://" + o + ".obs." + bucketRegion + ".myhuaweicloud.com"}
			num = num + 1
			dataLen = dataLen + 1
		} else if region == "all" {
			bucketACL := GetBucketAcl(o)
			objectsKeyNum, objectsSizeSum := getBucketObjectSum(o, bucketRegion, obsLsObjectNumber)
			data[num] = []string{SN, o, bucketACL, objectsKeyNum, objectsSizeSum, bucketRegion, "https://" + o + ".obs." + bucketRegion + ".myhuaweicloud.com"}
			dataLen = dataLen + 1
		}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if dataLen == 0 {
		log.Info("没发现存储桶 (No Buckets Found)")
	} else {
		Caption := "OBS 资源 (OBS resources)"
		cloud.PrintTable(td, Caption)
		util.WriteTimestamp(TimestampType)
	}
	cmdutil.WriteCacheFile(td, "huawei", "obs", "all", "all")
}

func PrintBucketsListHistory(region string) {
	cmdutil.PrintOSSCacheFile(header, region, "huawei", "obs", "all")
}

func PrintBucketsList(region string, lsFlushCache bool, obsLsObjectNumber string) {
	if lsFlushCache {
		PrintBucketsListRealTime(region, obsLsObjectNumber)
	} else {
		oldTimestamp := util.ReadTimestamp(TimestampType)
		if oldTimestamp == 0 {
			PrintBucketsListRealTime(region, obsLsObjectNumber)
		} else if util.IsFlushCache(oldTimestamp) {
			PrintBucketsListRealTime(region, obsLsObjectNumber)
		} else {
			util.TimeDifference(oldTimestamp)
			PrintBucketsListHistory(region)
		}
	}
}
