package alioss

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

type Bucket = cloud.Resource

type Object struct {
	BucketName   string
	ObjectNumber int
	ObjectSize   int64
}

type Acl struct {
	BucketName string
	Acl        string
}

type objectContents struct {
	Key          string
	Size         int64
	LastModified string
}

type error interface {
	Error() string
}

var (
	objectNum     int
	ObjectSize    int64
	objects       []objectContents
	TimestampType = util.ReturnTimestampType("alibaba", "oss")
	header        = []string{"序号 (SN)", "名称 (Name)", "存储桶 ACL (Bucket ACL)", "对象数量 (Object Number)", "存储桶大小 (Bucket Size)", "区域 (Region)", "存储桶地址 (Bucket URL)"}
)

func (o *OSSCollector) ListBuckets() ([]Bucket, error) {
	region := cloud.GetGlobalRegions()[0]
	o.OSSClient(region)
	var size = 1000
	var out []Bucket
	marker := oss.Marker("")
	var err error
	for {
		var lbr oss.ListBucketsResult
		lbr, err = o.Client.ListBuckets(oss.MaxKeys(size), marker)
		marker = oss.Marker(lbr.NextMarker)
		for _, bucket := range lbr.Buckets {
			obj := Bucket{Name: bucket.Name,
				Region: bucket.Location[4:],
			}
			out = append(out, obj)
		}
		if !lbr.IsTruncated {
			break
		}
	}
	errutil.HandleErrNoExit(err)
	return out, err
}

func (o *OSSCollector) ListObjects(bucketName string, ossLsObjectNumber string) ([]Object, []objectContents) {
	var (
		size    int
		out     []Object
		Buckets []Bucket
	)
	if ossLsObjectNumber == "all" {
		size = 1000
	} else {
		var err error
		size, err = strconv.Atoi(ossLsObjectNumber)
		errutil.HandleErr(err)
	}

	marker := oss.Marker("")
	OSSCollector := &OSSCollector{}
	Buckets, _ = OSSCollector.ListBuckets()
	if bucketName != "all" {
		for _, obj := range Buckets {
			if obj.Name == bucketName {
				Buckets = nil
				Buckets = append(Buckets, obj)
			}
		}
	}
	for _, j := range Buckets {
		BucketName := j.Name
		region := j.Region
		o.OSSClient(region)
		bucket, err := o.Client.Bucket(BucketName)
		errutil.HandleErr(err)
		objects = nil
		getAllObjects(bucket, marker, size, ossLsObjectNumber)
		log.Debugf("在 %s 存储桶中找到了 %d 个对象 (Found %d Objects in %s Bucket)", BucketName, objectNum, objectNum, BucketName)
		obj := Object{
			BucketName:   BucketName,
			ObjectNumber: objectNum,
			ObjectSize:   ObjectSize,
		}
		objectNum = 0
		ObjectSize = 0
		out = append(out, obj)
	}
	return out, objects
}

func getAllObjects(bucket *oss.Bucket, marker oss.Option, size int, ossLsObjectNumber string) {
	lor, err := bucket.ListObjects(oss.MaxKeys(size), marker)
	errutil.HandleErr(err)
	marker = oss.Marker(lor.NextMarker)
	objectNum = objectNum + len(lor.Objects)
	for _, k := range lor.Objects {
		ObjectSize = ObjectSize + k.Size
		obj := objectContents{
			Key:          k.Key,
			Size:         k.Size,
			LastModified: k.LastModified.Format("2006-01-02 15:04:05"),
		}
		objects = append(objects, obj)
	}
	log.Tracef("Next Marker: %s", lor.NextMarker)
	NextMarker := lor.NextMarker
	if objectNum == 100000 {
		var name bool
		prompt := &survey.Confirm{
			Message: "已查询到 10w 条对象，是否继续？如果继续可能会耗费较长时间。(Found up to 100,000 objects, want to continue? If you continue, it may take a long time)",
			Default: true,
		}
		_ = survey.AskOne(prompt, &name)
		if !name {
			NextMarker = ""
			log.Infoln("已停止继续查询对象，你还可以通过 -n 参数指定你想要查询对象的数量。(Has stopped continuing to query objects. You can specify the number of objects to query with the -n parameter.)")
		}
	}
	if ossLsObjectNumber != "all" {
		ossLsObjectNumberInt, err := strconv.Atoi(ossLsObjectNumber)
		errutil.HandleErr(err)
		if objectNum >= ossLsObjectNumberInt {
			NextMarker = ""
			objectNum = ossLsObjectNumberInt
			objects = objects[0:objectNum]
		}
	}
	if NextMarker != "" {
		getAllObjects(bucket, marker, size, ossLsObjectNumber)
	}
}

func (o *OSSCollector) GetBucketACL() []Acl {
	OSSCollector := &OSSCollector{}
	Buckets, _ := OSSCollector.ListBuckets()
	var out []Acl
	for _, j := range Buckets {
		BucketName := j.Name
		region := j.Region
		o.OSSClient(region)
		gbar, err := o.Client.GetBucketACL(BucketName)
		errutil.HandleErr(err)

		BucketACL := gbar.ACL
		if BucketACL == "private" {
			BucketACL = "私有 (Private)"
		} else if BucketACL == "public-read" {
			BucketACL = "公共读 (Public Read)"
		} else if BucketACL == "public-read-write" {
			BucketACL = "公共读写 (Public Read Write)"
		}

		obj := Acl{
			BucketName: BucketName,
			Acl:        BucketACL,
		}
		out = append(out, obj)
	}
	return out
}

func PrintBucketsListRealTime(region string, ossLsObjectNumber string) {
	OSSCollector := &OSSCollector{}
	Buckets, _ := OSSCollector.ListBuckets()
	log.Debugf("获取到 %d 条 OSS Bucket 信息 (Obtained %d OSS Bucket information)", len(Buckets), len(Buckets))

	Objects, _ := OSSCollector.ListObjects("all", ossLsObjectNumber)
	ACL := OSSCollector.GetBucketACL()

	var num = 0
	for _, o := range Buckets {
		if region == "all" {
			num = len(Buckets)
		} else {
			if region == o.Region {
				num = num + 1
			}
		}
	}
	var data = make([][]string, num)
	num = 0
	for i, o := range Buckets {
		if region == "all" {
			SN := strconv.Itoa(i + 1)
			ObjectNumber := strconv.Itoa(Objects[i].ObjectNumber)
			ObjectSize := formatFileSize(Objects[i].ObjectSize)
			BucketACL := ACL[i].Acl
			BucketURL := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", o.Name, o.Region)
			data[i] = []string{SN, o.Name, BucketACL, ObjectNumber, ObjectSize, o.Region, BucketURL}
		} else {
			if region == o.Region {
				ObjectNumber := strconv.Itoa(Objects[i].ObjectNumber)
				ObjectSize := formatFileSize(Objects[i].ObjectSize)
				BucketACL := ACL[i].Acl
				BucketURL := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", o.Name, o.Region)
				data[num] = []string{strconv.Itoa(num + 1), o.Name, BucketACL, ObjectNumber, ObjectSize, o.Region, BucketURL}
				num = num + 1
			}
		}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("没发现存储桶 (No Buckets Found)")
	} else {
		Caption := "OSS 资源 (OSS resources)"
		cloud.PrintTable(td, Caption)
	}
	cmdutil.WriteCacheFile(td, "alibaba", "oss", region, "all")
	util.WriteTimestamp(TimestampType)
}

func PrintBucketsListHistory(region string) {
	cmdutil.PrintOSSCacheFile(header, region, "alibaba", "OSS")
}

func PrintBucketsList(region string, lsFlushCache bool, ossLsObjectNumber string) {
	if lsFlushCache {
		PrintBucketsListRealTime(region, ossLsObjectNumber)
	} else {
		oldTimestamp := util.ReadTimestamp(TimestampType)
		if oldTimestamp == 0 {
			PrintBucketsListRealTime(region, ossLsObjectNumber)
		} else if util.IsFlushCache(oldTimestamp) {
			PrintBucketsListRealTime(region, ossLsObjectNumber)
		} else {
			util.TimeDifference(oldTimestamp)
			PrintBucketsListHistory(region)
		}
	}
}

func ReturnBucketList() []string {
	var (
		buckets  []string
		ossCache []pubutil.OSSCache
	)
	OSSCollector := &OSSCollector{}
	ossCache = cmdutil.ReadOSSCache("alibaba")
	if len(ossCache) == 0 {
		BucketsList, _ := OSSCollector.ListBuckets()
		for _, v := range BucketsList {
			buckets = append(buckets, v.Name)
		}
	} else {
		for _, v := range ossCache {
			buckets = append(buckets, v.Name)
		}
	}
	return buckets
}
