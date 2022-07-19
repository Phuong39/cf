package tencentvpc

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcvm"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"strconv"
	"strings"
	"time"
)

//此模块用于新增安全组(规则)并绑定对应实例

var (
	CVMCacheFilePath = cmdutil.ReturnCacheFile("tencent", "CVM")
)

type securityGroupPolicySet struct {
	Protocol  string
	Port      string
	CidrBlock string
	Action    string
}

//创建安全组并返回安全组id
func CreateSecurityGroup(region string) *string {
	request := vpc.NewCreateSecurityGroupRequest()
	request.SetScheme("https")
	groupName := strconv.FormatInt(time.Now().Unix(), 10)
	request.GroupName = common.StringPtr(groupName)
	request.GroupDescription = common.StringPtr("cf_" + groupName)
	response, err := VPCClient(region).CreateSecurityGroup(request)
	util.HandleErr(err)
	securityGroupId := response.Response.SecurityGroup.SecurityGroupId
	log.Infof("得到VPC安全组 id 为 %s 地区为 %s (Vpc security group id %s region %s ):", *securityGroupId, region, *securityGroupId, region)
	return securityGroupId
}

//添加规则
func CreateSecurityGroupPolicies(region string, securityGroupId *string) {
	//选择出站/入站
	var rule string
	prompt := &survey.Select{
		Message: "选择你要添加的出站/入站规则 (Select the Egress/Ingress rule you want to add): ",
		Options: []string{"出站规则 (Egress)", "入站规则 (Ingress)"},
	}
	err := survey.AskOne(prompt, &rule)
	util.HandleErr(err)
	request := vpc.NewCreateSecurityGroupPoliciesRequest()
	request.SetScheme("https")
	request.SecurityGroupId = securityGroupId
	var qs = []*survey.Question{
		{
			Name:   "Protocol",
			Prompt: &survey.Input{Message: "Protocol (必须 Required) (TCP,UDP,ICMP,ICMPv6,ALL)" + ":"},
		},
		{
			Name:   "Port",
			Prompt: &survey.Input{Message: "Port (必须 Required) (ALL, (1,65535))" + ":"},
		},
		{
			Name:   "CidrBlock",
			Prompt: &survey.Input{Message: "CidrBlock (必须 Required) (0.0.0.0/0)" + ":"},
		},
		{
			Name:   "Action",
			Prompt: &survey.Input{Message: "Action (必须 Required) (ACCEPT/DROP)" + ":"},
		},
	}
	secGP := securityGroupPolicySet{}
	if rule == "出站规则 (Egress)" {
		err1 := survey.Ask(qs, &secGP)
		util.HandleErr(err1)
		request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
			Egress: []*vpc.SecurityGroupPolicy{
				{
					Protocol:  common.StringPtr(secGP.Protocol),
					Port:      common.StringPtr(secGP.Port),
					CidrBlock: common.StringPtr(secGP.CidrBlock),
					Action:    common.StringPtr(secGP.Action),
				},
			},
		}
		_, err2 := VPCClient(region).CreateSecurityGroupPolicies(request)
		util.HandleErr(err2)
	} else {
		err1 := survey.Ask(qs, &secGP)
		util.HandleErr(err1)
		request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
			Ingress: []*vpc.SecurityGroupPolicy{
				{
					Protocol:  common.StringPtr(secGP.Protocol),
					Port:      common.StringPtr(secGP.Port),
					CidrBlock: common.StringPtr(secGP.CidrBlock),
					Action:    common.StringPtr(secGP.Action),
				},
			},
		}
		_, err2 := VPCClient(region).CreateSecurityGroupPolicies(request)
		util.HandleErr(err2)
	}
}

//解绑安全组
func DisassociateSecurityGroups(region string, securityGroupId *string, instanceIds string) {
	request := cvm.NewDisassociateSecurityGroupsRequest()
	request.SetScheme("https")
	request.SecurityGroupIds = common.StringPtrs([]string{*securityGroupId})
	request.InstanceIds = common.StringPtrs([]string{instanceIds})
	_, err := tencentcvm.CVMClient(region).DisassociateSecurityGroups(request)
	util.HandleErr(err)
}

//删除创建的安全组
func DeleteSecurityGroup(region string, securityGroupId *string) {
	request := vpc.NewDeleteSecurityGroupRequest()
	request.SetScheme("https")
	request.SecurityGroupId = securityGroupId
	_, err := VPCClient(region).DeleteSecurityGroup(request)
	util.HandleErr(err)
	log.Debugf("区域 %s 下的安全组 id 为 %s 已删除 (Region %s Security group id %s is deleted)", region, *securityGroupId, region, *securityGroupId)
}

//绑定安全组规则->CVM
func AssociateSecurityGroups(oldregion string, region string, securityGroupId *string, instancesMap map[string]string) {
	request := cvm.NewAssociateSecurityGroupsRequest()
	request.SetScheme("https")
	request.SecurityGroupIds = common.StringPtrs([]string{*securityGroupId})
	//根据传来的区域组成实例列表
	var InstancesList []string
	for k, v := range instancesMap {
		if v == oldregion {
			InstancesList = append(InstancesList, k)
		}
	}
	//选择实例机器
	log.Warnln(InstancesList)
	var InstanceId string
	prompt := &survey.Select{
		Message: "选择你要绑定的腾讯云实例 (Select the tencent cloud instance you want to bind): ",
		Options: InstancesList,
	}
	err := survey.AskOne(prompt, &InstanceId)
	util.HandleErr(err)
	//设置实例机器 需要跟安全组区域对应
	request.InstanceIds = common.StringPtrs([]string{InstanceId})
	_, err1 := tencentcvm.CVMClient(region).AssociateSecurityGroups(request)
	util.HandleErr(err1)
	log.Debugf("成功将安全组 %s 绑定至实例 %s (Success security group bound %s to the instance %s)", *securityGroupId, InstanceId, *securityGroupId, InstanceId)
}

func VPCAdd(regionList []string, instancesMap map[string]string) {
	//根据已有的实例区域创建
	var region string
	newregionList := util.RemoveDuplicatesAndEmpty(regionList)
	prompt := &survey.Select{
		Message: "选择你要创建的安全组区域 (Select the security group you want to create): ",
		Options: newregionList,
	}
	err := survey.AskOne(prompt, &region)
	util.HandleErr(err)
	str := strings.Split(region, "-")
	newRegion := str[0] + "-" + str[1]
	securityGroupId := CreateSecurityGroup(newRegion)
	CreateSecurityGroupPolicies(newRegion, securityGroupId)
	AssociateSecurityGroups(region, newRegion, securityGroupId, instancesMap)
}

func VPCDel(instancesList []string, instancesMap map[string]string, securityGroupId *string) {
	//先解绑 再删除
	var InstanceId string
	prompt := &survey.Select{
		Message: "选择你要解绑的腾讯云实例 (Select the tencent cloud instance you want to unbind): ",
		Options: instancesList,
	}
	err := survey.AskOne(prompt, &InstanceId)
	util.HandleErr(err)
	str := strings.Split(instancesMap[InstanceId], "-")
	region := str[0] + "-" + str[1]
	DisassociateSecurityGroups(region, securityGroupId, InstanceId)
	DeleteSecurityGroup(region, securityGroupId)
}

func VPCControl(op string, securityGroupId *string) {
	cacheData := cmdutil.ReadCacheFile(CVMCacheFilePath, "tencent", "CVM")
	var InstancesList []string
	var InstancesMap map[string]string
	var regionList []string
	InstancesMap = make(map[string]string)
	for _, i := range cacheData {
		InstancesList = append(InstancesList, i[1])
		InstancesMap[i[1]] = i[8]
		regionList = append(regionList, i[8])
	}
	switch op {
	case "add":
		{
			VPCAdd(regionList, InstancesMap)
		}
	case "del":
		{
			VPCDel(InstancesList, InstancesMap, securityGroupId)
		}
	}
}
