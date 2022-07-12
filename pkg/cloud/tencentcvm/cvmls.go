package tencentcvm

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"strconv"
)

var (
	CVMCacheFilePath = cmdutil.ReturnCVMCacheFile()
	header           = []string{"序号 (SN)", "实例 ID (Instance ID)", "系统名称 (OS Name)", "系统类型 (OS Type)", "状态 (Status)", "私有 IP (Private Ip Address)", "公网 IP (Public Ip Address)", "区域 ID (Region ID)"}
)

func DescribeInstances(region string, running bool, SpecifiedInstanceID string) []Instances {
	var out []Instances
	request := cvm.NewDescribeInstancesRequest()
	request.SetScheme("https")
	if region != "all" {
		//
		//log.Info("spid:", SpecifiedInstanceID)
	}
	response, err := CVMClient(region).DescribeInstances(request)
	util.HandleErr(err)
	InstancesList := response.Response.InstanceSet
	log.Tracef("正在 %s 区域中查找实例 (Looking for instances in the %s region)", region, region)
	if len(InstancesList) != 0 {
		log.Debugf("在 %s 区域下找到 %d 个实例 (Found %d instances in %s region)", region, len(InstancesList), len(InstancesList), region)
		var PrivateIpAddressList []string
		var PublicIpAddressList []string
		var PrivateIpAddress string
		var PublicIpAddress string
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
			obj := Instances{
				InstanceId:       *v.InstanceId,
				OSName:           *v.OsName,
				InstanceType:     *v.InstanceType,
				InstanceState:    *v.InstanceState,
				PrivateIpAddress: PrivateIpAddress,
				PublicIpAddress:  PublicIpAddress,
				Zone:             *v.Placement.Zone,
			}
			out = append(out, obj)
		}
	}
	//log.Warnln(out)
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
		data[i] = []string{SN, o.InstanceId, o.OSName, o.InstanceType, o.InstanceState, o.PrivateIpAddress, o.PublicIpAddress, o.Zone}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 CVM，可能是因为当前访问凭证权限不够 (No CVM found, Probably because the current Access Key do not have enough permissions)")
		cmdutil.WriteCacheFile(td, CVMCacheFilePath)
	} else {
		Caption := "CVM 资源 (CVM resources)"
		cloud.PrintTable(td, Caption)
		cmdutil.WriteCacheFile(td, CVMCacheFilePath)
	}
}

func PrintInstancesListHistory(region string, running bool, specifiedInstanceID string) {
	if cmdutil.FileExists(CVMCacheFilePath) {
		cmdutil.PrintCFCacheFile(CVMCacheFilePath, header, region, specifiedInstanceID, "tencent")
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
