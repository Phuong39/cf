package alioss

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
)

var (
	objectsHeader = []string{"序号 (SN)", "名称 (Key)", "大小 (Size)", "上次修改时间 (LastModified)"}
)

func PrintObjectsList() {
	buckets := ReturnBucketList()
	if len(buckets) == 0 {
		log.Info("没发现存储桶 (No Buckets Found)")
	} else {
		var bucketName string
		buckets = append(buckets, "exit")
		prompt := &survey.Select{
			Message: "选择一个存储桶 (Choose a bucket): ",
			Options: buckets,
		}
		err := survey.AskOne(prompt, &bucketName)
		errutil.HandleErr(err)
		if bucketName == "exit" {
			os.Exit(0)
		} else {
			OSSCollector := &OSSCollector{}
			objectsSum, objects := OSSCollector.ListObjects(bucketName)
			objectSize := formatFileSize(objectsSum[0].ObjectSize)
			var objectsData = make([][]string, len(objects))
			for i, o := range objects {
				objectsData[i] = []string{strconv.Itoa(i + 1), o.Key, formatFileSize(o.Size), o.LastModified}
			}
			var td = cloud.TableData{Header: objectsHeader, Body: objectsData}
			if len(objectsData) == 0 {
				log.Info("没发现对象 (No Objects Found)")
			} else {
				log.Infof("对象合计大小 %s (Total object size %s)", objectSize, objectSize)
				if len(objectsData) > 100 {
					var comfirm bool
					objectsDataLen := strconv.Itoa(len(objectsData))
					prompt := &survey.Confirm{
						Message: "存储桶中有 " + objectsDataLen + " 个对象，确定要列出吗？(There are " + objectsDataLen + "objects in the bucket, are you sure you want to list?)",
						Default: false,
					}
					err = survey.AskOne(prompt, &comfirm)
					errutil.HandleErr(err)
					if comfirm == false {
						os.Exit(0)
					}
				}
				Caption := "对象资源 (Objects resources)"
				cloud.PrintTable(td, Caption)
			}
		}
	}
}
