package aliram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/teamssix/cf/pkg/util/errutil"

	"github.com/teamssix/cf/pkg/cloud/alibaba/aliecs"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alioss"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alirds"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	log "github.com/sirupsen/logrus"

	"github.com/teamssix/cf/pkg/cloud"
)

var header = []string{"序号 (SN)", "策略名称 (PolicyName)", "描述 (Description)"}
var header2 = []string{"序号 (SN)", "可执行的操作 (Available actions)", "描述 (Description)"}
var SN = 1
var data2 [][]string

const (
	osslsAction        = "cf alibaba oss ls"
	osslsDescription   = "列出 OSS 资源"
	ossgetAction       = "cf alibaba oss obj get"
	ossgetDescription  = "下载 OSS 资源"
	ecslsAction        = "cf alibaba ecs ls"
	ecslsDescription   = "列出 ECS 资源"
	ecsexecAction      = "cf alibaba ecs exec"
	ecsexecDescription = "在 ECS 上执行命令"
	rdslsAction        = "cf alibaba rds ls"
	rdslsDescription   = "列出 RDS 资源"
	consoleAction      = "cf alibaba console"
	consoleDescription = "接管控制台"
)

func ListPermissions() {
	// 获取当前AK的用户名
	userName := getCallerIdentity()
	log.Infof("当前用户名为 %s (Current username is %s)", userName, userName)
	var data [][]string
	// 如果是root用户，直接返回root权限
	if userName == "root" {
		data = append(data, []string{"1", "AdministratorAccess", "管理所有阿里云资源的权限"})
		var td = cloud.TableData{Header: header, Body: data}
		Caption := "当前凭证具备的权限 (Permissions owned)"
		cloud.PrintTable(td, Caption)
		fmt.Println()
		data2 = appendData(osslsAction, osslsDescription)
		data2 = appendData(ossgetAction, ossgetDescription)
		data2 = appendData(ecslsAction, ecslsDescription)
		data2 = appendData(ecsexecAction, ecsexecDescription)
		data2 = appendData(rdslsAction, rdslsDescription)
		data2 = appendData(consoleAction, consoleDescription)
		var td2 = cloud.TableData{Header: header2, Body: data2}
		Caption2 := "当前凭证可以执行的操作 (Available actions)"
		cloud.PrintTable(td2, Caption2)
	} else {
		// 如果不是root用户，获取当前用户的权限列表
		data, err := listAllPoliciesForUser(userName)
		if err == nil {
			if len(data) == 0 {
				log.Infoln("当前凭证没有任何权限 (The current access key does not have any permissions)")
			} else {
				var td = cloud.TableData{Header: header, Body: data}
				Caption := "当前凭证具备的权限 (Permissions owned)"
				cloud.PrintTable(td, Caption)
				fmt.Println()
				for _, o := range data {
					switch {
					case strings.Contains(o[1], "AdministratorAccess"):
						data2 = appendData(osslsAction, osslsDescription)
						data2 = appendData(ecslsAction, ecslsDescription)
						data2 = appendData(ecsexecAction, ecsexecDescription)
						data2 = appendData(rdslsAction, rdslsDescription)
						data2 = appendData(consoleAction, consoleDescription)
					case strings.Contains(o[1], "ReadOnlyAccess"):
						data2 = appendData(osslsAction, osslsDescription)
						data2 = appendData(ecslsAction, ecslsDescription)
						data2 = appendData(rdslsAction, rdslsDescription)
					case strings.Contains(o[1], "AliyunOSSFullAccess"):
						data2 = appendData(osslsAction, osslsDescription)
						data2 = appendData(ossgetAction, ossgetDescription)
					case strings.Contains(o[1], "AliyunOSSReadOnlyAccess"):
						data2 = appendData(osslsAction, osslsDescription)
					case strings.Contains(o[1], "AliyunECSFullAccess"):
						data2 = appendData(ecslsAction, ecslsDescription)
						data2 = appendData(ecsexecAction, ecsexecDescription)
					case strings.Contains(o[1], "AliyunECSReadOnlyAccess"):
						data2 = appendData(ecslsAction, ecslsDescription)
					case strings.Contains(o[1], "AliyunECSAssistantFullAccess"):
						data2 = appendData(ecsexecAction, ecsexecDescription)
					case strings.Contains(o[1], "AliyunRDSReadOnlyAccess"):
						data2 = appendData(rdslsAction, rdslsDescription)
					case strings.Contains(o[1], "AliyunRAMFullAccess"):
						data2 = appendData(consoleAction, consoleDescription)
					}
				}
				if len(data2) == 0 {
					log.Infoln("当前凭证没有可以执行的操作 (Not available actions)")
				} else {
					var td2 = cloud.TableData{Header: header2, Body: data2}
					Caption2 := "当前凭证可以执行的操作 (Available actions)"
					cloud.PrintTable(td2, Caption2)
				}
			}
		} else if strings.Contains(err.Error(), "ErrorCode: NoPermission") {
			log.Debugln("当前凭证不具备 RAM 读权限 (No RAM read permissions)")
			obj1, obj2 := traversalPermissions()
			var data1 = make([][]string, len(obj1))
			var data2 = make([][]string, len(obj2))
			if len(obj1) == 0 {
				log.Infoln("当前凭证没有遍历到任何权限 (No permissions were found for the current access key)")
			} else {
				for i, o := range obj1 {
					SN := strconv.Itoa(i + 1)
					data1[i] = []string{SN, o[0], o[1]}
				}
				var td1 = cloud.TableData{Header: header, Body: data1}
				Caption1 := "当前凭证具备的权限 (Permissions owned)"
				cloud.PrintTable(td1, Caption1)
				fmt.Println()
				for j, o := range obj2 {
					SN := strconv.Itoa(j + 1)
					data2[j] = []string{SN, o[0], o[1]}
				}
				var td2 = cloud.TableData{Header: header2, Body: data2}
				Caption2 := "当前凭证可以执行的操作 (Available actions)"
				cloud.PrintTable(td2, Caption2)
			}
		} else {
			log.Debugln(err)
		}
	}
}

func getCallerIdentity() string {
	request := sts.CreateGetCallerIdentityRequest()
	request.Scheme = "https"
	response, err := STSClient().GetCallerIdentity(request)
	errutil.HandleErr(err)
	accountArn := response.Arn
	var userName string
	if accountArn[len(accountArn)-4:] == "root" {
		userName = "root"
	} else {
		userName = strings.Split(accountArn, "/")[1]
	}
	log.Debugf("获得到当前凭证的用户名为 %s (The user name to get the current credentials is %s)", userName, userName)
	return userName
}

func listAllPoliciesForUser(userName string) ([][]string, error) {
	data, err := listPoliciesForUser(userName)
	if err != nil {
		return nil, err
	} else {
		groups := listGroupsForUser(userName)
		if len(groups) > 0 {
			for _, g := range groups {
				for _, i := range listPoliciesForGroup(g) {
					data = append(data, i)
				}
			}
		}
		return data, err
	}
}

func listPoliciesForUser(userName string) ([][]string, error) {
	request := ram.CreateListPoliciesForUserRequest()
	request.Scheme = "https"
	request.UserName = userName
	response, err := RAMClient().ListPoliciesForUser(request)
	if err == nil {
		log.Debugf("成功获取到 %s 用户的权限信息 (Successfully obtained permission information for %s user)", userName, userName)
		var data [][]string
		for n, i := range response.Policies.Policy {
			SN := strconv.Itoa(n + 1)
			data = append(data, []string{SN, i.PolicyName, i.Description})
		}
		return data, err
	} else {
		return nil, err
	}
}

func listPoliciesForGroup(groupName string) [][]string {
	// 获取用户组权限
	request := ram.CreateListPoliciesForGroupRequest()
	request.Scheme = "https"
	request.GroupName = groupName
	response, err := RAMClient().ListPoliciesForGroup(request)
	errutil.HandleErr(err)
	var data [][]string
	for n, i := range response.Policies.Policy {
		SN := strconv.Itoa(n + 1)
		data = append(data, []string{SN, i.PolicyName, i.Description})
	}
	return data
}

func listGroupsForUser(userName string) []string {
	request := ram.CreateListGroupsForUserRequest()
	request.Scheme = "https"
	request.UserName = userName
	response, err := RAMClient().ListGroupsForUser(request)
	errutil.HandleErr(err)
	var groups []string
	for _, g := range response.Groups.Group {
		groups = append(groups, g.GroupName)
	}
	return groups
}

func traversalPermissions() ([][]string, [][]string) {
	var obj1 [][]string
	var obj2 [][]string
	// 1. cf alibaba oss get && cf alibaba oss ls
	log.Debugln("正在判断 oss get 权限 (Determining the permission of oss get)")
	OSSCollector := &alioss.OSSCollector{}
	OSSCollector.OSSClient("cn-hangzhou")
	tempBucketName := "teamssix-e2mxrjpkwadnybqzzzihlmxnsowmcpakshhakq2anj6j8ez"
	err1 := OSSCollector.Client.CreateBucket(tempBucketName)
	if err1 == nil {
		obj1 = append(obj1, []string{"AliyunOSSFullAccess", "管理对象存储服务(OSS)权限"})
		obj2 = append(obj2, []string{osslsAction, osslsDescription})
		obj2 = append(obj2, []string{ossgetAction, ossgetDescription})
		err1_2 := OSSCollector.Client.DeleteBucket(tempBucketName)
		if err1_2 != nil {
			log.Traceln(err1_2.Error())
		}
	} else {
		log.Traceln(err1.Error())
		log.Debugln("正在判断 oss ls 权限 (Determining the permission of oss ls)")
		_, err11 := OSSCollector.ListBuckets("all", "all")
		if err11 == nil {
			obj1 = append(obj1, []string{"AliyunOSSReadOnlyAccess", "只读访问对象存储服务(OSS)的权限"})
			obj2 = append(obj2, []string{osslsAction, osslsDescription})
		} else {
			log.Traceln(err11.Error())
		}
	}
	// 2. cf alibaba ecs ls
	log.Debugln("正在判断 ecs ls 权限 (Determining the permission of ecs ls)")
	request := ecs.CreateDescribeVpcsRequest()
	request.Scheme = "https"
	_, err2 := aliecs.ECSClient("cn-beijing").DescribeVpcs(request)
	if err2 == nil {
		obj1 = append(obj1, []string{"AliyunECSReadOnlyAccess", "只读访问云服务器服务(ECS)的权限"})
		obj2 = append(obj2, []string{ecslsAction, ecslsDescription})
	} else {
		log.Traceln(err2.Error())
	}
	// 3. cf alibaba ecs exec
	log.Debugln("正在判断 ecs exec 权限 (Determining the permission of ecs exec)")
	request3 := ecs.CreateInvokeCommandRequest()
	request3.Scheme = "https"
	request3.CommandId = "abcdefghijklmn"
	request3.InstanceId = &[]string{"abcdefghijklmn"}
	_, err3 := aliecs.ECSClient("cn-beijing").InvokeCommand(request3)
	if !strings.Contains(err3.Error(), "ErrorCode: Forbidden.RAM") {
		obj1 = append(obj1, []string{"AliyunECSAssistantFullAccess", "管理 ECS 云助手服务的权限"})
		obj2 = append(obj2, []string{ecsexecAction, ecsexecDescription})
	} else {
		log.Traceln(err3.Error())
	}
	// 4. cf alibaba rds ls
	log.Debugln("正在判断 rds ls 权限 (Determining the permission of rds ls)")
	_, err4 := alirds.DescribeDBInstances("cn-beijing", true, "all", "all", "")
	if err4 == nil {
		obj1 = append(obj1, []string{"AliyunRDSReadOnlyAccess", "只读访问云数据库服务(RDS)的权限"})
		obj2 = append(obj2, []string{rdslsAction, rdslsDescription})
	} else {
		log.Traceln(err4.Error())
	}
	// 5. cf alibaba console
	log.Debugln("正在判断 console 权限 (Determining the permission of console)")
	request5 := ram.CreateDetachPolicyFromUserRequest()
	request5.Scheme = "https"
	request5.PolicyType = "System"
	request5.PolicyName = "test"
	request5.UserName = "test"
	_, err5 := RAMClient().DetachPolicyFromUser(request5)
	if !strings.Contains(err5.Error(), "ErrorCode: NoPermission") {
		obj1 = append(obj1, []string{"AliyunRAMFullAccess", "管理访问控制(RAM)的权限，即管理用户以及授权的权限"})
		obj2 = append(obj2, []string{consoleAction, consoleDescription})
	} else {
		log.Traceln(err5.Error())
	}
	return obj1, obj2
}

func appendData(action string, description string) [][]string {
	var actionList []string
	for _, o := range data2 {
		actionList = append(actionList, o[1])
	}
	sort.Strings(actionList)
	index := sort.SearchStrings(actionList, action)
	if index < len(actionList) && actionList[index] == action {
		log.Tracef("当前 data2 中已存在 %s (%s already exists in the current data2 array)", action, action)
	} else {
		data2 = append(data2, []string{strconv.Itoa(SN), action, description})
		SN = SN + 1
	}
	return data2
}
