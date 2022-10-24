package tencentcam

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcvm"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentlh"
	"github.com/teamssix/cf/pkg/util/errutil"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	lh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"sort"
	"strconv"
	"strings"
)

var header = []string{"序号 (SN)", "策略名称 (PolicyName)", "描述 (Description)"}
var header2 = []string{"序号 (SN)", "可执行的操作 (Available actions)", "描述 (Description)"}
var SN = 1
var data2 [][]string

const (
	cvmLsAction               = "cf tencent cvm ls"
	cvmLsDescription          = "列出 CVM 资源"
	cvmExecAction             = "cf tencent cvm exec"
	cvmExecDescription        = "在 CVM 上执行命令"
	lhLsAction                = "cf tencent lh ls"
	lhLsDescription           = "列出轻量计算服务资源"
	lhExecAction              = "cf tencent lh exec"
	lhExecDescription         = "在轻量计算服务上执行命令"
	lhSSHAction               = "cf tencent lh ssh"
	lhSSHDescription          = "在轻量计算服务上执行 ssh 相关命令"
	vpcLsAction               = "cf tencent vpc ls"
	vpcLsDescription          = "列出 VPC 资源"
	vpcControlAction          = "cf tencent vpc control"
	vpcControlDescription     = "添加或删除 VPC 策略"
	agentUninstallAction      = "cf tencent uninstall"
	agentUninstallDescription = "卸载云镜"
	consoleAction             = "cf tencent console"
	consoleDescription        = "接管控制台"
)

func ListPermissions() {
	userType, accountId := getCallerIdentity()
	log.Infof("当前用户类型为 %s (Current user type is %s)", userType, userType)
	var data [][]string
	if userType == "Root" {
		data = append(data, []string{"1", "AdministratorAccess", "该策略允许您管理账户内所有用户及其权限、财务相关的信息、云服务资产。"})
		var td = cloud.TableData{Header: header, Body: data}
		Caption := "当前凭证具备的权限 (Permissions owned)"
		cloud.PrintTable(td, Caption)
		fmt.Println()
		data2 = appendData(cvmLsAction, cvmLsDescription)
		data2 = appendData(cvmExecAction, cvmExecDescription)
		data2 = appendData(lhLsAction, lhLsDescription)
		data2 = appendData(lhExecAction, lhExecDescription)
		data2 = appendData(lhSSHAction, lhSSHDescription)
		data2 = appendData(vpcLsAction, vpcLsDescription)
		data2 = appendData(vpcControlAction, vpcControlDescription)
		data2 = appendData(agentUninstallAction, agentUninstallDescription)
		data2 = appendData(consoleAction, consoleDescription)
		var td2 = cloud.TableData{Header: header2, Body: data2}
		Caption2 := "当前凭证可以执行的操作 (Available actions)"
		cloud.PrintTable(td2, Caption2)
	} else {
		data, err := listAllPoliciesForUser(accountId)
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
						data2 = appendData(cvmLsAction, cvmLsDescription)
						data2 = appendData(cvmExecAction, cvmExecDescription)
						data2 = appendData(lhLsAction, lhLsDescription)
						data2 = appendData(lhExecAction, lhExecDescription)
						data2 = appendData(lhSSHAction, lhSSHDescription)
						data2 = appendData(vpcLsAction, vpcLsDescription)
						data2 = appendData(vpcControlAction, vpcControlDescription)
						data2 = appendData(agentUninstallAction, agentUninstallDescription)
						data2 = appendData(consoleAction, consoleDescription)
					case strings.Contains(o[1], "ReadOnlyAccess"):
						data2 = appendData(cvmLsAction, cvmLsDescription)
						data2 = appendData(lhLsAction, lhLsDescription)
						data2 = appendData(vpcLsAction, vpcLsDescription)
					case strings.Contains(o[1], "QcloudCVMFullAccess"):
						data2 = appendData(cvmLsAction, cvmLsDescription)
						data2 = appendData(cvmExecAction, cvmExecDescription)
					case strings.Contains(o[1], "QcloudCVMReadOnlyAccess"):
						data2 = appendData(cvmLsAction, cvmLsDescription)
					case strings.Contains(o[1], "QcloudCamFullAccess"):
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
		} else if strings.Contains(err.Error(), "Code=AuthFailure.UnauthorizedOperation") {
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

func getCallerIdentity() (string, uint64) {
	request := sts.NewGetCallerIdentityRequest()
	response, err := STSClient().GetCallerIdentity(request)
	errutil.HandleErr(err)
	accountType := response.Response.Type
	accountId, _ := strconv.Atoi(*response.Response.UserId)
	return *accountType, uint64(accountId)
}

func listAllPoliciesForUser(accountId uint64) ([][]string, error) {
	data, err := listPoliciesForUser(accountId)
	if err != nil {
		return nil, err
	} else {
		groups := listGroupsForUser(accountId)
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

func listPoliciesForUser(accountId uint64) ([][]string, error) {
	request := cam.NewListAttachedUserPoliciesRequest()
	request.TargetUin = common.Uint64Ptr(accountId)
	response, err := CAMClient().ListAttachedUserPolicies(request)
	if err == nil {
		log.Debugf("成功获取到 %d 用户的权限信息 (Successfully obtained permission information for %d user)", accountId, accountId)
		var data [][]string
		for n, i := range response.Response.List {
			SN := strconv.Itoa(n + 1)
			data = append(data, []string{SN, *i.PolicyName, *i.Remark})
		}
		return data, err
	} else {
		return nil, err
	}
}

func listPoliciesForGroup(groupName uint64) [][]string {
	request := cam.NewListAttachedGroupPoliciesRequest()
	request.TargetGroupId = common.Uint64Ptr(groupName)
	response, err := CAMClient().ListAttachedGroupPolicies(request)
	errutil.HandleErr(err)
	var data [][]string
	for n, i := range response.Response.List {
		SN := strconv.Itoa(n + 1)
		data = append(data, []string{SN, *i.PolicyName, *i.Remark})
	}
	return data
}

func listGroupsForUser(accountId uint64) []uint64 {
	request := cam.NewListGroupsForUserRequest()
	request.SubUin = common.Uint64Ptr(accountId)
	response, err := CAMClient().ListGroupsForUser(request)
	errutil.HandleErr(err)
	var groups []uint64
	for _, g := range response.Response.GroupInfo {
		groups = append(groups, *g.GroupId)
	}
	return groups
}

func traversalPermissions() ([][]string, [][]string) {
	var obj1 [][]string
	var obj2 [][]string
	// 1. cf tencent cvm ls
	log.Debugln("正在判断 cvm ls 权限 (Determining the permission of cvm ls)")
	request1 := cvm.NewDescribeInstancesRequest()
	_, err1 := tencentcvm.CVMClient("ap-beijing").DescribeInstances(request1)
	if err1 == nil {
		obj1 = append(obj1, []string{"QcloudCVMReadOnlyAccess", "云服务器（CVM）相关资源只读访问权限"})
		obj2 = append(obj2, []string{cvmLsAction, cvmLsDescription})
	} else {
		log.Traceln(err1.Error())
	}

	// 2. cf tencent cvm exec && cf tencent lh exec
	log.Debugln("正在判断 cvm exec 权限 (Determining the permission of cvm exec)")
	request2 := tat.NewPreviewReplacedCommandContentRequest()
	_, err2 := tencentcvm.TATClient("ap-beijing").PreviewReplacedCommandContent(request2)
	if err2 == nil {
		obj1 = append(obj1, []string{"QcloudTATFullAccess", "腾讯云自动化助手（TAT）全读写访问权限"})
		obj2 = append(obj2, []string{cvmExecAction, cvmExecDescription})
		obj2 = append(obj2, []string{lhExecAction, lhExecDescription})
	} else {
		log.Traceln(err2.Error())
	}

	// 3. cf tencent lh ls
	log.Debugln("正在判断 lh ls 权限 (Determining the permission of lh ls)")
	request3 := lh.NewDescribeInstancesRequest()
	_, err3 := tencentlh.LHClient("ap-beijing").DescribeInstances(request3)
	if err3 == nil {
		obj1 = append(obj1, []string{"QcloudLighthouseReadOnlyAccess", "轻量应用服务器（Lighthouse）只读访问权限"})
		obj2 = append(obj2, []string{lhLsAction, lhLsDescription})
	} else {
		log.Traceln(err3.Error())
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
