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

func ReturnOSSCacheFile() string {
	config := GetAliCredential()
	var ossCacheFile string
	AccessKeyId := config.AccessKeyId
	if AccessKeyId == "" {
		ossCacheFile = ""
	} else {
		ossCacheFile = ReturnCacheDict() + "/" + AccessKeyId[len(AccessKeyId)-6:] + "_oss.json"
	}
	return ossCacheFile
}

func ReturnECSCacheFile() string {
	config := GetAliCredential()
	var ecsCacheFile string
	AccessKeyId := config.AccessKeyId
	if AccessKeyId == "" {
		ecsCacheFile = ""
	} else {
		ecsCacheFile = ReturnCacheDict() + "/" + AccessKeyId[len(AccessKeyId)-6:] + "_ecs.json"
	}
	return ecsCacheFile
}

func ReturnRDSCacheFile() string {
	config := GetAliCredential()
	var rdsCacheFile string
	AccessKeyId := config.AccessKeyId
	if AccessKeyId == "" {
		rdsCacheFile = ""
	} else {
		rdsCacheFile = ReturnCacheDict() + "/" + AccessKeyId[len(AccessKeyId)-6:] + "_rds.json"
	}
	return rdsCacheFile
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

func ReadCacheFile(filePath string) [][]string {
	if !FileExists(filePath) {
		log.Debugf("%s 文件不存在 (%s file does not exist)", filePath, filePath)
		if filePath == ReturnOSSCacheFile() {
			log.Warnln("需要先使用 [cf oss ls] 命令获取 OSS 资源 (You need to use the [cf oss ls] command to get the OSS resources first)")
		} else if filePath == ReturnECSCacheFile() {
			log.Warnln("需要先使用 [cf ecs ls] 命令获取 ECS 资源 (You need to use the [cf ecs ls] command to get the ECS resources first)")
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

func PrintOSSCacheFile(filePath string, header []string, region string) {
	data := ReadCacheFile(filePath)
	if region == "all" {
		PrintTable(data, header, "OSS")
	} else {
		var dataRegion [][]string
		for _, i := range data {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, "OSS")
	}
}

func PrintECSCacheFile(filePath string, header []string, region string, specifiedInstanceID string) {
	data := ReadCacheFile(filePath)
	switch {
	case region == "all" && specifiedInstanceID == "all":
		PrintTable(data, header, "ECS")
	case region != "all" && specifiedInstanceID == "all":
		var dataRegion [][]string
		for _, i := range data {
			if i[8] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, "ECS")
	case region == "all" && specifiedInstanceID != "all":
		var dataSpecifiedInstanceID [][]string
		for _, i := range data {
			if i[1] == specifiedInstanceID {
				dataSpecifiedInstanceID = append(dataSpecifiedInstanceID, i)
			}
		}
		PrintTable(dataSpecifiedInstanceID, header, "ECS")
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
		PrintTable(dataSpecifiedInstanceID, header, "ECS")
	}
}

func PrintRDSCacheFile(filePath string, header []string, region string, specifiedDBInstanceID string, engine string) {
	data := ReadCacheFile(filePath)
	switch {
	case region == "all" && specifiedDBInstanceID == "all" && engine == "all":
		PrintTable(data, header, "RDS")
	case region == "all" && specifiedDBInstanceID == "all" && engine != "all":
		var dataEngine [][]string
		for _, i := range data {
			if i[2] == engine {
				dataEngine = append(dataEngine, i)
			}
		}
		PrintTable(dataEngine, header, "RDS")
	case region == "all" && specifiedDBInstanceID != "all" && engine == "all":
		var dataSpecifiedDBInstanceID [][]string
		for _, i := range data {
			if i[1] == specifiedDBInstanceID {
				dataSpecifiedDBInstanceID = append(dataSpecifiedDBInstanceID, i)
			}
		}
		PrintTable(dataSpecifiedDBInstanceID, header, "RDS")
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
		PrintTable(dataSpecifiedDBInstanceID, header, "RDS")
	case region != "all" && specifiedDBInstanceID == "all" && engine == "all":
		var dataRegion [][]string
		for _, i := range data {
			if i[5] == region {
				dataRegion = append(dataRegion, i)
			}
		}
		PrintTable(dataRegion, header, "RDS")
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
		PrintTable(dataRegion, header, "RDS")
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
		PrintTable(dataRegion, header, "RDS")
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
		PrintTable(dataRegion, header, "RDS")
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
