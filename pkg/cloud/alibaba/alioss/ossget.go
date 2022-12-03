package alioss

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/schollz/progressbar/v3"

	log "github.com/sirupsen/logrus"
)

func getObject(bucketName string, objectKey string, outputPath string, ossLsRegion string) {
	if objectKey[len(objectKey)-1:] == "/" {
		pubutil.CreateFolder(returnBucketFileName(outputPath, bucketName, objectKey))
	} else {
		log.Infof("正在下载 %s 存储桶里的 %s 对象 (Downloading %s objects from %s bucket)", bucketName, objectKey, bucketName, objectKey)
		var (
			objectSize int64
			region     string
		)
		OSSCollector := &OSSCollector{}
		Buckets, _ := OSSCollector.ListBuckets(bucketName, ossLsRegion)
		for _, v := range Buckets {
			if v.Name == bucketName {
				region = v.Region
			}
		}
		fd, body, oserr, outputFile := OSSCollector.ReturnBucket(bucketName, objectKey, outputPath, region)
		_, objects := OSSCollector.ListObjects(bucketName, "all", ossLsRegion)
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
}

func DownloadAllObjects(bucketName string, outputPath string, ossLsRegion string, ossDownloadNumber string) {
	var (
		objectKey  string
		region     string
		objectList []string
	)
	OSSCollector := &OSSCollector{}
	objectList = append(objectList, "all")
	Buckets, _ := OSSCollector.ListBuckets(bucketName, ossLsRegion)
	for _, v := range Buckets {
		if v.Name == bucketName {
			region = v.Region
		}
	}
	_, objects := OSSCollector.ListObjects(bucketName, ossDownloadNumber, ossLsRegion)
	if len(objects) == 0 {
		log.Warnf("在 %s 存储桶中没有发现对象 (No object found in %s bucket)", bucketName, bucketName)
	} else {
		for _, o := range objects {
			objectList = append(objectList, o.Key)
		}
		objectList = append(objectList, "exit")
		sort.Strings(objectList)
		prompt := &survey.Select{
			Message: "选择一个对象 (Choose a object): ",
			Options: objectList,
		}
		survey.AskOne(prompt, &objectKey)
		if objectKey == "all" {
			log.Infof("正在下载 %s 存储桶内的所有对象…… (Downloading all objects in bucket %s...)", bucketName, bucketName)
			bar := returnBar((int64(len(objectList) - 2)))
			for _, j := range objects {
				if j.Key[len(j.Key)-1:] == "/" {
					bar.Add(1)
					pubutil.CreateFolder(returnBucketFileName(outputPath, bucketName, j.Key))
				} else {
					bar.Add(1)
					fd, body, _, _ := OSSCollector.ReturnBucket(bucketName, j.Key, outputPath, region)
					io.Copy(fd, body)
					body.Close()
					defer fd.Close()
				}
			}
			log.Infof("对象已被保存到 %s 目录下 (The object has been saved to the %s directory)", outputPath, outputPath)
		} else if objectKey == "exit" {
			os.Exit(0)
		} else {
			if objectKey[len(objectKey)-1:] == "/" {
				pubutil.CreateFolder(returnBucketFileName(outputPath, bucketName, objectKey))
			} else {
				getObject(bucketName, objectKey, outputPath, ossLsRegion)
			}
		}
	}
}

func DownloadObjects(bucketName string, objectKey string, outputPath string, ossLsRegion string, ossDownloadNumber string) {
	if outputPath == "./result" {
		pubutil.CreateFolder("./result")
	}
	if bucketName == "all" {
		var bucketList []string
		buckets := ReturnBucketList(bucketName, ossLsRegion)
		if len(buckets) == 0 {
			log.Info("没发现存储桶 (No Buckets Found)")
		} else {
			bucketList = append(bucketList, "all")
			for _, v := range buckets {
				bucketList = append(bucketList, v)
			}
			bucketList = append(bucketList, "exit")
			var SelectBucketName string
			sort.Strings(bucketList)
			prompt := &survey.Select{
				Message: "选择一个存储桶 (Choose a bucket): ",
				Options: bucketList,
			}
			err := survey.AskOne(prompt, &SelectBucketName)
			errutil.HandleErr(err)
			if SelectBucketName == "all" {
				for _, v := range buckets {
					if objectKey == "all" {
						DownloadAllObjects(v, outputPath, ossLsRegion, ossDownloadNumber)
					} else {
						getObject(v, objectKey, outputPath, ossLsRegion)
					}
				}
			} else if SelectBucketName == "exit" {
				os.Exit(0)
			} else {
				if objectKey == "all" {
					DownloadAllObjects(SelectBucketName, outputPath, ossLsRegion, ossDownloadNumber)
				} else {
					getObject(SelectBucketName, objectKey, outputPath, ossLsRegion)
				}
			}
		}
	} else {
		if objectKey == "all" {
			DownloadAllObjects(bucketName, outputPath, ossLsRegion, ossDownloadNumber)
		} else {
			getObject(bucketName, objectKey, outputPath, ossLsRegion)
		}
	}
}

func (o *OSSCollector) ReturnBucket(bucketName string, objectKey string, outputPath string, region string) (*os.File, io.ReadCloser, error, string) {
	o.OSSClient(region)
	bucket, err := o.Client.Bucket(bucketName)
	errutil.HandleErr(err)
	outputFile := returnBucketFileName(outputPath, bucketName, objectKey)
	fd, oserr := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0660)
	errutil.HandleErr(oserr)
	body, err := bucket.GetObject(objectKey)
	errutil.HandleErr(err)
	return fd, body, oserr, outputFile
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

func returnBucketFileName(outputPath string, bucketName string, objectName string) string {
	outputBucketFile := filepath.Join(outputPath, bucketName)
	pubutil.CreateFolder(outputBucketFile)
	outputFileName := filepath.Join(outputBucketFile, objectName)
	return outputFileName
}
