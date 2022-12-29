package awsec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"

	"github.com/teamssix/cf/pkg/util/errutil"
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

var (
	DescribeInstancesOut []Instances
	TimestampType        = util.ReturnTimestampType("aws", "ec2")
	header               = []string{"序号 (SN)", "实例 ID (Instance ID)", "实例名称 (Instance Name)", "系统名称 (OS Name)", "系统类型 (OS Type)", "状态 (Status)", "私有 IP (Private IP)", "公网 IP (Public IP)", "区域 ID (Region ID)"}
)

func DescribeInstances(region string, running bool, SpecifiedInstanceID string, NextToken string) []Instances {
	var (
		err error
		//InstancesList []Instances
		result *ec2.DescribeInstancesOutput
		svc    *ec2.EC2
	)
	log.Infof("正在 %s 区域中查找实例 (Looking for instances in the %s region)", region, region)
	if region == "all" {
		svc = EC2Client("all")
	} else {
		svc = EC2Client(region)
	}
	if NextToken == "" {
		result, err = svc.DescribeInstances(nil)
	} else {
		DescribeInstancesInput := ec2.DescribeInstancesInput{
			NextToken: &NextToken,
		}
		result, err = svc.DescribeInstances(&DescribeInstancesInput)
	}

	errutil.HandleErr(err)
	for _, i := range result.Reservations {
		InstancesList := i.Instances
		if len(InstancesList) != 0 {
			for _, i := range InstancesList {
				var (
					InstanceName     string
					PrivateIpAddress string
					PublicIpAddress  string
				)
				for _, tag := range i.Tags {
					if *tag.Key == "Name" {
						InstanceName = *tag.Value
					}
				}
				if i.PrivateIpAddress == nil {
					PrivateIpAddress = ""
				} else {
					PrivateIpAddress = *i.PrivateIpAddress
				}

				if i.PublicIpAddress == nil {
					PublicIpAddress = ""
				} else {
					PublicIpAddress = *i.PublicIpAddress
				}

				obj := Instances{
					InstanceId:       *i.InstanceId,
					InstanceName:     InstanceName,
					OSName:           *i.PlatformDetails,
					OSType:           *i.InstanceType,
					Status:           *i.State.Name,
					PrivateIpAddress: PrivateIpAddress,
					PublicIpAddress:  PublicIpAddress,
					RegionId:         *i.Placement.AvailabilityZone,
				}
				DescribeInstancesOut = append(DescribeInstancesOut, obj)
			}
		}
	}
	if NextToken != "" {
		NextToken = *result.NextToken
		log.Tracef("Next Token: %s", NextToken)
		_ = DescribeInstances(region, running, SpecifiedInstanceID, NextToken)
	}
	return DescribeInstancesOut
}

func ReturnInstancesList(region string, running bool, specifiedInstanceID string, ec2LsAllRegions bool) []Instances {
	var (
		InstancesList []Instances
		Instance      []Instances
		instanceNum   int
	)
	if region == "all" {
		for _, j := range GetEC2Regions() {
			instanceNum = len(InstancesList)
			region := *j.RegionName
			Instance = DescribeInstances(region, running, specifiedInstanceID, "")
			DescribeInstancesOut = nil
			for _, i := range Instance {
				InstancesList = append(InstancesList, i)
			}
			instanceNum = len(InstancesList) - instanceNum
			if instanceNum != 0 {
				log.Warnf("在 %s 区域下找到 %d 个实例 (Found %d instances in %s region)", region, instanceNum, instanceNum, region)
			}
		}
	} else {
		InstancesList = DescribeInstances(region, running, specifiedInstanceID, "")
		instanceNum = len(InstancesList)
		if instanceNum != 0 {
			log.Warnf("在 %s 区域下找到 %d 个实例 (Found %d instances in %s region)", region, len(InstancesList), len(InstancesList), region)
		}
	}
	return InstancesList
}

func PrintInstancesListRealTime(region string, running bool, specifiedInstanceID string, ec2LsAllRegions bool) {
	InstancesList := ReturnInstancesList(region, running, specifiedInstanceID, ec2LsAllRegions)
	var data = make([][]string, len(InstancesList))
	for i, o := range InstancesList {
		if specifiedInstanceID == "all" {
			SN := strconv.Itoa(i + 1)
			data[i] = []string{SN, o.InstanceId, o.InstanceName, o.OSName, o.OSType, o.Status, o.PrivateIpAddress, o.PublicIpAddress, o.RegionId}
		} else if specifiedInstanceID == o.InstanceId {
			SN := strconv.Itoa(i + 1)
			data[i] = []string{SN, o.InstanceId, o.InstanceName, o.OSName, o.OSType, o.Status, o.PrivateIpAddress, o.PublicIpAddress, o.RegionId}
		}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 EC2 资源，可能是因为当前访问密钥权限不够 (No EC2 instances found, Probably because the current Access Key do not have enough permissions)")
	} else {
		Caption := "EC2 资源 (EC2 resources)"
		cloud.PrintTable(td, Caption)
		util.WriteTimestamp(TimestampType)
	}
	cmdutil.WriteCacheFile(td, "aws", "ec2", region, specifiedInstanceID)
}

func PrintInstancesListHistory(region string, running bool, specifiedInstanceID string) {
	cmdutil.PrintECSCacheFile(header, region, specifiedInstanceID, "aws", "ec2", running)
}

func PrintInstancesList(region string, running bool, specifiedInstanceID string, ec2FlushCache bool, ec2LsAllRegions bool) {
	if ec2FlushCache {
		PrintInstancesListRealTime(region, running, specifiedInstanceID, ec2LsAllRegions)
	} else {
		oldTimestamp := util.ReadTimestamp(TimestampType)
		if oldTimestamp == 0 {
			PrintInstancesListRealTime(region, running, specifiedInstanceID, ec2LsAllRegions)
		} else if util.IsFlushCache(oldTimestamp) {
			PrintInstancesListRealTime(region, running, specifiedInstanceID, ec2LsAllRegions)
		} else {
			util.TimeDifference(oldTimestamp)
			PrintInstancesListHistory(region, running, specifiedInstanceID)
		}
	}
}
