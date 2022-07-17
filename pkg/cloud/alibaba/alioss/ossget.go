package alioss

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/schollz/progressbar/v3"

	log "github.com/sirupsen/logrus"

	"github.com/teamssix/cf/pkg/util"
)

func getObject(bucketName string, objectKey string, outputPath string) {
	log.Infof("正在下载 %s 存储桶里的 %s 对象 (Downloading %s objects from %s bucket)", bucketName, objectKey, bucketName, objectKey)
	var objectSize int64
	OSSCollector := &OSSCollector{}
	objects, fd, body, oserr, outputFile := OSSCollector.ReturnBucket(bucketName, objectKey, outputPath)
	for _, obj := range objects {
		if objectKey == obj.Key {
			objectSize = obj.Size
		}
	}
	bar := returnBar(objectSize)
	io.Copy(io.MultiWriter(fd, bar), body)
	body.Close()
	defer fd.Close()
	if oserr == nil {
		log.Infof("对象已被保存到 %s (The object has been saved to %s)", outputFile, outputFile)
	}
}

func DownloadAllObjects(bucketName string, outputPath string) {
	var (
		objectKey  string
		objectList []string
	)
	OSSCollector := &OSSCollector{}
	objectList = append(objectList, "all")
	_, objects := OSSCollector.ListObjects(bucketName)
	for _, o := range objects {
		objectList = append(objectList, o.Key)
	}
	prompt := &survey.Select{
		Message: "选择一个对象 (Choose a object): ",
		Options: objectList,
	}
	survey.AskOne(prompt, &objectKey)
	if objectKey == "all" {
		bar := returnBar((int64(len(objectList) - 1)))
		for _, j := range objects {
			bar.Add(1)
			_, fd, body, _, _ := OSSCollector.ReturnBucket(bucketName, j.Key, outputPath)
			io.Copy(fd, body)
			body.Close()
			defer fd.Close()
		}

		log.Infof("对象已被保存到 %s 目录下 (The object has been saved to the %s directory)", outputPath, outputPath)
	} else {
		getObject(bucketName, objectKey, outputPath)
	}
}

func DownloadObjects(bucketName string, objectKey string, outputPath string) {
	if objectKey == "all" {
		DownloadAllObjects(bucketName, outputPath)
	} else {
		getObject(bucketName, objectKey, outputPath)
	}
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func returnOutputFile(objectKey string, outputPath string) string {
	var (
		outputFile string
		objectName string
	)
	if strings.Contains(objectKey, "/") {
		objectKeyList := strings.Split(objectKey, "/")
		objectName = objectKeyList[len(objectKeyList)-1]
	} else {
		objectName = objectKey
	}

	if IsDir(outputPath) {
		if outputPath[len(outputPath)-1:] != "/" {
			outputFile = outputPath + "/" + objectName
		} else {
			outputFile = outputPath + objectName
		}
	} else {
		outputFile = outputPath
	}
	log.Debugf("下载保存路径为 %s (The save path is %s)", outputFile, outputFile)
	return outputFile
}

func (o *OSSCollector) ReturnBucket(bucketName string, objectKey string, outputPath string) ([]objectContents, *os.File, io.ReadCloser, error, string) {
	var (
		err    error
		region string
	)
	Buckets, _ := o.ListBuckets()
	for _, obj := range Buckets {
		if obj.Name == bucketName {
			region = obj.Region
		}
	}
	if region == "" {
		log.Errorln("未找到该 Bucket (This Bucket not found)")
		os.Exit(0)
	}
	o.OSSClient(region)
	bucket, err := o.Client.Bucket(bucketName)
	util.HandleErr(err)
	outputFile := returnOutputFile(objectKey, outputPath)
	fd, oserr := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0660)
	util.HandleErr(oserr)
	body, err := bucket.GetObject(objectKey)
	util.HandleErr(err)
	_, objects := o.ListObjects(bucketName)
	return objects, fd, body, oserr, outputFile
}

func returnBar(replen int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(replen,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}
