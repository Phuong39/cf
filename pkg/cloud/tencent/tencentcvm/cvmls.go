package tencentcvm

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/teamssix/cf/pkg/util/errutil"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var (
	header       = []string{"序号 (SN)", "实例 ID (Instance ID)", "实例名称 (Instance Name)", "系统名称 (OS Name)", "系统类型 (OS Type)", "状态 (Status)", "私有 IP (Private IP)", "公网 IP (Public IP)", "区域 ID (Region ID)"}
	LinuxSet     = []string{"CentOS", "Ubuntu", "Debian", "OpenSUSE", "SUSE", "CoreOS", "FreeBSD", "Kylin", "UnionTech", "TencentOS", "Other Linux"}
	InstancesOut []Instances
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
}

func DescribeInstances(region string, running bool, specifiedInstanceID string, offSet int64) []Instances {
	request := cvm.NewDescribeInstancesRequest()
	request.Offset = common.Int64Ptr(offSet)
	request.Limit = common.Int64Ptr(100)
	request.SetScheme("https")
	if running {
		request.Filters = []*cvm.Filter{
			{
				Name:   common.StringPtr("instance-state"),
				Values: common.StringPtrs([]string{"RUNNING"}),
			},
		}
	}
	if specifiedInstanceID != "all" {
		request.InstanceIds = common.StringPtrs([]string{specifiedInstanceID})
	}
	response, err := CVMClient(region).DescribeInstances(request)
	errutil.HandleErr(err)
	InstancesList := response.Response.InstanceSet
	log.Infof("正在 %s 区域中查找实例 (Looking for instances in the %s region)", region, region)
	InstancesTotalCount := *response.Response.TotalCount
	if len(InstancesList) != 0 {
		log.Infof("在 %s 区域下找到 %d 个实例 (Found %d instances in %s region)", region, len(InstancesList), len(InstancesList), region)
		var (
			PrivateIpAddressList []string
			PublicIpAddressList  []string
			PrivateIpAddress     string
			PublicIpAddress      string
			OSType               string
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
			newOSname := strings.Split(*v.OsName, " ")[0]
			if find(LinuxSet, newOSname) {
				OSType = "linux"
			} else {
				OSType = "windows"
			}
			obj := Instances{
				InstanceId:       *v.InstanceId,
				InstanceName:     *v.InstanceName,
				OSName:           *v.OsName,
				OSType:           OSType,
				Status:           *v.InstanceState,
				PrivateIpAddress: PrivateIpAddress,
				PublicIpAddress:  PublicIpAddress,
				RegionId:         *v.Placement.Zone,
			}
			InstancesOut = append(InstancesOut, obj)
		}
	}
	if InstancesTotalCount > int64(len(InstancesOut)) {
		_ = DescribeInstances(region, running, specifiedInstanceID, int64(len(InstancesOut)))
	}
	return InstancesOut
}

func ReturnInstancesList(region string, running bool, specifiedInstanceID string) []Instances {
	var InstancesList []Instances
	var Instance []Instances
	if region == "all" {
		for _, j := range GetCVMRegions() {
			region := *j.Region
			Instance = DescribeInstances(region, running, specifiedInstanceID, 0)
			InstancesOut = nil
			for _, i := range Instance {
				InstancesList = append(InstancesList, i)
			}
		}
	} else {
		InstancesList = DescribeInstances(region, running, specifiedInstanceID, 0)
	}
	return InstancesList
}

func PrintInstancesListRealTime(region string, running bool, specifiedInstanceID string) {
	InstancesList := ReturnInstancesList(region, running, specifiedInstanceID)
	var data = make([][]string, len(InstancesList))
	for i, o := range InstancesList {
		SN := strconv.Itoa(i + 1)
		data[i] = []string{SN, o.InstanceId, o.InstanceName, o.OSName, o.OSType, o.Status, o.PrivateIpAddress, o.PublicIpAddress, o.RegionId}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 CVM 实例 (No CVM instances found)")
	} else {
		Caption := "CVM 资源 (CVM resources)"
		cloud.PrintTable(td, Caption)
	}
	cmdutil.WriteCacheFile(td, "tencent", "cvm", region, specifiedInstanceID)
}

func PrintInstancesListHistory(region string, running bool, specifiedInstanceID string) {
	cmdutil.PrintECSCacheFile(header, region, specifiedInstanceID, "tencent", "CVM", running)
}

func PrintInstancesList(region string, running bool, specifiedInstanceID string, cvmFlushCache bool) {
	if cvmFlushCache {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	} else {
		PrintInstancesListHistory(region, running, specifiedInstanceID)
	}
}
