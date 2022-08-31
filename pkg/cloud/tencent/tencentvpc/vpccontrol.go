package tencentvpc

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcvm"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"os"
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

// 创建安全组并返回安全组id
func CreateSecurityGroup(region string) *string {
	request := vpc.NewCreateSecurityGroupRequest()
	request.SetScheme("https")
	groupName := strconv.FormatInt(time.Now().Unix(), 10)
	request.GroupName = common.StringPtr(groupName)
	request.GroupDescription = common.StringPtr("cf_" + groupName)
	response, err := VPCClient(region).CreateSecurityGroup(request)
	errutil.HandleErr(err)
	securityGroupId := response.Response.SecurityGroup.SecurityGroupId
	log.Debugf("得到 VPC 安全组 id 为 %s ，区域为 %s (Vpc security group id %s ,region %s ):", *securityGroupId, region, *securityGroupId, region)
	return securityGroupId
}

// 添加规则
func CreateSecurityGroupPolicies(region string, securityGroupId *string) {
	//选择出站/入站
	var rule string
	prompt := &survey.Select{
		Message: "选择你要添加的出站/入站规则 (Select the Egress/Ingress rule you want to add): ",
		Options: []string{"出站规则 (Egress)", "入站规则 (Ingress)"},
	}
	err := survey.AskOne(prompt, &rule)
	errutil.HandleErr(err)
	request := vpc.NewCreateSecurityGroupPoliciesRequest()
	request.SetScheme("https")
	request.SecurityGroupId = securityGroupId
	var qs = []*survey.Question{
		{
			Name:   "Protocol",
			Prompt: &survey.Input{Message: "Protocol (必须 Required) ([TCP],UDP,ICMP,ICMPv6,ALL)" + ":"},
		},
		{
			Name:   "Port",
			Prompt: &survey.Input{Message: "Port (必须 Required) ([ALL], (1,65535))" + ":"},
		},
		{
			Name:   "CidrBlock",
			Prompt: &survey.Input{Message: "CidrBlock (必须 Required) ([0.0.0.0/0])" + ":"},
		},
		{
			Name:   "Action",
			Prompt: &survey.Input{Message: "Action (必须 Required) ([ACCEPT]/DROP)" + ":"},
		},
	}
	secGP := securityGroupPolicySet{}
	err1 := survey.Ask(qs, &secGP)
	errutil.HandleErr(err1)
	switch {
	case secGP.Protocol == "":
		secGP.Protocol = "TCP"
	case secGP.Port == "":
		secGP.Port = "ALL"
	case secGP.CidrBlock == "":
		secGP.CidrBlock = "0.0.0.0/0"
	case secGP.Action == "":
		secGP.Action = "ACCEPT"
	}
	if rule == "出站规则 (Egress)" {
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
		errutil.HandleErr(err2)
	} else {
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
		errutil.HandleErr(err2)
	}
}

// 解绑安全组
func DisassociateSecurityGroups(region string, securityGroupId *string, instanceIds string) {
	request := cvm.NewDisassociateSecurityGroupsRequest()
	request.SetScheme("https")
	request.SecurityGroupIds = common.StringPtrs([]string{*securityGroupId})
	request.InstanceIds = common.StringPtrs([]string{instanceIds})
	_, err := tencentcvm.CVMClient(region).DisassociateSecurityGroups(request)
	errutil.HandleErr(err)
}

// 删除创建的安全组
func DeleteSecurityGroup(region string, securityGroupId *string) {
	request := vpc.NewDeleteSecurityGroupRequest()
	request.SetScheme("https")
	request.SecurityGroupId = securityGroupId
	_, err := VPCClient(region).DeleteSecurityGroup(request)
	errutil.HandleErr(err)
	if err == nil {
		log.Infof("%s 安全组已删除 (%s Security group is deleted)", region, region)
	}
}

// 绑定安全组规则->CVM
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
	var InstanceId string
	prompt := &survey.Select{
		Message: "选择你要绑定的腾讯云实例 (Select the tencent cloud instance you want to bind): ",
		Options: InstancesList,
	}
	err := survey.AskOne(prompt, &InstanceId)
	errutil.HandleErr(err)
	//设置实例机器 需要跟安全组区域对应
	request.InstanceIds = common.StringPtrs([]string{InstanceId})
	_, err1 := tencentcvm.CVMClient(region).AssociateSecurityGroups(request)
	errutil.HandleErr(err1)
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
	errutil.HandleErr(err)
	str := strings.Split(region, "-")
	newRegion := str[0] + "-" + str[1]
	securityGroupId := CreateSecurityGroup(newRegion)
	CreateSecurityGroupPolicies(newRegion, securityGroupId)
	AssociateSecurityGroups(region, newRegion, securityGroupId, instancesMap)
}

func VPCDel(instancesMap map[string]string, CVMCacheData [][]string) {
	var (
		securityGroupIdString string
		securityGroupId       *string
		vpcList               []string
		InstanceId            string
		InstanceIdList        []string
	)
	if len(getVPCCacheData()) == 0 {
		log.Infoln("正在查询 VPC 信息 (Searching for VPC information)")
		*securityGroupId = "all"
		PrintVPCSecurityGroupPoliciesListRealTime("all", securityGroupId)
		if len(getVPCCacheData()) != 0 {
			fmt.Println("")
		}
	}
	cacheData := getVPCCacheData()
	for _, j := range cacheData {
		vpcList = append(vpcList, j[1])
	}
	vpcList = RemoveRepByLoop(vpcList)

	prompt1 := &survey.Select{
		Message: "选择你要删除的 VPC ID (Select the VPC ID you want to delete):",
		Options: vpcList,
	}
	err := survey.AskOne(prompt1, &securityGroupIdString)
	errutil.HandleErr(err)
	securityGroupId = &securityGroupIdString
	for _, i := range CVMCacheData {
		if strings.Contains(i[9], *securityGroupId) {
			InstanceIdList = append(InstanceIdList, i[1])
		}
	}

	//先解绑 再删除
	if len(InstanceIdList) == 0 {
		str := strings.Split(instancesMap[InstanceId], "-")
		region := str[0] + "-" + str[1]
		DeleteSecurityGroup(region, securityGroupId)
	} else {
		var region string
		for _, i := range InstanceIdList {
			str := strings.Split(instancesMap[i], "-")
			region = str[0] + "-" + str[1]
			DisassociateSecurityGroups(region, securityGroupId, i)
		}
		DeleteSecurityGroup(region, securityGroupId)
	}
}

func RemoveRepByLoop(slc []string) []string {
	result := []string{}
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}

func VPCControl() {
	var (
		InstancesList []string
		InstancesMap  map[string]string
		regionList    []string
		op            string
		opList        = []string{"新增 (add)", "删除 (del)"}
	)
	if len(getCVMCacheData()) == 0 {
		log.Infoln("正在查询 CVM 实例信息 (Searching for CVM instance information)")
		tencentcvm.PrintInstancesListRealTime("all", false, "all")
		if len(getCVMCacheData()) == 0 {
			os.Exit(0)
		}
	}
	cacheData := getCVMCacheData()
	InstancesMap = make(map[string]string)
	for _, i := range cacheData {
		InstancesList = append(InstancesList, i[1])
		InstancesMap[i[1]] = i[8]
		regionList = append(regionList, i[8])
	}
	prompt := &survey.Select{
		Message: "选择你要执行的操作 (Select your action): ",
		Options: opList,
	}
	err := survey.AskOne(prompt, &op)
	errutil.HandleErr(err)
	switch op {
	case "新增 (add)":
		{
			VPCAdd(regionList, InstancesMap)
		}
	case "删除 (del)":
		{
			VPCDel(InstancesMap, cacheData)
		}
	}
}

func getCVMCacheData() [][]string {
	cacheData := cmdutil.ReadCacheFile(CVMCacheFilePath, "tencent", "CVM")
	return cacheData
}

func getVPCCacheData() [][]string {
	cacheData := cmdutil.ReadCacheFile(VPCCacheFilePath, "tencent", "VPC")
	return cacheData
}
