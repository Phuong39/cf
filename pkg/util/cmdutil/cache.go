package cmdutil

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
	"githubu.com/teamssix/cf/pkg/cloud"
	"githubu.com/teamssix/cf/pkg/util"
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
		ossCacheFile = ReturnCacheDict() + "/" + AccessKeyId[:6] + "_oss.json"
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
		ecsCacheFile = ReturnCacheDict() + "/" + AccessKeyId[:6] + "_ecs.json"
	}
	return ecsCacheFile
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
	log.Debugln("写入数据到文件 (Write data to a file): " + filePath)
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
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
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
}

func PrintCacheFile(filePath string, header []string, region string, specifiedInstanceID string) {
	data := ReadCacheFile(filePath)
	if filePath == ReturnOSSCacheFile() {
		if region == "all" {
			var td = cloud.TableData{Header: header, Body: data}
			if len(data) == 0 {
				log.Info("没有存储桶 (No Bucket)")
			} else {
				Caption := "OSS 资源 (OSS resources)"
				cloud.PrintTable(td, Caption)
			}
		} else {
			var dataRegion [][]string
			for _, i := range data {
				if i[5] == region {
					dataRegion = append(dataRegion, i)
				}
			}
			var td = cloud.TableData{Header: header, Body: dataRegion}
			if len(dataRegion) == 0 {
				log.Info("该区域下没有存储桶 (No Bucket was found in the region)")
			} else {
				Caption := "OSS 资源 (OSS resources)"
				cloud.PrintTable(td, Caption)
			}
		}
	} else if filePath == ReturnECSCacheFile() {
		if region == "all" && specifiedInstanceID == "all" {
			var td = cloud.TableData{Header: header, Body: data}
			if len(data) == 0 {
				log.Info("未发现 ECS，可能是因为当前访问凭证权限不够 (No ECS found, Probably because the current Access Key do not have enough permissions)")
			} else {
				Caption := "ECS 资源 (ECS resources)"
				cloud.PrintTable(td, Caption)
			}
		} else if region != "all" && specifiedInstanceID == "all" {
			var dataRegion [][]string
			for _, i := range data {
				if i[8] == region {
					dataRegion = append(dataRegion, i)
				}
			}
			var td = cloud.TableData{Header: header, Body: dataRegion}
			if len(dataRegion) == 0 {
				log.Infof("在 %s 区域下未发现 ECS，可能是因为当前访问凭证权限不够 (No ECS was found in %s region, Probably because the current Access Key do not have enough permissions)", region, region)
			} else {
				Caption := "ECS 资源 (ECS resources)"
				cloud.PrintTable(td, Caption)
			}
		} else if region == "all" && specifiedInstanceID != "all" {
			var dataSpecifiedInstanceID [][]string
			for _, i := range data {
				if i[1] == specifiedInstanceID {
					dataSpecifiedInstanceID = append(dataSpecifiedInstanceID, i)
				}
			}
			var td = cloud.TableData{Header: header, Body: dataSpecifiedInstanceID}
			if len(dataSpecifiedInstanceID) == 0 {
				log.Infof("未发现实例 ID 为 %s 的 ECS，可能是因为当前访问凭证权限不够 (No ECS with %s instance ID found, Probably because the current Access Key do not have enough permissions)", specifiedInstanceID, specifiedInstanceID)
			} else {
				Caption := "ECS 资源 (ECS resources)"
				cloud.PrintTable(td, Caption)
			}
		} else {
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
			var td = cloud.TableData{Header: header, Body: dataSpecifiedInstanceID}
			if len(dataSpecifiedInstanceID) == 0 {
				log.Infof("在 %s 区域下未发现实例 ID 为 %s 的 ECS，可能是因为当前访问凭证权限不够 (No ECS with instance ID %s found in %s region, Probably because the current Access Key do not have enough permissions)", region, specifiedInstanceID, specifiedInstanceID, region)
			} else {
				Caption := "ECS 资源 (ECS resources)"
				cloud.PrintTable(td, Caption)
			}
		}

	}
}
