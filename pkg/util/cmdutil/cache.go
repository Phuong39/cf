package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/teamssix/cf/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
)

func ReturnCacheDict() string {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	cacheDict := home + "/cache"
	return cacheDict
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

func createCacheDict() {
	cacheDict := ReturnCacheDict()
	if FileExists(cacheDict) == false {
		log.Traceln("创建缓存目录 (Create cache directory): " + cacheDict)
		err := os.MkdirAll(cacheDict, 0700)
		util.HandleErr(err)
	}
}

func WriteCacheFile(td cloud.TableData, filePath string) {
	log.Debugln("写入数据到缓存文件 (Write data to a cache file): " + filePath)
	filePtr, err := os.Create(filePath)
	util.HandleErr(err)
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(td.Body)
	util.HandleErr(err)
}

func ReadCacheFile(filePath string, provider string, resourceType string) [][]string {
	if !FileExists(filePath) {
		log.Debugf("%s 文件不存在 (%s file does not exist)", filePath, filePath)
		if filePath == ReturnCacheFile(provider, resourceType) {
			log.Warnf("需要先使用 cf 获取 %s 资源 (You need to use the cf to get the %s resources first)", resourceType, resourceType)
		}
		os.Exit(0)
	}
	log.Debugln("读取文件 (read file): " + filePath)
	filePtr, err := os.Open(filePath)
	util.HandleErr(err)
	defer filePtr.Close()
	var data [][]string
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&data)
	util.HandleErr(err)
	return data
}

func PrintOSSCacheFile(filePath string, header []string, region string, provider string, resourceType string) {
	data := ReadCacheFile(filePath, provider, resourceType)
	if region == "all" {
		PrintTable(data, header, resourceType)
	} else {
		var dataRegion [][]string
		for _, i := range data {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, resourceType)
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
