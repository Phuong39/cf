package tencentvpc

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcvm"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type securityGroupPolicyList struct {
	securityGroupId *string
	region          string
	secGPList       *vpc.SecurityGroupPolicySet
}

var (
	VPCCacheFilePath = cmdutil.ReturnCacheFile("tencent", "VPC")
	vpcHeader        = []string{"序号 (SN)", "安全组 ID (Security Group ID)", "类型 (Type)", "动作 (Action)", "协议 (Protocol)", "端口 (Port)", "网段或IP (CidrBlock)", "安全组规则描述 (PolicyDescription)", "区域 ID (Region ID)"}
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
		log.Debugf("未在 %s 区域下查找到安全组(Can not find vpc security group in %s)", region, region)
	}
	return nil, nil
}

//查看安全组的规则
func DescribeSecurityGroupPolicies(region string, securityGroupId *string) (*vpc.SecurityGroupPolicySet, error) {
	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	request.SetScheme("https")
	request.SecurityGroupId = securityGroupId
	response, err := VPCClient(region).DescribeSecurityGroupPolicies(request)
	util.HandleErr(err)
	return response.Response.SecurityGroupPolicySet, err
}

//处理全部或指定地区的安全组规则
func ReturnVPCSecurityGroupPoliciesList(region string, securityGroupId *string) []securityGroupPolicyList {
	var (
		secGroupPolicy     securityGroupPolicyList
		secGroupPolicyList []securityGroupPolicyList
		secGPList          []*vpc.SecurityGroupPolicySet
	)
	switch {
	case region == "all" && *securityGroupId == "all":
		for _, j := range tencentcvm.GetCVMRegions() {
			specifiedRegion := *j.Region
			vpcSecG, _ := DescribeSecurityGroups(specifiedRegion)
			if vpcSecG != nil {
				for _, i := range vpcSecG {
					vpcSecGP, _ := DescribeSecurityGroupPolicies(specifiedRegion, i)
					secGroupPolicy = securityGroupPolicyList{
						securityGroupId: i,
						region:          specifiedRegion,
						secGPList:       vpcSecGP,
					}
					secGroupPolicyList = append(secGroupPolicyList, secGroupPolicy)
				}
			}
		}
	case region == "all" && *securityGroupId != "all":
		for _, j := range tencentcvm.GetCVMRegions() {
			specifiedRegion := *j.Region
			vpcSecG, _ := DescribeSecurityGroups(specifiedRegion)
			if vpcSecG != nil {
				for _, i := range vpcSecG {
					vpcSecGP, _ := DescribeSecurityGroupPolicies(specifiedRegion, i)
					if *i == *securityGroupId {
						secGroupPolicy = securityGroupPolicyList{
							securityGroupId: i,
							region:          specifiedRegion,
							secGPList:       vpcSecGP,
						}
						secGroupPolicyList = append(secGroupPolicyList, secGroupPolicy)
					}
				}
			}
		}
	case region != "all" && *securityGroupId == "all":
		vpcSecG, _ := DescribeSecurityGroups(region)
		if vpcSecG != nil {
			for _, i := range vpcSecG {
				vpcSecGP, _ := DescribeSecurityGroupPolicies(region, i)
				secGroupPolicy := securityGroupPolicyList{
					securityGroupId: i,
					region:          region,
					secGPList:       vpcSecGP,
				}
				secGroupPolicyList = append(secGroupPolicyList, secGroupPolicy)
			}
		}
	default:
		vpcSecGP, _ := DescribeSecurityGroupPolicies(region, securityGroupId)
		secGPList = append(secGPList, vpcSecGP)
		secGroupPolicy := securityGroupPolicyList{
			securityGroupId: securityGroupId,
			region:          region,
			secGPList:       vpcSecGP,
		}
		secGroupPolicyList = append(secGroupPolicyList, secGroupPolicy)
	}
	return secGroupPolicyList
}

func PrintVPCSecurityGroupPoliciesListRealTime(region string, securityGroupId *string) {
	var (
		datalen int
		number  int
	)
	secGroupPolicyList := ReturnVPCSecurityGroupPoliciesList(region, securityGroupId)

	for _, o := range secGroupPolicyList {
		for range o.secGPList.Egress {
			datalen = datalen + 1
		}
		for range o.secGPList.Ingress {
			datalen = datalen + 1
		}
	}
	var data = make([][]string, datalen)
	for _, o := range secGroupPolicyList {
		for _, j := range o.secGPList.Egress {
			SN := strconv.Itoa(number + 1)
			data[number] = []string{SN, *o.securityGroupId, "出站规则 (Egress)", *j.Action, *j.Protocol, *j.Port, *j.CidrBlock, *j.PolicyDescription, o.region}
			number = number + 1
		}
		for _, j := range o.secGPList.Ingress {
			SN := strconv.Itoa(number + 1)
			data[number] = []string{SN, *o.securityGroupId, "入站规则 (Ingress)", *j.Action, *j.Protocol, *j.Port, *j.CidrBlock, *j.PolicyDescription, o.region}
			number = number + 1
		}
	}
	var td = cloud.TableData{Header: vpcHeader, Body: data}
	if len(data) == 0 {
		log.Info("未发现 VPC 安全组规则 (Can not find vpc security group egress)")
		cmdutil.WriteCacheFile(td, VPCCacheFilePath, region, *securityGroupId)
	} else {
		eCaption := "VPC 安全组规则 (vpc security group egress)"
		cloud.PrintTable(td, eCaption)
		cmdutil.WriteCacheFile(td, VPCCacheFilePath, region, *securityGroupId)
	}
}

func PrintVPCSecurityGroupPoliciesListHistory(region string, securityGroupId *string) {
	if cmdutil.FileExists(VPCCacheFilePath) {
		cmdutil.PrintECSCacheFile(VPCCacheFilePath, vpcHeader, region, *securityGroupId, "tencent", "VPC")
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
