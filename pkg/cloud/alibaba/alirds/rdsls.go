package alirds

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/teamssix/cf/pkg/util/errutil"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

var (
	DescribeDBInstancesOut []DBInstances
	TimestampType          = util.ReturnTimestampType("alibaba", "rds")
	header                 = []string{"序号 (SN)", "数据库 ID (DB ID)", "数据库类型 (DB Engine)", "数据库版本 (DB Engine Version)", "数据库状态 (DB Staus)", "区域 ID (Region ID)"}
)

type DBInstances struct {
	DBInstanceId     string
	Engine           string
	EngineVersion    string
	DBInstanceStatus string
	RegionId         string
}

func DescribeDBInstances(region string, running bool, specifiedDBInstanceID string, engine string, NextToken string) ([]DBInstances, error) {
	request := rds.CreateDescribeDBInstancesRequest()
	request.PageSize = requests.NewInteger(100)
	request.Scheme = "https"
	if NextToken != "" {
		request.NextToken = NextToken
	}
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
	errutil.HandleErrNoExit(err)
	DBInstancesList := response.Items.DBInstance
	log.Infof("正在 %s 区域中查找数据库实例 (Looking for DBInstances in the %s region)", region, region)
	if len(DBInstancesList) != 0 {
		log.Infof("在 %s 区域下找到 %d 个数据库实例 (Found %d DBInstances in %s region)", region, len(DBInstancesList), len(DBInstancesList), region)
		for _, i := range DBInstancesList {
			obj := DBInstances{
				DBInstanceId:     i.DBInstanceId,
				Engine:           i.Engine,
				EngineVersion:    i.EngineVersion,
				RegionId:         i.RegionId,
				DBInstanceStatus: i.DBInstanceStatus,
			}
			DescribeDBInstancesOut = append(DescribeDBInstancesOut, obj)
		}
	}
	NextToken = response.NextToken
	if NextToken != "" {
		log.Tracef("Next Token: %s", NextToken)
		_, _ = DescribeDBInstances(region, running, specifiedDBInstanceID, engine, NextToken)
	}
	return DescribeDBInstancesOut, err
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
			DBInstance, _ = DescribeDBInstances(j, running, specifiedDBInstanceID, engine, "")
			DescribeDBInstancesOut = nil
			for _, i := range DBInstance {
				DBInstancesList = append(DBInstancesList, i)
			}
		}
	} else {
		DBInstancesList, _ = DescribeDBInstances(region, running, specifiedDBInstanceID, engine, "")
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
	} else {
		Caption := "RDS 资源 (RDS resources)"
		cloud.PrintTable(td, Caption)
		util.WriteTimestamp(TimestampType)
	}
	cmdutil.WriteCacheFile(td, "alibaba", "rds", region, specifiedDBInstanceID)
}

func PrintDBInstancesListHistory(region string, running bool, specifiedDBInstanceID string, engine string) {
	cmdutil.PrintRDSCacheFile(header, region, specifiedDBInstanceID, engine, "alibaba", "RDS")
}

func PrintDBInstancesList(region string, running bool, specifiedDBInstanceID string, engine string, lsFlushCache bool) {
	if lsFlushCache {
		PrintDBInstancesListRealTime(region, running, specifiedDBInstanceID, engine)
	} else {
		oldTimestamp := util.ReadTimestamp(TimestampType)
		if oldTimestamp == 0 {
			PrintDBInstancesListRealTime(region, running, specifiedDBInstanceID, engine)
		} else if util.IsFlushCache(oldTimestamp) {
			PrintDBInstancesListRealTime(region, running, specifiedDBInstanceID, engine)
		} else {
			util.TimeDifference(oldTimestamp)
			PrintDBInstancesListHistory(region, running, specifiedDBInstanceID, engine)
		}
	}
}
