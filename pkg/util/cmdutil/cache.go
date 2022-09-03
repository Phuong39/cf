package cmdutil

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"os"
)

func ReturnCacheDict() string {
	home, err := pubutil.GetCFHomeDir()
	errutil.HandleErr(err)
	return home
}

func ReturnCacheFile(provider string, resourceType string) string {
	config := GetConfig(provider)
	var ossCacheFile string
	AccessKeyId := config.AccessKeyId
	if AccessKeyId == "" {
		ossCacheFile = ""
	} else {
		ossCacheFile = ReturnCacheDict() + "/" + AccessKeyId[len(AccessKeyId)-6:] + "_" + resourceType + ".json"
	}
	return ossCacheFile
}

func WriteCacheFile(td cloud.TableData, provider string, serviceType string, region string, id string) {
	AccessKeyId := GetConfig(provider).AccessKeyId
	if len(td.Body) == 0 {
		database.DeleteOSSCache(AccessKeyId)
	} else if region == "all" && id == "all" {
		log.Debugln("写入数据到缓存数据库 (Write data to a cache database)")
		if provider == "alibaba" {
			if serviceType == "oss" {
				var OSSCacheList []pubutil.OSSCache
				for _, v := range td.Body {
					OSSCache := pubutil.OSSCache{
						AccessKeyId:  AccessKeyId,
						SN:           v[0],
						Name:         v[1],
						BucketACL:    v[2],
						ObjectNumber: v[3],
						ObjectSize:   v[4],
						Region:       v[5],
						BucketURL:    v[6],
					}
					OSSCacheList = append(OSSCacheList, OSSCache)
				}
				database.InsertOSSCache(OSSCacheList)

			}
		}
	} else {
		log.Debugln("由于数据不是全部数据，所以不写入缓存文件 (Since the data is not all data, it is not written to the cache file)")
	}
}

func ReadCacheFile(filePath string, provider string, resourceType string) [][]string {
	if !pubutil.FileExists(filePath) {
		log.Debugf("%s 文件不存在 (%s file does not exist)", filePath, filePath)
		if filePath == ReturnCacheFile(provider, resourceType) {
			log.Warnf("需要先使用 cf 获取 %s 资源 (You need to use the cf to get the %s resources first)", resourceType, resourceType)
		}
		os.Exit(0)
	}
	log.Debugln("读取文件 (read file): " + filePath)
	filePtr, err := os.Open(filePath)
	errutil.HandleErr(err)
	defer filePtr.Close()
	var data [][]string
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&data)
	errutil.HandleErr(err)
	return data
}

func PrintOSSCacheFile(header []string, region string, provider string, resourceType string) {
	OSSCache := database.SelectOSSCache(provider)
	var data [][]string
	if len(OSSCache) > 0 {
		if region == "all" {
			for _, v := range OSSCache {
				dataSingle := []string{v.SN, v.Name, v.BucketACL, v.ObjectNumber, v.ObjectSize, v.Region, v.BucketURL}
				data = append(data, dataSingle)
			}
		} else {
			for _, v := range OSSCache {
				if v.Region == region {
					dataSingle := []string{v.SN, v.Name, v.BucketACL, v.ObjectNumber, v.ObjectSize, v.Region, v.BucketURL}
					data = append(data, dataSingle)
				}
			}
		}
		PrintTable(data, header, resourceType)
	} else {
		PrintTable(data, header, resourceType)
	}
}

func PrintECSCacheFile(filePath string, header []string, region string, specifiedInstanceID string, provider string, resourceType string) {
	data := ReadCacheFile(filePath, provider, resourceType)
	switch {
	case region == "all" && specifiedInstanceID == "all":
		PrintTable(data, header, resourceType)
	case region != "all" && specifiedInstanceID == "all":
		var dataRegion [][]string
		for _, i := range data {
			if i[8] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, resourceType)
	case region == "all" && specifiedInstanceID != "all":
		var dataSpecifiedInstanceID [][]string
		for _, i := range data {
			if i[1] == specifiedInstanceID {
				dataSpecifiedInstanceID = append(dataSpecifiedInstanceID, i)
			}
		}
		PrintTable(dataSpecifiedInstanceID, header, resourceType)
	case region != "all" && specifiedInstanceID != "all":
		var dataRegion [][]string
		for _, i := range data {
			if i[8] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		var dataSpecifiedInstanceID [][]string
		for _, i := range dataRegion {
			if i[1] == specifiedInstanceID {
				dataSpecifiedInstanceID = append(dataSpecifiedInstanceID, i)
			}
		}
		PrintTable(dataSpecifiedInstanceID, header, resourceType)
	}
}

func PrintRDSCacheFile(filePath string, header []string, region string, specifiedDBInstanceID string, engine string, provider string, resourceType string) {
	data := ReadCacheFile(filePath, provider, resourceType)
	switch {
	case region == "all" && specifiedDBInstanceID == "all" && engine == "all":
		PrintTable(data, header, resourceType)
	case region == "all" && specifiedDBInstanceID == "all" && engine != "all":
		var dataEngine [][]string
		for _, i := range data {
			if i[2] == engine {
				dataEngine = append(dataEngine, i)
			}
		}
		PrintTable(dataEngine, header, resourceType)
	case region == "all" && specifiedDBInstanceID != "all" && engine == "all":
		var dataSpecifiedDBInstanceID [][]string
		for _, i := range data {
			if i[1] == specifiedDBInstanceID {
				dataSpecifiedDBInstanceID = append(dataSpecifiedDBInstanceID, i)
			}
		}
		PrintTable(dataSpecifiedDBInstanceID, header, resourceType)
	case region == "all" && specifiedDBInstanceID != "all" && engine != "all":
		var dataEngine [][]string
		for _, i := range data {
			if i[2] == engine {
				dataEngine = append(dataEngine, i)
			}
		}
		var dataSpecifiedDBInstanceID [][]string
		for _, i := range dataEngine {
			if i[1] == specifiedDBInstanceID {
				dataSpecifiedDBInstanceID = append(dataSpecifiedDBInstanceID, i)
			}
		}
		PrintTable(dataSpecifiedDBInstanceID, header, resourceType)
	case region != "all" && specifiedDBInstanceID == "all" && engine == "all":
		var dataRegion [][]string
		for _, i := range data {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, resourceType)
	case region != "all" && specifiedDBInstanceID == "all" && engine != "all":
		var dataEngine [][]string
		for _, i := range data {
			if i[2] == engine {
				dataEngine = append(dataEngine, i)
			}
		}
		var dataRegion [][]string
		for _, i := range dataEngine {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, resourceType)
	case region != "all" && specifiedDBInstanceID != "all" && engine == "all":
		var dataSpecifiedDBInstanceID [][]string
		for _, i := range data {
			if i[1] == specifiedDBInstanceID {
				dataSpecifiedDBInstanceID = append(dataSpecifiedDBInstanceID, i)
			}
		}
		var dataRegion [][]string
		for _, i := range dataSpecifiedDBInstanceID {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, resourceType)
	case region != "all" && specifiedDBInstanceID != "all" && engine != "all":
		var dataEngine [][]string
		for _, i := range data {
			if i[2] == engine {
				dataEngine = append(dataEngine, i)
			}
		}
		var dataSpecifiedDBInstanceID [][]string
		for _, i := range dataEngine {
			if i[1] == specifiedDBInstanceID {
				dataSpecifiedDBInstanceID = append(dataSpecifiedDBInstanceID, i)
			}
		}
		var dataRegion [][]string
		for _, i := range dataSpecifiedDBInstanceID {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, resourceType)
	}
}

func PrintTable(data [][]string, header []string, resourceType string) {
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info(fmt.Sprintf("未发现 %s (No %s found)", resourceType, resourceType))
	} else {
		Caption := fmt.Sprintf("%s 资源 (%s resources)", resourceType, resourceType)
		cloud.PrintTable(td, Caption)
	}
}

func PrintSSHCacheFile(filePath string, header []string, provide string, resourceType string) {
	data := ReadCacheFile(filePath, provide, resourceType)
	PrintTable(data, header, resourceType)
}
