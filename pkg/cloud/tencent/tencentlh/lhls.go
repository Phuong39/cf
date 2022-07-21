package tencentlh

import (
	"encoding/json"
	lh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	cwp "github.com/teamssix/cf/pkg/cloud/tencent/tencentcwp"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var (
	LHCacheFilePath = cmdutil.ReturnCacheFile("tencent", "LH")
	header          = []string{"序号 (SN)", "实例ID (Instance ID)", "云镜ID (UUID)", "云镜状态 (CWP Status)", "实例名称 (Instance Name)", "系统名称 (OS Name)", "状态 (Status)", "私有 IP (Private Ip Address)", "公网 IP (Public Ip Address)", "区域 ID (Region ID)"}
)

type Instances struct {
	InstanceId       string
	InstanceName     string
	OSName           string
	Status           string
	PrivateIpAddress string
	PublicIpAddress  string
	RegionId         string
	CWPStatus        string
	UUID             string
}

func DescribeInstances(region string, running bool, SpecifiedInstanceID string) []Instances {
	var out []Instances
	request := lh.NewDescribeInstancesRequest()
	if running {
		request.Filters = []*lh.Filter{
			{
				Name:   common.StringPtr("instance-state"),
				Values: common.StringPtrs([]string{"RUNNING"}),
			},
		}
	}
	if SpecifiedInstanceID != "all" {
		request.InstanceIds = common.StringPtrs([]string{SpecifiedInstanceID})
	}
	response, err := LHClient(region).DescribeInstances(request)
	util.HandleErr(err)
	InstancesList := response.Response.InstanceSet
	log.Tracef("正在 %s 区域中查找实例 (Looking for instances in the %s region)", region, region)
	if len(InstancesList) != 0 {
		log.Debugf("在 %s 区域下找到 %d 个实例 (Found %d instances in %s region)", region, len(InstancesList), len(InstancesList), region)
		var (
			PrivateIpAddressList []string
			PublicIpAddressList  []string
			PrivateIpAddress     string
			PublicIpAddress      string
		)
		for _, v := range InstancesList {
			for _, m := range v.PrivateAddresses {
				PrivateIpAddressList = append(PrivateIpAddressList, *m)
			}
			for _, m := range v.PublicAddresses {
				PublicIpAddressList = append(PublicIpAddressList, *m)
			}
			a, _ := json.Marshal(PrivateIpAddressList)
			if len(PrivateIpAddressList) == 1 {
				PrivateIpAddress = PrivateIpAddressList[0]
			} else {
				PrivateIpAddress = string(a)
			}
			b, _ := json.Marshal(PublicIpAddressList)
			if len(PublicIpAddressList) == 1 {
				PublicIpAddress = PublicIpAddressList[0]
			} else {
				PublicIpAddress = string(b)
			}

			util.HandleErr(err)
			CWPStatus, CWPUUID := cwp.DescribeMachineCWPStatus("LH", *v.Uuid)

			obj := Instances{
				InstanceId:       *v.InstanceId,
				InstanceName:     *v.InstanceName,
				OSName:           *v.OsName,
				Status:           *v.InstanceState,
				PrivateIpAddress: PrivateIpAddress,
				PublicIpAddress:  PublicIpAddress,
				RegionId:         *v.Zone,
				CWPStatus:        *CWPStatus,
				UUID:             *CWPUUID,
			}
			out = append(out, obj)
		}
	}
	return out
}

func ReturnInstancesList(region string, running bool, specifiedInstanceID string) []Instances {
	var InstancesList []Instances
	var Instance []Instances
	if region == "all" {
		for _, j := range GetLHRegions() {
			region := *j.Region
			Instance = DescribeInstances(region, running, specifiedInstanceID)
			for _, i := range Instance {
				InstancesList = append(InstancesList, i)
			}
		}
	} else {
		InstancesList = DescribeInstances(region, running, specifiedInstanceID)
	}
	return InstancesList
}

func PrintInstancesListRealTime(region string, running bool, specifiedInstanceID string) {
	InstancesList := ReturnInstancesList(region, running, specifiedInstanceID)
	var data = make([][]string, len(InstancesList))
	for i, o := range InstancesList {
		SN := strconv.Itoa(i + 1)
		data[i] = []string{SN, o.InstanceId, o.UUID, o.CWPStatus, o.InstanceName, o.OSName, o.Status, o.PrivateIpAddress, o.PublicIpAddress, o.RegionId}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 LH 实例 (No LH instances found)")
		cmdutil.WriteCacheFile(td, LHCacheFilePath, region, specifiedInstanceID)
	} else {
		Caption := "LH 资源 (LH resources)"
		cloud.PrintTable(td, Caption)
		cmdutil.WriteCacheFile(td, LHCacheFilePath, region, specifiedInstanceID)
	}
}

func PrintInstancesListHistory(region string, running bool, specifiedInstanceID string) {
	if cmdutil.FileExists(LHCacheFilePath) {
		cmdutil.PrintECSCacheFile(LHCacheFilePath, header, region, specifiedInstanceID, "tencent", "LH")
	} else {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	}
}

func PrintInstancesList(region string, running bool, specifiedInstanceID string, lhFlushCache bool) {
	if lhFlushCache {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	} else {
		PrintInstancesListHistory(region, running, specifiedInstanceID)
	}
}
