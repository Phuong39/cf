package alirds

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

var (
	RDSCacheFilePath = cmdutil.ReturnRDSCacheFile("alibaba")
	header           = []string{"序号 (SN)", "数据库 ID (DB ID)", "数据库类型 (DB Engine)", "数据库版本 (DB Engine Version)", "数据库状态 (DB Staus)", "区域 ID (Region ID)"}
)

type DBInstances struct {
	DBInstanceId     string
	Engine           string
	EngineVersion    string
	DBInstanceStatus string
	RegionId         string
}

func DescribeDBInstances(region string, running bool, specifiedDBInstanceID string, engine string) []DBInstances {
	var out []DBInstances
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	if running == true {
		request.DBInstanceStatus = "Running"
	}
	if specifiedDBInstanceID != "all" {
		request.DBInstanceId = specifiedDBInstanceID
	}
	if engine != "all" {
		request.Engine = engine
	}
	response, err := RDSClient(region).DescribeDBInstances(request)
	util.HandleErrNoExit(err)
	DBInstancesList := response.Items.DBInstance
	log.Tracef("正在 %s 区域中查找数据库实例 (Looking for DBInstances in the %s region)", region, region)
	if len(DBInstancesList) != 0 {
		log.Debugf("在 %s 区域下找到 %d 个数据库实例 (Found %d DBInstances in %s region)", region, len(DBInstancesList), len(DBInstancesList), region)
		for _, i := range DBInstancesList {
			obj := DBInstances{
				DBInstanceId:     i.DBInstanceId,
				Engine:           i.Engine,
				EngineVersion:    i.EngineVersion,
				RegionId:         i.RegionId,
				DBInstanceStatus: i.DBInstanceStatus,
			}
			out = append(out, obj)
		}
	}
	return out
}

func ReturnDBInstancesList(region string, running bool, specifiedDBInstanceID string, engine string) []DBInstances {
	var DBInstancesList []DBInstances
	var DBInstance []DBInstances
	if region == "all" {
		var RegionsList []string
		for _, i := range GetRDSRegions() {
			RegionsList = append(RegionsList, i.RegionId)
		}
		RegionsList = RemoveRepeatedElement(RegionsList)
		for _, j := range RegionsList {
			DBInstance = DescribeDBInstances(j, running, specifiedDBInstanceID, engine)
			for _, i := range DBInstance {
				DBInstancesList = append(DBInstancesList, i)
			}
		}
	} else {
		DBInstancesList = DescribeDBInstances(region, running, specifiedDBInstanceID, engine)
	}
	return DBInstancesList
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func PrintDBInstancesListRealTime(region string, running bool, specifiedDBInstanceID string, engine string) {
	DBInstancesList := ReturnDBInstancesList(region, running, specifiedDBInstanceID, engine)
	var data = make([][]string, len(DBInstancesList))
	for i, o := range DBInstancesList {
		SN := strconv.Itoa(i + 1)
		data[i] = []string{SN, o.DBInstanceId, o.Engine, o.EngineVersion, o.DBInstanceStatus, o.RegionId}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 RDS (No RDS found)")
		cmdutil.WriteCacheFile(td, RDSCacheFilePath)
	} else {
		Caption := "RDS 资源 (RDS resources)"
		cloud.PrintTable(td, Caption)
		cmdutil.WriteCacheFile(td, RDSCacheFilePath)
	}
}

func PrintDBInstancesListHistory(region string, running bool, specifiedDBInstanceID string, engine string) {
	if cmdutil.FileExists(RDSCacheFilePath) {
		cmdutil.PrintRDSCacheFile(RDSCacheFilePath, header, region, specifiedDBInstanceID, engine, "alibaba")
	} else {
		PrintDBInstancesListRealTime(region, running, specifiedDBInstanceID, engine)
	}
}

func PrintDBInstancesList(region string, running bool, specifiedDBInstanceID string, engine string, lsFlushCache bool) {
	if lsFlushCache {
		PrintDBInstancesListRealTime(region, running, specifiedDBInstanceID, engine)
	} else {
		PrintDBInstancesListHistory(region, running, specifiedDBInstanceID, engine)
	}
}
