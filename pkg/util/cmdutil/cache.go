package cmdutil

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
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
	if serviceType == "lh" {
		serviceType = "ecs"
	}
	if len(td.Body) == 0 {
		if serviceType == "oss" {
			database.DeleteOSSCache(AccessKeyId)
		} else if serviceType == "ecs" {
			database.DeleteECSCache(AccessKeyId)
		} else if serviceType == "rds" {
			database.DeleteRDSCache(AccessKeyId)
		}
	} else if region == "all" && id == "all" {
		log.Debugln("写入数据到缓存数据库 (Write data to a cache database)")
		switch {
		case serviceType == "oss":
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
		case serviceType == "ecs":
			var ECSCacheList []pubutil.ECSCache
			for _, v := range td.Body {
				ECSCache := pubutil.ECSCache{
					AccessKeyId:      AccessKeyId,
					SN:               v[0],
					InstanceId:       v[1],
					InstanceName:     v[2],
					OSName:           v[3],
					OSType:           v[4],
					Status:           v[5],
					PrivateIpAddress: v[6],
					PublicIpAddress:  v[7],
					RegionId:         v[8],
				}
				ECSCacheList = append(ECSCacheList, ECSCache)
			}
			database.InsertECSCache(ECSCacheList)
		case serviceType == "rds":
			var RDSCacheList []pubutil.RDSCache
			for _, v := range td.Body {
				RDSCache := pubutil.RDSCache{
					AccessKeyId:      AccessKeyId,
					SN:               v[0],
					DBInstanceId:     v[1],
					Engine:           v[2],
					EngineVersion:    v[3],
					DBInstanceStatus: v[4],
					RegionId:         v[5],
				}
				RDSCacheList = append(RDSCacheList, RDSCache)
			}
			database.InsertRDSCache(RDSCacheList)
		}
	} else {
		log.Debugln("由于数据不是全部数据，所以不写入缓存文件 (Since the data is not all data, it is not written to the cache file)")
	}
}

func ReadOSSCache(provider string) []pubutil.OSSCache {
	return database.SelectOSSCache(provider)
}

func ReadECSCache(provider string) []pubutil.ECSCache {
	return database.SelectECSCache(provider)
}

func PrintOSSCacheFile(header []string, region string, provider string, resourceType string) {
	OSSCache := database.SelectOSSCache(provider)
	var data [][]string
	OSSCache = database.SelectOSSCacheFilter(provider, region)
	for _, v := range OSSCache {
		dataSingle := []string{v.SN, v.Name, v.BucketACL, v.ObjectNumber, v.ObjectSize, v.Region, v.BucketURL}
		data = append(data, dataSingle)
	}
	PrintTable(data, header, resourceType)
}

func PrintECSCacheFile(header []string, region string, specifiedInstanceID string, provider string, resourceType string, running bool) {
	ECSCache := database.SelectECSCache(provider)
	var data [][]string
	ECSCache = database.SelectEcsCacheFilter(provider, region, specifiedInstanceID, running)
	for _, v := range ECSCache {
		dataSingle := []string{v.SN, v.InstanceId, v.InstanceName, v.OSName, v.OSType, v.Status, v.PrivateIpAddress, v.PublicIpAddress, v.RegionId}
		data = append(data, dataSingle)
	}
	PrintTable(data, header, resourceType)
}

func PrintRDSCacheFile(header []string, region string, specifiedDBInstanceID string, engine string, provider string, resourceType string) {
	RDSCache := database.SelectRDSCache(provider)
	var data [][]string
	RDSCache = database.SelectRDSCacheFilter(provider, region, specifiedDBInstanceID, engine)
	for _, v := range RDSCache {
		dataSingle := []string{v.SN, v.DBInstanceId, v.Engine, v.EngineVersion, v.DBInstanceStatus, v.RegionId}
		data = append(data, dataSingle)
	}
	PrintTable(data, header, resourceType)
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
