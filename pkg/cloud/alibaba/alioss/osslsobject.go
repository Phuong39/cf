package alioss

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"

	log "github.com/sirupsen/logrus"
)

var (
	objectsHeader = []string{"序号 (SN)", "名称 (Key)", "大小 (Size)", "上次修改时间 (LastModified)"}
)

func PrintObjectsList(ossLsObjectNumber string, ossLsBucket string, ossLsRegion string) {
	buckets := ReturnBucketList(ossLsBucket, ossLsRegion)
	if len(buckets) == 0 {
		log.Info("没发现存储桶 (No Buckets Found)")
	} else {
		var bucketName string
		buckets = append(buckets, "exit")
		sort.Strings(buckets)
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
			objectsSum, objects := OSSCollector.ListObjects(bucketName, ossLsObjectNumber, ossLsRegion)
			objectSize := pubutil.FormatFileSize(objectsSum[0].ObjectSize)
			var objectsData = make([][]string, len(objects))
			for i, o := range objects {
				objectsData[i] = []string{strconv.Itoa(i + 1), o.Key, pubutil.FormatFileSize(o.Size), o.LastModified}
			}
			if len(objectsData) == 0 {
				log.Info("没发现对象 (No Objects Found)")
			} else {
				log.Infof("对象合计大小 %s (Total object size %s)", objectSize, objectSize)
				if len(objectsData) > 100 {
					// 输出前 100 条结果
					objectsData100 := objectsData[0:100]
					printTd(objectsData100)
					objectsDataLen := strconv.Itoa(len(objectsData))
					log.Warnf("当前存储桶中有 %s 个对象，为了更好的输出效果，CF 只会列出前 100 个对象。(There are currently %s objects in the bucket, for better output, CF will only list the first 100 objects.)", objectsDataLen, objectsDataLen)

					// 超过 100 时，将输出结果写入到文件中去
					home, _ := pubutil.GetCFHomeDir()
					cacheFolder := filepath.Join(home, "output")
					pubutil.CreateFolder(cacheFolder)
					cacheFileName := filepath.Join(cacheFolder, "ossObjectList-"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")
					var cacheFileData []string
					cacheFileData = append(cacheFileData, "SN, Key, Size, LastModified\n")
					for _, v := range objectsData {
						cacheFileData = append(cacheFileData, v[0]+","+v[1]+","+v[2]+","+v[3]+"\n")
					}
					file, err := os.OpenFile(cacheFileName, os.O_WRONLY|os.O_CREATE, 0666)
					errutil.HandleErr(err)
					defer file.Close()
					w := bufio.NewWriter(file)
					for _, v := range cacheFileData {
						w.WriteString(v)
					}
					w.Flush()
					log.Infof("全部对象列表已保存到 %s 文件中，如果您想查看全部对象，可打开该文件进行查看。(The full list of objects has been saved to the %s file, so if you want to see all the objects, you can open this file to view them.)", cacheFileName, cacheFileName)
				} else {
					printTd(objectsData)
				}
			}
		}
	}
}

func printTd(objectsData [][]string) {
	var td = cloud.TableData{Header: objectsHeader, Body: objectsData}
	Caption := "对象资源 (Objects resources)"
	cloud.PrintTable(td, Caption)
}
