package tencentcvm

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var (
	CVMCacheFilePath = cmdutil.ReturnCacheFile("tencent", "CVM")
	header           = []string{"序号 (SN)", "实例 ID (Instance ID)", "实例名称 (Instance Name)", "系统名称 (OS Name)", "系统类型 (OS Type)", "状态 (Status)", "私有 IP (Private Ip Address)", "公网 IP (Public Ip Address)", "区域 ID (Region ID)", "绑定的安全组 (Security Group Id)"}
)

type Instances struct {
	InstanceId       string
	InstanceName     string
	OSName           string
	OSType           string
	Status           string
	PrivateIpAddress string
	PublicIpAddress  string
	RegionId         string
	SecurityGroupIds string
}

func DescribeInstances(region string, running bool, SpecifiedInstanceID string) []Instances {
	var out []Instances
	request := cvm.NewDescribeInstancesRequest()
	request.SetScheme("https")
	if running {
		request.Filters = []*cvm.Filter{
			{
				Name:   common.StringPtr("instance-state"),
				Values: common.StringPtrs([]string{"RUNNING"}),
			},
		}
	}
	if SpecifiedInstanceID != "all" {
		request.InstanceIds = common.StringPtrs([]string{SpecifiedInstanceID})
	}
	response, err := CVMClient(region).DescribeInstances(request)
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
			SecurityGroupIdList  []string
		)
		for _, v := range InstancesList {
			for _, m := range v.PrivateIpAddresses {
				PrivateIpAddressList = append(PrivateIpAddressList, *m)
			}
			for _, m := range v.PublicIpAddresses {
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
			for _, i := range v.SecurityGroupIds {
				SecurityGroupIdList = append(SecurityGroupIdList, *i)
			}
			b, err := json.Marshal(SecurityGroupIdList)
			util.HandleErr(err)
			SecurityGroupIds := string(b)
			obj := Instances{
				InstanceId:       *v.InstanceId,
				InstanceName:     *v.InstanceName,
				OSName:           *v.OsName,
				OSType:           *v.InstanceType,
				Status:           *v.InstanceState,
				PrivateIpAddress: PrivateIpAddress,
				PublicIpAddress:  PublicIpAddress,
				RegionId:         *v.Placement.Zone,
				SecurityGroupIds: SecurityGroupIds,
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
		for _, j := range GetCVMRegions() {
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
		data[i] = []string{SN, o.InstanceId, o.InstanceName, o.OSName, o.OSType, o.Status, o.PrivateIpAddress, o.PublicIpAddress, o.RegionId, o.SecurityGroupIds}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 CVM 实例 (No CVM instances found)")
		cmdutil.WriteCacheFile(td, CVMCacheFilePath, region, specifiedInstanceID)
	} else {
		Caption := "CVM 资源 (CVM resources)"
		cloud.PrintTable(td, Caption)
		cmdutil.WriteCacheFile(td, CVMCacheFilePath, region, specifiedInstanceID)
	}
}

func PrintInstancesListHistory(region string, running bool, specifiedInstanceID string) {
	if cmdutil.FileExists(CVMCacheFilePath) {
		cmdutil.PrintECSCacheFile(CVMCacheFilePath, header, region, specifiedInstanceID, "tencent", "CVM")
	} else {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	}
}

func PrintInstancesList(region string, running bool, specifiedInstanceID string, cvmFlushCache bool) {
	if cvmFlushCache {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	} else {
		PrintInstancesListHistory(region, running, specifiedInstanceID)
	}
}
