package alioss

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
)

var (
	objectsHeader = []string{"序号 (SN)", "名称 (Key)", "大小 (Size)", "上次修改时间 (LastModified)"}
)

func PrintObjectsList(bucketName string) {
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
		Caption := "对象资源 (Objects resources)"
		cloud.PrintTable(td, Caption)
	}
}
