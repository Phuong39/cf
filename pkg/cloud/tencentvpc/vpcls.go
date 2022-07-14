package tencentvpc

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"strconv"
)

var (
	VPCCacheFilePath = cmdutil.ReturnVPCCacheFile("tencent")
	securityGp       = []string{"序号 (SN)", "安全组实例ID (SecurityGroupId)", "区域 ID (Region ID)"}
	eGress           = []string{"序号 (SN)", "动作 (Action)", "协议 (Protocol)", "端口 (Port)", "网段或IP (CidrBlock)", "安全组实例ID (SecurityGroupId)", "安全组规则描述 (PolicyDescription)", "区域 ID (Region ID)"}
	inGress          = []string{"序号 (SN)", "动作 (Action)", "协议 (Protocol)", "端口 (Port)", "网段或IP (CidrBlock)", "安全组实例ID (SecurityGroupId)", "安全组规则描述 (PolicyDescription)", "区域 ID (Region ID)"}
)

//查看安全组
func DescribeSecurityGroups(region string) ([]*string, []string) {
	var vpcSecGP []*string
	var regions []string
	request := vpc.NewDescribeSecurityGroupsRequest()
	request.SetScheme("https")
	response, err := VPCClient(region).DescribeSecurityGroups(request)
	util.HandleErr(err)
	log.Tracef("正在 %s 区域中查找安全组 (Looking for vpc security group in the %s region)", region, region)
	securityGroup := response.Response
	if *securityGroup.TotalCount > 0 {
		log.Debugf("在 %s 区域下找到安全组 (Found instances in %s region)", region, region)
		for _, i := range securityGroup.SecurityGroupSet {
			vpcSecGP = append(vpcSecGP, i.SecurityGroupId)
		}
		regions = append(regions, region)
		return vpcSecGP, regions
	} else {
		log.Debugf("未在 %s 区域下查找到安全组！(Can not find vpc security group in %s!)", region, region)
	}
	return nil, nil
}

//查看安全组的规则
func DescribeSecurityGroupPolicies(region string, securityGroupId *string) *vpc.SecurityGroupPolicySet {
	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	request.SetScheme("https")
	request.SecurityGroupId = securityGroupId
	response, err := VPCClient(region).DescribeSecurityGroupPolicies(request)
	util.HandleErr(err)
	return response.Response.SecurityGroupPolicySet
}

//处理全部或指定地区的安全组规则
func ReturnVPCSecurityGroupPoliciesList(region string, securityGroupId *string) []*vpc.SecurityGroupPolicySet {
	var secGPList []*vpc.SecurityGroupPolicySet
	if region == "all" && *securityGroupId == "all" {
		for _, j := range GetVPCRegions() {
			region := *j.Region
			vpcSecG, _ := DescribeSecurityGroups(region)
			if vpcSecG != nil {
				for _, i := range vpcSecG {
					vpcSecGP := DescribeSecurityGroupPolicies(region, i)
					secGPList = append(secGPList, vpcSecGP)
				}
			}
		}
	} else {
		vpcSecGP := DescribeSecurityGroupPolicies(region, securityGroupId)
		secGPList = append(secGPList, vpcSecGP)
	}
	return secGPList
}

//处理全部或指定地区的安全组安全组ID
func ReturnVPCSecurityGroupList(region string, securityGroupId *string) ([]*string, []string) {
	var vpcSecGP []*string
	var regions []string
	if region == "all" && *securityGroupId == "all" {
		for _, j := range GetVPCRegions() {
			region := *j.Region
			vpcSecG, org := DescribeSecurityGroups(region)
			if vpcSecG != nil {
				for _, i := range vpcSecG {
					vpcSecGP = append(vpcSecGP, i)
				}
				for _, i := range org {
					regions = append(regions, i)
				}
			}
		}
	} else {
		vpcSecGP, regions = DescribeSecurityGroups(region)
	}
	return vpcSecGP, regions
}

func PrintVPCSecurityGroupPoliciesListRealTime(region string, securityGroupId *string) {
	VPCList := ReturnVPCSecurityGroupPoliciesList(region, securityGroupId)
	var eGressData = make([][]string, len(VPCList))
	var inGressdata = make([][]string, len(VPCList))
	for i, o := range VPCList {
		SN := strconv.Itoa(i + 1)
		for _, j := range o.Egress {
			eGressData[i] = []string{SN, *j.Action, *j.Protocol, *j.Port, *j.CidrBlock, *securityGroupId, *j.PolicyDescription, region}
		}
		for _, j := range o.Ingress {
			inGressdata[i] = []string{SN, *j.Action, *j.Protocol, *j.Port, *j.CidrBlock, *securityGroupId, *j.PolicyDescription, region}
		}
	}
	var etd = cloud.TableData{Header: eGress, Body: eGressData}
	var itd = cloud.TableData{Header: inGress, Body: inGressdata}
	if len(eGressData) == 0 {
		log.Info("未发现VPC出站规则 (Can not find vpc security group egress)")
		cmdutil.WriteCacheFile(etd, VPCCacheFilePath)
	} else {
		eCaption := "VPC安全组出站规则 (vpc security group egress)"
		cloud.PrintTable(etd, eCaption)
		cmdutil.WriteCacheFile(etd, VPCCacheFilePath)
	}
	if len(inGressdata) == 0 {
		log.Info("未发现VPC安全组入站规则 (Can not find vpc security group ingress)")
		cmdutil.WriteCacheFile(itd, VPCCacheFilePath)
	} else {
		iCaption := "VPC安全组入站规则 (VPC security group ingress)"
		cloud.PrintTable(itd, iCaption)
		cmdutil.WriteCacheFile(itd, VPCCacheFilePath)
	}
}

func PrintVPCSecurityGroupListRealTime(region string, securityGroupId *string) {
	VPCList, regions := ReturnVPCSecurityGroupList(region, securityGroupId)
	var td = make([][]string, len(VPCList))
	for i, o := range VPCList {
		for _, j := range regions {
			SN := strconv.Itoa(i + 1)
			//FIXME region变量覆盖
			td[i] = []string{SN, *o, j}
		}
	}
	var data = cloud.TableData{Header: securityGp, Body: td}
	if len(td) == 0 {
		log.Info("未发现VPC安全组实例ID (Can not vpc security group id)")
		cmdutil.WriteCacheFile(data, VPCCacheFilePath)
	} else {
		Caption := "VPC安全组实例ID (VPC security group id)"
		cloud.PrintTable(data, Caption)
		cmdutil.WriteCacheFile(data, VPCCacheFilePath)
	}
}

//安全组规则历史
func PrintVPCSecurityGroupPoliciesListHistory(region string, securityGroupId *string) {
	if cmdutil.FileExists(VPCCacheFilePath) {
		cmdutil.PrintCFCacheFile(VPCCacheFilePath, eGress, region, *securityGroupId, "tencent")
		cmdutil.PrintCFCacheFile(VPCCacheFilePath, inGress, region, *securityGroupId, "tencent")
	} else {
		PrintVPCSecurityGroupPoliciesListRealTime(region, securityGroupId)
	}
}

func PrintVPCSecurityGroupPoliciesList(region string, securityGroupId *string, vpcFlushCache bool) {
	if vpcFlushCache {
		PrintVPCSecurityGroupPoliciesListRealTime(region, securityGroupId)
	} else {
		PrintVPCSecurityGroupPoliciesListHistory(region, securityGroupId)
	}
}

//安全组ID历史
func PrintVPCSecurityGroupListHistory(region string, securityGroupId *string) {
	if cmdutil.FileExists(VPCCacheFilePath) {
		cmdutil.PrintCFCacheFile(VPCCacheFilePath, eGress, region, *securityGroupId, "tencent")
		cmdutil.PrintCFCacheFile(VPCCacheFilePath, inGress, region, *securityGroupId, "tencent")
	} else {
		PrintVPCSecurityGroupPoliciesListRealTime(region, securityGroupId)
	}
}

func PrintVPCSecurityGroupList(region string, securityGroupId *string, vpcFlushCache bool) {
	if vpcFlushCache {
		PrintVPCSecurityGroupListRealTime(region, securityGroupId)
	} else {
		PrintVPCSecurityGroupListHistory(region, securityGroupId)
	}
}
