package huaweiiam

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	ecsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	rdsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	log "github.com/sirupsen/logrus"
	funk "github.com/thoas/go-funk"
)

var header = []string{"序号 (SN)", "策略名称 (PolicyName)", "描述 (Description)"}
var header2 = []string{"序号 (SN)", "可执行的操作 (Available actions)", "描述 (Description)"}
var SN = 1
var data2 [][]string
var rolesId []string

type BucketOperationsSample struct {
	bucketName string
	location   string
	obsClient  *obs.ObsClient
}

const (
	obslsAction        = "cf huawei obs ls"
	obslsDescription   = "列出 OBS 资源"
	obsgetAction       = "cf huawei obs obj get"
	obsgetDescription  = "下载 OBS 资源"
	ecslsAction        = "cf huawei ecs ls"
	ecslsDescription   = "列出 ECS 资源"
	rdslsAction        = "cf huawei rds ls"
	rdslsDescription   = "列出 RDS 资源"
	consoleAction      = "cf huawei console"
	consoleDescription = "接管控制台"
)

func ListPermissions() {
	var data [][]string
	userName, userId := getCallerIdentity()
	log.Infof("当前用户名为 %s (Current username is %s)", userName, userName)
	groupsName := getUserGroup(userId)
	_, domainsName := getDomainId()
	log.Infof("当前租户名/原华为云账户为 %s (Huawei Cloud account username is %s)", domainsName[0], domainsName[0])
	if funk.Contains(groupsName, "admin") {
		data = append(data, []string{"1", "All Administrator", "全部云服务管理员"})
		var td = cloud.TableData{Header: header, Body: data}
		Caption := "当前凭证具备的权限 (Permissions owned)"
		cloud.PrintTable(td, Caption)
		fmt.Println()
		data2 = appendData(obslsAction, obslsDescription)
		data2 = appendData(obsgetAction, obsgetDescription)
		data2 = appendData(ecslsAction, ecslsDescription)
		data2 = appendData(rdslsAction, rdslsDescription)
		data2 = appendData(consoleAction, consoleDescription)
		var td2 = cloud.TableData{Header: header2, Body: data2}
		Caption2 := "当前凭证可以执行的操作 (Available actions)"
		cloud.PrintTable(td2, Caption2)
	} else {
		// 如果不是admin组用户，就需要查找当前用户的所有权限
		allowPermissions := listAllPermissionsAction()
		if len(allowPermissions) != 0 {
			var td = cloud.TableData{Header: header, Body: data}
			Caption := "当前凭证具备的权限 (Permissions owned)"
			cloud.PrintTable(td, Caption)
			fmt.Println()
			if funk.Contains(allowPermissions, "*:*:*") {
				data2 = appendData(obslsAction, obslsDescription)
				data2 = appendData(obsgetAction, obsgetDescription)
				data2 = appendData(ecslsAction, ecslsDescription)
				data2 = appendData(rdslsAction, rdslsDescription)
				data2 = appendData(consoleAction, consoleDescription)
			} else if funk.Contains(allowPermissions, "obs:*:*") {
				data2 = appendData(obslsAction, obslsDescription)
				data2 = appendData(obsgetAction, obsgetDescription)
			} else if funk.Contains(allowPermissions, "ecs:*:*") || funk.Contains(allowPermissions, "ecs:*:get*") || funk.Contains(allowPermissions, "ecs:*:list*") {
				data2 = appendData(ecslsAction, ecslsDescription)
			} else if funk.Contains(allowPermissions, "rds:*:*") || funk.Contains(allowPermissions, "rds:*:get*") || funk.Contains(allowPermissions, "rds:*:list*") {
				data2 = appendData(rdslsAction, rdslsDescription)
			} else if funk.Contains(allowPermissions, "iam:users:*") {
				data2 = appendData(consoleAction, consoleDescription)
			}
			if len(data2) == 0 {
				log.Infoln("当前凭证没有可以执行的 OBS,ECS,RDS,IAM 操作 (Not available OBS,ECS,RDS,IAM actions)")
			} else {
				var td2 = cloud.TableData{Header: header2, Body: data2}
				Caption2 := "当前凭证可以执行的操作 (Available actions)"
				cloud.PrintTable(td2, Caption2)
			}
		} else if len(allowPermissions) == 0 {
			log.Debugln("当前凭证不具备 IAM 读权限 (No IAM read permissions)")
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
		}

	}
}

func getCallerIdentity() (string, string) {
	// 查找当前用户的UserId
	showPermanentAccessKeyRequestContent := &iamModel.ShowPermanentAccessKeyRequest{}
	showPermanentAccessKeyRequestContent.AccessKey = cmdutil.GetConfig("huawei").AccessKeyId
	showPermanentAccessKeyRequestResponse, err := IAMClient().ShowPermanentAccessKey(showPermanentAccessKeyRequestContent)
	if err != nil {
		log.Traceln(err)
	}
	userId := showPermanentAccessKeyRequestResponse.Credential.UserId
	// 查找当前用户的UserName
	showUserRequestContent := &iamModel.ShowUserRequest{}
	showUserRequestContent.UserId = userId
	showUserRequestResponse, err := IAMClient().ShowUser(showUserRequestContent)
	if err != nil {
		log.Traceln(err)
	}
	userName := showUserRequestResponse.User.Name
	log.Debugf("获得到当前凭证的用户名为 %s (The user name to get the current credentials is %s)", userName, userName)
	return userName, userId
}

func listAllPermissionsAction() []string {
	allowPermissionsRow := []string{}
	for _, roleId := range rolesId {
		keystoneShowPermissionRequestContent := &iamModel.KeystoneShowPermissionRequest{}
		keystoneShowPermissionRequestContent.RoleId = roleId
		keystoneShowPermissionRequestResponse, err := IAMClient().KeystoneShowPermission(keystoneShowPermissionRequestContent)
		if err != nil {
			log.Traceln(err)
		}
		for _, action := range *&keystoneShowPermissionRequestResponse.Role.Policy.Statement {
			if action.Effect.Value() == "Allow" {
				allowPermissionsRow = append(allowPermissionsRow, action.Action...)

			}
		}
	}
	allowPermissions := funk.UniqString(allowPermissionsRow)
	return allowPermissions
}

func ListDomainPermissionsForGroup(userId string, groupsId, domainsId []string) ([][]string, error) {
	data := [][]string{}
	for _, groupId := range groupsId {
		// 查找当前用户所在的用户组的权限
		for _, domainId := range domainsId {
			keystoneListDomainPermissionsForGroupRequestContent := &iamModel.KeystoneListDomainPermissionsForGroupRequest{}
			keystoneListDomainPermissionsForGroupRequestContent.DomainId = domainId
			keystoneListDomainPermissionsForGroupRequestContent.GroupId = groupId
			keystoneListDomainPermissionsForGroupRequestResponse, err := IAMClient().KeystoneListDomainPermissionsForGroup(keystoneListDomainPermissionsForGroupRequestContent)
			if err != nil {
				log.Traceln(err)
			}
			roles := keystoneListDomainPermissionsForGroupRequestResponse.Roles
			for sn, role := range *roles {
				rolesId = append(rolesId, role.Id)
				SN := strconv.Itoa(sn + 1)
				// 如果是自定义策略，没有DescriptionCn字段
				if role.DescriptionCn == nil {
					data = append(data, []string{SN, *role.DisplayName, *role.Description})
					log.Debugf("获得到当前用户所在的用户组全局的权限为 %s (The domain permission to get the current user's group is %s)", *role.DisplayName, *role.DisplayName)
				} else {
					data = append(data, []string{SN, *role.DisplayName, *role.DescriptionCn})
					log.Debugf("获得到当前用户所在的用户组全局的权限为 %s (The domain permission to get the current user's group is %s)", *role.DisplayName, *role.DisplayName)
				}
			}
		}
	}
	return data, nil
}

func ListAllProjectPermissionsForGroup(userId string, groupsId, domainsId []string) ([][]string, error) {
	data := [][]string{}
	for _, groupId := range groupsId {
		// 查找当前用户所在的用户组的权限
		for _, domainId := range domainsId {
			keystoneListAllProjectPermissionsForGroupRequestContent := &iamModel.KeystoneListAllProjectPermissionsForGroupRequest{}
			keystoneListAllProjectPermissionsForGroupRequestContent.DomainId = domainId
			keystoneListAllProjectPermissionsForGroupRequestContent.GroupId = groupId
			keystoneListAllProjectPermissionsForGroupRequestResponse, err := IAMClient().KeystoneListAllProjectPermissionsForGroup(keystoneListAllProjectPermissionsForGroupRequestContent)
			if err != nil {
				log.Traceln(err)
			}
			roles := keystoneListAllProjectPermissionsForGroupRequestResponse.Roles
			for sn, role := range *roles {
				rolesId = append(rolesId, role.Id)
				// 如果是自定义策略，没有DescriptionCn字段
				if role.DescriptionCn == nil {
					SN := strconv.Itoa(sn + 1)
					data = append(data, []string{SN, *role.DisplayName, *role.Description})
					log.Debugf("获得到当前用户所在的用户组所有资源的权限为 %s (The all project permission to get the current user's group is %s)", *role.DisplayName, *role.DisplayName)
				} else {
					SN := strconv.Itoa(sn + 1)
					data = append(data, []string{SN, *role.DisplayName, *role.DescriptionCn})
					log.Debugf("获得到当前用户所在的用户组所有资源的权限为 %s (The all project permission to get the current user's group is %s)", *role.DisplayName, *role.DisplayName)

				}
			}
		}
	}
	return data, nil
}

func traversalPermissions() ([][]string, [][]string) {
	// 遍历权限
	var obj1 [][]string
	var obj2 [][]string
	domainId, _ := getDomainId()
	huaweiConfig := cmdutil.GetConfig("huawei")
	region := "cn-east-3"
	log.Debugln("正在判断 obs get 权限 (Determining the permission of obs get)")
	// 判断是否有obs obj get
	var (
		endpoint   = "https://obs." + region + ".myhuaweicloud.com"
		ak         = huaweiConfig.AccessKeyId
		sk         = huaweiConfig.AccessKeySecret
		bucketName = util.GetRandomString(10)
		location   = region
	)
	sample := newBucketOperationsSample(ak, sk, endpoint, bucketName, location)
	if sample.CreateBucket() == true {
		log.Debugln("当前用户有 obs 管理员权限 (The current user has obs get permission)")
		obj1 = append(obj1, []string{"OBS Administrator", "对象存储服务管理员"})
		obj2 = append(obj2, []string{obslsAction, obslsDescription})
		obj2 = append(obj2, []string{obsgetAction, obsgetDescription})
		sample.DeleteBucket()
	} else {
		log.Debugln("正在判断 obs ls 权限 (Determining the permission of obs ls)")
		if sample.ListBuckets() == true {
			obj1 = append(obj1, []string{"OBS OperateAccess", "具有对象存储服务(OBS)查看桶列表等对象基本操作权限"})
			obj2 = append(obj2, []string{obslsAction, obslsDescription})
		}
	}
	log.Debugln("正在判断 ecs ls 权限 (Determining the permission of ecs ls)")
	listServersDetailsRequestContent := &ecsModel.ListServersDetailsRequest{}
	listServersDetailsRequestResponse, err := ECSClient().ListServersDetails(listServersDetailsRequestContent)
	if err == nil {
		log.Debugln(listServersDetailsRequestResponse)
		obj1 = append(obj1, []string{"ECS ReadOnlyAccess", "弹性云服务器的只读访问权限"})
		obj2 = append(obj2, []string{ecslsAction, ecslsDescription})
	} else {
		log.Traceln(err)
	}
	log.Debugln("正在判断 rds ls 权限 (Determining the permission of rds ls)")
	listInstancesRequestContent := &rdsModel.ListInstancesRequest{}
	listInstancesRequestResponse, err := RDSClient().ListInstances(listInstancesRequestContent)
	if err == nil {
		log.Debugln(listInstancesRequestResponse)
		obj1 = append(obj1, []string{"RDS ReadOnlyAccess", "关系型数据库服务资源只读权限"})
		obj2 = append(obj2, []string{rdslsAction, rdslsDescription})
	} else {
		log.Traceln(err)
	}

	log.Debugln("正在判断 console 权限 (Determining the permission of console)")
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	// 从切片中随机选择一个元素
	createUserRequestContent := &iamModel.CreateUserRequest{}
	passwordUser := util.GenerateRandomPasswords()
	userbody := &iamModel.CreateUserOption{
		Name:     "ggA2K2e4yxqN",
		DomainId: domainId[0],
		Password: &passwordUser,
	}
	createUserRequestContent.Body = &iamModel.CreateUserRequestBody{
		User: userbody,
	}
	createUserRequestResponse, err := IAMClient().CreateUser(createUserRequestContent)
	if err == nil {
		log.Debugln(createUserRequestResponse)
		obj1 = append(obj1, []string{"Security Administrator", "统一身份认证服务(除切换角色外)所有权限"})
		obj2 = append(obj2, []string{consoleAction, consoleDescription})
		newUserId := createUserRequestResponse.User.Id
		keystoneDeleteUserRequestContent := &iamModel.KeystoneDeleteUserRequest{}
		keystoneDeleteUserRequestContent.UserId = newUserId
		keystoneDeleteUserRequestResponse, err := IAMClient().KeystoneDeleteUser(keystoneDeleteUserRequestContent)
		if err == nil {
			log.Debugln(keystoneDeleteUserRequestResponse)
		} else {
			log.Traceln(err)
		}
	} else {
		log.Traceln(err)
	}
	return obj1, obj2
}

func getUserGroup(userId string) []string {
	// 查找当前用户所在的用户组的ID
	keystoneListGroupsForUserRequestContent := &iamModel.KeystoneListGroupsForUserRequest{}
	keystoneListGroupsForUserRequestContent.UserId = userId
	keystoneListGroupsForUserRequestResponse, err := IAMClient().KeystoneListGroupsForUser(keystoneListGroupsForUserRequestContent)
	if err != nil {
		log.Traceln(err)
	}

	groups := keystoneListGroupsForUserRequestResponse.Groups
	groupsName := []string{}
	for _, group := range *groups {
		groupsName = append(groupsName, group.Name)
		log.Debugf("获得到当前用户所在的用户组为 %s (The user group to get the current user is %s)", group.Name, group.Name)
	}
	return groupsName
}

func getDomainId() ([]string, []string) {
	keystoneListAuthDomainsRequestContent := &iamModel.KeystoneListAuthDomainsRequest{}
	keystoneListAuthDomainsRequestResponse, err := IAMClient().KeystoneListAuthDomains(keystoneListAuthDomainsRequestContent)
	if err != nil {
		log.Traceln(err)
	}
	domains := keystoneListAuthDomainsRequestResponse.Domains
	domainsId := []string{}
	domainsName := []string{}
	for _, domain := range *domains {
		domainsId = append(domainsId, domain.Id)
		domainsName = append(domainsName, domain.Name)
		log.Debugf("查询IAM用户可以访问到的账户(Root Account)为 %s  (The domain to get the current user is %s)", domain.Name, domain.Name)
	}
	return domainsId, domainsName
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

func newBucketOperationsSample(ak, sk, endpoint, bucketName, location string) *BucketOperationsSample {
	obsClient, err := obs.New(ak, sk, endpoint)
	if err != nil {
		log.Traceln(err)
	}
	return &BucketOperationsSample{obsClient: obsClient, bucketName: bucketName, location: location}
}

func (sample BucketOperationsSample) CreateBucket() bool {
	input := &obs.CreateBucketInput{}
	input.Bucket = sample.bucketName
	input.Location = sample.location
	_, err := sample.obsClient.CreateBucket(input)
	if err != nil {
		log.Traceln(err)
		return false
	}
	log.Debugf("Create bucket:%s successfully!", sample.bucketName)
	return true
}

func (sample BucketOperationsSample) DeleteBucket() bool {
	_, err := sample.obsClient.DeleteBucket(sample.bucketName)
	if err != nil {
		log.Traceln(err)
		return false
	}
	log.Debugf("Delete bucket %s successfully!", sample.bucketName)
	return true
}

func (sample BucketOperationsSample) ListBuckets() bool {
	input := &obs.CreateSignedUrlInput{}
	input.Method = obs.HttpMethodGet
	input.Expires = 3600
	output, err := sample.obsClient.CreateSignedUrl(input)
	if err == nil {
		log.Debugf("using temporary signature url: %s", output.SignedUrl)
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			log.Debugf("Code:" + obsError.Code)
			log.Debugf("Message:" + obsError.Message)
		} else {
			log.Traceln(err)
		}
	}

	listBucketsOutput, err := sample.obsClient.ListBucketsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err == nil {
		log.Debugf("Owner.DisplayName:%s, Owner.ID:%s", listBucketsOutput.Owner.DisplayName, listBucketsOutput.Owner.ID)
		return true
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			log.Debugf("Code:" + obsError.Code)
			log.Debugf("Message:" + obsError.Message)
		} else {
			log.Traceln(err)
		}
		return false
	}

}
