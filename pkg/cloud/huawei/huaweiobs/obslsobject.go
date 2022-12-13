package huaweiobs

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/errutil"
	"strconv"
)

var (
	objectsKey  []string
	objectsSize []int64
)

func ListObjects(bucket string, region string, obsLsObjectNumber string, NextMarker string) {
	var (
		obsLsObjectNumberInt int
		MaxKeys              int
		err                  error
		IsTruncated          bool
		input                *obs.ListObjectsInput
	)
	if obsLsObjectNumber != "all" {
		obsLsObjectNumberInt, err = strconv.Atoi(obsLsObjectNumber)
		errutil.HandleErr(err)
	}
	if obsLsObjectNumberInt == 0 || obsLsObjectNumberInt > 1000 {
		MaxKeys = 1000
	} else {
		MaxKeys = obsLsObjectNumberInt
	}
	if NextMarker == "" {
		input = &obs.ListObjectsInput{
			Bucket: bucket,
			ListObjsInput: obs.ListObjsInput{
				MaxKeys: MaxKeys,
			},
		}
	} else {
		input = &obs.ListObjectsInput{
			Bucket: bucket,
			Marker: NextMarker,
			ListObjsInput: obs.ListObjsInput{
				MaxKeys: MaxKeys,
			},
		}
	}
	result, err := obsClient("all").ListObjects(input)
	errutil.HandleErr(err)
	IsTruncated = result.IsTruncated
	for _, v := range result.Contents {
		objectsKey = append(objectsKey, v.Key)
		objectsSize = append(objectsSize, v.Size)
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
	if obsLsObjectNumber != "all" {
		obsLsObjectNumberInt, err = strconv.Atoi(obsLsObjectNumber)
		errutil.HandleErr(err)
		obsLsObjectNumberInt := int(obsLsObjectNumberInt)
		if obsLsObjectNumberInt >= objectNum {
			IsTruncated = false
			objectNum = obsLsObjectNumberInt
			objectsKey = objectsKey[0:objectNum]
			objectsSize = objectsSize[0:objectNum]
		}
	}
	if IsTruncated {
		ListObjects(bucket, region, obsLsObjectNumber, result.NextMarker)
	}
}
